package microx

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/wrap"
)

func TestRedisRateLimiter(t *testing.T) {
	Convey("TestRedisRateLimiter", t, func() {
		r, err := NewRedisRateLimiterWithOptions(&RedisRateLimiterOptions{
			Redis: wrap.RedisClientWrapperOptions{
				Retry: micro.RetryOptions{
					Attempts: 3,
					Delay:    time.Millisecond * 500,
				},
			},
			DefaultQPS: 3,
			QPS: map[string]int{
				"key1": 2,
				"key2": 0,
			},
		})
		So(err, ShouldBeNil)

		{
			now := time.Now()
			So(r.Allow(context.Background(), "key1"), ShouldBeNil)
			So(r.Allow(context.Background(), "key1"), ShouldBeNil)
			So(r.Allow(context.Background(), "key1"), ShouldEqual, micro.ErrFlowControl)
			So(r.Wait(context.Background(), "key1"), ShouldBeNil)
			So(time.Now().Sub(now), ShouldBeGreaterThan, 20*time.Millisecond)
		}
		{
			now := time.Now()
			So(r.Allow(context.Background(), "key3"), ShouldBeNil)
			So(r.Allow(context.Background(), "key3"), ShouldBeNil)
			So(r.Allow(context.Background(), "key3"), ShouldBeNil)
			So(r.Allow(context.Background(), "key3"), ShouldEqual, micro.ErrFlowControl)
			So(r.Wait(context.Background(), "key3"), ShouldBeNil)
			So(time.Now().Sub(now), ShouldBeGreaterThan, 20*time.Millisecond)
		}

		for i := 0; i < 10; i++ {
			So(r.Allow(context.Background(), "key2"), ShouldBeNil)
		}
		now := time.Now()
		for i := 0; i < 10; i++ {
			So(r.Wait(context.Background(), "key2"), ShouldBeNil)
		}
		So(time.Now().Sub(now), ShouldBeLessThan, 20*time.Millisecond)
	})
}
