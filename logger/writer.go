package logger

import (
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type Writer interface {
	Write(v interface{}) error
}

func NewWriterWithConfig(conf *config.Config, opts ...refx.Option) (Writer, error) {
	switch conf.GetString("type") {
	case "RotateFile":
		return NewRotateFileWriterWithConfig(conf, opts...)
	case "Stdout":
		return NewStdoutWriter(), nil
	}

	return nil, errors.Errorf("no such type %v", conf.GetString("Type"))
}
