import re


def parse_info():
    infos = []
    for line in open("../cast/cast.go"):
        res = re.match(r"func To(.*?)E\(v interface{}\) \((.*?), error\).*", line)
        if not res:
            continue
        infos.append({"name": res.group(1), "type": res.group(2)})
    return infos


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


def main():
    infos = parse_info()
    generate_bind(infos)


if __name__ == "__main__":
    main()
