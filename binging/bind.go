package binging

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refex"
)

func Bind(c Getter, v interface{}) error {
	b, err := Compile(v)
	if err != nil {
		return errors.Wrap(err, "bind failed")
	}
	return b.Bind(c, v)
}

func MustCompile(v interface{}) *Binder {
	rule, err := Compile(v)
	if err != nil {
		panic(err)
	}
	return rule
}

func Compile(v interface{}) (*Binder, error) {
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	infos := map[string]info{}
	if err := interfaceToBindRecursive(infos, rt, ""); err != nil {
		return nil, errors.Wrap(err, "compile failed")
	}
	return &Binder{infos: infos}, nil
}

type info struct {
	key      string
	dftVal   interface{}
	required bool
}

type Binder struct {
	infos map[string]info
}

func (b *Binder) Bind(c Getter, v interface{}) error {
	return bindRecursive(c, b.infos, v, "", "")
}

func bindRecursive(c Getter, infos map[string]info, v interface{}, prefix1 string, prefix2 string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &Error{Code: ErrInvalidDstType, Err: errors.Errorf("invalid value [%v] type or value is nil", reflect.TypeOf(v)), Key: prefix1}
	}
	rv = rv.Elem()

	info := infos[prefix1]
	if rv.Kind() == reflect.Struct {
		for i := 0; i < rv.NumField(); i++ {
			key := rv.Type().Field(i).Name
			switch rv.Field(i).Type().Kind() {
			case reflect.Ptr:
				if rv.Field(i).IsNil() {
					nv := reflect.New(rv.Field(i).Type().Elem())
					rv.Field(i).Set(nv)
				}
				if rv.Type().Field(i).Anonymous {
					if err := bindRecursive(c, infos, rv.Field(i).Interface(), prefix1, prefix2); err != nil {
						return err
					}
				} else {
					if err := bindRecursive(c, infos, rv.Field(i).Interface(), prefixAppendKey(prefix1, key), prefixAppendKey(prefix2, info.key)); err != nil {
						return err
					}
				}
			default:
				if rv.Type().Field(i).Anonymous {
					if err := bindRecursive(c, infos, rv.Field(i).Addr().Interface(), prefix1, prefix2); err != nil {
						return err
					}
				} else {
					if err := bindRecursive(c, infos, rv.Field(i).Addr().Interface(), prefixAppendKey(prefix1, key), prefixAppendKey(prefix2, info.key)); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}

	prefix2 = prefixAppendKey(prefix2, info.key)
	src, ok := c.Get(prefix2)
	if !ok {
		if info.required {
			return &Error{Code: ErrMissingRequiredField, Key: prefix1, Err: errors.New("prefix not exists")}
		}
		if info.dftVal != nil {
			src = info.dftVal
		}
	}
	if src != nil {
		if err := refex.SetInterface(v, src); err != nil {
			return &Error{Code: ErrInvalidFormat, Key: prefix1, Err: errors.Wrap(err, "set interface failed")}
		}
	}

	return nil
}

func interfaceToBindRecursive(infos map[string]info, rt reflect.Type, prefix string) error {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	for i := 0; i < rt.NumField(); i++ {
		t := rt.Field(i).Type
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		tag := rt.Field(i).Tag
		if tag.Get("bind") == "-" {
			continue
		}

		key := rt.Field(i).Name
		info, err := parseTag(key, tag.Get("bind"), tag.Get("dft"), rt.Field(i).Type)
		if err != nil {
			return errors.Wrap(err, "bind failed")
		}
		infos[prefixAppendKey(prefix, key)] = *info
		if t.Kind() == reflect.Struct {
			if rt.Field(i).Anonymous {
				if err := interfaceToBindRecursive(infos, t, prefix); err != nil {
					return err
				}
			} else {
				if err := interfaceToBindRecursive(infos, t, prefixAppendKey(prefix, key)); err != nil {
					return err
				}
			}
			continue
		}
	}

	return nil
}

func prefixAppendKey(prefix string, key string) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%v.%v", prefix, key)
}

func parseTag(key string, bTag string, dTag string, rt reflect.Type) (*info, error) {
	info := &info{}
	vals := strings.Split(bTag, ";")
	for _, val := range vals {
		val = strings.TrimSpace(val)
		if val == "required" {
			info.required = true
		} else {
			info.key = val
		}
	}
	if info.key == "" {
		info.key = key
	}
	if dTag == "" {
		return info, nil
	}
	if rt.Kind() == reflect.Ptr {
		val := reflect.New(rt.Elem())
		if err := refex.SetInterface(val.Interface(), dTag); err != nil {
			return nil, errors.Wrapf(err, "set interface failed. type [%v] tag [%v]", rt, dTag)
		}
		info.dftVal = val.Elem().Interface()
	} else {
		val := reflect.New(rt)
		if err := refex.SetInterface(val.Interface(), dTag); err != nil {
			return nil, errors.Wrapf(err, "set interface failed. type [%v] tag [%v]", rt, dTag)
		}
		info.dftVal = val.Elem().Interface()
	}

	return info, nil
}
