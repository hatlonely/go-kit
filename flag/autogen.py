#!/usr/bin/env python3

import re
import os


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
            res = re.match(r"func \(f \*Flag\) ((\w+)\((.*?)\)(.*?)) {", line)
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

bind_header = """// this file is auto generate by autogen.py. do not edit!
package flag

import (
	"net"
	"reflect"
	"time"

	"github.com/hatlonely/go-kit/cast"
)
"""

bind_item_tpl = """
func (f *Flag) {name}Var(p *{type}, name string, defaultValue {type}, usage string) {{
	var v {type}
	*p = defaultValue
	if err := f.addFlag(name, usage, reflect.TypeOf(v), false, "", cast.ToString(defaultValue), false); err != nil {{
		panic(err)
	}}
	info, _ := f.GetInfo(name)
	info.OnParse = func(val string) error {{
		return cast.SetInterface(p, val)
	}}
}}

func (f *Flag) {name}(name string, defaultValue {type}, usage string) *{type} {{
	var v {type}
	f.{name}Var(&v, name, defaultValue, usage)
	return &v
}}
"""


def generate_bind(infos):
    out = open("autogen_bind.go", "w")
    out.write(bind_header)
    items = ""
    for info in infos:
        items += bind_item_tpl.format(**info)
    out.write(items)
    out.close()


global_export_header = """// this file is auto generate by autogen.py. do not edit!
package flag

import (
	"net"
	"os"
	"reflect"
	"time"
)

var gflag = NewFlag(os.Args[0])

func Instance() *Flag {
	return gflag
}

func Parse() error {
	return ParseArgs(os.Args[1:])
}
"""

global_export_tpl1 = """
func {define} {{
	return gflag.{name}({params})
}}
"""

global_export_tpl2 = """
func {define} {{
	gflag.{name}({params})
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


get_header = """// this file is auto generate by autogen.py. do not edit!
package flag

import (
	"net"
	"time"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/cast"
)
"""

get_item_tpl = """
func (f *Flag) Get{name}E(key string) ({type}, error) {{
	v, ok := f.Get(key)
	if !ok {{
		var res {type}
		return res, errors.Errorf("get failed: key not exist. key: [%v]", key)
	}}
	return cast.To{name}E(v)
}}

func (f *Flag) Get{name}(key string) {type} {{
	var res {type}
	v, err := f.Get{name}E(key)
	if err != nil {{
		return res
	}}
	return v
}}

func (f *Flag) Get{name}D(key string, val {type}) {type} {{
	v, err := f.Get{name}E(key)
	if err != nil {{
		return val
	}}
	return v
}}

func (f *Flag) Get{name}P(key string) {type} {{
	v, err := f.Get{name}E(key)
	if err != nil {{
		panic(err)
	}}
	return v
}}
"""


def generate_get(infos):
    out = open("autogen_get.go", "w")
    out.write(get_header)
    for info in infos:
        out.write(get_item_tpl.format(**info))
    out.close()


def main():
    generate_bind([*parse_info("../cast/cast.go"), *parse_info("../cast/autogen_to_slice.go")])
    generate_get([*parse_info("../cast/cast.go"), *parse_info("../cast/autogen_to_slice.go")])
    generate_global_export(parse_fun())


if __name__ == "__main__":
    main()
