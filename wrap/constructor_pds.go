package wrap

import (
	"time"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"

	pds "github.com/alibabacloud-go/pds-sdk/client"
)

type PDSClientWrapperOptions struct {
	Retry              micro.RetryOptions
	Wrapper            WrapperOptions
	PDS                pds.Config
	RateLimiter        micro.RateLimiterOptions
	ParallelController micro.ParallelControllerOptions
}

func NewPDSClientWrapperWithOptions(options *PDSClientWrapperOptions, opts ...refx.Option) (*PDSClientWrapper, error) {
	var w PDSClientWrapper
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

	if options.PDS.AccessKeyId != nil && *options.PDS.AccessKeyId != "" {
		client, err := pds.NewClient(&options.PDS)
		if err != nil {
			return nil, errors.Wrap(err, "pds.NewClient failed")
		}
		if _, err := client.ListDrives(new(pds.ListDriveRequest)); err != nil {
			return nil, errors.Wrap(err, "pds.Client.ListDrives failed")
		}
		w.obj = client
	} else {
		res, err := alics.ECSMetaDataRamSecurityCredentials()
		if err != nil {
			return nil, errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentials failed")
		}
		options.PDS.SetAccessKeyId(res.AccessKeyID).SetAccessKeySecret(res.AccessKeySecret).SetSecurityToken(res.SecurityToken)
		client, err := pds.NewClient(&options.PDS)
		if err != nil {
			return nil, errors.Wrap(err, "pds.NewClient failed")
		}
		if _, err := client.ListDrives(new(pds.ListDriveRequest)); err != nil {
			return nil, errors.Wrap(err, "pds.Client.ListDrives failed")
		}
		w.obj = client
		go w.UpdateCredentialByECSRole(res, &options.PDS)
	}

	return &w, nil
}

func NewPDSClientWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*PDSClientWrapper, error) {
	var options PDSClientWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewPDSClientWrapperWithOptions(&options, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "NewPDSClientWrapperWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Wrapper"), w.OnWrapperChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Retry"), w.OnRetryChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiter"), w.OnRateLimiterChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("ParallelController"), w.OnParallelControllerChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("PDS"), func(cfg *config.Config) error {
		var options pds.Config
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}

		if options.AccessKeyId != nil && *options.AccessKeyId != "" {
			client, err := pds.NewClient(&options)
			if err != nil {
				return errors.Wrap(err, "pds.NewClient failed")
			}
			if _, err := client.ListDrives(new(pds.ListDriveRequest)); err != nil {
				return errors.Wrap(err, "pds.Client.ListDrives failed")
			}
			w.obj = client
			return nil
		}

		res, err := alics.ECSMetaDataRamSecurityCredentials()
		if err != nil {
			return errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentials failed")
		}
		options.SetAccessKeyId(res.AccessKeyID).SetAccessKeySecret(res.AccessKeySecret).SetSecurityToken(res.SecurityToken)
		client, err := pds.NewClient(&options)
		if err != nil {
			return errors.Wrap(err, "pds.NewClient failed")
		}
		if _, err := client.ListDrives(new(pds.ListDriveRequest)); err != nil {
			return errors.Wrap(err, "pds.Client.ListDrives failed")
		}
		w.obj = client
		go w.UpdateCredentialByECSRole(res, &options)
		return nil
	})

	return w, err
}

func (w *PDSClientWrapper) UpdateCredentialByECSRole(res *alics.ECSMetaDataRamSecurityCredentialsRes, options *pds.Config) {
	for {
		oldClient := w.obj

		d := res.ExpirationTime.Sub(time.Now()) - 25*time.Minute
		if d < 0 {
			d = 30 * time.Second
		}
		<-time.After(d)
		// 检查 client 是否被其他协程更新
		if w.obj != oldClient {
			break
		}

		res, err := alics.ECSMetaDataRamSecurityCredentials()
		if err != nil {
			log.Errorf("alics.ECSMetaDataRamSecurityCredentials failed. err: [%+v]", err)
			continue
		}
		options.SetAccessKeyId(res.AccessKeyID).SetAccessKeySecret(res.AccessKeySecret).SetSecurityToken(res.SecurityToken)
		client, err := pds.NewClient(options)
		if err != nil {
			log.Errorf("pds.NewClient failed. err: [%+v]", err)
			continue
		}
		if _, err := client.ListDrives(new(pds.ListDriveRequest)); err != nil {
			log.Errorf("pds.Client.ListDrives failed. err: [%+v]", err)
			continue
		}
		w.obj = client

		// 检查 client 是否被其他协程更新
		if w.obj != oldClient {
			break
		}
		w.obj = client
	}
}
