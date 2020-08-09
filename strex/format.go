package strex

import (
	"encoding/json"
	"regexp"
	"strings"
)

func MustJsonMarshal(v interface{}) string {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func MustJsonMarshalIndent(v interface{}) string {
	buf, err := json.MarshalIndent(v, "  ", "  ")
	if err != nil {
		panic(err)
	}
	return string(buf)
}

var spacePattern = regexp.MustCompile(`\s+`)

func FormatSpace(str string) string {
	return strings.TrimSpace(spacePattern.ReplaceAllString(str, " "))
}
