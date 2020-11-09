// this file is auto generate by autogen.py. do not edit!
package config

import (
	"net"
	"sync/atomic"
	"time"
)

type AtomicBool struct {
	v atomic.Value
}

func NewAtomicBool(v bool) *AtomicBool {
	var av atomic.Value
	av.Store(v)
	return &AtomicBool{v: av}
}

func (a *AtomicBool) Get() bool {
	return a.v.Load().(bool)
}

func (a *AtomicBool) Set(v bool) {
	a.v.Store(v)
}

type AtomicInt struct {
	v atomic.Value
}

func NewAtomicInt(v int) *AtomicInt {
	var av atomic.Value
	av.Store(v)
	return &AtomicInt{v: av}
}

func (a *AtomicInt) Get() int {
	return a.v.Load().(int)
}

func (a *AtomicInt) Set(v int) {
	a.v.Store(v)
}

type AtomicUint struct {
	v atomic.Value
}

func NewAtomicUint(v uint) *AtomicUint {
	var av atomic.Value
	av.Store(v)
	return &AtomicUint{v: av}
}

func (a *AtomicUint) Get() uint {
	return a.v.Load().(uint)
}

func (a *AtomicUint) Set(v uint) {
	a.v.Store(v)
}

type AtomicInt64 struct {
	v atomic.Value
}

func NewAtomicInt64(v int64) *AtomicInt64 {
	var av atomic.Value
	av.Store(v)
	return &AtomicInt64{v: av}
}

func (a *AtomicInt64) Get() int64 {
	return a.v.Load().(int64)
}

func (a *AtomicInt64) Set(v int64) {
	a.v.Store(v)
}

type AtomicInt32 struct {
	v atomic.Value
}

func NewAtomicInt32(v int32) *AtomicInt32 {
	var av atomic.Value
	av.Store(v)
	return &AtomicInt32{v: av}
}

func (a *AtomicInt32) Get() int32 {
	return a.v.Load().(int32)
}

func (a *AtomicInt32) Set(v int32) {
	a.v.Store(v)
}

type AtomicInt16 struct {
	v atomic.Value
}

func NewAtomicInt16(v int16) *AtomicInt16 {
	var av atomic.Value
	av.Store(v)
	return &AtomicInt16{v: av}
}

func (a *AtomicInt16) Get() int16 {
	return a.v.Load().(int16)
}

func (a *AtomicInt16) Set(v int16) {
	a.v.Store(v)
}

type AtomicInt8 struct {
	v atomic.Value
}

func NewAtomicInt8(v int8) *AtomicInt8 {
	var av atomic.Value
	av.Store(v)
	return &AtomicInt8{v: av}
}

func (a *AtomicInt8) Get() int8 {
	return a.v.Load().(int8)
}

func (a *AtomicInt8) Set(v int8) {
	a.v.Store(v)
}

type AtomicUint64 struct {
	v atomic.Value
}

func NewAtomicUint64(v uint64) *AtomicUint64 {
	var av atomic.Value
	av.Store(v)
	return &AtomicUint64{v: av}
}

func (a *AtomicUint64) Get() uint64 {
	return a.v.Load().(uint64)
}

func (a *AtomicUint64) Set(v uint64) {
	a.v.Store(v)
}

type AtomicUint32 struct {
	v atomic.Value
}

func NewAtomicUint32(v uint32) *AtomicUint32 {
	var av atomic.Value
	av.Store(v)
	return &AtomicUint32{v: av}
}

func (a *AtomicUint32) Get() uint32 {
	return a.v.Load().(uint32)
}

func (a *AtomicUint32) Set(v uint32) {
	a.v.Store(v)
}

type AtomicUint16 struct {
	v atomic.Value
}

func NewAtomicUint16(v uint16) *AtomicUint16 {
	var av atomic.Value
	av.Store(v)
	return &AtomicUint16{v: av}
}

func (a *AtomicUint16) Get() uint16 {
	return a.v.Load().(uint16)
}

func (a *AtomicUint16) Set(v uint16) {
	a.v.Store(v)
}

type AtomicUint8 struct {
	v atomic.Value
}

