#!/usr/bin/env python3

import re


def parse_info():
    infos = []
    for line in open("cast.go"):
        res = re.match(r"func To(.*?)E\(v interface{}\) \((.*?), error\).*", line)
        if not res:
            continue
        infos.append({"name": res.group(1), "type": res.group(2)})
    return infos


set_interface_header = """// this file is auto generate by autogen.py. do not edit!
package cast

import (
	"fmt"
	"net"
	"reflect"
	"time"
)
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

set_interface_item_tpl = """	case *{type}:
		v, err := To{name}E(src)
		if err != nil {{
			return err
		}}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
"""

def generate_set_interface(infos):
    out = open("autogen_set_interface.go", "w")
    out.write(set_interface_header)
    items = ""
    for info in infos:
        items += set_interface_item_tpl.format(**info)
    out.write(set_interface_tpl.format(items=items))
    out.close()


cast_header = """// this file is auto generate by autogen.py. do not edit!
package cast

import (
	"net"
	"time"
)
"""

cast_item_tpl = """
func To{name}(val interface{{}}) {type} {{
	if v, err := To{name}E(val); err == nil {{
		return v
	}}
	var v {type}
	return v
}}

func To{name}D(val interface{{}}, defaultValue {type}) {type} {{
	if v, err := To{name}E(val); err == nil {{
		return v
	}}
	return defaultValue
}}

func To{name}P(val interface{{}}) {type} {{
	v, err := To{name}E(val)
	if err != nil {{
		panic(err)
	}}
	return v
}}
"""


def generate_cast(infos):
    out = open("autogen_cast.go", "w")
    out.write(cast_header)
    items = ""
    for info in infos:
        items += cast_item_tpl.format(**info)
    out.write(items)
    out.close()


def main():
    infos = parse_info()
    generate_set_interface(infos)
    generate_cast(infos)


if __name__ == "__main__":
    main()
