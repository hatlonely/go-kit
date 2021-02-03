package rpcx

type Logger interface {
	Info(v interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
}
