package wrap

import (
	"time"

	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

type MNSOptions struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
}

type MNSClientWrapperOptions struct {
	Retry              micro.RetryOptions
	Wrapper            WrapperOptions
	MNS                MNSOptions
	RateLimiter        micro.RateLimiterOptions
	ParallelController micro.ParallelControllerOptions
}

func NewMNSClientWrapperWithOptions(options *MNSClientWrapperOptions, opts ...refx.Option) (*MNSClientWrapper, error) {
	var w MNSClientWrapper
	var err error

	w.options = &options.Wrapper
	w.retry, err = micro.NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "micro.NewRetryWithOptions failed")
	}
	w.rateLimiter, err = micro.NewRateLimiterWithOptions(&options.RateLimiter, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "micro.NewRateLimiterWithOptions failed")
	}
	w.parallelController, err = micro.NewParallelControllerWithOptions(&options.ParallelController, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "micro.NewParallelControllerWithOptions failed")
	}
	if w.options.EnableMetric {
		w.CreateMetric(w.options)
	}

	if options.MNS.AccessKeyID != "" {
		client := ali_mns.NewAliMNSClient(options.MNS.Endpoint, options.MNS.AccessKeyID, options.MNS.AccessKeySecret)
		w.obj = client
	} else {
		res, err := alics.ECSMetaDataRamSecurityCredentials()
		if err != nil {
			return nil, errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentials failed")
		}
		client := ali_mns.NewAliMNSClientWithToken(options.MNS.Endpoint, res.AccessKeyID, res.AccessKeySecret, res.SecurityToken)
		w.obj = client
		go w.UpdateCredentialByECSRole(res, &options.MNS)
	}

	return &w, nil
}

func NewMNSClientWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*MNSClientWrapper, error) {
	var options MNSClientWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewMNSClientWrapperWithOptions(&options, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "NewMNSClientWrapperWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Wrapper"), w.OnWrapperChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Retry"), w.OnRetryChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiter"), w.OnRateLimiterChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("ParallelController"), w.OnParallelControllerChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("MNS"), func(cfg *config.Config) error {
		var options MNSOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}

		if options.AccessKeyID != "" {
			client := ali_mns.NewAliMNSClient(options.Endpoint, options.AccessKeyID, options.AccessKeySecret)
			w.obj = client
			return nil
		}

		res, err := alics.ECSMetaDataRamSecurityCredentials()
		if err != nil {
			return errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentials failed")
		}
		client := ali_mns.NewAliMNSClientWithToken(options.Endpoint, res.AccessKeyID, res.AccessKeySecret, res.SecurityToken)
		w.obj = client
		go w.UpdateCredentialByECSRole(res, &options)

		return nil
	})

	return w, err
}

func (w *MNSClientWrapper) UpdateCredentialByECSRole(res *alics.ECSMetaDataRamSecurityCredentialsRes, options *MNSOptions) {
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
		client := ali_mns.NewAliMNSClientWithToken(options.Endpoint, res.AccessKeyID, res.AccessKeySecret, res.SecurityToken)

		// 检查 client 是否被其他协程更新
		if w.obj != oldClient {
			break
		}
		w.obj = client
	}
}

func NewMNSTopic(name string, w *MNSClientWrapper, qps ...int32) *MNSTopicWrapper {
	topic := ali_mns.NewMNSTopic(name, w.obj, qps...)
	return &MNSTopicWrapper{
		obj:                topic.(*ali_mns.MNSTopic),
		retry:              w.retry,
		options:            w.options,
		durationMetric:     w.durationMetric,
		inflightMetric:     w.inflightMetric,
		rateLimiter:        w.rateLimiter,
		parallelController: w.parallelController,
	}
}

func NewMNSQueue(name string, w *MNSClientWrapper, qps ...int32) *MNSQueueWrapper {
	queue := ali_mns.NewMNSQueue(name, w.obj, qps...)
	return &MNSQueueWrapper{
		obj:                queue.(*ali_mns.MNSQueue),
		retry:              w.retry,
		options:            w.options,
		durationMetric:     w.durationMetric,
		inflightMetric:     w.inflightMetric,
		rateLimiter:        w.rateLimiter,
		parallelController: w.parallelController,
	}
}
