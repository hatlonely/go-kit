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

func NewStdoutTextLogger() *Logger {
	log, err := NewLoggerWithOptions(&Options{
		Level: "Debug",
		Writers: []WriterOptions{{
			Type: "Stdout",
			Options: &StdoutWriterOptions{
				Formatter: FormatterOptions{
					Type: "Text",
					Options: &TextFormatOptions{
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

func NewStdoutJsonLogger() *Logger {
	log, err := NewLoggerWithOptions(&Options{
		Level: "Debug",
		Writers: []WriterOptions{{
			Type: "Stdout",
			Options: &StdoutWriterOptions{
				Formatter: FormatterOptions{
					Type: "Json",
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
	return NewLoggerWithOptions(&options, opts...)
}

func NewLoggerWithOptions(options *Options, opts ...refx.Option) (*Logger, error) {
	var writers []Writer
	for _, writerOpt := range options.Writers {
		writer, err := NewWriterWithOptions(&writerOpt, opts...)
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
		flatMap: options.FlatMap,
	}, nil
}

type Options struct {
	Level   string `rule:"x in ['Info', 'Warn', 'Debug', 'Error']"`
	FlatMap bool
	Writers []WriterOptions
}

func NewLogger(level Level, writers ...Writer) *Logger {
	return &Logger{
		level:          level,
		writers:        writers,
		withValues:     map[string]interface{}{},
		withValueFuncs: map[string]func() interface{}{},
	}
}

type Logger struct {
	writers []Writer
	level   Level
	flatMap bool

	withValues     map[string]interface{}
	withValueFuncs map[string]func() interface{}
}

func (l *Logger) With(key string, val interface{}) *Logger {
	log := &Logger{
		writers:        l.writers,
		level:          l.level,
		flatMap:        l.flatMap,
		withValues:     map[string]interface{}{},
		withValueFuncs: map[string]func() interface{}{},
	}

	for key, val := range l.withValues {
		log.withValues[key] = val
	}
	for key, val := range l.withValueFuncs {
		log.withValueFuncs[key] = val
	}

	log.withValues[key] = val
	return log
}

func (l *Logger) WithFunc(key string, val func() interface{}) *Logger {
	log := &Logger{
		writers:        l.writers,
		level:          l.level,
		flatMap:        l.flatMap,
		withValues:     map[string]interface{}{},
		withValueFuncs: map[string]func() interface{}{},
	}

	for key, val := range l.withValues {
		log.withValues[key] = val
	}
	for key, val := range l.withValueFuncs {
		log.withValueFuncs[key] = val
	}

	log.withValueFuncs[key] = val
	return log
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
	kvs := map[string]interface{}{
		"timestamp": now.Unix(),
		"time":      now.Format(time.RFC3339Nano),
		"level":     level.String(),
		"file":      fmt.Sprintf("%s:%v", path.Base(file), line),
		"caller":    fun,
	}
	if l.flatMap {
		for key, val := range v.(map[string]interface{}) {
			kvs[key] = val
		}
	} else {
		kvs["data"] = v
	}
	for key, val := range l.withValues {
		kvs[key] = val
	}
	for key, val := range l.withValueFuncs {
		kvs[key] = val()
	}

	for _, writer := range l.writers {
		_ = writer.Write(kvs)
	}
}

func (l *Logger) Close() error {
	var err error
	for _, w := range l.writers {
		err = w.Close()
	}
	return err
}
