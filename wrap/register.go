package wrap

import (
	"fmt"
	"net/http"
	"strings"

	alierr "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-tablestore-go-sdk/v5/tablestore"

	"github.com/hatlonely/go-kit/micro"
)

func init() {
	RegisterErrCode(&alierr.ServerError{}, func(err error) string {
		e := err.(*alierr.ServerError)
		return fmt.Sprintf("pop_%v_%v", e.HttpStatus(), e.ErrorCode())
	})
	RegisterErrCode(oss.ServiceError{}, func(err error) string {
		e := err.(oss.ServiceError)
		return fmt.Sprintf("oss_%v_%v", e.StatusCode, e.Code)
	})
	RegisterErrCode(&tablestore.OtsError{}, func(err error) string {
		e := err.(*tablestore.OtsError)
		return fmt.Sprintf("ots_%v_%v", e.HttpStatusCode, e.Code)
	})

	micro.RegisterRetryRetryIf("POP", func(err error) bool {
		switch e := err.(type) {
		case *alierr.ServerError:
			if e.HttpStatus() >= http.StatusInternalServerError {
				return true
			}
			if strings.Contains(e.Error(), "timeout") {
				return true
			}
		}
		return false
	})
	micro.RegisterRetryRetryIf("OSS", func(err error) bool {
		switch e := err.(type) {
		case oss.ServiceError:
			if e.StatusCode >= http.StatusInternalServerError {
				return true
			}
			if strings.Contains(e.Error(), "timeout") {
				return true
			}
			if strings.Contains(e.Error(), "connection reset by peer") {
				return true
			}
			if strings.Contains(e.Error(), "use of closed network connection") {
				return true
			}
		}
		return false
	})
	micro.RegisterRetryRetryIf("OTS", func(err error) bool {
		switch e := err.(type) {
		case *tablestore.OtsError:
			if e.HttpStatusCode >= http.StatusInternalServerError {
				return true
			}
			if strings.Contains(e.Error(), "timeout") {
				return true
			}
		}
		return false
	})
}
