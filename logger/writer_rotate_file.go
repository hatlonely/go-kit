package logger

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type RotateFileWriter struct {
	out       *rotatelogs.RotateLogs
	formatter Formatter
}

func NewRotateFileWriterWithConfig(cfg *config.Config, opts ...refx.Option) (*RotateFileWriter, error) {
	options := defaultRotateFileWriterOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, err
	}
	return NewRotateFileWriterWithOptions(&options)
}

func NewRotateFileWriterWithOptions(options *RotateFileWriterOptions, opts ...refx.Option) (*RotateFileWriter, error) {
	formatter, err := NewFormatterWithOptions(&options.Formatter, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "NewFormatterWithOptions failed")
	}

	abs, err := filepath.Abs(options.Filename)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
		return nil, errors.Wrapf(err, "os.MkdirAll [%s] failed", filepath.Dir(abs))
	}

	out, err := rotatelogs.New(
		abs+".%Y%m%d%H",
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithLinkName(abs),
		rotatelogs.WithMaxAge(options.MaxAge),
	)
	if err != nil {
		return nil, err
	}

	return &RotateFileWriter{
		out:       out,
		formatter: formatter,
	}, nil
}

func NewRotateFileWriter(opts ...RotateFileWriterOption) (*RotateFileWriter, error) {
	options := defaultRotateFileWriterOptions
	for _, opt := range opts {
		opt(&options)
	}

	return NewRotateFileWriterWithOptions(&options)
}

func (w *RotateFileWriter) Write(info *Info) error {
	buf, err := w.formatter.Format(info)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(buf)
	b.WriteString("\n")
	_, err = w.out.Write(b.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (w *RotateFileWriter) Close() error {
	return w.out.Close()
}

type RotateFileWriterOptions struct {
	Filename  string        `bind:"required"`
	MaxAge    time.Duration `dft:"24h"`
	Formatter FormatterOptions
}

var defaultRotateFileWriterOptions = RotateFileWriterOptions{
	MaxAge: time.Hour,
	Formatter: FormatterOptions{
		Type: "Json",
	},
}

type RotateFileWriterOption func(*RotateFileWriterOptions)

func WithRotateFilename(filename string) RotateFileWriterOption {
	return func(options *RotateFileWriterOptions) {
		options.Filename = filename
	}
}

func WithRotateMaxAge(maxAge time.Duration) RotateFileWriterOption {
	return func(options *RotateFileWriterOptions) {
		options.MaxAge = maxAge
	}
}
