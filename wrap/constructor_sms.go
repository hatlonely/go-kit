package wrap

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

type SMSOptions struct {
	RegionID        string
	AccessKeyID     string
	AccessKeySecret string
}

type SMSClientWrapperOptions struct {
	Retry              micro.RetryOptions
	Wrapper            WrapperOptions
	SMS                SMSOptions
	RateLimiter        micro.RateLimiterOptions
	ParallelController micro.ParallelControllerOptions
}

func NewSMSClientWithOptions(options *SMSOptions) (*dysmsapi.Client, error) {
	if options.AccessKeyID != "" {
		client, err := dysmsapi.NewClientWithAccessKey(options.RegionID, options.AccessKeyID, options.AccessKeySecret)
		if err != nil {
			return nil, errors.Wrap(err, "dysmsapi.NewClientWithAccessKey failed")
		}
		if _, err := client.QuerySmsTemplate(dysmsapi.CreateQuerySmsTemplateRequest()); err != nil {
			return nil, errors.Wrap(err, "dysmsapi.Client.QuerySmsTemplate failed")
		}
		return client, nil
	}

	role, err := alics.ECSMetaDataRamSecurityCredentialsRole()
	if err != nil {
		return nil, errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentialsRole failed")
	}
	client, err := dysmsapi.NewClientWithEcsRamRole(options.RegionID, role)
	if err != nil {
		return nil, errors.Wrap(err, "dysmsapi.NewClientWithEcsRamRole failed")
	}
	if _, err := client.QuerySmsTemplate(dysmsapi.CreateQuerySmsTemplateRequest()); err != nil {
		return nil, errors.Wrap(err, "dysmsapi.Client.QuerySmsTemplate failed")
	}
	return client, nil
}

func NewSMSClientWrapperWithOptions(options *SMSClientWrapperOptions, opts ...refx.Option) (*SMSClientWrapper, error) {
	var w SMSClientWrapper
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
	w.obj, err = NewSMSClientWithOptions(&options.SMS)
	if err != nil {
		return nil, errors.WithMessage(err, "NewSMSClientWithOptions failed")
	}
	return &w, nil
}

func NewSMSClientWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*SMSClientWrapper, error) {
	var options SMSClientWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewSMSClientWrapperWithOptions(&options, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "NewSMSClientWrapperWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Wrapper"), w.OnWrapperChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Retry"), w.OnRetryChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiterGroup"), w.OnRateLimiterChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("ParallelController"), w.OnParallelControllerChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("SMS"), func(cfg *config.Config) error {
		var options SMSOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		client, err := NewSMSClientWithOptions(&options)
		if err != nil {
			return errors.WithMessage(err, "NewSMSClientWithOptions failed")
		}
		w.obj = client
		return nil
	})

	return w, err
}
