package micro

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewLocalShareRateLimiterWithOptions(t *testing.T) {
	Convey("TestNewLocalShareRateLimiterWithOptions", t, func() {
		Convey("cost should be positive", func() {
			_, err := NewLocalShareRateLimiterWithOptions(&LocalShareRateLimiterOptions{
				Interval: 20 * time.Millisecond,
				Burst:    4,
				Cost: map[string]int{
					"key1": 0,
				},
			})
			So(err, ShouldNotBeNil)
		})
	})
}

func TestLocalShareRateLimiter(t *testing.T) {
	Convey("TestLocalGroupRateLimiter", t, func() {
		r, err := NewLocalShareRateLimiterWithOptions(&LocalShareRateLimiterOptions{
			Interval: 20 * time.Millisecond,
			Burst:    4,
			Cost: map[string]int{
				"key1": 2,
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
