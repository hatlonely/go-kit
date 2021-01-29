package refx

import (
	"reflect"

	"github.com/pkg/errors"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

type ConstructorInfo struct {
	HasParam    bool
	ReturnError bool
	ParamType   reflect.Type
	FuncValue   reflect.Value
}

func NewConstructorInfo(constructor interface{}, implement reflect.Type) (*ConstructorInfo, error) {
	rt := reflect.TypeOf(constructor)
	if rt.Kind() != reflect.Func {
		return nil, errors.New("constructor should be a function type")
	}
	if rt.NumIn() > 1 {
		return nil, errors.New("constructor parameters number should not greater than 1")
	}
	if rt.NumOut() > 2 {
		return nil, errors.New("constructor return number should not greater than 2")
	}
	if rt.NumOut() == 0 || !rt.Out(0).Implements(implement) {
		return nil, errors.New("constructor should return a RateLimiterGroup")
	}
	if rt.NumOut() == 2 && rt.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return nil, errors.New("constructor should return an error")
	}

	info := &ConstructorInfo{
		HasParam:    rt.NumIn() == 1,
		ReturnError: rt.NumOut() == 2,
		FuncValue:   reflect.ValueOf(constructor),
	}

	if rt.NumIn() > 0 {
		info.ParamType = rt.In(0)
	}

	return info, nil
}
