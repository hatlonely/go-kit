package microx

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/wrap"
)

// BenchmarkRedisTimedParallelController_Acquire-12    	   41810	     29157 ns/op
func BenchmarkRedisTimedParallelController_Acquire(b *testing.B) {
	ctl, _ := NewRedisTimedParallelControllerWithOptions(&RedisTimedParallelControllerOptions{
		Redis: wrap.RedisClientWrapperOptions{
			Redis: wrap.RedisOptions{
				Addr: "127.0.0.1:6379",
			},
			Retry: micro.RetryOptions{
				Attempts: 1,
			},
		},
		Prefix:          "pc",
		DefaultMaxToken: 999999999,
		Interval:        time.Second,
		Expiration:      5 * time.Second,
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token, _ := ctl.Acquire(context.Background(), "key1")
			_ = ctl.Release(context.Background(), "key1", token)
		}
	})
}

func TestRedisTimedParallelController_Acquire_Release(t *testing.T) {
	Convey("TestRedisTimedParallelController_Acquire_Release", t, func(c C) {
		for i := 1; i < 10; i++ {
			ctl, err := NewRedisTimedParallelControllerWithOptions(&RedisTimedParallelControllerOptions{
				Redis: wrap.RedisClientWrapperOptions{
					Retry: micro.RetryOptions{
						Attempts: 1,
					},
				},
				Prefix: "rtpc",
				MaxToken: map[string]int{
					"key1": i,
				},
				Interval:   time.Millisecond,
				Expiration: time.Second,
			})
			c.So(err, ShouldBeNil)
			ctl.client.Del(context.Background(), "rtpc_key1")
			var wg sync.WaitGroup
			var m int64
			for j := 0; j < 20; j++ {
				wg.Add(1)
				go func(i int) {
					for k := 0; k < 100; k++ {
						res, err := ctl.Acquire(context.Background(), "key1")
						c.So(err, ShouldBeNil)
						c.So(res, ShouldNotEqual, 0)
						atomic.AddInt64(&m, 1)
						c.So(m, ShouldBeLessThanOrEqualTo, i)
						c.So(m, ShouldBeGreaterThanOrEqualTo, 0)
						atomic.AddInt64(&m, -1)
						c.So(m, ShouldBeLessThanOrEqualTo, i)
						c.So(m, ShouldBeGreaterThanOrEqualTo, 0)
						err = ctl.Release(context.Background(), "key1", res)
						c.So(err, ShouldBeNil)
					}
					wg.Done()
				}(i)
			}
			wg.Wait()
		}
	})
}

func TestRedisTimedParallelController_TryAcquire(t *testing.T) {
	Convey("TestRedisTimedParallelController_TryAcquire", t, func() {
		ctl, err := NewRedisTimedParallelControllerWithOptions(&RedisTimedParallelControllerOptions{
			Redis: wrap.RedisClientWrapperOptions{
				Retry: micro.RetryOptions{
					Attempts: 1,
				},
			},
			Prefix: "rtpc",
			MaxToken: map[string]int{
				"key1": 2,
			},
			Interval:   time.Second,
			Expiration: time.Second * 10,
		})
		So(err, ShouldBeNil)
		ctl.client.Del(context.Background(), "rtpc_key1")
		token1, err := ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldBeNil)
		So(token1, ShouldNotEqual, 0)
		token2, err := ctl.Acquire(context.Background(), "key1")
		So(err, ShouldBeNil)
		So(token2, ShouldNotEqual, 0)
		token3, err := ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldEqual, micro.ErrParallelControl)
		So(token3, ShouldEqual, 0)
		So(ctl.Release(context.Background(), "key1", token1), ShouldBeNil)
		token4, err := ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldBeNil)
		So(token4, ShouldNotEqual, 0)
		token5, err := ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldEqual, micro.ErrParallelControl)
		So(token5, ShouldEqual, 0)
	})
}

func TestRedisTimedParallelController_ContextCancel(t *testing.T) {
	Convey("TestRedisTimedParallelController_ContextCancel", t, func() {
		ctl, err := NewRedisTimedParallelControllerWithOptions(&RedisTimedParallelControllerOptions{
			Redis: wrap.RedisClientWrapperOptions{
				Retry: micro.RetryOptions{
					Attempts: 1,
				},
			},
			Prefix: "rtpc",
			MaxToken: map[string]int{
				"key1": 2,
			},
			Interval:   time.Second,
			Expiration: time.Second * 10,
		})
		So(err, ShouldBeNil)
		ctl.client.Del(context.Background(), "rtpc_key1")
		Convey("key not match", func() {
			for i := 0; i < 10; i++ {
				res, err := ctl.Acquire(context.Background(), "key2")
				So(err, ShouldBeNil)
				So(res, ShouldEqual, 0)
			}
			for i := 0; i < 10; i++ {
				res, err := ctl.TryAcquire(context.Background(), "key2")
				So(err, ShouldBeNil)
				So(res, ShouldEqual, 0)
			}
			for i := 0; i < 10; i++ {
				err := ctl.Release(context.Background(), "key2", 0)
				So(err, ShouldBeNil)
			}
		})

		Convey("context cancel", func() {
			_, _ = ctl.Acquire(context.Background(), "key1")
			_, _ = ctl.Acquire(context.Background(), "key1")
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			_, err := ctl.Acquire(ctx, "key1")
			So(err, ShouldEqual, micro.ErrContextCancel)
		})
	})
}
