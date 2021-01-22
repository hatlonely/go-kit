package wrap

import (
	"github.com/hatlonely/go-kit/logger"
)

var log Logger

func init() {
	log = logger.NewStdoutTextLogger()
}

type Logger interface {
	Errorf(format string, args ...interface{})
}

func SetLogger(logger Logger) {
	log = logger
}
