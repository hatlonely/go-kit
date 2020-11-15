package refx

import (
	"reflect"
	"time"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/cast"
)

func SetDefaultValue(v interface{}) error {
	if reflect.ValueOf(v).IsNil() {
		return nil
	}

	if reflect.TypeOf(v).Kind() != reflect.Ptr || reflect.TypeOf(v).Elem().Kind() != reflect.Struct {
		return nil
	}
	rt := reflect.TypeOf(v).Elem()
	rv := reflect.ValueOf(v).Elem()
	for i := 0; i < rv.NumField(); i++ {
		tag := rt.Field(i).Tag.Get("dft")
		if !rv.Field(i).CanSet() {
			continue
		}

		if rt.Field(i).Type.Kind() == reflect.Struct && rt.Field(i).Type != reflect.TypeOf(time.Time{}) {
			if err := SetDefaultValue(rv.Field(i).Addr().Interface()); err != nil {
				return errors.WithMessage(err, "SetDefaultValue failed")
			}
			continue
		}
		if tag == "" {
			continue
		}
		if err := cast.SetInterface(rv.Field(i).Addr().Interface(), tag); err != nil {
			return errors.WithMessage(err, "SetInterface failed")
		}
	}

	return nil
}
