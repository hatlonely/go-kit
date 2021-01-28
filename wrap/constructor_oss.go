package wrap

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

func init() {
	RegisterErrCode(oss.ServiceError{}, func(err error) string {
		e := err.(oss.ServiceError)
		return fmt.Sprintf("oss_%v_%v", e.StatusCode, e.Code)
	})
	RegisterRetryRetryIf("OSS", func(err error) bool {
		switch e := err.(type) {
		case oss.ServiceError:
			if e.StatusCode >= http.StatusInternalServerError {
				return true
			}
			if strings.Contains(e.Error(), "timeout") {
				return true
			}
		}
		return false
	})
}

type OSSOptions struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
}

type OSSClientWrapperOptions struct {
	Retry            RetryOptions
	Wrapper          WrapperOptions
	OSS              OSSOptions
	RateLimiterGroup RateLimiterGroupOptions
}

func NewOSSClientWrapperWithOptions(options *OSSClientWrapperOptions, opts ...refx.Option) (*OSSClientWrapper, error) {
	var w OSSClientWrapper
	var err error

	w.options = &options.Wrapper
	w.retry, err = NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}
	w.rateLimiterGroup, err = NewRateLimiterGroupWithOptions(&options.RateLimiterGroup, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "NewRateLimiterGroupWithOptions failed")
	}
	if w.options.EnableMetric {
		w.CreateMetric(w.options)
	}

	if options.OSS.AccessKeyID != "" {
		client, err := oss.New(options.OSS.Endpoint, options.OSS.AccessKeyID, options.OSS.AccessKeySecret)
		if err != nil {
			return nil, errors.Wrap(err, "oss.New failed")
		}
		if _, err := client.ListBuckets(); err != nil {
			return nil, errors.Wrap(err, "oss.Client.ListBuckets failed")
		}
		w.obj = client
	} else {
		res, err := alics.ECSMetaDataRamSecurityCredentials()
		if err != nil {
			return nil, errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentials failed")
		}
		client, err := oss.New(options.OSS.Endpoint, res.AccessKeyID, res.AccessKeySecret, oss.SecurityToken(res.SecurityToken))
		if err != nil {
			return nil, errors.Wrap(err, "oss.New failed")
		}
		if _, err := client.ListBuckets(); err != nil {
			return nil, errors.Wrap(err, "oss.Client.ListBuckets failed")
		}
		w.obj = client
		go w.UpdateCredentialByECSRole(res, &options.OSS)
	}

	return &w, nil
}

func NewOSSClientWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*OSSClientWrapper, error) {
	var options OSSClientWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewOSSClientWrapperWithOptions(&options, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "NewOSSClientWrapperWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Wrapper"), w.OnWrapperChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Retry"), w.OnRetryChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiterGroup"), w.OnRateLimiterGroupChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("OSS"), func(cfg *config.Config) error {
		var options OSSOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}

		if options.AccessKeyID != "" {
			client, err := oss.New(options.Endpoint, options.AccessKeyID, options.AccessKeySecret)
			if err != nil {
				return errors.Wrap(err, "oss.New failed")
			}
			if _, err := client.ListBuckets(); err != nil {
				return errors.Wrap(err, "oss.Client.ListBuckets failed")
			}
			w.obj = client
			return nil
		}

		res, err := alics.ECSMetaDataRamSecurityCredentials()
		if err != nil {
			return errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentials failed")
		}
		client, err := oss.New(options.Endpoint, res.AccessKeyID, res.AccessKeySecret, oss.SecurityToken(res.SecurityToken))
		if err != nil {
			return errors.Wrap(err, "oss.New failed")
		}
		if _, err := client.ListBuckets(); err != nil {
			return errors.Wrap(err, "oss.Client.ListBuckets failed")
		}
		w.obj = client
		go w.UpdateCredentialByECSRole(res, &options)

		return nil
	})

	return w, err
}

func (w *OSSClientWrapper) UpdateCredentialByECSRole(res *alics.ECSMetaDataRamSecurityCredentialsRes, options *OSSOptions) {
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
		client, err := oss.New(options.Endpoint, res.AccessKeyID, res.AccessKeySecret, oss.SecurityToken(res.SecurityToken))
		if err != nil {
			log.Errorf("oss.New failed. err: [%+v]", err)
			continue
		}
		if _, err := client.ListBuckets(); err != nil {
			log.Errorf("oss.Client.ListBuckets failed. err: [%v]", err)
			continue
		}

		// 检查 client 是否被其他协程更新
		if w.obj != oldClient {
			break
		}
		w.obj = client
	}
}
