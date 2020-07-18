package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLocalProvider(t *testing.T) {
	Convey("", t, func() {
		provider, err := NewLocalProvider("test.txt")
		So(err, ShouldBeNil)
		{
			So(provider.Dump([]byte("hello world")), ShouldBeNil)
			buf, err := ioutil.ReadFile("test.txt")
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, "hello world")
		}
		{
			buf, err := provider.Load()
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, "hello world")
		}

		{
			ctx, cancel := context.WithCancel(context.Background())
			So(provider.EventLoop(ctx), ShouldBeNil)

			buf, err := provider.Load()
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, "hello world")

			for i := 0; i < 5; i++ {
				So(ioutil.WriteFile("test.txt", []byte(fmt.Sprintf("hello world %v", i)), 0644), ShouldBeNil)
				<-provider.Events()
				buf, err = provider.Load()
				So(err, ShouldBeNil)
				So(string(buf), ShouldEqual, fmt.Sprintf("hello world %v", i))
			}

			cancel()
		}

		_ = os.Remove("test.txt")
	})
}
