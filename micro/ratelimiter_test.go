package micro

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/refx"
)

func TestNewRateLimiterWithOptions(t *testing.T) {
	Convey("TestNewRateLimiterWithOptions", t, func() {
		Convey("empty options case", func() {
			r, err := NewRateLimiterWithOptions(&RateLimiterOptions{})
			So(err, ShouldBeNil)
			So(r, ShouldBeNil)
		})

		Convey("local group rate limiter", func() {
			r, err := NewRateLimiterWithOptions(&RateLimiterOptions{
				Type: "LocalGroup",
				Options: &LocalGroupRateLimiterOptions{
					"key1": {
						Interval: time.Second,
						Burst:    2,
					},
				},
			})
			So(err, ShouldBeNil)
			So(r, ShouldHaveSameTypeAs, &LocalGroupRateLimiter{})
		})

		Convey("local group rate limiter form config", func() {
			r, err := NewRateLimiterWithOptions(&RateLimiterOptions{
				Type: "LocalGroup",
				Options: map[string]interface{}{
					"key1": map[string]interface{}{
						"interval": "2s",
						"burst":    2,
					},
				},
			}, refx.WithCamelName())
			So(err, ShouldBeNil)
			So(r, ShouldHaveSameTypeAs, &LocalGroupRateLimiter{})
			So(r.(*LocalGroupRateLimiter).limiters["key1"].Burst(), ShouldEqual, 2)
		})
	})
}
