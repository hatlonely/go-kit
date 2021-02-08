package micro

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLocalTimedSemaphoreParallelController_Acquire_Release(t *testing.T) {
	Convey("TestLocalTimedSemaphoreParallelController_Acquire_Release", t, func(c C) {
		for i := 1; i < 10; i++ {
			ctl, err := NewLocalTimedSemaphoreParallelControllerWithOptions(&LocalTimedSemaphoreParallelControllerOptions{
				"key1": {
					MaxToken: i,
					Timeout:  time.Second,
				},
			})
			c.So(err, ShouldBeNil)
			var wg sync.WaitGroup
			var m int64
			for j := 0; j < 20; j++ {
				wg.Add(1)
				go func(i int) {
					for k := 0; k < 100; k++ {
						token, err := ctl.Acquire(context.Background(), "key1")
						c.So(token, ShouldNotEqual, 0)
						c.So(err, ShouldBeNil)
						atomic.AddInt64(&m, 1)
						c.So(m, ShouldBeLessThanOrEqualTo, i)
						c.So(m, ShouldBeGreaterThanOrEqualTo, 0)
						atomic.AddInt64(&m, -1)
						c.So(m, ShouldBeLessThanOrEqualTo, i)
						c.So(m, ShouldBeGreaterThanOrEqualTo, 0)
						err = ctl.Release(context.Background(), "key1", token)
						c.So(err, ShouldBeNil)
					}
					wg.Done()
				}(i)
			}
			wg.Wait()
		}
	})
}

func TestLocalTimedSemaphoreParallelController_TryAcquire(t *testing.T) {
	Convey("TestLocalTimedSemaphoreParallelController_TryAcquire", t, func() {
		ctl, err := NewLocalTimedSemaphoreParallelControllerWithOptions(&LocalTimedSemaphoreParallelControllerOptions{
			"key1": {
				MaxToken: 2,
				Timeout:  time.Second,
			},
		})
		So(err, ShouldBeNil)
		token1, err := ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldBeNil)
		So(token1, ShouldNotEqual, 0)
		token2, err := ctl.Acquire(context.Background(), "key1")
		fmt.Println(token2, err)
		So(err, ShouldBeNil)
		So(token2, ShouldNotEqual, 0)
		token, err := ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldEqual, ErrParallelControl)
		So(token, ShouldEqual, 0)
		So(ctl.Release(context.Background(), "key1", token1), ShouldBeNil)
		token3, err := ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldBeNil)
		So(token3, ShouldNotEqual, 0)
		token, err = ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldEqual, ErrParallelControl)
		So(token, ShouldEqual, 0)
	})
}

func TestLocalTimedSemaphoreParallelController_Context_Cancel(t *testing.T) {
	Convey("TestLocalTimedSemaphoreParallelController_Context_Cancel", t, func() {
		ctl, err := NewLocalTimedSemaphoreParallelControllerWithOptions(&LocalTimedSemaphoreParallelControllerOptions{
			"key1": {
				MaxToken: 2,
				Timeout:  time.Second,
			},
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
