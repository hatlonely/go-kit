package logger

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRotateFileWriter(t *testing.T) {
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

	Convey("TestRotateFileWriter", t, func() {
		writer, err := NewWriterWithOptions(&WriterOptions{
			Type: "RotateFile",
			Options: &RotateFileWriterOptions{
				Level: "Info",
				Formatter: FormatterOptions{
					Type: "Json",
				},
				MaxAge:   24 * time.Hour,
				Filename: "log/test.log",
			},
		})

		So(err, ShouldBeNil)
		So(writer.Write(info), ShouldBeNil)
		So(writer.Close(), ShouldBeNil)
	})
}
