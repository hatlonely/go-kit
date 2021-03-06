package logger

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

func RegisterFormatter(key string, constructor interface{}) {
	if _, ok := formatterConstructorMap[key]; ok {
		panic(fmt.Sprintf("formatter type [%v] is already registered", key))
	}

	info, err := refx.NewConstructor(constructor, reflect.TypeOf((*Formatter)(nil)).Elem())
	refx.Must(err)

	formatterConstructorMap[key] = info
}

var formatterConstructorMap = map[string]*refx.Constructor{}

type Formatter interface {
	Format(info *Info) ([]byte, error)
}

type FormatterOptions struct {
	Type    string
	Options interface{}
}

func NewFormatterWithOptions(options *FormatterOptions, opts ...refx.Option) (Formatter, error) {
	constructor, ok := formatterConstructorMap[options.Type]
	if !ok {
		return nil, errors.Errorf("unregistered formatter type: [%v]", options.Type)
	}

	result, err := constructor.Call(options.Options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "constructor.Call failed")
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, errors.Wrapf(result[1].Interface().(error), "NewFormatter failed. type: [%v]", options.Type)
		}
		return result[0].Interface().(Formatter), nil
	}

	return result[0].Interface().(Formatter), nil
}
