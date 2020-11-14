package flag

import (
	"reflect"
	"strings"

	"github.com/hatlonely/go-kit/cast"
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
		kvs:             map[string]string{},
	}
}

type Flag struct {
	name            string
	arguments       []string
	keyFlagInfosMap map[string]*Info
	shorthandKeyMap map[string]string
	nameKeyMap      map[string]string
	kvs             map[string]string
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

func (f *Flag) set(key string, val string) error {
	if k, ok := f.shorthandKeyMap[key]; ok {
		f.kvs[k] = val
		f.keyFlagInfosMap[k].Assigned = true
		if fun := f.keyFlagInfosMap[k].OnParse; fun != nil {
			return fun(val)
		}
	} else {
		f.kvs[key] = val
		f.keyFlagInfosMap[key].Assigned = true
		if fun := f.keyFlagInfosMap[key].OnParse; fun != nil {
			return fun(val)
		}
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
	if v, ok := f.kvs[key]; ok {
		return v, true
	}
	return nil, false
}
