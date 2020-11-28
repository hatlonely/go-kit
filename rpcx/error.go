//go:generate protoc -I. --gofast_out=plugins=grpc,paths=source_relative:. error.proto

package rpcx

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewError(err error, status codes.Code, code string, message string) *Error {
	return &Error{
		err:  err,
		code: status,
		Detail: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

func NewErrorf(status codes.Code, code string, format string, args ...interface{}) *Error {
	str := fmt.Sprintf(format, args...)
	return &Error{
		err:  errors.New(str),
		code: status,
		Detail: &ErrorDetail{
			Code:    code,
			Message: str,
		},
	}
}

func NewInternalError(err error) *Error {
	return &Error{
		err:  err,
		code: codes.Internal,
		Detail: &ErrorDetail{
			Status:  http.StatusInternalServerError,
			Code:    "InternalError",
			Message: err.Error(),
		},
	}
}

type Error struct {
	err    error
	code   codes.Code
	Detail *ErrorDetail
}

func (e *Error) SetRequestID(requestID string) *Error {
	e.Detail.RequestID = requestID
	return e
}

func (e *Error) SetRefer(refer string) *Error {
	e.Detail.Refer = refer
	return e
}

func (e *Error) SetStatus(status int64) *Error {
	e.Detail.Status = status
	return e
}

func (e *Error) Error() string {
	return e.Detail.Message
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
	s, _ := status.New(e.code, e.err.Error()).WithDetails(e.Detail)
	return s
}
