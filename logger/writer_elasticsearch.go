package logger

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/olivere/elastic/v7"
	uuid "github.com/satori/go.uuid"

	"github.com/hatlonely/go-kit/cast"
)

type ElasticSearchWriter struct {
	esCli    *elastic.Client
	messages chan map[string]interface{}
	wg       sync.WaitGroup

	options *ElasticSearchWriterOptions
}

type ElasticSearchOptions struct {
	URI                       string `dft:"http://elasticsearch:9200"`
	EnableSniff               bool
	Username                  string
	Password                  string
	EnableHealthCheck         bool          `dft:"true"`
	HealthCheckInterval       time.Duration `dft:"60s"`
	HealthCheckTimeout        time.Duration `dft:"5s"`
	HealthCheckTimeoutStartUp time.Duration `dft:"5s"`
}

type ElasticSearchWriterOptions struct {
	Index      string
	IDField    string
	Timeout    time.Duration `dft:"200ms"`
	MsgChanLen int           `dft:"200"`
	WorkerNum  int           `dft:"1"`

	ES ElasticSearchOptions
}

func NewElasticSearchWriterWithOptions(options *ElasticSearchWriterOptions) (*ElasticSearchWriter, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(options.ES.URI),
		elastic.SetSniff(options.ES.EnableSniff),
		elastic.SetBasicAuth(options.ES.Username, options.ES.Password),
		elastic.SetHealthcheck(options.ES.EnableHealthCheck),
		elastic.SetHealthcheckInterval(options.ES.HealthCheckInterval),
		elastic.SetHealthcheckTimeout(options.ES.HealthCheckTimeout),
		elastic.SetHealthcheckTimeoutStartup(options.ES.HealthCheckTimeoutStartUp),
	)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()
	if _, _, err := client.Ping(options.ES.URI).Do(ctx); err != nil {
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
