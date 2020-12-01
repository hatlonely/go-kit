// this file is auto generate by autogen.py. do not edit!
package flag

import (
	"net"
	"reflect"
	"time"

	"github.com/hatlonely/go-kit/cast"
)

func (f *Flag) BoolVar(p *bool, name string, defaultValue bool, usage string) {
	var v bool
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Bool(name string, defaultValue bool, usage string) *bool {
	var v bool
	f.BoolVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) IntVar(p *int, name string, defaultValue int, usage string) {
	var v int
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Int(name string, defaultValue int, usage string) *int {
	var v int
	f.IntVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) UintVar(p *uint, name string, defaultValue uint, usage string) {
	var v uint
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Uint(name string, defaultValue uint, usage string) *uint {
	var v uint
	f.UintVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Int64Var(p *int64, name string, defaultValue int64, usage string) {
	var v int64
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Int64(name string, defaultValue int64, usage string) *int64 {
	var v int64
	f.Int64Var(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Int32Var(p *int32, name string, defaultValue int32, usage string) {
	var v int32
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Int32(name string, defaultValue int32, usage string) *int32 {
	var v int32
	f.Int32Var(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Int16Var(p *int16, name string, defaultValue int16, usage string) {
	var v int16
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Int16(name string, defaultValue int16, usage string) *int16 {
	var v int16
	f.Int16Var(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Int8Var(p *int8, name string, defaultValue int8, usage string) {
	var v int8
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Int8(name string, defaultValue int8, usage string) *int8 {
	var v int8
	f.Int8Var(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Uint64Var(p *uint64, name string, defaultValue uint64, usage string) {
	var v uint64
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Uint64(name string, defaultValue uint64, usage string) *uint64 {
	var v uint64
	f.Uint64Var(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Uint32Var(p *uint32, name string, defaultValue uint32, usage string) {
	var v uint32
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Uint32(name string, defaultValue uint32, usage string) *uint32 {
	var v uint32
	f.Uint32Var(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Uint16Var(p *uint16, name string, defaultValue uint16, usage string) {
	var v uint16
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Uint16(name string, defaultValue uint16, usage string) *uint16 {
	var v uint16
	f.Uint16Var(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Uint8Var(p *uint8, name string, defaultValue uint8, usage string) {
	var v uint8
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Uint8(name string, defaultValue uint8, usage string) *uint8 {
	var v uint8
	f.Uint8Var(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Float64Var(p *float64, name string, defaultValue float64, usage string) {
	var v float64
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Float64(name string, defaultValue float64, usage string) *float64 {
	var v float64
	f.Float64Var(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Float32Var(p *float32, name string, defaultValue float32, usage string) {
	var v float32
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Float32(name string, defaultValue float32, usage string) *float32 {
	var v float32
	f.Float32Var(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) StringVar(p *string, name string, defaultValue string, usage string) {
	var v string
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) String(name string, defaultValue string, usage string) *string {
	var v string
	f.StringVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) DurationVar(p *time.Duration, name string, defaultValue time.Duration, usage string) {
	var v time.Duration
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Duration(name string, defaultValue time.Duration, usage string) *time.Duration {
	var v time.Duration
	f.DurationVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) TimeVar(p *time.Time, name string, defaultValue time.Time, usage string) {
	var v time.Time
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Time(name string, defaultValue time.Time, usage string) *time.Time {
	var v time.Time
	f.TimeVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) IPVar(p *net.IP, name string, defaultValue net.IP, usage string) {
	var v net.IP
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) IP(name string, defaultValue net.IP, usage string) *net.IP {
	var v net.IP
	f.IPVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) BoolSliceVar(p *[]bool, name string, defaultValue []bool, usage string) {
	var v []bool
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) BoolSlice(name string, defaultValue []bool, usage string) *[]bool {
	var v []bool
	f.BoolSliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) IntSliceVar(p *[]int, name string, defaultValue []int, usage string) {
	var v []int
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) IntSlice(name string, defaultValue []int, usage string) *[]int {
	var v []int
	f.IntSliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) UintSliceVar(p *[]uint, name string, defaultValue []uint, usage string) {
	var v []uint
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) UintSlice(name string, defaultValue []uint, usage string) *[]uint {
	var v []uint
	f.UintSliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Int64SliceVar(p *[]int64, name string, defaultValue []int64, usage string) {
	var v []int64
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Int64Slice(name string, defaultValue []int64, usage string) *[]int64 {
	var v []int64
	f.Int64SliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Int32SliceVar(p *[]int32, name string, defaultValue []int32, usage string) {
	var v []int32
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Int32Slice(name string, defaultValue []int32, usage string) *[]int32 {
	var v []int32
	f.Int32SliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Int16SliceVar(p *[]int16, name string, defaultValue []int16, usage string) {
	var v []int16
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Int16Slice(name string, defaultValue []int16, usage string) *[]int16 {
	var v []int16
	f.Int16SliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Int8SliceVar(p *[]int8, name string, defaultValue []int8, usage string) {
	var v []int8
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Int8Slice(name string, defaultValue []int8, usage string) *[]int8 {
	var v []int8
	f.Int8SliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Uint64SliceVar(p *[]uint64, name string, defaultValue []uint64, usage string) {
	var v []uint64
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Uint64Slice(name string, defaultValue []uint64, usage string) *[]uint64 {
	var v []uint64
	f.Uint64SliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Uint32SliceVar(p *[]uint32, name string, defaultValue []uint32, usage string) {
	var v []uint32
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Uint32Slice(name string, defaultValue []uint32, usage string) *[]uint32 {
	var v []uint32
	f.Uint32SliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Uint16SliceVar(p *[]uint16, name string, defaultValue []uint16, usage string) {
	var v []uint16
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Uint16Slice(name string, defaultValue []uint16, usage string) *[]uint16 {
	var v []uint16
	f.Uint16SliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Uint8SliceVar(p *[]uint8, name string, defaultValue []uint8, usage string) {
	var v []uint8
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Uint8Slice(name string, defaultValue []uint8, usage string) *[]uint8 {
	var v []uint8
	f.Uint8SliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Float64SliceVar(p *[]float64, name string, defaultValue []float64, usage string) {
	var v []float64
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Float64Slice(name string, defaultValue []float64, usage string) *[]float64 {
	var v []float64
	f.Float64SliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) Float32SliceVar(p *[]float32, name string, defaultValue []float32, usage string) {
	var v []float32
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) Float32Slice(name string, defaultValue []float32, usage string) *[]float32 {
	var v []float32
	f.Float32SliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) StringSliceVar(p *[]string, name string, defaultValue []string, usage string) {
	var v []string
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) StringSlice(name string, defaultValue []string, usage string) *[]string {
	var v []string
	f.StringSliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) DurationSliceVar(p *[]time.Duration, name string, defaultValue []time.Duration, usage string) {
	var v []time.Duration
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) DurationSlice(name string, defaultValue []time.Duration, usage string) *[]time.Duration {
	var v []time.Duration
	f.DurationSliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) TimeSliceVar(p *[]time.Time, name string, defaultValue []time.Time, usage string) {
	var v []time.Time
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) TimeSlice(name string, defaultValue []time.Time, usage string) *[]time.Time {
	var v []time.Time
	f.TimeSliceVar(&v, name, defaultValue, usage)
	return &v
}

func (f *Flag) IPSliceVar(p *[]net.IP, name string, defaultValue []net.IP, usage string) {
	var v []net.IP
	*p = defaultValue
	if err := f.addFlagWithOptions(&AddOptionOptions{
		Name:         name,
		Usage:        usage,
		Type:         reflect.TypeOf(v),
		DefaultValue: cast.ToString(defaultValue),
	}); err != nil {
		panic(err)
	}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {
		return cast.SetInterface(p, val)
	}
}

func (f *Flag) IPSlice(name string, defaultValue []net.IP, usage string) *[]net.IP {
	var v []net.IP
	f.IPSliceVar(&v, name, defaultValue, usage)
	return &v
}
