package flag

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/cast"
)

type AddFlagOptions struct {
	shorthand    string
	rtype        reflect.Type
	required     bool
	defaultValue string
}

type AddFlagOption func(*AddFlagOptions)

func Required() AddFlagOption {
	return func(o *AddFlagOptions) {
		o.required = true
	}
}

func DefaultValue(val string) AddFlagOption {
	return func(o *AddFlagOptions) {
		o.defaultValue = val
	}
}

func Shorthand(shorthand string) AddFlagOption {
	return func(o *AddFlagOptions) {
		o.shorthand = shorthand
	}
}

func Type(rtype reflect.Type) AddFlagOption {
	return func(o *AddFlagOptions) {
		o.rtype = rtype
	}
}

func (f *Flag) AddFlag(name string, usage string, opts ...AddFlagOption) error {
	o := &AddFlagOptions{}
	for _, opt := range opts {
		opt(o)
	}
	return f.addFlag(name, usage, o.rtype, o.required, o.shorthand, o.defaultValue, false)
}

func (f *Flag) AddArgument(name string, usage string, opts ...AddFlagOption) error {
	o := &AddFlagOptions{}
	for _, opt := range opts {
		opt(o)
	}
	return f.addFlag(name, usage, o.rtype, o.required, o.shorthand, o.defaultValue, true)
}

func (f *Flag) BindFlag(v interface{}, name string, usage string, rtype reflect.Type, required bool, shorthand string, defaultValue string, isArgument bool) error {
	if err := f.addFlag(name, usage, rtype, required, shorthand, defaultValue, isArgument); err != nil {
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

func (f *Flag) addFlag(name string, usage string, rtype reflect.Type, required bool, shorthand string, defaultValue string, isArgument bool) error {
	if _, ok := f.flagInfos[name]; ok {
		return errors.Errorf("conflict flag [%v]", name)
	}
	if shorthand != "" {
		if _, ok := f.shorthand[shorthand]; ok {
			return errors.Errorf("conflict shorthand [%v]", shorthand)
		}
	}

	if defaultValue != "" {
		required = false
	}
	info := &Info{
		Type:         rtype,
		Name:         name,
		Required:     required,
		Assigned:     false,
		Shorthand:    shorthand,
		Usage:        usage,
		DefaultValue: defaultValue,
		IsArgument:   isArgument,
	}

	if defaultValue != "" {
		f.kvs[name] = defaultValue
	}

	f.flagInfos[name] = info
	if isArgument {
		f.arguments = append(f.arguments, name)
	} else {
		if shorthand != "" {
			f.shorthand[shorthand] = name
		}
	}

	return nil
}
