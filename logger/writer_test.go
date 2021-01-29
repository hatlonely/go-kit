package logger

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStdoutWriter(t *testing.T) {
	Convey("TestStdoutWriter", t, func() {
		writer, err := NewWriterWithOptions(&WriterOptions{
			Type: "Stdout",
			Options: &StdoutWriterOptions{
				Formatter: FormatterOptions{
					Type: "Json",
				},
			},
		})

		So(err, ShouldBeNil)
		_ = writer.Write(map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
		})
	})
}

func TestRotateFileWriter(t *testing.T) {
	Convey("TestRotateFileWriter", t, func() {
		writer, err := NewWriterWithOptions(&WriterOptions{
			Type: "RotateFile",
			Options: &RotateFileWriterOptions{
				Formatter: FormatterOptions{
					Type: "Json",
				},
				MaxAge:   24 * time.Hour,
				Filename: "log/test.log",
			},
		})

		So(err, ShouldBeNil)
		_ = writer.Write(map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
		})
	})
}
