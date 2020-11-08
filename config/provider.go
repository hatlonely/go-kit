package config

import (
	"context"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

type Provider interface {
	Events() <-chan struct{}
	Errors() <-chan error
	Load() ([]byte, error)
	Dump(buf []byte) error
	EventLoop(ctx context.Context) error
}

func NewProviderWithConfig(cfg *Config, opts ...refx.Option) (Provider, error) {
	var options ProviderOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "cfg.Unmarshal failed.")
	}
	return NewProviderWithOptions(&options)
}

func NewProviderWithOptions(options *ProviderOptions) (Provider, error) {
	switch options.Type {
	case "Local":
		return NewLocalProviderWithOptions(&options.LocalProvider)
	case "OTS":
		return NewOTSProviderWithOptions(&options.OTSProvider)
	case "OTSLegacy":
		return NewOTSLegacyProviderWithOptions(&options.OTSLegacyProvider)
	}
	return nil, errors.Errorf("unsupported provider type [%v]", options.Type)
}

type ProviderOptions struct {
	Type              string
	LocalProvider     LocalProviderOptions
	OTSProvider       OTSProviderOptions
	OTSLegacyProvider OTSLegacyProviderOptions
}
