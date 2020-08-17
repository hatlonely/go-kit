package cast

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestToBool(t *testing.T) {
	Convey("TestToBoolE", t, func() {
		for _, v := range []interface{}{
			"true", "t", "T", "True", 1, true,
		} {
			So(ToBool(v), ShouldBeTrue)
		}
		for _, v := range []interface{}{
			"false", "f", "F", "false", 0, false,
		} {
			So(ToBool(v), ShouldBeFalse)
		}
	})
}

func TestToInt(t *testing.T) {
	Convey("TestToInt", t, func() {
		for _, v := range []interface{}{
			123, "123", int64(123), 123.0,
		} {
			So(ToInt(v), ShouldEqual, 123)
		}
	})
}
