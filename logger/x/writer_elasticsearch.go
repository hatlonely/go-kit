package loggerx

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/hatlonely/go-kit/cast"
	"github.com/hatlonely/go-kit/logger"
	"github.com/hatlonely/go-kit/strx"
	"github.com/hatlonely/go-kit/wrap"
)

type ElasticSearchWriterOptions struct {
	ES         wrap.ESClientWrapperOptions
	Level      string `dft:"Debug"`
	Index      string
	IDField    string
	Timeout    time.Duration `dft:"200ms"`
	MsgChanLen int           `dft:"200"`
	WorkerNum  int           `dft:"1"`
}

func NewElasticSearchWriterWithOptions(options *ElasticSearchWriterOptions) (*ElasticSearchWriter, error) {
	level, err := logger.LevelString(options.Level)
	if err != nil {
		return nil, errors.Wrap(err, "LevelToString failed")
	}

	client, err := wrap.NewESClientWrapperWithOptions(&options.ES)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()
	if _, _, err := client.Ping(options.ES.ES.URI).Do(ctx); err != nil {
		return nil, err
	}

	w := &ElasticSearchWriter{
		level:    level,
		esCli:    client,
		messages: make(chan *logger.Info, options.MsgChanLen),
		options:  options,
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

type ElasticSearchWriter struct {
	options *ElasticSearchWriterOptions
	level   logger.Level

	esCli    *wrap.ESClientWrapper
	messages chan *logger.Info
	wg       sync.WaitGroup
}

func (w *ElasticSearchWriter) Write(info *logger.Info) error {
	if info.Level < w.level {
		return nil
	}

	w.messages <- info
	return nil
}

func (w *ElasticSearchWriter) Close() error {
	close(w.messages)
	w.wg.Wait()

	return nil
}

func (w *ElasticSearchWriter) work() {
	for info := range w.messages {
		if _, ok := info.Fields[w.options.IDField]; !ok {
			info.Fields[w.options.IDField] = uuid.NewV4().String()
		}
		ctx, cancel := context.WithTimeout(context.Background(), w.options.Timeout)
		if _, err := w.esCli.Index().Index(w.options.Index).Id(cast.ToString(info.Fields[w.options.IDField])).BodyJson(info).Do(ctx); err != nil {
			fmt.Printf("ElasticSearchWriter write log failed. err: [%+v], kvs: [%v]\n", err, strx.JsonMarshal(info))
		}
		cancel()
	}
}
