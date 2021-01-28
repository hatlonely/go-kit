// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type STSClientWrapper struct {
	obj              *sts.Client
	retry            *Retry
	options          *WrapperOptions
	durationMetric   *prometheus.HistogramVec
	totalMetric      *prometheus.CounterVec
	rateLimiterGroup RateLimiterGroup
}

func (w *STSClientWrapper) Unwrap() *sts.Client {
	return w.obj
}

func (w *STSClientWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *STSClientWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *STSClientWrapper) OnRateLimiterGroupChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options RateLimiterGroupOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		rateLimiterGroup, err := NewRateLimiterGroupWithOptions(&options, opts...)
		if err != nil {
			return errors.Wrap(err, "NewRateLimiterGroupWithOptions failed")
		}
		w.rateLimiterGroup = rateLimiterGroup
		return nil
	}
}

func (w *STSClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "sts_Client_durationMs",
		Help:        "sts Client response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "sts_Client_total",
		Help:        "sts Client request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
}

func (w *STSClientWrapper) AssumeRole(ctx context.Context, request *sts.AssumeRoleRequest) (*sts.AssumeRoleResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *sts.AssumeRoleResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.AssumeRole"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "sts.Client.AssumeRole")
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
			defer func() {
				w.totalMetric.WithLabelValues("sts.Client.AssumeRole", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("sts.Client.AssumeRole", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AssumeRole(request)
		return err
	})
	return response, err
}

func (w *STSClientWrapper) AssumeRoleWithCallback(request *sts.AssumeRoleRequest, callback func(response *sts.AssumeRoleResponse, err error)) <-chan int {
	res0 := w.obj.AssumeRoleWithCallback(request, callback)
	return res0
}

func (w *STSClientWrapper) AssumeRoleWithChan(request *sts.AssumeRoleRequest) (<-chan *sts.AssumeRoleResponse, <-chan error) {
	res0, res1 := w.obj.AssumeRoleWithChan(request)
	return res0, res1
}

func (w *STSClientWrapper) AssumeRoleWithSAML(ctx context.Context, request *sts.AssumeRoleWithSAMLRequest) (*sts.AssumeRoleWithSAMLResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *sts.AssumeRoleWithSAMLResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.AssumeRoleWithSAML"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "sts.Client.AssumeRoleWithSAML")
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
			defer func() {
				w.totalMetric.WithLabelValues("sts.Client.AssumeRoleWithSAML", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("sts.Client.AssumeRoleWithSAML", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AssumeRoleWithSAML(request)
		return err
	})
	return response, err
}

func (w *STSClientWrapper) AssumeRoleWithSAMLWithCallback(request *sts.AssumeRoleWithSAMLRequest, callback func(response *sts.AssumeRoleWithSAMLResponse, err error)) <-chan int {
	res0 := w.obj.AssumeRoleWithSAMLWithCallback(request, callback)
	return res0
}

func (w *STSClientWrapper) AssumeRoleWithSAMLWithChan(request *sts.AssumeRoleWithSAMLRequest) (<-chan *sts.AssumeRoleWithSAMLResponse, <-chan error) {
	res0, res1 := w.obj.AssumeRoleWithSAMLWithChan(request)
	return res0, res1
}

func (w *STSClientWrapper) GetCallerIdentity(ctx context.Context, request *sts.GetCallerIdentityRequest) (*sts.GetCallerIdentityResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *sts.GetCallerIdentityResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetCallerIdentity"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "sts.Client.GetCallerIdentity")
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
			defer func() {
				w.totalMetric.WithLabelValues("sts.Client.GetCallerIdentity", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("sts.Client.GetCallerIdentity", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetCallerIdentity(request)
		return err
	})
	return response, err
}

func (w *STSClientWrapper) GetCallerIdentityWithCallback(request *sts.GetCallerIdentityRequest, callback func(response *sts.GetCallerIdentityResponse, err error)) <-chan int {
	res0 := w.obj.GetCallerIdentityWithCallback(request, callback)
	return res0
}

func (w *STSClientWrapper) GetCallerIdentityWithChan(request *sts.GetCallerIdentityRequest) (<-chan *sts.GetCallerIdentityResponse, <-chan error) {
	res0, res1 := w.obj.GetCallerIdentityWithChan(request)
	return res0, res1
}
