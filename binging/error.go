package binging

import (
	"fmt"
	"io"
)

type Error struct {
	Key  string
	Code ErrCode
	Err  error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code [%v] key [%v] : %v", e.Code, e.Key, e.Err)
	}
	return fmt.Sprintf("code [%v] key [%v]", e.Code, e.Key)
}

func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", e.Err)
			_, _ = io.WriteString(s, e.Error())
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}

//go:generate stringer -type ErrCode -trimprefix Err
type ErrCode int

const (
	ErrInvalidDstType       ErrCode = 1
	ErrMissingRequiredField ErrCode = 2
	ErrInvalidFormat        ErrCode = 3
)
