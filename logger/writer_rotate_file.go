package logger

import (
	"encoding/json"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type RotateFileWriter struct {
	out *rotatelogs.RotateLogs
}

func NewRotateFileWriter(filename string, maxAge time.Duration) (*RotateFileWriter, error) {
	abs, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	out, err := rotatelogs.New(
		abs+".%Y%m%d%H",
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithLinkName(abs),
		rotatelogs.WithMaxAge(maxAge),
	)
	if err != nil {
		return nil, err
	}

	return &RotateFileWriter{
		out: out,
	}, nil
}

func (r *RotateFileWriter) Write(v interface{}) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = r.out.Write(buf)
	if err != nil {
		return err
	}
	_, err = r.out.Write([]byte("\n"))
	if err != nil {
		return err
	}
	return nil
}
