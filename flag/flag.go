package flag

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/cast"
	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/strx"
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
		keyFlagInfosMap: map[string]*Info{},
		shorthandKeyMap: map[string]string{},
		nameKeyMap:      map[string]string{},
	}
}

type Flag struct {
	name            string
	arguments       []string
	keyFlagInfosMap map[string]*Info
	shorthandKeyMap map[string]string
	nameKeyMap      map[string]string
	kvs             interface{}
}

func (f *Flag) GetInfo(key string) (*Info, bool) {
	if info, ok := f.keyFlagInfosMap[key]; ok {
		return info, true
	}
	if k, ok := f.shorthandKeyMap[key]; ok {
		return f.keyFlagInfosMap[k], true
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

func (f *Flag) set(key string, val string) error {
	key = f.findKey(key)
	if err := refx.InterfaceSet(&f.kvs, key, val); err != nil {
		return errors.WithMessage(err, "InterfaceSet failed")
	}
	f.keyFlagInfosMap[key].Assigned = true
	if fun := f.keyFlagInfosMap[key].OnParse; fun != nil {
		return fun(val)
	}
	return nil
}

func (f *Flag) Set(key string, val interface{}) error {
	key = strx.KebabName(key)
	return f.set(key, cast.ToString(val))
}

func (f *Flag) Get(key string) (interface{}, bool) {
	key = strx.KebabName(key)
	key = strings.Replace(key, ".", "-", -1)
	key = strings.Replace(key, "--", "-", -1)

	if val, err := refx.InterfaceGet(f.kvs, key); err != nil {
		return nil, false
	} else {
		return val, true
	}
}
