package config

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestACMProvider(t *testing.T) {
	Convey("TestACMProvider", t, func() {
		provider, err := NewACMProviderWithOptions(&ACMProviderOptions{
			Endpoint:        "test-endpoint",
			Namespace:       "test-namespace",
			AccessKeyID:     "xx",
			AccessKeySecret: "xx",
			Timeout:         5 * time.Second,
			DataID:          "test-data-id",
			Group:           "test-group",
		})
		So(err, ShouldBeNil)
		fmt.Println(provider.Load())
	})
}
