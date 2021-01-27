package wrap

import (
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type OTSOptions struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	InstanceName    string
}

type OTSTableStoreClientWrapperOptions struct {
	Retry            RetryOptions
	Wrapper          WrapperOptions
	OTS              OTSOptions
	RateLimiterGroup RateLimiterGroupOptions
}

func NewOTSTableStoreClientWrapperWithOptions(options *OTSTableStoreClientWrapperOptions) (*OTSTableStoreClientWrapper, error) {
	retry, err := NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}

	rateLimiter, err := NewRateLimiterGroup(&options.RateLimiterGroup)
	if err != nil {
		return nil, errors.Wrap(err, "NewRateLimiterGroup failed")
	}

	if options.OTS.AccessKeyID != "" {
		client := tablestore.NewClient(options.OTS.Endpoint, options.OTS.InstanceName, options.OTS.AccessKeyID, options.OTS.AccessKeySecret)
		if _, err := client.ListTable(); err != nil {
			return nil, errors.Wrap(err, "tablestore.TableStoreClient.ListTable failed")
		}
		return &OTSTableStoreClientWrapper{
			obj:     client,
			retry:   retry,
			options: &options.Wrapper,
		}, nil
	}

	res, err := alics.ECSMetaDataRamSecurityCredentials()
	if err != nil {
		return nil, errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentials failed")
	}
	client := tablestore.NewClientWithConfig(options.OTS.Endpoint, options.OTS.InstanceName, res.AccessKeyID, res.AccessKeySecret, res.SecurityToken, nil)
	if _, err := client.ListTable(); err != nil {
		return nil, errors.Wrap(err, "tablestore.TableStoreClient.ListTable failed")
	}

	w := &OTSTableStoreClientWrapper{
		obj:              client,
		retry:            retry,
		options:          &options.Wrapper,
		rateLimiterGroup: rateLimiter,
	}

	if w.options.EnableMetric {
		w.CreateMetric(w.options)
	}

	go w.UpdateCredentialByECSRole(res, &options.OTS)

	return w, nil
}

func NewOTSTableStoreClientWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*OTSTableStoreClientWrapper, error) {
	var options OTSTableStoreClientWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewOTSTableStoreClientWrapperWithOptions(&options)
	if err != nil {
		return nil, errors.Wrap(err, "NewOTSTableStoreClientWrapperWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Wrapper"), w.OnWrapperChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Retry"), w.OnRetryChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiterGroup"), w.OnRateLimiterGroupChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("OTS"), func(cfg *config.Config) error {
		var options OTSOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}

		if options.AccessKeyID != "" {
			client := tablestore.NewClient(options.Endpoint, options.InstanceName, options.AccessKeyID, options.AccessKeySecret)
			if _, err := client.ListTable(); err != nil {
				return errors.Wrap(err, "tablestore.TableStoreClient.ListTable failed")
			}
			return nil
		}

		res, err := alics.ECSMetaDataRamSecurityCredentials()
		if err != nil {
			return errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentials failed")
		}
		client := tablestore.NewClientWithConfig(options.Endpoint, options.InstanceName, res.AccessKeyID, res.AccessKeySecret, res.SecurityToken, nil)
		if _, err := client.ListTable(); err != nil {
			return errors.Wrap(err, "tablestore.TableStoreClient.ListTable failed")
		}
		w.obj = client
		go w.UpdateCredentialByECSRole(res, &options)

		return nil
	})

	return w, err
}

func (w *OTSTableStoreClientWrapper) UpdateCredentialByECSRole(res *alics.ECSMetaDataRamSecurityCredentialsRes, options *OTSOptions) {
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
		client := tablestore.NewClientWithConfig(options.Endpoint, options.InstanceName, res.AccessKeyID, res.AccessKeySecret, res.SecurityToken, nil)
		if _, err := client.ListTable(); err != nil {
			log.Errorf("tablestore.TableStoreClient.ListTable failed. err: [%+v]", err)
			continue
		}

		// 检查 client 是否被其他协程更新
		if w.obj != oldClient {
			break
		}
		w.obj = client
	}
}
