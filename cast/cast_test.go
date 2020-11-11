package cast

import (
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestToBool(t *testing.T) {
	Convey("TestToBool", t, func() {
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

func TestToInt8(t *testing.T) {
	Convey("TestToInt8", t, func() {
		for _, v := range []interface{}{
			123, "123", int64(123), 123.0,
		} {
			So(ToInt8(v), ShouldEqual, 123)
		}
	})
}

func TestToInt16(t *testing.T) {
	Convey("TestToInt16", t, func() {
		for _, v := range []interface{}{
			12345, "12345", int64(12345), 12345.123,
		} {
			So(ToInt16(v), ShouldEqual, 12345)
		}
	})
}

func TestToInt32(t *testing.T) {
	Convey("TestToInt32", t, func() {
		for _, v := range []interface{}{
			1234567890, "1234567890", int64(1234567890), 1234567890.123,
		} {
			So(ToInt32(v), ShouldEqual, 1234567890)
		}
	})
}

func TestToInt64(t *testing.T) {
	Convey("TestToInt64", t, func() {
		for _, v := range []interface{}{
			1234567890123456789, "1234567890123456789", int64(1234567890123456789),
		} {
			So(ToInt64(v), ShouldEqual, 1234567890123456789)
		}
	})
}

func TestToUint(t *testing.T) {
	Convey("TestToUint", t, func() {
		for _, v := range []interface{}{
			123, "123", int64(123), 123.0,
		} {
			So(ToUint(v), ShouldEqual, 123)
		}
	})
}

func TestToUint8(t *testing.T) {
	Convey("TestToUint", t, func() {
		for _, v := range []interface{}{
			123, "123", int64(123), 123.0,
		} {
			So(ToUint8(v), ShouldEqual, 123)
		}
	})
}

func TestToUint16(t *testing.T) {
	Convey("TestToUint16", t, func() {
		for _, v := range []interface{}{
			12345, "12345", int64(12345), 12345.123,
		} {
			So(ToUint16(v), ShouldEqual, 12345)
		}
	})
}

func TestToUint32(t *testing.T) {
	Convey("TestToUint32", t, func() {
		for _, v := range []interface{}{
			1234567890, "1234567890", int64(1234567890), 1234567890.123,
		} {
			So(ToUint32(v), ShouldEqual, 1234567890)
		}
	})
}

func TestToUint64(t *testing.T) {
	Convey("TestToUint64", t, func() {
		for _, v := range []interface{}{
			1234567890123456789, "1234567890123456789", int64(1234567890123456789),
		} {
			So(ToUint64(v), ShouldEqual, 1234567890123456789)
		}
	})
}

func TestToFloat32(t *testing.T) {
	Convey("TestToFloat32", t, func() {
		for _, v := range []interface{}{
			123.456, "123.456",
		} {
			So(ToFloat32(v), ShouldEqual, 123.456)
		}
	})
}

func TestToFloat64(t *testing.T) {
	Convey("TestToFloat64", t, func() {
		for _, v := range []interface{}{
			123.456, "123.456",
		} {
			So(ToFloat64(v), ShouldEqual, 123.456)
		}
	})
}

func TestToString(t *testing.T) {
	Convey("TestToString", t, func() {
		for _, v := range []interface{}{
			123.456, "123.456",
		} {
			So(ToString(v), ShouldEqual, "123.456")
		}
	})
}

func TestToTime(t *testing.T) {
	Convey("TestToTime", t, func() {
		for _, v := range []interface{}{
			"2020-11-12T00:58:55+08:00",
			"2020-11-11 16:58:55",
			1605113935,
		} {
			So(ToTime(v), ShouldEqual, time.Unix(1605113935, 0))
		}
	})
}

func TestToDuration(t *testing.T) {
	Convey("TestToDuration", t, func() {
		for _, v := range []interface{}{
			20 * time.Millisecond,
			"20ms",
		} {
			So(ToDuration(v), ShouldEqual, 20*time.Millisecond)
		}
	})
}

func TestToIP(t *testing.T) {
	Convey("TestToIP", t, func() {
		for _, v := range []interface{}{
			"192.168.0.1", net.ParseIP("192.168.0.1"),
		} {
			So(ToIP(v), ShouldEqual, net.ParseIP("192.168.0.1"))
		}
	})
}
