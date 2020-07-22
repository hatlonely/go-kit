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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {
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
