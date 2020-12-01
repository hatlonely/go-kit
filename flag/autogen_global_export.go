// this file is auto generate by autogen.py. do not edit!
package flag

import (
	"net"
	"os"
	"time"

	"github.com/hatlonely/go-kit/refx"
)

var gflag = NewFlag(os.Args[0])

func Instance() *Flag {
	return gflag
}

func Parse(opts ...ParseOption) error {
	return ParseArgs(os.Args[1:], opts...)
}

func GetInfo(key string) (*Info, bool) {
	return gflag.GetInfo(key)
}

func Set(key string, val string, opts ...ParseOption) error {
	return gflag.Set(key, val, opts...)
}

func SetWithOptions(key string, val string, options *ParseOptions) error {
	return gflag.SetWithOptions(key, val, options)
}

func Get(key string) (interface{}, bool) {
	return gflag.Get(key)
}

func GetBoolE(key string) (bool, error) {
	return gflag.GetBoolE(key)
}

func GetBool(key string) bool {
	return gflag.GetBool(key)
}

func GetBoolD(key string, val bool) bool {
	return gflag.GetBoolD(key, val)
}

func GetBoolP(key string) bool {
	return gflag.GetBoolP(key)
}

func GetIntE(key string) (int, error) {
	return gflag.GetIntE(key)
}

func GetInt(key string) int {
	return gflag.GetInt(key)
}

func GetIntD(key string, val int) int {
	return gflag.GetIntD(key, val)
}

func GetIntP(key string) int {
	return gflag.GetIntP(key)
}

func GetUintE(key string) (uint, error) {
	return gflag.GetUintE(key)
}

func GetUint(key string) uint {
	return gflag.GetUint(key)
}

func GetUintD(key string, val uint) uint {
	return gflag.GetUintD(key, val)
}

func GetUintP(key string) uint {
	return gflag.GetUintP(key)
}

func GetInt64E(key string) (int64, error) {
	return gflag.GetInt64E(key)
}

func GetInt64(key string) int64 {
	return gflag.GetInt64(key)
}

func GetInt64D(key string, val int64) int64 {
	return gflag.GetInt64D(key, val)
}

func GetInt64P(key string) int64 {
	return gflag.GetInt64P(key)
}

func GetInt32E(key string) (int32, error) {
	return gflag.GetInt32E(key)
}

func GetInt32(key string) int32 {
	return gflag.GetInt32(key)
}

func GetInt32D(key string, val int32) int32 {
	return gflag.GetInt32D(key, val)
}

func GetInt32P(key string) int32 {
	return gflag.GetInt32P(key)
}

func GetInt16E(key string) (int16, error) {
	return gflag.GetInt16E(key)
}

func GetInt16(key string) int16 {
	return gflag.GetInt16(key)
}

func GetInt16D(key string, val int16) int16 {
	return gflag.GetInt16D(key, val)
}

func GetInt16P(key string) int16 {
	return gflag.GetInt16P(key)
}

func GetInt8E(key string) (int8, error) {
	return gflag.GetInt8E(key)
}

func GetInt8(key string) int8 {
	return gflag.GetInt8(key)
}

func GetInt8D(key string, val int8) int8 {
	return gflag.GetInt8D(key, val)
}

func GetInt8P(key string) int8 {
	return gflag.GetInt8P(key)
}

func GetUint64E(key string) (uint64, error) {
	return gflag.GetUint64E(key)
}

func GetUint64(key string) uint64 {
	return gflag.GetUint64(key)
}

func GetUint64D(key string, val uint64) uint64 {
	return gflag.GetUint64D(key, val)
}

func GetUint64P(key string) uint64 {
	return gflag.GetUint64P(key)
}

func GetUint32E(key string) (uint32, error) {
	return gflag.GetUint32E(key)
}

func GetUint32(key string) uint32 {
	return gflag.GetUint32(key)
}

func GetUint32D(key string, val uint32) uint32 {
	return gflag.GetUint32D(key, val)
}

func GetUint32P(key string) uint32 {
	return gflag.GetUint32P(key)
}

func GetUint16E(key string) (uint16, error) {
	return gflag.GetUint16E(key)
}

func GetUint16(key string) uint16 {
	return gflag.GetUint16(key)
}

func GetUint16D(key string, val uint16) uint16 {
	return gflag.GetUint16D(key, val)
}

func GetUint16P(key string) uint16 {
	return gflag.GetUint16P(key)
}

func GetUint8E(key string) (uint8, error) {
	return gflag.GetUint8E(key)
}

func GetUint8(key string) uint8 {
	return gflag.GetUint8(key)
}

func GetUint8D(key string, val uint8) uint8 {
	return gflag.GetUint8D(key, val)
}

func GetUint8P(key string) uint8 {
	return gflag.GetUint8P(key)
}

func GetFloat64E(key string) (float64, error) {
	return gflag.GetFloat64E(key)
}

func GetFloat64(key string) float64 {
	return gflag.GetFloat64(key)
}

