package flag

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/strx"
)

func (f *Flag) Struct(v interface{}, opts ...refx.Option) error {
	var options refx.Options
	for _, opt := range opts {
		opt(&options)
	}
	return f.bindStructRecursive(v, "", "", &options)
}

func (f *Flag) bindStructRecursive(v interface{}, prefixKey, prefixName string, options *refx.Options) error {
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
				if err := f.bindStructRecursive(rv.Field(i).Addr().Interface(), prefixKey, prefixName, options); err != nil {
					return err
				}
			} else {
				if err := f.bindStructRecursive(rv.Field(i).Addr().Interface(), prefixAppendKey(prefixKey, refx.FormatKeyWithOptions(key, options)), prefixAppendName(prefixName, strx.KebabName(key)), options); err != nil {
					return err
				}
			}
		} else if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct && typ != reflect.TypeOf(&time.Time{}) {
			rv.Field(i).Set(reflect.New(rv.Field(i).Type().Elem()))
			if rv.Type().Field(i).Anonymous {
				if err := f.bindStructRecursive(rv.Field(i).Interface(), prefixKey, prefixName, options); err != nil {
					return err
				}
			} else {
				if err := f.bindStructRecursive(rv.Field(i).Interface(), prefixAppendKey(prefixKey, refx.FormatKeyWithOptions(key, options)), prefixAppendName(prefixName, strx.KebabName(key)), options); err != nil {
					return err
				}
			}
		} else {
			options, err := parseTag(tag, key, prefixKey, prefixName, typ, options)
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

func parseTag(tag string, key string, prefixKey string, prefixName string, typ reflect.Type, ropt *refx.Options) (*AddFlagOptions, error) {
	var options AddFlagOptions
	tag = strings.TrimSpace(tag)
	for _, field := range strings.Split(tag, ";") {
		field = strings.Trim(field, " ")
		if field == "required" { // required
			options.Required = true
		} else if strings.HasPrefix(field, "--") { // --int-option, -i
			names := strings.Split(field, ",")
			options.Name = strings.Trim(names[0], " ")[2:]
			if len(names) > 2 {
				return nil, errors.Errorf("expected name field format is '--<name>[, -<shorthand>]', got [%v]", field)
			} else if len(names) == 2 {
				options.Shorthand = strings.TrimSpace(names[1])
				if !strings.HasPrefix(options.Shorthand, "-") {
					return nil, errors.Errorf("expected name field format is '--<name>[, -<shorthand>]', got [%v]", field)
				}
				options.Shorthand = options.Shorthand[1:]
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
		} else { // pos
			options.Name = strings.Trim(field, " ")
			options.IsArgument = true
		}
	}

	options.Key = prefixAppendKey(prefixKey, refx.FormatKeyWithOptions(key, ropt))
	if options.Name == "" {
		options.Name = strx.KebabName(key)
	}
	options.Name = prefixAppendName(prefixName, options.Name)
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
