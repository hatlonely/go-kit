package flag

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/cast"
	"github.com/hatlonely/go-kit/refx"
)

type AddFlagOptions struct {
	Type         reflect.Type
	Usage        string
	Name         string
	Key          string
	Shorthand    string
	DefaultValue string
	IsArgument   bool
	Required     bool
}

type AddFlagOption func(*AddFlagOptions)

func Required() AddFlagOption {
	return func(options *AddFlagOptions) {
		options.Required = true
	}
}

func DefaultValue(val string) AddFlagOption {
	return func(options *AddFlagOptions) {
		options.DefaultValue = val
	}
}

func Shorthand(shorthand string) AddFlagOption {
	return func(options *AddFlagOptions) {
		options.Shorthand = shorthand
	}
}

func Type(v interface{}) AddFlagOption {
	return func(options *AddFlagOptions) {
		options.Type = reflect.TypeOf(v)
	}
}

func Key(key string) AddFlagOption {
	return func(options *AddFlagOptions) {
		options.Key = key
	}
}

func (f *Flag) AddOption(name string, usage string, opts ...AddFlagOption) {
	options := &AddFlagOptions{
		Name:       name,
		Usage:      usage,
		IsArgument: false,
	}
	for _, opt := range opts {
		opt(options)
	}
	if err := f.addFlagWithOptions(options); err != nil {
		panic(err)
	}
}

func (f *Flag) AddArgument(name string, usage string, opts ...AddFlagOption) {
	options := &AddFlagOptions{
		Name:       name,
		Usage:      usage,
		IsArgument: true,
	}
	for _, opt := range opts {
		opt(options)
	}
	if err := f.addFlagWithOptions(options); err != nil {
		panic(err)
	}
}

func (f *Flag) BindFlagWithOptions(v interface{}, options *AddFlagOptions) error {
	if err := f.addFlagWithOptions(options); err != nil {
		return err
	}
	info, _ := f.GetInfo(options.Key)
	if options.DefaultValue != "" {
		if err := cast.SetInterface(v, options.DefaultValue); err != nil {
			return err
		}
	}
	info.OnParse = func(val string) error {
		if val != "" {
			return cast.SetInterface(v, val)
		}
		return nil
	}
	return nil
}

func (f *Flag) addFlagWithOptions(options *AddFlagOptions) error {
	if _, ok := f.keyInfoMap[options.Name]; ok {
		return errors.Errorf("conflict flag [%v]", options.Name)
	}
	if options.Shorthand != "" {
		if _, ok := f.shorthandKeyMap[options.Shorthand]; ok {
			return errors.Errorf("conflict shorthand [%v]", options.Shorthand)
		}
	}
	if options.DefaultValue != "" {
		options.Required = false
	}
	if options.Key == "" {
		options.Key = options.Name
	}
	if options.Type == nil {
		options.Type = reflect.TypeOf("")
	}
	info := &Info{
		Type:         options.Type,
		Usage:        options.Usage,
		Name:         options.Name,
		Key:          options.Key,
		Shorthand:    options.Shorthand,
		DefaultValue: options.DefaultValue,
		IsArgument:   options.IsArgument,
		Required:     options.Required,
		Assigned:     false,
	}

	if options.DefaultValue != "" {
		if err := refx.InterfaceSet(&f.root, options.Key, options.DefaultValue); err != nil {
			return errors.WithMessage(err, "InterfaceSet failed")
		}
	}

	f.keyInfoMap[options.Key] = info

	if options.Shorthand != "" {
		f.shorthandKeyMap[options.Shorthand] = options.Key
	}
	f.nameKeyMap[options.Name] = options.Key
	if options.IsArgument {
		f.arguments = append(f.arguments, options.Name)
	} else {
		f.options = append(f.options, options.Name)
	}

	return nil
}
