package wrap

import (
	"fmt"
	"net/http"
	"strings"

	alierr "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

func init() {
	RegisterRateLimiterGroup("LocalGroup", NewLocalGroupRateLimiterWithOptions)
	RegisterRateLimiterGroup("LocalShare", NewLocalShareRateLimiterWithOptions)

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

	RegisterRetryRetryIf("POP", func(err error) bool {
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
	RegisterRetryRetryIf("OSS", func(err error) bool {
		switch e := err.(type) {
		case oss.ServiceError:
			if e.StatusCode >= http.StatusInternalServerError {
				return true
			}
			if strings.Contains(e.Error(), "timeout") {
				return true
			}
		}
		return false
	})
	RegisterRetryRetryIf("OTS", func(err error) bool {
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
