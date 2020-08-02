package kv

import (
	"encoding/json"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type KV struct {
	store   Store
	options *Options
}

var ErrNotFound = errors.New("NotFound")

type Options struct {
	LoadFunc     func(key interface{}) (interface{}, error)
	KeyCompress  func(v interface{}) ([]byte, error)
	ValMarshal   func(v interface{}) ([]byte, error)
	ValUnmarshal func(buf []byte, v interface{}) error
}

var defaultKVOptions = Options{
	KeyCompress: func(v interface{}) ([]byte, error) {
		return json.Marshal(v)
	},
	ValMarshal: func(v interface{}) ([]byte, error) {
		return json.Marshal(v)
	},
	ValUnmarshal: func(buf []byte, v interface{}) error {
		return json.Unmarshal(buf, v)
	},
}

type Option func(options *Options)

func WithLoadFunc(loadFunc func(key interface{}) (interface{}, error)) Option {
	return func(options *Options) {
		options.LoadFunc = loadFunc
	}
}

func WithStringKey() Option {
	return func(options *Options) {
		options.KeyCompress = func(v interface{}) ([]byte, error) {
			return []byte(v.(string)), nil
		}
	}
}

func WithStringVal() Option {
	return func(options *Options) {
		options.ValMarshal = func(v interface{}) ([]byte, error) {
			switch v.(type) {
			case string:
				return []byte(v.(string)), nil
			case *string:
				val := v.(*string)
				return []byte(*val), nil
			}
			return nil, errors.New("unsupported val type")
		}
		options.ValUnmarshal = func(buf []byte, v interface{}) error {
			reflect.ValueOf(v).Elem().SetString(string(buf))
			return nil
		}
	}
}

func WithProtoVal() Option {
	return func(options *Options) {
		options.ValMarshal = func(v interface{}) ([]byte, error) {
			return proto.Marshal(v.(proto.Message))
		}
		options.ValUnmarshal = func(buf []byte, v interface{}) error {
			return proto.Unmarshal(buf, v.(proto.Message))
		}
	}
}

func WithKeyCompress(keyMarshal func(v interface{}) ([]byte, error)) Option {
	return func(options *Options) {
		options.KeyCompress = keyMarshal
	}
}

func WithValMarshal(valMarshal func(v interface{}) ([]byte, error)) Option {
	return func(options *Options) {
		options.ValMarshal = valMarshal
	}
}

func WithValUnmarshal(valUnmarshal func(buf []byte, v interface{}) error) Option {
	return func(options *Options) {
		options.ValUnmarshal = valUnmarshal
	}
}

func NewKV(store Store, opts ...Option) *KV {
	options := defaultKVOptions
	for _, opt := range opts {
		opt(&options)
	}

	return &KV{
		store:   store,
		options: &options,
	}
}

func (c *KV) Set(key interface{}, val interface{}) error {
	keyb, err := c.options.KeyCompress(key)
	if err != nil {
		return err
	}
	valb, err := c.options.ValMarshal(val)
	if err != nil {
		return err
	}
	if err := c.store.Set(keyb, valb); err != nil {
		return err
	}
	return nil
}

func (c *KV) Get(key interface{}, val interface{}) error {
	keyb, err := c.options.KeyCompress(key)
	if err != nil {
		return err
	}
	valb, err := c.store.Get(keyb)
	if err != nil {
		return err
	}
	if valb == nil {
		if c.options.LoadFunc == nil {
			return nil
		}

		v, err := c.options.LoadFunc(key)
		if err != nil {
			return err
		}
		if v == nil {
			if err := c.store.Set(keyb, []byte{}); err != nil {
				return err
			}
			return ErrNotFound
		}
		valb, err = c.options.ValMarshal(v)
		if err := c.store.Set(keyb, valb); err != nil {
			return err
		}
	}

	if len(valb) == 0 {
		return ErrNotFound
	}

	if err := c.options.ValUnmarshal(valb, val); err != nil {
		return err
	}

	return nil
}
