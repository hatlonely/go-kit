package wrap

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/pkg/errors"
)

type ACMConfigClientWrapperOptions struct {
	Retry   RetryOptions
	Wrapper WrapperOptions
	ACM     constant.ClientConfig
}

func NewACMConfigClientWrapperWithOptions(options *ACMConfigClientWrapperOptions) (*ACMConfigClientWrapper, error) {
	retry, err := NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": options.ACM,
	})
	if err != nil {
		return nil, errors.Wrap(err, "clients.CreateConfigClient failed")
	}

	return &ACMConfigClientWrapper{
		obj:     client.(*config_client.ConfigClient),
		retry:   retry,
		options: &options.Wrapper,
	}, nil
}
