package logger

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

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

//go:generate stringer -type Level -trimprefix Level -case=toUpper
type Level int

const (
	LevelDebug Level = 1
	LevelInfo  Level = 2
	LevelWarn  Level = 3
	LevelError Level = 4
)

func (l *Logger) Info(v interface{}) {
	l.Log(LevelInfo, v)
}

func (l *Logger) Warn(v interface{}) {
	l.Log(LevelWarn, v)
}

func (l *Logger) Log(level Level, v interface{}) {
	if level < l.level {
		return
	}
	pc, file, line, _ := runtime.Caller(2)
	fun := runtime.FuncForPC(pc).Name()

	now := time.Now()

	for _, writer := range l.writers {
		writer.Write(map[string]interface{}{
			"timestamp": now.Unix(),
			"time":      time.Now().Format(time.RFC3339Nano),
			"level":     level,
			"data":      v,
			"file":      fmt.Sprintf("%s:%v", path.Base(file), line),
			"caller":    fun,
		})
	}
}
