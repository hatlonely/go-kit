package config

import (
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGlobalGet(t *testing.T) {
	Convey("TestConfigGet", t, func() {
		CreateTestFile3()

		err := InitWithSimpleFile("test.json")
		So(err, ShouldBeNil)

		So(GetBool("boolVal"), ShouldBeTrue)
		So(GetInt("intVal"), ShouldEqual, 11)
		So(GetInt8("int8Val"), ShouldEqual, 22)
		So(GetInt16("int16Val"), ShouldEqual, 33)
		So(GetInt32("int32Val"), ShouldEqual, 44)
		So(GetInt64("int64Val"), ShouldEqual, 55)
		So(GetUint("uintVal"), ShouldEqual, 66)
		So(GetUint8("uint8Val"), ShouldEqual, 77)
		So(GetUint16("uint16Val"), ShouldEqual, 88)
		So(GetUint32("uint32Val"), ShouldEqual, 99)
		So(GetUint64("uint64Val"), ShouldEqual, 1010)
		So(GetFloat32("float32Val"), ShouldBeBetween, 32.31, 32.33)
		So(GetFloat64("float64Val"), ShouldBeBetween, 64.63, 64.65)
		So(GetString("stringVal"), ShouldEqual, "world")
		So(GetDuration("durationVal"), ShouldEqual, 6*time.Second)
		So(GetTime("timeVal"), ShouldEqual, time.Unix(1604930127, 0))
		So(GetIP("ipVal"), ShouldEqual, net.ParseIP("106.11.34.43"))

		So(GetBool("x"), ShouldBeFalse)
		So(GetInt("x"), ShouldEqual, 0)
		So(GetInt8("x"), ShouldEqual, 0)
		So(GetInt16("x"), ShouldEqual, 0)
		So(GetInt32("x"), ShouldEqual, 0)
		So(GetInt64("x"), ShouldEqual, 0)
		So(GetUint("x"), ShouldEqual, 0)
		So(GetUint8("x"), ShouldEqual, 0)
		So(GetUint16("x"), ShouldEqual, 0)
		So(GetUint32("x"), ShouldEqual, 0)
		So(GetUint64("x"), ShouldEqual, 0)
		So(GetFloat32("x"), ShouldEqual, 0.0)
		So(GetFloat64("x"), ShouldEqual, 0.0)
		So(GetString("x"), ShouldEqual, "")
		So(GetDuration("x"), ShouldEqual, 0)
		So(GetTime("x"), ShouldEqual, time.Time{})
		So(GetIP("x"), ShouldEqual, net.IP{})

		So(GetBoolP("boolVal"), ShouldBeTrue)
		So(GetIntP("intVal"), ShouldEqual, 11)
		So(GetInt8P("int8Val"), ShouldEqual, 22)
		So(GetInt16P("int16Val"), ShouldEqual, 33)
		So(GetInt32P("int32Val"), ShouldEqual, 44)
		So(GetInt64P("int64Val"), ShouldEqual, 55)
		So(GetUintP("uintVal"), ShouldEqual, 66)
		So(GetUint8P("uint8Val"), ShouldEqual, 77)
		So(GetUint16P("uint16Val"), ShouldEqual, 88)
		So(GetUint32P("uint32Val"), ShouldEqual, 99)
		So(GetUint64P("uint64Val"), ShouldEqual, 1010)
		So(GetFloat32P("float32Val"), ShouldBeBetween, 32.31, 32.33)
		So(GetFloat64P("float64Val"), ShouldBeBetween, 64.63, 64.65)
		So(GetStringP("stringVal"), ShouldEqual, "world")
		So(GetDurationP("durationVal"), ShouldEqual, 6*time.Second)
		So(GetTimeP("timeVal"), ShouldEqual, time.Unix(1604930127, 0))
		So(GetIPP("ipVal"), ShouldEqual, net.ParseIP("106.11.34.43"))

		So(func() { GetBoolP("x") }, ShouldPanic)
		So(func() { GetIntP("x") }, ShouldPanic)
		So(func() { GetInt8P("x") }, ShouldPanic)
		So(func() { GetInt16P("x") }, ShouldPanic)
		So(func() { GetInt32P("x") }, ShouldPanic)
		So(func() { GetInt64P("x") }, ShouldPanic)
		So(func() { GetUintP("x") }, ShouldPanic)
		So(func() { GetUint8P("x") }, ShouldPanic)
		So(func() { GetUint16P("x") }, ShouldPanic)
		So(func() { GetUint32P("x") }, ShouldPanic)
		So(func() { GetUint64P("x") }, ShouldPanic)
		So(func() { GetFloat32P("x") }, ShouldPanic)
		So(func() { GetFloat64P("x") }, ShouldPanic)
		So(func() { GetStringP("x") }, ShouldPanic)
		So(func() { GetDurationP("x") }, ShouldPanic)
		So(func() { GetTimeP("x") }, ShouldPanic)
		So(func() { GetIPP("x") }, ShouldPanic)

		So(GetBoolD("boolVal", false), ShouldBeTrue)
		So(GetIntD("intVal", 1), ShouldEqual, 11)
		So(GetInt8D("int8Val", 2), ShouldEqual, 22)
		So(GetInt16D("int16Val", 3), ShouldEqual, 33)
		So(GetInt32D("int32Val", 4), ShouldEqual, 44)
		So(GetInt64D("int64Val", 5), ShouldEqual, 55)
		So(GetUintD("uintVal", 6), ShouldEqual, 66)
		So(GetUint8D("uint8Val", 7), ShouldEqual, 77)
		So(GetUint16D("uint16Val", 8), ShouldEqual, 88)
		So(GetUint32D("uint32Val", 9), ShouldEqual, 99)
		So(GetUint64D("uint64Val", 10), ShouldEqual, 1010)
		So(GetFloat32D("float32Val", 3.2), ShouldBeBetween, 32.31, 32.33)
		So(GetFloat64D("float64Val", 6.4), ShouldBeBetween, 64.63, 64.65)
		So(GetStringD("stringVal", "hello"), ShouldEqual, "world")
		So(GetDurationD("durationVal", 3*time.Second), ShouldEqual, 6*time.Second)
		So(GetTimeD("timeVal", time.Unix(1604930126, 0)), ShouldEqual, time.Unix(1604930127, 0))
		So(GetIPD("ipVal", net.ParseIP("106.11.34.42")), ShouldEqual, net.ParseIP("106.11.34.43"))

		So(GetBoolD("x", false), ShouldBeFalse)
		So(GetIntD("x", 1), ShouldEqual, 1)
		So(GetInt8D("x", 2), ShouldEqual, 2)
		So(GetInt16D("x", 3), ShouldEqual, 3)
		So(GetInt32D("x", 4), ShouldEqual, 4)
		So(GetInt64D("x", 5), ShouldEqual, 5)
		So(GetUintD("x", 6), ShouldEqual, 6)
		So(GetUint8D("x", 7), ShouldEqual, 7)
		So(GetUint16D("x", 8), ShouldEqual, 8)
		So(GetUint32D("x", 9), ShouldEqual, 9)
		So(GetUint64D("x", 10), ShouldEqual, 10)
		So(GetFloat32D("x", 3.2), ShouldBeBetween, 3.1, 3.3)
		So(GetFloat64D("x", 6.4), ShouldBeBetween, 6.3, 6.5)
		So(GetStringD("x", "hello"), ShouldEqual, "hello")
		So(GetDurationD("x", 3*time.Second), ShouldEqual, 3*time.Second)
		So(GetTimeD("x", time.Unix(1604930126, 0)), ShouldEqual, time.Unix(1604930126, 0))
		So(GetIPD("x", net.ParseIP("106.11.34.42")), ShouldEqual, net.ParseIP("106.11.34.42"))

		{
			val, err := GetBoolE("boolVal")
			So(err, ShouldBeNil)
			So(val, ShouldBeTrue)
		}
		{
			val, err := GetIntE("intVal")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 11)
		}
		{
			val, err := GetInt8E("int8Val")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 22)
		}
		{
			val, err := GetInt16E("int16Val")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 33)
		}
		{
			val, err := GetInt32E("int32Val")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 44)
		}
		{
			val, err := GetInt64E("int64Val")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 55)
		}
		{
			val, err := GetUintE("uintVal")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 66)
		}
		{
			val, err := GetUint8E("uint8Val")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 77)
		}
		{
			val, err := GetUint16E("uint16Val")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 88)
		}
		{
			val, err := GetUint32E("uint32Val")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 99)
		}
		{
			val, err := GetUint64E("uint64Val")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 1010)
		}
		{
			val, err := GetFloat32E("float32Val")
			So(err, ShouldBeNil)
			So(val, ShouldBeBetween, 32.31, 32.33)
		}
		{
			val, err := GetFloat64E("float64Val")
			So(err, ShouldBeNil)
			So(val, ShouldBeBetween, 64.63, 64.65)
		}
		{
			val, err := GetStringE("stringVal")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "world")
		}
		{
			val, err := GetDurationE("durationVal")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 6*time.Second)
		}
		{
			val, err := GetTimeE("timeVal")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, time.Unix(1604930127, 0))
		}
		{
			val, err := GetIPE("ipVal")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, net.ParseIP("106.11.34.43"))
		}

		DeleteTestFile()
	})
}

