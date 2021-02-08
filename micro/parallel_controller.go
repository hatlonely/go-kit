package micro

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

var ErrParallelControl = errors.New("parallel control")
var ErrContextCancel = errors.New("context cancel")

type ParallelController interface {
	// 尝试获取 token，如果已空，返回 ErrParallelControl
	// 1. 成功返回 nil
	// 2. 已空，返回 ErrParallelControl
	// 3. 其他未知错误返回 err
	TryAcquire(ctx context.Context, key string) (int, error)
	// 获取 token，如果已空，阻塞等待
	// 1. 成功返回 nil
	// 2. 未知错误返回 err
	Acquire(ctx context.Context, key string) (int, error)
	// 返还 token，如果已满，直接丢弃
	// 1. 成功返回 nil
	// 2. 未知错误返回 err
	Release(ctx context.Context, key string, token int) error
}

func RegisterParallelController(key string, constructor interface{}) {
	if _, ok := parallelControllerGroupConstructorMap[key]; ok {
		panic(fmt.Sprintf("ratelimiter type [%v] is already registered", key))
	}

	info, err := refx.NewConstructor(constructor, reflect.TypeOf((*ParallelController)(nil)).Elem())
	refx.Must(err)

	parallelControllerGroupConstructorMap[key] = info
}

var parallelControllerGroupConstructorMap = map[string]*refx.Constructor{}

func NewParallelControllerWithOptions(options *ParallelControllerOptions, opts ...refx.Option) (ParallelController, error) {
	if options.Type == "" {
		return nil, nil
	}

	constructor, ok := parallelControllerGroupConstructorMap[options.Type]
	if !ok {
		return nil, errors.Errorf("unregistered ParallelController type: [%v]", options.Type)
	}

	result, err := constructor.Call(options.Options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "constructor.Call failed")
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, errors.Wrapf(result[1].Interface().(error), "NewParallelControllerWithOptions failed. type: [%v]", options.Type)
		}
		return result[0].Interface().(ParallelController), nil
	}

	return result[0].Interface().(ParallelController), nil
}

type ParallelControllerOptions struct {
	Type    string
	Options interface{}
}
