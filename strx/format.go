package strx

import (
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func JsonMarshal(v interface{}) string {
	buf, _ := json.Marshal(v)
	return string(buf)
}

func JsonMarshalIndent(v interface{}) string {
	buf, _ := json.MarshalIndent(v, "  ", "  ")
	return string(buf)
}

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
