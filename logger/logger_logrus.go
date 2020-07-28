package logger

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	log *logrus.Logger
}

func NewLogrusLogger(filename string, maxAge time.Duration) (*LogrusLogger, error) {
	log := logrus.New()

	abs, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	log.Formatter = &LogrusFormatter{}

	out, err := rotatelogs.New(
		abs+".%Y%m%d%H",
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithLinkName(abs),
		rotatelogs.WithMaxAge(maxAge),
	)
	if err != nil {
		return nil, err
	}
	log.Out = out

	return &LogrusLogger{
		log: log,
	}, nil
}

func (l *LogrusLogger) Info(v interface{}) {
	l.Log(logrus.InfoLevel, v)
}

func (l *LogrusLogger) Warn(v interface{}) {
	l.Log(logrus.WarnLevel, v)
}

func (l *LogrusLogger) Log(level logrus.Level, v interface{}) {
	pc, file, line, _ := runtime.Caller(2)
	fun := runtime.FuncForPC(pc).Name()

	now := time.Now()
	l.log.Log(level, map[string]interface{}{
		"timestamp": now.Unix(),
		"time":      time.Now().Format(time.RFC3339Nano),
		"level":     level,
		"data":      v,
		"file":      fmt.Sprintf("%s:%v", path.Base(file), line),
		"caller":    fun,
	})
}

type LogrusFormatter struct{}

func (f *LogrusFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	buf := []byte(entry.Message)
	return append(buf, '\n'), nil
}
