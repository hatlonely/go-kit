package logger

import (
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
)

type Writer interface {
	Write(v interface{}) error
}

func NewWriterWithConfig(conf *config.Config) (Writer, error) {
	switch conf.GetString("type") {
	case "RotateFile":
		return NewRotateFileWriterWithConfig(conf)
	case "Stdout":
		return NewStdoutWriter(), nil
	}

	return nil, errors.Errorf("no such type %v", conf.GetString("Type"))
}
