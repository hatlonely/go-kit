package logger

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

func RegisterWriter(key string, constructor interface{}) {
	if _, ok := writerConstructorMap[key]; ok {
		panic(fmt.Sprintf("writer type [%v] is already registered", key))
	}

	info, err := refx.NewConstructor(constructor, reflect.TypeOf((*Writer)(nil)).Elem())
	refx.Must(err)

	writerConstructorMap[key] = info
}

var writerConstructorMap = map[string]*refx.Constructor{}

type Writer interface {
	Write(info *Info) error
	Close() error
}

type WriterOptions struct {
	Type    string
	Options interface{}
}

func NewWriterWithOptions(options *WriterOptions, opts ...refx.Option) (Writer, error) {
	constructor, ok := writerConstructorMap[options.Type]
	if !ok {
		return nil, errors.Errorf("unregistered writer type: [%v]", options.Type)
	}

	result, err := constructor.Call(options.Options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "constructor.Call failed")
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, errors.Wrapf(result[1].Interface().(error), "NewWriter failed. type: [%v]", options.Type)
		}
		return result[0].Interface().(Writer), nil
	}

	return result[0].Interface().(Writer), nil
}

func NewWriterWithConfig(cfg *config.Config, opts ...refx.Option) (Writer, error) {
	var options WriterOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "cfg.Unmarshal failed.")
	}
	return NewWriterWithOptions(&options, opts...)
}
