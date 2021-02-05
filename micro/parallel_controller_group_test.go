package micro

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewParallelControllerGroupWithOptions(t *testing.T) {
	Convey("TestNewParallelControllerGroupWithOptions", t, func() {
		Convey("empty options case", func() {
			r, err := NewParallelControllerWithOptions(&ParallelControllerGroupOptions{})
			So(err, ShouldBeNil)
			So(r, ShouldBeNil)
		})

		r, err := NewParallelControllerWithOptions(&ParallelControllerGroupOptions{
			Type: "LocalGroup",
			Options: &LocalParallelControllerOptions{
				"key1": 2,
			},
		})
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)

		So(r.GetToken(context.Background(), "key1"), ShouldBeNil)
		So(r.GetToken(context.Background(), "key1"), ShouldBeNil)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		So(r.GetToken(ctx, "key1"), ShouldNotBeNil)
	})
}
