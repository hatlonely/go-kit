package loggerx

import (
	"github.com/hatlonely/go-kit/logger"
)

func init() {
	logger.RegisterWriter("ElasticSearch", NewElasticSearchWriterWithOptions)
}
