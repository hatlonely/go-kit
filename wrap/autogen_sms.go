// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

func NewSMSClientWrapper(
	obj *dysmsapi.Client,
	retry *micro.Retry,
	options *WrapperOptions,
	durationMetric *prometheus.HistogramVec,
	inflightMetric *prometheus.GaugeVec,
	rateLimiter micro.RateLimiter,
	parallelController micro.ParallelController) *SMSClientWrapper {
	return &SMSClientWrapper{
		obj:                obj,
		retry:              retry,
		options:            options,
		durationMetric:     durationMetric,
		inflightMetric:     inflightMetric,
		rateLimiter:        rateLimiter,
		parallelController: parallelController,
	}
}

type SMSClientWrapper struct {
	obj                *dysmsapi.Client
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *SMSClientWrapper) Unwrap() *dysmsapi.Client {
	return w.obj
}

func (w *SMSClientWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *SMSClientWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *SMSClientWrapper) OnRateLimiterChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *SMSClientWrapper) OnParallelControllerChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options micro.ParallelControllerOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		parallelController, err := micro.NewParallelControllerWithOptions(&options, opts...)
		if err != nil {
			return errors.Wrap(err, "NewParallelControllerWithOptions failed")
		}
		w.parallelController = parallelController
		return nil
	}
}

func (w *SMSClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "dysmsapi_Client_durationMs",
		Help:        "dysmsapi Client response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "dysmsapi_Client_inflight",
		Help:        "dysmsapi Client inflight",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "custom"})
}

func (w *SMSClientWrapper) AddShortUrl(ctx context.Context, request *dysmsapi.AddShortUrlRequest) (*dysmsapi.AddShortUrlResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.AddShortUrlResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AddShortUrl", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AddShortUrl", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AddShortUrl", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.AddShortUrl")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.AddShortUrl", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.AddShortUrl", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.AddShortUrl", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AddShortUrl(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) AddShortUrlWithCallback(request *dysmsapi.AddShortUrlRequest, callback func(response *dysmsapi.AddShortUrlResponse, err error)) <-chan int {
	res0 := w.obj.AddShortUrlWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) AddShortUrlWithChan(request *dysmsapi.AddShortUrlRequest) (<-chan *dysmsapi.AddShortUrlResponse, <-chan error) {
	res0, res1 := w.obj.AddShortUrlWithChan(request)
	return res0, res1
}

func (w *SMSClientWrapper) AddSmsSign(ctx context.Context, request *dysmsapi.AddSmsSignRequest) (*dysmsapi.AddSmsSignResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.AddSmsSignResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AddSmsSign", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AddSmsSign", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AddSmsSign", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.AddSmsSign")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.AddSmsSign", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.AddSmsSign", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.AddSmsSign", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AddSmsSign(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) AddSmsSignWithCallback(request *dysmsapi.AddSmsSignRequest, callback func(response *dysmsapi.AddSmsSignResponse, err error)) <-chan int {
	res0 := w.obj.AddSmsSignWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) AddSmsSignWithChan(request *dysmsapi.AddSmsSignRequest) (<-chan *dysmsapi.AddSmsSignResponse, <-chan error) {
	res0, res1 := w.obj.AddSmsSignWithChan(request)
	return res0, res1
}

func (w *SMSClientWrapper) AddSmsTemplate(ctx context.Context, request *dysmsapi.AddSmsTemplateRequest) (*dysmsapi.AddSmsTemplateResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.AddSmsTemplateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AddSmsTemplate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AddSmsTemplate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AddSmsTemplate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.AddSmsTemplate")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.AddSmsTemplate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.AddSmsTemplate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.AddSmsTemplate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AddSmsTemplate(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) AddSmsTemplateWithCallback(request *dysmsapi.AddSmsTemplateRequest, callback func(response *dysmsapi.AddSmsTemplateResponse, err error)) <-chan int {
	res0 := w.obj.AddSmsTemplateWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) AddSmsTemplateWithChan(request *dysmsapi.AddSmsTemplateRequest) (<-chan *dysmsapi.AddSmsTemplateResponse, <-chan error) {
	res0, res1 := w.obj.AddSmsTemplateWithChan(request)
	return res0, res1
}

