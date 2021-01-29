package config

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

func RegisterDecoder(key string, constructor interface{}) {
	if _, ok := decoderConstructorMap[key]; ok {
		panic(fmt.Sprintf("decoder type [%v] is already registered", key))
	}

	info, err := refx.NewConstructor(constructor, reflect.TypeOf((*Decoder)(nil)).Elem())
	refx.Must(err)

	decoderConstructorMap[key] = info
}

var decoderConstructorMap = map[string]*refx.Constructor{}

type Decoder interface {
	Decode(buf []byte) (*Storage, error)
	Encode(storage *Storage) ([]byte, error)
}

func NewDecoderWithOptions(options *DecoderOptions, opts ...refx.Option) (Decoder, error) {
	constructor, ok := decoderConstructorMap[options.Type]
	if !ok {
		return nil, errors.Errorf("unsupported decoder type: [%v]", options.Type)
	}

	result, err := constructor.Call(options.Options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "constructor.Call failed")
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, errors.Wrapf(result[1].Interface().(error), "NewDecoder failed. type: [%v]", options.Type)
		}
		return result[0].Interface().(Decoder), nil
	}

	return result[0].Interface().(Decoder), nil
}

func NewDecoderWithConfig(cfg *Config, opts ...refx.Option) (Decoder, error) {
	var options DecoderOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.WithMessage(err, "cfg.Unmarshal failed.")
	}
	return NewDecoderWithOptions(&options)
}

func NewDecoder(typ string) (Decoder, error) {
	return NewDecoderWithOptions(&DecoderOptions{Type: typ})
}

type DecoderOptions struct {
	Type    string
	Options interface{}
}