func TestGlobalBinds(t *testing.T) {
	Convey("TestGlobalBinds", t, func() {
		CreateTestFile2()

		boolVal := Bool("boolVal")
		intVal := Int("intVal")
		int8Val := Int8("int8Val")
		int16Val := Int16("int16Val")
		int32Val := Int32("int32Val")
		int64Val := Int64("int64Val")
		uintVal := Uint("uintVal")
		uint8Val := Uint8("uint8Val")
		uint16Val := Uint16("uint16Val")
		uint32Val := Uint32("uint32Val")
		uint64Val := Uint64("uint64Val")
		float32Val := Float32("float32Val")
		float64Val := Float64("float64Val")
		stringVal := String("stringVal")
		durationVal := Duration("durationVal")
		timeVal := Time("timeVal")
		ipVal := IP("ipVal")

		err := InitWithSimpleFile("test.json")
		So(err, ShouldBeNil)
		So(Watch(), ShouldBeNil)
		defer Stop()

		So(boolVal.Get(), ShouldEqual, false)
		So(intVal.Get(), ShouldEqual, 1)
		So(int8Val.Get(), ShouldEqual, 2)
		So(int16Val.Get(), ShouldEqual, 3)
		So(int32Val.Get(), ShouldEqual, 4)
		So(int64Val.Get(), ShouldEqual, 5)
		So(uintVal.Get(), ShouldEqual, 6)
		So(uint8Val.Get(), ShouldEqual, 7)
		So(uint16Val.Get(), ShouldEqual, 8)
		So(uint32Val.Get(), ShouldEqual, 9)
		So(uint64Val.Get(), ShouldEqual, 10)
		So(float32Val.Get(), ShouldBeBetween, 3.1, 3.3)
		So(float64Val.Get(), ShouldBeBetween, 6.3, 6.5)
		So(stringVal.Get(), ShouldEqual, "hello")
		So(durationVal.Get(), ShouldEqual, 3*time.Second)
		So(timeVal.Get(), ShouldEqual, time.Unix(1604930126, 0))
		So(ipVal.Get(), ShouldEqual, net.ParseIP("106.11.34.42"))

		CreateTestFile3()
		time.Sleep(time.Second)

		So(boolVal.Get(), ShouldEqual, true)
		So(intVal.Get(), ShouldEqual, 11)
		So(int8Val.Get(), ShouldEqual, 22)
		So(int16Val.Get(), ShouldEqual, 33)
		So(int32Val.Get(), ShouldEqual, 44)
		So(int64Val.Get(), ShouldEqual, 55)
		So(uintVal.Get(), ShouldEqual, 66)
		So(uint8Val.Get(), ShouldEqual, 77)
		So(uint16Val.Get(), ShouldEqual, 88)
		So(uint32Val.Get(), ShouldEqual, 99)
		So(uint64Val.Get(), ShouldEqual, 1010)
		So(float32Val.Get(), ShouldBeBetween, 32.31, 32.33)
		So(float64Val.Get(), ShouldBeBetween, 64.63, 64.65)
		So(stringVal.Get(), ShouldEqual, "world")
		So(durationVal.Get(), ShouldEqual, 6*time.Second)
		So(timeVal.Get(), ShouldEqual, time.Unix(1604930127, 0))
		So(ipVal.Get(), ShouldEqual, net.ParseIP("106.11.34.43"))

		DeleteTestFile()
	})
}

