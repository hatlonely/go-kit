package flag

import (
	"net"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGlobalFlagBind(t *testing.T) {
	Convey("TestGlobalFlagBind", t, func() {
		gflag = NewFlag("test")

		boolVal := Bool("bool-val", false, "")
		intVal := Int("int-val", 1, "")
		int8Val := Int8("int8-val", 2, "")
		int16Val := Int16("int16-val", 3, "")
		int32Val := Int32("int32-val", 4, "")
		int64Val := Int64("int64-val", 5, "")
		uintVal := Uint("uint-val", 6, "")
		uint8Val := Uint8("uint8-val", 7, "")
		uint16Val := Uint16("uint16-val", 8, "")
		uint32Val := Uint32("uint32-val", 9, "")
		uint64Val := Uint64("uint64-val", 10, "")
		float32Val := Float32("float32-val", 3.2, "")
		float64Val := Float64("float64-val", 6.4, "")
		stringVal := String("string-val", "hello", "")
		timeVal := Time("time-val", time.Unix(1604930126, 0), "")
		durationVal := Duration("duration-val", 3*time.Second, "")
		ipVal := IP("ip-val", net.ParseIP("106.11.34.42"), "")
		boolSliceVal := BoolSlice("bool-slice-val", []bool{true, true, false}, "")
		intSliceVal := IntSlice("int-slice-val", []int{1, 2, 3}, "")
		int8SliceVal := Int8Slice("int8-slice-val", []int8{1, 2, 3}, "")
		int16SliceVal := Int16Slice("int16-slice-val", []int16{1, 2, 3}, "")
		int32SliceVal := Int32Slice("int32-slice-val", []int32{1, 2, 3}, "")
		int64SliceVal := Int64Slice("int64-slice-val", []int64{1, 2, 3}, "")
		uintSliceVal := UintSlice("uint-slice-val", []uint{1, 2, 3}, "")
		uint8SliceVal := Uint8Slice("uint8-slice-val", []uint8{1, 2, 3}, "")
		uint16SliceVal := Uint16Slice("uint16-slice-val", []uint16{1, 2, 3}, "")
		uint32SliceVal := Uint32Slice("uint32-slice-val", []uint32{1, 2, 3}, "")
		uint64SliceVal := Uint64Slice("uint64-slice-val", []uint64{1, 2, 3}, "")
		float32SliceVal := Float32Slice("float32-slice-val", []float32{1.1, 2.2, 3.3}, "")
		float64SliceVal := Float64Slice("float64-slice-val", []float64{1.1, 2.2, 3.3}, "")
		stringSliceVal := StringSlice("string-slice-val", []string{"1", "2", "3"}, "")
		timeSliceVal := TimeSlice("time-slice-val", []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)}, "")
		durationSliceVal := DurationSlice("duration-slice-val", []time.Duration{time.Second, 2 * time.Second, 3 * time.Second}, "")
		ipSliceVal := IPSlice("ip-slice-val", []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")}, "")

		So(ParseArgs(strings.Split("--bool-val "+
			"--int-val 11 --int8-val 22 --int16-val 33 --int32-val 44 --int64-val 55 "+
			"--uint-val 66 --uint8-val 77 --uint16-val 88 --uint32-val 99 --uint64-val 1010 "+
			"--float32-val 32.32 --float64-val 64.64 --string-val world "+
			"--duration-val 6s --time-val 2020-11-09T21:55:27+08:00 --ip-val 106.11.34.43 "+
			"--bool-slice-val true,false,true "+
			"--int-slice-val 11,22,33 --int8-slice-val 11,22,33 --int16-slice-val 11,22,33 --int32-slice-val 11,22,33 --int64-slice-val 11,22,33 "+
			"--uint-slice-val 11,22,33 --uint8-slice-val 11,22,33 --uint16-slice-val 11,22,33 --uint32-slice-val 11,22,33 --uint64-slice-val 11,22,33 "+
			"--float32-slice-val 11.11,22.22,33.33 --float64-slice-val 11.11,22.22,33.33 --string-slice-val 11,22,33 "+
			"--duration-slice-val 4s,5s,6s --ip-slice-val 192.168.0.3,192.168.0.4 "+
			"--time-slice-val 1970-01-01T08:00:04+0800,1970-01-01T08:00:05+08:00,1970-01-01T08:00:06+08:00", " ")), ShouldBeNil)

		So(*boolVal, ShouldEqual, true)
		So(*intVal, ShouldEqual, 11)
		So(*int8Val, ShouldEqual, 22)
		So(*int16Val, ShouldEqual, 33)
		So(*int32Val, ShouldEqual, 44)
		So(*int64Val, ShouldEqual, 55)
		So(*uintVal, ShouldEqual, 66)
		So(*uint8Val, ShouldEqual, 77)
		So(*uint16Val, ShouldEqual, 88)
		So(*uint32Val, ShouldEqual, 99)
		So(*uint64Val, ShouldEqual, 1010)
		So(*float32Val, ShouldBeBetween, 32.31, 32.33)
		So(*float64Val, ShouldBeBetween, 64.63, 64.65)
		So(*stringVal, ShouldEqual, "world")
		So(*durationVal, ShouldEqual, 6*time.Second)
		So(*timeVal, ShouldEqual, time.Unix(1604930127, 0))
		So(*ipVal, ShouldEqual, net.ParseIP("106.11.34.43"))
		So(*boolSliceVal, ShouldResemble, []bool{true, false, true})
		So(*intSliceVal, ShouldResemble, []int{11, 22, 33})
		So(*int8SliceVal, ShouldResemble, []int8{11, 22, 33})
		So(*int16SliceVal, ShouldResemble, []int16{11, 22, 33})
		So(*int32SliceVal, ShouldResemble, []int32{11, 22, 33})
		So(*int64SliceVal, ShouldResemble, []int64{11, 22, 33})
		So(*uintSliceVal, ShouldResemble, []uint{11, 22, 33})
		So(*uint8SliceVal, ShouldResemble, []uint8{11, 22, 33})
		So(*uint16SliceVal, ShouldResemble, []uint16{11, 22, 33})
		So(*uint32SliceVal, ShouldResemble, []uint32{11, 22, 33})
		So(*uint64SliceVal, ShouldResemble, []uint64{11, 22, 33})
		So(*float32SliceVal, ShouldResemble, []float32{11.11, 22.22, 33.33})
		So(*float64SliceVal, ShouldResemble, []float64{11.11, 22.22, 33.33})
		So(*stringSliceVal, ShouldResemble, []string{"11", "22", "33"})
		So(*timeSliceVal, ShouldResemble, []time.Time{time.Unix(4, 0), time.Unix(5, 0), time.Unix(6, 0)})
		So(*durationSliceVal, ShouldResemble, []time.Duration{4 * time.Second, 5 * time.Second, 6 * time.Second})
		So(*ipSliceVal, ShouldResemble, []net.IP{net.ParseIP("192.168.0.3"), net.ParseIP("192.168.0.4")})
	})
}

