// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

func NewKMSClientWrapper(
	obj *kms.Client,
	retry *micro.Retry,
	options *WrapperOptions,
	durationMetric *prometheus.HistogramVec,
	inflightMetric *prometheus.GaugeVec,
	rateLimiter micro.RateLimiter,
	parallelController micro.ParallelController) *KMSClientWrapper {
	return &KMSClientWrapper{
		obj:                obj,
		retry:              retry,
		options:            options,
		durationMetric:     durationMetric,
		inflightMetric:     inflightMetric,
		rateLimiter:        rateLimiter,
		parallelController: parallelController,
	}
}

type KMSClientWrapper struct {
	obj                *kms.Client
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *KMSClientWrapper) Unwrap() *kms.Client {
	return w.obj
}

func (w *KMSClientWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *KMSClientWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *KMSClientWrapper) OnRateLimiterChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *KMSClientWrapper) OnParallelControllerChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *KMSClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        fmt.Sprintf("%s_kms_Client_durationMs", options.Name),
		Help:        "kms Client response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        fmt.Sprintf("%s_kms_Client_inflight", options.Name),
		Help:        "kms Client inflight",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "custom"})
}

func (w *KMSClientWrapper) AsymmetricDecrypt(ctx context.Context, request *kms.AsymmetricDecryptRequest) (*kms.AsymmetricDecryptResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.AsymmetricDecryptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AsymmetricDecrypt", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AsymmetricDecrypt", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AsymmetricDecrypt", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricDecrypt", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.AsymmetricDecrypt", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.AsymmetricDecrypt", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.AsymmetricDecrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AsymmetricDecrypt(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) AsymmetricDecryptWithCallback(request *kms.AsymmetricDecryptRequest, callback func(response *kms.AsymmetricDecryptResponse, err error)) <-chan int {
	res0 := w.obj.AsymmetricDecryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) AsymmetricDecryptWithChan(request *kms.AsymmetricDecryptRequest) (<-chan *kms.AsymmetricDecryptResponse, <-chan error) {
	res0, res1 := w.obj.AsymmetricDecryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) AsymmetricEncrypt(ctx context.Context, request *kms.AsymmetricEncryptRequest) (*kms.AsymmetricEncryptResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.AsymmetricEncryptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AsymmetricEncrypt", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AsymmetricEncrypt", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AsymmetricEncrypt", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricEncrypt", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.AsymmetricEncrypt", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.AsymmetricEncrypt", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.AsymmetricEncrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AsymmetricEncrypt(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) AsymmetricEncryptWithCallback(request *kms.AsymmetricEncryptRequest, callback func(response *kms.AsymmetricEncryptResponse, err error)) <-chan int {
	res0 := w.obj.AsymmetricEncryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) AsymmetricEncryptWithChan(request *kms.AsymmetricEncryptRequest) (<-chan *kms.AsymmetricEncryptResponse, <-chan error) {
	res0, res1 := w.obj.AsymmetricEncryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) AsymmetricSign(ctx context.Context, request *kms.AsymmetricSignRequest) (*kms.AsymmetricSignResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.AsymmetricSignResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AsymmetricSign", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AsymmetricSign", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AsymmetricSign", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricSign", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.AsymmetricSign", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.AsymmetricSign", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.AsymmetricSign", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AsymmetricSign(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) AsymmetricSignWithCallback(request *kms.AsymmetricSignRequest, callback func(response *kms.AsymmetricSignResponse, err error)) <-chan int {
	res0 := w.obj.AsymmetricSignWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) AsymmetricSignWithChan(request *kms.AsymmetricSignRequest) (<-chan *kms.AsymmetricSignResponse, <-chan error) {
	res0, res1 := w.obj.AsymmetricSignWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) AsymmetricVerify(ctx context.Context, request *kms.AsymmetricVerifyRequest) (*kms.AsymmetricVerifyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.AsymmetricVerifyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AsymmetricVerify", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AsymmetricVerify", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AsymmetricVerify", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricVerify", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.AsymmetricVerify", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.AsymmetricVerify", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.AsymmetricVerify", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AsymmetricVerify(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) AsymmetricVerifyWithCallback(request *kms.AsymmetricVerifyRequest, callback func(response *kms.AsymmetricVerifyResponse, err error)) <-chan int {
	res0 := w.obj.AsymmetricVerifyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) AsymmetricVerifyWithChan(request *kms.AsymmetricVerifyRequest) (<-chan *kms.AsymmetricVerifyResponse, <-chan error) {
	res0, res1 := w.obj.AsymmetricVerifyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CancelKeyDeletion(ctx context.Context, request *kms.CancelKeyDeletionRequest) (*kms.CancelKeyDeletionResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.CancelKeyDeletionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CancelKeyDeletion", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CancelKeyDeletion", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CancelKeyDeletion", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.CancelKeyDeletion", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.CancelKeyDeletion", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.CancelKeyDeletion", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.CancelKeyDeletion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CancelKeyDeletion(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CancelKeyDeletionWithCallback(request *kms.CancelKeyDeletionRequest, callback func(response *kms.CancelKeyDeletionResponse, err error)) <-chan int {
	res0 := w.obj.CancelKeyDeletionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CancelKeyDeletionWithChan(request *kms.CancelKeyDeletionRequest) (<-chan *kms.CancelKeyDeletionResponse, <-chan error) {
	res0, res1 := w.obj.CancelKeyDeletionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CertificatePrivateKeyDecrypt(ctx context.Context, request *kms.CertificatePrivateKeyDecryptRequest) (*kms.CertificatePrivateKeyDecryptResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.CertificatePrivateKeyDecryptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CertificatePrivateKeyDecrypt", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CertificatePrivateKeyDecrypt", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CertificatePrivateKeyDecrypt", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePrivateKeyDecrypt", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.CertificatePrivateKeyDecrypt", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.CertificatePrivateKeyDecrypt", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.CertificatePrivateKeyDecrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CertificatePrivateKeyDecrypt(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CertificatePrivateKeyDecryptWithCallback(request *kms.CertificatePrivateKeyDecryptRequest, callback func(response *kms.CertificatePrivateKeyDecryptResponse, err error)) <-chan int {
	res0 := w.obj.CertificatePrivateKeyDecryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CertificatePrivateKeyDecryptWithChan(request *kms.CertificatePrivateKeyDecryptRequest) (<-chan *kms.CertificatePrivateKeyDecryptResponse, <-chan error) {
	res0, res1 := w.obj.CertificatePrivateKeyDecryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CertificatePrivateKeySign(ctx context.Context, request *kms.CertificatePrivateKeySignRequest) (*kms.CertificatePrivateKeySignResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.CertificatePrivateKeySignResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CertificatePrivateKeySign", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CertificatePrivateKeySign", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CertificatePrivateKeySign", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePrivateKeySign", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.CertificatePrivateKeySign", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.CertificatePrivateKeySign", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.CertificatePrivateKeySign", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CertificatePrivateKeySign(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CertificatePrivateKeySignWithCallback(request *kms.CertificatePrivateKeySignRequest, callback func(response *kms.CertificatePrivateKeySignResponse, err error)) <-chan int {
	res0 := w.obj.CertificatePrivateKeySignWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CertificatePrivateKeySignWithChan(request *kms.CertificatePrivateKeySignRequest) (<-chan *kms.CertificatePrivateKeySignResponse, <-chan error) {
	res0, res1 := w.obj.CertificatePrivateKeySignWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CertificatePublicKeyEncrypt(ctx context.Context, request *kms.CertificatePublicKeyEncryptRequest) (*kms.CertificatePublicKeyEncryptResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.CertificatePublicKeyEncryptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CertificatePublicKeyEncrypt", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CertificatePublicKeyEncrypt", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CertificatePublicKeyEncrypt", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePublicKeyEncrypt", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.CertificatePublicKeyEncrypt", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.CertificatePublicKeyEncrypt", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.CertificatePublicKeyEncrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CertificatePublicKeyEncrypt(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CertificatePublicKeyEncryptWithCallback(request *kms.CertificatePublicKeyEncryptRequest, callback func(response *kms.CertificatePublicKeyEncryptResponse, err error)) <-chan int {
	res0 := w.obj.CertificatePublicKeyEncryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CertificatePublicKeyEncryptWithChan(request *kms.CertificatePublicKeyEncryptRequest) (<-chan *kms.CertificatePublicKeyEncryptResponse, <-chan error) {
	res0, res1 := w.obj.CertificatePublicKeyEncryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CertificatePublicKeyVerify(ctx context.Context, request *kms.CertificatePublicKeyVerifyRequest) (*kms.CertificatePublicKeyVerifyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.CertificatePublicKeyVerifyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CertificatePublicKeyVerify", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CertificatePublicKeyVerify", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CertificatePublicKeyVerify", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePublicKeyVerify", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.CertificatePublicKeyVerify", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.CertificatePublicKeyVerify", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.CertificatePublicKeyVerify", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CertificatePublicKeyVerify(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CertificatePublicKeyVerifyWithCallback(request *kms.CertificatePublicKeyVerifyRequest, callback func(response *kms.CertificatePublicKeyVerifyResponse, err error)) <-chan int {
	res0 := w.obj.CertificatePublicKeyVerifyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CertificatePublicKeyVerifyWithChan(request *kms.CertificatePublicKeyVerifyRequest) (<-chan *kms.CertificatePublicKeyVerifyResponse, <-chan error) {
	res0, res1 := w.obj.CertificatePublicKeyVerifyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateAlias(ctx context.Context, request *kms.CreateAliasRequest) (*kms.CreateAliasResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.CreateAliasResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateAlias", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateAlias", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateAlias", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.CreateAlias", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.CreateAlias", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.CreateAlias", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.CreateAlias", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateAlias(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateAliasWithCallback(request *kms.CreateAliasRequest, callback func(response *kms.CreateAliasResponse, err error)) <-chan int {
	res0 := w.obj.CreateAliasWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateAliasWithChan(request *kms.CreateAliasRequest) (<-chan *kms.CreateAliasResponse, <-chan error) {
	res0, res1 := w.obj.CreateAliasWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateCertificate(ctx context.Context, request *kms.CreateCertificateRequest) (*kms.CreateCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.CreateCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateCertificate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateCertificate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateCertificate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.CreateCertificate", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.CreateCertificate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.CreateCertificate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.CreateCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateCertificate(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateCertificateWithCallback(request *kms.CreateCertificateRequest, callback func(response *kms.CreateCertificateResponse, err error)) <-chan int {
	res0 := w.obj.CreateCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateCertificateWithChan(request *kms.CreateCertificateRequest) (<-chan *kms.CreateCertificateResponse, <-chan error) {
	res0, res1 := w.obj.CreateCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateKey(ctx context.Context, request *kms.CreateKeyRequest) (*kms.CreateKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.CreateKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateKey", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateKey", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateKey", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.CreateKey", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.CreateKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.CreateKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.CreateKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateKey(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateKeyVersion(ctx context.Context, request *kms.CreateKeyVersionRequest) (*kms.CreateKeyVersionResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.CreateKeyVersionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateKeyVersion", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateKeyVersion", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateKeyVersion", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.CreateKeyVersion", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.CreateKeyVersion", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.CreateKeyVersion", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.CreateKeyVersion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateKeyVersion(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateKeyVersionWithCallback(request *kms.CreateKeyVersionRequest, callback func(response *kms.CreateKeyVersionResponse, err error)) <-chan int {
	res0 := w.obj.CreateKeyVersionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateKeyVersionWithChan(request *kms.CreateKeyVersionRequest) (<-chan *kms.CreateKeyVersionResponse, <-chan error) {
	res0, res1 := w.obj.CreateKeyVersionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateKeyWithCallback(request *kms.CreateKeyRequest, callback func(response *kms.CreateKeyResponse, err error)) <-chan int {
	res0 := w.obj.CreateKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateKeyWithChan(request *kms.CreateKeyRequest) (<-chan *kms.CreateKeyResponse, <-chan error) {
	res0, res1 := w.obj.CreateKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateSecret(ctx context.Context, request *kms.CreateSecretRequest) (*kms.CreateSecretResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.CreateSecretResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateSecret", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateSecret", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateSecret", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.CreateSecret", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.CreateSecret", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.CreateSecret", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.CreateSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateSecret(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateSecretWithCallback(request *kms.CreateSecretRequest, callback func(response *kms.CreateSecretResponse, err error)) <-chan int {
	res0 := w.obj.CreateSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateSecretWithChan(request *kms.CreateSecretRequest) (<-chan *kms.CreateSecretResponse, <-chan error) {
	res0, res1 := w.obj.CreateSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) Decrypt(ctx context.Context, request *kms.DecryptRequest) (*kms.DecryptResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.DecryptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.Decrypt", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.Decrypt", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.Decrypt", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.Decrypt", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.Decrypt", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.Decrypt", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.Decrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.Decrypt(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DecryptWithCallback(request *kms.DecryptRequest, callback func(response *kms.DecryptResponse, err error)) <-chan int {
	res0 := w.obj.DecryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DecryptWithChan(request *kms.DecryptRequest) (<-chan *kms.DecryptResponse, <-chan error) {
	res0, res1 := w.obj.DecryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DeleteAlias(ctx context.Context, request *kms.DeleteAliasRequest) (*kms.DeleteAliasResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.DeleteAliasResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteAlias", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteAlias", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteAlias", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteAlias", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.DeleteAlias", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.DeleteAlias", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.DeleteAlias", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteAlias(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DeleteAliasWithCallback(request *kms.DeleteAliasRequest, callback func(response *kms.DeleteAliasResponse, err error)) <-chan int {
	res0 := w.obj.DeleteAliasWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DeleteAliasWithChan(request *kms.DeleteAliasRequest) (<-chan *kms.DeleteAliasResponse, <-chan error) {
	res0, res1 := w.obj.DeleteAliasWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DeleteCertificate(ctx context.Context, request *kms.DeleteCertificateRequest) (*kms.DeleteCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.DeleteCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteCertificate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteCertificate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteCertificate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteCertificate", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.DeleteCertificate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.DeleteCertificate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.DeleteCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteCertificate(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DeleteCertificateWithCallback(request *kms.DeleteCertificateRequest, callback func(response *kms.DeleteCertificateResponse, err error)) <-chan int {
	res0 := w.obj.DeleteCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DeleteCertificateWithChan(request *kms.DeleteCertificateRequest) (<-chan *kms.DeleteCertificateResponse, <-chan error) {
	res0, res1 := w.obj.DeleteCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DeleteKeyMaterial(ctx context.Context, request *kms.DeleteKeyMaterialRequest) (*kms.DeleteKeyMaterialResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.DeleteKeyMaterialResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteKeyMaterial", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteKeyMaterial", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteKeyMaterial", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteKeyMaterial", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.DeleteKeyMaterial", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.DeleteKeyMaterial", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.DeleteKeyMaterial", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteKeyMaterial(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DeleteKeyMaterialWithCallback(request *kms.DeleteKeyMaterialRequest, callback func(response *kms.DeleteKeyMaterialResponse, err error)) <-chan int {
	res0 := w.obj.DeleteKeyMaterialWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DeleteKeyMaterialWithChan(request *kms.DeleteKeyMaterialRequest) (<-chan *kms.DeleteKeyMaterialResponse, <-chan error) {
	res0, res1 := w.obj.DeleteKeyMaterialWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DeleteSecret(ctx context.Context, request *kms.DeleteSecretRequest) (*kms.DeleteSecretResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.DeleteSecretResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteSecret", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteSecret", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteSecret", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteSecret", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.DeleteSecret", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.DeleteSecret", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.DeleteSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteSecret(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DeleteSecretWithCallback(request *kms.DeleteSecretRequest, callback func(response *kms.DeleteSecretResponse, err error)) <-chan int {
	res0 := w.obj.DeleteSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DeleteSecretWithChan(request *kms.DeleteSecretRequest) (<-chan *kms.DeleteSecretResponse, <-chan error) {
	res0, res1 := w.obj.DeleteSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeAccountKmsStatus(ctx context.Context, request *kms.DescribeAccountKmsStatusRequest) (*kms.DescribeAccountKmsStatusResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.DescribeAccountKmsStatusResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DescribeAccountKmsStatus", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DescribeAccountKmsStatus", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DescribeAccountKmsStatus", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeAccountKmsStatus", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.DescribeAccountKmsStatus", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.DescribeAccountKmsStatus", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.DescribeAccountKmsStatus", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DescribeAccountKmsStatus(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeAccountKmsStatusWithCallback(request *kms.DescribeAccountKmsStatusRequest, callback func(response *kms.DescribeAccountKmsStatusResponse, err error)) <-chan int {
	res0 := w.obj.DescribeAccountKmsStatusWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeAccountKmsStatusWithChan(request *kms.DescribeAccountKmsStatusRequest) (<-chan *kms.DescribeAccountKmsStatusResponse, <-chan error) {
	res0, res1 := w.obj.DescribeAccountKmsStatusWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeCertificate(ctx context.Context, request *kms.DescribeCertificateRequest) (*kms.DescribeCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.DescribeCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DescribeCertificate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DescribeCertificate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DescribeCertificate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeCertificate", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.DescribeCertificate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.DescribeCertificate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.DescribeCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DescribeCertificate(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeCertificateWithCallback(request *kms.DescribeCertificateRequest, callback func(response *kms.DescribeCertificateResponse, err error)) <-chan int {
	res0 := w.obj.DescribeCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeCertificateWithChan(request *kms.DescribeCertificateRequest) (<-chan *kms.DescribeCertificateResponse, <-chan error) {
	res0, res1 := w.obj.DescribeCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeKey(ctx context.Context, request *kms.DescribeKeyRequest) (*kms.DescribeKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.DescribeKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DescribeKey", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DescribeKey", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DescribeKey", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeKey", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.DescribeKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.DescribeKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.DescribeKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DescribeKey(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeKeyVersion(ctx context.Context, request *kms.DescribeKeyVersionRequest) (*kms.DescribeKeyVersionResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.DescribeKeyVersionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DescribeKeyVersion", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DescribeKeyVersion", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DescribeKeyVersion", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeKeyVersion", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.DescribeKeyVersion", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.DescribeKeyVersion", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.DescribeKeyVersion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DescribeKeyVersion(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeKeyVersionWithCallback(request *kms.DescribeKeyVersionRequest, callback func(response *kms.DescribeKeyVersionResponse, err error)) <-chan int {
	res0 := w.obj.DescribeKeyVersionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeKeyVersionWithChan(request *kms.DescribeKeyVersionRequest) (<-chan *kms.DescribeKeyVersionResponse, <-chan error) {
	res0, res1 := w.obj.DescribeKeyVersionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeKeyWithCallback(request *kms.DescribeKeyRequest, callback func(response *kms.DescribeKeyResponse, err error)) <-chan int {
	res0 := w.obj.DescribeKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeKeyWithChan(request *kms.DescribeKeyRequest) (<-chan *kms.DescribeKeyResponse, <-chan error) {
	res0, res1 := w.obj.DescribeKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeRegions(ctx context.Context, request *kms.DescribeRegionsRequest) (*kms.DescribeRegionsResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.DescribeRegionsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DescribeRegions", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DescribeRegions", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DescribeRegions", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeRegions", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.DescribeRegions", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.DescribeRegions", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.DescribeRegions", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DescribeRegions(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeRegionsWithCallback(request *kms.DescribeRegionsRequest, callback func(response *kms.DescribeRegionsResponse, err error)) <-chan int {
	res0 := w.obj.DescribeRegionsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeRegionsWithChan(request *kms.DescribeRegionsRequest) (<-chan *kms.DescribeRegionsResponse, <-chan error) {
	res0, res1 := w.obj.DescribeRegionsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeSecret(ctx context.Context, request *kms.DescribeSecretRequest) (*kms.DescribeSecretResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.DescribeSecretResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DescribeSecret", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DescribeSecret", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DescribeSecret", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeSecret", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.DescribeSecret", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.DescribeSecret", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.DescribeSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DescribeSecret(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeSecretWithCallback(request *kms.DescribeSecretRequest, callback func(response *kms.DescribeSecretResponse, err error)) <-chan int {
	res0 := w.obj.DescribeSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeSecretWithChan(request *kms.DescribeSecretRequest) (<-chan *kms.DescribeSecretResponse, <-chan error) {
	res0, res1 := w.obj.DescribeSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeService(ctx context.Context, request *kms.DescribeServiceRequest) (*kms.DescribeServiceResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.DescribeServiceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DescribeService", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DescribeService", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DescribeService", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeService", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.DescribeService", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.DescribeService", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.DescribeService", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DescribeService(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeServiceWithCallback(request *kms.DescribeServiceRequest, callback func(response *kms.DescribeServiceResponse, err error)) <-chan int {
	res0 := w.obj.DescribeServiceWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeServiceWithChan(request *kms.DescribeServiceRequest) (<-chan *kms.DescribeServiceResponse, <-chan error) {
	res0, res1 := w.obj.DescribeServiceWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DisableKey(ctx context.Context, request *kms.DisableKeyRequest) (*kms.DisableKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.DisableKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DisableKey", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DisableKey", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DisableKey", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.DisableKey", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.DisableKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.DisableKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.DisableKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DisableKey(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DisableKeyWithCallback(request *kms.DisableKeyRequest, callback func(response *kms.DisableKeyResponse, err error)) <-chan int {
	res0 := w.obj.DisableKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DisableKeyWithChan(request *kms.DisableKeyRequest) (<-chan *kms.DisableKeyResponse, <-chan error) {
	res0, res1 := w.obj.DisableKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) EnableKey(ctx context.Context, request *kms.EnableKeyRequest) (*kms.EnableKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.EnableKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.EnableKey", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.EnableKey", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.EnableKey", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.EnableKey", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.EnableKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.EnableKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.EnableKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.EnableKey(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) EnableKeyWithCallback(request *kms.EnableKeyRequest, callback func(response *kms.EnableKeyResponse, err error)) <-chan int {
	res0 := w.obj.EnableKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) EnableKeyWithChan(request *kms.EnableKeyRequest) (<-chan *kms.EnableKeyResponse, <-chan error) {
	res0, res1 := w.obj.EnableKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) Encrypt(ctx context.Context, request *kms.EncryptRequest) (*kms.EncryptResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.EncryptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.Encrypt", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.Encrypt", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.Encrypt", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.Encrypt", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.Encrypt", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.Encrypt", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.Encrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.Encrypt(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) EncryptWithCallback(request *kms.EncryptRequest, callback func(response *kms.EncryptResponse, err error)) <-chan int {
	res0 := w.obj.EncryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) EncryptWithChan(request *kms.EncryptRequest) (<-chan *kms.EncryptResponse, <-chan error) {
	res0, res1 := w.obj.EncryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ExportCertificate(ctx context.Context, request *kms.ExportCertificateRequest) (*kms.ExportCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ExportCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ExportCertificate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ExportCertificate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ExportCertificate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ExportCertificate", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ExportCertificate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ExportCertificate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ExportCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ExportCertificate(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ExportCertificateWithCallback(request *kms.ExportCertificateRequest, callback func(response *kms.ExportCertificateResponse, err error)) <-chan int {
	res0 := w.obj.ExportCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ExportCertificateWithChan(request *kms.ExportCertificateRequest) (<-chan *kms.ExportCertificateResponse, <-chan error) {
	res0, res1 := w.obj.ExportCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ExportDataKey(ctx context.Context, request *kms.ExportDataKeyRequest) (*kms.ExportDataKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ExportDataKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ExportDataKey", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ExportDataKey", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ExportDataKey", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ExportDataKey", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ExportDataKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ExportDataKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ExportDataKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ExportDataKey(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ExportDataKeyWithCallback(request *kms.ExportDataKeyRequest, callback func(response *kms.ExportDataKeyResponse, err error)) <-chan int {
	res0 := w.obj.ExportDataKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ExportDataKeyWithChan(request *kms.ExportDataKeyRequest) (<-chan *kms.ExportDataKeyResponse, <-chan error) {
	res0, res1 := w.obj.ExportDataKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GenerateAndExportDataKey(ctx context.Context, request *kms.GenerateAndExportDataKeyRequest) (*kms.GenerateAndExportDataKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.GenerateAndExportDataKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GenerateAndExportDataKey", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GenerateAndExportDataKey", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GenerateAndExportDataKey", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateAndExportDataKey", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.GenerateAndExportDataKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.GenerateAndExportDataKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.GenerateAndExportDataKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GenerateAndExportDataKey(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GenerateAndExportDataKeyWithCallback(request *kms.GenerateAndExportDataKeyRequest, callback func(response *kms.GenerateAndExportDataKeyResponse, err error)) <-chan int {
	res0 := w.obj.GenerateAndExportDataKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GenerateAndExportDataKeyWithChan(request *kms.GenerateAndExportDataKeyRequest) (<-chan *kms.GenerateAndExportDataKeyResponse, <-chan error) {
	res0, res1 := w.obj.GenerateAndExportDataKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GenerateDataKey(ctx context.Context, request *kms.GenerateDataKeyRequest) (*kms.GenerateDataKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.GenerateDataKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GenerateDataKey", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GenerateDataKey", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GenerateDataKey", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateDataKey", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.GenerateDataKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.GenerateDataKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.GenerateDataKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GenerateDataKey(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GenerateDataKeyWithCallback(request *kms.GenerateDataKeyRequest, callback func(response *kms.GenerateDataKeyResponse, err error)) <-chan int {
	res0 := w.obj.GenerateDataKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GenerateDataKeyWithChan(request *kms.GenerateDataKeyRequest) (<-chan *kms.GenerateDataKeyResponse, <-chan error) {
	res0, res1 := w.obj.GenerateDataKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GenerateDataKeyWithoutPlaintext(ctx context.Context, request *kms.GenerateDataKeyWithoutPlaintextRequest) (*kms.GenerateDataKeyWithoutPlaintextResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.GenerateDataKeyWithoutPlaintextResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GenerateDataKeyWithoutPlaintext", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GenerateDataKeyWithoutPlaintext", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GenerateDataKeyWithoutPlaintext", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateDataKeyWithoutPlaintext", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.GenerateDataKeyWithoutPlaintext", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.GenerateDataKeyWithoutPlaintext", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.GenerateDataKeyWithoutPlaintext", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GenerateDataKeyWithoutPlaintext(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GenerateDataKeyWithoutPlaintextWithCallback(request *kms.GenerateDataKeyWithoutPlaintextRequest, callback func(response *kms.GenerateDataKeyWithoutPlaintextResponse, err error)) <-chan int {
	res0 := w.obj.GenerateDataKeyWithoutPlaintextWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GenerateDataKeyWithoutPlaintextWithChan(request *kms.GenerateDataKeyWithoutPlaintextRequest) (<-chan *kms.GenerateDataKeyWithoutPlaintextResponse, <-chan error) {
	res0, res1 := w.obj.GenerateDataKeyWithoutPlaintextWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetCertificate(ctx context.Context, request *kms.GetCertificateRequest) (*kms.GetCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.GetCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetCertificate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetCertificate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetCertificate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.GetCertificate", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.GetCertificate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.GetCertificate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.GetCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetCertificate(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetCertificateWithCallback(request *kms.GetCertificateRequest, callback func(response *kms.GetCertificateResponse, err error)) <-chan int {
	res0 := w.obj.GetCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetCertificateWithChan(request *kms.GetCertificateRequest) (<-chan *kms.GetCertificateResponse, <-chan error) {
	res0, res1 := w.obj.GetCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetParametersForImport(ctx context.Context, request *kms.GetParametersForImportRequest) (*kms.GetParametersForImportResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.GetParametersForImportResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetParametersForImport", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetParametersForImport", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetParametersForImport", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.GetParametersForImport", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.GetParametersForImport", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.GetParametersForImport", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.GetParametersForImport", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetParametersForImport(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetParametersForImportWithCallback(request *kms.GetParametersForImportRequest, callback func(response *kms.GetParametersForImportResponse, err error)) <-chan int {
	res0 := w.obj.GetParametersForImportWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetParametersForImportWithChan(request *kms.GetParametersForImportRequest) (<-chan *kms.GetParametersForImportResponse, <-chan error) {
	res0, res1 := w.obj.GetParametersForImportWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetPublicKey(ctx context.Context, request *kms.GetPublicKeyRequest) (*kms.GetPublicKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.GetPublicKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetPublicKey", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetPublicKey", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetPublicKey", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.GetPublicKey", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.GetPublicKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.GetPublicKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.GetPublicKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetPublicKey(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetPublicKeyWithCallback(request *kms.GetPublicKeyRequest, callback func(response *kms.GetPublicKeyResponse, err error)) <-chan int {
	res0 := w.obj.GetPublicKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetPublicKeyWithChan(request *kms.GetPublicKeyRequest) (<-chan *kms.GetPublicKeyResponse, <-chan error) {
	res0, res1 := w.obj.GetPublicKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetRandomPassword(ctx context.Context, request *kms.GetRandomPasswordRequest) (*kms.GetRandomPasswordResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.GetRandomPasswordResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetRandomPassword", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetRandomPassword", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetRandomPassword", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.GetRandomPassword", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.GetRandomPassword", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.GetRandomPassword", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.GetRandomPassword", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetRandomPassword(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetRandomPasswordWithCallback(request *kms.GetRandomPasswordRequest, callback func(response *kms.GetRandomPasswordResponse, err error)) <-chan int {
	res0 := w.obj.GetRandomPasswordWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetRandomPasswordWithChan(request *kms.GetRandomPasswordRequest) (<-chan *kms.GetRandomPasswordResponse, <-chan error) {
	res0, res1 := w.obj.GetRandomPasswordWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetSecretValue(ctx context.Context, request *kms.GetSecretValueRequest) (*kms.GetSecretValueResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.GetSecretValueResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetSecretValue", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetSecretValue", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetSecretValue", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.GetSecretValue", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.GetSecretValue", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.GetSecretValue", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.GetSecretValue", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetSecretValue(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetSecretValueWithCallback(request *kms.GetSecretValueRequest, callback func(response *kms.GetSecretValueResponse, err error)) <-chan int {
	res0 := w.obj.GetSecretValueWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetSecretValueWithChan(request *kms.GetSecretValueRequest) (<-chan *kms.GetSecretValueResponse, <-chan error) {
	res0, res1 := w.obj.GetSecretValueWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ImportCertificate(ctx context.Context, request *kms.ImportCertificateRequest) (*kms.ImportCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ImportCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ImportCertificate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ImportCertificate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ImportCertificate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ImportCertificate", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ImportCertificate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ImportCertificate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ImportCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ImportCertificate(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ImportCertificateWithCallback(request *kms.ImportCertificateRequest, callback func(response *kms.ImportCertificateResponse, err error)) <-chan int {
	res0 := w.obj.ImportCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ImportCertificateWithChan(request *kms.ImportCertificateRequest) (<-chan *kms.ImportCertificateResponse, <-chan error) {
	res0, res1 := w.obj.ImportCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ImportEncryptionCertificate(ctx context.Context, request *kms.ImportEncryptionCertificateRequest) (*kms.ImportEncryptionCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ImportEncryptionCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ImportEncryptionCertificate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ImportEncryptionCertificate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ImportEncryptionCertificate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ImportEncryptionCertificate", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ImportEncryptionCertificate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ImportEncryptionCertificate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ImportEncryptionCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ImportEncryptionCertificate(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ImportEncryptionCertificateWithCallback(request *kms.ImportEncryptionCertificateRequest, callback func(response *kms.ImportEncryptionCertificateResponse, err error)) <-chan int {
	res0 := w.obj.ImportEncryptionCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ImportEncryptionCertificateWithChan(request *kms.ImportEncryptionCertificateRequest) (<-chan *kms.ImportEncryptionCertificateResponse, <-chan error) {
	res0, res1 := w.obj.ImportEncryptionCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ImportKeyMaterial(ctx context.Context, request *kms.ImportKeyMaterialRequest) (*kms.ImportKeyMaterialResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ImportKeyMaterialResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ImportKeyMaterial", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ImportKeyMaterial", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ImportKeyMaterial", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ImportKeyMaterial", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ImportKeyMaterial", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ImportKeyMaterial", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ImportKeyMaterial", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ImportKeyMaterial(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ImportKeyMaterialWithCallback(request *kms.ImportKeyMaterialRequest, callback func(response *kms.ImportKeyMaterialResponse, err error)) <-chan int {
	res0 := w.obj.ImportKeyMaterialWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ImportKeyMaterialWithChan(request *kms.ImportKeyMaterialRequest) (<-chan *kms.ImportKeyMaterialResponse, <-chan error) {
	res0, res1 := w.obj.ImportKeyMaterialWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListAliases(ctx context.Context, request *kms.ListAliasesRequest) (*kms.ListAliasesResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ListAliasesResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListAliases", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListAliases", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListAliases", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ListAliases", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ListAliases", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ListAliases", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ListAliases", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListAliases(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListAliasesByKeyId(ctx context.Context, request *kms.ListAliasesByKeyIdRequest) (*kms.ListAliasesByKeyIdResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ListAliasesByKeyIdResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListAliasesByKeyId", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListAliasesByKeyId", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListAliasesByKeyId", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ListAliasesByKeyId", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ListAliasesByKeyId", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ListAliasesByKeyId", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ListAliasesByKeyId", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListAliasesByKeyId(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListAliasesByKeyIdWithCallback(request *kms.ListAliasesByKeyIdRequest, callback func(response *kms.ListAliasesByKeyIdResponse, err error)) <-chan int {
	res0 := w.obj.ListAliasesByKeyIdWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListAliasesByKeyIdWithChan(request *kms.ListAliasesByKeyIdRequest) (<-chan *kms.ListAliasesByKeyIdResponse, <-chan error) {
	res0, res1 := w.obj.ListAliasesByKeyIdWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListAliasesWithCallback(request *kms.ListAliasesRequest, callback func(response *kms.ListAliasesResponse, err error)) <-chan int {
	res0 := w.obj.ListAliasesWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListAliasesWithChan(request *kms.ListAliasesRequest) (<-chan *kms.ListAliasesResponse, <-chan error) {
	res0, res1 := w.obj.ListAliasesWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListCertificates(ctx context.Context, request *kms.ListCertificatesRequest) (*kms.ListCertificatesResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ListCertificatesResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListCertificates", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListCertificates", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListCertificates", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ListCertificates", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ListCertificates", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ListCertificates", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ListCertificates", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListCertificates(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListCertificatesWithCallback(request *kms.ListCertificatesRequest, callback func(response *kms.ListCertificatesResponse, err error)) <-chan int {
	res0 := w.obj.ListCertificatesWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListCertificatesWithChan(request *kms.ListCertificatesRequest) (<-chan *kms.ListCertificatesResponse, <-chan error) {
	res0, res1 := w.obj.ListCertificatesWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListKeyVersions(ctx context.Context, request *kms.ListKeyVersionsRequest) (*kms.ListKeyVersionsResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ListKeyVersionsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListKeyVersions", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListKeyVersions", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListKeyVersions", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ListKeyVersions", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ListKeyVersions", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ListKeyVersions", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ListKeyVersions", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListKeyVersions(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListKeyVersionsWithCallback(request *kms.ListKeyVersionsRequest, callback func(response *kms.ListKeyVersionsResponse, err error)) <-chan int {
	res0 := w.obj.ListKeyVersionsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListKeyVersionsWithChan(request *kms.ListKeyVersionsRequest) (<-chan *kms.ListKeyVersionsResponse, <-chan error) {
	res0, res1 := w.obj.ListKeyVersionsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListKeys(ctx context.Context, request *kms.ListKeysRequest) (*kms.ListKeysResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ListKeysResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListKeys", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListKeys", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListKeys", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ListKeys", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ListKeys", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ListKeys", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ListKeys", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListKeys(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListKeysWithCallback(request *kms.ListKeysRequest, callback func(response *kms.ListKeysResponse, err error)) <-chan int {
	res0 := w.obj.ListKeysWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListKeysWithChan(request *kms.ListKeysRequest) (<-chan *kms.ListKeysResponse, <-chan error) {
	res0, res1 := w.obj.ListKeysWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListResourceTags(ctx context.Context, request *kms.ListResourceTagsRequest) (*kms.ListResourceTagsResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ListResourceTagsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListResourceTags", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListResourceTags", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListResourceTags", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ListResourceTags", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ListResourceTags", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ListResourceTags", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ListResourceTags", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListResourceTags(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListResourceTagsWithCallback(request *kms.ListResourceTagsRequest, callback func(response *kms.ListResourceTagsResponse, err error)) <-chan int {
	res0 := w.obj.ListResourceTagsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListResourceTagsWithChan(request *kms.ListResourceTagsRequest) (<-chan *kms.ListResourceTagsResponse, <-chan error) {
	res0, res1 := w.obj.ListResourceTagsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListSecretVersionIds(ctx context.Context, request *kms.ListSecretVersionIdsRequest) (*kms.ListSecretVersionIdsResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ListSecretVersionIdsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListSecretVersionIds", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListSecretVersionIds", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListSecretVersionIds", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ListSecretVersionIds", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ListSecretVersionIds", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ListSecretVersionIds", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ListSecretVersionIds", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListSecretVersionIds(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListSecretVersionIdsWithCallback(request *kms.ListSecretVersionIdsRequest, callback func(response *kms.ListSecretVersionIdsResponse, err error)) <-chan int {
	res0 := w.obj.ListSecretVersionIdsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListSecretVersionIdsWithChan(request *kms.ListSecretVersionIdsRequest) (<-chan *kms.ListSecretVersionIdsResponse, <-chan error) {
	res0, res1 := w.obj.ListSecretVersionIdsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListSecrets(ctx context.Context, request *kms.ListSecretsRequest) (*kms.ListSecretsResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ListSecretsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListSecrets", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListSecrets", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListSecrets", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ListSecrets", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ListSecrets", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ListSecrets", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ListSecrets", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListSecrets(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListSecretsWithCallback(request *kms.ListSecretsRequest, callback func(response *kms.ListSecretsResponse, err error)) <-chan int {
	res0 := w.obj.ListSecretsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListSecretsWithChan(request *kms.ListSecretsRequest) (<-chan *kms.ListSecretsResponse, <-chan error) {
	res0, res1 := w.obj.ListSecretsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) OpenKmsService(ctx context.Context, request *kms.OpenKmsServiceRequest) (*kms.OpenKmsServiceResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.OpenKmsServiceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.OpenKmsService", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.OpenKmsService", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.OpenKmsService", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.OpenKmsService", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.OpenKmsService", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.OpenKmsService", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.OpenKmsService", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.OpenKmsService(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) OpenKmsServiceWithCallback(request *kms.OpenKmsServiceRequest, callback func(response *kms.OpenKmsServiceResponse, err error)) <-chan int {
	res0 := w.obj.OpenKmsServiceWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) OpenKmsServiceWithChan(request *kms.OpenKmsServiceRequest) (<-chan *kms.OpenKmsServiceResponse, <-chan error) {
	res0, res1 := w.obj.OpenKmsServiceWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) PutSecretValue(ctx context.Context, request *kms.PutSecretValueRequest) (*kms.PutSecretValueResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.PutSecretValueResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.PutSecretValue", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.PutSecretValue", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.PutSecretValue", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.PutSecretValue", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.PutSecretValue", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.PutSecretValue", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.PutSecretValue", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.PutSecretValue(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) PutSecretValueWithCallback(request *kms.PutSecretValueRequest, callback func(response *kms.PutSecretValueResponse, err error)) <-chan int {
	res0 := w.obj.PutSecretValueWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) PutSecretValueWithChan(request *kms.PutSecretValueRequest) (<-chan *kms.PutSecretValueResponse, <-chan error) {
	res0, res1 := w.obj.PutSecretValueWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ReEncrypt(ctx context.Context, request *kms.ReEncryptRequest) (*kms.ReEncryptResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ReEncryptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ReEncrypt", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ReEncrypt", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ReEncrypt", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ReEncrypt", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ReEncrypt", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ReEncrypt", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ReEncrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ReEncrypt(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ReEncryptWithCallback(request *kms.ReEncryptRequest, callback func(response *kms.ReEncryptResponse, err error)) <-chan int {
	res0 := w.obj.ReEncryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ReEncryptWithChan(request *kms.ReEncryptRequest) (<-chan *kms.ReEncryptResponse, <-chan error) {
	res0, res1 := w.obj.ReEncryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) RestoreSecret(ctx context.Context, request *kms.RestoreSecretRequest) (*kms.RestoreSecretResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.RestoreSecretResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.RestoreSecret", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.RestoreSecret", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.RestoreSecret", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.RestoreSecret", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.RestoreSecret", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.RestoreSecret", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.RestoreSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.RestoreSecret(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) RestoreSecretWithCallback(request *kms.RestoreSecretRequest, callback func(response *kms.RestoreSecretResponse, err error)) <-chan int {
	res0 := w.obj.RestoreSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) RestoreSecretWithChan(request *kms.RestoreSecretRequest) (<-chan *kms.RestoreSecretResponse, <-chan error) {
	res0, res1 := w.obj.RestoreSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) RotateSecret(ctx context.Context, request *kms.RotateSecretRequest) (*kms.RotateSecretResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.RotateSecretResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.RotateSecret", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.RotateSecret", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.RotateSecret", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.RotateSecret", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.RotateSecret", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.RotateSecret", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.RotateSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.RotateSecret(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) RotateSecretWithCallback(request *kms.RotateSecretRequest, callback func(response *kms.RotateSecretResponse, err error)) <-chan int {
	res0 := w.obj.RotateSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) RotateSecretWithChan(request *kms.RotateSecretRequest) (<-chan *kms.RotateSecretResponse, <-chan error) {
	res0, res1 := w.obj.RotateSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ScheduleKeyDeletion(ctx context.Context, request *kms.ScheduleKeyDeletionRequest) (*kms.ScheduleKeyDeletionResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.ScheduleKeyDeletionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ScheduleKeyDeletion", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ScheduleKeyDeletion", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ScheduleKeyDeletion", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.ScheduleKeyDeletion", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.ScheduleKeyDeletion", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.ScheduleKeyDeletion", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.ScheduleKeyDeletion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ScheduleKeyDeletion(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ScheduleKeyDeletionWithCallback(request *kms.ScheduleKeyDeletionRequest, callback func(response *kms.ScheduleKeyDeletionResponse, err error)) <-chan int {
	res0 := w.obj.ScheduleKeyDeletionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ScheduleKeyDeletionWithChan(request *kms.ScheduleKeyDeletionRequest) (<-chan *kms.ScheduleKeyDeletionResponse, <-chan error) {
	res0, res1 := w.obj.ScheduleKeyDeletionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) TagResource(ctx context.Context, request *kms.TagResourceRequest) (*kms.TagResourceResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.TagResourceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.TagResource", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.TagResource", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.TagResource", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.TagResource", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.TagResource", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.TagResource", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.TagResource", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.TagResource(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) TagResourceWithCallback(request *kms.TagResourceRequest, callback func(response *kms.TagResourceResponse, err error)) <-chan int {
	res0 := w.obj.TagResourceWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) TagResourceWithChan(request *kms.TagResourceRequest) (<-chan *kms.TagResourceResponse, <-chan error) {
	res0, res1 := w.obj.TagResourceWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UntagResource(ctx context.Context, request *kms.UntagResourceRequest) (*kms.UntagResourceResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.UntagResourceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UntagResource", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UntagResource", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UntagResource", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.UntagResource", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.UntagResource", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.UntagResource", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.UntagResource", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UntagResource(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UntagResourceWithCallback(request *kms.UntagResourceRequest, callback func(response *kms.UntagResourceResponse, err error)) <-chan int {
	res0 := w.obj.UntagResourceWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UntagResourceWithChan(request *kms.UntagResourceRequest) (<-chan *kms.UntagResourceResponse, <-chan error) {
	res0, res1 := w.obj.UntagResourceWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateAlias(ctx context.Context, request *kms.UpdateAliasRequest) (*kms.UpdateAliasResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.UpdateAliasResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateAlias", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateAlias", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateAlias", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateAlias", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.UpdateAlias", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.UpdateAlias", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.UpdateAlias", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateAlias(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateAliasWithCallback(request *kms.UpdateAliasRequest, callback func(response *kms.UpdateAliasResponse, err error)) <-chan int {
	res0 := w.obj.UpdateAliasWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateAliasWithChan(request *kms.UpdateAliasRequest) (<-chan *kms.UpdateAliasResponse, <-chan error) {
	res0, res1 := w.obj.UpdateAliasWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateCertificateStatus(ctx context.Context, request *kms.UpdateCertificateStatusRequest) (*kms.UpdateCertificateStatusResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.UpdateCertificateStatusResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateCertificateStatus", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateCertificateStatus", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateCertificateStatus", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateCertificateStatus", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.UpdateCertificateStatus", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.UpdateCertificateStatus", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.UpdateCertificateStatus", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateCertificateStatus(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateCertificateStatusWithCallback(request *kms.UpdateCertificateStatusRequest, callback func(response *kms.UpdateCertificateStatusResponse, err error)) <-chan int {
	res0 := w.obj.UpdateCertificateStatusWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateCertificateStatusWithChan(request *kms.UpdateCertificateStatusRequest) (<-chan *kms.UpdateCertificateStatusResponse, <-chan error) {
	res0, res1 := w.obj.UpdateCertificateStatusWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateKeyDescription(ctx context.Context, request *kms.UpdateKeyDescriptionRequest) (*kms.UpdateKeyDescriptionResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.UpdateKeyDescriptionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateKeyDescription", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateKeyDescription", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateKeyDescription", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateKeyDescription", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.UpdateKeyDescription", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.UpdateKeyDescription", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.UpdateKeyDescription", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateKeyDescription(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateKeyDescriptionWithCallback(request *kms.UpdateKeyDescriptionRequest, callback func(response *kms.UpdateKeyDescriptionResponse, err error)) <-chan int {
	res0 := w.obj.UpdateKeyDescriptionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateKeyDescriptionWithChan(request *kms.UpdateKeyDescriptionRequest) (<-chan *kms.UpdateKeyDescriptionResponse, <-chan error) {
	res0, res1 := w.obj.UpdateKeyDescriptionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateRotationPolicy(ctx context.Context, request *kms.UpdateRotationPolicyRequest) (*kms.UpdateRotationPolicyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.UpdateRotationPolicyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateRotationPolicy", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateRotationPolicy", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateRotationPolicy", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateRotationPolicy", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.UpdateRotationPolicy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.UpdateRotationPolicy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.UpdateRotationPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateRotationPolicy(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateRotationPolicyWithCallback(request *kms.UpdateRotationPolicyRequest, callback func(response *kms.UpdateRotationPolicyResponse, err error)) <-chan int {
	res0 := w.obj.UpdateRotationPolicyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateRotationPolicyWithChan(request *kms.UpdateRotationPolicyRequest) (<-chan *kms.UpdateRotationPolicyResponse, <-chan error) {
	res0, res1 := w.obj.UpdateRotationPolicyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateSecret(ctx context.Context, request *kms.UpdateSecretRequest) (*kms.UpdateSecretResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.UpdateSecretResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateSecret", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateSecret", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateSecret", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateSecret", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.UpdateSecret", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.UpdateSecret", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.UpdateSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateSecret(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateSecretRotationPolicy(ctx context.Context, request *kms.UpdateSecretRotationPolicyRequest) (*kms.UpdateSecretRotationPolicyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.UpdateSecretRotationPolicyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateSecretRotationPolicy", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateSecretRotationPolicy", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateSecretRotationPolicy", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateSecretRotationPolicy", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.UpdateSecretRotationPolicy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.UpdateSecretRotationPolicy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.UpdateSecretRotationPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateSecretRotationPolicy(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateSecretRotationPolicyWithCallback(request *kms.UpdateSecretRotationPolicyRequest, callback func(response *kms.UpdateSecretRotationPolicyResponse, err error)) <-chan int {
	res0 := w.obj.UpdateSecretRotationPolicyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateSecretRotationPolicyWithChan(request *kms.UpdateSecretRotationPolicyRequest) (<-chan *kms.UpdateSecretRotationPolicyResponse, <-chan error) {
	res0, res1 := w.obj.UpdateSecretRotationPolicyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateSecretVersionStage(ctx context.Context, request *kms.UpdateSecretVersionStageRequest) (*kms.UpdateSecretVersionStageResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.UpdateSecretVersionStageResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateSecretVersionStage", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateSecretVersionStage", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateSecretVersionStage", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateSecretVersionStage", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.UpdateSecretVersionStage", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.UpdateSecretVersionStage", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.UpdateSecretVersionStage", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateSecretVersionStage(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateSecretVersionStageWithCallback(request *kms.UpdateSecretVersionStageRequest, callback func(response *kms.UpdateSecretVersionStageResponse, err error)) <-chan int {
	res0 := w.obj.UpdateSecretVersionStageWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateSecretVersionStageWithChan(request *kms.UpdateSecretVersionStageRequest) (<-chan *kms.UpdateSecretVersionStageResponse, <-chan error) {
	res0, res1 := w.obj.UpdateSecretVersionStageWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateSecretWithCallback(request *kms.UpdateSecretRequest, callback func(response *kms.UpdateSecretResponse, err error)) <-chan int {
	res0 := w.obj.UpdateSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateSecretWithChan(request *kms.UpdateSecretRequest) (<-chan *kms.UpdateSecretResponse, <-chan error) {
	res0, res1 := w.obj.UpdateSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UploadCertificate(ctx context.Context, request *kms.UploadCertificateRequest) (*kms.UploadCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *kms.UploadCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UploadCertificate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UploadCertificate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UploadCertificate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "kms.Client.UploadCertificate", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("kms.Client.UploadCertificate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("kms.Client.UploadCertificate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("kms.Client.UploadCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UploadCertificate(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UploadCertificateWithCallback(request *kms.UploadCertificateRequest, callback func(response *kms.UploadCertificateResponse, err error)) <-chan int {
	res0 := w.obj.UploadCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UploadCertificateWithChan(request *kms.UploadCertificateRequest) (<-chan *kms.UploadCertificateResponse, <-chan error) {
	res0, res1 := w.obj.UploadCertificateWithChan(request)
	return res0, res1
}
