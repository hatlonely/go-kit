// this file is auto generate by autogen.py. do not edit!
package config

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/hatlonely/go-kit/refx"
)

var gcfg = &Config{
	itemHandlers: map[string][]OnChangeHandler{},
}

func Init(filename string, opts ...refx.Option) error {
	cfg, err := NewConfigWithBaseFile(filename, opts...)
	if err != nil {
		return err
	}
	cfg.itemHandlers = gcfg.itemHandlers
	gcfg = cfg
	return nil
}

func InitWithSimpleFile(filename string, opts ...SimpleFileOption) error {
	cfg, err := NewConfigWithSimpleFile(filename, opts...)
	if err != nil {
		return err
	}
	cfg.itemHandlers = gcfg.itemHandlers
	gcfg = cfg
	return nil
}

func GetComponent() (Provider, *Storage, Decoder, Cipher) {
	return gcfg.GetComponent()
}

func Get(key string) (interface{}, bool) {
	return gcfg.Get(key)
}

func GetE(key string) (interface{}, error) {
	return gcfg.GetE(key)
}

func UnsafeSet(key string, val interface{}) error {
	return gcfg.UnsafeSet(key, val)
}

func Unmarshal(v interface{}, opts ...refx.Option) error {
	return gcfg.Unmarshal(v, opts...)
}

func Sub(key string) *Config {
	return gcfg.Sub(key)
}

func SubArr(key string) ([]*Config, error) {
	return gcfg.SubArr(key)
}

func SubMap(key string) (map[string]*Config, error) {
	return gcfg.SubMap(key)
}

func Stop() {
	gcfg.Stop()
}

func Watch() error {
	return gcfg.Watch()
}

func AddOnChangeHandler(handler OnChangeHandler) {
	gcfg.AddOnChangeHandler(handler)
}

func AddOnItemChangeHandler(key string, handler OnChangeHandler) {
	gcfg.AddOnItemChangeHandler(key, handler)
}

func TransformWithOptions(options *Options, transformOptions *TransformOptions) (*Config, error) {
	return gcfg.TransformWithOptions(options, transformOptions)
}

func Transform(options *Options, opts ...TransformOption) (*Config, error) {
	return gcfg.Transform(options, opts...)
}

func Bytes() ([]byte, error) {
	return gcfg.Bytes()
}

func Save() error {
	return gcfg.Save()
}

func Diff(o *Config) {
	gcfg.Diff(o)
}

func ToString() string {
	return gcfg.ToString()
}

func ToJsonString() string {
	return gcfg.ToJsonString()
}

func SetLogger(log Logger) {
	gcfg.SetLogger(log)
}

func GetBool(key string) bool {
	return gcfg.GetBool(key)
}

func GetInt(key string) int {
	return gcfg.GetInt(key)
}

func GetUint(key string) uint {
	return gcfg.GetUint(key)
}

func GetInt64(key string) int64 {
	return gcfg.GetInt64(key)
}

func GetInt32(key string) int32 {
	return gcfg.GetInt32(key)
}

func GetInt16(key string) int16 {
	return gcfg.GetInt16(key)
}

func GetInt8(key string) int8 {
	return gcfg.GetInt8(key)
}

func GetUint64(key string) uint64 {
	return gcfg.GetUint64(key)
}

func GetUint32(key string) uint32 {
	return gcfg.GetUint32(key)
}

func GetUint16(key string) uint16 {
	return gcfg.GetUint16(key)
}

func GetUint8(key string) uint8 {
	return gcfg.GetUint8(key)
}

func GetFloat64(key string) float64 {
	return gcfg.GetFloat64(key)
}

func GetFloat32(key string) float32 {
	return gcfg.GetFloat32(key)
}

func GetString(key string) string {
	return gcfg.GetString(key)
}

func GetDuration(key string) time.Duration {
	return gcfg.GetDuration(key)
}

func GetTime(key string) time.Time {
	return gcfg.GetTime(key)
}

func GetIP(key string) net.IP {
	return gcfg.GetIP(key)
}

func GetBoolE(key string) (bool, error) {
	return gcfg.GetBoolE(key)
}

func GetIntE(key string) (int, error) {
	return gcfg.GetIntE(key)
}

func GetUintE(key string) (uint, error) {
	return gcfg.GetUintE(key)
}

func GetInt64E(key string) (int64, error) {
	return gcfg.GetInt64E(key)
}

func GetInt32E(key string) (int32, error) {
	return gcfg.GetInt32E(key)
}

func GetInt16E(key string) (int16, error) {
	return gcfg.GetInt16E(key)
}

