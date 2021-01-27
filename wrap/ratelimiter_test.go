package wrap

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/time/rate"

	"github.com/hatlonely/go-kit/refx"
)

func Test1(t *testing.T) {
	Convey("Test", t, func() {
		limiter := rate.NewLimiter(rate.Every(200*time.Millisecond), 100)
		limiter.Allow()
	})

	Convey("Test2", t, func() {
		v := map[string]interface{}{
			"type": "LocalGroup",
			"localGroupRateLimiter": map[string]interface{}{
				"DB.First": map[string]interface{}{
					"interval": "2s",
					"burst":    1,
				},
			},
		}

		var options RateLimiterGroupOptions

		So(refx.InterfaceToStruct(v, &options, refx.WithCamelName()), ShouldBeNil)

		rateLimiterGroup, err := NewRateLimiterGroup(&options)
		So(err, ShouldBeNil)

		for i := 0; i < 10; i++ {
			rateLimiterGroup.Wait(context.Background(), "DB.First")
			fmt.Println("hello world")
		}
		fmt.Println(options)
	})
}
