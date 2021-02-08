package micro

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLocalSemaphoreParallelController_Acquire_Release(t *testing.T) {
	Convey("TestLocalSemaphoreParallelController_Acquire_Release", t, func(c C) {
		for i := 1; i < 10; i++ {
			ctl, err := NewLocalSemaphoreParallelControllerWithOptions(&LocalSemaphoreParallelControllerOptions{
				"key1": i,
			})
			c.So(err, ShouldBeNil)
			var wg sync.WaitGroup
			var m int64
			for j := 0; j < 20; j++ {
				wg.Add(1)
				go func(i int) {
					for k := 0; k < 100; k++ {
						res, err := ctl.Acquire(context.Background(), "key1")
						c.So(res, ShouldEqual, 0)
						c.So(err, ShouldBeNil)
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

func TestLocalSemaphoreParallelController_TryAcquire(t *testing.T) {
	Convey("TestLocalSemaphoreParallelController_TryAcquire", t, func() {
		ctl, err := NewLocalSemaphoreParallelControllerWithOptions(&LocalSemaphoreParallelControllerOptions{
			"key1": 2,
		})
		So(err, ShouldBeNil)
		res, err := ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldBeNil)
		So(res, ShouldEqual, 0)
		res, err = ctl.Acquire(context.Background(), "key1")
		So(err, ShouldBeNil)
		So(res, ShouldEqual, 0)
		res, err = ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldEqual, ErrParallelControl)
		So(res, ShouldEqual, 0)
		So(ctl.Release(context.Background(), "key1", res), ShouldBeNil)
		res, err = ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldBeNil)
		So(res, ShouldEqual, 0)
		res, err = ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldEqual, ErrParallelControl)
		So(res, ShouldEqual, 0)
	})
}

func TestLocalSemaphoreParallelController_Context_Cancel(t *testing.T) {
	Convey("TestLocalSemaphoreParallelController_Context_Cancel", t, func() {
		ctl, err := NewLocalSemaphoreParallelControllerWithOptions(&LocalSemaphoreParallelControllerOptions{
			"key1": 2,
		})
		So(err, ShouldBeNil)
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
			So(err, ShouldEqual, ErrContextCancel)
		})
	})
}
