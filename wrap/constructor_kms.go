package wrap

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

type KMSOptions struct {
	RegionID        string
	AccessKeyID     string
	AccessKeySecret string
}

type KMSClientWrapperOptions struct {
	Retry              micro.RetryOptions
	Wrapper            WrapperOptions
	KMS                KMSOptions
	RateLimiter        micro.RateLimiterOptions
	ParallelController micro.ParallelControllerOptions
}

func NewKMSClientWithOptions(options *KMSOptions) (*kms.Client, error) {
	if options.AccessKeyID != "" {
		client, err := kms.NewClientWithAccessKey(options.RegionID, options.AccessKeyID, options.AccessKeySecret)
		if err != nil {
			return nil, errors.Wrap(err, "kms.NewClientWithAccessKey failed")
		}
		if _, err := client.ListKeys(kms.CreateListKeysRequest()); err != nil {
			return nil, errors.Wrap(err, "kms.Client.ListKeys failed")
		}
		return client, nil
	}

	role, err := alics.ECSMetaDataRamSecurityCredentialsRole()
	if err != nil {
		return nil, errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentialsRole failed")
	}
	client, err := kms.NewClientWithEcsRamRole(options.RegionID, role)
	if err != nil {
		return nil, errors.Wrap(err, "kms.NewClientWithEcsRamRole failed")
	}
	if _, err := client.ListKeys(kms.CreateListKeysRequest()); err != nil {
		return nil, errors.Wrap(err, "kms.Client.ListKeys failed")
	}
	return client, nil
}

func NewKMSClientWrapperWithOptions(options *KMSClientWrapperOptions, opts ...refx.Option) (*KMSClientWrapper, error) {
	var w KMSClientWrapper
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
	w.obj, err = NewKMSClientWithOptions(&options.KMS)
	if err != nil {
		return nil, errors.WithMessage(err, "NewKMSClientWithOptions failed")
	}
	return &w, nil
}

func NewKMSClientWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*KMSClientWrapper, error) {
	var options KMSClientWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewKMSClientWrapperWithOptions(&options, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "NewKMSClientWrapperWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Wrapper"), w.OnWrapperChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Retry"), w.OnRetryChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiterGroup"), w.OnRateLimiterChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("ParallelController"), w.OnParallelControllerChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("KMS"), func(cfg *config.Config) error {
		var options KMSOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		client, err := NewKMSClientWithOptions(&options)
		if err != nil {
			return errors.WithMessage(err, "NewKMSClientWithOptions failed")
		}
		w.obj = client
		return nil
	})

	return w, err
}
