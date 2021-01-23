package refx

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	playgroundValidator "gopkg.in/go-playground/validator.v9"

	"github.com/hatlonely/go-kit/cast"
	"github.com/hatlonely/go-kit/strx"
	"github.com/hatlonely/go-kit/validator"
)

type Options struct {
	CamelName           bool
	SnakeName           bool
	KebabName           bool
	PascalName          bool
	DisableDefaultValue bool

	validators []func(v interface{}) error
}

type Option func(options *Options)

func WithCamelName() Option {
	return func(options *Options) {
		options.CamelName = true
	}
}

func WithSnakeName() Option {
	return func(options *Options) {
		options.SnakeName = true
	}
}

func WithKebabName() Option {
	return func(options *Options) {
		options.KebabName = true
	}
}

func WithPascalName() Option {
	return func(options *Options) {
		options.PascalName = true
	}
}

func WithPlaygroundValidator() Option {
	validate := playgroundValidator.New()
	return func(options *Options) {
		options.validators = append(options.validators, validate.Struct)
	}
}

func WithDefaultValidator() Option {
	return func(options *Options) {
		options.validators = append(options.validators, validator.Validate)
	}
}

func WithDisableDefaultValue() Option {
	return func(opt *Options) {
		opt.DisableDefaultValue = true
	}
}

func NewOptions(opts ...Option) *Options {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}
	return &options
}

func (o *Options) FormatKey(str string) string {
	if o.CamelName {
		return strx.CamelName(str)
	}
	if o.SnakeName {
		return strx.SnakeName(str)
	}
	if o.KebabName {
		return strx.KebabName(str)
	}
	if o.PascalName {
		return strx.PascalName(str)
	}
	return str
}

func (o *Options) Validate(v interface{}) error {
	for _, validate := range o.validators {
		if err := validate(v); err != nil {
			return errors.Wrap(err, "validate failed")
		}
	}
	return nil
}

func InterfaceGet(v interface{}, key string) (interface{}, error) {
	return interfaceGetRecursive(v, key, "")
}

func InterfaceSet(pv *interface{}, key string, val interface{}) error {
	return interfaceSetRecursive(pv, key, val, "")
}

func InterfaceToStruct(src interface{}, dst interface{}, opts ...Option) error {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	return InterfaceToStructWithOptions(src, dst, &options)
}

func InterfaceToStructWithOptions(src interface{}, dst interface{}, options *Options) error {
	if err := interfaceToStructRecursive(src, dst, "", options); err != nil {
		return err
	}

	return options.Validate(dst)
}

