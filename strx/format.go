package strx

import (
	"encoding/json"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var sortJson = jsoniter.ConfigCompatibleWithStandardLibrary

func JsonMarshal(v interface{}) string {
	buf, _ := jsoniter.Marshal(v)
	return string(buf)
}

func JsonMarshalIndent(v interface{}) string {
	buf, _ := jsoniter.MarshalIndent(v, "", "  ")
	return string(buf)
}

func JsonMarshalSortKeys(v interface{}) string {
	buf, _ := sortJson.Marshal(v)
	return string(buf)
}

func JsonMarshalIndentSortKeys(v interface{}) string {
	// encoding/json 无法 marshal map[interface{}]interface{}
	if buf, err := json.MarshalIndent(v, "", "  "); err == nil {
		return string(buf)
	}
	// jsoniter marshal indent 在 sortKeys indent 格式有问题
	buf, _ := sortJson.MarshalIndent(v, "", "  ")
	return string(buf)
}

var spacePattern = regexp.MustCompile(`\s+`)

func FormatSpace(str string) string {
	return strings.TrimSpace(spacePattern.ReplaceAllString(str, " "))
}
