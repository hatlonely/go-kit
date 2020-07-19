#!/usr/bin/env python3

import re

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

item_tpl = """	case *{type}:
		v, err := To{name}E(src)
		if err != nil {{
			return err
		}}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
"""

def generate_set_interface():
    infos = []
    for line in open("cast.go"):
        res = re.match(r"func To(.*?)E\(v interface{}\) \((.*?), error\).*", line)
        if not res:
            continue
        infos.append({"name": res.group(1), "type": res.group(2)})
    items = ""
    for info in infos:
        items += item_tpl.format(**info)
    out = open("set_interface.go", "w")
    out.write(set_interface_header)
    out.write(set_interface_tpl.format(items=items))
    out.close()


def main():
    generate_set_interface()


if __name__ == "__main__":
    main()