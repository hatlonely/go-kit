package logger

import (
	"context"
	"reflect"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey"
	"github.com/olivere/elastic/v7"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/cli"
)

func TestElasticSearchWriter_Write(t *testing.T) {
	patches := ApplyFunc(cli.NewElasticSearchWithOptions, func(options *cli.ElasticSearchOptions) (*elastic.Client, error) {
		return &elastic.Client{}, nil
	}).ApplyMethod(reflect.TypeOf(&elastic.IndexService{}), "Do", func(
		s *elastic.IndexService, ctx context.Context) (*elastic.IndexResponse, error) {
		return nil, nil
	})

	defer patches.Reset()

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