func NewAtomicUint8(v uint8) *AtomicUint8 {
	var av atomic.Value
	av.Store(v)
	return &AtomicUint8{v: av}
}

func (a *AtomicUint8) Get() uint8 {
	return a.v.Load().(uint8)
}

func (a *AtomicUint8) Set(v uint8) {
	a.v.Store(v)
}

type AtomicFloat64 struct {
	v atomic.Value
}

func NewAtomicFloat64(v float64) *AtomicFloat64 {
	var av atomic.Value
	av.Store(v)
	return &AtomicFloat64{v: av}
}

func (a *AtomicFloat64) Get() float64 {
	return a.v.Load().(float64)
}

func (a *AtomicFloat64) Set(v float64) {
	a.v.Store(v)
}

type AtomicFloat32 struct {
	v atomic.Value
}

func NewAtomicFloat32(v float32) *AtomicFloat32 {
	var av atomic.Value
	av.Store(v)
	return &AtomicFloat32{v: av}
}

func (a *AtomicFloat32) Get() float32 {
	return a.v.Load().(float32)
}

func (a *AtomicFloat32) Set(v float32) {
	a.v.Store(v)
}

type AtomicString struct {
	v atomic.Value
}

func NewAtomicString(v string) *AtomicString {
	var av atomic.Value
	av.Store(v)
	return &AtomicString{v: av}
}

func (a *AtomicString) Get() string {
	return a.v.Load().(string)
}

func (a *AtomicString) Set(v string) {
	a.v.Store(v)
}

type AtomicDuration struct {
	v atomic.Value
}

func NewAtomicDuration(v time.Duration) *AtomicDuration {
	var av atomic.Value
	av.Store(v)
	return &AtomicDuration{v: av}
}

func (a *AtomicDuration) Get() time.Duration {
	return a.v.Load().(time.Duration)
}

func (a *AtomicDuration) Set(v time.Duration) {
	a.v.Store(v)
}

type AtomicTime struct {
	v atomic.Value
}

func NewAtomicTime(v time.Time) *AtomicTime {
	var av atomic.Value
	av.Store(v)
	return &AtomicTime{v: av}
}

func (a *AtomicTime) Get() time.Time {
	return a.v.Load().(time.Time)
}

func (a *AtomicTime) Set(v time.Time) {
	a.v.Store(v)
}

type AtomicIP struct {
	v atomic.Value
}

func NewAtomicIP(v net.IP) *AtomicIP {
	var av atomic.Value
	av.Store(v)
	return &AtomicIP{v: av}
}

func (a *AtomicIP) Get() net.IP {
	return a.v.Load().(net.IP)
}

func (a *AtomicIP) Set(v net.IP) {
	a.v.Store(v)
}

