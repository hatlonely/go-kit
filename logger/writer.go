package logger

import (
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
)

type Writer interface {
	Write(v interface{}) error
}

func NewWriterWithConfig(conf *config.Config) (Writer, error) {
	switch conf.GetString("Type") {
	case "RotateFile":
		return NewRotateFileWriterWithConfig(conf)
	}

	return nil, errors.Errorf("no such type %v", conf.GetString("Type"))
}
