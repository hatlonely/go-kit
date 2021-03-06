package loggerx

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hatlonely/go-kit/logger"
)

type MarkdownFormatterOptions struct {
	DefaultTitle string
	TitleField   string
}

func NewMarkdownFormatterWithOptions(options *MarkdownFormatterOptions) *MarkdownFormatter {
	return &MarkdownFormatter{
		options: options,
	}
}

type MarkdownFormatter struct {
	options *MarkdownFormatterOptions
}

func (f *MarkdownFormatter) Format(info *logger.Info) ([]byte, error) {
	var buf bytes.Buffer

	title := f.options.DefaultTitle
	if f.options.TitleField != "" && info.Fields[f.options.TitleField] != nil {
		title = fmt.Sprintf("%v", info.Fields[f.options.TitleField])
	}
	buf.WriteString(fmt.Sprintf("## 【%v】 %v\n", info.Level.String(), title))

	buf.WriteString(fmt.Sprintf("【**时间**】 %v\n\n", info.Time.Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf("【**文件**】 `%v`\n\n", info.File))
	buf.WriteString(fmt.Sprintf("【**调用**】 `%v`\n\n", info.Caller))

	for key, val := range info.Fields {
		buf.WriteString(fmt.Sprintf("【**`%v`**】 `%v`\n\n", key, val))
	}

	buf.WriteString("```json\n")
	d, _ := json.MarshalIndent(info.Data, "", "  ")
	buf.Write(d)
	buf.WriteString("\n```\n")

	return buf.Bytes(), nil
}