func InterfaceDiff(v1 interface{}, v2 interface{}) ([]string, error) {
	var keys []string
	if err := InterfaceTravel(v1, func(key string, val1 interface{}) error {
		val2, err := InterfaceGet(v2, key)
		if err != nil || val1 != val2 {
			keys = append(keys, key)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return keys, nil
}

// todo remove value if map or slice is empty, use recursive implement
func InterfaceDel(pv *interface{}, key string) error {
	info, prev, err := getLastToken(key)
	if err != nil {
		return err
	}
	v, err := InterfaceGet(*pv, prev)
	if err != nil {
		return nil
	}
	if v == nil {
		return nil
	}
	if info.mod == ArrMod {
		val, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("unsupport slice type. prefix: [%v], type: [%v]", prev, reflect.TypeOf(*pv))
		}
		if info.idx >= len(val) {
			return nil
		}
		val = append(val[:info.idx], val[info.idx+1:]...)
		return InterfaceSet(pv, prev, val)
	}

	switch v.(type) {
	case map[string]interface{}:
		delete(v.(map[string]interface{}), info.key)
	case map[interface{}]interface{}:
		delete(v.(map[interface{}]interface{}), info.key)
	default:
		return fmt.Errorf("unsupport slice type. prefix: [%v], type: [%v]", prev, reflect.TypeOf(*pv))
	}

	return nil
}

func interfaceSetRecursive(pv *interface{}, key string, val interface{}, prefix string) error {
	if key == "" {
		*pv = val
		return nil
	}
	info, next, err := getToken(key)
	if err != nil {
		return fmt.Errorf("get token failed. prefix: [%v], err: [%v]", prefix, err)
	}
	if info.mod == ArrMod {
		if *pv == nil {
			*pv = []interface{}{}
		}
		v, ok := (*pv).([]interface{})
		if !ok {
			return fmt.Errorf("unsupport slice type. prefix: [%v], type: [%v]", prefix, reflect.TypeOf(*pv))
		}
		if info.idx > len(v) {
			return fmt.Errorf("index out of bounds. prefix: [%v], index: [%v]", prefix, info.idx)
		}
		if info.idx < len(v) {
			sub := v[info.idx]
			if err := interfaceSetRecursive(&sub, next, val, prefixAppendIdx(prefix, info.idx)); err != nil {
				return err
			}
			v[info.idx] = sub
			return nil
		}

		var sub interface{}
		if err := interfaceSetRecursive(&sub, next, val, prefixAppendIdx(prefix, info.idx)); err != nil {
			return err
		}
		v = append(v, sub)
		*pv = v
		return nil
	}

	if *pv == nil {
		*pv = map[string]interface{}{}
	}
	switch (*pv).(type) {
	case map[string]interface{}:
		sub := (*pv).(map[string]interface{})[info.key]
		if err := interfaceSetRecursive(&sub, next, val, prefixAppendKey(prefix, info.key)); err != nil {
			return err
		}
		(*pv).(map[string]interface{})[info.key] = sub
	case map[interface{}]interface{}:
		sub := (*pv).(map[interface{}]interface{})[info.key]
		if err := interfaceSetRecursive(&sub, next, val, prefixAppendKey(prefix, info.key)); err != nil {
			return err
		}
		(*pv).(map[interface{}]interface{})[info.key] = sub
	default:
		return fmt.Errorf("unsupport map type. prefix: [%v], type: [%v]", prefix, reflect.TypeOf(*pv))
	}

	return nil
}

func interfaceGetRecursive(v interface{}, key string, prefix string) (interface{}, error) {
	if v == nil {
		return nil, fmt.Errorf("no such key. prefix: [%v], key: [%v]", prefix, key)
	}
	if key == "" {
		return v, nil
	}
	info, next, err := getToken(key)
	if err != nil {
		return nil, fmt.Errorf("get token failed. prefix: [%v], err: [%v]", prefix, err)
	}

	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)

	if info.mod == ArrMod {
		switch rt.Kind() {
		case reflect.Slice:
			if info.idx >= rv.Len() {
				return nil, errors.Errorf("index out of bounds. prefix: [%v], index: [%v]", prefix, info.idx)
			}
			return interfaceGetRecursive(rv.Index(info.idx).Interface(), next, prefixAppendIdx(prefix, info.idx))
		default:
			return nil, errors.Errorf("node is not a slice. prefix: [%v], type: [%v]", prefix, reflect.TypeOf(v))
		}
	}

	switch rt.Kind() {
	case reflect.Map:
		if rv.MapIndex(reflect.ValueOf(info.key)).Kind() == reflect.Invalid {
			return nil, errors.Errorf("no such key. prefix [%v], key [%v]", prefix, key)
		}
		return interfaceGetRecursive(rv.MapIndex(reflect.ValueOf(info.key)).Interface(), next, prefixAppendKey(prefix, info.key))
	default:
		return nil, errors.Errorf("node is not a map. prefix: [%v], type: [%v]", prefix, reflect.TypeOf(v))
	}
}

func interfaceToStructRecursive(src interface{}, dst interface{}, prefix string, options *Options) error {
	if reflect.ValueOf(dst).Kind() != reflect.Ptr {
		return errors.Errorf("invalid dst type [%v] or dst is nil. prefix: [%v]", reflect.TypeOf(dst), prefix)
	}
	if dst == nil {
		return errors.Errorf("dst should not be nil. prefix: [%v]", prefix)
	}

	rv := reflect.ValueOf(dst).Elem()
	rt := reflect.TypeOf(dst).Elem()

	if rt.Kind() == reflect.Struct && !options.DisableDefaultValue {
		if err := SetDefaultValue(dst); err != nil {
			return errors.WithMessagef(err, "SetDefaultValue failed. prefix: [%v]", prefix)
		}
	}

	if src == nil {
		return nil
	}

	switch dst.(type) {
	case *time.Time, **regexp.Regexp:
		if err := cast.SetInterface(dst, src); err != nil {
			return errors.WithMessagef(err, "cast.SetInterface failed. prefix: [%v]", prefix)
		}
		return nil
	}

	if rt.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rt.Elem()))
		}
		return interfaceToStructRecursive(src, rv.Interface(), prefix, options)
	}
	if rt.Kind() == reflect.Interface {
		rv.Set(reflect.ValueOf(src))
		return nil
	}

	srv := reflect.ValueOf(src)
	srt := reflect.TypeOf(src)
	switch rt.Kind() {
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			key := options.FormatKey(rt.Field(i).Name)
			var val interface{}
			switch srt.Kind() {
			case reflect.Map:
				if srv.MapIndex(reflect.ValueOf(key)).Kind() == reflect.Invalid {
					continue
				}
				val = srv.MapIndex(reflect.ValueOf(key)).Interface()
			default:
				return errors.Errorf("cannot convert type [%v] to [%v]. prefix: [%v]", reflect.TypeOf(src), reflect.TypeOf(dst), prefix)
			}

			if rv.Field(i).Kind() == reflect.Ptr && rv.Field(i).IsNil() {
				nv := reflect.New(rv.Field(i).Type().Elem())
				rv.Field(i).Set(nv)
			}
			if err := interfaceToStructRecursive(val, rv.Field(i).Addr().Interface(), prefixAppendKey(prefix, key), options); err != nil {
				return err
			}
		}
	case reflect.Slice:
		if srt.Kind() != reflect.Slice {
			return errors.Errorf("cannot convert type [%v] to [%v]. prefix: [%v]", reflect.TypeOf(src), reflect.TypeOf(dst), prefix)
		}
		rv.Set(reflect.MakeSlice(rt, 0, srv.Len()))

		for idx := 0; idx < srv.Len(); idx++ {
			nv := reflect.New(rt.Elem())
			if err := interfaceToStructRecursive(srv.Index(idx).Interface(), nv.Interface(), prefixAppendIdx(prefix, idx), options); err != nil {
				return err
			}
			rv.Set(reflect.Append(rv, nv.Elem()))
		}
	case reflect.Map:
		if rv.IsNil() {
			rv.Set(reflect.MakeMap(rt))
		}
		for _, key := range srv.MapKeys() {
			newKey := reflect.New(rt.Key())
			if err := cast.SetInterface(newKey.Interface(), key.Interface()); err != nil {
				return errors.WithMessage(err, "cast.SetInterface failed")
			}

			val := srv.MapIndex(key).Interface()
			nv := reflect.New(rt.Elem())
			if err := interfaceToStructRecursive(val, nv.Interface(), prefixAppendKey(prefix, cast.ToString(newKey.Elem().Interface())), options); err != nil {
				return err
			}
			rv.SetMapIndex(newKey.Elem(), nv.Elem())
		}
	default:
		if err := cast.SetInterface(dst, src); err != nil {
			return errors.WithMessagef(err, "cast.SetInterface failed. prefix: [%v]", prefix)
		}
	}

	return nil
}