func GetFloat64D(key string, val float64) float64 {
	return gflag.GetFloat64D(key, val)
}

func GetFloat64P(key string) float64 {
	return gflag.GetFloat64P(key)
}

func GetFloat32E(key string) (float32, error) {
	return gflag.GetFloat32E(key)
}

func GetFloat32(key string) float32 {
	return gflag.GetFloat32(key)
}

func GetFloat32D(key string, val float32) float32 {
	return gflag.GetFloat32D(key, val)
}

func GetFloat32P(key string) float32 {
	return gflag.GetFloat32P(key)
}

func GetStringE(key string) (string, error) {
	return gflag.GetStringE(key)
}

func GetString(key string) string {
	return gflag.GetString(key)
}

func GetStringD(key string, val string) string {
	return gflag.GetStringD(key, val)
}

func GetStringP(key string) string {
	return gflag.GetStringP(key)
}

func GetDurationE(key string) (time.Duration, error) {
	return gflag.GetDurationE(key)
}

func GetDuration(key string) time.Duration {
	return gflag.GetDuration(key)
}

func GetDurationD(key string, val time.Duration) time.Duration {
	return gflag.GetDurationD(key, val)
}

func GetDurationP(key string) time.Duration {
	return gflag.GetDurationP(key)
}

func GetTimeE(key string) (time.Time, error) {
	return gflag.GetTimeE(key)
}

func GetTime(key string) time.Time {
	return gflag.GetTime(key)
}

func GetTimeD(key string, val time.Time) time.Time {
	return gflag.GetTimeD(key, val)
}

func GetTimeP(key string) time.Time {
	return gflag.GetTimeP(key)
}

func GetIPE(key string) (net.IP, error) {
	return gflag.GetIPE(key)
}

func GetIP(key string) net.IP {
	return gflag.GetIP(key)
}

func GetIPD(key string, val net.IP) net.IP {
	return gflag.GetIPD(key, val)
}

func GetIPP(key string) net.IP {
	return gflag.GetIPP(key)
}

func GetBoolSliceE(key string) ([]bool, error) {
	return gflag.GetBoolSliceE(key)
}

func GetBoolSlice(key string) []bool {
	return gflag.GetBoolSlice(key)
}

func GetBoolSliceD(key string, val []bool) []bool {
	return gflag.GetBoolSliceD(key, val)
}

func GetBoolSliceP(key string) []bool {
	return gflag.GetBoolSliceP(key)
}

func GetIntSliceE(key string) ([]int, error) {
	return gflag.GetIntSliceE(key)
}

func GetIntSlice(key string) []int {
	return gflag.GetIntSlice(key)
}

func GetIntSliceD(key string, val []int) []int {
	return gflag.GetIntSliceD(key, val)
}

func GetIntSliceP(key string) []int {
	return gflag.GetIntSliceP(key)
}

func GetUintSliceE(key string) ([]uint, error) {
	return gflag.GetUintSliceE(key)
}

func GetUintSlice(key string) []uint {
	return gflag.GetUintSlice(key)
}

func GetUintSliceD(key string, val []uint) []uint {
	return gflag.GetUintSliceD(key, val)
}

func GetUintSliceP(key string) []uint {
	return gflag.GetUintSliceP(key)
}

func GetInt64SliceE(key string) ([]int64, error) {
	return gflag.GetInt64SliceE(key)
}

func GetInt64Slice(key string) []int64 {
	return gflag.GetInt64Slice(key)
}

func GetInt64SliceD(key string, val []int64) []int64 {
	return gflag.GetInt64SliceD(key, val)
}

func GetInt64SliceP(key string) []int64 {
	return gflag.GetInt64SliceP(key)
}

func GetInt32SliceE(key string) ([]int32, error) {
	return gflag.GetInt32SliceE(key)
}

func GetInt32Slice(key string) []int32 {
	return gflag.GetInt32Slice(key)
}

func GetInt32SliceD(key string, val []int32) []int32 {
	return gflag.GetInt32SliceD(key, val)
}

func GetInt32SliceP(key string) []int32 {
	return gflag.GetInt32SliceP(key)
}

func GetInt16SliceE(key string) ([]int16, error) {
	return gflag.GetInt16SliceE(key)
}

func GetInt16Slice(key string) []int16 {
	return gflag.GetInt16Slice(key)
}

func GetInt16SliceD(key string, val []int16) []int16 {
	return gflag.GetInt16SliceD(key, val)
}

func GetInt16SliceP(key string) []int16 {
	return gflag.GetInt16SliceP(key)
}

func GetInt8SliceE(key string) ([]int8, error) {
	return gflag.GetInt8SliceE(key)
}

func GetInt8Slice(key string) []int8 {
	return gflag.GetInt8Slice(key)
}

func GetInt8SliceD(key string, val []int8) []int8 {
	return gflag.GetInt8SliceD(key, val)
}

func GetInt8SliceP(key string) []int8 {
	return gflag.GetInt8SliceP(key)
}

