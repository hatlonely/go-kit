#!/usr/bin/env python3

import os
import re

get_header = """// this file is generate by autogen.py, do not edit
package config

import (
	"fmt"
	"net"
	"reflect"
	"time"
)
"""

cast_tpl = """
func To{name}E(v interface{{}}) ({type}, error) {{
	return cast.To{name}E(v)
}}
"""

gete_tpl = """
func (c *Config) Get{name}E(key string) ({type}, error) {{
	v, err := c.Get(key)
	if err != nil {{
		var res {type}
		return res, err
	}}
	return To{name}E(v)
}}
"""

get_tpl = """
func (c *Config) Get{name}(key string) {type} {{
	v, _ := c.Get{name}E(key)
	return v
}}
"""

getp_tpl = """
func (c *Config) Get{name}P(key string) {type} {{
	v, err := c.Get{name}E(key)
	if err != nil {{
		panic(err)
	}}
	return v
}}
"""

getd_tpl = """
func (c *Config) Get{name}D(key string, dftVal {type}) {type} {{
	v, err := c.Get{name}E(key)
	if err != nil {{
		return dftVal
	}}
	return v
}}
"""

set_interface_tpl = """
func SetInterface(dst interface{{}}, src interface{{}}) error {{
	switch dst.(type) {{
{items}
	default:
		return fmt.Errorf("unsupport dst type [%v]", reflect.TypeOf(dst))
	}}

	return nil
}}
"""

item_tpl = """	case *{type}:
		v, err := To{name}E(src)
		if err != nil {{
			return err
		}}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
"""

bind_header = """// this file is generate by autogen.py, do not edit
package config

import (
	"net"
	"sync/atomic"
	"time"
)
"""

atomic_type_tpl = """
type Atomic{name} struct {{
	v atomic.Value
}}

func NewAtomic{name}(v {type}) *Atomic{name} {{
	var av atomic.Value
	av.Store(v)
	return &Atomic{name}{{v: av}}
}}

func (a *Atomic{name}) Get() {type} {{
	return a.v.Load().({type})
}}

func (a *Atomic{name}) Set(v {type}) {{
	a.v.Store(v)
}}
"""

bind_var_tpl = """
func (c *Config) {name}Var(key string, av *Atomic{name}, opts ...BindOption) {{
	options := &BindOptions{{}}
	for _, opt := range opts {{
		opt(options)
	}}

	var v {type}
	if c.storage != nil {{
		v = c.Get{name}(key)
	}}
	av.Set(v)
	c.AddOnItemChangeHandler(key, func(conf *Config) {{
		var err error
		v, err = c.Get{name}E(key)
		if err != nil {{
			if options.OnFail != nil {{
				options.OnFail(err)
			}}
			return
		}}
		av.Set(v)
		if options.OnSucc != nil {{
			options.OnSucc(c.Sub(""))
		}}
	}})
}}
"""

bind_tpl = """
func (c *Config) {name}(key string, opts ...BindOption) *Atomic{name} {{
	var v Atomic{name}
	c.{name}Var(key, &v, opts...)
	return &v
}}
"""

global_header = """// this file is generate by autogen.py, do not edit
package config

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

var gconf = &Config{
	itemHandlers: map[string][]OnChangeHandler{},
	log:          logrus.New(),
}

func Init(filename string) error {
	conf, err := NewConfigWithBaseFile(filename)
	if err != nil {
		return err
	}
	gconf.storage = conf.storage
	gconf.decoder = conf.decoder
	gconf.provider = conf.provider
	return nil
}
"""

global_tpl1 = """
func {define} {{
	return gconf.{name}({params})
}}
"""

global_tpl2 = """
func {define} {{
	gconf.{name}({params})
}}
"""


def generate_global_go():
    global_out = open("autogen_global.go", "w")
    funs = []
    for name in os.listdir("."):
        if not name.endswith(".go"):
            continue
        for line in open(name):
            res = re.match(r"func \(c \*Config\) ((\w+)\((.*?)\)(.*?)) {", line)
            if not res:
                continue
            params = [i.strip().split(" ")[0] for i in res.group(3).split(",")]
            funs.append({
                "define": res.group(1),
                "name": res.group(2),
                "params": ", ".join([i if i != "opts" else "opts..." for i in params]),
                "return": res.group(4)
            })
    global_out.write(global_header)
    for fun in funs:
        if fun["return"]:
            global_out.write(global_tpl1.format(**fun))
        else:
            global_out.write(global_tpl2.format(**fun))
    global_out.close()


def generate_get_bind():
    infos = []
    for line in open("cast.go"):
        res = re.match(r"func To(.*?)E\(v interface{}\) \((.*?), error\).*", line)
        if not res:
            continue
        infos.append({"name": res.group(1), "type": res.group(2)})

    get_out = open("autogen_get.go", "w")
    get_out.write(get_header)
    for info in infos:
        get_out.write(get_tpl.format(**info))
    for info in infos:
        get_out.write(gete_tpl.format(**info))
    for info in infos:
        get_out.write(getp_tpl.format(**info))
    for info in infos:
        get_out.write(getd_tpl.format(**info))
    items = ""
    for info in infos:
        items += item_tpl.format(**info)
    get_out.write(set_interface_tpl.format(items=items))
    get_out.close()

    bind_out = open("autogen_bind.go", "w")
    bind_out.write(bind_header)
    for info in infos:
        bind_out.write(atomic_type_tpl.format(**info))
    for info in infos:
        bind_out.write(bind_var_tpl.format(**info))
    for info in infos:
        bind_out.write(bind_tpl.format(**info))
    bind_out.close()


def main():
    generate_get_bind()
    generate_global_go()


if __name__ == "__main__":
    main()
