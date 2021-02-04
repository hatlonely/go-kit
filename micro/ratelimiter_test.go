package micro

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewRateLimiterWithOptions(t *testing.T) {
	Convey("TestNewRateLimiterWithOptions", t, func() {
		Convey("empty options case", func() {
			r, err := NewRateLimiterWithOptions(&RateLimiterOptions{})
			So(err, ShouldBeNil)
			So(r, ShouldBeNil)
		})

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
		So(r, ShouldNotBeNil)

		So(r.Allow(context.Background(), "key1"), ShouldBeNil)
	})
}
