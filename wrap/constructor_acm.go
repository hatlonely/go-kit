package wrap

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

type ACMConfigClientWrapperOptions struct {
	Retry              micro.RetryOptions
	Wrapper            WrapperOptions
	ACM                constant.ClientConfig
	RateLimiter        micro.RateLimiterOptions
	ParallelController micro.ParallelControllerOptions
}

func NewACMConfigClientWrapperWithOptions(options *ACMConfigClientWrapperOptions, opts ...refx.Option) (*ACMConfigClientWrapper, error) {
	var w ACMConfigClientWrapper
	var err error

	w.options = &options.Wrapper
	w.retry, err = micro.NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "micro.NewRetryWithOptions failed")
	}
	w.rateLimiter, err = micro.NewRateLimiterWithOptions(&options.RateLimiter, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "micro.NewRateLimiterWithOptions failed")
	}
	w.parallelController, err = micro.NewParallelControllerWithOptions(&options.ParallelController, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "micro.NewParallelControllerWithOptions failed")
	}
	if w.options.EnableMetric {
		w.CreateMetric(w.options)
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": options.ACM,
	})
	if err != nil {
		return nil, errors.Wrap(err, "clients.CreateConfigClient failed")
	}
	w.obj = client.(*config_client.ConfigClient)

	if w.options.EnableMetric {
		w.CreateMetric(w.options)
	}

	return &w, nil
}

func NewACMConfigClientWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*ACMConfigClientWrapper, error) {
	var options ACMConfigClientWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewACMConfigClientWrapperWithOptions(&options, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "NewACMConfigClientWrapperWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Wrapper"), w.OnWrapperChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Retry"), w.OnRetryChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiter"), w.OnRateLimiterChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("ParallelController"), w.OnParallelControllerChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("ACM"), func(cfg *config.Config) error {
		var options constant.ClientConfig
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		client, err := clients.CreateConfigClient(map[string]interface{}{
			"clientConfig": options,
		})
		if err != nil {
			return errors.Wrap(err, "clients.CreateConfigClient failed")
		}
		w.obj = client.(*config_client.ConfigClient)
		return nil
	})

	return w, err
}
