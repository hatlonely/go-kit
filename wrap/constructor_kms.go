package wrap

import (
	"fmt"
	"net/http"
	"strings"

	alierr "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

func init() {
	RegisterErrCode(&alierr.ServerError{}, func(err error) string {
		e := err.(*alierr.ServerError)
		return fmt.Sprintf("pop_%v_%v", e.HttpStatus(), e.ErrorCode())
	})
	RegisterRetryRetryIf("POP", func(err error) bool {
		switch e := err.(type) {
		case *alierr.ServerError:
			if e.HttpStatus() >= http.StatusInternalServerError {
				return true
			}
			if strings.Contains(e.Error(), "timeout") {
				return true
			}
		}
		return false
	})
}

type KMSOptions struct {
	RegionID        string
	AccessKeyID     string
	AccessKeySecret string
}

type KMSClientWrapperOptions struct {
	Retry            RetryOptions
	Wrapper          WrapperOptions
	KMS              KMSOptions
	RateLimiterGroup RateLimiterGroupOptions
}

func NewKMSClientWrapperWithOptions(options *KMSClientWrapperOptions) (*KMSClientWrapper, error) {
	retry, err := NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}
	rateLimiterGroup, err := NewRateLimiterGroup(&options.RateLimiterGroup)
	if err != nil {
		return nil, errors.Wrap(err, "NewRateLimiterGroup failed")
	}

	if options.KMS.AccessKeyID != "" {
		client, err := kms.NewClientWithAccessKey(options.KMS.RegionID, options.KMS.AccessKeyID, options.KMS.AccessKeySecret)
		if err != nil {
			return nil, errors.Wrap(err, "kms.NewClientWithAccessKey failed")
		}
		if _, err := client.ListKeys(kms.CreateListKeysRequest()); err != nil {
			return nil, errors.Wrap(err, "kms.Client.ListKeys failed")
		}
		return &KMSClientWrapper{
			obj:              client,
			retry:            retry,
			options:          &options.Wrapper,
			rateLimiterGroup: rateLimiterGroup,
		}, nil
	}

	role, err := alics.ECSMetaDataRamSecurityCredentialsRole()
	if err != nil {
		return nil, errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentialsRole failed")
	}

	client, err := kms.NewClientWithEcsRamRole(options.KMS.RegionID, role)
	if err != nil {
		return nil, errors.Wrap(err, "kms.NewClientWithEcsRamRole failed")
	}
	if _, err := client.ListKeys(kms.CreateListKeysRequest()); err != nil {
		return nil, errors.Wrap(err, "kms.Client.ListKeys failed")
	}

	w := &KMSClientWrapper{
		obj:              client,
		retry:            retry,
		options:          &options.Wrapper,
		rateLimiterGroup: rateLimiterGroup,
	}

	if w.options.EnableMetric {
		w.CreateMetric(w.options)
	}

	return w, nil
}

func NewKMSClientWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*KMSClientWrapper, error) {
	var options KMSClientWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewKMSClientWrapperWithOptions(&options)
	if err != nil {
		return nil, errors.Wrap(err, "NewKMSClientWrapperWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Wrapper"), w.OnWrapperChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Retry"), w.OnRetryChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiterGroup"), w.OnRateLimiterGroupChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("KMS"), func(cfg *config.Config) error {
		var options KMSOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}

		if options.AccessKeyID != "" {
			client, err := kms.NewClientWithAccessKey(options.RegionID, options.AccessKeyID, options.AccessKeySecret)
			if err != nil {
				return errors.Wrap(err, "kms.NewClientWithAccessKey failed")
			}
			if _, err := client.ListKeys(kms.CreateListKeysRequest()); err != nil {
				return errors.Wrap(err, "kms.Client.ListKeys failed")
			}
			w.obj = client
			return nil
		}

		role, err := alics.ECSMetaDataRamSecurityCredentialsRole()
		if err != nil {
			return errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentialsRole failed")
		}
		client, err := kms.NewClientWithEcsRamRole(options.RegionID, role)
		if err != nil {
			return errors.Wrap(err, "kms.NewClientWithEcsRamRole failed")
		}
		if _, err := client.ListKeys(kms.CreateListKeysRequest()); err != nil {
			return errors.Wrap(err, "kms.Client.ListKeys failed")
		}
		w.obj = client

		return nil
	})

	return w, err
}