func (c *Config) BoolVar(key string, av *AtomicBool, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v bool
	if c.storage != nil {
		v = c.GetBool(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetBoolE(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) IntVar(key string, av *AtomicInt, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v int
	if c.storage != nil {
		v = c.GetInt(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetIntE(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) UintVar(key string, av *AtomicUint, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v uint
	if c.storage != nil {
		v = c.GetUint(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetUintE(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) Int64Var(key string, av *AtomicInt64, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v int64
	if c.storage != nil {
		v = c.GetInt64(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetInt64E(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) Int32Var(key string, av *AtomicInt32, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v int32
	if c.storage != nil {
		v = c.GetInt32(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetInt32E(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) Int16Var(key string, av *AtomicInt16, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v int16
	if c.storage != nil {
		v = c.GetInt16(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetInt16E(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) Int8Var(key string, av *AtomicInt8, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v int8
	if c.storage != nil {
		v = c.GetInt8(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetInt8E(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) Uint64Var(key string, av *AtomicUint64, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v uint64
	if c.storage != nil {
		v = c.GetUint64(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetUint64E(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) Uint32Var(key string, av *AtomicUint32, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v uint32
	if c.storage != nil {
		v = c.GetUint32(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetUint32E(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) Uint16Var(key string, av *AtomicUint16, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v uint16
	if c.storage != nil {
		v = c.GetUint16(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetUint16E(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) Uint8Var(key string, av *AtomicUint8, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v uint8
	if c.storage != nil {
		v = c.GetUint8(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetUint8E(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) Float64Var(key string, av *AtomicFloat64, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v float64
	if c.storage != nil {
		v = c.GetFloat64(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetFloat64E(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) Float32Var(key string, av *AtomicFloat32, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v float32
	if c.storage != nil {
		v = c.GetFloat32(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetFloat32E(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) StringVar(key string, av *AtomicString, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v string
	if c.storage != nil {
		v = c.GetString(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetStringE(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) DurationVar(key string, av *AtomicDuration, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v time.Duration
	if c.storage != nil {
		v = c.GetDuration(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetDurationE(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) TimeVar(key string, av *AtomicTime, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v time.Time
	if c.storage != nil {
		v = c.GetTime(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetTimeE(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) IPVar(key string, av *AtomicIP, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v net.IP
	if c.storage != nil {
		v = c.GetIP(key)
	}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		var err error
		v, err = c.GetIPE(key)
		if err != nil {
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Set(v)
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}

func (c *Config) Bool(key string, opts ...BindOption) *AtomicBool {
	var v AtomicBool
	c.BoolVar(key, &v, opts...)
	return &v
}

func (c *Config) Int(key string, opts ...BindOption) *AtomicInt {
	var v AtomicInt
	c.IntVar(key, &v, opts...)
	return &v
}

func (c *Config) Uint(key string, opts ...BindOption) *AtomicUint {
	var v AtomicUint
	c.UintVar(key, &v, opts...)
	return &v
}

func (c *Config) Int64(key string, opts ...BindOption) *AtomicInt64 {
	var v AtomicInt64
	c.Int64Var(key, &v, opts...)
	return &v
}

func (c *Config) Int32(key string, opts ...BindOption) *AtomicInt32 {
	var v AtomicInt32
	c.Int32Var(key, &v, opts...)
	return &v
}

func (c *Config) Int16(key string, opts ...BindOption) *AtomicInt16 {
	var v AtomicInt16
	c.Int16Var(key, &v, opts...)
	return &v
}

func (c *Config) Int8(key string, opts ...BindOption) *AtomicInt8 {
	var v AtomicInt8
	c.Int8Var(key, &v, opts...)
	return &v
}

func (c *Config) Uint64(key string, opts ...BindOption) *AtomicUint64 {
	var v AtomicUint64
	c.Uint64Var(key, &v, opts...)
	return &v
}

func (c *Config) Uint32(key string, opts ...BindOption) *AtomicUint32 {
	var v AtomicUint32
	c.Uint32Var(key, &v, opts...)
	return &v
}

func (c *Config) Uint16(key string, opts ...BindOption) *AtomicUint16 {
	var v AtomicUint16
	c.Uint16Var(key, &v, opts...)
	return &v
}

func (c *Config) Uint8(key string, opts ...BindOption) *AtomicUint8 {
	var v AtomicUint8
	c.Uint8Var(key, &v, opts...)
	return &v
}

func (c *Config) Float64(key string, opts ...BindOption) *AtomicFloat64 {
	var v AtomicFloat64
	c.Float64Var(key, &v, opts...)
	return &v
}

func (c *Config) Float32(key string, opts ...BindOption) *AtomicFloat32 {
	var v AtomicFloat32
	c.Float32Var(key, &v, opts...)
	return &v
}

func (c *Config) String(key string, opts ...BindOption) *AtomicString {
	var v AtomicString
	c.StringVar(key, &v, opts...)
	return &v
}

func (c *Config) Duration(key string, opts ...BindOption) *AtomicDuration {
	var v AtomicDuration
	c.DurationVar(key, &v, opts...)
	return &v
}

func (c *Config) Time(key string, opts ...BindOption) *AtomicTime {
	var v AtomicTime
	c.TimeVar(key, &v, opts...)
	return &v
}

func (c *Config) IP(key string, opts ...BindOption) *AtomicIP {
	var v AtomicIP
	c.IPVar(key, &v, opts...)
	return &v
}
