package config

import (
	"net"
	"os"
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

func CreateTestFile2() {
	fp, _ := os.Create("test.json")
	_, _ = fp.WriteString(`{
  "boolVal": false,
  "intVal": 1,
  "int8Val": 2,
  "int16Val": 3,
  "int32Val": 4,
  "int64Val": 5,
  "uintVal": 6,
  "uint8Val": 7,
  "uint16Val": 8,
  "uint32Val": 9,
  "uint64Val": 10,
  "float32Val": 3.2,
  "float64Val": 6.4,
}`)
	_ = fp.Close()
}

func CreateTestFile3() {
	fp, _ := os.Create("test.json")
	_, _ = fp.WriteString(`{
  "boolVal": true,
  "intVal": 11,
  "int8Val": 22,
  "int16Val": 33,
  "int32Val": 44,
  "int64Val": 55,
  "uintVal": 66,
  "uint8Val": 77,
  "uint16Val": 88,
  "uint32Val": 99,
  "uint64Val": 1010,
  "float32Val": 32.32,
  "float64Val": 64.64,
}`)
	_ = fp.Close()
}

func TestIntVar(t *testing.T) {
	Convey("TestIntVar", t, func(c C) {
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
		cfg, err := NewConfigWithSimpleFile("test.json")
		So(err, ShouldBeNil)
		So(cfg.Watch(), ShouldBeNil)
		defer cfg.Stop()

		cfg.BoolVar("boolVal", &boolVal, OnSucc(func(cfg *Config) {
			c.So(cfg.GetBool(""), ShouldEqual, true)
		}))
		cfg.IntVar("intVal", &intVal, OnSucc(func(cfg *Config) {
			c.So(cfg.GetInt(""), ShouldEqual, 11)
		}))
		cfg.Int8Var("int8Val", &int8Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetInt8(""), ShouldEqual, 22)
		}))
		cfg.Int16Var("int16Val", &int16Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetInt16(""), ShouldEqual, 33)
		}))
		cfg.Int32Var("int32Val", &int32Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetInt32(""), ShouldEqual, 44)
		}))
		cfg.Int64Var("int64Val", &int64Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetInt64(""), ShouldEqual, 55)
		}))
		cfg.UintVar("uintVal", &uintVal, OnSucc(func(cfg *Config) {
			c.So(cfg.GetUint(""), ShouldEqual, 66)
		}))
		cfg.Uint8Var("uint8Val", &uint8Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetUint8(""), ShouldEqual, 77)
		}))
		cfg.Uint16Var("uint16Val", &uint16Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetUint16(""), ShouldEqual, 88)
		}))
		cfg.Uint32Var("uint32Val", &uint32Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetUint32(""), ShouldEqual, 99)
		}))
		cfg.Uint64Var("uint64Val", &uint64Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetUint64(""), ShouldEqual, 1010)
		}))
		cfg.Float32Var("float32Val", &float32Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetFloat32(""), ShouldBeBetween, 32.31, 32.33)
		}))
		cfg.Float64Var("float64Val", &float64Val, OnSucc(func(cfg *Config) {
			c.So(cfg.GetFloat64(""), ShouldBeBetween, 64.63, 64.65)
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

		DeleteTestFile()
	})
}
