package logger

import (
	"fmt"
	"path"
	"runtime"
	"time"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

func NewStdoutLogger() *Logger {
	log, err := NewLoggerWithOptions(&Options{
		Level: "Debug",
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
		}},
	})
	if err != nil {
		panic(err)
	}
	return log
}

func NewLoggerWithConfig(cfg *config.Config, opts ...refx.Option) (*Logger, error) {
	var options Options
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "cfg.Unmarshal failed.")
	}
	return NewLoggerWithOptions(&options)
}

func NewLoggerWithOptions(options *Options) (*Logger, error) {
	var writers []Writer
	for _, writerOpt := range options.Writers {
		writer, err := NewWriterWithOptions(&writerOpt)
		if err != nil {
			return nil, errors.WithMessage(err, "NewWriterWithOptions failed")
		}
		writers = append(writers, writer)
	}
	level, err := LevelString(options.Level)
	if err != nil {
		return nil, errors.WithMessage(err, "LevelToString failed")
	}
	return &Logger{
		level:   level,
		writers: writers,
	}, nil
}

type Options struct {
	Level   string `rule:"x in ['Info', 'Warn', 'Debug', 'Error']"`
	Writers []WriterOptions
}

func NewLogger(level Level, writers ...Writer) *Logger {
	return &Logger{
		level:   level,
		writers: writers,
	}
}

type Logger struct {
	writers []Writer
	level   Level
}

//go:generate enumer -type Level -trimprefix Level -text
type Level int

const (
	LevelDebug Level = 1
	LevelInfo  Level = 2
	LevelWarn  Level = 3
	LevelError Level = 4
)

func (l *Logger) Debug(v interface{}) {
	l.Log(LevelDebug, v)
}

func (l *Logger) Info(v interface{}) {
	l.Log(LevelInfo, v)
}

func (l *Logger) Warn(v interface{}) {
	l.Log(LevelWarn, v)
}

func (l *Logger) Error(v interface{}) {
	l.Log(LevelError, v)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logf(LevelDebug, format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logf(LevelInfo, format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logf(LevelWarn, format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logf(LevelError, format, args...)
}

func (l *Logger) Logf(level Level, format string, args ...interface{}) {
	l.Log(level, fmt.Sprintf(format, args...))
}

func (l *Logger) Log(level Level, v interface{}) {
	if level < l.level {
		return
	}
	pc, file, line, _ := runtime.Caller(2)
	fun := runtime.FuncForPC(pc).Name()

	now := time.Now()
	for _, writer := range l.writers {
		_ = writer.Write(map[string]interface{}{
			"timestamp": now.Unix(),
			"time":      now.Format(time.RFC3339Nano),
			"level":     level.String(),
			"data":      v,
			"file":      fmt.Sprintf("%s:%v", path.Base(file), line),
			"caller":    fun,
		})
	}
}