func GetInt8E(key string) (int8, error) {
	return gcfg.GetInt8E(key)
}

func GetUint64E(key string) (uint64, error) {
	return gcfg.GetUint64E(key)
}

func GetUint32E(key string) (uint32, error) {
	return gcfg.GetUint32E(key)
}

func GetUint16E(key string) (uint16, error) {
	return gcfg.GetUint16E(key)
}

func GetUint8E(key string) (uint8, error) {
	return gcfg.GetUint8E(key)
}

func GetFloat64E(key string) (float64, error) {
	return gcfg.GetFloat64E(key)
}

func GetFloat32E(key string) (float32, error) {
	return gcfg.GetFloat32E(key)
}

func GetStringE(key string) (string, error) {
	return gcfg.GetStringE(key)
}

func GetDurationE(key string) (time.Duration, error) {
	return gcfg.GetDurationE(key)
}

func GetTimeE(key string) (time.Time, error) {
	return gcfg.GetTimeE(key)
}

func GetIPE(key string) (net.IP, error) {
	return gcfg.GetIPE(key)
}

func GetBoolP(key string) bool {
	return gcfg.GetBoolP(key)
}

func GetIntP(key string) int {
	return gcfg.GetIntP(key)
}

func GetUintP(key string) uint {
	return gcfg.GetUintP(key)
}

func GetInt64P(key string) int64 {
	return gcfg.GetInt64P(key)
}

func GetInt32P(key string) int32 {
	return gcfg.GetInt32P(key)
}

func GetInt16P(key string) int16 {
	return gcfg.GetInt16P(key)
}

func GetInt8P(key string) int8 {
	return gcfg.GetInt8P(key)
}

func GetUint64P(key string) uint64 {
	return gcfg.GetUint64P(key)
}

func GetUint32P(key string) uint32 {
	return gcfg.GetUint32P(key)
}

func GetUint16P(key string) uint16 {
	return gcfg.GetUint16P(key)
}

func GetUint8P(key string) uint8 {
	return gcfg.GetUint8P(key)
}

func GetFloat64P(key string) float64 {
	return gcfg.GetFloat64P(key)
}

func GetFloat32P(key string) float32 {
	return gcfg.GetFloat32P(key)
}

func GetStringP(key string) string {
	return gcfg.GetStringP(key)
}

func GetDurationP(key string) time.Duration {
	return gcfg.GetDurationP(key)
}

func GetTimeP(key string) time.Time {
	return gcfg.GetTimeP(key)
}

func GetIPP(key string) net.IP {
	return gcfg.GetIPP(key)
}

func GetBoolD(key string, dftVal bool) bool {
	return gcfg.GetBoolD(key, dftVal)
}

func GetIntD(key string, dftVal int) int {
	return gcfg.GetIntD(key, dftVal)
}

func GetUintD(key string, dftVal uint) uint {
	return gcfg.GetUintD(key, dftVal)
}

func GetInt64D(key string, dftVal int64) int64 {
	return gcfg.GetInt64D(key, dftVal)
}

func GetInt32D(key string, dftVal int32) int32 {
	return gcfg.GetInt32D(key, dftVal)
}

func GetInt16D(key string, dftVal int16) int16 {
	return gcfg.GetInt16D(key, dftVal)
}

func GetInt8D(key string, dftVal int8) int8 {
	return gcfg.GetInt8D(key, dftVal)
}

func GetUint64D(key string, dftVal uint64) uint64 {
	return gcfg.GetUint64D(key, dftVal)
}

func GetUint32D(key string, dftVal uint32) uint32 {
	return gcfg.GetUint32D(key, dftVal)
}

func GetUint16D(key string, dftVal uint16) uint16 {
	return gcfg.GetUint16D(key, dftVal)
}

func GetUint8D(key string, dftVal uint8) uint8 {
	return gcfg.GetUint8D(key, dftVal)
}

func GetFloat64D(key string, dftVal float64) float64 {
	return gcfg.GetFloat64D(key, dftVal)
}

func GetFloat32D(key string, dftVal float32) float32 {
	return gcfg.GetFloat32D(key, dftVal)
}

func GetStringD(key string, dftVal string) string {
	return gcfg.GetStringD(key, dftVal)
}

func GetDurationD(key string, dftVal time.Duration) time.Duration {
	return gcfg.GetDurationD(key, dftVal)
}

func GetTimeD(key string, dftVal time.Time) time.Time {
	return gcfg.GetTimeD(key, dftVal)
}

func GetIPD(key string, dftVal net.IP) net.IP {
	return gcfg.GetIPD(key, dftVal)
}

