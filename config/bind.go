package config

import (
	"reflect"
	"sync/atomic"
)

type BindOptions struct {
	OnSucc func(c *Config)
	OnFail func(err error)
}

type BindOption func(options *BindOptions)

func OnSucc(fun func(c *Config)) BindOption {
	return func(options *BindOptions) {
		options.OnSucc = fun
	}
}

func OnFail(fun func(err error)) BindOption {
	return func(options *BindOptions) {
		options.OnFail = fun
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
		if err := c.Sub(key).Unmarshal(val.Interface()); err == nil {
			av.Store(val.Elem().Interface())
		} else {
			c.log.Warnf("bind var failed. key: [%v], err: [%v]", key, err)
			if options.OnFail != nil {
				options.OnFail(err)
			}
		}
	}
	c.AddOnItemChangeHandler(key, func(conf *Config) {
		val := reflect.New(reflect.TypeOf(v))
		if err := conf.Sub(key).Unmarshal(val.Interface()); err != nil {
			c.log.Warnf("bind var failed. key: [%v], err: [%v]", key, err)
			if options.OnFail != nil {
				options.OnFail(err)
			}
			return
		}
		av.Store(val.Elem().Interface())
		if options.OnSucc != nil {
			options.OnSucc(c.Sub(""))
		}
	})
}
