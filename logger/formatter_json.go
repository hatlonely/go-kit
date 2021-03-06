package logger

import (
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/hatlonely/go-kit/strx"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type JsonFormatterOptions struct {
	FlatMap       bool
	PascalNameKey bool
}

func NewJsonFormatterWithOptions(options *JsonFormatterOptions) *JsonFormatter {
	timeKey := "time"
	timestampKey := "timestamp"
	levelKey := "level"
	fileKey := "file"
	callerKey := "caller"
	dataKey := "data"

	if options.PascalNameKey {
		for _, key := range []*string{&timeKey, &timestampKey, &levelKey, &fileKey, &callerKey, &dataKey} {
			*key = strx.PascalName(*key)
		}
	}

	return &JsonFormatter{
		options:      options,
		timeKey:      timeKey,
		timestampKey: timestampKey,
		levelKey:     levelKey,
		fileKey:      fileKey,
		callerKey:    callerKey,
		dataKey:      dataKey,
	}
}

type JsonFormatter struct {
	options      *JsonFormatterOptions
	timeKey      string
	timestampKey string
	levelKey     string
	fileKey      string
	callerKey    string
	dataKey      string
}

func (f JsonFormatter) Format(info *Info) ([]byte, error) {
	kvs := map[string]interface{}{
		f.timestampKey: info.Time.Unix(),
		f.timeKey:      info.Time.Format(time.RFC3339Nano),
		f.levelKey:     info.Level.String(),
		f.fileKey:      info.File,
		f.callerKey:    info.Caller,
	}

	for key, val := range info.Fields {
		kvs[key] = val
	}

	if f.options.FlatMap {
		for key, val := range info.Data.(map[string]interface{}) {
			kvs[key] = val
		}
	} else {
		kvs[f.dataKey] = info.Data
	}

	return json.Marshal(kvs)
}
