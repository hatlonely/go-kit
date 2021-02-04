package micro

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test1(t *testing.T) {
	Convey("test", t, func() {
		pc, err := NewLocalParallelControllerGroupWithOptions(&LocalParallelControllerGroupOptions{
			"key1": 1,
			"key2": 2,
		})
		So(err, ShouldBeNil)

		var wg sync.WaitGroup

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				for j := 0; j < 10; j++ {
					pc.GetToken(context.Background(), "key1")
					fmt.Println("GetToken key1", i)
					time.Sleep(time.Second)
					fmt.Println("PutToken key1", i)
					pc.PutToken(context.Background(), "key1")
				}
				wg.Done()
			}(i)
		}

		wg.Wait()
	})
}
