// this file is generate by autogen.py, do not edit
package config

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

var gconf = &Config{
	itemHandlers: map[string][]OnChangeHandler{},
	log:          logrus.New(),
}

func Init(filename string) error {
	conf, err := NewConfigWithBaseFile(filename)
	if err != nil {
		return err
	}
	gconf.storage = conf.storage
	gconf.decoder = conf.decoder
	gconf.provider = conf.provider
	return nil
}

func GetComponent() (Provider, *Storage, Decoder, Cipher) {
	return gconf.GetComponent()
}

func Get(key string) (interface{}, error) {
	return gconf.Get(key)
}

func UnsafeSet(key string, val interface{}) error {
	return gconf.UnsafeSet(key, val)
}

func Unmarshal(v interface{}) error {
	return gconf.Unmarshal(v)
}

func Sub(key string) *Config {
	return gconf.Sub(key)
}

func SubArr(key string) ([]*Config, error) {
	return gconf.SubArr(key)
}

func SubMap(key string) (map[string]*Config, error) {
	return gconf.SubMap(key)
}

func Stop() {
	gconf.Stop()
}

func Watch() error {
	return gconf.Watch()
}

func AddOnChangeHandler(handler OnChangeHandler) {
	gconf.AddOnChangeHandler(handler)
}

func AddOnItemChangeHandler(key string, handler OnChangeHandler) {
	gconf.AddOnItemChangeHandler(key, handler)
}

func GetBool(key string) bool {
	return gconf.GetBool(key)
}

func GetInt(key string) int {
	return gconf.GetInt(key)
}

func GetUint(key string) uint {
	return gconf.GetUint(key)
}

func GetInt64(key string) int64 {
	return gconf.GetInt64(key)
}

func GetInt32(key string) int32 {
	return gconf.GetInt32(key)
}

func GetInt16(key string) int16 {
	return gconf.GetInt16(key)
}

func GetInt8(key string) int8 {
	return gconf.GetInt8(key)
}

func GetUint64(key string) uint64 {
	return gconf.GetUint64(key)
}

func GetUint32(key string) uint32 {
	return gconf.GetUint32(key)
}

func GetUint16(key string) uint16 {
	return gconf.GetUint16(key)
}

func GetUint8(key string) uint8 {
	return gconf.GetUint8(key)
}

func GetFloat64(key string) float64 {
	return gconf.GetFloat64(key)
}

func GetFloat32(key string) float32 {
	return gconf.GetFloat32(key)
}

func GetString(key string) string {
	return gconf.GetString(key)
}

func GetDuration(key string) time.Duration {
	return gconf.GetDuration(key)
}

func GetTime(key string) time.Time {
	return gconf.GetTime(key)
}

func GetIP(key string) net.IP {
	return gconf.GetIP(key)
}

func GetBoolE(key string) (bool, error) {
	return gconf.GetBoolE(key)
}

func GetIntE(key string) (int, error) {
	return gconf.GetIntE(key)
}

func GetUintE(key string) (uint, error) {
	return gconf.GetUintE(key)
}

func GetInt64E(key string) (int64, error) {
	return gconf.GetInt64E(key)
}

func GetInt32E(key string) (int32, error) {
	return gconf.GetInt32E(key)
}

func GetInt16E(key string) (int16, error) {
	return gconf.GetInt16E(key)
}

func GetInt8E(key string) (int8, error) {
	return gconf.GetInt8E(key)
}

func GetUint64E(key string) (uint64, error) {
	return gconf.GetUint64E(key)
}

func GetUint32E(key string) (uint32, error) {
	return gconf.GetUint32E(key)
}

func GetUint16E(key string) (uint16, error) {
	return gconf.GetUint16E(key)
}

func GetUint8E(key string) (uint8, error) {
	return gconf.GetUint8E(key)
}

func GetFloat64E(key string) (float64, error) {
	return gconf.GetFloat64E(key)
}

func GetFloat32E(key string) (float32, error) {
	return gconf.GetFloat32E(key)
}

func GetStringE(key string) (string, error) {
	return gconf.GetStringE(key)
}

