//go:generate protoc -I. -I../rpc-api --gofast_out . --gofast_opt paths=source_relative error.proto
//go:generate protoc -I. -I../rpc-api --go_out . --go_opt paths=source_relative error.proto

package rpcx

import (
	"fmt"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewError(err error, rpcCode codes.Code, code string, message string) *Error {
	return &Error{
		err:  err,
		Code: rpcCode,
		Detail: &ErrorDetail{
			Status:  int32(runtime.HTTPStatusFromCode(rpcCode)),
			Code:    code,
			Message: message,
		},
	}
}

func NewErrorf(err error, rpcCode codes.Code, code string, format string, args ...interface{}) *Error {
	str := fmt.Sprintf(format, args...)
	if err == nil {
		err = errors.New(str)
	}
	return &Error{
		err:  err,
		Code: rpcCode,
		Detail: &ErrorDetail{
			Status:  int32(runtime.HTTPStatusFromCode(rpcCode)),
			Code:    code,
			Message: str,
		},
	}
}

func NewInternalError(err error) *Error {
	return &Error{
		err:  err,
		Code: codes.Internal,
		Detail: &ErrorDetail{
			Status:  http.StatusInternalServerError,
			Code:    "InternalError",
			Message: "unknown error",
		},
	}
}

type Error struct {
	err    error
	Code   codes.Code
	Detail *ErrorDetail `json:"detail,omitempty"`
}

func (e *Error) SetRequestID(requestID string) *Error {
	e.Detail.RequestID = requestID
	return e
}

func (e *Error) SetRefer(refer string) *Error {
	e.Detail.Refer = refer
	return e
}

func (e *Error) SetStatus(status int) *Error {
	e.Detail.Status = int32(status)
	return e
}

func (e *Error) SetMessage(message string) *Error {
	e.Detail.Message = message
	return e
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", e.err)
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}

func (e *Error) ToStatus() *status.Status {
	s, _ := status.New(e.Code, e.err.Error()).WithDetails(e.Detail)
	return s
}