func InterfaceTravel(v interface{}, fun func(key string, val interface{}) error) error {
	return interfaceTravelRecursive(v, fun, "")
}

// return error if fun is error immediately
func interfaceTravelRecursive(v interface{}, fun func(key string, val interface{}) error, prefix string) error {
	if v == nil {
		return nil
	}

	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)

	switch rt.Kind() {
	case reflect.Ptr:
		return errors.Errorf("key [%v], unsupported type [%v]", prefix, reflect.TypeOf(v))
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			val := rv.MapIndex(key).Interface()
			if err := interfaceTravelRecursive(val, fun, prefixAppendKey(prefix, cast.ToString(key.Interface()))); err != nil {
				return err
			}
		}
	case reflect.Slice:
		for idx := 0; idx < rv.Len(); idx++ {
			if err := interfaceTravelRecursive(rv.Index(idx).Interface(), fun, prefixAppendIdx(prefix, idx)); err != nil {
				return err
			}
		}
	default:
		if err := fun(prefix, v); err != nil {
			return err
		}
	}

	return nil
}

const MapMod = 1
const ArrMod = 2

type KeyInfo struct {
	key string
	idx int
	mod int
}

func prefixAppendKey(prefix string, key string) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%v.%v", prefix, key)
}

func prefixAppendIdx(prefix string, idx int) string {
	if prefix == "" {
		return fmt.Sprintf("[%v]", idx)
	}
	return fmt.Sprintf("%v[%v]", prefix, idx)
}