func (w *SMSClientWrapper) CreateShortParam(ctx context.Context, request *dysmsapi.CreateShortParamRequest) (*dysmsapi.CreateShortParamResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.CreateShortParamResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateShortParam", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateShortParam", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateShortParam", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.CreateShortParam")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.CreateShortParam", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.CreateShortParam", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.CreateShortParam", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateShortParam(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) CreateShortParamWithCallback(request *dysmsapi.CreateShortParamRequest, callback func(response *dysmsapi.CreateShortParamResponse, err error)) <-chan int {
	res0 := w.obj.CreateShortParamWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) CreateShortParamWithChan(request *dysmsapi.CreateShortParamRequest) (<-chan *dysmsapi.CreateShortParamResponse, <-chan error) {
	res0, res1 := w.obj.CreateShortParamWithChan(request)
	return res0, res1
}

func (w *SMSClientWrapper) DeleteShortUrl(ctx context.Context, request *dysmsapi.DeleteShortUrlRequest) (*dysmsapi.DeleteShortUrlResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.DeleteShortUrlResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteShortUrl", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteShortUrl", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteShortUrl", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.DeleteShortUrl")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.DeleteShortUrl", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.DeleteShortUrl", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.DeleteShortUrl", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteShortUrl(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) DeleteShortUrlWithCallback(request *dysmsapi.DeleteShortUrlRequest, callback func(response *dysmsapi.DeleteShortUrlResponse, err error)) <-chan int {
	res0 := w.obj.DeleteShortUrlWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) DeleteShortUrlWithChan(request *dysmsapi.DeleteShortUrlRequest) (<-chan *dysmsapi.DeleteShortUrlResponse, <-chan error) {
	res0, res1 := w.obj.DeleteShortUrlWithChan(request)
	return res0, res1
}

func (w *SMSClientWrapper) DeleteSmsSign(ctx context.Context, request *dysmsapi.DeleteSmsSignRequest) (*dysmsapi.DeleteSmsSignResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.DeleteSmsSignResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteSmsSign", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteSmsSign", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteSmsSign", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.DeleteSmsSign")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.DeleteSmsSign", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.DeleteSmsSign", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.DeleteSmsSign", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteSmsSign(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) DeleteSmsSignWithCallback(request *dysmsapi.DeleteSmsSignRequest, callback func(response *dysmsapi.DeleteSmsSignResponse, err error)) <-chan int {
	res0 := w.obj.DeleteSmsSignWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) DeleteSmsSignWithChan(request *dysmsapi.DeleteSmsSignRequest) (<-chan *dysmsapi.DeleteSmsSignResponse, <-chan error) {
	res0, res1 := w.obj.DeleteSmsSignWithChan(request)
	return res0, res1
}

func (w *SMSClientWrapper) DeleteSmsTemplate(ctx context.Context, request *dysmsapi.DeleteSmsTemplateRequest) (*dysmsapi.DeleteSmsTemplateResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.DeleteSmsTemplateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteSmsTemplate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteSmsTemplate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteSmsTemplate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.DeleteSmsTemplate")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.DeleteSmsTemplate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.DeleteSmsTemplate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.DeleteSmsTemplate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteSmsTemplate(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) DeleteSmsTemplateWithCallback(request *dysmsapi.DeleteSmsTemplateRequest, callback func(response *dysmsapi.DeleteSmsTemplateResponse, err error)) <-chan int {
	res0 := w.obj.DeleteSmsTemplateWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) DeleteSmsTemplateWithChan(request *dysmsapi.DeleteSmsTemplateRequest) (<-chan *dysmsapi.DeleteSmsTemplateResponse, <-chan error) {
	res0, res1 := w.obj.DeleteSmsTemplateWithChan(request)
	return res0, res1
}