func TestGlobalBindVars(t *testing.T) {
	Convey("TestBindVars", t, func(c C) {
		CreateTestFile2()

		var boolVal AtomicBool
		var intVal AtomicInt
		var int8Val AtomicInt8
		var int16Val AtomicInt16
		var int32Val AtomicInt32
		var int64Val AtomicInt64
		var uintVal AtomicUint
		var uint8Val AtomicUint8
		var uint16Val AtomicUint16
		var uint32Val AtomicUint32
		var uint64Val AtomicUint64
		var float32Val AtomicFloat32
		var float64Val AtomicFloat64
		var stringVal AtomicString
		var durationVal AtomicDuration
		var timeVal AtomicTime
		var ipVal AtomicIP
		err := InitWithSimpleFile("test.json")
		So(err, ShouldBeNil)
		So(Watch(), ShouldBeNil)
		defer Stop()

		BoolVar("boolVal", &boolVal, OnSucc(func(cfg *Config) {
			c.So(cfg.GetBool(""), ShouldEqual, true)
		}))
		IntVar("intVal", &intVal, OnSucc(func(cfg *Config) {
			c.So(cfg.GetInt(""), ShouldEqual, 11)
		}))
		Int8Var("int8Val", &int8Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetInt8(""), ShouldEqual, 22)
		}))
		Int16Var("int16Val", &int16Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetInt16(""), ShouldEqual, 33)
		}))
		Int32Var("int32Val", &int32Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetInt32(""), ShouldEqual, 44)
		}))
		Int64Var("int64Val", &int64Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetInt64(""), ShouldEqual, 55)
		}))
		UintVar("uintVal", &uintVal, OnSucc(func(cfg *Config) {
			c.So(cfg.GetUint(""), ShouldEqual, 66)
		}))
		Uint8Var("uint8Val", &uint8Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetUint8(""), ShouldEqual, 77)
		}))
		Uint16Var("uint16Val", &uint16Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetUint16(""), ShouldEqual, 88)
		}))
		Uint32Var("uint32Val", &uint32Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetUint32(""), ShouldEqual, 99)
		}))
		Uint64Var("uint64Val", &uint64Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetUint64(""), ShouldEqual, 1010)
		}))
		Float32Var("float32Val", &float32Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetFloat32(""), ShouldBeBetween, 32.31, 32.33)
		}))
		Float64Var("float64Val", &float64Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetFloat64(""), ShouldBeBetween, 64.63, 64.65)
		}))
		StringVar("stringVal", &stringVal, OnSucc(func(cfg *Config) {
			c.So(cfg.GetString(""), ShouldEqual, "world")
		}))
		DurationVar("durationVal", &durationVal, OnSucc(func(cfg *Config) {
			c.So(cfg.GetDuration(""), ShouldEqual, 6*time.Second)
		}))
		TimeVar("timeVal", &timeVal, OnSucc(func(cfg *Config) {
			c.So(cfg.GetTime(""), ShouldEqual, time.Unix(1604930127, 0))
		}))
		IPVar("ipVal", &ipVal, OnSucc(func(cfg *Config) {
			c.So(cfg.GetIP(""), ShouldEqual, net.ParseIP("106.11.34.43"))
		}))

		So(boolVal.Get(), ShouldEqual, false)
		So(intVal.Get(), ShouldEqual, 1)
		So(int8Val.Get(), ShouldEqual, 2)
		So(int16Val.Get(), ShouldEqual, 3)
		So(int32Val.Get(), ShouldEqual, 4)
		So(int64Val.Get(), ShouldEqual, 5)
		So(uintVal.Get(), ShouldEqual, 6)
		So(uint8Val.Get(), ShouldEqual, 7)
		So(uint16Val.Get(), ShouldEqual, 8)
		So(uint32Val.Get(), ShouldEqual, 9)
		So(uint64Val.Get(), ShouldEqual, 10)
		So(float32Val.Get(), ShouldBeBetween, 3.1, 3.3)
		So(float64Val.Get(), ShouldBeBetween, 6.3, 6.5)
		So(stringVal.Get(), ShouldEqual, "hello")
		So(durationVal.Get(), ShouldEqual, 3*time.Second)
		So(timeVal.Get(), ShouldEqual, time.Unix(1604930126, 0))
		So(ipVal.Get(), ShouldEqual, net.ParseIP("106.11.34.42"))

		CreateTestFile3()
		time.Sleep(time.Second)

		So(boolVal.Get(), ShouldEqual, true)
		So(intVal.Get(), ShouldEqual, 11)
		So(int8Val.Get(), ShouldEqual, 22)
		So(int16Val.Get(), ShouldEqual, 33)
		So(int32Val.Get(), ShouldEqual, 44)
		So(int64Val.Get(), ShouldEqual, 55)
		So(uintVal.Get(), ShouldEqual, 66)
		So(uint8Val.Get(), ShouldEqual, 77)
		So(uint16Val.Get(), ShouldEqual, 88)
		So(uint32Val.Get(), ShouldEqual, 99)
		So(uint64Val.Get(), ShouldEqual, 1010)
		So(float32Val.Get(), ShouldBeBetween, 32.31, 32.33)
		So(float64Val.Get(), ShouldBeBetween, 64.63, 64.65)
		So(stringVal.Get(), ShouldEqual, "world")
		So(durationVal.Get(), ShouldEqual, 6*time.Second)
		So(timeVal.Get(), ShouldEqual, time.Unix(1604930127, 0))
		So(ipVal.Get(), ShouldEqual, net.ParseIP("106.11.34.43"))

		DeleteTestFile()
	})
}
