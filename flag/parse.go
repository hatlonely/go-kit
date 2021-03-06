package flag

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type ParseOptions struct {
	JsonVal bool
}

type ParseOption func(options *ParseOptions)

func WithJsonVal() ParseOption {
	return func(options *ParseOptions) {
		options.JsonVal = true
	}
}

func (f *Flag) ParseArgs(args []string, opts ...ParseOption) error {
	var options ParseOptions
	for _, opt := range opts {
		opt(&options)
	}

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
			key, val := key[0:idx], key[idx+1:]
			if err := f.SetWithOptions(key, val, &options); err != nil {
				return errors.WithMessage(err, fmt.Sprintf("parse failed, key [%v]", key))
			}
		} else if info, ok := f.GetInfo(key); ok {
			if info.Type != reflect.TypeOf(false) { // 选项不是 bool，后面必有一个值
				if i+1 >= len(args) {
					return errors.Errorf("miss argument for non boolean key [%v]", key)
				}
				if err := f.SetWithOptions(key, args[i+1], &options); err != nil {
					return errors.WithMessage(err, fmt.Sprintf("parse failed, key [%v]", key))
				}
				i++
			} else { // 选项为 bool 类型，如果后面的值为合法的 bool 值，否则设置为 true
				val := "true"
				if i+1 < len(args) && isBoolValue(args[i+1]) {
					val = args[i+1]
					i++
				}
				if err := f.SetWithOptions(key, val, &options); err != nil {
					return errors.WithMessage(err, "parse failed")
				}
			}
		} else if f.allBoolFlag(key) { // -aux 全是 bool 选项，-aux 和 -a -u -x 等效
			for i := 0; i < len(key); i++ {
				k := key[i : i+1]
				if err := f.SetWithOptions(k, "true", &options); err != nil {
					return errors.WithMessage(err, "parse failed")
				}
			}
		} else if _, ok := f.GetInfo(key[0:1]); ok { // -p123456 和 -p 123456 等效
			key, val := key[0:1], key[1:]
			if err := f.SetWithOptions(key, val, &options); err != nil {
				return errors.WithMessage(err, "parse failed")
			}
		} else {
			if i+1 >= len(args) {
				return errors.Errorf("miss argument for non boolean key [%v]", key)
			}
			if err := f.SetWithOptions(key, args[i+1], &options); err != nil {
				return errors.WithMessage(err, fmt.Sprintf("parse failed, key [%v]", key))
			}
			i++
		}
	}

	for i, val := range arguments {
		if i >= len(f.arguments) {
			break
		}
		if err := f.SetWithOptions(f.nameKeyMap[f.arguments[i]], val, &options); err != nil {
			return errors.WithMessage(err, "parse failed")
		}
	}

	for key, info := range f.keyInfoMap {
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
