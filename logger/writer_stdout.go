package logger

import (
	"bytes"
	"os"

	"github.com/pkg/errors"
)

func NewStdoutWriter() *StdoutWriter {
	writer, _ := NewStdoutWriterWithOptions(&StdoutWriterOptions{
		Formatter: FormatterOptions{
			Type: "Text",
			TextFormat: TextFormatOptions{
				Format: "{{.time}} [{{.level}}] [{{.caller}}:{{.file}}] {{.data}}",
			},
		},
	})
	return writer
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

func (w *StdoutWriter) Write(kvs map[string]interface{}) error {
	buf, err := w.formatter.Format(kvs)
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
