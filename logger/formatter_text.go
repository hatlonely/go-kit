package logger

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

func NewTextFormatterWithOptions(options *TextFormatOptions) (*TextFormatter, error) {
	tpl, err := template.New(`__Logger_Text_Formatter__`).Parse(options.Format)
	if err != nil {
		return nil, errors.Wrapf(err, "template.Parse format [%v] failed", options.Format)
	}
	return &TextFormatter{
		tpl: tpl,
	}, nil
}

type TextFormatter struct {
	tpl *template.Template
}

func (f *TextFormatter) Format(kvs map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer

	err := f.tpl.Execute(&buf, kvs)
	return buf.Bytes(), err
}

type TextFormatOptions struct {
	Format string `dft:"{{.time}} [{{.level}}] [{{.caller}}:{{.file}}] {{.data}}"`
}
