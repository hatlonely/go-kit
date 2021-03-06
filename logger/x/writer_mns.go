package loggerx

import (
	"github.com/hatlonely/go-kit/logger"
)

type MNSWriterOptions struct {
	Level      string
	MsgChanLen int `dft:"200"`
	WorkerNum  int `dft:"1"`
}

func NewMNSWriterWithOptions(options *MNSWriterOptions) *MNSWriter {
	return &MNSWriter{}
}

type MNSWriter struct {
	options *MNSWriterOptions
}

func (w *MNSWriter) send(info *logger.Info) error {
	return nil
}
