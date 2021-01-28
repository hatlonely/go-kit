package bind

import (
	"fmt"
	"reflect"
	"regexp"
	"time"

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

	switch v.(type) {
	case *time.Time, **regexp.Regexp:
	default:
		count := 0
		switch rv.Kind() {
		case reflect.Ptr:
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			return bindRecursive(rv.Interface(), prefix, getters, options)
		case reflect.Struct:
			if !options.DisableDefaultValue {
				if err := refx.SetDefaultValue(v); err != nil {
					return 0, errors.Wrap(err, "SetDefaultValue failed")
				}
			}

			for i := 0; i < rv.NumField(); i++ {
				// ignore unexported field
				if !rv.Field(i).CanInterface() {
					continue
				}
				key := rv.Type().Field(i).Name
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
			return count, nil
		case reflect.Slice:
			for i := 0; ; i++ {
				nv := reflect.New(rv.Type().Elem())
				if cnt, err := bindRecursive(nv.Interface(), prefixAppendIdx(prefix, i), getters, options); err != nil {
					return 0, err
				} else if cnt != 0 {
					count += cnt
					rv.Set(reflect.Append(rv, nv.Elem()))
				} else {
					return count, nil
				}
			}
		case reflect.Map:
			if rv.IsNil() {
				rv.Set(reflect.MakeMap(rv.Type()))
			}
			for _, getter := range getters {
				if getter == nil {
					continue
				}
				src, ok := getter.Get(prefix)
				if !ok {
					continue
				}
				srv := reflect.ValueOf(src)
				for _, key := range srv.MapKeys() {
					newKey := reflect.New(rv.Type().Key())
					if err := cast.SetInterface(newKey.Interface(), key.Interface()); err != nil {
						return 0, errors.WithMessage(err, "cast.SetInterface failed")
					}

					nv := reflect.New(rv.Type().Elem())
					if cnt, err := bindRecursive(nv.Interface(), prefixAppendKey(prefix, cast.ToString(newKey.Elem().Interface())), getters, options); err != nil {
						return 0, err
					} else {
						count += cnt
					}
					rv.SetMapIndex(newKey.Elem(), nv.Elem())
				}
			}
			return count, nil
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
		switch v.(type) {
		case *interface{}:
			rv.Set(reflect.ValueOf(src))
		default:
			if err := cast.SetInterface(v, src); err != nil {
				return 0, &Error{Code: ErrInvalidFormat, Key: prefix, Err: errors.Wrap(err, "cast.SetInterface failed")}
			}
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
