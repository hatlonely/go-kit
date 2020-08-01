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
	switch conf.GetString("type") {
	case "Local":
		return NewLocalProvider(conf.GetString("file"))
	case "OTS":
		return NewOTSProviderWithAccessKey(
			conf.GetString("accessKeyID"),
			conf.GetString("accessKeySecret"),
			conf.GetString("endpoint"),
			conf.GetString("instance"),
			conf.GetString("table"),
			conf.GetString("key"),
			conf.GetDuration("interval"),
		)
	}
	return nil, fmt.Errorf("unsupport provider type. type: [%v]", conf.GetString("Type"))
}
