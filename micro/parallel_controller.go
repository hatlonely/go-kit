package micro

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

type ParallelController interface {
	PutToken(ctx context.Context, key string) error
	GetToken(ctx context.Context, key string) error
	TryGetToken(ctx context.Context, key string) bool
	TryPutToken(ctx context.Context, key string) bool
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
