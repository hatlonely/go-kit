package config

import (
	"context"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/pkg/errors"
)

type ACMProviderOptions struct {
	Endpoint        string
	Namespace       string
	AccessKeyID     string
	AccessKeySecret string
	Timeout         time.Duration
	DataID          string
	Group           string
}

type ACMProvider struct {
	client  config_client.IConfigClient
	options *ACMProviderOptions

	buf    []byte
	events chan struct{}
	errors chan error
}

func NewACMProviderWithOptions(options *ACMProviderOptions) (*ACMProvider, error) {
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

	p := &ACMProvider{
		client:  client,
		options: options,
	}

	buf, err := p.load()
	if err != nil {
		return nil, errors.Wrap(err, "provider.load failed")
	}
	p.buf = buf

	return p, nil
}

func (p *ACMProvider) Events() <-chan struct{} {
	return p.events
}

func (p *ACMProvider) Errors() <-chan error {
	return p.errors
}

func (p *ACMProvider) Load() ([]byte, error) {
	return p.buf, nil
}

func (p *ACMProvider) load() ([]byte, error) {
	content, err := p.client.GetConfig(vo.ConfigParam{
		DataId: p.options.DataID,
		Group:  p.options.Group,
	})
	if err != nil {
		return nil, errors.Wrap(err, "acm.GetConfig failed")
	}
	return []byte(content), nil
}

func (p *ACMProvider) Dump(buf []byte) error {
	ok, err := p.client.PublishConfig(vo.ConfigParam{
		DataId:  p.options.DataID,
		Group:   p.options.Group,
		Content: string(buf),
	})
	if err != nil {
		return errors.Wrap(err, "acm.PublishConfig failed")
	}
	if !ok {
		return errors.Errorf("acm.PublishConfig not success")
	}
	return nil
}

func (p *ACMProvider) EventLoop(ctx context.Context) error {
	return p.client.ListenConfig(vo.ConfigParam{
		DataId: p.options.DataID,
		Group:  p.options.Group,
		OnChange: func(namespace, group, dataId, data string) {
			p.buf = []byte(data)
		},
	})
}
