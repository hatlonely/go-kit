package logger

import (
	"bytes"
	"encoding/json"
	"os"
)

func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{}
}

type StdoutWriter struct{}

func (r *StdoutWriter) Write(v interface{}) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(buf)
	b.WriteString("\n")
	_, err = os.Stdout.Write(b.Bytes())
	if err != nil {
		return err
	}
	return nil
}
