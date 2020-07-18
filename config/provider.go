package config

import (
	"context"
	"fmt"
)

type Provider interface {
	Events() <-chan struct{}
	Errors() <-chan error
	Load() ([]byte, error)
	Dump(buf []byte) error
	EventLoop(ctx context.Context) error
}

func NewProviderWithConfig(conf *Config) (Provider, error) {
	switch conf.GetString("Type") {
	case "Local":
		return NewLocalProvider(conf.GetString("File"))
	case "OTS":
		return NewOTSProviderWithAccessKey(
			conf.GetString("AccessKeyID"),
			conf.GetString("AccessKeySecret"),
			conf.GetString("Endpoint"),
			conf.GetString("Instance"),
			conf.GetString("Table"),
			conf.GetString("Key"),
			conf.GetDuration("Interval"),
		)
	}
	return nil, fmt.Errorf("unsupport provider type. type: [%v]", conf.GetString("Type"))
}