func TestGlobalFlagBindVar(t *testing.T) {
	Convey("TestGlobalFlagBindVar", t, func() {
		gflag = NewFlag("test")

		var boolVal bool
		var intVal int
		var int8Val int8
		var int16Val int16
		var int32Val int32
		var int64Val int64
		var uintVal uint
		var uint8Val uint8
		var uint16Val uint16
		var uint32Val uint32
		var uint64Val uint64
		var float32Val float32
		var float64Val float64
		var stringVal string
		var timeVal time.Time
		var durationVal time.Duration
		var ipVal net.IP
		var boolSliceVal []bool
		var intSliceVal []int
		var int8SliceVal []int8
		var int16SliceVal []int16
		var int32SliceVal []int32
		var int64SliceVal []int64
		var uintSliceVal []uint
		var uint8SliceVal []uint8
		var uint16SliceVal []uint16
		var uint32SliceVal []uint32
		var uint64SliceVal []uint64
		var float32SliceVal []float32
		var float64SliceVal []float64
		var stringSliceVal []string
		var timeSliceVal []time.Time
		var durationSliceVal []time.Duration
		var ipSliceVal []net.IP

		BoolVar(&boolVal, "bool-val", false, "")
		IntVar(&intVal, "int-val", 1, "")
		Int8Var(&int8Val, "int8-val", 2, "")
		Int16Var(&int16Val, "int16-val", 3, "")
		Int32Var(&int32Val, "int32-val", 4, "")
		Int64Var(&int64Val, "int64-val", 5, "")
		UintVar(&uintVal, "uint-val", 6, "")
		Uint8Var(&uint8Val, "uint8-val", 7, "")
		Uint16Var(&uint16Val, "uint16-val", 8, "")
		Uint32Var(&uint32Val, "uint32-val", 9, "")
		Uint64Var(&uint64Val, "uint64-val", 10, "")
		Float32Var(&float32Val, "float32-val", 3.2, "")
		Float64Var(&float64Val, "float64-val", 6.4, "")
		StringVar(&stringVal, "string-val", "hello", "")
		TimeVar(&timeVal, "time-val", time.Unix(1604930126, 0), "")
		DurationVar(&durationVal, "duration-val", 3*time.Second, "")
		IPVar(&ipVal, "ip-val", net.ParseIP("106.11.34.42"), "")
		BoolSliceVar(&boolSliceVal, "bool-slice-val", []bool{true, true, false}, "")
		IntSliceVar(&intSliceVal, "int-slice-val", []int{1, 2, 3}, "")
		Int8SliceVar(&int8SliceVal, "int8-slice-val", []int8{1, 2, 3}, "")
		Int16SliceVar(&int16SliceVal, "int16-slice-val", []int16{1, 2, 3}, "")
		Int32SliceVar(&int32SliceVal, "int32-slice-val", []int32{1, 2, 3}, "")
		Int64SliceVar(&int64SliceVal, "int64-slice-val", []int64{1, 2, 3}, "")
		UintSliceVar(&uintSliceVal, "uint-slice-val", []uint{1, 2, 3}, "")
		Uint8SliceVar(&uint8SliceVal, "uint8-slice-val", []uint8{1, 2, 3}, "")
		Uint16SliceVar(&uint16SliceVal, "uint16-slice-val", []uint16{1, 2, 3}, "")
		Uint32SliceVar(&uint32SliceVal, "uint32-slice-val", []uint32{1, 2, 3}, "")
		Uint64SliceVar(&uint64SliceVal, "uint64-slice-val", []uint64{1, 2, 3}, "")
		Float32SliceVar(&float32SliceVal, "float32-slice-val", []float32{1.1, 2.2, 3.3}, "")
		Float64SliceVar(&float64SliceVal, "float64-slice-val", []float64{1.1, 2.2, 3.3}, "")
		StringSliceVar(&stringSliceVal, "string-slice-val", []string{"1", "2", "3"}, "")
		TimeSliceVar(&timeSliceVal, "time-slice-val", []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)}, "")
		DurationSliceVar(&durationSliceVal, "duration-slice-val", []time.Duration{time.Second, 2 * time.Second, 3 * time.Second}, "")
		IPSliceVar(&ipSliceVal, "ip-slice-val", []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")}, "")

		So(ParseArgs(strings.Split("--bool-val "+
			"--int-val 11 --int8-val 22 --int16-val 33 --int32-val 44 --int64-val 55 "+
			"--uint-val 66 --uint8-val 77 --uint16-val 88 --uint32-val 99 --uint64-val 1010 "+
			"--float32-val 32.32 --float64-val 64.64 --string-val world "+
			"--duration-val 6s --time-val 2020-11-09T21:55:27+08:00 --ip-val 106.11.34.43 "+
			"--bool-slice-val true,false,true "+
			"--int-slice-val 11,22,33 --int8-slice-val 11,22,33 --int16-slice-val 11,22,33 --int32-slice-val 11,22,33 --int64-slice-val 11,22,33 "+
			"--uint-slice-val 11,22,33 --uint8-slice-val 11,22,33 --uint16-slice-val 11,22,33 --uint32-slice-val 11,22,33 --uint64-slice-val 11,22,33 "+
			"--float32-slice-val 11.11,22.22,33.33 --float64-slice-val 11.11,22.22,33.33 --string-slice-val 11,22,33 "+
			"--duration-slice-val 4s,5s,6s --ip-slice-val 192.168.0.3,192.168.0.4 "+
			"--time-slice-val 1970-01-01T08:00:04+0800,1970-01-01T08:00:05+08:00,1970-01-01T08:00:06+08:00", " ")), ShouldBeNil)

		So(boolVal, ShouldEqual, true)
		So(intVal, ShouldEqual, 11)
		So(int8Val, ShouldEqual, 22)
		So(int16Val, ShouldEqual, 33)
		So(int32Val, ShouldEqual, 44)
		So(int64Val, ShouldEqual, 55)
		So(uintVal, ShouldEqual, 66)
		So(uint8Val, ShouldEqual, 77)
		So(uint16Val, ShouldEqual, 88)
		So(uint32Val, ShouldEqual, 99)
		So(uint64Val, ShouldEqual, 1010)
		So(float32Val, ShouldBeBetween, 32.31, 32.33)
		So(float64Val, ShouldBeBetween, 64.63, 64.65)
		So(stringVal, ShouldEqual, "world")
		So(durationVal, ShouldEqual, 6*time.Second)
		So(timeVal, ShouldEqual, time.Unix(1604930127, 0))
		So(ipVal, ShouldEqual, net.ParseIP("106.11.34.43"))
		So(boolSliceVal, ShouldResemble, []bool{true, false, true})
		So(intSliceVal, ShouldResemble, []int{11, 22, 33})
		So(int8SliceVal, ShouldResemble, []int8{11, 22, 33})
		So(int16SliceVal, ShouldResemble, []int16{11, 22, 33})
		So(int32SliceVal, ShouldResemble, []int32{11, 22, 33})
		So(int64SliceVal, ShouldResemble, []int64{11, 22, 33})
		So(uintSliceVal, ShouldResemble, []uint{11, 22, 33})
		So(uint8SliceVal, ShouldResemble, []uint8{11, 22, 33})
		So(uint16SliceVal, ShouldResemble, []uint16{11, 22, 33})
		So(uint32SliceVal, ShouldResemble, []uint32{11, 22, 33})
		So(uint64SliceVal, ShouldResemble, []uint64{11, 22, 33})
		So(float32SliceVal, ShouldResemble, []float32{11.11, 22.22, 33.33})
		So(float64SliceVal, ShouldResemble, []float64{11.11, 22.22, 33.33})
		So(stringSliceVal, ShouldResemble, []string{"11", "22", "33"})
		So(timeSliceVal, ShouldResemble, []time.Time{time.Unix(4, 0), time.Unix(5, 0), time.Unix(6, 0)})
		So(durationSliceVal, ShouldResemble, []time.Duration{4 * time.Second, 5 * time.Second, 6 * time.Second})
		So(ipSliceVal, ShouldResemble, []net.IP{net.ParseIP("192.168.0.3"), net.ParseIP("192.168.0.4")})
	})
}

