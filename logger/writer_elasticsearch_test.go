package logger

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/cli"
)

func TestElasticSearchWriter_Write(t *testing.T) {
	Convey("", t, func() {
		writer, err := NewWriterWithOptions(&WriterOptions{
			Type: "ElasticSearch",
			ElasticSearchWriter: ElasticSearchWriterOptions{
				Index:      "hatlonely",
				IDField:    "requestID",
				Timeout:    200 * time.Millisecond,
				MsgChanLen: 200,
				WorkerNum:  2,
				ElasticSearch: cli.ElasticSearchOptions{
					URI: "http://127.0.0.1:9200",
				},
			},
		})

		So(err, ShouldBeNil)
		_ = writer.Write(map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
		})
		_ = writer.Write(map[string]interface{}{
			"key3": "val3",
			"key4": "val4",
		})

		defer writer.Close()
	})
}
