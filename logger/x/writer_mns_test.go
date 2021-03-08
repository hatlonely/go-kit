package loggerx

import (
	"reflect"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey"
	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/logger"
	"github.com/hatlonely/go-kit/wrap"
)

func TestMNSWriter_Write(t *testing.T) {
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

	Convey("TestEmailWriter_Write", t, func(c C) {
		patches := ApplyMethod(reflect.TypeOf(&ali_mns.MNSTopic{}), "PublishMessage", func(
			topic *ali_mns.MNSTopic, message ali_mns.MessagePublishRequest) (resp ali_mns.MessageSendResponse, err error) {
			return ali_mns.MessageSendResponse{}, nil
		})
		defer patches.Reset()

		w, err := NewMNSWriterWithOptions(&MNSWriterOptions{
			Level:      "Info",
			WorkerNum:  1,
			MsgChanLen: 200,
			Formatter: logger.FormatterOptions{
				Type: "Html",
				Options: &HtmlFormatterOptions{
					DefaultTitle: "测试告警",
				},
			},
			MNS: wrap.MNSClientWrapperOptions{
				MNS: wrap.MNSOptions{
					Endpoint:        "http://xx.mns.cn-shanghai.aliyuncs.com/",
					AccessKeyID:     "xx",
					AccessKeySecret: "xx",
				},
			},
			Topic:   "xx",
			Timeout: 3 * time.Second,
		})
		So(err, ShouldBeNil)

		So(w.Write(info), ShouldBeNil)
		So(w.Close(), ShouldBeNil)
	})
}
