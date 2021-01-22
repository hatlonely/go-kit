//go:generate python3 autogen.py
package cast

import (
	"encoding/json"
	"net"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func ToBoolE(v interface{}) (bool, error) {
	return cast.ToBoolE(v)
}

func ToIntE(v interface{}) (int, error) {
	return cast.ToIntE(v)
}

func ToUintE(v interface{}) (uint, error) {
	return cast.ToUintE(v)
}

func ToInt64E(v interface{}) (int64, error) {
	return cast.ToInt64E(v)
}

func ToInt32E(v interface{}) (int32, error) {
	return cast.ToInt32E(v)
}

func ToInt16E(v interface{}) (int16, error) {
	return cast.ToInt16E(v)
}

func ToInt8E(v interface{}) (int8, error) {
	return cast.ToInt8E(v)
}

func ToUint64E(v interface{}) (uint64, error) {
	return cast.ToUint64E(v)
}

func ToUint32E(v interface{}) (uint32, error) {
	return cast.ToUint32E(v)
}

func ToUint16E(v interface{}) (uint16, error) {
	return cast.ToUint16E(v)
}

func ToUint8E(v interface{}) (uint8, error) {
	return cast.ToUint8E(v)
}

func ToFloat64E(v interface{}) (float64, error) {
	return cast.ToFloat64E(v)
}

func ToFloat32E(v interface{}) (float32, error) {
	return cast.ToFloat32E(v)
}

func ToStringE(v interface{}) (string, error) {
	return cast.ToStringE(v)
}

func ToDurationE(v interface{}) (time.Duration, error) {
	return cast.ToDurationE(v)
}

func ToTimeE(v interface{}) (time.Time, error) {
	return cast.ToTimeE(v)
}

func ToIPE(v interface{}) (net.IP, error) {
	switch val := v.(type) {
	case string:
		ip := net.ParseIP(val)
		if ip == nil {
			return nil, errors.Errorf("pase ip failed. val [%v]", v)
		}
		return ip, nil
	case net.IP:
		return val, nil
	default:
		return nil, errors.Errorf("convert type [%v] to ip failed", reflect.TypeOf(v))
	}
}

func ToMapStringStringE(v interface{}) (map[string]string, error) {
	switch val := v.(type) {
	case string:
		m := map[string]string{}
		if err := json.Unmarshal([]byte(val), &m); err != nil {
			return nil, errors.Wrap(err, "json.Unmarshal failed")
		}
		return m, nil
	case map[string]string:
		return val, nil
	default:
		return nil, errors.Errorf("unsupported type [%v] to map[string]string", reflect.TypeOf(v))
	}
}
