package config

import (
	"reflect"
	"sync/atomic"

	"github.com/hatlonely/go-kit/refx"
)

type BindOptions struct {
	OnSucc      func(cfg *Config)
	OnFail      func(err error)
	RefxOptions []refx.Option
}

type BindOption func(options *BindOptions)

func OnSucc(fun func(cfg *Config)) BindOption {
	return func(options *BindOptions) {
		options.OnSucc = fun
	}
}

func OnFail(fun func(err error)) BindOption {
	return func(options *BindOptions) {
		options.OnFail = fun
	}
}

func WithUnmarshalOptions(opts ...refx.Option) BindOption {
	return func(options *BindOptions) {
		options.RefxOptions = opts
	}
}

func (c *Config) Bind(key string, v interface{}, opts ...BindOption) *atomic.Value {
	var av atomic.Value
	c.BindVar(key, v, &av, opts...)
	return &av
}

func (c *Config) BindVar(key string, v interface{}, av *atomic.Value, opts ...BindOption) {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	val := reflect.New(reflect.TypeOf(v))
	if c.storage != nil {
		if err := c.Sub(key).Unmarshal(val.Interface(), options.RefxOptions...); err == nil {
			av.Store(val.Elem().Interface())
		} else {
			c.log.Warnf("bind var failed. key: [%v], err: [%v]", key, err)
			if options.OnFail != nil {
				options.OnFail(err)
			}
		}
	}
	c.AddOnItemChangeHandler(key, func(cfg *Config) error {
		val := reflect.New(reflect.TypeOf(v))
		if err := cfg.Unmarshal(val.Interface(), options.RefxOptions...); err != nil {
			c.log.Warnf("bind var failed. key: [%v], err: [%v]", key, err)
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return err
		}
		av.Store(val.Elem().Interface())
		if options.OnSucc != nil {
			options.OnSucc(cfg)
		}
		return nil
	})
}