func (w *SMSClientWrapper) ModifySmsSign(ctx context.Context, request *dysmsapi.ModifySmsSignRequest) (*dysmsapi.ModifySmsSignResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.ModifySmsSignResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ModifySmsSign", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ModifySmsSign", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ModifySmsSign", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.ModifySmsSign")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.ModifySmsSign", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.ModifySmsSign", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.ModifySmsSign", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ModifySmsSign(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) ModifySmsSignWithCallback(request *dysmsapi.ModifySmsSignRequest, callback func(response *dysmsapi.ModifySmsSignResponse, err error)) <-chan int {
	res0 := w.obj.ModifySmsSignWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) ModifySmsSignWithChan(request *dysmsapi.ModifySmsSignRequest) (<-chan *dysmsapi.ModifySmsSignResponse, <-chan error) {
	res0, res1 := w.obj.ModifySmsSignWithChan(request)
	return res0, res1
}

func (w *SMSClientWrapper) ModifySmsTemplate(ctx context.Context, request *dysmsapi.ModifySmsTemplateRequest) (*dysmsapi.ModifySmsTemplateResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.ModifySmsTemplateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ModifySmsTemplate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ModifySmsTemplate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ModifySmsTemplate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.ModifySmsTemplate")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.ModifySmsTemplate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.ModifySmsTemplate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.ModifySmsTemplate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ModifySmsTemplate(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) ModifySmsTemplateWithCallback(request *dysmsapi.ModifySmsTemplateRequest, callback func(response *dysmsapi.ModifySmsTemplateResponse, err error)) <-chan int {
	res0 := w.obj.ModifySmsTemplateWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) ModifySmsTemplateWithChan(request *dysmsapi.ModifySmsTemplateRequest) (<-chan *dysmsapi.ModifySmsTemplateResponse, <-chan error) {
	res0, res1 := w.obj.ModifySmsTemplateWithChan(request)
	return res0, res1
}

func (w *SMSClientWrapper) QuerySendDetails(ctx context.Context, request *dysmsapi.QuerySendDetailsRequest) (*dysmsapi.QuerySendDetailsResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.QuerySendDetailsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.QuerySendDetails", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.QuerySendDetails", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.QuerySendDetails", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.QuerySendDetails")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.QuerySendDetails", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.QuerySendDetails", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.QuerySendDetails", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.QuerySendDetails(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) QuerySendDetailsWithCallback(request *dysmsapi.QuerySendDetailsRequest, callback func(response *dysmsapi.QuerySendDetailsResponse, err error)) <-chan int {
	res0 := w.obj.QuerySendDetailsWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) QuerySendDetailsWithChan(request *dysmsapi.QuerySendDetailsRequest) (<-chan *dysmsapi.QuerySendDetailsResponse, <-chan error) {
	res0, res1 := w.obj.QuerySendDetailsWithChan(request)
	return res0, res1
}

func (w *SMSClientWrapper) QueryShortUrl(ctx context.Context, request *dysmsapi.QueryShortUrlRequest) (*dysmsapi.QueryShortUrlResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.QueryShortUrlResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.QueryShortUrl", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.QueryShortUrl", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.QueryShortUrl", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.QueryShortUrl")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.QueryShortUrl", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.QueryShortUrl", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.QueryShortUrl", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.QueryShortUrl(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) QueryShortUrlWithCallback(request *dysmsapi.QueryShortUrlRequest, callback func(response *dysmsapi.QueryShortUrlResponse, err error)) <-chan int {
	res0 := w.obj.QueryShortUrlWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) QueryShortUrlWithChan(request *dysmsapi.QueryShortUrlRequest) (<-chan *dysmsapi.QueryShortUrlResponse, <-chan error) {
	res0, res1 := w.obj.QueryShortUrlWithChan(request)
	return res0, res1
}

