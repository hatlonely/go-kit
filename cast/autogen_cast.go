// this file is auto generate by autogen.py. do not edit!
package cast

import (
	"net"
	"regexp"
	"time"
)

func ToBool(val interface{}) bool {
	if v, err := ToBoolE(val); err == nil {
		return v
	}
	var v bool
	return v
}

func ToBoolD(val interface{}, defaultValue bool) bool {
	if v, err := ToBoolE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToBoolP(val interface{}) bool {
	v, err := ToBoolE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToInt(val interface{}) int {
	if v, err := ToIntE(val); err == nil {
		return v
	}
	var v int
	return v
}

func ToIntD(val interface{}, defaultValue int) int {
	if v, err := ToIntE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToIntP(val interface{}) int {
	v, err := ToIntE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToUint(val interface{}) uint {
	if v, err := ToUintE(val); err == nil {
		return v
	}
	var v uint
	return v
}

func ToUintD(val interface{}, defaultValue uint) uint {
	if v, err := ToUintE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToUintP(val interface{}) uint {
	v, err := ToUintE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToInt64(val interface{}) int64 {
	if v, err := ToInt64E(val); err == nil {
		return v
	}
	var v int64
	return v
}

func ToInt64D(val interface{}, defaultValue int64) int64 {
	if v, err := ToInt64E(val); err == nil {
		return v
	}
	return defaultValue
}

func ToInt64P(val interface{}) int64 {
	v, err := ToInt64E(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToInt32(val interface{}) int32 {
	if v, err := ToInt32E(val); err == nil {
		return v
	}
	var v int32
	return v
}

func ToInt32D(val interface{}, defaultValue int32) int32 {
	if v, err := ToInt32E(val); err == nil {
		return v
	}
	return defaultValue
}

func ToInt32P(val interface{}) int32 {
	v, err := ToInt32E(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToInt16(val interface{}) int16 {
	if v, err := ToInt16E(val); err == nil {
		return v
	}
	var v int16
	return v
}

func ToInt16D(val interface{}, defaultValue int16) int16 {
	if v, err := ToInt16E(val); err == nil {
		return v
	}
	return defaultValue
}

func ToInt16P(val interface{}) int16 {
	v, err := ToInt16E(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToInt8(val interface{}) int8 {
	if v, err := ToInt8E(val); err == nil {
		return v
	}
	var v int8
	return v
}

func ToInt8D(val interface{}, defaultValue int8) int8 {
	if v, err := ToInt8E(val); err == nil {
		return v
	}
	return defaultValue
}

func ToInt8P(val interface{}) int8 {
	v, err := ToInt8E(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToUint64(val interface{}) uint64 {
	if v, err := ToUint64E(val); err == nil {
		return v
	}
	var v uint64
	return v
}

func ToUint64D(val interface{}, defaultValue uint64) uint64 {
	if v, err := ToUint64E(val); err == nil {
		return v
	}
	return defaultValue
}

func ToUint64P(val interface{}) uint64 {
	v, err := ToUint64E(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToUint32(val interface{}) uint32 {
	if v, err := ToUint32E(val); err == nil {
		return v
	}
	var v uint32
	return v
}

func ToUint32D(val interface{}, defaultValue uint32) uint32 {
	if v, err := ToUint32E(val); err == nil {
		return v
	}
	return defaultValue
}

func ToUint32P(val interface{}) uint32 {
	v, err := ToUint32E(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToUint16(val interface{}) uint16 {
	if v, err := ToUint16E(val); err == nil {
		return v
	}
	var v uint16
	return v
}

func ToUint16D(val interface{}, defaultValue uint16) uint16 {
	if v, err := ToUint16E(val); err == nil {
		return v
	}
	return defaultValue
}

func ToUint16P(val interface{}) uint16 {
	v, err := ToUint16E(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToUint8(val interface{}) uint8 {
	if v, err := ToUint8E(val); err == nil {
		return v
	}
	var v uint8
	return v
}

func ToUint8D(val interface{}, defaultValue uint8) uint8 {
	if v, err := ToUint8E(val); err == nil {
		return v
	}
	return defaultValue
}

func ToUint8P(val interface{}) uint8 {
	v, err := ToUint8E(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToFloat64(val interface{}) float64 {
	if v, err := ToFloat64E(val); err == nil {
		return v
	}
	var v float64
	return v
}

func ToFloat64D(val interface{}, defaultValue float64) float64 {
	if v, err := ToFloat64E(val); err == nil {
		return v
	}
	return defaultValue
}

func ToFloat64P(val interface{}) float64 {
	v, err := ToFloat64E(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToFloat32(val interface{}) float32 {
	if v, err := ToFloat32E(val); err == nil {
		return v
	}
	var v float32
	return v
}

func ToFloat32D(val interface{}, defaultValue float32) float32 {
	if v, err := ToFloat32E(val); err == nil {
		return v
	}
	return defaultValue
}

func ToFloat32P(val interface{}) float32 {
	v, err := ToFloat32E(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToString(val interface{}) string {
	if v, err := ToStringE(val); err == nil {
		return v
	}
	var v string
	return v
}

func ToStringD(val interface{}, defaultValue string) string {
	if v, err := ToStringE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToStringP(val interface{}) string {
	v, err := ToStringE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToDuration(val interface{}) time.Duration {
	if v, err := ToDurationE(val); err == nil {
		return v
	}
	var v time.Duration
	return v
}

func ToDurationD(val interface{}, defaultValue time.Duration) time.Duration {
	if v, err := ToDurationE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToDurationP(val interface{}) time.Duration {
	v, err := ToDurationE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToTime(val interface{}) time.Time {
	if v, err := ToTimeE(val); err == nil {
		return v
	}
	var v time.Time
	return v
}

func ToTimeD(val interface{}, defaultValue time.Time) time.Time {
	if v, err := ToTimeE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToTimeP(val interface{}) time.Time {
	v, err := ToTimeE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToIP(val interface{}) net.IP {
	if v, err := ToIPE(val); err == nil {
		return v
	}
	var v net.IP
	return v
}

func ToIPD(val interface{}, defaultValue net.IP) net.IP {
	if v, err := ToIPE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToIPP(val interface{}) net.IP {
	v, err := ToIPE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToRegex(val interface{}) *regexp.Regexp {
	if v, err := ToRegexE(val); err == nil {
		return v
	}
	var v *regexp.Regexp
	return v
}

func ToRegexD(val interface{}, defaultValue *regexp.Regexp) *regexp.Regexp {
	if v, err := ToRegexE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToRegexP(val interface{}) *regexp.Regexp {
	v, err := ToRegexE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToMapStringString(val interface{}) map[string]string {
	if v, err := ToMapStringStringE(val); err == nil {
		return v
	}
	var v map[string]string
	return v
}

func ToMapStringStringD(val interface{}, defaultValue map[string]string) map[string]string {
	if v, err := ToMapStringStringE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToMapStringStringP(val interface{}) map[string]string {
	v, err := ToMapStringStringE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToBoolSlice(val interface{}) []bool {
	if v, err := ToBoolSliceE(val); err == nil {
		return v
	}
	var v []bool
	return v
}

func ToBoolSliceD(val interface{}, defaultValue []bool) []bool {
	if v, err := ToBoolSliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToBoolSliceP(val interface{}) []bool {
	v, err := ToBoolSliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToIntSlice(val interface{}) []int {
	if v, err := ToIntSliceE(val); err == nil {
		return v
	}
	var v []int
	return v
}

func ToIntSliceD(val interface{}, defaultValue []int) []int {
	if v, err := ToIntSliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToIntSliceP(val interface{}) []int {
	v, err := ToIntSliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToUintSlice(val interface{}) []uint {
	if v, err := ToUintSliceE(val); err == nil {
		return v
	}
	var v []uint
	return v
}

func ToUintSliceD(val interface{}, defaultValue []uint) []uint {
	if v, err := ToUintSliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToUintSliceP(val interface{}) []uint {
	v, err := ToUintSliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToInt64Slice(val interface{}) []int64 {
	if v, err := ToInt64SliceE(val); err == nil {
		return v
	}
	var v []int64
	return v
}

func ToInt64SliceD(val interface{}, defaultValue []int64) []int64 {
	if v, err := ToInt64SliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToInt64SliceP(val interface{}) []int64 {
	v, err := ToInt64SliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToInt32Slice(val interface{}) []int32 {
	if v, err := ToInt32SliceE(val); err == nil {
		return v
	}
	var v []int32
	return v
}

func ToInt32SliceD(val interface{}, defaultValue []int32) []int32 {
	if v, err := ToInt32SliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToInt32SliceP(val interface{}) []int32 {
	v, err := ToInt32SliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToInt16Slice(val interface{}) []int16 {
	if v, err := ToInt16SliceE(val); err == nil {
		return v
	}
	var v []int16
	return v
}

func ToInt16SliceD(val interface{}, defaultValue []int16) []int16 {
	if v, err := ToInt16SliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToInt16SliceP(val interface{}) []int16 {
	v, err := ToInt16SliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToInt8Slice(val interface{}) []int8 {
	if v, err := ToInt8SliceE(val); err == nil {
		return v
	}
	var v []int8
	return v
}

func ToInt8SliceD(val interface{}, defaultValue []int8) []int8 {
	if v, err := ToInt8SliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToInt8SliceP(val interface{}) []int8 {
	v, err := ToInt8SliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToUint64Slice(val interface{}) []uint64 {
	if v, err := ToUint64SliceE(val); err == nil {
		return v
	}
	var v []uint64
	return v
}

func ToUint64SliceD(val interface{}, defaultValue []uint64) []uint64 {
	if v, err := ToUint64SliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToUint64SliceP(val interface{}) []uint64 {
	v, err := ToUint64SliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToUint32Slice(val interface{}) []uint32 {
	if v, err := ToUint32SliceE(val); err == nil {
		return v
	}
	var v []uint32
	return v
}

func ToUint32SliceD(val interface{}, defaultValue []uint32) []uint32 {
	if v, err := ToUint32SliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToUint32SliceP(val interface{}) []uint32 {
	v, err := ToUint32SliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToUint16Slice(val interface{}) []uint16 {
	if v, err := ToUint16SliceE(val); err == nil {
		return v
	}
	var v []uint16
	return v
}

func ToUint16SliceD(val interface{}, defaultValue []uint16) []uint16 {
	if v, err := ToUint16SliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToUint16SliceP(val interface{}) []uint16 {
	v, err := ToUint16SliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToUint8Slice(val interface{}) []uint8 {
	if v, err := ToUint8SliceE(val); err == nil {
		return v
	}
	var v []uint8
	return v
}

func ToUint8SliceD(val interface{}, defaultValue []uint8) []uint8 {
	if v, err := ToUint8SliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToUint8SliceP(val interface{}) []uint8 {
	v, err := ToUint8SliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToFloat64Slice(val interface{}) []float64 {
	if v, err := ToFloat64SliceE(val); err == nil {
		return v
	}
	var v []float64
	return v
}

func ToFloat64SliceD(val interface{}, defaultValue []float64) []float64 {
	if v, err := ToFloat64SliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToFloat64SliceP(val interface{}) []float64 {
	v, err := ToFloat64SliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToFloat32Slice(val interface{}) []float32 {
	if v, err := ToFloat32SliceE(val); err == nil {
		return v
	}
	var v []float32
	return v
}

func ToFloat32SliceD(val interface{}, defaultValue []float32) []float32 {
	if v, err := ToFloat32SliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToFloat32SliceP(val interface{}) []float32 {
	v, err := ToFloat32SliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToStringSlice(val interface{}) []string {
	if v, err := ToStringSliceE(val); err == nil {
		return v
	}
	var v []string
	return v
}

func ToStringSliceD(val interface{}, defaultValue []string) []string {
	if v, err := ToStringSliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToStringSliceP(val interface{}) []string {
	v, err := ToStringSliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToDurationSlice(val interface{}) []time.Duration {
	if v, err := ToDurationSliceE(val); err == nil {
		return v
	}
	var v []time.Duration
	return v
}

func ToDurationSliceD(val interface{}, defaultValue []time.Duration) []time.Duration {
	if v, err := ToDurationSliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToDurationSliceP(val interface{}) []time.Duration {
	v, err := ToDurationSliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToTimeSlice(val interface{}) []time.Time {
	if v, err := ToTimeSliceE(val); err == nil {
		return v
	}
	var v []time.Time
	return v
}

func ToTimeSliceD(val interface{}, defaultValue []time.Time) []time.Time {
	if v, err := ToTimeSliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToTimeSliceP(val interface{}) []time.Time {
	v, err := ToTimeSliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToIPSlice(val interface{}) []net.IP {
	if v, err := ToIPSliceE(val); err == nil {
		return v
	}
	var v []net.IP
	return v
}

func ToIPSliceD(val interface{}, defaultValue []net.IP) []net.IP {
	if v, err := ToIPSliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToIPSliceP(val interface{}) []net.IP {
	v, err := ToIPSliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToRegexSlice(val interface{}) []*regexp.Regexp {
	if v, err := ToRegexSliceE(val); err == nil {
		return v
	}
	var v []*regexp.Regexp
	return v
}

func ToRegexSliceD(val interface{}, defaultValue []*regexp.Regexp) []*regexp.Regexp {
	if v, err := ToRegexSliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToRegexSliceP(val interface{}) []*regexp.Regexp {
	v, err := ToRegexSliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}

func ToMapStringStringSlice(val interface{}) []map[string]string {
	if v, err := ToMapStringStringSliceE(val); err == nil {
		return v
	}
	var v []map[string]string
	return v
}

func ToMapStringStringSliceD(val interface{}, defaultValue []map[string]string) []map[string]string {
	if v, err := ToMapStringStringSliceE(val); err == nil {
		return v
	}
	return defaultValue
}

func ToMapStringStringSliceP(val interface{}) []map[string]string {
	v, err := ToMapStringStringSliceE(val)
	if err != nil {
		panic(err)
	}
	return v
}
