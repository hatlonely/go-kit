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
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

type ACMConfigClientWrapper struct {
	obj            *config_client.ConfigClient
	retry          *micro.Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	inflightMetric *prometheus.GaugeVec
	rateLimiter    micro.RateLimiter
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
		var options micro.RetryOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		retry, err := micro.NewRetryWithOptions(&options)
		if err != nil {
			return errors.Wrap(err, "NewRetryWithOptions failed")
		}
		w.retry = retry
		return nil
	}
}

func (w *ACMConfigClientWrapper) OnRateLimiterChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options micro.RateLimiterOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		rateLimiter, err := micro.NewRateLimiterWithOptions(&options, opts...)
		if err != nil {
			return errors.Wrap(err, "NewRateLimiterWithOptions failed")
		}
		w.rateLimiter = rateLimiter
		return nil
	}
}

func (w *ACMConfigClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "config_client_ConfigClient_durationMs",
		Help:        "config_client ConfigClient response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "config_client_ConfigClient_inflight",
		Help:        "config_client ConfigClient inflight",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "custom"})
}

func (w *ACMConfigClientWrapper) CancelListenConfig(ctx context.Context, param vo.ConfigParam) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ConfigClient.CancelListenConfig"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.CancelListenConfig")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("config_client.ConfigClient.CancelListenConfig", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("config_client.ConfigClient.CancelListenConfig", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("config_client.ConfigClient.CancelListenConfig", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.CancelListenConfig(param)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *ACMConfigClientWrapper) DeleteConfig(ctx context.Context, param vo.ConfigParam) (bool, error) {
	ctxOptions := FromContext(ctx)
	var deleted bool
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ConfigClient.DeleteConfig"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.DeleteConfig")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("config_client.ConfigClient.DeleteConfig", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("config_client.ConfigClient.DeleteConfig", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("config_client.ConfigClient.DeleteConfig", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		deleted, err = w.obj.DeleteConfig(param)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return deleted, err
}

func (w *ACMConfigClientWrapper) GetConfig(ctx context.Context, param vo.ConfigParam) (string, error) {
	ctxOptions := FromContext(ctx)
	var content string
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ConfigClient.GetConfig"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.GetConfig")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("config_client.ConfigClient.GetConfig", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("config_client.ConfigClient.GetConfig", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("config_client.ConfigClient.GetConfig", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		content, err = w.obj.GetConfig(param)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return content, err
}

func (w *ACMConfigClientWrapper) ListenConfig(ctx context.Context, param vo.ConfigParam) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ConfigClient.ListenConfig"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.ListenConfig")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("config_client.ConfigClient.ListenConfig", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("config_client.ConfigClient.ListenConfig", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("config_client.ConfigClient.ListenConfig", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.ListenConfig(param)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *ACMConfigClientWrapper) PublishConfig(ctx context.Context, param vo.ConfigParam) (bool, error) {
	ctxOptions := FromContext(ctx)
	var published bool
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ConfigClient.PublishConfig"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.PublishConfig")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("config_client.ConfigClient.PublishConfig", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("config_client.ConfigClient.PublishConfig", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("config_client.ConfigClient.PublishConfig", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		published, err = w.obj.PublishConfig(param)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return published, err
}

func (w *ACMConfigClientWrapper) SearchConfig(ctx context.Context, param vo.SearchConfigParm) (*model.ConfigPage, error) {
	ctxOptions := FromContext(ctx)
	var res0 *model.ConfigPage
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ConfigClient.SearchConfig"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.SearchConfig")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("config_client.ConfigClient.SearchConfig", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("config_client.ConfigClient.SearchConfig", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("config_client.ConfigClient.SearchConfig", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.SearchConfig(param)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}
