package loggerx

import (
	"context"
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/hatlonely/go-kit/cast"
	"github.com/hatlonely/go-kit/wrap"
)

type ElasticSearchWriter struct {
	esCli    *wrap.ESClientWrapper
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

	ESClientWrapper wrap.ESClientWrapperOptions
}

func NewElasticSearchWriterWithOptions(options *ElasticSearchWriterOptions) (*ElasticSearchWriter, error) {
	client, err := wrap.NewESClientWrapperWithOptions(&options.ESClientWrapper)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()
	if _, _, err := client.Ping(options.ESClientWrapper.ES.URI).Do(ctx); err != nil {
		return nil, err
	}

	w := &ElasticSearchWriter{
		esCli:    client,
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
		ctx, cancel := context.WithTimeout(context.Background(), w.options.Timeout)
		if _, err := w.esCli.Index().Index(w.options.Index).Id(cast.ToString(kvs[w.options.IDField])).BodyJson(kvs).Do(ctx); err != nil {
			fmt.Printf("ElasticSearchWriter write log failed, err: [%v]\n", err)
		}
		cancel()
	}
}