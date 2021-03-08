package loggerx

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/logger"
	"github.com/hatlonely/go-kit/strx"
	"github.com/hatlonely/go-kit/wrap"
)

type MNSWriterOptions struct {
	MNS        wrap.MNSClientWrapperOptions
	Level      string `dft:"Debug"`
	MsgChanLen int    `dft:"200"`
	WorkerNum  int    `dft:"1"`
	Formatter  logger.FormatterOptions
	Topic      string
	Timeout    time.Duration `dft:"3s"`
}

func NewMNSWriterWithOptions(options *MNSWriterOptions) (*MNSWriter, error) {
	level, err := logger.LevelString(options.Level)
	if err != nil {
		return nil, errors.Wrap(err, "LevelToString failed")
	}

	formatter, err := logger.NewFormatterWithOptions(&options.Formatter)
	if err != nil {
		return nil, errors.WithMessage(err, "logger.NewFormatterWithOptions failed")
	}

	client, err := wrap.NewMNSClientWrapperWithOptions(&options.MNS)
	if err != nil {
		return nil, errors.WithMessage(err, "wrap.NewMNSClientWrapperWithOptions failed")
	}

	topic := wrap.NewMNSTopic(options.Topic, client)

	w := &MNSWriter{
		options:   options,
		level:     level,
		formatter: formatter,
		topic:     topic,
		messages:  make(chan *logger.Info, options.MsgChanLen),
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

type MNSWriter struct {
	options   *MNSWriterOptions
	level     logger.Level
	formatter logger.Formatter

	messages chan *logger.Info
	wg       sync.WaitGroup

	topic *wrap.MNSTopicWrapper
}

func (w *MNSWriter) Write(info *logger.Info) error {
	if info.Level < w.level {
		return nil
	}

	w.messages <- info
	return nil
}

func (w *MNSWriter) Close() error {
	close(w.messages)
	w.wg.Wait()

	return nil
}

func (w *MNSWriter) work() {
	for info := range w.messages {
		if err := w.send(info); err != nil {
			fmt.Printf("MNSWriter write log failed. err: [%+v], kvs: [%v]\n", err, strx.JsonMarshal(info))
		}
	}
}

func (w *MNSWriter) send(info *logger.Info) error {
	buf, err := w.formatter.Format(info)
	if err != nil {
		buf, _ = json.MarshalIndent(info, "", "  ")
	}
	ctx, cancel := context.WithTimeout(context.Background(), w.options.Timeout)
	defer cancel()
	_, err = w.topic.PublishMessage(ctx, ali_mns.MessagePublishRequest{
		MessageBody: string(buf),
	})
	if err != nil {
		return errors.Wrap(err, "topic.PublishMessage failed")
	}
	return nil
}
