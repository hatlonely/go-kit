package micro

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLocalLocker_Lock_UnLock(t *testing.T) {
	Convey("TestLocalLocker_Lock", t, func(c C) {
		l := NewLocalLocker()
		var wg sync.WaitGroup
		var m int64
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				for j := 0; j < 100; j++ {
					c.So(l.Lock(context.Background(), "key1"), ShouldBeNil)
					atomic.AddInt64(&m, 1)
					c.So(m, ShouldEqual, 1)
					atomic.AddInt64(&m, -1)
					c.So(m, ShouldEqual, 0)
					c.So(l.Unlock(context.Background(), "key1"), ShouldBeNil)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	})
}

func TestLocalLocker_TryLock(t *testing.T) {
	Convey("TestLocalLocker_TryLock", t, func() {
		l := NewLocalLocker()
		So(l.TryLock(context.Background(), "key1"), ShouldBeNil)
		So(l.TryLock(context.Background(), "key1"), ShouldEqual, ErrLocked)
		So(l.Unlock(context.Background(), "key1"), ShouldBeNil)
		So(l.TryLock(context.Background(), "key1"), ShouldBeNil)
	})
}
