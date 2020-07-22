// this file is auto generate by autogen.py. do not edit!
package flag

import (
	"net"
	"os"
	"reflect"
	"time"
)

var gflag = NewFlag(os.Args[0])

func Parse() error {
	return ParseArgs(os.Args[1:])
}

func GetInfo(key string) (*Info, bool) {
	return gflag.GetInfo(key)
}

func Set(key string, val interface{}) error {
	return gflag.Set(key, val)
}

func Get(key string, val string) (interface{}, bool) {
	return gflag.Get(key, val)
}

func AddFlag(name string, usage string, opts ...AddFlagOption) error {
	return gflag.AddFlag(name, usage, opts...)
}

func AddArgument(name string, usage string, opts ...AddFlagOption) error {
	return gflag.AddArgument(name, usage, opts...)
}

func BindFlag(v interface{}, name string, usage string, rtype reflect.Type, required bool, shorthand string, defaultValue string, isArgument bool) error {
	return gflag.BindFlag(v, name, usage, rtype, required, shorthand, defaultValue, isArgument)
}

func Bind(v interface{}) error {
	return gflag.Bind(v)
}

func Usage() string {
	return gflag.Usage()
}

func BoolVar(p *bool, name string, defaultValue bool, usage string) {
	gflag.BoolVar(p, name, defaultValue, usage)
}

func Bool(name string, defaultValue bool, usage string) *bool {
	return gflag.Bool(name, defaultValue, usage)
}

func IntVar(p *int, name string, defaultValue int, usage string) {
	gflag.IntVar(p, name, defaultValue, usage)
}

func Int(name string, defaultValue int, usage string) *int {
	return gflag.Int(name, defaultValue, usage)
}

func UintVar(p *uint, name string, defaultValue uint, usage string) {
	gflag.UintVar(p, name, defaultValue, usage)
}

func Uint(name string, defaultValue uint, usage string) *uint {
	return gflag.Uint(name, defaultValue, usage)
}

func Int64Var(p *int64, name string, defaultValue int64, usage string) {
	gflag.Int64Var(p, name, defaultValue, usage)
}

func Int64(name string, defaultValue int64, usage string) *int64 {
	return gflag.Int64(name, defaultValue, usage)
}

func Int32Var(p *int32, name string, defaultValue int32, usage string) {
	gflag.Int32Var(p, name, defaultValue, usage)
}

func Int32(name string, defaultValue int32, usage string) *int32 {
	return gflag.Int32(name, defaultValue, usage)
}

func Int16Var(p *int16, name string, defaultValue int16, usage string) {
	gflag.Int16Var(p, name, defaultValue, usage)
}

func Int16(name string, defaultValue int16, usage string) *int16 {
	return gflag.Int16(name, defaultValue, usage)
}

func Int8Var(p *int8, name string, defaultValue int8, usage string) {
	gflag.Int8Var(p, name, defaultValue, usage)
}

func Int8(name string, defaultValue int8, usage string) *int8 {
	return gflag.Int8(name, defaultValue, usage)
}

func Uint64Var(p *uint64, name string, defaultValue uint64, usage string) {
	gflag.Uint64Var(p, name, defaultValue, usage)
}

func Uint64(name string, defaultValue uint64, usage string) *uint64 {
	return gflag.Uint64(name, defaultValue, usage)
}

func Uint32Var(p *uint32, name string, defaultValue uint32, usage string) {
	gflag.Uint32Var(p, name, defaultValue, usage)
}

func Uint32(name string, defaultValue uint32, usage string) *uint32 {
	return gflag.Uint32(name, defaultValue, usage)
}

func Uint16Var(p *uint16, name string, defaultValue uint16, usage string) {
	gflag.Uint16Var(p, name, defaultValue, usage)
}

func Uint16(name string, defaultValue uint16, usage string) *uint16 {
	return gflag.Uint16(name, defaultValue, usage)
}

func Uint8Var(p *uint8, name string, defaultValue uint8, usage string) {
	gflag.Uint8Var(p, name, defaultValue, usage)
}

func Uint8(name string, defaultValue uint8, usage string) *uint8 {
	return gflag.Uint8(name, defaultValue, usage)
}

func Float64Var(p *float64, name string, defaultValue float64, usage string) {
	gflag.Float64Var(p, name, defaultValue, usage)
}

func Float64(name string, defaultValue float64, usage string) *float64 {
	return gflag.Float64(name, defaultValue, usage)
}

func Float32Var(p *float32, name string, defaultValue float32, usage string) {
	gflag.Float32Var(p, name, defaultValue, usage)
}

func Float32(name string, defaultValue float32, usage string) *float32 {
	return gflag.Float32(name, defaultValue, usage)
}

func StringVar(p *string, name string, defaultValue string, usage string) {
	gflag.StringVar(p, name, defaultValue, usage)
}

func String(name string, defaultValue string, usage string) *string {
	return gflag.String(name, defaultValue, usage)
}

func DurationVar(p *time.Duration, name string, defaultValue time.Duration, usage string) {
	gflag.DurationVar(p, name, defaultValue, usage)
}

func Duration(name string, defaultValue time.Duration, usage string) *time.Duration {
	return gflag.Duration(name, defaultValue, usage)
}

func TimeVar(p *time.Time, name string, defaultValue time.Time, usage string) {
	gflag.TimeVar(p, name, defaultValue, usage)
}

func Time(name string, defaultValue time.Time, usage string) *time.Time {
	return gflag.Time(name, defaultValue, usage)
}

func IPVar(p *net.IP, name string, defaultValue net.IP, usage string) {
	gflag.IPVar(p, name, defaultValue, usage)
}

