package loggerx

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hatlonely/go-kit/logger"
)

type HtmlFormatterOptions struct {
	DefaultTitle string
	TitleField   string
}

func NewHtmlFormatterWithOptions(options *HtmlFormatterOptions) *HtmlFormatter {
	return &HtmlFormatter{
		options: options,
	}
}

type HtmlFormatter struct {
	options *HtmlFormatterOptions
}

func (f *HtmlFormatter) Format(info *logger.Info) ([]byte, error) {
	var buf bytes.Buffer

	title := f.options.DefaultTitle
	if f.options.TitleField != "" && info.Fields[f.options.TitleField] != nil {
		title = fmt.Sprintf("%v", info.Fields[f.options.TitleField])
	}
	buf.WriteString(fmt.Sprintf("<h2>【%v】 %v</h2>", info.Level.String(), title))

	buf.WriteString("<ul>\n")
	buf.WriteString(fmt.Sprintf("<li>【时间】 <code>%v</code></li>\n", info.Time.Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf("<li>【文件】 <code>%v</code></li>\n", info.File))
	buf.WriteString(fmt.Sprintf("<li>【调用】 <code>%v</code></li>\n", info.Caller))

	for key, val := range info.Fields {
		buf.WriteString(fmt.Sprintf("<li>【<code>%v</code>】 <code>%v</code></li>\n", key, val))
	}
	buf.WriteString("</ul>")

	buf.WriteString("<pre><code>\n")
	d, _ := json.MarshalIndent(info.Data, "", "  ")
	buf.Write(d)
	buf.WriteString("\n</code></pre>\n")

	fmt.Println(buf.String())
	return buf.Bytes(), nil
}
