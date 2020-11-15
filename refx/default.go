package refx

import (
	"reflect"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/cast"
)

var mutex sync.RWMutex
var defaultValueMap = map[reflect.Type]reflect.Value{}

func SetDefaultValue(v interface{}) error {
	if reflect.ValueOf(v).IsNil() {
		return nil
	}

	rt := reflect.TypeOf(v)
	mutex.RLock()
	if rv, ok := defaultValueMap[rt]; ok {
		mutex.RUnlock()
		reflect.ValueOf(v).Elem().Set(rv)
		return nil
	}
	mutex.RUnlock()

	if err := setDefaultValue(v); err != nil {
		return err
	}

	mutex.Lock()
	defaultValueMap[rt] = reflect.ValueOf(v).Elem()
	mutex.Unlock()
	return nil
}

func setDefaultValue(v interface{}) error {
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
