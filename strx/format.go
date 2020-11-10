package strx

import (
	"encoding/json"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// encoding/json not support map[interface{}]interface{}
func JsonMarshal(v interface{}) string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(v)
	return string(buf)
}

// json-iterator not support prefix
func JsonMarshalIndent(v interface{}) string {
	buf, _ := json.MarshalIndent(v, "  ", "  ")
	return string(buf)
}

func MustJsonMarshal(v interface{}) string {
	buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(v)
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
