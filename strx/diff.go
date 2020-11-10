package strx

import (
	"encoding/json"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

func Diff(text1, text2 string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, true)
	return dmp.DiffPrettyText(diffs)
}

func JsonDiff(text1, text2 string) string {
	differ := gojsondiff.New()
	diff, err := differ.Compare([]byte(text1), []byte(text2))
	if err != nil {
		return err.Error()
	}
	var v interface{}
	if err := json.Unmarshal([]byte(text1), &v); err != nil {
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
