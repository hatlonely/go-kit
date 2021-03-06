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
						Template: "{{.time}} [{{.level}}] [{{.caller}}:{{.file}}] {{.data}}",
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
		return nil, errors.Wrap(err, "LevelToString failed")
	}
	return &Logger{
		level:   level,
		writers: writers,
	}, nil
}

type Options struct {
	Level   string `rule:"x in ['Info', 'Warn', 'Debug', 'Error', 'Fatal']"`
	Writers []WriterOptions
}

func NewLogger(level Level, writers ...Writer) *Logger {
	return &Logger{
		level:   level,
		writers: writers,
	}
}

type Logger struct {
	parent  *Logger
	writers []Writer
	level   Level
	flatMap bool

	nodeType NodeType
	key      string
	val      interface{}
	fun      func() interface{}
	kvs      map[string]interface{}
}

func (l *Logger) With(key string, val interface{}) *Logger {
	return &Logger{
		parent:   l,
		nodeType: NodeTypeVal,
		level:    l.level,
		flatMap:  l.flatMap,
		key:      key,
		val:      val,
	}
}

func (l *Logger) WithFunc(key string, val func() interface{}) *Logger {
	return &Logger{
		parent:   l,
		nodeType: NodeTypeFunc,
		level:    l.level,
		flatMap:  l.flatMap,
		key:      key,
		fun:      val,
	}
}

func (l *Logger) WithFields(kvs map[string]interface{}) *Logger {
	return &Logger{
		parent:   l,
		nodeType: NodeTypeMap,
		level:    l.level,
		flatMap:  l.flatMap,
		kvs:      kvs,
	}
}

type NodeType int

const (
	NodeTypeVal  NodeType = 1
	NodeTypeFunc NodeType = 2
	NodeTypeMap  NodeType = 3
)

//go:generate enumer -type Level -trimprefix Level -text
type Level int

const (
	LevelDebug Level = 1
	LevelInfo  Level = 2
	LevelWarn  Level = 3
	LevelError Level = 4
	LevelFatal Level = 5
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

func (l *Logger) Fatal(v interface{}) {
	l.Log(LevelFatal, v)
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

type Info struct {
	Time   time.Time
	Level  Level
	File   string
	Caller string
	Fields map[string]interface{}
	Data   interface{}
}

func (l *Logger) Log(level Level, v interface{}) {
	if level < l.level {
		return
	}

	fields := map[string]interface{}{}
	node := l
	for node.parent != nil {
		switch node.nodeType {
		case NodeTypeVal:
			fields[node.key] = node.val
		case NodeTypeFunc:
			fields[node.key] = node.fun()
		case NodeTypeMap:
			for k, v := range node.kvs {
				fields[k] = v
			}
		}
		node = node.parent
	}

	pc, file, line, _ := runtime.Caller(2)
	for _, writer := range node.writers {
		_ = writer.Write(&Info{
			Time:   time.Now(),
			Level:  level,
			File:   fmt.Sprintf("%s:%v", path.Base(file), line),
			Caller: runtime.FuncForPC(pc).Name(),
			Fields: fields,
			Data:   v,
		})
	}
}

func (l *Logger) Close() error {
	var err error
	for _, w := range l.writers {
		err = w.Close()
	}
	return err
}
