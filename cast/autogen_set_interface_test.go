package cast

import (
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetInterfaceType(t *testing.T) {
	Convey("TestSetInterfaceType", t, func() {
		{
			var val bool
			So(SetInterface(&val, "true"), ShouldBeNil)
			So(val, ShouldBeTrue)
		}
		{
			var val int
			So(SetInterface(&val, "10"), ShouldBeNil)
			So(val, ShouldEqual, 10)
		}
		{
			var val int8
			So(SetInterface(&val, "10"), ShouldBeNil)
			So(val, ShouldEqual, 10)
		}
		{
			var val int16
			So(SetInterface(&val, "10"), ShouldBeNil)
			So(val, ShouldEqual, 10)
		}
		{
			var val int32
			So(SetInterface(&val, "10"), ShouldBeNil)
			So(val, ShouldEqual, 10)
		}
		{
			var val int64
			So(SetInterface(&val, "10"), ShouldBeNil)
			So(val, ShouldEqual, 10)
		}
		{
			var val uint
			So(SetInterface(&val, "10"), ShouldBeNil)
			So(val, ShouldEqual, 10)
		}
		{
			var val uint8
			So(SetInterface(&val, "10"), ShouldBeNil)
			So(val, ShouldEqual, 10)
		}
		{
			var val uint16
			So(SetInterface(&val, "10"), ShouldBeNil)
			So(val, ShouldEqual, 10)
		}
		{
			var val uint32
			So(SetInterface(&val, "10"), ShouldBeNil)
			So(val, ShouldEqual, 10)
		}
		{
			var val uint64
			So(SetInterface(&val, "10"), ShouldBeNil)
			So(val, ShouldEqual, 10)
		}
		{
			var val float32
			So(SetInterface(&val, "1.23"), ShouldBeNil)
			So(val, ShouldBeBetween, 1.22, 1.24)
		}
		{
			var val float64
			So(SetInterface(&val, "123.456"), ShouldBeNil)
			So(val, ShouldAlmostEqual, 123.456)
		}
		{
			var val string
			So(SetInterface(&val, 123), ShouldBeNil)
			So(val, ShouldEqual, "123")
		}
		{
			var val time.Duration
			So(SetInterface(&val, "5s"), ShouldBeNil)
			So(val, ShouldEqual, 5*time.Second)
		}
		{
			var val net.IP
			So(SetInterface(&val, "192.168.0.1"), ShouldBeNil)
			So(val, ShouldResemble, net.ParseIP("192.168.0.1"))
		}
		{
			var val time.Time
			So(SetInterface(&val, "2020-11-12T12:48:15+08:00"), ShouldBeNil)
			So(val, ShouldEqual, time.Unix(1605156495, 0))
		}
	})
}

func TestSetInterfaceSlice(t *testing.T) {
	Convey("TestSetInterfaceSlice", t, func() {
		{
			var val []bool
			So(SetInterface(&val, "true,false"), ShouldBeNil)
			So(val, ShouldResemble, []bool{true, false})
		}
		{
			var val []int
			So(SetInterface(&val, "1,2,3"), ShouldBeNil)
			So(val, ShouldResemble, []int{1, 2, 3})
		}
		{
			var val []int8
			So(SetInterface(&val, "1,2,3"), ShouldBeNil)
			So(val, ShouldResemble, []int8{1, 2, 3})
		}
		{
			var val []int16
			So(SetInterface(&val, "1,2,3"), ShouldBeNil)
			So(val, ShouldResemble, []int16{1, 2, 3})
		}
		{
			var val []int32
			So(SetInterface(&val, "1,2,3"), ShouldBeNil)
			So(val, ShouldResemble, []int32{1, 2, 3})
		}
		{
			var val []int64
			So(SetInterface(&val, "1,2,3"), ShouldBeNil)
			So(val, ShouldResemble, []int64{1, 2, 3})
		}
		{
			var val []uint
			So(SetInterface(&val, "1,2,3"), ShouldBeNil)
			So(val, ShouldResemble, []uint{1, 2, 3})
		}
		{
			var val []uint8
			So(SetInterface(&val, "1,2,3"), ShouldBeNil)
			So(val, ShouldResemble, []uint8{1, 2, 3})
		}
		{
			var val []uint16
			So(SetInterface(&val, "1,2,3"), ShouldBeNil)
			So(val, ShouldResemble, []uint16{1, 2, 3})
		}
		{
			var val []uint32
			So(SetInterface(&val, "1,2,3"), ShouldBeNil)
			So(val, ShouldResemble, []uint32{1, 2, 3})
		}
		{
			var val []uint64
			So(SetInterface(&val, "1,2,3"), ShouldBeNil)
			So(val, ShouldResemble, []uint64{1, 2, 3})
		}
		{
			var val []float32
			So(SetInterface(&val, "123.456,234.567"), ShouldBeNil)
			So(nearlyEqualFloat32Slice(val, []float32{123.456, 234.567}), ShouldBeTrue)
		}
		{
			var val []float64
			So(SetInterface(&val, "123.456,234.567"), ShouldBeNil)
			So(nearlyEqualFloat64Slice(val, []float64{123.456, 234.567}), ShouldBeTrue)
		}
		{
			var val []string
			So(SetInterface(&val, "1,2,3"), ShouldBeNil)
			So(val, ShouldResemble, []string{"1", "2", "3"})
		}
		{
			var val []time.Duration
			So(SetInterface(&val, "1s,2ms,3m"), ShouldBeNil)
			So(val, ShouldResemble, []time.Duration{time.Second, 2 * time.Millisecond, 3 * time.Minute})
		}
		{
			var val []time.Time
			So(SetInterface(&val, "1970-01-01 08:00:01 +0800 CST,1970-01-01T08:00:02+08:00,1970-01-01T08:00:03+08:00"), ShouldBeNil)
			So(val, ShouldResemble, []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)})
		}
		{
			var val []net.IP
			So(SetInterface(&val, "192.168.0.1,192.168.0.2"), ShouldBeNil)
			So(val, ShouldResemble, []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")})
		}
	})
}

func TestSetInterfaceError(t *testing.T) {
	Convey("TestSetInterfaceError", t, func() {
		{
			var val bool
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val int
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val int8
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val int16
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val int32
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val int64
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val uint
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val uint8
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val uint16
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val uint32
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val uint64
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val float32
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val float64
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val time.Duration
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val net.IP
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val time.Time
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}

		{
			var val []bool
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []int
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []int8
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []int16
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []int32
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []int64
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []uint
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []uint8
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []uint16
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []uint32
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []uint64
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []float32
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []float64
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []time.Duration
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []time.Time
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
		{
			var val []net.IP
			So(SetInterface(&val, "x"), ShouldNotBeNil)
		}
	})
}
