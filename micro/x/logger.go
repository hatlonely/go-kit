package microx

type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}
