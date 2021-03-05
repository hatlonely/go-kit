package loggerx

import (
	"github.com/hatlonely/go-kit/logger"
)

func init() {
	logger.RegisterFormatter("Markdown", NewMarkdownFormatterWithOptions)

	logger.RegisterWriter("ElasticSearch", NewElasticSearchWriterWithOptions)
}
