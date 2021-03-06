package logger

import (
	"bytes"
	"os"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

func NewStdoutWriter() *StdoutWriter {
	writer, _ := NewStdoutWriterWithOptions(&StdoutWriterOptions{
		Formatter: FormatterOptions{
			Type: "Text",
			Options: &TextFormatOptions{
				Template: "{{.Time}} [{{.Level}}] [{{.Caller}}:{{.File}}] [{{.Fields}}] {{.Data}}",
			},
		},
	})
	return writer
}

func NewStdoutWriterWithOptions(options *StdoutWriterOptions, opts ...refx.Option) (*StdoutWriter, error) {
	formatter, err := NewFormatterWithOptions(&options.Formatter, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "NewFormatterWithOptions failed")
	}
	return &StdoutWriter{formatter: formatter}, nil
}

type StdoutWriter struct {
	formatter Formatter
}

func (w *StdoutWriter) Write(info *Info) error {
	buf, err := w.formatter.Format(info)
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

func (w *StdoutWriter) Close() error {
	return nil
}

type StdoutWriterOptions struct {
	Formatter FormatterOptions
}
