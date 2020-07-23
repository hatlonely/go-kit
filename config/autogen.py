#!/usr/bin/env python3

import os
import re


def parse_info(filename):
    infos = []
    for line in open(filename):
        res = re.match(r"func To(.*?)E\(v interface{}\) \((.*?), error\).*", line)
        if not res:
            continue
        infos.append({"name": res.group(1), "type": res.group(2)})
    return infos


def parse_fun():
    fun = []
    for name in os.listdir("."):
        if not name.endswith(".go"):
            continue
        for line in open(name):
            res = re.match(r"func \(c \*Config\) ((\w+)\((.*?)\)(.*?)) {", line)
            if not res:
                continue
            params = [i.strip().split(" ")[0] for i in res.group(3).split(",")]
            info = {
                "define": res.group(1),
                "name": res.group(2),
                "params": ", ".join([i if i != "opts" else "opts..." for i in params]),
                "return": res.group(4)
            }
            if info["name"][0].islower():
                continue
            fun.append(info)
    return fun


get_header = """// this file is auto generate by autogen.py. do not edit!
package config

import (
	"net"
	"time"
)
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

def generate_get_bind(infos):
    out = open("autogen_get.go", "w")
    out.write(get_header)
    for info in infos:
        out.write(get_tpl.format(**info))
    for info in infos:
        out.write(gete_tpl.format(**info))
    for info in infos:
        out.write(getp_tpl.format(**info))
    for info in infos:
        out.write(getd_tpl.format(**info))
    out.close()


bind_header = """// this file is auto generate by autogen.py. do not edit!
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


def generate_bind(infos):
    out = open("autogen_bind.go", "w")
    out.write(bind_header)
    for info in infos:
        out.write(atomic_type_tpl.format(**info))
    for info in infos:
        out.write(bind_var_tpl.format(**info))
    for info in infos:
        out.write(bind_tpl.format(**info))
    out.close()


global_export_header = """// this file is auto generate by autogen.py. do not edit!
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

global_export_tpl1 = """
func {define} {{
	return gconf.{name}({params})
}}
"""

global_export_tpl2 = """
func {define} {{
	gconf.{name}({params})
}}
"""


def generate_global_export(infos):
    out = open("autogen_global_export.go", "w")
    out.write(global_export_header)
    items = ""
    for info in infos:
        if info["return"]:
            out.write(global_export_tpl1.format(**info))
        else:
            out.write(global_export_tpl2.format(**info))
    out.write(items)
    out.close()


def main():
    generate_get_bind(parse_info("../cast/cast.go"))
    generate_global_export(parse_fun())


if __name__ == "__main__":
    main()
