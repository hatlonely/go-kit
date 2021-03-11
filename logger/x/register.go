package loggerx

import (
	"github.com/hatlonely/go-kit/logger"
)

func init() {
	logger.RegisterFormatter("Markdown", NewMarkdownFormatterWithOptions)
	logger.RegisterFormatter("Html", NewHtmlFormatterWithOptions)

	logger.RegisterWriter("ElasticSearch", NewElasticSearchWriterWithOptions)
	logger.RegisterWriter("DingTalk", NewDingTalkWriterWithOptions)
	logger.RegisterWriter("Email", NewEmailWriterWithOptions)
	logger.RegisterWriter("MNS", NewMNSWriterWithOptions)
}
