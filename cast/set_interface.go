// this file is auto generate by autogen.py. do not edit!
package cast

import (
	"fmt"
	"net"
	"reflect"
	"time"
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

	default:
		return fmt.Errorf("unsupport dst type [%v]", reflect.TypeOf(dst))
	}

	return nil
}
