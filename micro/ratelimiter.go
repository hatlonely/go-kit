package micro

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

// https://konghq.com/blog/how-to-design-a-scalable-rate-limiting-algorithm/

var ErrFlowControl = errors.New("flow control")

type RateLimiter interface {
	Wait(ctx context.Context, key string) error
	WaitN(ctx context.Context, key string, n int) (err error)
	Allow(key string) error
}

func RegisterRateLimiter(key string, constructor interface{}) {
	if _, ok := rateLimiterConstructorMap[key]; ok {
		panic(fmt.Sprintf("ratelimiter type [%v] is already registered", key))
	}

	info, err := refx.NewConstructor(constructor, reflect.TypeOf((*RateLimiter)(nil)).Elem())
	refx.Must(err)

	rateLimiterConstructorMap[key] = info
}

var rateLimiterConstructorMap = map[string]*refx.Constructor{}

func NewRateLimiterGroupWithOptions(options *RateLimiterOptions, opts ...refx.Option) (RateLimiter, error) {
	if options.Type == "" {
		return nil, nil
	}

	constructor, ok := rateLimiterConstructorMap[options.Type]
	if !ok {
		return nil, errors.Errorf("unregistered RateLimiter type: [%v]", options.Type)
	}

	result, err := constructor.Call(options.Options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "constructor.Call failed")
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, errors.Wrapf(result[1].Interface().(error), "NewRateLimiterGroupWithOptions failed. type: [%v]", options.Type)
		}
		return result[0].Interface().(RateLimiter), nil
	}

	return result[0].Interface().(RateLimiter), nil
}

type RateLimiterOptions struct {
	Type    string
	Options interface{}
}
