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

// BenchmarkRedisLocker_Lock-12    	    2778	    413410 ns/op
func BenchmarkRedisLocker_Lock(b *testing.B) {
	l, _ := NewRedisLockerWithOptions(&RedisLockerOptions{
		Redis: wrap.RedisClientWrapperOptions{
			Redis: wrap.RedisOptions{
				Addr: "127.0.0.1:6379",
			},
			Retry: micro.RetryOptions{
				Attempts: 1,
			},
		},
		Prefix:     "redis_locker",
		Expiration: 2 * time.Second,
		RenewTime:  1 * time.Second,
		Interval:   time.Microsecond,
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = l.Lock(context.Background(), "key")
			_ = l.Unlock(context.Background(), "key")
		}
	})
}

func TestRedisLocker_Lock_UnLock(t *testing.T) {
	Convey("TestRedisLocker_Lock_UnLock", t, func(c C) {
		l, err := NewRedisLockerWithOptions(&RedisLockerOptions{
			Prefix:     "redis_locker_test",
			Expiration: 30 * time.Second,
			RenewTime:  10 * time.Second,
			Interval:   0,
		})
		So(err, ShouldBeNil)
		var wg sync.WaitGroup
		var m int64
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				for j := 0; j < 100; j++ {
					c.So(l.Lock(context.Background(), "key1"), ShouldBeNil)
					atomic.AddInt64(&m, 1)
					c.So(m, ShouldEqual, 1)
					time.Sleep(time.Millisecond * 10)
					atomic.AddInt64(&m, -1)
					c.So(m, ShouldEqual, 0)
					c.So(l.Unlock(context.Background(), "key1"), ShouldBeNil)
					time.Sleep(time.Millisecond * 10)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
	})
}

func TestRedisLocker_TryLock(t *testing.T) {
	Convey("TestRedisLocker_TryLock", t, func() {
		l, err := NewRedisLockerWithOptions(&RedisLockerOptions{
			Redis: wrap.RedisClientWrapperOptions{
				Retry: micro.RetryOptions{
					Attempts: 1,
				},
			},
			Prefix:     "redis_locker_test",
			Expiration: 30 * time.Second,
			RenewTime:  10 * time.Second,
			Interval:   10 * time.Second,
		})
		So(err, ShouldBeNil)
		So(l.TryLock(context.Background(), "key1"), ShouldBeNil)
		So(l.TryLock(context.Background(), "key1"), ShouldEqual, micro.ErrLocked)
		So(l.Unlock(context.Background(), "key1"), ShouldBeNil)
		So(l.TryLock(context.Background(), "key1"), ShouldBeNil)
		So(l.Unlock(context.Background(), "key1"), ShouldBeNil)
	})
}
