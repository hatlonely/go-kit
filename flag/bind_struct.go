package flag

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hatlonely/go-kit/strex"
)

func (f *Flag) Struct(v interface{}) error {
	return f.bindStructRecursive(v, "")
}

func (f *Flag) bindStructRecursive(v interface{}, prefix string) error {
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
		t := rt.Field(i).Type

		if tag == "-" {
			continue
		}
		if t.Kind() == reflect.Struct && t != reflect.TypeOf(time.Time{}) {
			key := tag
			if key == "" {
				key = strex.KebabName(rt.Field(i).Name)
			}
			if prefix != "" {
				key = prefix + "-" + key
			}
			if err := f.bindStructRecursive(rv.Field(i).Addr().Interface(), key); err != nil {
				return err
			}
		} else if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct && t != reflect.TypeOf(&time.Time{}) {
			rv.Field(i).Set(reflect.New(rv.Field(i).Type().Elem()))
			key := tag
			if key == "" {
				key = strex.KebabName(rt.Field(i).Name)
			}
			if prefix != "" {
				key = prefix + "-" + key
			}
			if err := f.bindStructRecursive(rv.Field(i).Interface(), key); err != nil {
				return err
			}
		} else {
			name, shorthand, usage, required, defaultValue, isArgument, err := parseTag(tag)
			if err != nil {
				return err
			}
			if name == "" {
				name = strex.KebabName(rt.Field(i).Name)
			}
			if prefix != "" {
				name = prefix + "-" + name
			}
			if err := f.BindFlag(rv.Field(i).Addr().Interface(), name, usage, rv.Field(i).Type(), required, shorthand, defaultValue, isArgument); err != nil {
				return err
			}
		}
	}

	return nil
}

func parseTag(tag string) (name string, shorthand string, usage string, required bool, defaultValue string, isArgument bool, err error) {
	if strings.Trim(tag, " ") == "" {
		isArgument = false
		return
	}
	for _, field := range strings.Split(tag, ";") {
		field = strings.Trim(field, " ")
		if field == "required" { // required
			required = true
		} else if strings.HasPrefix(field, "--") { // --int-option, -i
			names := strings.Split(field, ",")
			name = strings.Trim(names[0], " ")[2:]
			if len(names) > 2 {
				err = fmt.Errorf("expected name field format is '--<name>[, -<shorthand>]', got [%v]", field)
				return
			} else if len(names) == 2 {
				shorthand = strings.Trim(names[1], " ")
				if !strings.HasPrefix(shorthand, "-") {
					err = fmt.Errorf("expected name field format is '--<name>[, -<shorthand>]', got [%v]", field)
					return
				}
				shorthand = shorthand[1:]
			}
		} else if strings.Contains(field, ":") { // default: 10; usage: int flag
			kvs := strings.Split(field, ":")
			if len(kvs) != 2 {
				err = fmt.Errorf("expected format '<key>:<value>', got [%v]", field)
				return
			}
			key := strings.Trim(kvs[0], " ")
			val := strings.Trim(kvs[1], " ")
			switch key {
			case "default":
				defaultValue = val
			case "usage":
				usage = val
			}
		} else { // pos
			name = strings.Trim(field, " ")
			isArgument = true
		}
	}

	return
}
