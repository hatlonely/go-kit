package micro

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

type ParallelControllerGroup interface {
	PutToken(ctx context.Context, key string) error
	GetToken(ctx context.Context, key string) error
}

func RegisterParallelControllerGroup(key string, constructor interface{}) {
	if _, ok := parallelControllerGroupConstructorMap[key]; ok {
		panic(fmt.Sprintf("ratelimiter type [%v] is already registered", key))
	}

	info, err := refx.NewConstructor(constructor, reflect.TypeOf((*ParallelControllerGroup)(nil)).Elem())
	refx.Must(err)

	parallelControllerGroupConstructorMap[key] = info
}

var parallelControllerGroupConstructorMap = map[string]*refx.Constructor{}

func NewParallelControllerGroupWithOptions(options *ParallelControllerGroupOptions, opts ...refx.Option) (ParallelControllerGroup, error) {
	if options.Type == "" {
		return nil, nil
	}

	constructor, ok := parallelControllerGroupConstructorMap[options.Type]
	if !ok {
		return nil, errors.Errorf("unregistered ParallelControllerGroup type: [%v]", options.Type)
	}

	result, err := constructor.Call(options.Options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "constructor.Call failed")
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, errors.Wrapf(result[1].Interface().(error), "NewParallelControllerGroupWithOptions failed. type: [%v]", options.Type)
		}
		return result[0].Interface().(ParallelControllerGroup), nil
	}

	return result[0].Interface().(ParallelControllerGroup), nil
}

type ParallelControllerGroupOptions struct {
	Type    string
	Options interface{}
}
