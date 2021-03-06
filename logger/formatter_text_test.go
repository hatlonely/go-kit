package logger

import (
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTextFormatter(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339Nano, "2020-11-07T16:07:36.65798+08:00")
	info := &Info{
		Time:   ts,
		Level:  LevelInfo,
		File:   "formatter_test.go:18",
		Caller: "TextFormatter",
		Fields: map[string]interface{}{
			"meta1": "val1",
			"meta2": "val2",
		},
		Data: map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
		},
	}

	Convey("TestTextFormatter", t, func() {
		formatter, err := NewTextFormatterWithOptions(&TextFormatOptions{
			Template: "{{.Time}} [{{.Level}}] [{{.Caller}}:{{.File}}] [{{.Fields}}] {{.Data}}",
		})

		So(err, ShouldBeNil)
		So(reflect.TypeOf(formatter).String(), ShouldEqual, "*logger.TextFormatter")

		buf, err := formatter.Format(info)
		So(err, ShouldBeNil)
		So(string(buf), ShouldEqual, `2020-11-07 16:07:36.65798 +0800 CST [Info] [TextFormatter:formatter_test.go:18] [map[meta1:val1 meta2:val2]] map[key1:val1 key2:val2]`)
	})
}
