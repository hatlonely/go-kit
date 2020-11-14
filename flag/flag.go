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
	Name         string
	Required     bool
	Assigned     bool
	Usage        string
	Shorthand    string
	DefaultValue string
	IsArgument   bool
	OnParse      func(val string) error
}

func NewFlag(name string) *Flag {
	return &Flag{
		name:      name,
		flagInfos: map[string]*Info{},
		shorthand: map[string]string{},
		kvs:       map[string]string{},
	}
}

type Flag struct {
	name      string
	flagInfos map[string]*Info
	arguments []string
	shorthand map[string]string
	kvs       map[string]string
}

func (f *Flag) GetInfo(key string) (*Info, bool) {
	if info, ok := f.flagInfos[key]; ok {
		return info, true
	}
	if k, ok := f.shorthand[key]; ok {
		return f.flagInfos[k], true
	}
	return nil, false
}

func (f *Flag) set(key string, val string) error {
	if k, ok := f.shorthand[key]; ok {
		f.kvs[k] = val
		f.flagInfos[k].Assigned = true
		if fun := f.flagInfos[k].OnParse; fun != nil {
			return fun(val)
		}
	} else {
		f.kvs[key] = val
		f.flagInfos[key].Assigned = true
		if fun := f.flagInfos[key].OnParse; fun != nil {
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
