package strex

import (
	"encoding/json"
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
