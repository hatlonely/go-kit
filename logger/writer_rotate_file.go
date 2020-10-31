package logger

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type RotateFileWriter struct {
	out *rotatelogs.RotateLogs
}

func NewRotateFileWriterWithConfig(cfg *config.Config, opts ...refx.Option) (*RotateFileWriter, error) {
	options := defaultRotateFileWriterOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, err
	}
	return NewRotateFileWriterWithOptions(&options)
}

func NewRotateFileWriterWithOptions(options *RotateFileWriterOptions) (*RotateFileWriter, error) {
	abs, err := filepath.Abs(options.Filename)
	if err != nil {
		return nil, err
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
		out: out,
	}, nil
}

func NewRotateFileWriter(opts ...RotateFileWriterOption) (*RotateFileWriter, error) {
	options := defaultRotateFileWriterOptions
	for _, opt := range opts {
		opt(&options)
	}

	return NewRotateFileWriterWithOptions(&RotateFileWriterOptions{
		Filename: options.Filename,
		MaxAge:   options.MaxAge,
	})
}

func (r *RotateFileWriter) Write(v interface{}) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(buf)
	b.WriteString("\n")
	_, err = r.out.Write(b.Bytes())
	if err != nil {
		return err
	}
	return nil
}

type RotateFileWriterOptions struct {
	Filename string
	MaxAge   time.Duration
}

var defaultRotateFileWriterOptions = RotateFileWriterOptions{
	MaxAge: time.Hour,
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
