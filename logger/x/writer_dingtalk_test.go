package loggerx

import (
	"testing"
	"time"

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
		w, err := NewDingTalkWriterWithOptions(&DingTalkWriterOptions{
			AccessToken:         "fbade6a40c8110b14f499f736a3f257ead441a0e75e1cb5471679ae6d79aea0c",
			Secret:              "SEC01ae9d8839d281c45d0e1eb3dab7ef167431873b5a5b3ba002e3324dc2b39b53",
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
