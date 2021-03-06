package logger

import (
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewWriterWithOptions(t *testing.T) {
	Convey("TestNewWriterWithOptions", t, func() {
		Convey("case stdout", func() {
			writer, err := NewWriterWithOptions(&WriterOptions{
				Type: "Stdout",
				Options: &StdoutWriterOptions{
					Formatter: FormatterOptions{
						Type:    "Json",
						Options: &JsonFormatterOptions{},
					},
				},
			})
			So(err, ShouldBeNil)
			So(reflect.TypeOf(writer).String(), ShouldEqual, "*logger.StdoutWriter")
		})

		Convey("case rotate file", func() {
			writer, err := NewWriterWithOptions(&WriterOptions{
				Type: "RotateFile",
				Options: &RotateFileWriterOptions{
					Filename: "tmp/test.log",
					MaxAge:   24 * time.Hour,
					Formatter: FormatterOptions{
						Type:    "Json",
						Options: &JsonFormatterOptions{},
					},
				},
			})
			So(err, ShouldBeNil)
			So(reflect.TypeOf(writer).String(), ShouldEqual, "*logger.RotateFileWriter")
		})
	})
}
