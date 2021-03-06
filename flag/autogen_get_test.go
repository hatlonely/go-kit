package flag

import (
	"net"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {
	Convey("TestGet", t, func() {
		flag := NewFlag("test")

		flag.Bool("bool-val", false, "")
		flag.Int("int-val", 1, "")
		flag.Int8("int8-val", 2, "")
		flag.Int16("int16-val", 3, "")
		flag.Int32("int32-val", 4, "")
		flag.Int64("int64-val", 5, "")
		flag.Uint("uint-val", 6, "")
		flag.Uint8("uint8-val", 7, "")
		flag.Uint16("uint16-val", 8, "")
		flag.Uint32("uint32-val", 9, "")
		flag.Uint64("uint64-val", 10, "")
		flag.Float32("float32-val", 3.2, "")
		flag.Float64("float64-val", 6.4, "")
		flag.String("string-val", "hello", "")
		flag.Time("time-val", time.Unix(1604930126, 0), "")
		flag.Duration("duration-val", 3*time.Second, "")
		flag.IP("ip-val", net.ParseIP("106.11.34.42"), "")

		So(flag.ParseArgs(strings.Split("--bool-val "+
			"--int-val 11 --int8-val 22 --int16-val 33 --int32-val 44 --int64-val 55 "+
			"--uint-val 66 --uint8-val 77 --uint16-val 88 --uint32-val 99 --uint64-val 1010 "+
			"--float32-val 32.32 --float64-val 64.64 --string-val world "+
			"--duration-val 6s --time-val 2020-11-09T21:55:27+08:00 --ip-val 106.11.34.43", " ")), ShouldBeNil)

		Convey("Get", func() {
			So(flag.GetBool("bool-val"), ShouldBeTrue)
			So(flag.GetInt("int-val"), ShouldEqual, 11)
			So(flag.GetInt8("int8-val"), ShouldEqual, 22)
			So(flag.GetInt16("int16-val"), ShouldEqual, 33)
			So(flag.GetInt32("int32-val"), ShouldEqual, 44)
			So(flag.GetInt64("int64-val"), ShouldEqual, 55)
			So(flag.GetUint("uint-val"), ShouldEqual, 66)
			So(flag.GetUint8("uint8-val"), ShouldEqual, 77)
			So(flag.GetUint16("uint16-val"), ShouldEqual, 88)
			So(flag.GetUint32("uint32-val"), ShouldEqual, 99)
			So(flag.GetUint64("uint64-val"), ShouldEqual, 1010)
			So(flag.GetFloat32("float32-val"), ShouldBeBetween, 32.31, 32.33)
			So(flag.GetFloat64("float64-val"), ShouldBeBetween, 64.63, 64.65)
			So(flag.GetString("string-val"), ShouldEqual, "world")
			So(flag.GetDuration("duration-val"), ShouldEqual, 6*time.Second)
			So(flag.GetTime("time-val"), ShouldEqual, time.Unix(1604930127, 0))
			So(flag.GetIP("ip-val"), ShouldEqual, net.ParseIP("106.11.34.43"))

			So(flag.GetBool("x"), ShouldBeFalse)
			So(flag.GetInt("x"), ShouldEqual, 0)
			So(flag.GetInt8("x"), ShouldEqual, 0)
			So(flag.GetInt16("x"), ShouldEqual, 0)
			So(flag.GetInt32("x"), ShouldEqual, 0)
			So(flag.GetInt64("x"), ShouldEqual, 0)
			So(flag.GetUint("x"), ShouldEqual, 0)
			So(flag.GetUint8("x"), ShouldEqual, 0)
			So(flag.GetUint16("x"), ShouldEqual, 0)
			So(flag.GetUint32("x"), ShouldEqual, 0)
			So(flag.GetUint64("x"), ShouldEqual, 0)
			So(flag.GetFloat32("x"), ShouldEqual, 0.0)
			So(flag.GetFloat64("x"), ShouldEqual, 0.0)
			So(flag.GetString("x"), ShouldEqual, "")
			So(flag.GetDuration("x"), ShouldEqual, 0)
			So(flag.GetTime("x"), ShouldEqual, time.Time{})
			So(flag.GetIP("x"), ShouldEqual, net.IP{})
		})

		Convey("GetP", func() {
			So(flag.GetBoolP("bool-val"), ShouldBeTrue)
			So(flag.GetIntP("int-val"), ShouldEqual, 11)
			So(flag.GetInt8P("int8-val"), ShouldEqual, 22)
			So(flag.GetInt16P("int16-val"), ShouldEqual, 33)
			So(flag.GetInt32P("int32-val"), ShouldEqual, 44)
			So(flag.GetInt64P("int64-val"), ShouldEqual, 55)
			So(flag.GetUintP("uint-val"), ShouldEqual, 66)
			So(flag.GetUint8P("uint8-val"), ShouldEqual, 77)
			So(flag.GetUint16P("uint16-val"), ShouldEqual, 88)
			So(flag.GetUint32P("uint32-val"), ShouldEqual, 99)
			So(flag.GetUint64P("uint64-val"), ShouldEqual, 1010)
			So(flag.GetFloat32P("float32-val"), ShouldBeBetween, 32.31, 32.33)
			So(flag.GetFloat64P("float64-val"), ShouldBeBetween, 64.63, 64.65)
			So(flag.GetStringP("string-val"), ShouldEqual, "world")
			So(flag.GetDurationP("duration-val"), ShouldEqual, 6*time.Second)
			So(flag.GetTimeP("time-val"), ShouldEqual, time.Unix(1604930127, 0))
			So(flag.GetIPP("ip-val"), ShouldEqual, net.ParseIP("106.11.34.43"))

			So(func() { flag.GetBoolP("x") }, ShouldPanic)
			So(func() { flag.GetIntP("x") }, ShouldPanic)
			So(func() { flag.GetInt8P("x") }, ShouldPanic)
			So(func() { flag.GetInt16P("x") }, ShouldPanic)
			So(func() { flag.GetInt32P("x") }, ShouldPanic)
			So(func() { flag.GetInt64P("x") }, ShouldPanic)
			So(func() { flag.GetUintP("x") }, ShouldPanic)
			So(func() { flag.GetUint8P("x") }, ShouldPanic)
			So(func() { flag.GetUint16P("x") }, ShouldPanic)
			So(func() { flag.GetUint32P("x") }, ShouldPanic)
			So(func() { flag.GetUint64P("x") }, ShouldPanic)
			So(func() { flag.GetFloat32P("x") }, ShouldPanic)
			So(func() { flag.GetFloat64P("x") }, ShouldPanic)
			So(func() { flag.GetStringP("x") }, ShouldPanic)
			So(func() { flag.GetDurationP("x") }, ShouldPanic)
			So(func() { flag.GetTimeP("x") }, ShouldPanic)
			So(func() { flag.GetIPP("x") }, ShouldPanic)
		})

		Convey("GetD", func() {
			So(flag.GetBoolD("bool-val", false), ShouldBeTrue)
			So(flag.GetIntD("int-val", 1), ShouldEqual, 11)
			So(flag.GetInt8D("int8-val", 2), ShouldEqual, 22)
			So(flag.GetInt16D("int16-val", 3), ShouldEqual, 33)
			So(flag.GetInt32D("int32-val", 4), ShouldEqual, 44)
			So(flag.GetInt64D("int64-val", 5), ShouldEqual, 55)
			So(flag.GetUintD("uint-val", 6), ShouldEqual, 66)
			So(flag.GetUint8D("uint8-val", 7), ShouldEqual, 77)
			So(flag.GetUint16D("uint16-val", 8), ShouldEqual, 88)
			So(flag.GetUint32D("uint32-val", 9), ShouldEqual, 99)
			So(flag.GetUint64D("uint64-val", 10), ShouldEqual, 1010)
			So(flag.GetFloat32D("float32-val", 3.2), ShouldBeBetween, 32.31, 32.33)
			So(flag.GetFloat64D("float64-val", 6.4), ShouldBeBetween, 64.63, 64.65)
			So(flag.GetStringD("string-val", "hello"), ShouldEqual, "world")
			So(flag.GetDurationD("duration-val", 3*time.Second), ShouldEqual, 6*time.Second)
			So(flag.GetTimeD("time-val", time.Unix(1604930126, 0)), ShouldEqual, time.Unix(1604930127, 0))
			So(flag.GetIPD("ip-val", net.ParseIP("106.11.34.42")), ShouldEqual, net.ParseIP("106.11.34.43"))

			So(flag.GetBoolD("x", false), ShouldBeFalse)
			So(flag.GetIntD("x", 1), ShouldEqual, 1)
			So(flag.GetInt8D("x", 2), ShouldEqual, 2)
			So(flag.GetInt16D("x", 3), ShouldEqual, 3)
			So(flag.GetInt32D("x", 4), ShouldEqual, 4)
			So(flag.GetInt64D("x", 5), ShouldEqual, 5)
			So(flag.GetUintD("x", 6), ShouldEqual, 6)
			So(flag.GetUint8D("x", 7), ShouldEqual, 7)
			So(flag.GetUint16D("x", 8), ShouldEqual, 8)
			So(flag.GetUint32D("x", 9), ShouldEqual, 9)
			So(flag.GetUint64D("x", 10), ShouldEqual, 10)
			So(flag.GetFloat32D("x", 3.2), ShouldBeBetween, 3.1, 3.3)
			So(flag.GetFloat64D("x", 6.4), ShouldBeBetween, 6.3, 6.5)
			So(flag.GetStringD("x", "hello"), ShouldEqual, "hello")
			So(flag.GetDurationD("x", 3*time.Second), ShouldEqual, 3*time.Second)
			So(flag.GetTimeD("x", time.Unix(1604930126, 0)), ShouldEqual, time.Unix(1604930126, 0))
			So(flag.GetIPD("x", net.ParseIP("106.11.34.42")), ShouldEqual, net.ParseIP("106.11.34.42"))
		})
	})
}

func TestGetSlice(t *testing.T) {
	Convey("TestGet", t, func() {
		flag := NewFlag("test")

		flag.Bool("bool-val", false, "")
		flag.Int("int-val", 1, "")
		flag.Int8("int8-val", 2, "")
		flag.Int16("int16-val", 3, "")
		flag.Int32("int32-val", 4, "")
		flag.Int64("int64-val", 5, "")
		flag.Uint("uint-val", 6, "")
		flag.Uint8("uint8-val", 7, "")
		flag.Uint16("uint16-val", 8, "")
		flag.Uint32("uint32-val", 9, "")
		flag.Uint64("uint64-val", 10, "")
		flag.Float32("float32-val", 3.2, "")
		flag.Float64("float64-val", 6.4, "")
		flag.String("string-val", "hello", "")
		flag.Time("time-val", time.Unix(1604930126, 0), "")
		flag.Duration("duration-val", 3*time.Second, "")
		flag.IP("ip-val", net.ParseIP("106.11.34.42"), "")

		So(flag.ParseArgs(strings.Split("--bool-slice-val true,false,true "+
			"--int-slice-val 11,22,33 --int8-slice-val 11,22,33 --int16-slice-val 11,22,33 --int32-slice-val 11,22,33 --int64-slice-val 11,22,33 "+
			"--uint-slice-val 11,22,33 --uint8-slice-val 11,22,33 --uint16-slice-val 11,22,33 --uint32-slice-val 11,22,33 --uint64-slice-val 11,22,33 "+
			"--float32-slice-val 11.11,22.22,33.33 --float64-slice-val 11.11,22.22,33.33 --string-slice-val 11,22,33 "+
			"--duration-slice-val 4s,5s,6s --ip-slice-val 192.168.0.3,192.168.0.4 "+
			"--time-slice-val 1970-01-01T08:00:04+0800,1970-01-01T08:00:05+08:00,1970-01-01T08:00:06+08:00", " ")), ShouldBeNil)

		Convey("Get", func() {
			So(flag.GetBoolSlice("bool-slice-val"), ShouldResemble, []bool{true, false, true})
			So(flag.GetIntSlice("int-slice-val"), ShouldResemble, []int{11, 22, 33})
			So(flag.GetInt8Slice("int8-slice-val"), ShouldResemble, []int8{11, 22, 33})
			So(flag.GetInt16Slice("int16-slice-val"), ShouldResemble, []int16{11, 22, 33})
			So(flag.GetInt32Slice("int32-slice-val"), ShouldResemble, []int32{11, 22, 33})
			So(flag.GetInt64Slice("int64-slice-val"), ShouldResemble, []int64{11, 22, 33})
			So(flag.GetUintSlice("uint-slice-val"), ShouldResemble, []uint{11, 22, 33})
			So(flag.GetUint8Slice("uint8-slice-val"), ShouldResemble, []uint8{11, 22, 33})
			So(flag.GetUint16Slice("uint16-slice-val"), ShouldResemble, []uint16{11, 22, 33})
			So(flag.GetUint32Slice("uint32-slice-val"), ShouldResemble, []uint32{11, 22, 33})
			So(flag.GetUint64Slice("uint64-slice-val"), ShouldResemble, []uint64{11, 22, 33})
			So(flag.GetFloat32Slice("float32-slice-val"), ShouldResemble, []float32{11.11, 22.22, 33.33})
			So(flag.GetFloat64Slice("float64-slice-val"), ShouldResemble, []float64{11.11, 22.22, 33.33})
			So(flag.GetStringSlice("string-slice-val"), ShouldResemble, []string{"11", "22", "33"})
			So(flag.GetDurationSlice("duration-slice-val"), ShouldResemble, []time.Duration{4 * time.Second, 5 * time.Second, 6 * time.Second})
			So(flag.GetTimeSlice("time-slice-val"), ShouldResemble, []time.Time{time.Unix(4, 0), time.Unix(5, 0), time.Unix(6, 0)})
			So(flag.GetIPSlice("ip-slice-val"), ShouldResemble, []net.IP{net.ParseIP("192.168.0.3"), net.ParseIP("192.168.0.4")})

			So(flag.GetBoolSlice("x"), ShouldBeEmpty)
			So(flag.GetIntSlice("x"), ShouldBeEmpty)
			So(flag.GetInt8Slice("x"), ShouldBeEmpty)
			So(flag.GetInt16Slice("x"), ShouldBeEmpty)
			So(flag.GetInt32Slice("x"), ShouldBeEmpty)
			So(flag.GetInt64Slice("x"), ShouldBeEmpty)
			So(flag.GetUintSlice("x"), ShouldBeEmpty)
			So(flag.GetUint8Slice("x"), ShouldBeEmpty)
			So(flag.GetUint16Slice("x"), ShouldBeEmpty)
			So(flag.GetUint32Slice("x"), ShouldBeEmpty)
			So(flag.GetUint64Slice("x"), ShouldBeEmpty)
			So(flag.GetFloat32Slice("x"), ShouldBeEmpty)
			So(flag.GetFloat64Slice("x"), ShouldBeEmpty)
			So(flag.GetStringSlice("x"), ShouldBeEmpty)
			So(flag.GetDurationSlice("x"), ShouldBeEmpty)
			So(flag.GetTimeSlice("x"), ShouldBeEmpty)
			So(flag.GetIPSlice("x"), ShouldBeEmpty)
		})

		Convey("GetP", func() {
			So(flag.GetBoolSliceP("bool-slice-val"), ShouldResemble, []bool{true, false, true})
			So(flag.GetIntSliceP("int-slice-val"), ShouldResemble, []int{11, 22, 33})
			So(flag.GetInt8SliceP("int8-slice-val"), ShouldResemble, []int8{11, 22, 33})
			So(flag.GetInt16SliceP("int16-slice-val"), ShouldResemble, []int16{11, 22, 33})
			So(flag.GetInt32SliceP("int32-slice-val"), ShouldResemble, []int32{11, 22, 33})
			So(flag.GetInt64SliceP("int64-slice-val"), ShouldResemble, []int64{11, 22, 33})
			So(flag.GetUintSliceP("uint-slice-val"), ShouldResemble, []uint{11, 22, 33})
			So(flag.GetUint8SliceP("uint8-slice-val"), ShouldResemble, []uint8{11, 22, 33})
			So(flag.GetUint16SliceP("uint16-slice-val"), ShouldResemble, []uint16{11, 22, 33})
			So(flag.GetUint32SliceP("uint32-slice-val"), ShouldResemble, []uint32{11, 22, 33})
			So(flag.GetUint64SliceP("uint64-slice-val"), ShouldResemble, []uint64{11, 22, 33})
			So(flag.GetFloat32SliceP("float32-slice-val"), ShouldResemble, []float32{11.11, 22.22, 33.33})
			So(flag.GetFloat64SliceP("float64-slice-val"), ShouldResemble, []float64{11.11, 22.22, 33.33})
			So(flag.GetStringSliceP("string-slice-val"), ShouldResemble, []string{"11", "22", "33"})
			So(flag.GetDurationSliceP("duration-slice-val"), ShouldResemble, []time.Duration{4 * time.Second, 5 * time.Second, 6 * time.Second})
			So(flag.GetTimeSliceP("time-slice-val"), ShouldResemble, []time.Time{time.Unix(4, 0), time.Unix(5, 0), time.Unix(6, 0)})
			So(flag.GetIPSliceP("ip-slice-val"), ShouldResemble, []net.IP{net.ParseIP("192.168.0.3"), net.ParseIP("192.168.0.4")})

			So(func() { flag.GetBoolSliceP("x") }, ShouldPanic)
			So(func() { flag.GetIntSliceP("x") }, ShouldPanic)
			So(func() { flag.GetInt8SliceP("x") }, ShouldPanic)
			So(func() { flag.GetInt16SliceP("x") }, ShouldPanic)
			So(func() { flag.GetInt32SliceP("x") }, ShouldPanic)
			So(func() { flag.GetInt64SliceP("x") }, ShouldPanic)
			So(func() { flag.GetUintSliceP("x") }, ShouldPanic)
			So(func() { flag.GetUint8SliceP("x") }, ShouldPanic)
			So(func() { flag.GetUint16SliceP("x") }, ShouldPanic)
			So(func() { flag.GetUint32SliceP("x") }, ShouldPanic)
			So(func() { flag.GetUint64SliceP("x") }, ShouldPanic)
			So(func() { flag.GetFloat32SliceP("x") }, ShouldPanic)
			So(func() { flag.GetFloat64SliceP("x") }, ShouldPanic)
			So(func() { flag.GetStringSliceP("x") }, ShouldPanic)
			So(func() { flag.GetDurationSliceP("x") }, ShouldPanic)
			So(func() { flag.GetTimeSliceP("x") }, ShouldPanic)
			So(func() { flag.GetIPSliceP("x") }, ShouldPanic)
		})

		Convey("GetD", func() {
			So(flag.GetBoolSliceD("bool-slice-val", []bool{true, true, false}), ShouldResemble, []bool{true, false, true})
			So(flag.GetIntSliceD("int-slice-val", []int{1, 2, 3}), ShouldResemble, []int{11, 22, 33})
			So(flag.GetInt8SliceD("int8-slice-val", []int8{1, 2, 3}), ShouldResemble, []int8{11, 22, 33})
			So(flag.GetInt16SliceD("int16-slice-val", []int16{1, 2, 3}), ShouldResemble, []int16{11, 22, 33})
			So(flag.GetInt32SliceD("int32-slice-val", []int32{1, 2, 3}), ShouldResemble, []int32{11, 22, 33})
			So(flag.GetInt64SliceD("int64-slice-val", []int64{1, 2, 3}), ShouldResemble, []int64{11, 22, 33})
			So(flag.GetUintSliceD("uint-slice-val", []uint{1, 2, 3}), ShouldResemble, []uint{11, 22, 33})
			So(flag.GetUint8SliceD("uint8-slice-val", []uint8{1, 2, 3}), ShouldResemble, []uint8{11, 22, 33})
			So(flag.GetUint16SliceD("uint16-slice-val", []uint16{1, 2, 3}), ShouldResemble, []uint16{11, 22, 33})
			So(flag.GetUint32SliceD("uint32-slice-val", []uint32{1, 2, 3}), ShouldResemble, []uint32{11, 22, 33})
			So(flag.GetUint64SliceD("uint64-slice-val", []uint64{1, 2, 3}), ShouldResemble, []uint64{11, 22, 33})
			So(flag.GetFloat32SliceD("float32-slice-val", []float32{1.1, 2.2, 3.3}), ShouldResemble, []float32{11.11, 22.22, 33.33})
			So(flag.GetFloat64SliceD("float64-slice-val", []float64{1.1, 2.2, 3.3}), ShouldResemble, []float64{11.11, 22.22, 33.33})
			So(flag.GetStringSliceD("string-slice-val", []string{"1", "2", "3"}), ShouldResemble, []string{"11", "22", "33"})
			So(flag.GetDurationSliceD("duration-slice-val", []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second}), ShouldResemble, []time.Duration{4 * time.Second, 5 * time.Second, 6 * time.Second})
			So(flag.GetTimeSliceD("time-slice-val", []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)}), ShouldResemble, []time.Time{time.Unix(4, 0), time.Unix(5, 0), time.Unix(6, 0)})
			So(flag.GetIPSliceD("ip-slice-val", []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")}), ShouldResemble, []net.IP{net.ParseIP("192.168.0.3"), net.ParseIP("192.168.0.4")})

			So(flag.GetBoolSliceD("x", []bool{true, true, false}), ShouldResemble, []bool{true, true, false})
			So(flag.GetIntSliceD("x", []int{1, 2, 3}), ShouldResemble, []int{1, 2, 3})
			So(flag.GetInt8SliceD("x", []int8{1, 2, 3}), ShouldResemble, []int8{1, 2, 3})
			So(flag.GetInt16SliceD("x", []int16{1, 2, 3}), ShouldResemble, []int16{1, 2, 3})
			So(flag.GetInt32SliceD("x", []int32{1, 2, 3}), ShouldResemble, []int32{1, 2, 3})
			So(flag.GetInt64SliceD("x", []int64{1, 2, 3}), ShouldResemble, []int64{1, 2, 3})
			So(flag.GetUintSliceD("x", []uint{1, 2, 3}), ShouldResemble, []uint{1, 2, 3})
			So(flag.GetUint8SliceD("x", []uint8{1, 2, 3}), ShouldResemble, []uint8{1, 2, 3})
			So(flag.GetUint16SliceD("x", []uint16{1, 2, 3}), ShouldResemble, []uint16{1, 2, 3})
			So(flag.GetUint32SliceD("x", []uint32{1, 2, 3}), ShouldResemble, []uint32{1, 2, 3})
			So(flag.GetUint64SliceD("x", []uint64{1, 2, 3}), ShouldResemble, []uint64{1, 2, 3})
			So(flag.GetFloat32SliceD("x", []float32{1.1, 2.2, 3.3}), ShouldResemble, []float32{1.1, 2.2, 3.3})
			So(flag.GetFloat64SliceD("x", []float64{1.1, 2.2, 3.3}), ShouldResemble, []float64{1.1, 2.2, 3.3})
			So(flag.GetStringSliceD("x", []string{"1", "2", "3"}), ShouldResemble, []string{"1", "2", "3"})
			So(flag.GetDurationSliceD("x", []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second}), ShouldResemble, []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second})
			So(flag.GetTimeSliceD("x", []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)}), ShouldResemble, []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)})
			So(flag.GetIPSliceD("x", []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")}), ShouldResemble, []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")})
		})
	})
}
