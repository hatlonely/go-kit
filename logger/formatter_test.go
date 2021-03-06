package logger

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewFormatterWithOptions(t *testing.T) {
	Convey("TestNewFormatterWithOptions", t, func() {
		Convey("case json", func() {
			formatter, err := NewFormatterWithOptions(&FormatterOptions{
				Type:    "Json",
				Options: &JsonFormatterOptions{},
			})
			So(err, ShouldBeNil)
			So(reflect.TypeOf(formatter).String(), ShouldEqual, "*logger.JsonFormatter")
		})

		Convey("case text", func() {
			formatter, err := NewFormatterWithOptions(&FormatterOptions{
				Type: "Text",
				Options: &TextFormatOptions{
					Template: "{{.Time}} [{{.Level}}] [{{.Caller}}:{{.File}}] [{{.Fields}}] {{.Data}}",
				},
			})
			So(err, ShouldBeNil)
			So(reflect.TypeOf(formatter).String(), ShouldEqual, "*logger.TextFormatter")
		})
	})
}
