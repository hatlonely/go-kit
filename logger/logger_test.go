package logger

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	w, err := NewRotateFileWriter("hello.info", 24*time.Hour)
	if err != nil {
		panic(err)
	}
	l := NewLogger(0, w)
	l.Info(map[string]string{
		"hello": "world",
		"key1":  "key2",
	})
}
