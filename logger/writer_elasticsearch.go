package logger

import (
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/cli"
)

type ElasticSearchWriter struct {
	esCli *elastic.Client
}

type ElasticSearchWriterOptions struct {
	ElasticSearch cli.ElasticSearchOptions
}

func NewElasticSearchWriterWithOptions(options *ElasticSearchWriterOptions) (*ElasticSearchWriter, error) {
	esCli, err := cli.NewElasticSearchWithOptions(&options.ElasticSearch)
	if err != nil {
		return nil, errors.Wrap(err, "cli.NewElasticSearchWithOptions failed")
	}
	return &ElasticSearchWriter{esCli: esCli}, nil
}