func (w *SMSClientWrapper) QuerySmsSign(ctx context.Context, request *dysmsapi.QuerySmsSignRequest) (*dysmsapi.QuerySmsSignResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.QuerySmsSignResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.QuerySmsSign", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.QuerySmsSign", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.QuerySmsSign", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.QuerySmsSign")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.QuerySmsSign", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.QuerySmsSign", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.QuerySmsSign", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.QuerySmsSign(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) QuerySmsSignWithCallback(request *dysmsapi.QuerySmsSignRequest, callback func(response *dysmsapi.QuerySmsSignResponse, err error)) <-chan int {
	res0 := w.obj.QuerySmsSignWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) QuerySmsSignWithChan(request *dysmsapi.QuerySmsSignRequest) (<-chan *dysmsapi.QuerySmsSignResponse, <-chan error) {
	res0, res1 := w.obj.QuerySmsSignWithChan(request)
	return res0, res1
}

func (w *SMSClientWrapper) QuerySmsTemplate(ctx context.Context, request *dysmsapi.QuerySmsTemplateRequest) (*dysmsapi.QuerySmsTemplateResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.QuerySmsTemplateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.QuerySmsTemplate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.QuerySmsTemplate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.QuerySmsTemplate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.QuerySmsTemplate")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.QuerySmsTemplate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.QuerySmsTemplate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.QuerySmsTemplate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.QuerySmsTemplate(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) QuerySmsTemplateWithCallback(request *dysmsapi.QuerySmsTemplateRequest, callback func(response *dysmsapi.QuerySmsTemplateResponse, err error)) <-chan int {
	res0 := w.obj.QuerySmsTemplateWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) QuerySmsTemplateWithChan(request *dysmsapi.QuerySmsTemplateRequest) (<-chan *dysmsapi.QuerySmsTemplateResponse, <-chan error) {
	res0, res1 := w.obj.QuerySmsTemplateWithChan(request)
	return res0, res1
}

func (w *SMSClientWrapper) SendBatchSms(ctx context.Context, request *dysmsapi.SendBatchSmsRequest) (*dysmsapi.SendBatchSmsResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.SendBatchSmsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SendBatchSms", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SendBatchSms", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SendBatchSms", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.SendBatchSms")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.SendBatchSms", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.SendBatchSms", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.SendBatchSms", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.SendBatchSms(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) SendBatchSmsWithCallback(request *dysmsapi.SendBatchSmsRequest, callback func(response *dysmsapi.SendBatchSmsResponse, err error)) <-chan int {
	res0 := w.obj.SendBatchSmsWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) SendBatchSmsWithChan(request *dysmsapi.SendBatchSmsRequest) (<-chan *dysmsapi.SendBatchSmsResponse, <-chan error) {
	res0, res1 := w.obj.SendBatchSmsWithChan(request)
	return res0, res1
}

func (w *SMSClientWrapper) SendSms(ctx context.Context, request *dysmsapi.SendSmsRequest) (*dysmsapi.SendSmsResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *dysmsapi.SendSmsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SendSms", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SendSms", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SendSms", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "dysmsapi.Client.SendSms")
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
			w.inflightMetric.WithLabelValues("dysmsapi.Client.SendSms", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("dysmsapi.Client.SendSms", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("dysmsapi.Client.SendSms", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.SendSms(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *SMSClientWrapper) SendSmsWithCallback(request *dysmsapi.SendSmsRequest, callback func(response *dysmsapi.SendSmsResponse, err error)) <-chan int {
	res0 := w.obj.SendSmsWithCallback(request, callback)
	return res0
}

func (w *SMSClientWrapper) SendSmsWithChan(request *dysmsapi.SendSmsRequest) (<-chan *dysmsapi.SendSmsResponse, <-chan error) {
	res0, res1 := w.obj.SendSmsWithChan(request)
	return res0, res1
}
