package loggerx

import (
	"context"
	"reflect"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey"
	"github.com/olivere/elastic/v7"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/logger"
	"github.com/hatlonely/go-kit/wrap"
)

func TestElasticSearchWriter_Write(t *testing.T) {
	patches := ApplyFunc(elastic.NewClient, func(options ...elastic.ClientOptionFunc) (*elastic.Client, error) {
		return &elastic.Client{}, nil
	}).ApplyMethod(reflect.TypeOf(&elastic.IndexService{}), "Do", func(
		s *elastic.IndexService, ctx context.Context) (*elastic.IndexResponse, error) {
		return nil, nil
	}).ApplyMethod(reflect.TypeOf(&elastic.PingService{}), "Do", func(
		s *elastic.PingService, ctx context.Context) (*elastic.PingResult, int, error) {
		return nil, 0, nil
	})

	defer patches.Reset()

	ts, _ := time.Parse(time.RFC3339Nano, "2020-11-07T16:07:36.65798+08:00")
	info := &logger.Info{
		Time:   ts,
		Level:  logger.LevelInfo,
		File:   "formatter_test.go:18",
		Caller: "TextFormatter",
		Fields: map[string]interface{}{
			"meta1": "val1",
			"meta2": "val2",
		},
		Data: map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
		},
	}
	Convey("TestElasticSearchWriter_Write", t, func() {
		writer, err := NewElasticSearchWriterWithOptions(&ElasticSearchWriterOptions{
			Level:      "Info",
			Index:      "hatlonely",
			IDField:    "requestID",
			Timeout:    200 * time.Millisecond,
			MsgChanLen: 200,
			WorkerNum:  2,
			ESClientWrapper: wrap.ESClientWrapperOptions{
				ES: wrap.ESOptions{
					URI: "http://127.0.0.1:9200",
				},
			},
		})

		So(err, ShouldBeNil)
		So(writer.Write(info), ShouldBeNil)
		So(writer.Close(), ShouldBeNil)
	})
}
