package cast

import (
	"fmt"
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestToBoolSlice(t *testing.T) {
	Convey("TestToBoolSlice", t, func() {
		So(ToBoolSlice([]bool{true, false, true}), ShouldResemble, []bool{true, false, true})
		So(ToBoolSlice([]interface{}{"true", false, 1, 0}), ShouldResemble, []bool{true, false, true, false})
		So(ToBoolSlice([]string{"true", "1", "false"}), ShouldResemble, []bool{true, true, false})
		So(ToBoolSlice("true,1,false"), ShouldResemble, []bool{true, true, false})
	})
}

func TestToIntSlice(t *testing.T) {
	Convey("TestToIntSlice", t, func() {
		So(ToIntSlice([]int{1, 2, 3}), ShouldResemble, []int{1, 2, 3})
		So(ToIntSlice([]interface{}{"1", 2, int64(3)}), ShouldResemble, []int{1, 2, 3})
		So(ToIntSlice([]string{"1", "2", "3"}), ShouldResemble, []int{1, 2, 3})
		So(ToIntSlice("1,2,3"), ShouldResemble, []int{1, 2, 3})
	})
}

func TestToInt8Slice(t *testing.T) {
	Convey("TestToInt8Slice", t, func() {
		So(ToInt8Slice([]int8{1, 2, 3}), ShouldResemble, []int8{1, 2, 3})
		So(ToInt8Slice([]interface{}{"1", 2, int64(3)}), ShouldResemble, []int8{1, 2, 3})
		So(ToInt8Slice([]string{"1", "2", "3"}), ShouldResemble, []int8{1, 2, 3})
		So(ToInt8Slice("1,2,3"), ShouldResemble, []int8{1, 2, 3})
	})
}

func TestToInt16Slice(t *testing.T) {
	Convey("TestToInt16Slice", t, func() {
		So(ToInt16Slice([]int16{1, 2, 3}), ShouldResemble, []int16{1, 2, 3})
		So(ToInt16Slice([]interface{}{"1", 2, int64(3)}), ShouldResemble, []int16{1, 2, 3})
		So(ToInt16Slice([]string{"1", "2", "3"}), ShouldResemble, []int16{1, 2, 3})
		So(ToInt16Slice("1,2,3"), ShouldResemble, []int16{1, 2, 3})
	})
}

func TestToInt32Slice(t *testing.T) {
	Convey("TestToInt32Slice", t, func() {
		So(ToInt32Slice([]int32{1, 2, 3}), ShouldResemble, []int32{1, 2, 3})
		So(ToInt32Slice([]interface{}{"1", 2, int64(3)}), ShouldResemble, []int32{1, 2, 3})
		So(ToInt32Slice([]string{"1", "2", "3"}), ShouldResemble, []int32{1, 2, 3})
		So(ToInt32Slice("1,2,3"), ShouldResemble, []int32{1, 2, 3})
	})
}

func TestToInt64Slice(t *testing.T) {
	Convey("TestToInt64Slice", t, func() {
		So(ToInt64Slice([]int64{1, 2, 3}), ShouldResemble, []int64{1, 2, 3})
		So(ToInt64Slice([]interface{}{"1", 2, int64(3)}), ShouldResemble, []int64{1, 2, 3})
		So(ToInt64Slice([]string{"1", "2", "3"}), ShouldResemble, []int64{1, 2, 3})
		So(ToInt64Slice("1,2,3"), ShouldResemble, []int64{1, 2, 3})
	})
}

func TestToUintSlice(t *testing.T) {
	Convey("TestToUintSlice", t, func() {
		So(ToUintSlice([]uint{1, 2, 3}), ShouldResemble, []uint{1, 2, 3})
		So(ToUintSlice([]interface{}{"1", 2, uint64(3)}), ShouldResemble, []uint{1, 2, 3})
		So(ToUintSlice([]string{"1", "2", "3"}), ShouldResemble, []uint{1, 2, 3})
		So(ToUintSlice("1,2,3"), ShouldResemble, []uint{1, 2, 3})
	})
}

func TestToUint8Slice(t *testing.T) {
	Convey("TestToUint8Slice", t, func() {
		So(ToUint8Slice([]uint8{1, 2, 3}), ShouldResemble, []uint8{1, 2, 3})
		So(ToUint8Slice([]interface{}{"1", 2, uint64(3)}), ShouldResemble, []uint8{1, 2, 3})
		So(ToUint8Slice([]string{"1", "2", "3"}), ShouldResemble, []uint8{1, 2, 3})
		So(ToUint8Slice("1,2,3"), ShouldResemble, []uint8{1, 2, 3})
	})
}

func TestToUint16Slice(t *testing.T) {
	Convey("TestToUint16Slice", t, func() {
		So(ToUint16Slice([]uint16{1, 2, 3}), ShouldResemble, []uint16{1, 2, 3})
		So(ToUint16Slice([]interface{}{"1", 2, uint64(3)}), ShouldResemble, []uint16{1, 2, 3})
		So(ToUint16Slice([]string{"1", "2", "3"}), ShouldResemble, []uint16{1, 2, 3})
		So(ToUint16Slice("1,2,3"), ShouldResemble, []uint16{1, 2, 3})
	})
}

func TestToUint32Slice(t *testing.T) {
	Convey("TestToUint32Slice", t, func() {
		So(ToUint32Slice([]uint32{1, 2, 3}), ShouldResemble, []uint32{1, 2, 3})
		So(ToUint32Slice([]interface{}{"1", 2, uint64(3)}), ShouldResemble, []uint32{1, 2, 3})
		So(ToUint32Slice([]string{"1", "2", "3"}), ShouldResemble, []uint32{1, 2, 3})
		So(ToUint32Slice("1,2,3"), ShouldResemble, []uint32{1, 2, 3})
	})
}

func TestToUint64Slice(t *testing.T) {
	Convey("TestToUint64Slice", t, func() {
		So(ToUint64Slice([]uint64{1, 2, 3}), ShouldResemble, []uint64{1, 2, 3})
		So(ToUint64Slice([]interface{}{"1", 2, uint64(3)}), ShouldResemble, []uint64{1, 2, 3})
		So(ToUint64Slice([]string{"1", "2", "3"}), ShouldResemble, []uint64{1, 2, 3})
		So(ToUint64Slice("1,2,3"), ShouldResemble, []uint64{1, 2, 3})
	})
}

func nearlyEqualFloat32Slice(s1, s2 []float32) bool {
	if len(s1) != len(s2) {
		fmt.Println(s1, s2)
		return false
	}
	for i := range s1 {
		if s1[i]-s2[i] < -0.01 || s1[i]-s2[i] > 0.01 {
			fmt.Println(s1[i], s2[i])
			return false
		}
	}
	return true
}

func nearlyEqualFloat64Slice(s1, s2 []float64) bool {
	if len(s1) != len(s2) {
		fmt.Println(s1, s2)
		return false
	}
	for i := range s1 {
		if s1[i]-s2[i] < -0.01 || s1[i]-s2[i] > 0.01 {
			fmt.Println(s1[i], s2[i])
			return false
		}
	}
	return true
}

func TestToFloat32Slice(t *testing.T) {
	Convey("TestToFloat32Slice", t, func() {
		So(nearlyEqualFloat32Slice(ToFloat32Slice([]float32{1.1, 2.2, 3.3}), []float32{1.1, 2.2, 3.3}), ShouldBeTrue)
		So(nearlyEqualFloat32Slice(ToFloat32Slice([]interface{}{"1.1", 2.2, "3.3"}), []float32{1.1, 2.2, 3.3}), ShouldBeTrue)
		So(nearlyEqualFloat32Slice(ToFloat32Slice([]string{"1.1", "2.2", "3.3"}), []float32{1.1, 2.2, 3.3}), ShouldBeTrue)
		So(nearlyEqualFloat32Slice(ToFloat32Slice("1.1,2.2,3.3"), []float32{1.1, 2.2, 3.3}), ShouldBeTrue)
	})
}

func TestToFloat64Slice(t *testing.T) {
	Convey("TestToFloat64Slice", t, func() {
		So(nearlyEqualFloat64Slice(ToFloat64Slice([]float64{1.1, 2.2, 3.3}), []float64{1.1, 2.2, 3.3}), ShouldBeTrue)
		So(nearlyEqualFloat64Slice(ToFloat64Slice([]interface{}{"1.1", 2.2, "3.3"}), []float64{1.1, 2.2, 3.3}), ShouldBeTrue)
		So(nearlyEqualFloat64Slice(ToFloat64Slice([]string{"1.1", "2.2", "3.3"}), []float64{1.1, 2.2, 3.3}), ShouldBeTrue)
		So(nearlyEqualFloat64Slice(ToFloat64Slice("1.1,2.2,3.3"), []float64{1.1, 2.2, 3.3}), ShouldBeTrue)
	})
}

func TestToDurationSlice(t *testing.T) {
	Convey("TestToDurationSlice", t, func() {
		So(ToDurationSlice([]time.Duration{1, 2, 3}), ShouldResemble, []time.Duration{1, 2, 3})
		So(ToDurationSlice([]interface{}{"1", 2, time.Duration(3)}), ShouldResemble, []time.Duration{1, 2, 3})
		So(ToDurationSlice([]string{"1s", "2", "3"}), ShouldResemble, []time.Duration{time.Second, 2, 3})
		So(ToDurationSlice("1s,2s,3s"), ShouldResemble, []time.Duration{time.Second, 2 * time.Second, 3 * time.Second})
	})
}

func TestToTimeSlice(t *testing.T) {
	Convey("TestToTimeSlice", t, func() {
		So(ToTimeSlice([]time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)}), ShouldResemble, []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)})
		So(ToTimeSlice([]interface{}{"1970-01-01 08:00:01 +0800 CST", 2, 3}), ShouldResemble, []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)})
		So(ToTimeSlice([]string{"1970-01-01 08:00:01 +0800 CST", "1970-01-01T08:00:02+08:00", "1970-01-01T08:00:03+08:00"}), ShouldResemble, []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)})
		So(ToTimeSlice("1970-01-01 08:00:01 +0800 CST,1970-01-01T08:00:02+08:00,1970-01-01T08:00:03+08:00"), ShouldResemble, []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)})
	})
}

func TestToIPSlice(t *testing.T) {
	Convey("TestToIPSlice", t, func() {
		So(ToIPSlice([]net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")}), ShouldResemble, []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")})
		So(ToIPSlice([]interface{}{net.ParseIP("192.168.0.1"), "192.168.0.2"}), ShouldResemble, []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")})
		So(ToIPSlice([]string{"192.168.0.1", "192.168.0.2"}), ShouldResemble, []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")})
		So(ToIPSlice("192.168.0.1,192.168.0.2"), ShouldResemble, []net.IP{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2")})
	})
}

func TestToStringSlice(t *testing.T) {
	Convey("TestToStringSlice", t, func() {
		So(ToStringSlice([]string{"1", "2", "3"}), ShouldResemble, []string{"1", "2", "3"})
		So(ToStringSlice([]interface{}{"1", 2, uint64(3)}), ShouldResemble, []string{"1", "2", "3"})
		So(ToStringSlice([]string{"1", "2", "3"}), ShouldResemble, []string{"1", "2", "3"})
		So(ToStringSlice("1,2,3"), ShouldResemble, []string{"1", "2", "3"})
	})
}
