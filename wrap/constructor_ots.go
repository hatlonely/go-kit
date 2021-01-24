package wrap

import (
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
)

type OTSTableStoreClientWrapperOptions struct {
	Retry   RetryOptions
	Wrapper WrapperOptions
	OTS     struct {
		Endpoint        string
		AccessKeyID     string
		AccessKeySecret string
		InstanceName    string
	}
}

func NewOTSTableStoreClientWrapperWithOptions(options *OTSTableStoreClientWrapperOptions) (*OTSTableStoreClientWrapper, error) {
	retry, err := NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}

	if options.OTS.AccessKeyID != "" {
		client := tablestore.NewClient(options.OTS.Endpoint, options.OTS.InstanceName, options.OTS.AccessKeyID, options.OTS.AccessKeySecret)
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
	client := tablestore.NewClientWithConfig(options.OTS.Endpoint, options.OTS.InstanceName, res.AccessKeyID, res.AccessKeySecret, res.SecurityToken, nil)

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
			wrapper.obj = tablestore.NewClientWithConfig(options.OTS.Endpoint, options.OTS.InstanceName, res.AccessKeyID, res.AccessKeySecret, res.SecurityToken, nil)
		}
	}()

	return wrapper, nil
}
