package config

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

func RegisterProvider(key string, constructor interface{}) {
	if _, ok := providerConstructorMap[key]; ok {
		panic(fmt.Sprintf("provider type [%v] is already registered", key))
	}

	info, err := refx.NewConstructor(constructor, reflect.TypeOf((*Provider)(nil)).Elem())
	refx.Must(err)

	providerConstructorMap[key] = info
}

var providerConstructorMap = map[string]*refx.Constructor{}

type Provider interface {
	Events() <-chan struct{}
	Errors() <-chan error
	Load() ([]byte, error)
	Dump(buf []byte) error
	EventLoop(ctx context.Context) error
}

func NewProviderWithConfig(cfg *Config, opts ...refx.Option) (Provider, error) {
	var options ProviderOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.WithMessage(err, "cfg.Unmarshal failed.")
	}
	return NewProviderWithOptions(&options)
}

func NewProviderWithOptions(options *ProviderOptions, opts ...refx.Option) (Provider, error) {
	constructor, ok := providerConstructorMap[options.Type]
	if !ok {
		return nil, errors.Errorf("unregistered provider type: [%v]", options.Type)
	}

	result, err := constructor.Call(options.Options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "constructor.Call failed")
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, errors.Wrapf(result[1].Interface().(error), "NewProvider failed. type: [%v]", options.Type)
		}
		return result[0].Interface().(Provider), nil
	}

	return result[0].Interface().(Provider), nil
}

type ProviderOptions struct {
	Type    string
	Options interface{}
}
