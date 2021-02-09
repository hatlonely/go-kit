package micro

import (
	"sync"
	"sync/atomic"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSemaphore(t *testing.T) {
	Convey("TestNewSemaphore", t, func() {
		_, err := NewSemaphore(0)
		So(err, ShouldNotBeNil)
	})
}

func TestSemaphore_Acquire_Release(t *testing.T) {
	Convey("TestSemaphore_Acquire_Release", t, func(c C) {
		for i := 1; i < 10; i++ {
			sem, err := NewSemaphore(i)
			So(err, ShouldBeNil)
			var wg sync.WaitGroup
			var m int64
			for j := 0; j < 20; j++ {
				wg.Add(1)
				go func(i int) {
					for k := 0; k < 200; k++ {
						sem.Acquire()
						atomic.AddInt64(&m, 1)
						c.So(m, ShouldBeLessThanOrEqualTo, i)
						c.So(m, ShouldBeGreaterThanOrEqualTo, 0)
						atomic.AddInt64(&m, -1)
						c.So(m, ShouldBeLessThanOrEqualTo, i)
						c.So(m, ShouldBeGreaterThanOrEqualTo, 0)
						sem.Release()
					}
					wg.Done()
				}(i)
			}
			wg.Wait()
		}
	})
}

func TestSemaphore_TryAcquire(t *testing.T) {
	Convey("TestSemaphore_TryAcquire", t, func() {
		sem, err := NewSemaphore(2)
		So(err, ShouldBeNil)
		So(sem.TryAcquire(), ShouldBeTrue)
		sem.Acquire()
		So(sem.TryAcquire(), ShouldBeFalse)
		sem.Release()
		So(sem.TryAcquire(), ShouldBeTrue)
		So(sem.TryAcquire(), ShouldBeFalse)
	})
}

func TestSemaphore_ForceRelease(t *testing.T) {
	Convey("TestSemaphore_ForceRelease", t, func() {
		sem, err := NewSemaphore(2)
		So(err, ShouldBeNil)
		So(sem.curToken, ShouldEqual, 0)
		sem.ForceRelease()
		So(sem.curToken, ShouldEqual, 0)
		sem.Acquire()
		So(sem.curToken, ShouldEqual, 1)
		sem.Acquire()
		So(sem.curToken, ShouldEqual, 2)
		So(sem.TryAcquire(), ShouldBeFalse)
		So(sem.curToken, ShouldEqual, 2)
		sem.ForceRelease()
		So(sem.curToken, ShouldEqual, 1)
		So(sem.TryAcquire(), ShouldBeTrue)
		So(sem.curToken, ShouldEqual, 2)
		So(sem.TryAcquire(), ShouldBeFalse)
		So(sem.curToken, ShouldEqual, 2)
		sem.ForceRelease()
		So(sem.curToken, ShouldEqual, 1)
		sem.ForceRelease()
		So(sem.curToken, ShouldEqual, 0)
		sem.ForceRelease()
		So(sem.curToken, ShouldEqual, 0)
	})
}
