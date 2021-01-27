package wrap

import (
	"reflect"
)

func RegisterErrCode(v interface{}, fun func(err error) string) {
	errCodeMap[reflect.TypeOf(v)] = fun
}

var errCodeMap = map[reflect.Type]func(err error) string{}

func ErrCode(err error) string {
	if err == nil {
		return "OK"
	}

	if fun, ok := errCodeMap[reflect.TypeOf(err)]; ok {
		return fun(err)
	}

	return "Unknown"
}
