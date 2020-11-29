package logger

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/hatlonely/go-kit/cast"
	"github.com/hatlonely/go-kit/cli"
)

type ElasticSearchWriter struct {
	esCli    *elastic.Client
	messages chan map[string]interface{}
	wg       sync.WaitGroup

	options *ElasticSearchWriterOptions
}

type ElasticSearchWriterOptions struct {
	Index      string
	IDField    string
	Timeout    time.Duration `dft:"200ms"`
	MsgChanLen int           `dft:"200"`
	WorkerNum  int           `dft:"1"`

	ElasticSearch cli.ElasticSearchOptions
}

func NewElasticSearchWriterWithOptions(options *ElasticSearchWriterOptions) (*ElasticSearchWriter, error) {
	esCli, err := cli.NewElasticSearchWithOptions(&options.ElasticSearch)
	if err != nil {
		return nil, errors.Wrap(err, "cli.NewElasticSearchWithOptions failed")
	}

	w := &ElasticSearchWriter{
		esCli:    esCli,
		messages: make(chan map[string]interface{}, options.MsgChanLen),
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

func (w *ElasticSearchWriter) Write(kvs map[string]interface{}) error {
	w.messages <- kvs
	return nil
}

func (w *ElasticSearchWriter) Close() error {
	close(w.messages)
	w.wg.Wait()

	return nil
}

func (w *ElasticSearchWriter) work() {
	for kvs := range w.messages {
		if _, ok := kvs[w.options.IDField]; !ok {
			kvs[w.options.IDField] = uuid.NewV4().String()
		}
		err := retry.Do(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), w.options.Timeout)
			defer cancel()
			if _, err := w.esCli.Index().Index(w.options.Index).Id(cast.ToString(kvs[w.options.IDField])).BodyJson(kvs).Do(ctx); err != nil {
				return err
			}
			return nil
		}, retry.Attempts(6), retry.DelayType(retry.BackOffDelay), retry.LastErrorOnly(true))
		if err != nil {
			fmt.Printf("ElasticSearchWriter write log failed, err: [%v]\n", err)
		}
	}
}
