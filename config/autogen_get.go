// this file is auto generate by autogen.py. do not edit!
package config

import (
	"net"
	"time"
	
	"github.com/hatlonely/go-kit/cast"
)

func (c *Config) GetBool(key string) bool {
	v, _ := c.GetBoolE(key)
	return v
}

func (c *Config) GetInt(key string) int {
	v, _ := c.GetIntE(key)
	return v
}

func (c *Config) GetUint(key string) uint {
	v, _ := c.GetUintE(key)
	return v
}

func (c *Config) GetInt64(key string) int64 {
	v, _ := c.GetInt64E(key)
	return v
}

func (c *Config) GetInt32(key string) int32 {
	v, _ := c.GetInt32E(key)
	return v
}

func (c *Config) GetInt16(key string) int16 {
	v, _ := c.GetInt16E(key)
	return v
}

func (c *Config) GetInt8(key string) int8 {
	v, _ := c.GetInt8E(key)
	return v
}

func (c *Config) GetUint64(key string) uint64 {
	v, _ := c.GetUint64E(key)
	return v
}

func (c *Config) GetUint32(key string) uint32 {
	v, _ := c.GetUint32E(key)
	return v
}

func (c *Config) GetUint16(key string) uint16 {
	v, _ := c.GetUint16E(key)
	return v
}

func (c *Config) GetUint8(key string) uint8 {
	v, _ := c.GetUint8E(key)
	return v
}

func (c *Config) GetFloat64(key string) float64 {
	v, _ := c.GetFloat64E(key)
	return v
}

func (c *Config) GetFloat32(key string) float32 {
	v, _ := c.GetFloat32E(key)
	return v
}

func (c *Config) GetString(key string) string {
	v, _ := c.GetStringE(key)
	return v
}

func (c *Config) GetDuration(key string) time.Duration {
	v, _ := c.GetDurationE(key)
	return v
}

func (c *Config) GetTime(key string) time.Time {
	v, _ := c.GetTimeE(key)
	return v
}

func (c *Config) GetIP(key string) net.IP {
	v, _ := c.GetIPE(key)
	return v
}

func (c *Config) GetBoolE(key string) (bool, error) {
	v, err := c.Get(key)
	if err != nil {
		var res bool
		return res, err
	}
	return cast.ToBoolE(v)
}

func (c *Config) GetIntE(key string) (int, error) {
	v, err := c.Get(key)
	if err != nil {
		var res int
		return res, err
	}
	return cast.ToIntE(v)
}

func (c *Config) GetUintE(key string) (uint, error) {
	v, err := c.Get(key)
	if err != nil {
		var res uint
		return res, err
	}
	return cast.ToUintE(v)
}

func (c *Config) GetInt64E(key string) (int64, error) {
	v, err := c.Get(key)
	if err != nil {
		var res int64
		return res, err
	}
	return cast.ToInt64E(v)
}

func (c *Config) GetInt32E(key string) (int32, error) {
	v, err := c.Get(key)
	if err != nil {
		var res int32
		return res, err
	}
	return cast.ToInt32E(v)
}

func (c *Config) GetInt16E(key string) (int16, error) {
	v, err := c.Get(key)
	if err != nil {
		var res int16
		return res, err
	}
	return cast.ToInt16E(v)
}

func (c *Config) GetInt8E(key string) (int8, error) {
	v, err := c.Get(key)
	if err != nil {
		var res int8
		return res, err
	}
	return cast.ToInt8E(v)
}

func (c *Config) GetUint64E(key string) (uint64, error) {
	v, err := c.Get(key)
	if err != nil {
		var res uint64
		return res, err
	}
	return cast.ToUint64E(v)
}

func (c *Config) GetUint32E(key string) (uint32, error) {
	v, err := c.Get(key)
	if err != nil {
		var res uint32
		return res, err
	}
	return cast.ToUint32E(v)
}

func (c *Config) GetUint16E(key string) (uint16, error) {
	v, err := c.Get(key)
	if err != nil {
		var res uint16
		return res, err
	}
	return cast.ToUint16E(v)
}

func (c *Config) GetUint8E(key string) (uint8, error) {
	v, err := c.Get(key)
	if err != nil {
		var res uint8
		return res, err
	}
	return cast.ToUint8E(v)
}

