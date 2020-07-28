package logger

type Writer interface {
	Write(v interface{}) error
}
