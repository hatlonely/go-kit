package micro

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/emirpasic/gods/sets/treeset"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTimedSemaphore_Acquire_Release(t *testing.T) {
	Convey("TestTimedSemaphore_Acquire_Release", t, func(c C) {
		for i := 1; i < 10; i++ {
			sem, err := NewTimedSemaphore(i, time.Second)
			So(err, ShouldBeNil)
			var wg sync.WaitGroup
			var m int64
			for j := 0; j < 20; j++ {
				wg.Add(1)
				go func(i int) {
					for k := 0; k < 200; k++ {
						res := sem.Acquire()
						atomic.AddInt64(&m, 1)
						c.So(m, ShouldBeLessThanOrEqualTo, i)
						c.So(m, ShouldBeGreaterThanOrEqualTo, 0)
						atomic.AddInt64(&m, -1)
						c.So(m, ShouldBeLessThanOrEqualTo, i)
						c.So(m, ShouldBeGreaterThanOrEqualTo, 0)
						sem.Release(res)
					}
					wg.Done()
				}(i)
			}
			wg.Wait()
		}
	})
}

func TestTimedSemaphore_TryAcquire(t *testing.T) {
	Convey("TestTimedSemaphore_TryAcquire", t, func() {
		sem, err := NewTimedSemaphore(2, time.Second)
		So(err, ShouldBeNil)

		token1, ok := sem.TryAcquire()
		So(ok, ShouldBeTrue)
		So(token1, ShouldNotEqual, 0)
		fmt.Println(TokenSetToString(sem.tokenSet))

		token2 := sem.Acquire()
		So(token2, ShouldNotEqual, 0)
		fmt.Println(TokenSetToString(sem.tokenSet))

		token, ok := sem.TryAcquire()
		So(ok, ShouldBeFalse)
		So(token, ShouldEqual, 0)
		fmt.Println(TokenSetToString(sem.tokenSet))

		So(sem.Release(token2), ShouldBeTrue)
		fmt.Println(TokenSetToString(sem.tokenSet))

		token3, ok := sem.TryAcquire()
		So(ok, ShouldBeTrue)
		So(token3, ShouldNotEqual, 0)
		fmt.Println(TokenSetToString(sem.tokenSet))

		token, ok = sem.TryAcquire()
		So(ok, ShouldBeFalse)
		So(token, ShouldEqual, 0)
		fmt.Println(TokenSetToString(sem.tokenSet))
	})
}

func TestTimedSemaphore_AcquireTimed(t *testing.T) {
	Convey("TestTimedSemaphore_Acquire", t, func() {
		sem, err := NewTimedSemaphore(2, time.Millisecond*100)
		So(err, ShouldBeNil)

		token1 := sem.Acquire()
		fmt.Println(TokenSetToString(sem.tokenSet))
		token2 := sem.Acquire()
		fmt.Println(TokenSetToString(sem.tokenSet))

		token3 := sem.Acquire()
		fmt.Println(TokenSetToString(sem.tokenSet))
		token4 := sem.Acquire()
		fmt.Println(TokenSetToString(sem.tokenSet))

		fmt.Println(token1, token2, token3, token4)
	})
}

func TokenSetToString(set *treeset.Set) string {
	it := set.Iterator()
	var res []string
	for it.Next() {
		res = append(res, fmt.Sprintf("%v", it.Value()))
	}
	return "[" + strings.Join(res, ", ") + "]"
}
