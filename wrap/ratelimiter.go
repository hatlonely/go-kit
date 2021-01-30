package wrap

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

func RegisterRateLimiterGroup(key string, constructor interface{}) {
	if _, ok := rateLimiterGroupConstructorMap[key]; ok {
		panic(fmt.Sprintf("ratelimiter type [%v] is already registered", key))
	}

	info, err := refx.NewConstructor(constructor, reflect.TypeOf((*RateLimiterGroup)(nil)).Elem())
	refx.Must(err)

	rateLimiterGroupConstructorMap[key] = info
}

var rateLimiterGroupConstructorMap = map[string]*refx.Constructor{}

type RateLimiterGroup interface {
	Wait(ctx context.Context, key string) error
}

func NewRateLimiterGroupWithOptions(options *RateLimiterGroupOptions, opts ...refx.Option) (RateLimiterGroup, error) {
	constructor, ok := rateLimiterGroupConstructorMap[options.Type]
	if !ok {
		return nil, errors.Errorf("unregistered RateLimiterGroup type: [%v]", options.Type)
	}

	result, err := constructor.Call(options.Options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "constructor.Call failed")
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, errors.Wrapf(result[1].Interface().(error), "NewRateLimiterGroupWithOptions failed. type: [%v]", options.Type)
		}
		return result[0].Interface().(RateLimiterGroup), nil
	}

	return result[0].Interface().(RateLimiterGroup), nil
}

type RateLimiterGroupOptions struct {
	Type    string
	Options interface{}
}
