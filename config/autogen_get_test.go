package config

import (
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfigGet(t *testing.T) {
	Convey("TestConfigGet", t, func() {
		CreateTestFile3()

		cfg, err := NewConfigWithSimpleFile("test.json")
		So(err, ShouldBeNil)

		So(cfg.GetBool("boolVal"), ShouldBeTrue)
		So(cfg.GetInt("intVal"), ShouldEqual, 11)
		So(cfg.GetInt8("int8Val"), ShouldEqual, 22)
		So(cfg.GetInt16("int16Val"), ShouldEqual, 33)
		So(cfg.GetInt32("int32Val"), ShouldEqual, 44)
		So(cfg.GetInt64("int64Val"), ShouldEqual, 55)
		So(cfg.GetUint("uintVal"), ShouldEqual, 66)
		So(cfg.GetUint8("uint8Val"), ShouldEqual, 77)
		So(cfg.GetUint16("uint16Val"), ShouldEqual, 88)
		So(cfg.GetUint32("uint32Val"), ShouldEqual, 99)
		So(cfg.GetUint64("uint64Val"), ShouldEqual, 1010)
		So(cfg.GetFloat32("float32Val"), ShouldBeBetween, 32.31, 32.33)
		So(cfg.GetFloat64("float64Val"), ShouldBeBetween, 64.63, 64.65)
		So(cfg.GetString("stringVal"), ShouldEqual, "world")
		So(cfg.GetDuration("durationVal"), ShouldEqual, 6*time.Second)
		So(cfg.GetTime("timeVal"), ShouldEqual, time.Unix(1604930127, 0))
		So(cfg.GetIP("ipVal"), ShouldEqual, net.ParseIP("106.11.34.43"))

		So(cfg.GetBool("x"), ShouldBeFalse)
		So(cfg.GetInt("x"), ShouldEqual, 0)
		So(cfg.GetInt8("x"), ShouldEqual, 0)
		So(cfg.GetInt16("x"), ShouldEqual, 0)
		So(cfg.GetInt32("x"), ShouldEqual, 0)
		So(cfg.GetInt64("x"), ShouldEqual, 0)
		So(cfg.GetUint("x"), ShouldEqual, 0)
		So(cfg.GetUint8("x"), ShouldEqual, 0)
		So(cfg.GetUint16("x"), ShouldEqual, 0)
		So(cfg.GetUint32("x"), ShouldEqual, 0)
		So(cfg.GetUint64("x"), ShouldEqual, 0)
		So(cfg.GetFloat32("x"), ShouldEqual, 0.0)
		So(cfg.GetFloat64("x"), ShouldEqual, 0.0)
		So(cfg.GetString("x"), ShouldEqual, "")
		So(cfg.GetDuration("x"), ShouldEqual, 0)
		So(cfg.GetTime("x"), ShouldEqual, time.Time{})
		So(cfg.GetIP("x"), ShouldEqual, net.IP{})

		DeleteTestFile()
	})
}

func TestConfigGetP(t *testing.T) {
	Convey("TestConfigGetE", t, func() {
		CreateTestFile3()

		cfg, err := NewConfigWithSimpleFile("test.json")
		So(err, ShouldBeNil)

		So(cfg.GetBoolP("boolVal"), ShouldBeTrue)
		So(cfg.GetIntP("intVal"), ShouldEqual, 11)
		So(cfg.GetInt8P("int8Val"), ShouldEqual, 22)
		So(cfg.GetInt16P("int16Val"), ShouldEqual, 33)
		So(cfg.GetInt32P("int32Val"), ShouldEqual, 44)
		So(cfg.GetInt64P("int64Val"), ShouldEqual, 55)
		So(cfg.GetUintP("uintVal"), ShouldEqual, 66)
		So(cfg.GetUint8P("uint8Val"), ShouldEqual, 77)
		So(cfg.GetUint16P("uint16Val"), ShouldEqual, 88)
		So(cfg.GetUint32P("uint32Val"), ShouldEqual, 99)
		So(cfg.GetUint64P("uint64Val"), ShouldEqual, 1010)
		So(cfg.GetFloat32P("float32Val"), ShouldBeBetween, 32.31, 32.33)
		So(cfg.GetFloat64P("float64Val"), ShouldBeBetween, 64.63, 64.65)
		So(cfg.GetStringP("stringVal"), ShouldEqual, "world")
		So(cfg.GetDurationP("durationVal"), ShouldEqual, 6*time.Second)
		So(cfg.GetTimeP("timeVal"), ShouldEqual, time.Unix(1604930127, 0))
		So(cfg.GetIPP("ipVal"), ShouldEqual, net.ParseIP("106.11.34.43"))

		So(func() { cfg.GetBoolP("x") }, ShouldPanic)
		So(func() { cfg.GetIntP("x") }, ShouldPanic)
		So(func() { cfg.GetInt8P("x") }, ShouldPanic)
		So(func() { cfg.GetInt16P("x") }, ShouldPanic)
		So(func() { cfg.GetInt32P("x") }, ShouldPanic)
		So(func() { cfg.GetInt64P("x") }, ShouldPanic)
		So(func() { cfg.GetUintP("x") }, ShouldPanic)
		So(func() { cfg.GetUint8P("x") }, ShouldPanic)
		So(func() { cfg.GetUint16P("x") }, ShouldPanic)
		So(func() { cfg.GetUint32P("x") }, ShouldPanic)
		So(func() { cfg.GetUint64P("x") }, ShouldPanic)
		So(func() { cfg.GetFloat32P("x") }, ShouldPanic)
		So(func() { cfg.GetFloat64P("x") }, ShouldPanic)
		So(func() { cfg.GetStringP("x") }, ShouldPanic)
		So(func() { cfg.GetDurationP("x") }, ShouldPanic)
		So(func() { cfg.GetTimeP("x") }, ShouldPanic)
		So(func() { cfg.GetIPP("x") }, ShouldPanic)

		DeleteTestFile()
	})
}
