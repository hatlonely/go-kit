//go:generate python3 autogen.py
package refex

import (
	"fmt"
	"net"
	"reflect"
	"time"

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
	switch v.(type) {
	case string:
		ip := net.ParseIP(v.(string))
		if ip == nil {
			return nil, fmt.Errorf("pase ip failed. val [%v]", v)
		}
		return ip, nil
	case net.IP:
		return v.(net.IP), nil
	default:
		return nil, fmt.Errorf("convert type [%v] to ip failed", reflect.TypeOf(v))
	}
}

func SetInterface(dst interface{}, src interface{}) error {
	switch dst.(type) {
	case *bool:
		v, err := ToBoolE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *int:
		v, err := ToIntE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *uint:
		v, err := ToUintE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *int64:
		v, err := ToInt64E(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *int32:
		v, err := ToInt32E(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *int16:
		v, err := ToInt16E(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *int8:
		v, err := ToInt8E(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *uint64:
		v, err := ToUint64E(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *uint32:
		v, err := ToUint32E(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *uint16:
		v, err := ToUint16E(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *uint8:
		v, err := ToUint8E(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *float64:
		v, err := ToFloat64E(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *float32:
		v, err := ToFloat32E(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *string:
		v, err := ToStringE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *time.Duration:
		v, err := ToDurationE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *time.Time:
		v, err := ToTimeE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *net.IP:
		v, err := ToIPE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))

	default:
		return fmt.Errorf("unsupport dst type [%v]", reflect.TypeOf(dst))
	}

	return nil
}
