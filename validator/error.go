package validator

import (
	"fmt"
	"io"
)

type Error struct {
	Code ErrCode
	Key  string
	Val  interface{}
	Tag  string
	Err  error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code [%v] key [%v] val [%v] tag [%v]: %v", e.Code, e.Key, e.Val, e.Tag, e.Err)
	}
	return fmt.Sprintf("code [%v] key [%v] val [%v] tag [%v]", e.Code, e.Key, e.Val, e.Tag)
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
	ErrEvalFailed   ErrCode = 1
	ErrRuleNotMatch ErrCode = 2
)
