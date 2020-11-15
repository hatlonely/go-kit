package flag_test

import (
	"net"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/flag"
)

func TestGlobalFlagBind(t *testing.T) {
	Convey("TestGlobalFlag", t, func() {
		boolVal := flag.Bool("bool-val", false, "")
		intVal := flag.Int("int-val", 1, "")
		int8Val := flag.Int8("int8-val", 2, "")
		int16Val := flag.Int16("int16-val", 3, "")
		int32Val := flag.Int32("int32-val", 4, "")
		int64Val := flag.Int64("int64-val", 5, "")
		uintVal := flag.Uint("uint-val", 6, "")
		uint8Val := flag.Uint8("uint8-val", 7, "")
		uint16Val := flag.Uint16("uint16-val", 8, "")
		uint32Val := flag.Uint32("uint32-val", 9, "")
		uint64Val := flag.Uint64("uint64-val", 10, "")
		float32Val := flag.Float32("float32-val", 3.2, "")
		float64Val := flag.Float64("float64-val", 6.4, "")
		stringVal := flag.String("string-val", "hello", "")
		timeVal := flag.Time("time-val", time.Unix(1604930126, 0), "")
		durationVal := flag.Duration("duration-val", 3*time.Second, "")
		ipVal := flag.IP("ip-val", net.ParseIP("106.11.34.42"), "")
		boolSliceVal := flag.BoolSlice("bool-slice-val", []bool{true, true, false}, "")
		intSliceVal := flag.IntSlice("int-slice-val", []int{1, 2, 3}, "")
		int8SliceVal := flag.Int8Slice("int8-slice-val", []int8{1, 2, 3}, "")
		int16SliceVal := flag.Int16Slice("int16-slice-val", []int16{1, 2, 3}, "")
		int32SliceVal := flag.Int32Slice("int32-slice-val", []int32{1, 2, 3}, "")
		int64SliceVal := flag.Int64Slice("int64-slice-val", []int64{1, 2, 3}, "")
		uintSliceVal := flag.UintSlice("uint-slice-val", []uint{1, 2, 3}, "")
		uint8SliceVal := flag.Uint8Slice("uint8-slice-val", []uint8{1, 2, 3}, "")
		uint16SliceVal := flag.Uint16Slice("uint16-slice-val", []uint16{1, 2, 3}, "")
		uint32SliceVal := flag.Uint32Slice("uint32-slice-val", []uint32{1, 2, 3}, "")
		uint64SliceVal := flag.Uint64Slice("uint64-slice-val", []uint64{1, 2, 3}, "")
		float32SliceVal := flag.Float32Slice("float32-slice-val", []float32{1.1, 2.2, 3.3}, "")
		float64SliceVal := flag.Float64Slice("float64-slice-val", []float64{1.1, 2.2, 3.3}, "")
		stringSliceVal := flag.StringSlice("string-slice-val", []string{"1", "2", "3"}, "")
		timeSliceVal := flag.TimeSlice("time-slice-val", []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)}, "")
		durationSliceVal := flag.DurationSlice("duration-slice-val", []time.Duration{time.Second, 2 * time.Second, 3 * time.Second}, "")
		ipSliceVal := flag.IPSlice("ip-slice-val", []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")}, "")

		So(flag.ParseArgs(strings.Split("--bool-val "+
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
	Convey("TestGlobalFlag", t, func() {
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

		flag.BoolVar(&boolVal, "bool-val", false, "")
		flag.IntVar(&intVal, "int-val", 1, "")
		flag.Int8Var(&int8Val, "int8-val", 2, "")
		flag.Int16Var(&int16Val, "int16-val", 3, "")
		flag.Int32Var(&int32Val, "int32-val", 4, "")
		flag.Int64Var(&int64Val, "int64-val", 5, "")
		flag.UintVar(&uintVal, "uint-val", 6, "")
		flag.Uint8Var(&uint8Val, "uint8-val", 7, "")
		flag.Uint16Var(&uint16Val, "uint16-val", 8, "")
		flag.Uint32Var(&uint32Val, "uint32-val", 9, "")
		flag.Uint64Var(&uint64Val, "uint64-val", 10, "")
		flag.Float32Var(&float32Val, "float32-val", 3.2, "")
		flag.Float64Var(&float64Val, "float64-val", 6.4, "")
		flag.StringVar(&stringVal, "string-val", "hello", "")
		flag.TimeVar(&timeVal, "time-val", time.Unix(1604930126, 0), "")
		flag.DurationVar(&durationVal, "duration-val", 3*time.Second, "")
		flag.IPVar(&ipVal, "ip-val", net.ParseIP("106.11.34.42"), "")
		flag.BoolSliceVar(&boolSliceVal, "bool-slice-val", []bool{true, true, false}, "")
		flag.IntSliceVar(&intSliceVal, "int-slice-val", []int{1, 2, 3}, "")
		flag.Int8SliceVar(&int8SliceVal, "int8-slice-val", []int8{1, 2, 3}, "")
		flag.Int16SliceVar(&int16SliceVal, "int16-slice-val", []int16{1, 2, 3}, "")
		flag.Int32SliceVar(&int32SliceVal, "int32-slice-val", []int32{1, 2, 3}, "")
		flag.Int64SliceVar(&int64SliceVal, "int64-slice-val", []int64{1, 2, 3}, "")
		flag.UintSliceVar(&uintSliceVal, "uint-slice-val", []uint{1, 2, 3}, "")
		flag.Uint8SliceVar(&uint8SliceVal, "uint8-slice-val", []uint8{1, 2, 3}, "")
		flag.Uint16SliceVar(&uint16SliceVal, "uint16-slice-val", []uint16{1, 2, 3}, "")
		flag.Uint32SliceVar(&uint32SliceVal, "uint32-slice-val", []uint32{1, 2, 3}, "")
		flag.Uint64SliceVar(&uint64SliceVal, "uint64-slice-val", []uint64{1, 2, 3}, "")
		flag.Float32SliceVar(&float32SliceVal, "float32-slice-val", []float32{1.1, 2.2, 3.3}, "")
		flag.Float64SliceVar(&float64SliceVal, "float64-slice-val", []float64{1.1, 2.2, 3.3}, "")
		flag.StringSliceVar(&stringSliceVal, "string-slice-val", []string{"1", "2", "3"}, "")
		flag.TimeSliceVar(&timeSliceVal, "time-slice-val", []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)}, "")
		flag.DurationSliceVar(&durationSliceVal, "duration-slice-val", []time.Duration{time.Second, 2 * time.Second, 3 * time.Second}, "")
		flag.IPSliceVar(&ipSliceVal, "ip-slice-val", []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")}, "")

		So(flag.ParseArgs(strings.Split("--bool-val "+
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
