package wrap

import (
	"fmt"
	"reflect"

	alierr "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

func RegisterErrCode(v interface{}, fun func(err error) string) {
	errCodeMap[reflect.TypeOf(v)] = fun
}

var errCodeMap = map[reflect.Type]func(err error) string{
	reflect.TypeOf(oss.ServiceError{}): func(err error) string {
		e := err.(oss.ServiceError)
		return fmt.Sprintf("oss_%v_%v", e.StatusCode, e.Code)
	},
	reflect.TypeOf(&tablestore.OtsError{}): func(err error) string {
		e := err.(*tablestore.OtsError)
		return fmt.Sprintf("ots_%v_%v", e.HttpStatusCode, e.Code)
	},
	reflect.TypeOf(&alierr.ServerError{}): func(err error) string {
		e := err.(*alierr.ServerError)
		return fmt.Sprintf("pop_%v_%v", e.HttpStatus(), e.ErrorCode())
	},
}

func ErrCode(err error) string {
	if err == nil {
		return "OK"
	}

	rt := reflect.TypeOf(err)
	if fun, ok := errCodeMap[rt]; ok {
		return fun(err)
	}

	return "Unknown"
}
