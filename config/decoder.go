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

	info, err := refx.NewConstructorInfo(constructor, reflect.TypeOf((*Decoder)(nil)).Elem())
	refx.Must(err)

	decoderConstructorMap[key] = info
}

var decoderConstructorMap = map[string]*refx.ConstructorInfo{}

type Decoder interface {
	Decode(buf []byte) (*Storage, error)
	Encode(storage *Storage) ([]byte, error)
}

func NewDecoderWithOptions(options *DecoderOptions, opts ...refx.Option) (Decoder, error) {
	constructor, ok := decoderConstructorMap[options.Type]
	if !ok {
		return nil, errors.Errorf("unsupported decoder type: [%v]", options.Type)
	}

	var result []reflect.Value
	if constructor.HasParam {
		if reflect.TypeOf(options.Options) == constructor.ParamType {
			result = constructor.FuncValue.Call([]reflect.Value{reflect.ValueOf(options.Options)})
		} else {
			params := reflect.New(constructor.ParamType)
			if err := refx.InterfaceToStruct(options.Options, params.Interface(), opts...); err != nil {
				return nil, errors.Wrap(err, "refx.InterfaceToStruct failed")
			}
			result = constructor.FuncValue.Call([]reflect.Value{params.Elem()})
		}
	} else {
		result = constructor.FuncValue.Call(nil)
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, result[1].Interface().(error)
		}
		return result[0].Interface().(Decoder), nil
	}

	return result[0].Interface().(Decoder), nil
}

func NewDecoderWithConfig(cfg *Config, opts ...refx.Option) (Decoder, error) {
	var options DecoderOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "cfg.Unmarshal failed.")
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
