package wrap

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

func RegisterRateLimiterGroup(key string, constructor interface{}) {
	if _, ok := rateLimiterGroupMap[key]; ok {
		panic(fmt.Sprintf("ratelimiter type [%v] is already registered", key))
	}

	rt := reflect.TypeOf(constructor)
	if rt.Kind() != reflect.Func {
		panic("constructor should be a function type")
	}
	if rt.NumIn() > 1 {
		panic("constructor parameters number should not greater than 1")
	}
	if rt.NumOut() > 2 {
		panic("constructor return number should not greater than 2")
	}
	if rt.NumOut() == 0 || !rt.Out(0).Implements(reflect.TypeOf((*RateLimiterGroup)(nil)).Elem()) {
		panic("constructor should return a RateLimiterGroup")
	}
	if rt.NumOut() == 2 && rt.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		panic("constructor should return an error")
	}

	info := &RateLimiterGroupConstructor{
		HasParam:    rt.NumIn() == 1,
		ReturnError: rt.NumOut() == 2,
		FuncValue:   reflect.ValueOf(constructor),
	}

	if rt.NumIn() > 0 {
		info.ParamType = rt.In(0)
	}

	rateLimiterGroupMap[key] = info
}

type RateLimiterGroupConstructor struct {
	HasParam    bool
	ReturnError bool
	ParamType   reflect.Type
	FuncValue   reflect.Value
}

var rateLimiterGroupMap = map[string]*RateLimiterGroupConstructor{}

type RateLimiterGroup interface {
	Wait(ctx context.Context, key string) error
}

func NewRateLimiterGroupWithOptions(options *RateLimiterGroupOptions, opts ...refx.Option) (RateLimiterGroup, error) {
	constructor, ok := rateLimiterGroupMap[options.Type]
	if !ok {
		return nil, errors.Errorf("unsupported rate limiter type: [%v]", options.Type)
	}

	var result []reflect.Value
	if constructor.HasParam {
		if reflect.TypeOf(options.Options) == constructor.ParamType {
			result = constructor.FuncValue.Call([]reflect.Value{reflect.ValueOf(options.Options)})
		} else {
			params := reflect.New(constructor.ParamType)
			if err := refx.InterfaceToStruct(options.Options, params.Interface(), opts...); err != nil {
				return nil, errors.Wrap(err, "refx.InterfaceToStruct failed")
			}
			result = constructor.FuncValue.Call([]reflect.Value{params.Elem()})
		}
	} else {
		result = constructor.FuncValue.Call(nil)
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, result[1].Interface().(error)
		}
		return result[0].Interface().(RateLimiterGroup), nil
	}

	return result[0].Interface().(RateLimiterGroup), nil
}

type RateLimiterGroupOptions struct {
	Type    string
	Options interface{}
}
