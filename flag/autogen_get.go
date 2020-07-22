// this file is auto generate by autogen.py. do not edit!
package flag

import (
	"net"
	"time"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/cast"
)

func (f *Flag) GetBoolE(key string) (bool, error) {
	v, ok := f.Get(key)
	if !ok {
		var res bool
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToBoolE(v)
}

func (f *Flag) GetBool(key string) bool {
	var res bool
	v, err := f.GetBoolE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetBoolD(key string, val bool) bool {
	v, err := f.GetBoolE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetBoolP(key string) bool {
	v, err := f.GetBoolE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetIntE(key string) (int, error) {
	v, ok := f.Get(key)
	if !ok {
		var res int
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToIntE(v)
}

func (f *Flag) GetInt(key string) int {
	var res int
	v, err := f.GetIntE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetIntD(key string, val int) int {
	v, err := f.GetIntE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetIntP(key string) int {
	v, err := f.GetIntE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetUintE(key string) (uint, error) {
	v, ok := f.Get(key)
	if !ok {
		var res uint
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToUintE(v)
}

func (f *Flag) GetUint(key string) uint {
	var res uint
	v, err := f.GetUintE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetUintD(key string, val uint) uint {
	v, err := f.GetUintE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetUintP(key string) uint {
	v, err := f.GetUintE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetInt64E(key string) (int64, error) {
	v, ok := f.Get(key)
	if !ok {
		var res int64
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToInt64E(v)
}

func (f *Flag) GetInt64(key string) int64 {
	var res int64
	v, err := f.GetInt64E(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetInt64D(key string, val int64) int64 {
	v, err := f.GetInt64E(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetInt64P(key string) int64 {
	v, err := f.GetInt64E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetInt32E(key string) (int32, error) {
	v, ok := f.Get(key)
	if !ok {
		var res int32
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToInt32E(v)
}

func (f *Flag) GetInt32(key string) int32 {
	var res int32
	v, err := f.GetInt32E(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetInt32D(key string, val int32) int32 {
	v, err := f.GetInt32E(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetInt32P(key string) int32 {
	v, err := f.GetInt32E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetInt16E(key string) (int16, error) {
	v, ok := f.Get(key)
	if !ok {
		var res int16
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToInt16E(v)
}

func (f *Flag) GetInt16(key string) int16 {
	var res int16
	v, err := f.GetInt16E(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetInt16D(key string, val int16) int16 {
	v, err := f.GetInt16E(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetInt16P(key string) int16 {
	v, err := f.GetInt16E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetInt8E(key string) (int8, error) {
	v, ok := f.Get(key)
	if !ok {
		var res int8
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToInt8E(v)
}

func (f *Flag) GetInt8(key string) int8 {
	var res int8
	v, err := f.GetInt8E(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetInt8D(key string, val int8) int8 {
	v, err := f.GetInt8E(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetInt8P(key string) int8 {
	v, err := f.GetInt8E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetUint64E(key string) (uint64, error) {
	v, ok := f.Get(key)
	if !ok {
		var res uint64
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToUint64E(v)
}

func (f *Flag) GetUint64(key string) uint64 {
	var res uint64
	v, err := f.GetUint64E(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetUint64D(key string, val uint64) uint64 {
	v, err := f.GetUint64E(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetUint64P(key string) uint64 {
	v, err := f.GetUint64E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetUint32E(key string) (uint32, error) {
	v, ok := f.Get(key)
	if !ok {
		var res uint32
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToUint32E(v)
}

func (f *Flag) GetUint32(key string) uint32 {
	var res uint32
	v, err := f.GetUint32E(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetUint32D(key string, val uint32) uint32 {
	v, err := f.GetUint32E(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetUint32P(key string) uint32 {
	v, err := f.GetUint32E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetUint16E(key string) (uint16, error) {
	v, ok := f.Get(key)
	if !ok {
		var res uint16
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToUint16E(v)
}

func (f *Flag) GetUint16(key string) uint16 {
	var res uint16
	v, err := f.GetUint16E(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetUint16D(key string, val uint16) uint16 {
	v, err := f.GetUint16E(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetUint16P(key string) uint16 {
	v, err := f.GetUint16E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetUint8E(key string) (uint8, error) {
	v, ok := f.Get(key)
	if !ok {
		var res uint8
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToUint8E(v)
}

func (f *Flag) GetUint8(key string) uint8 {
	var res uint8
	v, err := f.GetUint8E(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetUint8D(key string, val uint8) uint8 {
	v, err := f.GetUint8E(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetUint8P(key string) uint8 {
	v, err := f.GetUint8E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetFloat64E(key string) (float64, error) {
	v, ok := f.Get(key)
	if !ok {
		var res float64
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToFloat64E(v)
}

func (f *Flag) GetFloat64(key string) float64 {
	var res float64
	v, err := f.GetFloat64E(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetFloat64D(key string, val float64) float64 {
	v, err := f.GetFloat64E(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetFloat64P(key string) float64 {
	v, err := f.GetFloat64E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetFloat32E(key string) (float32, error) {
	v, ok := f.Get(key)
	if !ok {
		var res float32
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToFloat32E(v)
}

func (f *Flag) GetFloat32(key string) float32 {
	var res float32
	v, err := f.GetFloat32E(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetFloat32D(key string, val float32) float32 {
	v, err := f.GetFloat32E(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetFloat32P(key string) float32 {
	v, err := f.GetFloat32E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetStringE(key string) (string, error) {
	v, ok := f.Get(key)
	if !ok {
		var res string
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToStringE(v)
}

func (f *Flag) GetString(key string) string {
	var res string
	v, err := f.GetStringE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetStringD(key string, val string) string {
	v, err := f.GetStringE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetStringP(key string) string {
	v, err := f.GetStringE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetDurationE(key string) (time.Duration, error) {
	v, ok := f.Get(key)
	if !ok {
		var res time.Duration
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToDurationE(v)
}

func (f *Flag) GetDuration(key string) time.Duration {
	var res time.Duration
	v, err := f.GetDurationE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetDurationD(key string, val time.Duration) time.Duration {
	v, err := f.GetDurationE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetDurationP(key string) time.Duration {
	v, err := f.GetDurationE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetTimeE(key string) (time.Time, error) {
	v, ok := f.Get(key)
	if !ok {
		var res time.Time
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToTimeE(v)
}

func (f *Flag) GetTime(key string) time.Time {
	var res time.Time
	v, err := f.GetTimeE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetTimeD(key string, val time.Time) time.Time {
	v, err := f.GetTimeE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetTimeP(key string) time.Time {
	v, err := f.GetTimeE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetIPE(key string) (net.IP, error) {
	v, ok := f.Get(key)
	if !ok {
		var res net.IP
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToIPE(v)
}

func (f *Flag) GetIP(key string) net.IP {
	var res net.IP
	v, err := f.GetIPE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetIPD(key string, val net.IP) net.IP {
	v, err := f.GetIPE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetIPP(key string) net.IP {
	v, err := f.GetIPE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetBoolSliceE(key string) ([]bool, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []bool
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToBoolSliceE(v)
}

func (f *Flag) GetBoolSlice(key string) []bool {
	var res []bool
	v, err := f.GetBoolSliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetBoolSliceD(key string, val []bool) []bool {
	v, err := f.GetBoolSliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetBoolSliceP(key string) []bool {
	v, err := f.GetBoolSliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetIntSliceE(key string) ([]int, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []int
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToIntSliceE(v)
}

func (f *Flag) GetIntSlice(key string) []int {
	var res []int
	v, err := f.GetIntSliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetIntSliceD(key string, val []int) []int {
	v, err := f.GetIntSliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetIntSliceP(key string) []int {
	v, err := f.GetIntSliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetUintSliceE(key string) ([]uint, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []uint
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToUintSliceE(v)
}

func (f *Flag) GetUintSlice(key string) []uint {
	var res []uint
	v, err := f.GetUintSliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetUintSliceD(key string, val []uint) []uint {
	v, err := f.GetUintSliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetUintSliceP(key string) []uint {
	v, err := f.GetUintSliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetInt64SliceE(key string) ([]int64, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []int64
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToInt64SliceE(v)
}

func (f *Flag) GetInt64Slice(key string) []int64 {
	var res []int64
	v, err := f.GetInt64SliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetInt64SliceD(key string, val []int64) []int64 {
	v, err := f.GetInt64SliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetInt64SliceP(key string) []int64 {
	v, err := f.GetInt64SliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetInt32SliceE(key string) ([]int32, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []int32
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToInt32SliceE(v)
}

func (f *Flag) GetInt32Slice(key string) []int32 {
	var res []int32
	v, err := f.GetInt32SliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetInt32SliceD(key string, val []int32) []int32 {
	v, err := f.GetInt32SliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetInt32SliceP(key string) []int32 {
	v, err := f.GetInt32SliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetInt16SliceE(key string) ([]int16, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []int16
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToInt16SliceE(v)
}

func (f *Flag) GetInt16Slice(key string) []int16 {
	var res []int16
	v, err := f.GetInt16SliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetInt16SliceD(key string, val []int16) []int16 {
	v, err := f.GetInt16SliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetInt16SliceP(key string) []int16 {
	v, err := f.GetInt16SliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetInt8SliceE(key string) ([]int8, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []int8
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToInt8SliceE(v)
}

func (f *Flag) GetInt8Slice(key string) []int8 {
	var res []int8
	v, err := f.GetInt8SliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetInt8SliceD(key string, val []int8) []int8 {
	v, err := f.GetInt8SliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetInt8SliceP(key string) []int8 {
	v, err := f.GetInt8SliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetUint64SliceE(key string) ([]uint64, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []uint64
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToUint64SliceE(v)
}

func (f *Flag) GetUint64Slice(key string) []uint64 {
	var res []uint64
	v, err := f.GetUint64SliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetUint64SliceD(key string, val []uint64) []uint64 {
	v, err := f.GetUint64SliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetUint64SliceP(key string) []uint64 {
	v, err := f.GetUint64SliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetUint32SliceE(key string) ([]uint32, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []uint32
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToUint32SliceE(v)
}

func (f *Flag) GetUint32Slice(key string) []uint32 {
	var res []uint32
	v, err := f.GetUint32SliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetUint32SliceD(key string, val []uint32) []uint32 {
	v, err := f.GetUint32SliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetUint32SliceP(key string) []uint32 {
	v, err := f.GetUint32SliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetUint16SliceE(key string) ([]uint16, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []uint16
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToUint16SliceE(v)
}

func (f *Flag) GetUint16Slice(key string) []uint16 {
	var res []uint16
	v, err := f.GetUint16SliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetUint16SliceD(key string, val []uint16) []uint16 {
	v, err := f.GetUint16SliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetUint16SliceP(key string) []uint16 {
	v, err := f.GetUint16SliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetUint8SliceE(key string) ([]uint8, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []uint8
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToUint8SliceE(v)
}

func (f *Flag) GetUint8Slice(key string) []uint8 {
	var res []uint8
	v, err := f.GetUint8SliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetUint8SliceD(key string, val []uint8) []uint8 {
	v, err := f.GetUint8SliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetUint8SliceP(key string) []uint8 {
	v, err := f.GetUint8SliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetFloat64SliceE(key string) ([]float64, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []float64
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToFloat64SliceE(v)
}

func (f *Flag) GetFloat64Slice(key string) []float64 {
	var res []float64
	v, err := f.GetFloat64SliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetFloat64SliceD(key string, val []float64) []float64 {
	v, err := f.GetFloat64SliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetFloat64SliceP(key string) []float64 {
	v, err := f.GetFloat64SliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetFloat32SliceE(key string) ([]float32, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []float32
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToFloat32SliceE(v)
}

func (f *Flag) GetFloat32Slice(key string) []float32 {
	var res []float32
	v, err := f.GetFloat32SliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetFloat32SliceD(key string, val []float32) []float32 {
	v, err := f.GetFloat32SliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetFloat32SliceP(key string) []float32 {
	v, err := f.GetFloat32SliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetStringSliceE(key string) ([]string, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []string
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToStringSliceE(v)
}

func (f *Flag) GetStringSlice(key string) []string {
	var res []string
	v, err := f.GetStringSliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetStringSliceD(key string, val []string) []string {
	v, err := f.GetStringSliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetStringSliceP(key string) []string {
	v, err := f.GetStringSliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetDurationSliceE(key string) ([]time.Duration, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []time.Duration
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToDurationSliceE(v)
}

func (f *Flag) GetDurationSlice(key string) []time.Duration {
	var res []time.Duration
	v, err := f.GetDurationSliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetDurationSliceD(key string, val []time.Duration) []time.Duration {
	v, err := f.GetDurationSliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetDurationSliceP(key string) []time.Duration {
	v, err := f.GetDurationSliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetTimeSliceE(key string) ([]time.Time, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []time.Time
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToTimeSliceE(v)
}

func (f *Flag) GetTimeSlice(key string) []time.Time {
	var res []time.Time
	v, err := f.GetTimeSliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetTimeSliceD(key string, val []time.Time) []time.Time {
	v, err := f.GetTimeSliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetTimeSliceP(key string) []time.Time {
	v, err := f.GetTimeSliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *Flag) GetIPSliceE(key string) ([]net.IP, error) {
	v, ok := f.Get(key)
	if !ok {
		var res []net.IP
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}
	return cast.ToIPSliceE(v)
}

func (f *Flag) GetIPSlice(key string) []net.IP {
	var res []net.IP
	v, err := f.GetIPSliceE(key)
	if err != nil {
		return res
	}
	return v
}

func (f *Flag) GetIPSliceD(key string, val []net.IP) []net.IP {
	v, err := f.GetIPSliceE(key)
	if err != nil {
		return val
	}
	return v
}

func (f *Flag) GetIPSliceP(key string) []net.IP {
	v, err := f.GetIPSliceE(key)
	if err != nil {
		panic(err)
	}
	return v
}
