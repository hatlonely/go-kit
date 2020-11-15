package flag

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

type Options struct {
	Help       bool   `flag:"--help,-h; usage: show help info"`
	Version    bool   `flag:"--version,-v; usage: show version"`
	ConfigPath string `flag:"--config-path,-c; usage: config path"`
}

type Info struct {
	Type         reflect.Type
	Usage        string
	Name         string
	Key          string
	Shorthand    string
	DefaultValue string
	IsArgument   bool
	Required     bool
	Assigned     bool
	OnParse      func(val string) error
}

func NewFlag(name string) *Flag {
	return &Flag{
		name:            name,
		keyInfoMap:      map[string]*Info{},
		shorthandKeyMap: map[string]string{},
		nameKeyMap:      map[string]string{},
	}
}

type Flag struct {
	name            string
	arguments       []string
	keyInfoMap      map[string]*Info
	shorthandKeyMap map[string]string
	nameKeyMap      map[string]string
	root            interface{}
}

func (f *Flag) GetInfo(key string) (*Info, bool) {
	key = f.findKey(key)
	if info, ok := f.keyInfoMap[key]; ok {
		return info, true
	}
	return nil, false
}

func (f *Flag) findKey(key string) string {
	if k, ok := f.shorthandKeyMap[key]; ok {
		return k
	}
	if k, ok := f.nameKeyMap[key]; ok {
		return k
	}
	return key
}

func (f *Flag) Set(key string, val string) error {
	key = f.findKey(key)
	if err := refx.InterfaceSet(&f.root, key, val); err != nil {
		return errors.WithMessage(err, "InterfaceSet failed")
	}
	if _, ok := f.keyInfoMap[key]; ok {
		f.keyInfoMap[key].Assigned = true
		if fun := f.keyInfoMap[key].OnParse; fun != nil {
			return fun(val)
		}
	}
	return nil
}

func (f *Flag) Get(key string) (interface{}, bool) {
	key = f.findKey(key)
	if val, err := refx.InterfaceGet(f.root, key); err != nil {
		return nil, false
	} else {
		return val, true
	}
}
