package wrap

import (
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
)

type ACMConfigClientWrapperOptions struct {
	Retry           RetryOptions
	Wrapper         WrapperOptions
	Endpoint        string
	Namespace       string
	AccessKeyID     string
	AccessKeySecret string
	Timeout         time.Duration
}

func NewACMConfigClientWrapperWithOptions(options *ACMConfigClientWrapperOptions) (*ACMConfigClientWrapper, error) {
	retry, err := NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}

	if options.AccessKeyID != "" {
		client, err := clients.CreateConfigClient(map[string]interface{}{
			"clientConfig": constant.ClientConfig{
				Endpoint:    options.Endpoint,
				NamespaceId: options.Namespace,
				AccessKey:   options.AccessKeyID,
				SecretKey:   options.AccessKeySecret,
				TimeoutMs:   uint64(options.Timeout.Milliseconds()),
			},
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

	res, err := alics.ECSMetaDataRamSecurityCredentials()
	if err != nil {
		return nil, errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentials failed")
	}
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": constant.ClientConfig{
			Endpoint:    options.Endpoint,
			NamespaceId: options.Namespace,
			AccessKey:   res.AccessKeyID,
			SecretKey:   res.AccessKeySecret,
			TimeoutMs:   uint64(options.Timeout.Milliseconds()),
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "clients.CreateConfigClient failed")
	}

	wrapper := &ACMConfigClientWrapper{
		obj:     client.(*config_client.ConfigClient),
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
			client, err = clients.CreateConfigClient(map[string]interface{}{
				"clientConfig": constant.ClientConfig{
					Endpoint:    options.Endpoint,
					NamespaceId: options.Namespace,
					AccessKey:   res.AccessKeyID,
					SecretKey:   res.AccessKeySecret,
					TimeoutMs:   uint64(options.Timeout.Milliseconds()),
				},
			})
			if err != nil {
				log.Errorf("clients.CreateConfigClient failed. err: [%+v]", err)
				continue
			}
			wrapper.obj = client.(*config_client.ConfigClient)
		}
	}()

	return wrapper, nil
}
