package logger

import (
	"testing"
	"time"
)

func hello(l Logger) {

}

func TestLogrusLogger(t *testing.T) {
	l, err := NewLogrusLogger("hello.info", 24*time.Hour)
	if err != nil {
		panic(err)
	}
	l.Info(map[string]string{
		"hello": "world",
		"key1":  "key2",
	})
}
