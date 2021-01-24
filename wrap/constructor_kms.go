package wrap

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
)

type KMSClientWrapperOptions struct {
	Retry   RetryOptions
	Wrapper WrapperOptions
	KMS     struct {
		RegionID        string
		AccessKeyID     string
		AccessKeySecret string
	}
}

func NewKMSClientWrapperWithOptions(options *KMSClientWrapperOptions) (*KMSClientWrapper, error) {
	retry, err := NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}

	if options.KMS.AccessKeyID != "" {
		client, err := kms.NewClientWithAccessKey(options.KMS.RegionID, options.KMS.AccessKeyID, options.KMS.AccessKeySecret)
		if err != nil {
			return nil, errors.Wrap(err, "kms.NewClientWithAccessKey failed")
		}
		return &KMSClientWrapper{
			obj:     client,
			retry:   retry,
			options: &options.Wrapper,
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

	return &KMSClientWrapper{
		obj:     client,
		retry:   retry,
		options: &options.Wrapper,
	}, nil
}
