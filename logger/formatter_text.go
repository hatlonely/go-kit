package logger

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

func NewTextFormatterWithOptions(options *TextFormatOptions) (*TextFormatter, error) {
	tpl, err := template.New(`__Logger_Text_Formatter__`).Parse(options.Template)
	if err != nil {
		return nil, errors.Wrapf(err, "template.Parse format [%v] failed", options.Template)
	}
	return &TextFormatter{
		tpl: tpl,
	}, nil
}

type TextFormatter struct {
	tpl *template.Template
}

func (f *TextFormatter) Format(info *Info) ([]byte, error) {
	var buf bytes.Buffer

	err := f.tpl.Execute(&buf, info)
	return buf.Bytes(), err
}

type TextFormatOptions struct {
	Template string `dft:"{{.Time}} [{{.Level}}] [{{.Caller}}:{{.File}}] [{{.Fields}}] {{.Data}}"`
}