func Bind(key string, v interface{}, opts ...BindOption) *atomic.Value {
	return gcfg.Bind(key, v, opts...)
}

func BindVar(key string, v interface{}, av *atomic.Value, opts ...BindOption) {
	gcfg.BindVar(key, v, av, opts...)
}

func BoolVar(key string, av *AtomicBool, opts ...BindOption) {
	gcfg.BoolVar(key, av, opts...)
}

func IntVar(key string, av *AtomicInt, opts ...BindOption) {
	gcfg.IntVar(key, av, opts...)
}

func UintVar(key string, av *AtomicUint, opts ...BindOption) {
	gcfg.UintVar(key, av, opts...)
}

func Int64Var(key string, av *AtomicInt64, opts ...BindOption) {
	gcfg.Int64Var(key, av, opts...)
}

func Int32Var(key string, av *AtomicInt32, opts ...BindOption) {
	gcfg.Int32Var(key, av, opts...)
}

func Int16Var(key string, av *AtomicInt16, opts ...BindOption) {
	gcfg.Int16Var(key, av, opts...)
}

func Int8Var(key string, av *AtomicInt8, opts ...BindOption) {
	gcfg.Int8Var(key, av, opts...)
}

func Uint64Var(key string, av *AtomicUint64, opts ...BindOption) {
	gcfg.Uint64Var(key, av, opts...)
}

func Uint32Var(key string, av *AtomicUint32, opts ...BindOption) {
	gcfg.Uint32Var(key, av, opts...)
}

func Uint16Var(key string, av *AtomicUint16, opts ...BindOption) {
	gcfg.Uint16Var(key, av, opts...)
}

func Uint8Var(key string, av *AtomicUint8, opts ...BindOption) {
	gcfg.Uint8Var(key, av, opts...)
}

func Float64Var(key string, av *AtomicFloat64, opts ...BindOption) {
	gcfg.Float64Var(key, av, opts...)
}

func Float32Var(key string, av *AtomicFloat32, opts ...BindOption) {
	gcfg.Float32Var(key, av, opts...)
}

func StringVar(key string, av *AtomicString, opts ...BindOption) {
	gcfg.StringVar(key, av, opts...)
}

func DurationVar(key string, av *AtomicDuration, opts ...BindOption) {
	gcfg.DurationVar(key, av, opts...)
}

func TimeVar(key string, av *AtomicTime, opts ...BindOption) {
	gcfg.TimeVar(key, av, opts...)
}

func IPVar(key string, av *AtomicIP, opts ...BindOption) {
	gcfg.IPVar(key, av, opts...)
}

func Bool(key string, opts ...BindOption) *AtomicBool {
	return gcfg.Bool(key, opts...)
}

func Int(key string, opts ...BindOption) *AtomicInt {
	return gcfg.Int(key, opts...)
}

func Uint(key string, opts ...BindOption) *AtomicUint {
	return gcfg.Uint(key, opts...)
}

func Int64(key string, opts ...BindOption) *AtomicInt64 {
	return gcfg.Int64(key, opts...)
}

func Int32(key string, opts ...BindOption) *AtomicInt32 {
	return gcfg.Int32(key, opts...)
}

func Int16(key string, opts ...BindOption) *AtomicInt16 {
	return gcfg.Int16(key, opts...)
}

func Int8(key string, opts ...BindOption) *AtomicInt8 {
	return gcfg.Int8(key, opts...)
}

func Uint64(key string, opts ...BindOption) *AtomicUint64 {
	return gcfg.Uint64(key, opts...)
}

func Uint32(key string, opts ...BindOption) *AtomicUint32 {
	return gcfg.Uint32(key, opts...)
}

func Uint16(key string, opts ...BindOption) *AtomicUint16 {
	return gcfg.Uint16(key, opts...)
}

func Uint8(key string, opts ...BindOption) *AtomicUint8 {
	return gcfg.Uint8(key, opts...)
}

func Float64(key string, opts ...BindOption) *AtomicFloat64 {
	return gcfg.Float64(key, opts...)
}

func Float32(key string, opts ...BindOption) *AtomicFloat32 {
	return gcfg.Float32(key, opts...)
}

func String(key string, opts ...BindOption) *AtomicString {
	return gcfg.String(key, opts...)
}

func Duration(key string, opts ...BindOption) *AtomicDuration {
	return gcfg.Duration(key, opts...)
}

func Time(key string, opts ...BindOption) *AtomicTime {
	return gcfg.Time(key, opts...)
}

func IP(key string, opts ...BindOption) *AtomicIP {
	return gcfg.IP(key, opts...)
}
