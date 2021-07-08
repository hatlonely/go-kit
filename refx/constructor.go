package refx

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

func Must(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}

type Constructor struct {
	HasParam    bool
	HasOption   bool
	ReturnError bool
	IsInstance  bool
	Instance    reflect.Value
	ParamType   reflect.Type
	FuncValue   reflect.Value
}

func (c *Constructor) Call(v interface{}, opts ...Option) ([]reflect.Value, error) {
	if c.IsInstance {
		return []reflect.Value{c.Instance}, nil
	}

	var params []reflect.Value
	if c.HasParam {
		if reflect.TypeOf(v) == c.ParamType {
			params = append(params, reflect.ValueOf(v))
		} else {
			param := reflect.New(c.ParamType)
			if err := InterfaceToStruct(v, param.Interface(), opts...); err != nil {
				return nil, errors.WithMessage(err, "refx.InterfaceToStruct failed")
			}
			params = append(params, param.Elem())
		}
	}
	if c.HasOption {
		for _, opt := range opts {
			params = append(params, reflect.ValueOf(opt))
		}
	}
	return c.FuncValue.Call(params), nil
}

func NewConstructor(constructor interface{}, implement reflect.Type) (*Constructor, error) {
	rt := reflect.TypeOf(constructor)

	if rt.Implements(implement) {
		return &Constructor{
			IsInstance: true,
			Instance:   reflect.ValueOf(constructor),
		}, nil
	}

	if rt.Kind() != reflect.Func {
		return nil, errors.New("constructor should be a function type")
	}

	var info Constructor
	info.FuncValue = reflect.ValueOf(constructor)

	if rt.NumIn() > 2 {
		return nil, errors.New("constructor parameters number should not greater than 2")
	}
	if rt.NumIn() == 2 {
		if rt.In(1) != reflect.TypeOf([]Option{}) {
			panic("constructor parameters should be []refx.Option")
		}
		info.HasParam = true
		info.HasOption = true
	}
	if rt.NumIn() == 1 {
		if rt.In(0) != reflect.TypeOf([]Option{}) {
			info.HasParam = true
		} else {
			info.HasOption = true
		}
	}
	if rt.NumOut() > 2 {
		return nil, errors.New("constructor return number should not greater than 2")
	}
	if rt.NumOut() == 0 || !rt.Out(0).Implements(implement) {
		return nil, errors.Errorf("constructor should return a %s", implement)
	}
	if rt.NumOut() == 2 && rt.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return nil, errors.New("constructor should return an error")
	}

	info.ReturnError = rt.NumOut() == 2

	if info.HasParam {
		info.ParamType = rt.In(0)
	}

	return &info, nil
}

var constructorMap = map[reflect.Type]map[string]*Constructor{}

func Register(implement reflect.Type, key string, constructor interface{}) {
	if _, ok := constructorMap[implement]; !ok {
		constructorMap[implement] = map[string]*Constructor{}
	}

	if _, ok := constructorMap[implement][key]; ok {
		panic(fmt.Sprintf("%v type [%v] is already registered", implement.String(), key))
	}

	info, err := NewConstructor(constructor, implement)
	Must(err)

	constructorMap[implement][key] = info
}

func New(implement reflect.Type, options *TypeOptions, opts ...Option) (interface{}, error) {
	if options.Type == "" {
		return nil, nil
	}

	if _, ok := constructorMap[implement]; !ok {
		return nil, errors.Errorf("%v type not found", implement.String())
	}

	constructor, ok := constructorMap[implement][options.Type]
	if !ok {
		return nil, errors.Errorf("unregistered %v type: [%v]", implement.String(), options.Type)
	}

	result, err := constructor.Call(options.Options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "constructor.Call failed")
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, errors.Wrapf(result[1].Interface().(error), "New failed. type: [%v]", options.Type)
		}
		return result[0].Interface(), nil
	}

	return result[0].Interface(), nil
}

type TypeOptions struct {
	Type    string
	Options interface{}
}