func (c *Config) GetFloat64E(key string) (float64, error) {
	v, err := c.Get(key)
	if err != nil {
		var res float64
		return res, err
	}
	return cast.ToFloat64E(v)
}

func (c *Config) GetFloat32E(key string) (float32, error) {
	v, err := c.Get(key)
	if err != nil {
		var res float32
		return res, err
	}
	return cast.ToFloat32E(v)
}

func (c *Config) GetStringE(key string) (string, error) {
	v, err := c.Get(key)
	if err != nil {
		var res string
		return res, err
	}
	return cast.ToStringE(v)
}

func (c *Config) GetDurationE(key string) (time.Duration, error) {
	v, err := c.Get(key)
	if err != nil {
		var res time.Duration
		return res, err
	}
	return cast.ToDurationE(v)
}

func (c *Config) GetTimeE(key string) (time.Time, error) {
	v, err := c.Get(key)
	if err != nil {
		var res time.Time
		return res, err
	}
	return cast.ToTimeE(v)
}

func (c *Config) GetIPE(key string) (net.IP, error) {
	v, err := c.Get(key)
	if err != nil {
		var res net.IP
		return res, err
	}
	return cast.ToIPE(v)
}

func (c *Config) GetBoolP(key string) bool {
	v, err := c.GetBoolE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetIntP(key string) int {
	v, err := c.GetIntE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetUintP(key string) uint {
	v, err := c.GetUintE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetInt64P(key string) int64 {
	v, err := c.GetInt64E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetInt32P(key string) int32 {
	v, err := c.GetInt32E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetInt16P(key string) int16 {
	v, err := c.GetInt16E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetInt8P(key string) int8 {
	v, err := c.GetInt8E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetUint64P(key string) uint64 {
	v, err := c.GetUint64E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetUint32P(key string) uint32 {
	v, err := c.GetUint32E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetUint16P(key string) uint16 {
	v, err := c.GetUint16E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetUint8P(key string) uint8 {
	v, err := c.GetUint8E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetFloat64P(key string) float64 {
	v, err := c.GetFloat64E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetFloat32P(key string) float32 {
	v, err := c.GetFloat32E(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetStringP(key string) string {
	v, err := c.GetStringE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetDurationP(key string) time.Duration {
	v, err := c.GetDurationE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetTimeP(key string) time.Time {
	v, err := c.GetTimeE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetIPP(key string) net.IP {
	v, err := c.GetIPE(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Config) GetBoolD(key string, dftVal bool) bool {
	v, err := c.GetBoolE(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetIntD(key string, dftVal int) int {
	v, err := c.GetIntE(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetUintD(key string, dftVal uint) uint {
	v, err := c.GetUintE(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetInt64D(key string, dftVal int64) int64 {
	v, err := c.GetInt64E(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetInt32D(key string, dftVal int32) int32 {
	v, err := c.GetInt32E(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetInt16D(key string, dftVal int16) int16 {
	v, err := c.GetInt16E(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetInt8D(key string, dftVal int8) int8 {
	v, err := c.GetInt8E(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetUint64D(key string, dftVal uint64) uint64 {
	v, err := c.GetUint64E(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetUint32D(key string, dftVal uint32) uint32 {
	v, err := c.GetUint32E(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetUint16D(key string, dftVal uint16) uint16 {
	v, err := c.GetUint16E(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetUint8D(key string, dftVal uint8) uint8 {
	v, err := c.GetUint8E(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetFloat64D(key string, dftVal float64) float64 {
	v, err := c.GetFloat64E(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetFloat32D(key string, dftVal float32) float32 {
	v, err := c.GetFloat32E(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetStringD(key string, dftVal string) string {
	v, err := c.GetStringE(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetDurationD(key string, dftVal time.Duration) time.Duration {
	v, err := c.GetDurationE(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetTimeD(key string, dftVal time.Time) time.Time {
	v, err := c.GetTimeE(key)
	if err != nil {
		return dftVal
	}
	return v
}

func (c *Config) GetIPD(key string, dftVal net.IP) net.IP {
	v, err := c.GetIPE(key)
	if err != nil {
		return dftVal
	}
	return v
}