func IP(name string, defaultValue net.IP, usage string) *net.IP {
	return gflag.IP(name, defaultValue, usage)
}

func BoolSliceVar(p *[]bool, name string, defaultValue []bool, usage string) {
	gflag.BoolSliceVar(p, name, defaultValue, usage)
}

func BoolSlice(name string, defaultValue []bool, usage string) *[]bool {
	return gflag.BoolSlice(name, defaultValue, usage)
}

func IntSliceVar(p *[]int, name string, defaultValue []int, usage string) {
	gflag.IntSliceVar(p, name, defaultValue, usage)
}

func IntSlice(name string, defaultValue []int, usage string) *[]int {
	return gflag.IntSlice(name, defaultValue, usage)
}

func UintSliceVar(p *[]uint, name string, defaultValue []uint, usage string) {
	gflag.UintSliceVar(p, name, defaultValue, usage)
}

func UintSlice(name string, defaultValue []uint, usage string) *[]uint {
	return gflag.UintSlice(name, defaultValue, usage)
}

func Int64SliceVar(p *[]int64, name string, defaultValue []int64, usage string) {
	gflag.Int64SliceVar(p, name, defaultValue, usage)
}

func Int64Slice(name string, defaultValue []int64, usage string) *[]int64 {
	return gflag.Int64Slice(name, defaultValue, usage)
}

func Int32SliceVar(p *[]int32, name string, defaultValue []int32, usage string) {
	gflag.Int32SliceVar(p, name, defaultValue, usage)
}

func Int32Slice(name string, defaultValue []int32, usage string) *[]int32 {
	return gflag.Int32Slice(name, defaultValue, usage)
}

func Int16SliceVar(p *[]int16, name string, defaultValue []int16, usage string) {
	gflag.Int16SliceVar(p, name, defaultValue, usage)
}

func Int16Slice(name string, defaultValue []int16, usage string) *[]int16 {
	return gflag.Int16Slice(name, defaultValue, usage)
}

func Int8SliceVar(p *[]int8, name string, defaultValue []int8, usage string) {
	gflag.Int8SliceVar(p, name, defaultValue, usage)
}

func Int8Slice(name string, defaultValue []int8, usage string) *[]int8 {
	return gflag.Int8Slice(name, defaultValue, usage)
}

func Uint64SliceVar(p *[]uint64, name string, defaultValue []uint64, usage string) {
	gflag.Uint64SliceVar(p, name, defaultValue, usage)
}

func Uint64Slice(name string, defaultValue []uint64, usage string) *[]uint64 {
	return gflag.Uint64Slice(name, defaultValue, usage)
}

func Uint32SliceVar(p *[]uint32, name string, defaultValue []uint32, usage string) {
	gflag.Uint32SliceVar(p, name, defaultValue, usage)
}

func Uint32Slice(name string, defaultValue []uint32, usage string) *[]uint32 {
	return gflag.Uint32Slice(name, defaultValue, usage)
}

func Uint16SliceVar(p *[]uint16, name string, defaultValue []uint16, usage string) {
	gflag.Uint16SliceVar(p, name, defaultValue, usage)
}

func Uint16Slice(name string, defaultValue []uint16, usage string) *[]uint16 {
	return gflag.Uint16Slice(name, defaultValue, usage)
}

func Uint8SliceVar(p *[]uint8, name string, defaultValue []uint8, usage string) {
	gflag.Uint8SliceVar(p, name, defaultValue, usage)
}

func Uint8Slice(name string, defaultValue []uint8, usage string) *[]uint8 {
	return gflag.Uint8Slice(name, defaultValue, usage)
}

func Float64SliceVar(p *[]float64, name string, defaultValue []float64, usage string) {
	gflag.Float64SliceVar(p, name, defaultValue, usage)
}

func Float64Slice(name string, defaultValue []float64, usage string) *[]float64 {
	return gflag.Float64Slice(name, defaultValue, usage)
}

func Float32SliceVar(p *[]float32, name string, defaultValue []float32, usage string) {
	gflag.Float32SliceVar(p, name, defaultValue, usage)
}

func Float32Slice(name string, defaultValue []float32, usage string) *[]float32 {
	return gflag.Float32Slice(name, defaultValue, usage)
}

func StringSliceVar(p *[]string, name string, defaultValue []string, usage string) {
	gflag.StringSliceVar(p, name, defaultValue, usage)
}

func StringSlice(name string, defaultValue []string, usage string) *[]string {
	return gflag.StringSlice(name, defaultValue, usage)
}

func DurationSliceVar(p *[]time.Duration, name string, defaultValue []time.Duration, usage string) {
	gflag.DurationSliceVar(p, name, defaultValue, usage)
}

func DurationSlice(name string, defaultValue []time.Duration, usage string) *[]time.Duration {
	return gflag.DurationSlice(name, defaultValue, usage)
}

func TimeSliceVar(p *[]time.Time, name string, defaultValue []time.Time, usage string) {
	gflag.TimeSliceVar(p, name, defaultValue, usage)
}

func TimeSlice(name string, defaultValue []time.Time, usage string) *[]time.Time {
	return gflag.TimeSlice(name, defaultValue, usage)
}

func IPSliceVar(p *[]net.IP, name string, defaultValue []net.IP, usage string) {
	gflag.IPSliceVar(p, name, defaultValue, usage)
}

func IPSlice(name string, defaultValue []net.IP, usage string) *[]net.IP {
	return gflag.IPSlice(name, defaultValue, usage)
}

func ParseArgs(args []string) error {
	return gflag.ParseArgs(args)
}
