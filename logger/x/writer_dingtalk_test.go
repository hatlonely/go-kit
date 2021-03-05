package loggerx

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/logger"
)

func TestCalculateDingTalkSign(t *testing.T) {
	Convey("TestCalculateDingTalkSign", t, func() {
		So(calculateSign(1614834744445, "test-secret-1"), ShouldEqual, "1NfE5h9MRMRTk%2BKVmADFChpFCfmSV14v4CGFdE22S7Q%3D")
		So(calculateSign(1614834744446, "test-secret-2"), ShouldEqual, "i0xUVZpu3t75uz9OglaWdXhn7l%2FCRGINycYHheIcPY8%3D")
	})
}

func TestDingTalkWriter_Write(t *testing.T) {
	Convey("", t, func() {
		patches := ApplyMethod(reflect.TypeOf(&http.Client{}), "Do", func(cli *http.Client, req *http.Request) (*http.Response, error) {
			return &http.Response{
				Status:     http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(`{"errcode":0,"errmsg":"ok"}`))),
			}, nil
		})
		defer patches.Reset()

		w, err := NewDingTalkWriterWithOptions(&DingTalkWriterOptions{
			AccessToken:         "xx",
			Secret:              "xx",
			Title:               "测试",
			DialTimeout:         3 * time.Second,
			Timeout:             6 * time.Second,
			MaxIdleConnsPerHost: 2,
			WorkerNum:           1,
			MsgChanLen:          200,
			Formatter: logger.FormatterOptions{
				Type: "Markdown",
				Options: &MarkdownFormatterOptions{
					DefaultTitle: "测试告警",
				},
			},
		})
		So(err, ShouldBeNil)

		now := time.Now()
		So(w.Write(map[string]interface{}{
			"timestamp": now.Unix(),
			"time":      now.Format(time.RFC3339Nano),
			"level":     logger.LevelError.String(),
			"file":      "test-file",
			"caller":    "test-caller",
			"data": map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
			},
		}), ShouldBeNil)

		So(w.Close(), ShouldBeNil)
	})
}

func TestLoggerWithDingTalkWriter(t *testing.T) {
	Convey("", t, func() {
		patches := ApplyMethod(reflect.TypeOf(&http.Client{}), "Do", func(cli *http.Client, req *http.Request) (*http.Response, error) {
			return &http.Response{
				Status:     http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(`{"errcode":0,"errmsg":"ok"}`))),
			}, nil
		})
		defer patches.Reset()

		w, err := NewDingTalkWriterWithOptions(&DingTalkWriterOptions{
			AccessToken:         "xx",
			Secret:              "xx",
			Title:               "测试",
			DialTimeout:         3 * time.Second,
			Timeout:             6 * time.Second,
			MaxIdleConnsPerHost: 2,
			WorkerNum:           1,
			MsgChanLen:          200,
			Formatter: logger.FormatterOptions{
				Type: "Markdown",
				Options: &MarkdownFormatterOptions{
					DefaultTitle: "测试告警",
				},
			},
		})
		So(err, ShouldBeNil)

		log := logger.NewLogger(logger.LevelInfo, w)

		log.With("MetaKey1", "Val1").With("MetaKey2", "Vale").Error(map[string]interface{}{
			"key1":  "val1",
			"hello": "world",
		})

		So(log.Close(), ShouldBeNil)
	})
}
