package flag

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func (f *Flag) Parse(args []string) error {
	var arguments []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			arguments = append(arguments, arg)
			continue
		}
		key := arg[1:]
		if strings.HasPrefix(arg, "--") {
			key = arg[2:]
		}
		if strings.Contains(key, "=") { // 选项中含有等号，按照等号分割成 name val
			idx := strings.Index(key, "=")
			val := key[idx+1:]
			key := key[0:idx]
			_, ok := f.GetInfo(key)
			if !ok {
				return errors.Errorf("unknown key [%v]", key)
			}
			if err := f.Set(key, val); err != nil {
				return errors.WithMessage(err, fmt.Sprintf("parse failed, key [%v]", key))
			}
		} else if info, ok := f.GetInfo(key); ok {
			if info.Type != reflect.TypeOf(false) { // 选项不是 bool，后面必有一个值
				if i+1 >= len(args) {
					return errors.Errorf("miss argument for non boolean key [%v]", key)
				}
				if err := f.Set(key, args[i+1]); err != nil {
					return errors.WithMessage(err, fmt.Sprintf("parse failed, key [%v]", key))
				}
				i++
			} else { // 选项为 bool 类型，如果后面的值为合法的 bool 值，否则设置为 true
				val := "true"
				if i+1 < len(args) && isBoolValue(args[i+1]) {
					val = args[i+1]
					i++
				}
				if err := f.Set(key, val); err != nil {
					return errors.WithMessage(err, "parse failed")
				}
			}
		} else if f.allBoolFlag(key) { // -aux 全是 bool 选项，-aux 和 -a -u -x 等效
			for i := 0; i < len(key); i++ {
				k := key[i : i+1]
				if err := f.Set(k, "true"); err != nil {
					return errors.WithMessage(err, "parse failed")
				}
			}
		} else { // -p123456 和 -p 123456 等效
			key := key[0:1]
			val := key[1:]
			_, ok := f.GetInfo(key)
			if !ok {
				return errors.Errorf("unknown key [%v]", key)
			}
			if err := f.Set(key, val); err != nil {
				return errors.WithMessage(err, "parse failed")
			}
		}
	}

	for i, val := range arguments {
		if i >= len(f.arguments) {
			break
		}
		if err := f.Set(f.arguments[i], val); err != nil {
			return errors.WithMessage(err, "parse failed")
		}
	}

	for key, info := range f.flagInfos {
		if info.Required && !info.Assigned {
			return errors.Errorf("option [%v] is required, but not assigned", key)
		}
	}

	return nil
}

func isBoolValue(val string) bool {
	_, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return true
}

func (f *Flag) allBoolFlag(key string) bool {
	for i := 0; i < len(key); i++ {
		info, ok := f.GetInfo(key[i : i+1])
		if !ok || info.Type != reflect.TypeOf(false) {
			return false
		}
	}

	return true
}
