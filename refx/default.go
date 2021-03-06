package refx

import (
	"reflect"
	"regexp"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/cast"
)

var mutex sync.RWMutex
var defaultValueMap = map[reflect.Type]reflect.Value{}

func SetDefaultValueCopyP(v interface{}) {
	if err := SetDefaultValueCopy(v); err != nil {
		panic(err)
	}
}

func SetDefaultValueCopy(v interface{}) error {
	if reflect.ValueOf(v).IsNil() {
		return nil
	}

	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Struct {
		return errors.Errorf("expect a struct point. got [%v]", rt)
	}

	mutex.RLock()
	if rv, ok := defaultValueMap[rt]; ok {
		mutex.RUnlock()
		reflect.ValueOf(v).Elem().Set(rv)
		return nil
	}
	mutex.RUnlock()

	if err := SetDefaultValue(v); err != nil {
		return err
	}

	nv := reflect.New(reflect.TypeOf(v).Elem())
	nv.Elem().Set(reflect.ValueOf(v).Elem())

	mutex.Lock()
	defaultValueMap[rt] = nv.Elem()
	mutex.Unlock()
	return nil
}

func SetDefaultValueP(v interface{}) {
	if err := SetDefaultValue(v); err != nil {
		panic(err)
	}
}

func SetDefaultValue(v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr || reflect.TypeOf(v).Elem().Kind() != reflect.Struct {
		return errors.Errorf("expect a struct point. got [%v]", reflect.TypeOf(v))
	}
	rt := reflect.TypeOf(v).Elem()
	rv := reflect.ValueOf(v).Elem()
	for i := 0; i < rv.NumField(); i++ {
		tag := rt.Field(i).Tag.Get("dft")
		if !rv.Field(i).CanSet() {
			continue
		}
		if rt.Field(i).Type.Kind() == reflect.Struct && rt.Field(i).Type != reflect.TypeOf(time.Time{}) && rt.Field(i).Type != reflect.TypeOf(regexp.Regexp{}) {
			if err := SetDefaultValue(rv.Field(i).Addr().Interface()); err != nil {
				return errors.WithMessage(err, "SetDefaultValueCopy failed")
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
