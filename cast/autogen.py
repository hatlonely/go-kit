#!/usr/bin/env python3

import re


def parse_info(filename):
    infos = []
    for line in open(filename):
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


to_slice_header = """// this file is auto generate by autogen.py. do not edit!
package cast

import (
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
)
"""


to_slice_item_tpl = """
func To{name}SliceE(v interface{{}}) ([]{type}, error) {{
	switch v.(type) {{
	case []{type}:
		return v.([]{type}), nil
	case []interface{{}}:
		vs := make([]{type}, len(v.([]interface{{}})), 0)
		for _, i := range v.([]interface{{}}) {{
			val, err := To{name}E(i)
			if err != nil {{
				return nil, errors.WithMessage(err, "cast failed")
			}}
			vs = append(vs, val)
		}}
		return vs, nil
	case []string:
		vs := make([]{type}, len(v.([]interface{{}})), 0)
		for _, i := range v.([]interface{{}}) {{
			val, err := To{name}E(i)
			if err != nil {{
				return nil, errors.WithMessage(err, "cast failed")
			}}
			vs = append(vs, val)
		}}
		return vs, nil
	case string:
		var vs []{type}
		for _, i := range strings.Split(v.(string), ",") {{
			i = strings.TrimSpace(i)
			val, err := To{name}E(i)
			if err != nil {{
				return nil, errors.WithMessage(err, "cast failed")
			}}
			vs = append(vs, val)
		}}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []{type}", reflect.TypeOf(v))
	}}
}}
"""


to_string_slice = """
func ToStringSliceE(v interface{}) ([]string, error) {
	switch v.(type) {
	case []string:
		return v.([]string), nil
	case []interface{}:
		vs := make([]string, len(v.([]interface{})), 0)
		for _, i := range v.([]interface{}) {
			val, err := ToStringE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	case string:
		var vs []string
		for _, i := range strings.Split(v.(string), ",") {
			i = strings.TrimSpace(i)
			val, err := ToStringE(i)
			if err != nil {
				return nil, errors.WithMessage(err, "cast failed")
			}
			vs = append(vs, val)
		}
		return vs, nil
	default:
		return nil, errors.Errorf("type %v cannot convert []string", reflect.TypeOf(v))
	}
}
"""


def generate_to_slice(infos):
    out = open("autogen_to_slice.go", "w")
    out.write(to_slice_header)
    items = ""
    for info in infos:
        if info["type"] == "string":
            items += to_string_slice
        else:
            items += to_slice_item_tpl.format(**info)
    out.write(items)
    out.close()


def main():
    infos1 = parse_info("cast.go")
    generate_to_slice(infos1)
    infos2 = parse_info("autogen_to_slice.go")
    generate_set_interface([*infos1, *infos2])
    generate_cast([*infos1, *infos2])


if __name__ == "__main__":
    main()
