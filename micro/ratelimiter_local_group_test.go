package micro

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewLocalGroupRateLimiterWithOptions(t *testing.T) {
	Convey("TestNewLocalGroupRateLimiterWithOptions", t, func() {
		Convey("interval should be positive", func() {
			_, err := NewLocalGroupRateLimiterWithOptions(&LocalGroupRateLimiterOptions{
				"key1": {
					Interval: 0,
					Burst:    2,
				},
			})
			So(err, ShouldNotBeNil)
		})

		Convey("bucket should be positive", func() {
			_, err := NewLocalGroupRateLimiterWithOptions(&LocalGroupRateLimiterOptions{
				"key1": {
					Interval: time.Millisecond,
					Burst:    0,
				},
			})
			So(err, ShouldNotBeNil)
		})
	})
}

func TestLocalGroupRateLimiter(t *testing.T) {
	Convey("TestLocalGroupRateLimiter", t, func() {
		r, err := NewLocalGroupRateLimiterWithOptions(&LocalGroupRateLimiterOptions{
			"key1": {
				Interval: 20 * time.Millisecond,
				Burst:    2,
			},
		})
		So(err, ShouldBeNil)

		now := time.Now()
		So(r.Allow(context.Background(), "key1"), ShouldBeNil)
		So(r.Allow(context.Background(), "key1"), ShouldBeNil)
		So(r.Allow(context.Background(), "key1"), ShouldEqual, ErrFlowControl)
		So(r.Wait(context.Background(), "key1"), ShouldBeNil)
		So(time.Now().Sub(now), ShouldBeGreaterThan, 20*time.Millisecond)

		for i := 0; i < 10; i++ {
			So(r.Allow(context.Background(), "key"), ShouldBeNil)
		}
		now = time.Now()
		for i := 0; i < 10; i++ {
			So(r.Wait(context.Background(), "key"), ShouldBeNil)
		}
		So(time.Now().Sub(now), ShouldBeLessThan, 20*time.Millisecond)
	})
}
