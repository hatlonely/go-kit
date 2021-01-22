// this file is auto generate by autogen.py. do not edit!
package cast

import (
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func ToBoolSliceE(v interface{}) ([]bool, error) {
	switch v.(type) {
	case []bool:
		return v.([]bool), nil
	case []interface{}:
		vs := make([]bool, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToBoolE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]bool, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToBoolE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []bool
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToBoolE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []bool", reflect.TypeOf(v))
	}
}

func ToIntSliceE(v interface{}) ([]int, error) {
	switch v.(type) {
	case []int:
		return v.([]int), nil
	case []interface{}:
		vs := make([]int, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToIntE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]int, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToIntE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []int
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToIntE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []int", reflect.TypeOf(v))
	}
}

func ToUintSliceE(v interface{}) ([]uint, error) {
	switch v.(type) {
	case []uint:
		return v.([]uint), nil
	case []interface{}:
		vs := make([]uint, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToUintE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]uint, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToUintE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []uint
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToUintE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []uint", reflect.TypeOf(v))
	}
}

func ToInt64SliceE(v interface{}) ([]int64, error) {
	switch v.(type) {
	case []int64:
		return v.([]int64), nil
	case []interface{}:
		vs := make([]int64, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToInt64E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]int64, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToInt64E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []int64
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToInt64E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []int64", reflect.TypeOf(v))
	}
}

func ToInt32SliceE(v interface{}) ([]int32, error) {
	switch v.(type) {
	case []int32:
		return v.([]int32), nil
	case []interface{}:
		vs := make([]int32, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToInt32E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]int32, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToInt32E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []int32
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToInt32E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []int32", reflect.TypeOf(v))
	}
}

func ToInt16SliceE(v interface{}) ([]int16, error) {
	switch v.(type) {
	case []int16:
		return v.([]int16), nil
	case []interface{}:
		vs := make([]int16, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToInt16E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]int16, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToInt16E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []int16
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToInt16E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []int16", reflect.TypeOf(v))
	}
}

func ToInt8SliceE(v interface{}) ([]int8, error) {
	switch v.(type) {
	case []int8:
		return v.([]int8), nil
	case []interface{}:
		vs := make([]int8, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToInt8E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]int8, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToInt8E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []int8
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToInt8E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []int8", reflect.TypeOf(v))
	}
}

func ToUint64SliceE(v interface{}) ([]uint64, error) {
	switch v.(type) {
	case []uint64:
		return v.([]uint64), nil
	case []interface{}:
		vs := make([]uint64, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToUint64E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]uint64, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToUint64E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []uint64
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToUint64E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []uint64", reflect.TypeOf(v))
	}
}

func ToUint32SliceE(v interface{}) ([]uint32, error) {
	switch v.(type) {
	case []uint32:
		return v.([]uint32), nil
	case []interface{}:
		vs := make([]uint32, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToUint32E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]uint32, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToUint32E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []uint32
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToUint32E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []uint32", reflect.TypeOf(v))
	}
}

func ToUint16SliceE(v interface{}) ([]uint16, error) {
	switch v.(type) {
	case []uint16:
		return v.([]uint16), nil
	case []interface{}:
		vs := make([]uint16, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToUint16E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]uint16, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToUint16E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []uint16
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToUint16E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []uint16", reflect.TypeOf(v))
	}
}

func ToUint8SliceE(v interface{}) ([]uint8, error) {
	switch v.(type) {
	case []uint8:
		return v.([]uint8), nil
	case []interface{}:
		vs := make([]uint8, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToUint8E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]uint8, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToUint8E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []uint8
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToUint8E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []uint8", reflect.TypeOf(v))
	}
}

func ToFloat64SliceE(v interface{}) ([]float64, error) {
	switch v.(type) {
	case []float64:
		return v.([]float64), nil
	case []interface{}:
		vs := make([]float64, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToFloat64E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]float64, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToFloat64E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []float64
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToFloat64E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []float64", reflect.TypeOf(v))
	}
}

func ToFloat32SliceE(v interface{}) ([]float32, error) {
	switch v.(type) {
	case []float32:
		return v.([]float32), nil
	case []interface{}:
		vs := make([]float32, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToFloat32E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]float32, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToFloat32E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []float32
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToFloat32E(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []float32", reflect.TypeOf(v))
	}
}

func ToStringSliceE(v interface{}) ([]string, error) {
	switch v.(type) {
	case []string:
		return v.([]string), nil
	case []interface{}:
		vs := make([]string, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToStringE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []string
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToStringE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []string", reflect.TypeOf(v))
	}
}

func ToDurationSliceE(v interface{}) ([]time.Duration, error) {
	switch v.(type) {
	case []time.Duration:
		return v.([]time.Duration), nil
	case []interface{}:
		vs := make([]time.Duration, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToDurationE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]time.Duration, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToDurationE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []time.Duration
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToDurationE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []time.Duration", reflect.TypeOf(v))
	}
}

func ToTimeSliceE(v interface{}) ([]time.Time, error) {
	switch v.(type) {
	case []time.Time:
		return v.([]time.Time), nil
	case []interface{}:
		vs := make([]time.Time, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToTimeE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]time.Time, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToTimeE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []time.Time
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToTimeE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []time.Time", reflect.TypeOf(v))
	}
}

func ToIPSliceE(v interface{}) ([]net.IP, error) {
	switch v.(type) {
	case []net.IP:
		return v.([]net.IP), nil
	case []interface{}:
		vs := make([]net.IP, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToIPE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]net.IP, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToIPE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []net.IP
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToIPE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []net.IP", reflect.TypeOf(v))
	}
}

func ToMapStringStringSliceE(v interface{}) ([]map[string]string, error) {
	switch v.(type) {
	case []map[string]string:
		return v.([]map[string]string), nil
	case []interface{}:
		vs := make([]map[string]string, 0, len(v.([]interface{})))
		for _, i := range v.([]interface{}) {
			val, err := ToMapStringStringE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case []string:
		vs := make([]map[string]string, 0, len(v.([]string)))
		for _, i := range v.([]string) {
			val, err := ToMapStringStringE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []map[string]string
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToMapStringStringE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []map[string]string", reflect.TypeOf(v))
	}
}
