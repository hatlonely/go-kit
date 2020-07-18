package strex

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

func Diff(text1, text2 string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, true)
	return dmp.DiffPrettyText(diffs)
}
