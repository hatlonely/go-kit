package logger

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStdoutWriter(t *testing.T) {
	Convey("TestStdoutWriter", t, func() {
		//writer, err := NewWriterWithOptions(&WriterOptions{
		//	Type: "Stdout",
		//	Options: &StdoutWriterOptions{
		//		Formatter: FormatterOptions{
		//			Type: "Json",
		//		},
		//	},
		//})

		writer, err := NewStdoutWriterWithOptions(&StdoutWriterOptions{
			Formatter: FormatterOptions{
				Type:    "Json",
				Options: &JsonFormatterOptions{},
			},
		})
		So(err, ShouldBeNil)

		t, _ := time.Parse(time.RFC3339Nano, "2020-11-07T16:07:36.65798+08:00")
		_ = writer.Write(&Info{
			Time:  t,
			Level: LevelInfo,
			Data: map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
			},
			File:   "formatter_test.go:18",
			Caller: "TextFormatter",
		})
	})
}
