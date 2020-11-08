package config

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOTSLegacyProvider(t *testing.T) {
	Convey("TestOTSLegacyProvider", t, func() {
		provider, err := NewOTSLegacyProviderWithOptions(&OTSLegacyProviderOptions{
			Endpoint:        "https://hatlonely.cn-shanghai.ots.aliyuncs.com",
			Instance:        "hatlonely",
			AccessKeyID:     "xx",
			AccessKeySecret: "xx",
			Table:           "TestConfig",
			PrimaryKeys:     []string{"Key1", "Key2"},
			Interval:        1 * time.Second,
		})
		So(err, ShouldBeNil)
		buf, _ := provider.Load()
		fmt.Println(string(buf))
	})
}
