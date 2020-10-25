//go:generate protoc -I. --gofast_out=plugins=grpc,paths=source_relative:. error.proto

package rpcx

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewError(err error, status codes.Code, requestID string, code string, message string) error {
	return &Error{
		Err: err,
		Info: &EInfo{
			Status:    int64(status),
			RequestID: requestID,
			Code:      code,
			Message:   message,
		},
	}
}

func NewErrorf(status codes.Code, requestID string, code string, format string, args ...interface{}) error {
	str := fmt.Sprintf(format, args...)
	err := errors.New(str)
	return &Error{
		Err: err,
		Info: &EInfo{
			Status:    int64(status),
			RequestID: requestID,
			Code:      code,
			Message:   str,
		},
	}
}

func NewInternalError(err error, requestID string) *Error {
	return &Error{
		Err: err,
		Info: &EInfo{
			Status:    int64(codes.Internal),
			RequestID: requestID,
			Code:      "InternalError",
			Message:   err.Error(),
		},
	}
}

type Error struct {
	Err  error
	Info *EInfo
}

func (e *Error) Error() string {
	return e.Info.Message
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
	s, _ := status.New(codes.Code(e.Info.Status), e.Info.Message).WithDetails(e.Info)
	return s
}