func GetDurationE(key string) (time.Duration, error) {
	return gconf.GetDurationE(key)
}

func GetTimeE(key string) (time.Time, error) {
	return gconf.GetTimeE(key)
}

func GetIPE(key string) (net.IP, error) {
	return gconf.GetIPE(key)
}

func GetBoolP(key string) bool {
	return gconf.GetBoolP(key)
}

func GetIntP(key string) int {
	return gconf.GetIntP(key)
}

func GetUintP(key string) uint {
	return gconf.GetUintP(key)
}

func GetInt64P(key string) int64 {
	return gconf.GetInt64P(key)
}

func GetInt32P(key string) int32 {
	return gconf.GetInt32P(key)
}

func GetInt16P(key string) int16 {
	return gconf.GetInt16P(key)
}

func GetInt8P(key string) int8 {
	return gconf.GetInt8P(key)
}

func GetUint64P(key string) uint64 {
	return gconf.GetUint64P(key)
}

func GetUint32P(key string) uint32 {
	return gconf.GetUint32P(key)
}

func GetUint16P(key string) uint16 {
	return gconf.GetUint16P(key)
}

func GetUint8P(key string) uint8 {
	return gconf.GetUint8P(key)
}

func GetFloat64P(key string) float64 {
	return gconf.GetFloat64P(key)
}

func GetFloat32P(key string) float32 {
	return gconf.GetFloat32P(key)
}

func GetStringP(key string) string {
	return gconf.GetStringP(key)
}

func GetDurationP(key string) time.Duration {
	return gconf.GetDurationP(key)
}

func GetTimeP(key string) time.Time {
	return gconf.GetTimeP(key)
}

func GetIPP(key string) net.IP {
	return gconf.GetIPP(key)
}

func GetBoolD(key string, dftVal bool) bool {
	return gconf.GetBoolD(key, dftVal)
}

func GetIntD(key string, dftVal int) int {
	return gconf.GetIntD(key, dftVal)
}

func GetUintD(key string, dftVal uint) uint {
	return gconf.GetUintD(key, dftVal)
}

func GetInt64D(key string, dftVal int64) int64 {
	return gconf.GetInt64D(key, dftVal)
}

func GetInt32D(key string, dftVal int32) int32 {
	return gconf.GetInt32D(key, dftVal)
}

func GetInt16D(key string, dftVal int16) int16 {
	return gconf.GetInt16D(key, dftVal)
}

func GetInt8D(key string, dftVal int8) int8 {
	return gconf.GetInt8D(key, dftVal)
}

func GetUint64D(key string, dftVal uint64) uint64 {
	return gconf.GetUint64D(key, dftVal)
}

func GetUint32D(key string, dftVal uint32) uint32 {
	return gconf.GetUint32D(key, dftVal)
}

func GetUint16D(key string, dftVal uint16) uint16 {
	return gconf.GetUint16D(key, dftVal)
}

func GetUint8D(key string, dftVal uint8) uint8 {
	return gconf.GetUint8D(key, dftVal)
}

func GetFloat64D(key string, dftVal float64) float64 {
	return gconf.GetFloat64D(key, dftVal)
}

func GetFloat32D(key string, dftVal float32) float32 {
	return gconf.GetFloat32D(key, dftVal)
}

func GetStringD(key string, dftVal string) string {
	return gconf.GetStringD(key, dftVal)
}

func GetDurationD(key string, dftVal time.Duration) time.Duration {
	return gconf.GetDurationD(key, dftVal)
}

func GetTimeD(key string, dftVal time.Time) time.Time {
	return gconf.GetTimeD(key, dftVal)
}

func GetIPD(key string, dftVal net.IP) net.IP {
	return gconf.GetIPD(key, dftVal)
}

func Bind(key string, v interface{}, opts ...BindOption) *atomic.Value {
	return gconf.Bind(key, v, opts...)
}

func BindVar(key string, v interface{}, av *atomic.Value, opts ...BindOption) {
	gconf.BindVar(key, v, av, opts...)
}

func BoolVar(key string, av *AtomicBool, opts ...BindOption) {
	gconf.BoolVar(key, av, opts...)
}

