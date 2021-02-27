package microx

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/wrap"
)

// BenchmarkRedisRateLimiter_Wait-12    	  103843	     11429 ns/op
func BenchmarkRedisRateLimiter_Wait(b *testing.B) {
	r, _ := NewRedisRateLimiterWithOptions(&RedisRateLimiterOptions{
		Redis: wrap.RedisClientWrapperOptions{
			Redis: wrap.RedisOptions{
				Addr: "127.0.0.1:6379",
			},
			Retry: micro.RetryOptions{
				Attempts: 3,
				Delay:    time.Millisecond * 500,
			},
		},
		Window:     time.Second,
		DefaultQPS: 9999999999,
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Wait(context.Background(), "key1")
		}
	})
}

func TestRedisRateLimiter(t *testing.T) {
	Convey("TestRedisRateLimiter", t, func() {
		r, err := NewRedisRateLimiterWithOptions(&RedisRateLimiterOptions{
			Redis: wrap.RedisClientWrapperOptions{
				Retry: micro.RetryOptions{
					Attempts: 3,
					Delay:    time.Millisecond * 500,
				},
			},
			Window:     time.Second,
			DefaultQPS: 3,
			QPS: map[string]int{
				"key1": 2,
				"key2": 0,
			},
		})
		So(err, ShouldBeNil)

		{
			time.Sleep(time.Until(time.Now().Truncate(time.Second).Add(time.Second)))
			now := time.Now()
			So(r.Allow(context.Background(), "key1"), ShouldBeNil)
			So(r.Allow(context.Background(), "key1"), ShouldBeNil)
			So(r.Allow(context.Background(), "key1"), ShouldEqual, micro.ErrFlowControl)
			So(r.Wait(context.Background(), "key1"), ShouldBeNil)
			So(time.Now().Sub(now), ShouldBeGreaterThan, 20*time.Millisecond)
		}
		{
			time.Sleep(time.Until(time.Now().Truncate(time.Second).Add(time.Second)))
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
