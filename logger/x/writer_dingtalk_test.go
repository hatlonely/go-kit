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

func TestDingTalkWriter_Send(t *testing.T) {
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
			AccessToken:         "test-token",
			Secret:              "test-secret",
			Title:               "测试",
			DialTimeout:         3 * time.Second,
			Timeout:             6 * time.Second,
			MaxIdleConnsPerHost: 2,
			MsgChanLen:          200,
			Formatter: logger.FormatterOptions{
				Type: "Json",
			},
		})
		So(err, ShouldBeNil)

		So(w.send(map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
		}), ShouldBeNil)
	})
}