func TestGlobalFlagGet(t *testing.T) {
	Convey("TestGlobalFlagGet", t, func() {
		gflag = NewFlag("test")

		Bool("bool-val", false, "")
		Int("int-val", 1, "")
		Int8("int8-val", 2, "")
		Int16("int16-val", 3, "")
		Int32("int32-val", 4, "")
		Int64("int64-val", 5, "")
		Uint("uint-val", 6, "")
		Uint8("uint8-val", 7, "")
		Uint16("uint16-val", 8, "")
		Uint32("uint32-val", 9, "")
		Uint64("uint64-val", 10, "")
		Float32("float32-val", 3.2, "")
		Float64("float64-val", 6.4, "")
		String("string-val", "hello", "")
		Time("time-val", time.Unix(1604930126, 0), "")
		Duration("duration-val", 3*time.Second, "")
		IP("ip-val", net.ParseIP("106.11.34.42"), "")
		BoolSlice("bool-slice-val", []bool{true, true, false}, "")
		IntSlice("int-slice-val", []int{1, 2, 3}, "")
		Int8Slice("int8-slice-val", []int8{1, 2, 3}, "")
		Int16Slice("int16-slice-val", []int16{1, 2, 3}, "")
		Int32Slice("int32-slice-val", []int32{1, 2, 3}, "")
		Int64Slice("int64-slice-val", []int64{1, 2, 3}, "")
		UintSlice("uint-slice-val", []uint{1, 2, 3}, "")
		Uint8Slice("uint8-slice-val", []uint8{1, 2, 3}, "")
		Uint16Slice("uint16-slice-val", []uint16{1, 2, 3}, "")
		Uint32Slice("uint32-slice-val", []uint32{1, 2, 3}, "")
		Uint64Slice("uint64-slice-val", []uint64{1, 2, 3}, "")
		Float32Slice("float32-slice-val", []float32{1.1, 2.2, 3.3}, "")
		Float64Slice("float64-slice-val", []float64{1.1, 2.2, 3.3}, "")
		StringSlice("string-slice-val", []string{"1", "2", "3"}, "")
		TimeSlice("time-slice-val", []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)}, "")
		DurationSlice("duration-slice-val", []time.Duration{time.Second, 2 * time.Second, 3 * time.Second}, "")
		IPSlice("ip-slice-val", []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")}, "")

		So(ParseArgs(strings.Split("--bool-val "+
			"--int-val 11 --int8-val 22 --int16-val 33 --int32-val 44 --int64-val 55 "+
			"--uint-val 66 --uint8-val 77 --uint16-val 88 --uint32-val 99 --uint64-val 1010 "+
			"--float32-val 32.32 --float64-val 64.64 --string-val world "+
			"--duration-val 6s --time-val 2020-11-09T21:55:27+08:00 --ip-val 106.11.34.43 "+
			"--bool-slice-val true,false,true "+
			"--int-slice-val 11,22,33 --int8-slice-val 11,22,33 --int16-slice-val 11,22,33 --int32-slice-val 11,22,33 --int64-slice-val 11,22,33 "+
			"--uint-slice-val 11,22,33 --uint8-slice-val 11,22,33 --uint16-slice-val 11,22,33 --uint32-slice-val 11,22,33 --uint64-slice-val 11,22,33 "+
			"--float32-slice-val 11.11,22.22,33.33 --float64-slice-val 11.11,22.22,33.33 --string-slice-val 11,22,33 "+
			"--duration-slice-val 4s,5s,6s --ip-slice-val 192.168.0.3,192.168.0.4 "+
			"--time-slice-val 1970-01-01T08:00:04+0800,1970-01-01T08:00:05+08:00,1970-01-01T08:00:06+08:00", " ")), ShouldBeNil)

		So(GetBool("bool-val"), ShouldBeTrue)
		So(GetInt("int-val"), ShouldEqual, 11)
		So(GetInt8("int8-val"), ShouldEqual, 22)
		So(GetInt16("int16-val"), ShouldEqual, 33)
		So(GetInt32("int32-val"), ShouldEqual, 44)
		So(GetInt64("int64-val"), ShouldEqual, 55)
		So(GetUint("uint-val"), ShouldEqual, 66)
		So(GetUint8("uint8-val"), ShouldEqual, 77)
		So(GetUint16("uint16-val"), ShouldEqual, 88)
		So(GetUint32("uint32-val"), ShouldEqual, 99)
		So(GetUint64("uint64-val"), ShouldEqual, 1010)
		So(GetFloat32("float32-val"), ShouldBeBetween, 32.31, 32.33)
		So(GetFloat64("float64-val"), ShouldBeBetween, 64.63, 64.65)
		So(GetString("string-val"), ShouldEqual, "world")
		So(GetDuration("duration-val"), ShouldEqual, 6*time.Second)
		So(GetTime("time-val"), ShouldEqual, time.Unix(1604930127, 0))
		So(GetIP("ip-val"), ShouldEqual, net.ParseIP("106.11.34.43"))
		So(GetBoolSlice("bool-slice-val"), ShouldResemble, []bool{true, false, true})
		So(GetIntSlice("int-slice-val"), ShouldResemble, []int{11, 22, 33})
		So(GetInt8Slice("int8-slice-val"), ShouldResemble, []int8{11, 22, 33})
		So(GetInt16Slice("int16-slice-val"), ShouldResemble, []int16{11, 22, 33})
		So(GetInt32Slice("int32-slice-val"), ShouldResemble, []int32{11, 22, 33})
		So(GetInt64Slice("int64-slice-val"), ShouldResemble, []int64{11, 22, 33})
		So(GetUintSlice("uint-slice-val"), ShouldResemble, []uint{11, 22, 33})
		So(GetUint8Slice("uint8-slice-val"), ShouldResemble, []uint8{11, 22, 33})
		So(GetUint16Slice("uint16-slice-val"), ShouldResemble, []uint16{11, 22, 33})
		So(GetUint32Slice("uint32-slice-val"), ShouldResemble, []uint32{11, 22, 33})
		So(GetUint64Slice("uint64-slice-val"), ShouldResemble, []uint64{11, 22, 33})
		So(GetFloat32Slice("float32-slice-val"), ShouldResemble, []float32{11.11, 22.22, 33.33})
		So(GetFloat64Slice("float64-slice-val"), ShouldResemble, []float64{11.11, 22.22, 33.33})
		So(GetStringSlice("string-slice-val"), ShouldResemble, []string{"11", "22", "33"})
		So(GetDurationSlice("duration-slice-val"), ShouldResemble, []time.Duration{4 * time.Second, 5 * time.Second, 6 * time.Second})
		So(GetTimeSlice("time-slice-val"), ShouldResemble, []time.Time{time.Unix(4, 0), time.Unix(5, 0), time.Unix(6, 0)})
		So(GetIPSlice("ip-slice-val"), ShouldResemble, []net.IP{net.ParseIP("192.168.0.3"), net.ParseIP("192.168.0.4")})

		So(GetBoolP("bool-val"), ShouldBeTrue)
		So(GetIntP("int-val"), ShouldEqual, 11)
		So(GetInt8P("int8-val"), ShouldEqual, 22)
		So(GetInt16P("int16-val"), ShouldEqual, 33)
		So(GetInt32P("int32-val"), ShouldEqual, 44)
		So(GetInt64P("int64-val"), ShouldEqual, 55)
		So(GetUintP("uint-val"), ShouldEqual, 66)
		So(GetUint8P("uint8-val"), ShouldEqual, 77)
		So(GetUint16P("uint16-val"), ShouldEqual, 88)
		So(GetUint32P("uint32-val"), ShouldEqual, 99)
		So(GetUint64P("uint64-val"), ShouldEqual, 1010)
		So(GetFloat32P("float32-val"), ShouldBeBetween, 32.31, 32.33)
		So(GetFloat64P("float64-val"), ShouldBeBetween, 64.63, 64.65)
		So(GetStringP("string-val"), ShouldEqual, "world")
		So(GetDurationP("duration-val"), ShouldEqual, 6*time.Second)
		So(GetTimeP("time-val"), ShouldEqual, time.Unix(1604930127, 0))
		So(GetIPP("ip-val"), ShouldEqual, net.ParseIP("106.11.34.43"))
		So(GetBoolSliceP("bool-slice-val"), ShouldResemble, []bool{true, false, true})
		So(GetIntSliceP("int-slice-val"), ShouldResemble, []int{11, 22, 33})
		So(GetInt8SliceP("int8-slice-val"), ShouldResemble, []int8{11, 22, 33})
		So(GetInt16SliceP("int16-slice-val"), ShouldResemble, []int16{11, 22, 33})
		So(GetInt32SliceP("int32-slice-val"), ShouldResemble, []int32{11, 22, 33})
		So(GetInt64SliceP("int64-slice-val"), ShouldResemble, []int64{11, 22, 33})
		So(GetUintSliceP("uint-slice-val"), ShouldResemble, []uint{11, 22, 33})
		So(GetUint8SliceP("uint8-slice-val"), ShouldResemble, []uint8{11, 22, 33})
		So(GetUint16SliceP("uint16-slice-val"), ShouldResemble, []uint16{11, 22, 33})
		So(GetUint32SliceP("uint32-slice-val"), ShouldResemble, []uint32{11, 22, 33})
		So(GetUint64SliceP("uint64-slice-val"), ShouldResemble, []uint64{11, 22, 33})
		So(GetFloat32SliceP("float32-slice-val"), ShouldResemble, []float32{11.11, 22.22, 33.33})
		So(GetFloat64SliceP("float64-slice-val"), ShouldResemble, []float64{11.11, 22.22, 33.33})
		So(GetStringSliceP("string-slice-val"), ShouldResemble, []string{"11", "22", "33"})
		So(GetDurationSliceP("duration-slice-val"), ShouldResemble, []time.Duration{4 * time.Second, 5 * time.Second, 6 * time.Second})
		So(GetTimeSliceP("time-slice-val"), ShouldResemble, []time.Time{time.Unix(4, 0), time.Unix(5, 0), time.Unix(6, 0)})
		So(GetIPSliceP("ip-slice-val"), ShouldResemble, []net.IP{net.ParseIP("192.168.0.3"), net.ParseIP("192.168.0.4")})

		So(GetBoolD("bool-val", false), ShouldBeTrue)
		So(GetIntD("int-val", 1), ShouldEqual, 11)
		So(GetInt8D("int8-val", 2), ShouldEqual, 22)
		So(GetInt16D("int16-val", 3), ShouldEqual, 33)
		So(GetInt32D("int32-val", 4), ShouldEqual, 44)
		So(GetInt64D("int64-val", 5), ShouldEqual, 55)
		So(GetUintD("uint-val", 6), ShouldEqual, 66)
		So(GetUint8D("uint8-val", 7), ShouldEqual, 77)
		So(GetUint16D("uint16-val", 8), ShouldEqual, 88)
		So(GetUint32D("uint32-val", 9), ShouldEqual, 99)
		So(GetUint64D("uint64-val", 10), ShouldEqual, 1010)
		So(GetFloat32D("float32-val", 3.2), ShouldBeBetween, 32.31, 32.33)
		So(GetFloat64D("float64-val", 6.4), ShouldBeBetween, 64.63, 64.65)
		So(GetStringD("string-val", "hello"), ShouldEqual, "world")
		So(GetDurationD("duration-val", 3*time.Second), ShouldEqual, 6*time.Second)
		So(GetTimeD("time-val", time.Unix(1604930126, 0)), ShouldEqual, time.Unix(1604930127, 0))
		So(GetIPD("ip-val", net.ParseIP("106.11.34.42")), ShouldEqual, net.ParseIP("106.11.34.43"))
		So(GetBoolSliceD("x", []bool{true, true, false}), ShouldResemble, []bool{true, true, false})
		So(GetIntSliceD("x", []int{1, 2, 3}), ShouldResemble, []int{1, 2, 3})
		So(GetInt8SliceD("x", []int8{1, 2, 3}), ShouldResemble, []int8{1, 2, 3})
		So(GetInt16SliceD("x", []int16{1, 2, 3}), ShouldResemble, []int16{1, 2, 3})
		So(GetInt32SliceD("x", []int32{1, 2, 3}), ShouldResemble, []int32{1, 2, 3})
		So(GetInt64SliceD("x", []int64{1, 2, 3}), ShouldResemble, []int64{1, 2, 3})
		So(GetUintSliceD("x", []uint{1, 2, 3}), ShouldResemble, []uint{1, 2, 3})
		So(GetUint8SliceD("x", []uint8{1, 2, 3}), ShouldResemble, []uint8{1, 2, 3})
		So(GetUint16SliceD("x", []uint16{1, 2, 3}), ShouldResemble, []uint16{1, 2, 3})
		So(GetUint32SliceD("x", []uint32{1, 2, 3}), ShouldResemble, []uint32{1, 2, 3})
		So(GetUint64SliceD("x", []uint64{1, 2, 3}), ShouldResemble, []uint64{1, 2, 3})
		So(GetFloat32SliceD("x", []float32{1.1, 2.2, 3.3}), ShouldResemble, []float32{1.1, 2.2, 3.3})
		So(GetFloat64SliceD("x", []float64{1.1, 2.2, 3.3}), ShouldResemble, []float64{1.1, 2.2, 3.3})
		So(GetStringSliceD("x", []string{"1", "2", "3"}), ShouldResemble, []string{"1", "2", "3"})
		So(GetDurationSliceD("x", []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second}), ShouldResemble, []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second})
		So(GetTimeSliceD("x", []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)}), ShouldResemble, []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)})
		So(GetIPSliceD("x", []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")}), ShouldResemble, []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")})
	})
}
