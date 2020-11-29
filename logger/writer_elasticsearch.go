package logger

import (
	"context"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/cast"
	"github.com/hatlonely/go-kit/cli"
)

type ElasticSearchWriter struct {
	esCli *elastic.Client
	vals  chan map[string]interface{}

	options *ElasticSearchWriterOptions
}

type ElasticSearchWriterOptions struct {
	Index   string
	IDField string
	Timeout time.Duration

	ElasticSearch cli.ElasticSearchOptions
}

func NewElasticSearchWriterWithOptions(options *ElasticSearchWriterOptions) (*ElasticSearchWriter, error) {
	esCli, err := cli.NewElasticSearchWithOptions(&options.ElasticSearch)
	if err != nil {
		return nil, errors.Wrap(err, "cli.NewElasticSearchWithOptions failed")
	}
	return &ElasticSearchWriter{esCli: esCli}, nil
}

func (w *ElasticSearchWriter) Write(kvs map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), w.options.Timeout)
	defer cancel()
	if _, err := w.esCli.Index().Index(w.options.Index).Id(cast.ToString(kvs[w.options.IDField])).BodyJson(kvs).Do(ctx); err != nil {
		return errors.WithMessagef(err, "elasticsearch.PutDocument failed")
	}

	return nil
}
