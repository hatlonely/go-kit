package logger

import (
	"github.com/pkg/errors"
)

type Formatter interface {
	Format(v interface{}) ([]byte, error)
}

func NewFormatterWithOptions(options *FormatterOptions) (Formatter, error) {
	switch options.Type {
	case "Json":
		return JsonFormatter{}, nil
	case "Text":
		return NewTextFormatterWithOptions(&options.TextFormat)
	}
	return nil, errors.Errorf("unsupported formatter type [%v]", options.Type)
}

type FormatterOptions struct {
	Type       string `dft:"Json" rule:"x in ['Json', 'Text']"`
	TextFormat TextFormatOptions
}
