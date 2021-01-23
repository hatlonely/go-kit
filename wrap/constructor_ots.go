package wrap

import (
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
)

type OTSTableStoreClientWrapperOptions struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	InstanceName    string
	Retry           RetryOptions
	Wrapper         WrapperOptions
}

func NewOTSTableStoreClientWrapperWithOptions(options *OTSTableStoreClientWrapperOptions) (*OTSTableStoreClientWrapper, error) {
	retry, err := NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}

	if options.AccessKeyID != "" {
		client := tablestore.NewClient(options.Endpoint, options.InstanceName, options.AccessKeyID, options.AccessKeySecret)
		return &OTSTableStoreClientWrapper{
			obj:     client,
			retry:   retry,
			options: &options.Wrapper,
		}, nil
	}

	res, err := alics.ECSMetaDataRamSecurityCredentials()
	if err != nil {
		return nil, errors.Wrap(err, "ECSMetaDataRamSecurityCredentials failed")
	}
	client := tablestore.NewClientWithConfig(options.Endpoint, options.InstanceName, res.AccessKeyID, res.AccessKeySecret, res.SecurityToken, nil)

	wrapper := &OTSTableStoreClientWrapper{
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
				log.Errorf("ECSMetaDataRamSecurityCredentials failed. err: [%+v]", err)
				continue
			}
			wrapper.obj = tablestore.NewClientWithConfig(options.Endpoint, options.InstanceName, res.AccessKeyID, res.AccessKeySecret, res.SecurityToken, nil)
		}
	}()

	return wrapper, nil
}