func GetUint64SliceE(key string) ([]uint64, error) {
	return gflag.GetUint64SliceE(key)
}

func GetUint64Slice(key string) []uint64 {
	return gflag.GetUint64Slice(key)
}

func GetUint64SliceD(key string, val []uint64) []uint64 {
	return gflag.GetUint64SliceD(key, val)
}

func GetUint64SliceP(key string) []uint64 {
	return gflag.GetUint64SliceP(key)
}

func GetUint32SliceE(key string) ([]uint32, error) {
	return gflag.GetUint32SliceE(key)
}

func GetUint32Slice(key string) []uint32 {
	return gflag.GetUint32Slice(key)
}

func GetUint32SliceD(key string, val []uint32) []uint32 {
	return gflag.GetUint32SliceD(key, val)
}

func GetUint32SliceP(key string) []uint32 {
	return gflag.GetUint32SliceP(key)
}

func GetUint16SliceE(key string) ([]uint16, error) {
	return gflag.GetUint16SliceE(key)
}

func GetUint16Slice(key string) []uint16 {
	return gflag.GetUint16Slice(key)
}

func GetUint16SliceD(key string, val []uint16) []uint16 {
	return gflag.GetUint16SliceD(key, val)
}

func GetUint16SliceP(key string) []uint16 {
	return gflag.GetUint16SliceP(key)
}

func GetUint8SliceE(key string) ([]uint8, error) {
	return gflag.GetUint8SliceE(key)
}

func GetUint8Slice(key string) []uint8 {
	return gflag.GetUint8Slice(key)
}

func GetUint8SliceD(key string, val []uint8) []uint8 {
	return gflag.GetUint8SliceD(key, val)
}

func GetUint8SliceP(key string) []uint8 {
	return gflag.GetUint8SliceP(key)
}

func GetFloat64SliceE(key string) ([]float64, error) {
	return gflag.GetFloat64SliceE(key)
}

func GetFloat64Slice(key string) []float64 {
	return gflag.GetFloat64Slice(key)
}

func GetFloat64SliceD(key string, val []float64) []float64 {
	return gflag.GetFloat64SliceD(key, val)
}

func GetFloat64SliceP(key string) []float64 {
	return gflag.GetFloat64SliceP(key)
}

func GetFloat32SliceE(key string) ([]float32, error) {
	return gflag.GetFloat32SliceE(key)
}

func GetFloat32Slice(key string) []float32 {
	return gflag.GetFloat32Slice(key)
}

func GetFloat32SliceD(key string, val []float32) []float32 {
	return gflag.GetFloat32SliceD(key, val)
}

func GetFloat32SliceP(key string) []float32 {
	return gflag.GetFloat32SliceP(key)
}

func GetStringSliceE(key string) ([]string, error) {
	return gflag.GetStringSliceE(key)
}

func GetStringSlice(key string) []string {
	return gflag.GetStringSlice(key)
}

func GetStringSliceD(key string, val []string) []string {
	return gflag.GetStringSliceD(key, val)
}

func GetStringSliceP(key string) []string {
	return gflag.GetStringSliceP(key)
}

func GetDurationSliceE(key string) ([]time.Duration, error) {
	return gflag.GetDurationSliceE(key)
}

func GetDurationSlice(key string) []time.Duration {
	return gflag.GetDurationSlice(key)
}

func GetDurationSliceD(key string, val []time.Duration) []time.Duration {
	return gflag.GetDurationSliceD(key, val)
}

func GetDurationSliceP(key string) []time.Duration {
	return gflag.GetDurationSliceP(key)
}

func GetTimeSliceE(key string) ([]time.Time, error) {
	return gflag.GetTimeSliceE(key)
}

func GetTimeSlice(key string) []time.Time {
	return gflag.GetTimeSlice(key)
}

func GetTimeSliceD(key string, val []time.Time) []time.Time {
	return gflag.GetTimeSliceD(key, val)
}

func GetTimeSliceP(key string) []time.Time {
	return gflag.GetTimeSliceP(key)
}

func GetIPSliceE(key string) ([]net.IP, error) {
	return gflag.GetIPSliceE(key)
}

func GetIPSlice(key string) []net.IP {
	return gflag.GetIPSlice(key)
}

func GetIPSliceD(key string, val []net.IP) []net.IP {
	return gflag.GetIPSliceD(key, val)
}

func GetIPSliceP(key string) []net.IP {
	return gflag.GetIPSliceP(key)
}

func AddOption(name string, usage string, opts ...AddOptionOption) {
	gflag.AddOption(name, usage, opts...)
}

func AddArgument(name string, usage string, opts ...AddOptionOption) {
	gflag.AddArgument(name, usage, opts...)
}

func BindFlagWithOptions(v interface{}, options *AddOptionOptions) error {
	return gflag.BindFlagWithOptions(v, options)
}

func Struct(v interface{}, opts ...refx.Option) error {
	return gflag.Struct(v, opts...)
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

func ParseArgs(args []string, opts ...ParseOption) error {
	return gflag.ParseArgs(args, opts...)
}