func getLastToken(key string) (info KeyInfo, prev string, err error) {
	if key[len(key)-1] == ']' {
		pos := strings.LastIndex(key, "[")
		// "123]" => error
		if pos == -1 {
			return info, "", fmt.Errorf("miss '[' in key. key: [%v]", key)
		}
		sub := key[pos+1 : len(key)-1]
		// "[]" => error
		if sub == "" {
			return info, "", fmt.Errorf("idx should not be empty. key: [%v]", key)
		}
		// "[abc]" => error
		idx, err := strconv.Atoi(sub)
		if err != nil {
			return info, "", fmt.Errorf("idx to int fail. key: [%v], sub: [%v]", key, sub)
		}
		// "key[3]" => 3, "key"
		return KeyInfo{idx: idx, mod: ArrMod}, key[:pos], nil
	}
	pos := strings.LastIndex(key, ".")
	// "key" => "key", ""
	if pos == -1 {
		return KeyInfo{key: key, mod: MapMod}, "", nil
	}
	// "key1.key2." => error
	if key[pos+1:] == "" {
		return info, "", fmt.Errorf("key should not be empty. key: [%v]", key)
	}
	// "key1[3].key2" => "key2", "key1[3]"
	return KeyInfo{key: key[pos+1:], mod: MapMod}, key[:pos], nil
}

func getToken(key string) (info KeyInfo, next string, err error) {
	if key[0] == '[' {
		pos := strings.Index(key, "]")
		// "[123" => error
		if pos == -1 {
			return info, next, fmt.Errorf("miss ']' in key. key: [%v]", key)
		}
		// "[]" => error
		if key[1:pos] == "" {
			return info, next, fmt.Errorf("idx should not be empty. key: [%v]", key)
		}
		idx, err := strconv.Atoi(key[1:pos])
		// "[abc]" => error
		if err != nil {
			return info, next, fmt.Errorf("idx to int fail. key: [%v], sub: [%v]", key, key[1:pos])
		}
		// "[1].key" => "1", "key"
		if pos+1 < len(key) && key[pos+1] == '.' {
			return KeyInfo{idx: idx, mod: ArrMod}, key[pos+2:], nil
		}
		// "[1][2]" => 1, "[2]"
		return KeyInfo{idx: idx, mod: ArrMod}, key[pos+1:], nil
	}
	pos := strings.IndexAny(key, ".[")
	// "key" => "key", ""
	if pos == -1 {
		return KeyInfo{key: key, mod: MapMod}, "", nil
	}
	// "key[0]" => "key", "[0]"
	if key[pos] == '[' {
		return KeyInfo{key: key[:pos], mod: MapMod}, key[pos:], nil
	}
	// ".key1.key2" => error
	if key[:pos] == "" {
		return info, "", fmt.Errorf("key should not be empty. key: [%v]", key)
	}
	// "key1.key2.key3" => "key1", "key2.key3"
	return KeyInfo{key: key[:pos], mod: MapMod}, key[pos+1:], nil
}
