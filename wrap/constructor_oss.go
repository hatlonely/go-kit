package wrap

import (
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
)

type OSSClientWrapperOptions struct {
	Retry           RetryOptions
	Wrapper         WrapperOptions
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
}

func NewOSSClientWrapperWithOptions(options *OSSClientWrapperOptions) (*OSSClientWrapper, error) {
	retry, err := NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}

	if options.AccessKeyID != "" {
		client, err := oss.New(options.Endpoint, options.AccessKeyID, options.AccessKeySecret)
		if err != nil {
			return nil, errors.Wrap(err, "oss.New failed")
		}
		return &OSSClientWrapper{
			obj:     client,
			retry:   retry,
			options: &options.Wrapper,
		}, nil
	}

	res, err := alics.ECSMetaDataRamSecurityCredentials()
	if err != nil {
		return nil, errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentials failed")
	}
	client, err := oss.New(options.Endpoint, res.AccessKeyID, res.AccessKeySecret, oss.SecurityToken(res.SecurityToken))
	if err != nil {
		return nil, errors.Wrap(err, "oss.New failed")
	}
	wrapper := &OSSClientWrapper{
		obj:     client,
		retry:   retry,
		options: &options.Wrapper,
	}

	go func() {
		for {
			d := res.ExpirationTime.Sub(time.Now()) - 5*time.Minute
			if d < 0 {
				d = 5 * time.Second
			}
			<-time.After(d)
			res, err = alics.ECSMetaDataRamSecurityCredentials()
			if err != nil {
				log.Errorf("alics.ECSMetaDataRamSecurityCredentials failed. err: [%+v]", err)
				continue
			}
			client, err = oss.New(options.Endpoint, res.AccessKeyID, res.AccessKeySecret, oss.SecurityToken(res.SecurityToken))
			if err != nil {
				log.Errorf("clients.CreateConfigClient failed. err: [%+v]", err)
				continue
			}
			wrapper.obj = client
		}
	}()

	return wrapper, nil
}
