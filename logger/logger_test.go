package logger

import (
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLogger(t *testing.T) {
	Convey("TestLogger", t, func() {
		log, err := NewLoggerWithOptions(&Options{
			Level: "Info",
			Writers: []WriterOptions{{
				Type: "Stdout",
				StdoutWriter: StdoutWriterOptions{
					Formatter: FormatterOptions{
						Type: "Text",
						TextFormat: TextFormatOptions{
							Format: "{{.time}} [{{.level}}] [{{.caller}}:{{.file}}] {{.data}}",
						},
					},
				},
			}, {
				Type: "RotateFile",
				RotateFileWriter: RotateFileWriterOptions{
					MaxAge:   24 * time.Hour,
					Filename: "log/test.log",
					Formatter: FormatterOptions{
						Type: "Json",
					},
				},
			}},
		})
		So(err, ShouldBeNil)
		log.Info("hello world")
		log.Warn(map[string]string{
			"key1": "val1",
			"key2": "val2",
		})
	})
}

func TestNewStdoutLogger(t *testing.T) {
	Convey("TestNewStdoutLogger", t, func() {
		log := NewStdoutJsonLogger()
		log.Info("hello world")
		log.Warn("hello golang")
		log.WithFunc("key1", func() interface{} {
			return time.Now().UnixNano()
		}).Info("")
		log.WithFunc("key2", func() interface{} {
			return time.Now().UnixNano()
		}).Info("")
		log.With("Key1", "Val1").With("Key2", "Val2").Info("")
		log.With("Key3", "Val3").With("Key4", "Val4").Info("")
	})
}

func TestNewFlatMapLogger(t *testing.T) {
	Convey("TestNewFlatMapLogger", t, func() {
		log, _ := NewLoggerWithOptions(&Options{
			Level:   "Debug",
			FlatMap: true,
			Writers: []WriterOptions{{
				Type: "Stdout",
				StdoutWriter: StdoutWriterOptions{
					Formatter: FormatterOptions{
						Type: "Json",
					},
				},
			}},
		})

		log.Info(map[string]interface{}{"key1": "val1", "hello": "world"})
	})
}

func TestParallel(t *testing.T) {
	w, err := NewRotateFileWriter(WithRotateFilename("log/test.log"))
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
					"key1": "val1",
					"key2": "val2",
				})
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkHello(b *testing.B) {
	w, err := NewRotateFileWriter(WithRotateFilename("log/test.log"))
	if err != nil {
		panic(err)
	}
	l := NewLogger(0, w)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info(map[string]string{
				"key1": "val1",
				"key2": "val2",
			})
		}
	})
}
