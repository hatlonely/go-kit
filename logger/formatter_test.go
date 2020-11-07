package logger

import (
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJsonFormatter(t *testing.T) {
	Convey("TestJsonFormatter", t, func() {
		formatter, err := NewFormatterWithOptions(&FormatterOptions{
			Type: "Json",
		})
		So(err, ShouldBeNil)
		So(reflect.TypeOf(formatter).String(), ShouldEqual, "logger.JsonFormatter")
		buf, err := formatter.Format(map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
		})
		So(err, ShouldBeNil)
		So(string(buf), ShouldEqual, `{"key1":"val1","key2":"val2"}`)
	})
}

func TestTextFormatter(t *testing.T) {
	Convey("TestTextFormatter", t, func() {
		formatter, err := NewFormatterWithOptions(&FormatterOptions{
			Type: "Text",
			TextFormat: TextFormatOptions{
				Format: "{{.time}} [{{.level}}] [{{.caller}}:{{.file}}] {{.data}}",
			},
		})
		So(err, ShouldBeNil)
		So(reflect.TypeOf(formatter).String(), ShouldEqual, "*logger.TextFormatter")

		buf, err := formatter.Format(map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"time":      "2020-11-07T16:07:36.65798+08:00",
			"level":     "info",
			"data": map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
			},
			"file":   "formatter_test.go:18",
			"caller": "TextFormatter",
		})
		So(err, ShouldBeNil)
		So(string(buf), ShouldEqual, `2020-11-07T16:07:36.65798+08:00 [info] [TextFormatter:formatter_test.go:18] map[key1:val1 key2:val2]`)
	})
}
