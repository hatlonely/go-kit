package flag

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

func (f *Flag) Struct(v interface{}, opts ...refx.Option) error {
	var options refx.Options
	for _, opt := range opts {
		opt(&options)
	}
	return f.bindStructRecursive(v, "", &options)
}

func (f *Flag) bindStructRecursive(v interface{}, prefixKey string, options *refx.Options) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("expected a pointer, got [%v]", reflect.TypeOf(v))
	}

	rv := reflect.ValueOf(v).Elem()
	rt := reflect.TypeOf(v).Elem()

	if rt.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got [%v]", rt)
	}

	for i := 0; i < rv.NumField(); i++ {
		if !rv.Field(i).CanSet() {
			continue
		}
		tag := rt.Field(i).Tag.Get("flag")
		typ := rt.Field(i).Type
		key := rt.Field(i).Name

		if tag == "-" {
			continue
		}
		if typ.Kind() == reflect.Struct && typ != reflect.TypeOf(time.Time{}) {
			if rv.Type().Field(i).Anonymous {
				if err := f.bindStructRecursive(rv.Field(i).Addr().Interface(), prefixKey, options); err != nil {
					return err
				}
			} else {
				if err := f.bindStructRecursive(rv.Field(i).Addr().Interface(), prefixAppendKey(prefixKey, refx.FormatKeyWithOptions(key, options)), options); err != nil {
					return err
				}
			}
		} else if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct && typ != reflect.TypeOf(&time.Time{}) {
			rv.Field(i).Set(reflect.New(rv.Field(i).Type().Elem()))
			if rv.Type().Field(i).Anonymous {
				if err := f.bindStructRecursive(rv.Field(i).Interface(), prefixKey, options); err != nil {
					return err
				}
			} else {
				if err := f.bindStructRecursive(rv.Field(i).Interface(), prefixAppendKey(prefixKey, refx.FormatKeyWithOptions(key, options)), options); err != nil {
					return err
				}
			}
		} else {
			options, err := parseTag(tag, key, prefixKey, typ, options)
			if err != nil {
				return err
			}
			if err := f.BindFlagWithOptions(rv.Field(i).Addr().Interface(), options); err != nil {
				return err
			}
		}
	}

	return nil
}

var reKey = regexp.MustCompile(`[.\w_-]+`)

func parseTag(tag string, key string, prefixKey string, typ reflect.Type, ropt *refx.Options) (*AddFlagOptions, error) {
	var options AddFlagOptions
	options.Key = prefixAppendKey(prefixKey, refx.FormatKeyWithOptions(key, ropt))

	tag = strings.TrimSpace(tag)
	for _, field := range strings.Split(tag, ";") {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		if field == "required" { // required
			options.Required = true
		} else if field == "isArgument" {
			options.IsArgument = true
		} else if strings.HasPrefix(field, "-") { // --int-option, -i
			for _, name := range strings.Split(field, ",") {
				name = strings.TrimSpace(name)
				if strings.HasPrefix(name, "--") {
					name = name[2:]
					if !reKey.Match([]byte(name)) {
						return nil, errors.Errorf("invalid key format, key [%v], name [%v]", options.Key, name)
					}
					options.Name = name
					continue
				}
				if strings.HasPrefix(field, "-") {
					name = name[1:]
					if !reKey.Match([]byte(name)) {
						return nil, errors.Errorf("invalid key format, key [%v], name [%v]", options.Key, name)
					}
					options.Shorthand = name
					continue
				}
				return nil, errors.Errorf("invalid key format, key [%v], name [%v]", options.Key, name)
			}
		} else if strings.Contains(field, ":") { // default: 10; usage: int flag
			idx := strings.Index(field, ":")
			key := strings.Trim(field[:idx], " ")
			val := strings.Trim(field[idx+1:], " ")
			switch key {
			case "default":
				options.DefaultValue = val
			case "usage":
				options.Usage = val
			}
		} else if reKey.Match([]byte(field)) { // pos
			options.Name = field
			options.IsArgument = true
		} else {
			return nil, errors.Errorf("invalid key format, key [%v], field [%v]", options.Key, field)
		}
	}

	if options.Name == "" {
		options.Name = options.Key
	}
	options.Type = typ

	return &options, nil
}

func prefixAppendKey(prefix string, key string) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%v.%v", prefix, key)
}

func prefixAppendName(prefix string, key string) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%v-%v", prefix, key)
}