func IntVar(key string, av *AtomicInt, opts ...BindOption) {
	gconf.IntVar(key, av, opts...)
}

func UintVar(key string, av *AtomicUint, opts ...BindOption) {
	gconf.UintVar(key, av, opts...)
}

func Int64Var(key string, av *AtomicInt64, opts ...BindOption) {
	gconf.Int64Var(key, av, opts...)
}

func Int32Var(key string, av *AtomicInt32, opts ...BindOption) {
	gconf.Int32Var(key, av, opts...)
}

func Int16Var(key string, av *AtomicInt16, opts ...BindOption) {
	gconf.Int16Var(key, av, opts...)
}

func Int8Var(key string, av *AtomicInt8, opts ...BindOption) {
	gconf.Int8Var(key, av, opts...)
}

func Uint64Var(key string, av *AtomicUint64, opts ...BindOption) {
	gconf.Uint64Var(key, av, opts...)
}

func Uint32Var(key string, av *AtomicUint32, opts ...BindOption) {
	gconf.Uint32Var(key, av, opts...)
}

func Uint16Var(key string, av *AtomicUint16, opts ...BindOption) {
	gconf.Uint16Var(key, av, opts...)
}

func Uint8Var(key string, av *AtomicUint8, opts ...BindOption) {
	gconf.Uint8Var(key, av, opts...)
}

func Float64Var(key string, av *AtomicFloat64, opts ...BindOption) {
	gconf.Float64Var(key, av, opts...)
}

func Float32Var(key string, av *AtomicFloat32, opts ...BindOption) {
	gconf.Float32Var(key, av, opts...)
}

func StringVar(key string, av *AtomicString, opts ...BindOption) {
	gconf.StringVar(key, av, opts...)
}

func DurationVar(key string, av *AtomicDuration, opts ...BindOption) {
	gconf.DurationVar(key, av, opts...)
}

func TimeVar(key string, av *AtomicTime, opts ...BindOption) {
	gconf.TimeVar(key, av, opts...)
}

func IPVar(key string, av *AtomicIP, opts ...BindOption) {
	gconf.IPVar(key, av, opts...)
}

func Bool(key string, opts ...BindOption) *AtomicBool {
	return gconf.Bool(key, opts...)
}

func Int(key string, opts ...BindOption) *AtomicInt {
	return gconf.Int(key, opts...)
}

func Uint(key string, opts ...BindOption) *AtomicUint {
	return gconf.Uint(key, opts...)
}

func Int64(key string, opts ...BindOption) *AtomicInt64 {
	return gconf.Int64(key, opts...)
}

func Int32(key string, opts ...BindOption) *AtomicInt32 {
	return gconf.Int32(key, opts...)
}

func Int16(key string, opts ...BindOption) *AtomicInt16 {
	return gconf.Int16(key, opts...)
}

func Int8(key string, opts ...BindOption) *AtomicInt8 {
	return gconf.Int8(key, opts...)
}

func Uint64(key string, opts ...BindOption) *AtomicUint64 {
	return gconf.Uint64(key, opts...)
}

func Uint32(key string, opts ...BindOption) *AtomicUint32 {
	return gconf.Uint32(key, opts...)
}

func Uint16(key string, opts ...BindOption) *AtomicUint16 {
	return gconf.Uint16(key, opts...)
}

func Uint8(key string, opts ...BindOption) *AtomicUint8 {
	return gconf.Uint8(key, opts...)
}

func Float64(key string, opts ...BindOption) *AtomicFloat64 {
	return gconf.Float64(key, opts...)
}

func Float32(key string, opts ...BindOption) *AtomicFloat32 {
	return gconf.Float32(key, opts...)
}

func String(key string, opts ...BindOption) *AtomicString {
	return gconf.String(key, opts...)
}

func Duration(key string, opts ...BindOption) *AtomicDuration {
	return gconf.Duration(key, opts...)
}

func Time(key string, opts ...BindOption) *AtomicTime {
	return gconf.Time(key, opts...)
}

func IP(key string, opts ...BindOption) *AtomicIP {
	return gconf.IP(key, opts...)
}
