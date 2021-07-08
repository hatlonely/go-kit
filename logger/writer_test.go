package logger

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hatlonely/go-kit/refx"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewWriterWithOptions(t *testing.T) {
	Convey("TestNewWriterWithOptions", t, func() {
		Convey("case stdout", func() {
			writer, err := NewWriterWithOptions(&refx.TypeOptions{
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
			writer, err := NewWriterWithOptions(&refx.TypeOptions{
				Type: "RotateFile",
				Options: &RotateFileWriterOptions{
					Level:    "Info",
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

			_ = os.RemoveAll("tmp")
		})
	})
}
