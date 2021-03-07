package micro

import (
	"errors"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/refx"
)

func TestNewRetryWithOptions(t *testing.T) {
	Convey("TestNewRetryWithOptions", t, func() {
		var options RetryOptions
		refx.SetDefaultValueCopyP(&options)
		Convey("default", func() {
			So(options.Attempts, ShouldEqual, 3)
			So(options.Delay, ShouldEqual, time.Second)
			So(options.MaxDelay, ShouldEqual, 3*time.Minute)
			So(options.MaxJitter, ShouldEqual, time.Second)
			So(options.LastErrorOnly, ShouldBeFalse)
			So(options.DelayType, ShouldEqual, "BackOff")
			retry, err := NewRetryWithOptions(&options)
			So(err, ShouldBeNil)
			So(retry, ShouldNotBeNil)
		})

		Convey("unknown delay type", func() {
			options.DelayType = "unknown"
			_, err := NewRetryWithOptions(&options)
			So(err, ShouldNotBeNil)
		})

		Convey("unknown retry if", func() {
			options.RetryIf = "unknown"
			_, err := NewRetryWithOptions(&options)
			So(err, ShouldNotBeNil)
		})

		Convey("register retry if", func() {
			RegisterRetryRetryIf("test", func(err error) bool {
				return false
			})
			options.RetryIf = "test"
			_, err := NewRetryWithOptions(&options)
			So(err, ShouldBeNil)
		})
	})
}

func TestParseDelayType(t *testing.T) {
	Convey("TestParseDelayType", t, func() {
		fun, err := parseDelayType("Fixed,Random")
		So(err, ShouldBeNil)
		So(fun, ShouldNotBeNil)
	})
}

func TestRetry_Do(t *testing.T) {
	Convey("TestRetry_Do", t, func() {
		var options RetryOptions
		refx.SetDefaultValueCopyP(&options)
		retry, err := NewRetryWithOptions(&options)
		So(err, ShouldBeNil)

		i := 0
		err = retry.Do(func() error {
			i++
			if i%2 != 0 {
				return errors.New("timeout")
			}
			return nil
		})
		So(err, ShouldBeNil)
		So(i, ShouldEqual, 2)
	})
}
