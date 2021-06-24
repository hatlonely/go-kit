package rpcx

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

type MiddlewareOptions struct {
	Type    string
	Options string
}

func RegisterLocker(key string, constructor interface{}) {
	if _, ok := middlewareConstructorMap[key]; ok {
		panic(fmt.Sprintf("locker type [%v] is already registered", key))
	}

	info, err := refx.NewConstructor(constructor, reflect.TypeOf((*Middleware)(nil)).Elem())
	refx.Must(err)

	middlewareConstructorMap[key] = info
}

var middlewareConstructorMap = map[string]*refx.Constructor{}

func NewMiddlewareWithOptions(options *MiddlewareOptions, opts ...refx.Option) (Middleware, error) {
	if options.Type == "" {
		return nil, nil
	}

	constructor, ok := middlewareConstructorMap[options.Type]
	if !ok {
		return nil, errors.Errorf("unregistered Middleware type: [%v]", options.Type)
	}

	result, err := constructor.Call(options.Options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "constructor.Call failed")
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, errors.Wrapf(result[1].Interface().(error), "NewMiddlewareWithOptions failed. type: [%v]", options.Type)
		}
		return result[0].Interface().(Middleware), nil
	}

	return result[0].Interface().(Middleware), nil
}

type Middleware interface {
	HttpMiddleware
	GrpcMiddleware
}
