package config

import (
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAtomicBool(t *testing.T) {
	Convey("TestAtomicBool", t, func() {
		val := NewAtomicBool(true)
		So(val.Get(), ShouldEqual, true)
		val.Set(false)
		So(val.Get(), ShouldEqual, false)
	})
}

func TestAtomicInt(t *testing.T) {
	Convey("TestAtomicInt", t, func() {
		val := NewAtomicInt(10)
		So(val.Get(), ShouldEqual, 10)
		val.Set(20)
		So(val.Get(), ShouldEqual, 20)
	})
}

func TestAtomicUint(t *testing.T) {
	Convey("TestAtomicUint", t, func() {
		val := NewAtomicUint(10)
		So(val.Get(), ShouldEqual, 10)
		val.Set(20)
		So(val.Get(), ShouldEqual, 20)
	})
}

func TestAtomicInt64(t *testing.T) {
	Convey("TestAtomicInt64", t, func() {
		val := NewAtomicInt64(10)
		So(val.Get(), ShouldEqual, 10)
		val.Set(20)
		So(val.Get(), ShouldEqual, 20)
	})
}

func TestAtomicInt32(t *testing.T) {
	Convey("TestAtomicInt32", t, func() {
		val := NewAtomicInt32(10)
		So(val.Get(), ShouldEqual, 10)
		val.Set(20)
		So(val.Get(), ShouldEqual, 20)
	})
}

func TestAtomicInt16(t *testing.T) {
	Convey("TestAtomicInt16", t, func() {
		val := NewAtomicInt16(10)
		So(val.Get(), ShouldEqual, 10)
		val.Set(20)
		So(val.Get(), ShouldEqual, 20)
	})
}

func TestAtomicInt8(t *testing.T) {
	Convey("TestAtomicInt8", t, func() {
		val := NewAtomicInt8(10)
		So(val.Get(), ShouldEqual, 10)
		val.Set(20)
		So(val.Get(), ShouldEqual, 20)
	})
}

func TestAtomicUint64(t *testing.T) {
	Convey("TestAtomicUint64", t, func() {
		val := NewAtomicUint64(10)
		So(val.Get(), ShouldEqual, 10)
		val.Set(20)
		So(val.Get(), ShouldEqual, 20)
	})
}

func TestAtomicUint32(t *testing.T) {
	Convey("TestAtomicUint32", t, func() {
		val := NewAtomicUint32(10)
		So(val.Get(), ShouldEqual, 10)
		val.Set(20)
		So(val.Get(), ShouldEqual, 20)
	})
}

func TestAtomicUint16(t *testing.T) {
	Convey("TestAtomicUint16", t, func() {
		val := NewAtomicUint16(10)
		So(val.Get(), ShouldEqual, 10)
		val.Set(20)
		So(val.Get(), ShouldEqual, 20)
	})
}

func TestAtomicUint8(t *testing.T) {
	Convey("TestAtomicUint8", t, func() {
		val := NewAtomicUint8(10)
		So(val.Get(), ShouldEqual, 10)
		val.Set(20)
		So(val.Get(), ShouldEqual, 20)
	})
}

func TestAtomicFloat64(t *testing.T) {
	Convey("TestAtomicFloat64", t, func() {
		val := NewAtomicFloat64(10)
		So(val.Get(), ShouldAlmostEqual, 10)
		val.Set(20)
		So(val.Get(), ShouldAlmostEqual, 20)
	})
}

func TestAtomicFloat32(t *testing.T) {
	Convey("TestAtomicFloat32", t, func() {
		val := NewAtomicFloat32(10)
		So(val.Get(), ShouldAlmostEqual, 10)
		val.Set(20)
		So(val.Get(), ShouldAlmostEqual, 20)
	})
}

func TestAtomicString(t *testing.T) {
	Convey("TestAtomicString", t, func() {
		val := NewAtomicString("hello")
		So(val.Get(), ShouldEqual, "hello")
		val.Set("world")
		So(val.Get(), ShouldEqual, "world")
	})
}

func TestAtomicDuration(t *testing.T) {
	Convey("TestAtomicDuration", t, func() {
		val := NewAtomicDuration(3 * time.Second)
		So(val.Get(), ShouldEqual, 3*time.Second)
		val.Set(5 * time.Minute)
		So(val.Get(), ShouldEqual, 5*time.Minute)
	})
}

func TestAtomicTime(t *testing.T) {
	Convey("TestAtomicTime", t, func() {
		val := NewAtomicTime(time.Unix(1604922296, 0))
		So(val.Get(), ShouldEqual, time.Unix(1604922296, 0))
		val.Set(time.Unix(1604922297, 0))
		So(val.Get(), ShouldEqual, time.Unix(1604922297, 0))
	})
}

func TestAtomicIP(t *testing.T) {
	Convey("TestAtomicIP", t, func() {
		val := NewAtomicIP(net.ParseIP("106.11.34.42"))
		So(val.Get(), ShouldEqual, net.ParseIP("106.11.34.42"))
		val.Set(net.ParseIP("106.11.34.43"))
		So(val.Get(), ShouldEqual, net.ParseIP("106.11.34.43"))
	})
}
