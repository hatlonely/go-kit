package logger

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func NewJsonFormatter() JsonFormatter {
	return JsonFormatter{}
}

type JsonFormatter struct{}

func (f JsonFormatter) Format(kvs map[string]interface{}) ([]byte, error) {
	return json.Marshal(kvs)
}
