package loggerx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type MarkdownFormatterOptions struct {
	DefaultTitle string
	TitleField   string
}

func NewMarkdownFormatterWithOptions(options *MarkdownFormatterOptions) *MarkdownFormatter {
	return &MarkdownFormatter{
		options: options,
		buildKeys: map[string]bool{
			"time":      true,
			"timestamp": true,
			"level":     true,
			"file":      true,
			"caller":    true,
			"data":      true,
		},
	}
}

type MarkdownFormatter struct {
	options   *MarkdownFormatterOptions
	buildKeys map[string]bool
}

func (f *MarkdownFormatter) Format(kvs map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer

	title := f.options.DefaultTitle
	if f.options.TitleField != "" && kvs[f.options.TitleField] != nil {
		title = fmt.Sprintf("%v", kvs[f.options.TitleField])
	}
	buf.WriteString(fmt.Sprintf("## 【%v】 %v\n", kvs["level"], title))

	buf.WriteString(fmt.Sprintf("【**时间**】 %v\n\n", time.Unix(kvs["timestamp"].(int64), 0).Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf("【**文件**】 `%v`\n\n", kvs["file"]))
	buf.WriteString(fmt.Sprintf("【**调用**】 `%v`\n\n", kvs["caller"]))

	for key, val := range kvs {
		if f.buildKeys[key] {
			continue
		}
		buf.WriteString(fmt.Sprintf("【**`%v`**】 `%v`\n\n", key, val))
	}

	buf.WriteString("```json\n")
	d, _ := json.MarshalIndent(kvs["data"], "  ", "  ")
	buf.Write(d)
	buf.WriteString("\n```\n")

	return buf.Bytes(), nil
}
