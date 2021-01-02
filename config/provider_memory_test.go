package config

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMemoryProvider(t *testing.T) {
	Convey("TestMemoryProvider_Load", t, func() {
		provider := NewMemoryProviderWithOptions(&MemoryProviderOptions{
			Buffer: "hello world",
		})

		Convey("Events", func() {
			So(func() {
				provider.Events()
			}, ShouldPanic)
		})
		Convey("Errors", func() {
			So(func() {
				provider.Errors()
			}, ShouldPanic)
		})
		Convey("EventLoop", func() {
			So(func() {
				provider.EventLoop(context.Background())
			}, ShouldPanic)
		})
		Convey("Load", func() {
			buf, err := provider.Load()
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, "hello world")
		})
		Convey("Dump", func() {
			So(provider.Dump([]byte(`hello golang`)), ShouldBeNil)
			So(provider.buf, ShouldResemble, []byte(`hello golang`))
		})
	})
}
