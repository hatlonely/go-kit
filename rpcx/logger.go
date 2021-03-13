package rpcx

type Logger interface {
	Info(v interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}
