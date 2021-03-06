package logger

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJsonFormatter(t *testing.T) {
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

	Convey("case default", t, func() {
		formatter := NewJsonFormatterWithOptions(&JsonFormatterOptions{})
		buf, err := formatter.Format(info)
		So(err, ShouldBeNil)

		var v interface{}
		So(json.Unmarshal(buf, &v), ShouldBeNil)
		So(v, ShouldResemble, map[string]interface{}{
			"caller": "TextFormatter",
			"data": map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
			},
			"file":      "formatter_test.go:18",
			"level":     "Info",
			"meta1":     "val1",
			"meta2":     "val2",
			"time":      "2020-11-07T16:07:36.65798+08:00",
			"timestamp": 1.604736456e+09,
		})
	})

	Convey("case pascal name", t, func() {
		formatter := NewJsonFormatterWithOptions(&JsonFormatterOptions{
			PascalNameKey: true,
		})
		buf, err := formatter.Format(info)
		So(err, ShouldBeNil)

		var v interface{}
		So(json.Unmarshal(buf, &v), ShouldBeNil)
		So(v, ShouldResemble, map[string]interface{}{
			"Caller": "TextFormatter",
			"Data": map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
			},
			"File":      "formatter_test.go:18",
			"Level":     "Info",
			"meta1":     "val1",
			"meta2":     "val2",
			"Time":      "2020-11-07T16:07:36.65798+08:00",
			"Timestamp": 1.604736456e+09,
		})
	})

	Convey("case flat map", t, func() {
		formatter := NewJsonFormatterWithOptions(&JsonFormatterOptions{
			FlatMap: true,
		})
		buf, err := formatter.Format(info)
		So(err, ShouldBeNil)

		var v interface{}
		So(json.Unmarshal(buf, &v), ShouldBeNil)
		So(v, ShouldResemble, map[string]interface{}{
			"caller":    "TextFormatter",
			"key1":      "val1",
			"key2":      "val2",
			"file":      "formatter_test.go:18",
			"level":     "Info",
			"meta1":     "val1",
			"meta2":     "val2",
			"time":      "2020-11-07T16:07:36.65798+08:00",
			"timestamp": 1.604736456e+09,
		})
	})
}
