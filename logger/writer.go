package logger

import (
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type Writer interface {
	Write(kvs map[string]interface{}) error
}

func NewWriterWithConfig(cfg *config.Config, opts ...refx.Option) (Writer, error) {
	var options WriterOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "cfg.Unmarshal failed.")
	}
	return NewWriterWithOptions(&options)
}

func NewWriterWithOptions(options *WriterOptions) (Writer, error) {
	switch options.Type {
	case "RotateFile":
		return NewRotateFileWriterWithOptions(&options.RotateFileWriter)
	case "Stdout":
		return NewStdoutWriterWithOptions(&options.StdoutWriter)
	}
	return nil, errors.Errorf("unsupported writer type [%v]", options.Type)
}

type WriterOptions struct {
	Type             string `dft:"Stdout" rule:"x in ['RotateFile', 'Stdout']"`
	RotateFileWriter RotateFileWriterOptions
	StdoutWriter     StdoutWriterOptions
}
