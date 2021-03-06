package loggerx

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/logger"
	"github.com/hatlonely/go-kit/strx"
)

type DingTalkWriterOptions struct {
	Level string

	AccessToken string
	Secret      string
	Title       string

	DialTimeout         time.Duration `dft:"3s"`
	Timeout             time.Duration `dft:"6s"`
	MaxIdleConnsPerHost int           `dft:"2"`
	MsgChanLen          int           `dft:"200"`
	WorkerNum           int           `dft:"1"`

	Formatter logger.FormatterOptions
}

func NewDingTalkWriterWithOptions(options *DingTalkWriterOptions) (*DingTalkWriter, error) {
	level, err := logger.LevelString(options.Level)
	if err != nil {
		return nil, errors.Wrap(err, "LevelToString failed")
	}

	formatter, err := logger.NewFormatterWithOptions(&options.Formatter)
	if err != nil {
		return nil, errors.WithMessage(err, "logger.NewFormatterWithOptions failed")
	}

	w := &DingTalkWriter{
		options:   options,
		level:     level,
		messages:  make(chan *logger.Info, options.MsgChanLen),
		formatter: formatter,
		cli: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, netw, addr string) (net.Conn, error) {
					c, err := net.DialTimeout(netw, addr, options.DialTimeout)
					if err != nil {
						return nil, err
					}
					return c, nil
				},
				MaxIdleConnsPerHost: options.MaxIdleConnsPerHost,
			},
			Timeout: options.Timeout,
		},
	}

	for i := 0; i < options.WorkerNum; i++ {
		w.wg.Add(1)
		go func() {
			w.work()
			w.wg.Done()
		}()
	}

	return w, nil
}

type DingTalkWriter struct {
	options   *DingTalkWriterOptions
	level     logger.Level
	formatter logger.Formatter
	cli       *http.Client

	messages chan *logger.Info
	wg       sync.WaitGroup
}

type DingTalkMessage struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
}

type DingTalkError struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (w *DingTalkWriter) Write(info *logger.Info) error {
	if info.Level < w.level {
		return nil
	}

	w.messages <- info
	return nil
}

func (w *DingTalkWriter) Close() error {
	close(w.messages)
	w.wg.Wait()

	return nil
}

func (w *DingTalkWriter) work() {
	for info := range w.messages {
		if err := w.send(info); err != nil {
			fmt.Printf("DingTalkWriter write log failed. err: [%+v], kvs: [%v]\n", err, strx.JsonMarshal(info))
		}
	}
}

func calculateSign(timestamp int64, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(fmt.Sprintf("%d\n%s", timestamp, secret)))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return url.QueryEscape(sign)
}

func (w *DingTalkWriter) send(info *logger.Info) error {
	buf, err := w.formatter.Format(info)
	if err != nil {
		buf, err = json.Marshal(info)
	}

	var message DingTalkMessage
	message.MsgType = "markdown"
	message.Markdown.Title = w.options.Title
	message.Markdown.Text = string(buf)
	body, _ := json.Marshal(&message)

	timestamp := time.Now().UnixNano() / 1000000
	req, err := http.NewRequest("POST", fmt.Sprintf(
		"https://oapi.dingtalk.com/robot/send?access_token=%s&timestamp=%d&sign=%s",
		w.options.AccessToken, timestamp, calculateSign(timestamp, w.options.Secret),
	), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := w.cli.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rbody, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return errors.Errorf("status not ok. status: [%s], body: [%s]", res.Status, string(rbody))
	}

	var dingTalkError DingTalkError
	if err := json.Unmarshal(rbody, &dingTalkError); err != nil {
		return errors.Wrapf(err, "json.Unmarshal failed. body: [%s]", string(rbody))
	}

	if dingTalkError.ErrCode != 0 {
		return errors.New(dingTalkError.ErrMsg)
	}

	return nil
}
