// this file is auto generate by autogen.py. do not edit!
package cast

import (
	"net"
	"reflect"
	"regexp"
	"time"

	"github.com/pkg/errors"
)

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
	case **regexp.Regexp:
		v, err := ToRegexE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *map[string]string:
		v, err := ToMapStringStringE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]bool:
		v, err := ToBoolSliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]int:
		v, err := ToIntSliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]uint:
		v, err := ToUintSliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]int64:
		v, err := ToInt64SliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]int32:
		v, err := ToInt32SliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]int16:
		v, err := ToInt16SliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]int8:
		v, err := ToInt8SliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]uint64:
		v, err := ToUint64SliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]uint32:
		v, err := ToUint32SliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]uint16:
		v, err := ToUint16SliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]uint8:
		v, err := ToUint8SliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]float64:
		v, err := ToFloat64SliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]float32:
		v, err := ToFloat32SliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]string:
		v, err := ToStringSliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]time.Duration:
		v, err := ToDurationSliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]time.Time:
		v, err := ToTimeSliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]net.IP:
		v, err := ToIPSliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]*regexp.Regexp:
		v, err := ToRegexSliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
	case *[]map[string]string:
		v, err := ToMapStringStringSliceE(src)
		if err != nil {
			return err
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))

	default:
		return errors.Errorf("unsupport dst type [%v]", reflect.TypeOf(dst))
	}

	return nil
}
