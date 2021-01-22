package wrap

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
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
	client := tablestore.NewClient(options.Endpoint, options.InstanceName, res.AccessKeyID, res.AccessKeySecret)

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
			wrapper.obj = tablestore.NewClient(options.Endpoint, options.InstanceName, res.AccessKeyID, res.AccessKeySecret)
		}
	}()

	return wrapper, nil
}

type KMSClientWrapperOptions struct {
	RegionID        string
	AccessKeyID     string
	AccessKeySecret string
	Retry           RetryOptions
	Wrapper         WrapperOptions
}

func NewKMSClientWrapperWithOptions(options *KMSClientWrapperOptions) (*KMSClientWrapper, error) {
	retry, err := NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}

	if options.AccessKeyID != "" {
		client, err := kms.NewClientWithAccessKey(options.RegionID, options.AccessKeyID, options.AccessKeySecret)
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

	client, err := kms.NewClientWithEcsRamRole(options.RegionID, role)
	if err != nil {
		return nil, errors.Wrap(err, "kms.NewClientWithEcsRamRole failed")
	}

	return &KMSClientWrapper{
		obj:     client,
		retry:   retry,
		options: &options.Wrapper,
	}, nil
}

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
	client, err := oss.New(options.Endpoint, res.AccessKeyID, res.AccessKeySecret)
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
			client, err = oss.New(options.Endpoint, res.AccessKeyID, res.AccessKeySecret)
			if err != nil {
				log.Errorf("clients.CreateConfigClient failed. err: [%+v]", err)
				continue
			}
			wrapper.obj = client
		}
	}()

	return wrapper, nil
}
