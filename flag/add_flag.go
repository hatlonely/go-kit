package flag

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/cast"
)

type AddFlagOptions struct {
	Shorthand    string
	Type         reflect.Type
	Required     bool
	DefaultValue string
	Usage        string
	Name         string
	IsArgument   bool
}

type AddFlagOption func(*AddFlagOptions)

func Required() AddFlagOption {
	return func(o *AddFlagOptions) {
		o.Required = true
	}
}

func DefaultValue(val string) AddFlagOption {
	return func(o *AddFlagOptions) {
		o.DefaultValue = val
	}
}

func Shorthand(shorthand string) AddFlagOption {
	return func(o *AddFlagOptions) {
		o.Shorthand = shorthand
	}
}

func Type(v interface{}) AddFlagOption {
	return func(o *AddFlagOptions) {
		o.Type = reflect.TypeOf(v)
	}
}

func (f *Flag) AddFlag(name string, usage string, opts ...AddFlagOption) error {
	options := &AddFlagOptions{
		Name:       name,
		Usage:      usage,
		IsArgument: false,
	}
	for _, opt := range opts {
		opt(options)
	}
	return f.addFlagWithOptions(options)
}

func (f *Flag) AddArgument(name string, usage string, opts ...AddFlagOption) error {
	options := &AddFlagOptions{
		Name:       name,
		Usage:      usage,
		IsArgument: true,
	}
	for _, opt := range opts {
		opt(options)
	}
	return f.addFlagWithOptions(options)
}

func (f *Flag) BindFlag(v interface{}, name string, usage string, rtype reflect.Type, required bool, shorthand string, defaultValue string, isArgument bool) error {
	if err := f.addFlagWithOptions(&AddFlagOptions{
		Name:         name,
		Usage:        usage,
		Type:         rtype,
		Required:     required,
		Shorthand:    shorthand,
		DefaultValue: defaultValue,
		IsArgument:   isArgument,
	}); err != nil {
		return err
	}
	info, _ := f.GetInfo(name)
	if defaultValue != "" {
		if err := cast.SetInterface(v, defaultValue); err != nil {
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
	if _, ok := f.flagInfos[options.Name]; ok {
		return errors.Errorf("conflict flag [%v]", options.Name)
	}
	if options.Shorthand != "" {
		if _, ok := f.shorthand[options.Shorthand]; ok {
			return errors.Errorf("conflict shorthand [%v]", options.Shorthand)
		}
	}

	if options.DefaultValue != "" {
		options.Required = false
	}
	info := &Info{
		Type:         options.Type,
		Name:         options.Name,
		Required:     options.Required,
		Assigned:     false,
		Shorthand:    options.Shorthand,
		Usage:        options.Usage,
		DefaultValue: options.DefaultValue,
		IsArgument:   options.IsArgument,
	}

	if options.DefaultValue != "" {
		f.kvs[options.Name] = options.DefaultValue
	}

	f.flagInfos[options.Name] = info
	if options.IsArgument {
		f.arguments = append(f.arguments, options.Name)
	} else {
		if options.Shorthand != "" {
			f.shorthand[options.Shorthand] = options.Name
		}
	}

	return nil
}

func (f *Flag) addFlag(name string, usage string, rtype reflect.Type, required bool, shorthand string, defaultValue string, isArgument bool) error {
	return f.addFlagWithOptions(&AddFlagOptions{
		Name:         name,
		Usage:        usage,
		Type:         rtype,
		Required:     required,
		Shorthand:    shorthand,
		DefaultValue: defaultValue,
		IsArgument:   isArgument,
	})
}
