package bind

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/cast"
	"github.com/hatlonely/go-kit/refx"
)

func Bind(v interface{}, getters []Getter, opts ...refx.Option) error {
	options := refx.NewOptions(opts...)
	if _, err := bindRecursive(v, "", getters, options); err != nil {
		return err
	}

	return options.Validate(v)
}

func bindRecursive(v interface{}, prefix string, getters []Getter, options *refx.Options) (int, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return 0, &Error{Code: ErrInvalidDstType, Err: errors.Errorf("invalid value [%v] type or value is nil", reflect.TypeOf(v)), Key: prefix}
	}
	rv = rv.Elem()

	if rv.Kind() == reflect.Struct {
		if !options.DisableDefaultValue {
			if err := refx.SetDefaultValue(v); err != nil {
				return 0, errors.Wrap(err, "SetDefaultValue failed")
			}
		}

		count := 0
		for i := 0; i < rv.NumField(); i++ {
			// ignore unexported field
			if !rv.Field(i).CanInterface() {
				continue
			}
			key := rv.Type().Field(i).Name
			switch rv.Field(i).Type().Kind() {
			case reflect.Ptr:
				if rv.Field(i).IsNil() {
					nv := reflect.New(rv.Field(i).Type().Elem())
					rv.Field(i).Set(nv)
				}
				if rv.Type().Field(i).Anonymous {
					if cnt, err := bindRecursive(rv.Field(i).Interface(), prefix, getters, options); err != nil {
						return 0, err
					} else {
						count += cnt
					}
				} else {
					if cnt, err := bindRecursive(rv.Field(i).Interface(), prefixAppendKey(prefix, options.FormatKey(key)), getters, options); err != nil {
						return 0, err
					} else {
						count += cnt
					}
				}
			default:
				if rv.Type().Field(i).Anonymous {
					if cnt, err := bindRecursive(rv.Field(i).Addr().Interface(), prefix, getters, options); err != nil {
						return 0, err
					} else {
						count += cnt
					}
				} else {
					if cnt, err := bindRecursive(rv.Field(i).Addr().Interface(), prefixAppendKey(prefix, options.FormatKey(key)), getters, options); err != nil {
						return 0, err
					} else {
						count += cnt
					}
				}
			}
		}
		return count, nil
	} else if rv.Kind() == reflect.Slice {
		count := 0
		rt := rv.Type()
		for i := 0; ; i++ {
			switch rt.Elem().Kind() {
			case reflect.Ptr:
				nv := reflect.New(rt.Elem().Elem())
				if cnt, err := bindRecursive(nv.Interface(), prefixAppendIdx(prefix, i), getters, options); err != nil {
					return 0, err
				} else if cnt != 0 {
					count += cnt
					rv.Set(reflect.Append(rv, nv.Elem().Addr()))
				} else {
					return count, nil
				}
			default:
				nv := reflect.New(rt.Elem())
				if cnt, err := bindRecursive(nv.Interface(), prefixAppendIdx(prefix, i), getters, options); err != nil {
					return 0, err
				} else if cnt != 0 {
					count += cnt
					rv.Set(reflect.Append(rv, nv.Elem()))
				} else {
					return count, nil
				}
			}
		}
	}

	for _, getter := range getters {
		if getter == nil {
			continue
		}
		src, ok := getter.Get(prefix)
		if !ok {
			continue
		}
		if err := cast.SetInterface(v, src); err != nil {
			return 0, &Error{Code: ErrInvalidFormat, Key: prefix, Err: errors.Wrap(err, "set interface failed")}
		}
		return 1, nil
	}

	return 0, nil
}

func prefixAppendIdx(prefix string, idx int) string {
	if prefix == "" {
		return fmt.Sprintf("[%v]", idx)
	}
	return fmt.Sprintf("%v[%v]", prefix, idx)
}

func prefixAppendKey(prefix string, key string) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%v.%v", prefix, key)
}
