package logger

import (
	"bytes"
	"os"

	"github.com/pkg/errors"
)

func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{}
}

func NewStdoutWriterWithOptions(options *StdoutWriterOptions) (*StdoutWriter, error) {
	formatter, err := NewFormatterWithOptions(&options.Formatter)
	if err != nil {
		return nil, errors.WithMessage(err, "NewFormatterWithOptions failed")
	}
	return &StdoutWriter{formatter: formatter}, nil
}

type StdoutWriter struct {
	formatter Formatter
}

func (r *StdoutWriter) Write(v interface{}) error {
	buf, err := r.formatter.Format(v)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(buf)
	b.WriteString("\n")
	_, err = os.Stdout.Write(b.Bytes())
	if err != nil {
		return err
	}
	return nil
}

type StdoutWriterOptions struct {
	Formatter FormatterOptions
}
