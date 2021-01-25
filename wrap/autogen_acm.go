// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type ACMConfigClientWrapper struct {
	obj            *config_client.ConfigClient
	retry          *Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	totalMetric    *prometheus.CounterVec
}

func (w *ACMConfigClientWrapper) Unwrap() *config_client.ConfigClient {
	return w.obj
}

func (w *ACMConfigClientWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *ACMConfigClientWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options RetryOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		retry, err := NewRetryWithOptions(&options)
		if err != nil {
			return errors.Wrap(err, "NewRetryWithOptions failed")
		}
		w.retry = retry
		return nil
	}
}

func (w *ACMConfigClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "config_client_ConfigClient_durationMs",
		Help:        "config_client ConfigClient response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "config_client_ConfigClient_total",
		Help:        "config_client ConfigClient request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode"})
}

func (w *ACMConfigClientWrapper) CancelListenConfig(ctx context.Context, param vo.ConfigParam) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.CancelListenConfig")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("config_client.ConfigClient.CancelListenConfig", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("config_client.ConfigClient.CancelListenConfig", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		err = w.obj.CancelListenConfig(param)
		return err
	})
	return err
}

func (w *ACMConfigClientWrapper) DeleteConfig(ctx context.Context, param vo.ConfigParam) (bool, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.DeleteConfig")
		defer span.Finish()
	}

	var deleted bool
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("config_client.ConfigClient.DeleteConfig", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("config_client.ConfigClient.DeleteConfig", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		deleted, err = w.obj.DeleteConfig(param)
		return err
	})
	return deleted, err
}

func (w *ACMConfigClientWrapper) GetConfig(ctx context.Context, param vo.ConfigParam) (string, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.GetConfig")
		defer span.Finish()
	}

	var content string
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("config_client.ConfigClient.GetConfig", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("config_client.ConfigClient.GetConfig", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		content, err = w.obj.GetConfig(param)
		return err
	})
	return content, err
}

func (w *ACMConfigClientWrapper) ListenConfig(ctx context.Context, param vo.ConfigParam) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.ListenConfig")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("config_client.ConfigClient.ListenConfig", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("config_client.ConfigClient.ListenConfig", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		err = w.obj.ListenConfig(param)
		return err
	})
	return err
}

func (w *ACMConfigClientWrapper) PublishConfig(ctx context.Context, param vo.ConfigParam) (bool, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.PublishConfig")
		defer span.Finish()
	}

	var published bool
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("config_client.ConfigClient.PublishConfig", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("config_client.ConfigClient.PublishConfig", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		published, err = w.obj.PublishConfig(param)
		return err
	})
	return published, err
}

func (w *ACMConfigClientWrapper) SearchConfig(ctx context.Context, param vo.SearchConfigParm) (*model.ConfigPage, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.SearchConfig")
		defer span.Finish()
	}

	var res0 *model.ConfigPage
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("config_client.ConfigClient.SearchConfig", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("config_client.ConfigClient.SearchConfig", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.SearchConfig(param)
		return err
	})
	return res0, err
}
