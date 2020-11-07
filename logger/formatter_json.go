package logger

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type JsonFormatter struct{}

func (f JsonFormatter) Format(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
