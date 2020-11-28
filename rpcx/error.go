//go:generate protoc -I. --gofast_out=plugins=grpc,paths=source_relative:. error.proto

package rpcx

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewErrorWithoutRefer(err error, status codes.Code, code string, message string) error {
	return NewError(err, status, code, message, "")
}

func NewError(err error, status codes.Code, code string, message string, refer string) error {
	return &Error{
		Err: err,
		Detail: &ErrorDetail{
			Status:  int64(status),
			Code:    code,
			Message: message,
			Refer:   refer,
		},
	}
}

func NewErrorWithoutReferf(status codes.Code, code string, format string, args ...interface{}) error {
	return NewErrorf(status, code, "", format, args...)
}

func NewErrorf(status codes.Code, code string, refer string, format string, args ...interface{}) error {
	str := fmt.Sprintf(format, args...)
	err := errors.New(str)
	return &Error{
		Err: err,
		Detail: &ErrorDetail{
			Status:  int64(status),
			Code:    code,
			Message: str,
			Refer:   refer,
		},
	}
}

func NewInternalError(err error) *Error {
	return &Error{
		Err: err,
		Detail: &ErrorDetail{
			Status:  int64(codes.Internal),
			Code:    "InternalError",
			Message: err.Error(),
		},
	}
}

type Error struct {
	Err    error
	Detail *ErrorDetail
}

func (e *Error) SetRequestID(requestID string) *Error {
	e.Detail.RequestID = requestID
	return e
}

func (e *Error) Error() string {
	return e.Detail.Message
}

func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", e.Err)
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}

func (e *Error) ToStatus() *status.Status {
	s, _ := status.New(codes.Code(e.Detail.Status), e.Detail.Message).WithDetails(e.Detail)
	return s
}
