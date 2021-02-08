package micro

import (
	"context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewParallelControllerWithOptions(t *testing.T) {
	Convey("TestNewParallelControllerGroupWithOptions", t, func() {
		Convey("empty options case", func() {
			r, err := NewParallelControllerWithOptions(&ParallelControllerOptions{})
			So(err, ShouldBeNil)
			So(r, ShouldBeNil)
		})

		r, err := NewParallelControllerWithOptions(&ParallelControllerOptions{
			Type: "LocalSemaphore",
			Options: &LocalSemaphoreParallelControllerOptions{
				"key1": 2,
			},
		})
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)

		token1, err := r.Acquire(context.Background(), "key1")
		So(err, ShouldBeNil)
		token2, err := r.Acquire(context.Background(), "key1")
		So(err, ShouldBeNil)
		fmt.Println(token1, token2)
	})
}
