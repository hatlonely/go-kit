package logger

import (
	"sync"
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

func TestParallel(t *testing.T) {
	w, err := NewRotateFileWriter("hello.info", 24*time.Hour)
	if err != nil {
		panic(err)
	}
	l := NewLogger(0, w)

	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			for time.Now().Before(start.Add(10 * time.Second)) {
				l.Info(map[string]string{
					"hello": "world",
					"key1":  "key2",
				})
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkHello(b *testing.B) {
	w, err := NewRotateFileWriter("hello.info", 24*time.Hour)
	if err != nil {
		panic(err)
	}
	l := NewLogger(0, w)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info(map[string]string{
				"hello": "world",
				"key1":  "key2",
			})
		}
	})
}
