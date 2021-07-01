package logger

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type Writer interface {
	Write(info *Info) error
	Close() error
}

func RegisterWriter(key string, constructor interface{}) {
	refx.Register(reflect.TypeOf((*Writer)(nil)).Elem(), key, constructor)
}

func NewWriterWithOptions(options *refx.TypeOptions, opts ...refx.Option) (Writer, error) {
	v, err := refx.New(reflect.TypeOf((*Writer)(nil)).Elem(), options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "refx.New failed")
	}
	return v.(Writer), nil
}

func NewWriterWithConfig(cfg *config.Config, opts ...refx.Option) (Writer, error) {
	var options refx.TypeOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "cfg.Unmarshal failed.")
	}
	return NewWriterWithOptions(&options, opts...)
}
