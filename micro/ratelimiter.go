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
	// 立即返回是否限流
	// 1. 成功返回 nil
	// 2. 限流返回 ErrFlowControl
	// 3. 未知错误返回 err
	Allow(ctx context.Context, key string) error
	// 阻塞等待限流结果返回，或者立即返回错误
	// 1. 成功返回 nil
	// 2. 未知错误返回 err，比如网络错误，不可能返回 ErrFlowControl
	Wait(ctx context.Context, key string) error
	// 阻塞等待 n 个限流结果，或者立即返回错误
	// 1. 成功返回 nil
	// 2. 未知错误返回 err，比如网络错误，不可能返回 ErrFlowControl
	WaitN(ctx context.Context, key string, n int) (err error)
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

func NewRateLimiterWithOptions(options *RateLimiterOptions, opts ...refx.Option) (RateLimiter, error) {
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
			return nil, errors.Wrapf(result[1].Interface().(error), "NewRateLimiterWithOptions failed. type: [%v]", options.Type)
		}
		return result[0].Interface().(RateLimiter), nil
	}

	return result[0].Interface().(RateLimiter), nil
}

type RateLimiterOptions struct {
	Type    string
	Options interface{}
}
