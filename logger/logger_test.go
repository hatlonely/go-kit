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
	l := NewLogger(0, w, NewStdoutWriter())
	l.Info(map[string]string{
		"hello": "world",
		"key1":  "key2",
	})
}

func BenchmarkHello(b *testing.B) {
	w, err := NewRotateFileWriter("hello.info", 24*time.Hour)
	if err != nil {
		panic(err)
	}
	l := NewLogger(0, w, NewStdoutWriter())

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info(map[string]string{
				"hello": "world",
				"key1":  "key2",
			})
		}
	})
}
