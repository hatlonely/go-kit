package strx

import (
	"encoding/json"
	"reflect"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

// TODO 1.显示不太直观 2.省略相同的行
func Diff(text1, text2 string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, true)
	return dmp.DiffPrettyText(diffs)
}

// TODO 这个库不太友好，如果 text2 为空，显示不出 diff
func JsonDiff(text1, text2 string) string {
	if text1 == "null" {
		text1 = "{}"
	}
	if text2 == "null" {
		text2 = "{}"
	}
	var v interface{}
	if err := json.Unmarshal([]byte(text1), &v); err != nil {
		return Diff(text1, text2)
	} else if reflect.TypeOf(v).Kind() != reflect.Map && reflect.TypeOf(v).Kind() != reflect.Slice {
		return Diff(text1, text2)
	}
	if err := json.Unmarshal([]byte(text2), &v); err != nil {
		return Diff(text1, text2)
	} else if reflect.TypeOf(v).Kind() != reflect.Map && reflect.TypeOf(v).Kind() != reflect.Slice {
		return Diff(text1, text2)
	}

	differ := gojsondiff.New()
	diff, err := differ.Compare([]byte(text1), []byte(text2))
	if err != nil {
		return err.Error()
	}

	f := formatter.NewAsciiFormatter(v, formatter.AsciiFormatterConfig{
		ShowArrayIndex: true,
		Coloring:       true,
	})
	deltaJson, err := f.Format(diff)
	if err != nil {
		return err.Error()
	}
	return deltaJson
}
