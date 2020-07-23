package logger

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"

	"github.com/hatlonely/go-kit/strex"
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
	log.AddHook(&CallerHook{})
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
	l.log.Info(strex.MustJsonMarshal(v))
}

func (l *LogrusLogger) Warn(v interface{}) {
	l.log.Info(strex.MustJsonMarshal(v))
}

type CallerHook struct{}

func (hook *CallerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *CallerHook) Fire(entry *logrus.Entry) error {
	for i := 5; i < 20; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if strings.Contains(file, "logrus") {
			continue
		}
		funcName := runtime.FuncForPC(pc).Name()
		entry.Data["@file"] = fmt.Sprintf("%s:%v:%s", path.Base(file), line, path.Base(funcName))
		break
	}

	return nil
}

type LogrusFormatter struct{}

func (f *LogrusFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	buf := []byte(entry.Message)
	return append(buf, '\n'), nil
}
