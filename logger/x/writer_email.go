package loggerx

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/logger"
	"github.com/hatlonely/go-kit/strx"
)

type EmailWriterOptions struct {
	Level string `dft:"Debug"`

	Server   string
	Port     int
	Password string
	Username string

	ToList  []string
	Subject string

	MsgChanLen int `dft:"200"`
	WorkerNum  int `dft:"1"`

	Formatter logger.FormatterOptions
}

func NewEmailWriterWithOptions(options *EmailWriterOptions) (*EmailWriter, error) {
	level, err := logger.LevelString(options.Level)
	if err != nil {
		return nil, errors.Wrap(err, "LevelToString failed")
	}

	formatter, err := logger.NewFormatterWithOptions(&options.Formatter)
	if err != nil {
		return nil, errors.WithMessage(err, "logger.NewFormatterWithOptions failed")
	}

	w := &EmailWriter{
		options:   options,
		level:     level,
		messages:  make(chan *logger.Info, options.MsgChanLen),
		formatter: formatter,
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

type EmailWriter struct {
	options   *EmailWriterOptions
	level     logger.Level
	formatter logger.Formatter

	messages chan *logger.Info
	wg       sync.WaitGroup
}

func (w *EmailWriter) Write(info *logger.Info) error {
	if info.Level < w.level {
		return nil
	}

	w.messages <- info
	return nil
}

func (w *EmailWriter) Close() error {
	close(w.messages)
	w.wg.Wait()

	return nil
}

func (w *EmailWriter) work() {
	for info := range w.messages {
		buf, err := w.formatter.Format(info)
		if err != nil {
			buf, _ = json.MarshalIndent(info, "  ", "  ")
		}

		if err := w.send(w.options.ToList, w.options.Subject, string(buf)); err != nil {
			fmt.Printf("EmailWriter write log failed. err: [%+v], kvs: [%v]\n", err, strx.JsonMarshal(info))
		}
	}
}

func (w *EmailWriter) send(toList []string, subject string, body string) error {
	content := fmt.Sprintf(`From: %v
To: %v
Subject: %v
Content-Type: text/html; charset=UTF-8;
%v
`, w.options.Username, strings.Join(toList, ";"), subject, body)

	return smtp.SendMail(
		fmt.Sprintf("%v:%v", w.options.Server, w.options.Port),
		smtp.PlainAuth("", w.options.Username, w.options.Password, w.options.Server),
		w.options.Username,
		toList,
		[]byte(content),
	)
}
