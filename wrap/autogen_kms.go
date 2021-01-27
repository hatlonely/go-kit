// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type KMSClientWrapper struct {
	obj              *kms.Client
	retry            *Retry
	options          *WrapperOptions
	durationMetric   *prometheus.HistogramVec
	totalMetric      *prometheus.CounterVec
	rateLimiterGroup RateLimiterGroup
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

func (w *KMSClientWrapper) OnRateLimiterGroupChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options RateLimiterGroupOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		rateLimiterGroup, err := NewRateLimiterGroup(&options)
		if err != nil {
			return errors.Wrap(err, "NewRateLimiterGroup failed")
		}
		w.rateLimiterGroup = rateLimiterGroup
		return nil
	}
}

func (w *KMSClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "kms_Client_durationMs",
		Help:        "kms Client response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "kms_Client_total",
		Help:        "kms Client request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
}

func (w *KMSClientWrapper) AsymmetricDecrypt(ctx context.Context, request *kms.AsymmetricDecryptRequest) (*kms.AsymmetricDecryptResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricDecrypt")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.AsymmetricDecryptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.AsymmetricDecrypt"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.AsymmetricDecrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.AsymmetricDecrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AsymmetricDecrypt(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) AsymmetricDecryptWithCallback(ctx context.Context, request *kms.AsymmetricDecryptRequest, callback func(response *kms.AsymmetricDecryptResponse, err error)) <-chan int {
	res0 := w.obj.AsymmetricDecryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) AsymmetricDecryptWithChan(ctx context.Context, request *kms.AsymmetricDecryptRequest) (<-chan *kms.AsymmetricDecryptResponse, <-chan error) {
	res0, res1 := w.obj.AsymmetricDecryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) AsymmetricEncrypt(ctx context.Context, request *kms.AsymmetricEncryptRequest) (*kms.AsymmetricEncryptResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricEncrypt")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.AsymmetricEncryptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.AsymmetricEncrypt"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.AsymmetricEncrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.AsymmetricEncrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AsymmetricEncrypt(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) AsymmetricEncryptWithCallback(ctx context.Context, request *kms.AsymmetricEncryptRequest, callback func(response *kms.AsymmetricEncryptResponse, err error)) <-chan int {
	res0 := w.obj.AsymmetricEncryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) AsymmetricEncryptWithChan(ctx context.Context, request *kms.AsymmetricEncryptRequest) (<-chan *kms.AsymmetricEncryptResponse, <-chan error) {
	res0, res1 := w.obj.AsymmetricEncryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) AsymmetricSign(ctx context.Context, request *kms.AsymmetricSignRequest) (*kms.AsymmetricSignResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricSign")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.AsymmetricSignResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.AsymmetricSign"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.AsymmetricSign", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.AsymmetricSign", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AsymmetricSign(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) AsymmetricSignWithCallback(ctx context.Context, request *kms.AsymmetricSignRequest, callback func(response *kms.AsymmetricSignResponse, err error)) <-chan int {
	res0 := w.obj.AsymmetricSignWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) AsymmetricSignWithChan(ctx context.Context, request *kms.AsymmetricSignRequest) (<-chan *kms.AsymmetricSignResponse, <-chan error) {
	res0, res1 := w.obj.AsymmetricSignWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) AsymmetricVerify(ctx context.Context, request *kms.AsymmetricVerifyRequest) (*kms.AsymmetricVerifyResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricVerify")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.AsymmetricVerifyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.AsymmetricVerify"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.AsymmetricVerify", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.AsymmetricVerify", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AsymmetricVerify(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) AsymmetricVerifyWithCallback(ctx context.Context, request *kms.AsymmetricVerifyRequest, callback func(response *kms.AsymmetricVerifyResponse, err error)) <-chan int {
	res0 := w.obj.AsymmetricVerifyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) AsymmetricVerifyWithChan(ctx context.Context, request *kms.AsymmetricVerifyRequest) (<-chan *kms.AsymmetricVerifyResponse, <-chan error) {
	res0, res1 := w.obj.AsymmetricVerifyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CancelKeyDeletion(ctx context.Context, request *kms.CancelKeyDeletionRequest) (*kms.CancelKeyDeletionResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CancelKeyDeletion")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.CancelKeyDeletionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.CancelKeyDeletion"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.CancelKeyDeletion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.CancelKeyDeletion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CancelKeyDeletion(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CancelKeyDeletionWithCallback(ctx context.Context, request *kms.CancelKeyDeletionRequest, callback func(response *kms.CancelKeyDeletionResponse, err error)) <-chan int {
	res0 := w.obj.CancelKeyDeletionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CancelKeyDeletionWithChan(ctx context.Context, request *kms.CancelKeyDeletionRequest) (<-chan *kms.CancelKeyDeletionResponse, <-chan error) {
	res0, res1 := w.obj.CancelKeyDeletionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CertificatePrivateKeyDecrypt(ctx context.Context, request *kms.CertificatePrivateKeyDecryptRequest) (*kms.CertificatePrivateKeyDecryptResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePrivateKeyDecrypt")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.CertificatePrivateKeyDecryptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.CertificatePrivateKeyDecrypt"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.CertificatePrivateKeyDecrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.CertificatePrivateKeyDecrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CertificatePrivateKeyDecrypt(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CertificatePrivateKeyDecryptWithCallback(ctx context.Context, request *kms.CertificatePrivateKeyDecryptRequest, callback func(response *kms.CertificatePrivateKeyDecryptResponse, err error)) <-chan int {
	res0 := w.obj.CertificatePrivateKeyDecryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CertificatePrivateKeyDecryptWithChan(ctx context.Context, request *kms.CertificatePrivateKeyDecryptRequest) (<-chan *kms.CertificatePrivateKeyDecryptResponse, <-chan error) {
	res0, res1 := w.obj.CertificatePrivateKeyDecryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CertificatePrivateKeySign(ctx context.Context, request *kms.CertificatePrivateKeySignRequest) (*kms.CertificatePrivateKeySignResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePrivateKeySign")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.CertificatePrivateKeySignResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.CertificatePrivateKeySign"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.CertificatePrivateKeySign", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.CertificatePrivateKeySign", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CertificatePrivateKeySign(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CertificatePrivateKeySignWithCallback(ctx context.Context, request *kms.CertificatePrivateKeySignRequest, callback func(response *kms.CertificatePrivateKeySignResponse, err error)) <-chan int {
	res0 := w.obj.CertificatePrivateKeySignWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CertificatePrivateKeySignWithChan(ctx context.Context, request *kms.CertificatePrivateKeySignRequest) (<-chan *kms.CertificatePrivateKeySignResponse, <-chan error) {
	res0, res1 := w.obj.CertificatePrivateKeySignWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CertificatePublicKeyEncrypt(ctx context.Context, request *kms.CertificatePublicKeyEncryptRequest) (*kms.CertificatePublicKeyEncryptResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePublicKeyEncrypt")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.CertificatePublicKeyEncryptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.CertificatePublicKeyEncrypt"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.CertificatePublicKeyEncrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.CertificatePublicKeyEncrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CertificatePublicKeyEncrypt(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CertificatePublicKeyEncryptWithCallback(ctx context.Context, request *kms.CertificatePublicKeyEncryptRequest, callback func(response *kms.CertificatePublicKeyEncryptResponse, err error)) <-chan int {
	res0 := w.obj.CertificatePublicKeyEncryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CertificatePublicKeyEncryptWithChan(ctx context.Context, request *kms.CertificatePublicKeyEncryptRequest) (<-chan *kms.CertificatePublicKeyEncryptResponse, <-chan error) {
	res0, res1 := w.obj.CertificatePublicKeyEncryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CertificatePublicKeyVerify(ctx context.Context, request *kms.CertificatePublicKeyVerifyRequest) (*kms.CertificatePublicKeyVerifyResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePublicKeyVerify")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.CertificatePublicKeyVerifyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.CertificatePublicKeyVerify"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.CertificatePublicKeyVerify", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.CertificatePublicKeyVerify", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CertificatePublicKeyVerify(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CertificatePublicKeyVerifyWithCallback(ctx context.Context, request *kms.CertificatePublicKeyVerifyRequest, callback func(response *kms.CertificatePublicKeyVerifyResponse, err error)) <-chan int {
	res0 := w.obj.CertificatePublicKeyVerifyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CertificatePublicKeyVerifyWithChan(ctx context.Context, request *kms.CertificatePublicKeyVerifyRequest) (<-chan *kms.CertificatePublicKeyVerifyResponse, <-chan error) {
	res0, res1 := w.obj.CertificatePublicKeyVerifyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateAlias(ctx context.Context, request *kms.CreateAliasRequest) (*kms.CreateAliasResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateAlias")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.CreateAliasResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.CreateAlias"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.CreateAlias", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.CreateAlias", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateAlias(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateAliasWithCallback(ctx context.Context, request *kms.CreateAliasRequest, callback func(response *kms.CreateAliasResponse, err error)) <-chan int {
	res0 := w.obj.CreateAliasWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateAliasWithChan(ctx context.Context, request *kms.CreateAliasRequest) (<-chan *kms.CreateAliasResponse, <-chan error) {
	res0, res1 := w.obj.CreateAliasWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateCertificate(ctx context.Context, request *kms.CreateCertificateRequest) (*kms.CreateCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateCertificate")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.CreateCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.CreateCertificate"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.CreateCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.CreateCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateCertificateWithCallback(ctx context.Context, request *kms.CreateCertificateRequest, callback func(response *kms.CreateCertificateResponse, err error)) <-chan int {
	res0 := w.obj.CreateCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateCertificateWithChan(ctx context.Context, request *kms.CreateCertificateRequest) (<-chan *kms.CreateCertificateResponse, <-chan error) {
	res0, res1 := w.obj.CreateCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateKey(ctx context.Context, request *kms.CreateKeyRequest) (*kms.CreateKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateKey")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.CreateKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.CreateKey"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.CreateKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.CreateKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateKeyVersion(ctx context.Context, request *kms.CreateKeyVersionRequest) (*kms.CreateKeyVersionResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateKeyVersion")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.CreateKeyVersionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.CreateKeyVersion"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.CreateKeyVersion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.CreateKeyVersion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateKeyVersion(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateKeyVersionWithCallback(ctx context.Context, request *kms.CreateKeyVersionRequest, callback func(response *kms.CreateKeyVersionResponse, err error)) <-chan int {
	res0 := w.obj.CreateKeyVersionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateKeyVersionWithChan(ctx context.Context, request *kms.CreateKeyVersionRequest) (<-chan *kms.CreateKeyVersionResponse, <-chan error) {
	res0, res1 := w.obj.CreateKeyVersionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateKeyWithCallback(ctx context.Context, request *kms.CreateKeyRequest, callback func(response *kms.CreateKeyResponse, err error)) <-chan int {
	res0 := w.obj.CreateKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateKeyWithChan(ctx context.Context, request *kms.CreateKeyRequest) (<-chan *kms.CreateKeyResponse, <-chan error) {
	res0, res1 := w.obj.CreateKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateSecret(ctx context.Context, request *kms.CreateSecretRequest) (*kms.CreateSecretResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateSecret")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.CreateSecretResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.CreateSecret"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.CreateSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.CreateSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateSecret(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateSecretWithCallback(ctx context.Context, request *kms.CreateSecretRequest, callback func(response *kms.CreateSecretResponse, err error)) <-chan int {
	res0 := w.obj.CreateSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateSecretWithChan(ctx context.Context, request *kms.CreateSecretRequest) (<-chan *kms.CreateSecretResponse, <-chan error) {
	res0, res1 := w.obj.CreateSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) Decrypt(ctx context.Context, request *kms.DecryptRequest) (*kms.DecryptResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.Decrypt")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.DecryptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.Decrypt"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.Decrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.Decrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.Decrypt(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DecryptWithCallback(ctx context.Context, request *kms.DecryptRequest, callback func(response *kms.DecryptResponse, err error)) <-chan int {
	res0 := w.obj.DecryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DecryptWithChan(ctx context.Context, request *kms.DecryptRequest) (<-chan *kms.DecryptResponse, <-chan error) {
	res0, res1 := w.obj.DecryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DeleteAlias(ctx context.Context, request *kms.DeleteAliasRequest) (*kms.DeleteAliasResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteAlias")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.DeleteAliasResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DeleteAlias"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.DeleteAlias", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.DeleteAlias", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteAlias(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DeleteAliasWithCallback(ctx context.Context, request *kms.DeleteAliasRequest, callback func(response *kms.DeleteAliasResponse, err error)) <-chan int {
	res0 := w.obj.DeleteAliasWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DeleteAliasWithChan(ctx context.Context, request *kms.DeleteAliasRequest) (<-chan *kms.DeleteAliasResponse, <-chan error) {
	res0, res1 := w.obj.DeleteAliasWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DeleteCertificate(ctx context.Context, request *kms.DeleteCertificateRequest) (*kms.DeleteCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteCertificate")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.DeleteCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DeleteCertificate"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.DeleteCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.DeleteCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DeleteCertificateWithCallback(ctx context.Context, request *kms.DeleteCertificateRequest, callback func(response *kms.DeleteCertificateResponse, err error)) <-chan int {
	res0 := w.obj.DeleteCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DeleteCertificateWithChan(ctx context.Context, request *kms.DeleteCertificateRequest) (<-chan *kms.DeleteCertificateResponse, <-chan error) {
	res0, res1 := w.obj.DeleteCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DeleteKeyMaterial(ctx context.Context, request *kms.DeleteKeyMaterialRequest) (*kms.DeleteKeyMaterialResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteKeyMaterial")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.DeleteKeyMaterialResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DeleteKeyMaterial"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.DeleteKeyMaterial", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.DeleteKeyMaterial", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteKeyMaterial(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DeleteKeyMaterialWithCallback(ctx context.Context, request *kms.DeleteKeyMaterialRequest, callback func(response *kms.DeleteKeyMaterialResponse, err error)) <-chan int {
	res0 := w.obj.DeleteKeyMaterialWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DeleteKeyMaterialWithChan(ctx context.Context, request *kms.DeleteKeyMaterialRequest) (<-chan *kms.DeleteKeyMaterialResponse, <-chan error) {
	res0, res1 := w.obj.DeleteKeyMaterialWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DeleteSecret(ctx context.Context, request *kms.DeleteSecretRequest) (*kms.DeleteSecretResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteSecret")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.DeleteSecretResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DeleteSecret"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.DeleteSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.DeleteSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteSecret(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DeleteSecretWithCallback(ctx context.Context, request *kms.DeleteSecretRequest, callback func(response *kms.DeleteSecretResponse, err error)) <-chan int {
	res0 := w.obj.DeleteSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DeleteSecretWithChan(ctx context.Context, request *kms.DeleteSecretRequest) (<-chan *kms.DeleteSecretResponse, <-chan error) {
	res0, res1 := w.obj.DeleteSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeAccountKmsStatus(ctx context.Context, request *kms.DescribeAccountKmsStatusRequest) (*kms.DescribeAccountKmsStatusResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeAccountKmsStatus")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.DescribeAccountKmsStatusResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DescribeAccountKmsStatus"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.DescribeAccountKmsStatus", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.DescribeAccountKmsStatus", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DescribeAccountKmsStatus(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeAccountKmsStatusWithCallback(ctx context.Context, request *kms.DescribeAccountKmsStatusRequest, callback func(response *kms.DescribeAccountKmsStatusResponse, err error)) <-chan int {
	res0 := w.obj.DescribeAccountKmsStatusWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeAccountKmsStatusWithChan(ctx context.Context, request *kms.DescribeAccountKmsStatusRequest) (<-chan *kms.DescribeAccountKmsStatusResponse, <-chan error) {
	res0, res1 := w.obj.DescribeAccountKmsStatusWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeCertificate(ctx context.Context, request *kms.DescribeCertificateRequest) (*kms.DescribeCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeCertificate")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.DescribeCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DescribeCertificate"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.DescribeCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.DescribeCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DescribeCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeCertificateWithCallback(ctx context.Context, request *kms.DescribeCertificateRequest, callback func(response *kms.DescribeCertificateResponse, err error)) <-chan int {
	res0 := w.obj.DescribeCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeCertificateWithChan(ctx context.Context, request *kms.DescribeCertificateRequest) (<-chan *kms.DescribeCertificateResponse, <-chan error) {
	res0, res1 := w.obj.DescribeCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeKey(ctx context.Context, request *kms.DescribeKeyRequest) (*kms.DescribeKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeKey")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.DescribeKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DescribeKey"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.DescribeKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.DescribeKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DescribeKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeKeyVersion(ctx context.Context, request *kms.DescribeKeyVersionRequest) (*kms.DescribeKeyVersionResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeKeyVersion")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.DescribeKeyVersionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DescribeKeyVersion"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.DescribeKeyVersion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.DescribeKeyVersion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DescribeKeyVersion(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeKeyVersionWithCallback(ctx context.Context, request *kms.DescribeKeyVersionRequest, callback func(response *kms.DescribeKeyVersionResponse, err error)) <-chan int {
	res0 := w.obj.DescribeKeyVersionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeKeyVersionWithChan(ctx context.Context, request *kms.DescribeKeyVersionRequest) (<-chan *kms.DescribeKeyVersionResponse, <-chan error) {
	res0, res1 := w.obj.DescribeKeyVersionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeKeyWithCallback(ctx context.Context, request *kms.DescribeKeyRequest, callback func(response *kms.DescribeKeyResponse, err error)) <-chan int {
	res0 := w.obj.DescribeKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeKeyWithChan(ctx context.Context, request *kms.DescribeKeyRequest) (<-chan *kms.DescribeKeyResponse, <-chan error) {
	res0, res1 := w.obj.DescribeKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeRegions(ctx context.Context, request *kms.DescribeRegionsRequest) (*kms.DescribeRegionsResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeRegions")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.DescribeRegionsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DescribeRegions"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.DescribeRegions", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.DescribeRegions", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DescribeRegions(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeRegionsWithCallback(ctx context.Context, request *kms.DescribeRegionsRequest, callback func(response *kms.DescribeRegionsResponse, err error)) <-chan int {
	res0 := w.obj.DescribeRegionsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeRegionsWithChan(ctx context.Context, request *kms.DescribeRegionsRequest) (<-chan *kms.DescribeRegionsResponse, <-chan error) {
	res0, res1 := w.obj.DescribeRegionsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeSecret(ctx context.Context, request *kms.DescribeSecretRequest) (*kms.DescribeSecretResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeSecret")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.DescribeSecretResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DescribeSecret"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.DescribeSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.DescribeSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DescribeSecret(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeSecretWithCallback(ctx context.Context, request *kms.DescribeSecretRequest, callback func(response *kms.DescribeSecretResponse, err error)) <-chan int {
	res0 := w.obj.DescribeSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeSecretWithChan(ctx context.Context, request *kms.DescribeSecretRequest) (<-chan *kms.DescribeSecretResponse, <-chan error) {
	res0, res1 := w.obj.DescribeSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeService(ctx context.Context, request *kms.DescribeServiceRequest) (*kms.DescribeServiceResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeService")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.DescribeServiceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DescribeService"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.DescribeService", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.DescribeService", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DescribeService(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeServiceWithCallback(ctx context.Context, request *kms.DescribeServiceRequest, callback func(response *kms.DescribeServiceResponse, err error)) <-chan int {
	res0 := w.obj.DescribeServiceWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeServiceWithChan(ctx context.Context, request *kms.DescribeServiceRequest) (<-chan *kms.DescribeServiceResponse, <-chan error) {
	res0, res1 := w.obj.DescribeServiceWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DisableKey(ctx context.Context, request *kms.DisableKeyRequest) (*kms.DisableKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DisableKey")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.DisableKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DisableKey"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.DisableKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.DisableKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DisableKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DisableKeyWithCallback(ctx context.Context, request *kms.DisableKeyRequest, callback func(response *kms.DisableKeyResponse, err error)) <-chan int {
	res0 := w.obj.DisableKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DisableKeyWithChan(ctx context.Context, request *kms.DisableKeyRequest) (<-chan *kms.DisableKeyResponse, <-chan error) {
	res0, res1 := w.obj.DisableKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) EnableKey(ctx context.Context, request *kms.EnableKeyRequest) (*kms.EnableKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.EnableKey")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.EnableKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.EnableKey"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.EnableKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.EnableKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.EnableKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) EnableKeyWithCallback(ctx context.Context, request *kms.EnableKeyRequest, callback func(response *kms.EnableKeyResponse, err error)) <-chan int {
	res0 := w.obj.EnableKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) EnableKeyWithChan(ctx context.Context, request *kms.EnableKeyRequest) (<-chan *kms.EnableKeyResponse, <-chan error) {
	res0, res1 := w.obj.EnableKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) Encrypt(ctx context.Context, request *kms.EncryptRequest) (*kms.EncryptResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.Encrypt")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.EncryptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.Encrypt"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.Encrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.Encrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.Encrypt(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) EncryptWithCallback(ctx context.Context, request *kms.EncryptRequest, callback func(response *kms.EncryptResponse, err error)) <-chan int {
	res0 := w.obj.EncryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) EncryptWithChan(ctx context.Context, request *kms.EncryptRequest) (<-chan *kms.EncryptResponse, <-chan error) {
	res0, res1 := w.obj.EncryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ExportCertificate(ctx context.Context, request *kms.ExportCertificateRequest) (*kms.ExportCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ExportCertificate")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ExportCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ExportCertificate"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ExportCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ExportCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ExportCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ExportCertificateWithCallback(ctx context.Context, request *kms.ExportCertificateRequest, callback func(response *kms.ExportCertificateResponse, err error)) <-chan int {
	res0 := w.obj.ExportCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ExportCertificateWithChan(ctx context.Context, request *kms.ExportCertificateRequest) (<-chan *kms.ExportCertificateResponse, <-chan error) {
	res0, res1 := w.obj.ExportCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ExportDataKey(ctx context.Context, request *kms.ExportDataKeyRequest) (*kms.ExportDataKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ExportDataKey")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ExportDataKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ExportDataKey"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ExportDataKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ExportDataKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ExportDataKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ExportDataKeyWithCallback(ctx context.Context, request *kms.ExportDataKeyRequest, callback func(response *kms.ExportDataKeyResponse, err error)) <-chan int {
	res0 := w.obj.ExportDataKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ExportDataKeyWithChan(ctx context.Context, request *kms.ExportDataKeyRequest) (<-chan *kms.ExportDataKeyResponse, <-chan error) {
	res0, res1 := w.obj.ExportDataKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GenerateAndExportDataKey(ctx context.Context, request *kms.GenerateAndExportDataKeyRequest) (*kms.GenerateAndExportDataKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateAndExportDataKey")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.GenerateAndExportDataKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GenerateAndExportDataKey"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.GenerateAndExportDataKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.GenerateAndExportDataKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GenerateAndExportDataKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GenerateAndExportDataKeyWithCallback(ctx context.Context, request *kms.GenerateAndExportDataKeyRequest, callback func(response *kms.GenerateAndExportDataKeyResponse, err error)) <-chan int {
	res0 := w.obj.GenerateAndExportDataKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GenerateAndExportDataKeyWithChan(ctx context.Context, request *kms.GenerateAndExportDataKeyRequest) (<-chan *kms.GenerateAndExportDataKeyResponse, <-chan error) {
	res0, res1 := w.obj.GenerateAndExportDataKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GenerateDataKey(ctx context.Context, request *kms.GenerateDataKeyRequest) (*kms.GenerateDataKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateDataKey")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.GenerateDataKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GenerateDataKey"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.GenerateDataKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.GenerateDataKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GenerateDataKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GenerateDataKeyWithCallback(ctx context.Context, request *kms.GenerateDataKeyRequest, callback func(response *kms.GenerateDataKeyResponse, err error)) <-chan int {
	res0 := w.obj.GenerateDataKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GenerateDataKeyWithChan(ctx context.Context, request *kms.GenerateDataKeyRequest) (<-chan *kms.GenerateDataKeyResponse, <-chan error) {
	res0, res1 := w.obj.GenerateDataKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GenerateDataKeyWithoutPlaintext(ctx context.Context, request *kms.GenerateDataKeyWithoutPlaintextRequest) (*kms.GenerateDataKeyWithoutPlaintextResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateDataKeyWithoutPlaintext")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.GenerateDataKeyWithoutPlaintextResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GenerateDataKeyWithoutPlaintext"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.GenerateDataKeyWithoutPlaintext", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.GenerateDataKeyWithoutPlaintext", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GenerateDataKeyWithoutPlaintext(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GenerateDataKeyWithoutPlaintextWithCallback(ctx context.Context, request *kms.GenerateDataKeyWithoutPlaintextRequest, callback func(response *kms.GenerateDataKeyWithoutPlaintextResponse, err error)) <-chan int {
	res0 := w.obj.GenerateDataKeyWithoutPlaintextWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GenerateDataKeyWithoutPlaintextWithChan(ctx context.Context, request *kms.GenerateDataKeyWithoutPlaintextRequest) (<-chan *kms.GenerateDataKeyWithoutPlaintextResponse, <-chan error) {
	res0, res1 := w.obj.GenerateDataKeyWithoutPlaintextWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetCertificate(ctx context.Context, request *kms.GetCertificateRequest) (*kms.GetCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetCertificate")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.GetCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetCertificate"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.GetCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.GetCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetCertificateWithCallback(ctx context.Context, request *kms.GetCertificateRequest, callback func(response *kms.GetCertificateResponse, err error)) <-chan int {
	res0 := w.obj.GetCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetCertificateWithChan(ctx context.Context, request *kms.GetCertificateRequest) (<-chan *kms.GetCertificateResponse, <-chan error) {
	res0, res1 := w.obj.GetCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetParametersForImport(ctx context.Context, request *kms.GetParametersForImportRequest) (*kms.GetParametersForImportResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetParametersForImport")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.GetParametersForImportResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetParametersForImport"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.GetParametersForImport", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.GetParametersForImport", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetParametersForImport(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetParametersForImportWithCallback(ctx context.Context, request *kms.GetParametersForImportRequest, callback func(response *kms.GetParametersForImportResponse, err error)) <-chan int {
	res0 := w.obj.GetParametersForImportWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetParametersForImportWithChan(ctx context.Context, request *kms.GetParametersForImportRequest) (<-chan *kms.GetParametersForImportResponse, <-chan error) {
	res0, res1 := w.obj.GetParametersForImportWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetPublicKey(ctx context.Context, request *kms.GetPublicKeyRequest) (*kms.GetPublicKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetPublicKey")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.GetPublicKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetPublicKey"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.GetPublicKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.GetPublicKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetPublicKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetPublicKeyWithCallback(ctx context.Context, request *kms.GetPublicKeyRequest, callback func(response *kms.GetPublicKeyResponse, err error)) <-chan int {
	res0 := w.obj.GetPublicKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetPublicKeyWithChan(ctx context.Context, request *kms.GetPublicKeyRequest) (<-chan *kms.GetPublicKeyResponse, <-chan error) {
	res0, res1 := w.obj.GetPublicKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetRandomPassword(ctx context.Context, request *kms.GetRandomPasswordRequest) (*kms.GetRandomPasswordResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetRandomPassword")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.GetRandomPasswordResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetRandomPassword"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.GetRandomPassword", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.GetRandomPassword", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetRandomPassword(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetRandomPasswordWithCallback(ctx context.Context, request *kms.GetRandomPasswordRequest, callback func(response *kms.GetRandomPasswordResponse, err error)) <-chan int {
	res0 := w.obj.GetRandomPasswordWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetRandomPasswordWithChan(ctx context.Context, request *kms.GetRandomPasswordRequest) (<-chan *kms.GetRandomPasswordResponse, <-chan error) {
	res0, res1 := w.obj.GetRandomPasswordWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetSecretValue(ctx context.Context, request *kms.GetSecretValueRequest) (*kms.GetSecretValueResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetSecretValue")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.GetSecretValueResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetSecretValue"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.GetSecretValue", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.GetSecretValue", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetSecretValue(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetSecretValueWithCallback(ctx context.Context, request *kms.GetSecretValueRequest, callback func(response *kms.GetSecretValueResponse, err error)) <-chan int {
	res0 := w.obj.GetSecretValueWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetSecretValueWithChan(ctx context.Context, request *kms.GetSecretValueRequest) (<-chan *kms.GetSecretValueResponse, <-chan error) {
	res0, res1 := w.obj.GetSecretValueWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ImportCertificate(ctx context.Context, request *kms.ImportCertificateRequest) (*kms.ImportCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ImportCertificate")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ImportCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ImportCertificate"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ImportCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ImportCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ImportCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ImportCertificateWithCallback(ctx context.Context, request *kms.ImportCertificateRequest, callback func(response *kms.ImportCertificateResponse, err error)) <-chan int {
	res0 := w.obj.ImportCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ImportCertificateWithChan(ctx context.Context, request *kms.ImportCertificateRequest) (<-chan *kms.ImportCertificateResponse, <-chan error) {
	res0, res1 := w.obj.ImportCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ImportEncryptionCertificate(ctx context.Context, request *kms.ImportEncryptionCertificateRequest) (*kms.ImportEncryptionCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ImportEncryptionCertificate")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ImportEncryptionCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ImportEncryptionCertificate"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ImportEncryptionCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ImportEncryptionCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ImportEncryptionCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ImportEncryptionCertificateWithCallback(ctx context.Context, request *kms.ImportEncryptionCertificateRequest, callback func(response *kms.ImportEncryptionCertificateResponse, err error)) <-chan int {
	res0 := w.obj.ImportEncryptionCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ImportEncryptionCertificateWithChan(ctx context.Context, request *kms.ImportEncryptionCertificateRequest) (<-chan *kms.ImportEncryptionCertificateResponse, <-chan error) {
	res0, res1 := w.obj.ImportEncryptionCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ImportKeyMaterial(ctx context.Context, request *kms.ImportKeyMaterialRequest) (*kms.ImportKeyMaterialResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ImportKeyMaterial")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ImportKeyMaterialResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ImportKeyMaterial"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ImportKeyMaterial", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ImportKeyMaterial", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ImportKeyMaterial(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ImportKeyMaterialWithCallback(ctx context.Context, request *kms.ImportKeyMaterialRequest, callback func(response *kms.ImportKeyMaterialResponse, err error)) <-chan int {
	res0 := w.obj.ImportKeyMaterialWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ImportKeyMaterialWithChan(ctx context.Context, request *kms.ImportKeyMaterialRequest) (<-chan *kms.ImportKeyMaterialResponse, <-chan error) {
	res0, res1 := w.obj.ImportKeyMaterialWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListAliases(ctx context.Context, request *kms.ListAliasesRequest) (*kms.ListAliasesResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListAliases")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ListAliasesResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ListAliases"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ListAliases", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ListAliases", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListAliases(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListAliasesByKeyId(ctx context.Context, request *kms.ListAliasesByKeyIdRequest) (*kms.ListAliasesByKeyIdResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListAliasesByKeyId")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ListAliasesByKeyIdResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ListAliasesByKeyId"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ListAliasesByKeyId", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ListAliasesByKeyId", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListAliasesByKeyId(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListAliasesByKeyIdWithCallback(ctx context.Context, request *kms.ListAliasesByKeyIdRequest, callback func(response *kms.ListAliasesByKeyIdResponse, err error)) <-chan int {
	res0 := w.obj.ListAliasesByKeyIdWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListAliasesByKeyIdWithChan(ctx context.Context, request *kms.ListAliasesByKeyIdRequest) (<-chan *kms.ListAliasesByKeyIdResponse, <-chan error) {
	res0, res1 := w.obj.ListAliasesByKeyIdWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListAliasesWithCallback(ctx context.Context, request *kms.ListAliasesRequest, callback func(response *kms.ListAliasesResponse, err error)) <-chan int {
	res0 := w.obj.ListAliasesWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListAliasesWithChan(ctx context.Context, request *kms.ListAliasesRequest) (<-chan *kms.ListAliasesResponse, <-chan error) {
	res0, res1 := w.obj.ListAliasesWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListCertificates(ctx context.Context, request *kms.ListCertificatesRequest) (*kms.ListCertificatesResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListCertificates")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ListCertificatesResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ListCertificates"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ListCertificates", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ListCertificates", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListCertificates(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListCertificatesWithCallback(ctx context.Context, request *kms.ListCertificatesRequest, callback func(response *kms.ListCertificatesResponse, err error)) <-chan int {
	res0 := w.obj.ListCertificatesWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListCertificatesWithChan(ctx context.Context, request *kms.ListCertificatesRequest) (<-chan *kms.ListCertificatesResponse, <-chan error) {
	res0, res1 := w.obj.ListCertificatesWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListKeyVersions(ctx context.Context, request *kms.ListKeyVersionsRequest) (*kms.ListKeyVersionsResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListKeyVersions")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ListKeyVersionsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ListKeyVersions"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ListKeyVersions", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ListKeyVersions", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListKeyVersions(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListKeyVersionsWithCallback(ctx context.Context, request *kms.ListKeyVersionsRequest, callback func(response *kms.ListKeyVersionsResponse, err error)) <-chan int {
	res0 := w.obj.ListKeyVersionsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListKeyVersionsWithChan(ctx context.Context, request *kms.ListKeyVersionsRequest) (<-chan *kms.ListKeyVersionsResponse, <-chan error) {
	res0, res1 := w.obj.ListKeyVersionsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListKeys(ctx context.Context, request *kms.ListKeysRequest) (*kms.ListKeysResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListKeys")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ListKeysResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ListKeys"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ListKeys", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ListKeys", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListKeys(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListKeysWithCallback(ctx context.Context, request *kms.ListKeysRequest, callback func(response *kms.ListKeysResponse, err error)) <-chan int {
	res0 := w.obj.ListKeysWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListKeysWithChan(ctx context.Context, request *kms.ListKeysRequest) (<-chan *kms.ListKeysResponse, <-chan error) {
	res0, res1 := w.obj.ListKeysWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListResourceTags(ctx context.Context, request *kms.ListResourceTagsRequest) (*kms.ListResourceTagsResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListResourceTags")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ListResourceTagsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ListResourceTags"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ListResourceTags", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ListResourceTags", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListResourceTags(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListResourceTagsWithCallback(ctx context.Context, request *kms.ListResourceTagsRequest, callback func(response *kms.ListResourceTagsResponse, err error)) <-chan int {
	res0 := w.obj.ListResourceTagsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListResourceTagsWithChan(ctx context.Context, request *kms.ListResourceTagsRequest) (<-chan *kms.ListResourceTagsResponse, <-chan error) {
	res0, res1 := w.obj.ListResourceTagsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListSecretVersionIds(ctx context.Context, request *kms.ListSecretVersionIdsRequest) (*kms.ListSecretVersionIdsResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListSecretVersionIds")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ListSecretVersionIdsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ListSecretVersionIds"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ListSecretVersionIds", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ListSecretVersionIds", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListSecretVersionIds(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListSecretVersionIdsWithCallback(ctx context.Context, request *kms.ListSecretVersionIdsRequest, callback func(response *kms.ListSecretVersionIdsResponse, err error)) <-chan int {
	res0 := w.obj.ListSecretVersionIdsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListSecretVersionIdsWithChan(ctx context.Context, request *kms.ListSecretVersionIdsRequest) (<-chan *kms.ListSecretVersionIdsResponse, <-chan error) {
	res0, res1 := w.obj.ListSecretVersionIdsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListSecrets(ctx context.Context, request *kms.ListSecretsRequest) (*kms.ListSecretsResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListSecrets")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ListSecretsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ListSecrets"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ListSecrets", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ListSecrets", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListSecrets(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListSecretsWithCallback(ctx context.Context, request *kms.ListSecretsRequest, callback func(response *kms.ListSecretsResponse, err error)) <-chan int {
	res0 := w.obj.ListSecretsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListSecretsWithChan(ctx context.Context, request *kms.ListSecretsRequest) (<-chan *kms.ListSecretsResponse, <-chan error) {
	res0, res1 := w.obj.ListSecretsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) OpenKmsService(ctx context.Context, request *kms.OpenKmsServiceRequest) (*kms.OpenKmsServiceResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.OpenKmsService")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.OpenKmsServiceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.OpenKmsService"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.OpenKmsService", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.OpenKmsService", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.OpenKmsService(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) OpenKmsServiceWithCallback(ctx context.Context, request *kms.OpenKmsServiceRequest, callback func(response *kms.OpenKmsServiceResponse, err error)) <-chan int {
	res0 := w.obj.OpenKmsServiceWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) OpenKmsServiceWithChan(ctx context.Context, request *kms.OpenKmsServiceRequest) (<-chan *kms.OpenKmsServiceResponse, <-chan error) {
	res0, res1 := w.obj.OpenKmsServiceWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) PutSecretValue(ctx context.Context, request *kms.PutSecretValueRequest) (*kms.PutSecretValueResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.PutSecretValue")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.PutSecretValueResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.PutSecretValue"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.PutSecretValue", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.PutSecretValue", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.PutSecretValue(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) PutSecretValueWithCallback(ctx context.Context, request *kms.PutSecretValueRequest, callback func(response *kms.PutSecretValueResponse, err error)) <-chan int {
	res0 := w.obj.PutSecretValueWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) PutSecretValueWithChan(ctx context.Context, request *kms.PutSecretValueRequest) (<-chan *kms.PutSecretValueResponse, <-chan error) {
	res0, res1 := w.obj.PutSecretValueWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ReEncrypt(ctx context.Context, request *kms.ReEncryptRequest) (*kms.ReEncryptResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ReEncrypt")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ReEncryptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ReEncrypt"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ReEncrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ReEncrypt", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ReEncrypt(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ReEncryptWithCallback(ctx context.Context, request *kms.ReEncryptRequest, callback func(response *kms.ReEncryptResponse, err error)) <-chan int {
	res0 := w.obj.ReEncryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ReEncryptWithChan(ctx context.Context, request *kms.ReEncryptRequest) (<-chan *kms.ReEncryptResponse, <-chan error) {
	res0, res1 := w.obj.ReEncryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) RestoreSecret(ctx context.Context, request *kms.RestoreSecretRequest) (*kms.RestoreSecretResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.RestoreSecret")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.RestoreSecretResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.RestoreSecret"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.RestoreSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.RestoreSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.RestoreSecret(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) RestoreSecretWithCallback(ctx context.Context, request *kms.RestoreSecretRequest, callback func(response *kms.RestoreSecretResponse, err error)) <-chan int {
	res0 := w.obj.RestoreSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) RestoreSecretWithChan(ctx context.Context, request *kms.RestoreSecretRequest) (<-chan *kms.RestoreSecretResponse, <-chan error) {
	res0, res1 := w.obj.RestoreSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) RotateSecret(ctx context.Context, request *kms.RotateSecretRequest) (*kms.RotateSecretResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.RotateSecret")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.RotateSecretResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.RotateSecret"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.RotateSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.RotateSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.RotateSecret(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) RotateSecretWithCallback(ctx context.Context, request *kms.RotateSecretRequest, callback func(response *kms.RotateSecretResponse, err error)) <-chan int {
	res0 := w.obj.RotateSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) RotateSecretWithChan(ctx context.Context, request *kms.RotateSecretRequest) (<-chan *kms.RotateSecretResponse, <-chan error) {
	res0, res1 := w.obj.RotateSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ScheduleKeyDeletion(ctx context.Context, request *kms.ScheduleKeyDeletionRequest) (*kms.ScheduleKeyDeletionResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ScheduleKeyDeletion")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.ScheduleKeyDeletionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ScheduleKeyDeletion"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.ScheduleKeyDeletion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.ScheduleKeyDeletion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ScheduleKeyDeletion(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ScheduleKeyDeletionWithCallback(ctx context.Context, request *kms.ScheduleKeyDeletionRequest, callback func(response *kms.ScheduleKeyDeletionResponse, err error)) <-chan int {
	res0 := w.obj.ScheduleKeyDeletionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ScheduleKeyDeletionWithChan(ctx context.Context, request *kms.ScheduleKeyDeletionRequest) (<-chan *kms.ScheduleKeyDeletionResponse, <-chan error) {
	res0, res1 := w.obj.ScheduleKeyDeletionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) TagResource(ctx context.Context, request *kms.TagResourceRequest) (*kms.TagResourceResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.TagResource")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.TagResourceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.TagResource"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.TagResource", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.TagResource", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.TagResource(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) TagResourceWithCallback(ctx context.Context, request *kms.TagResourceRequest, callback func(response *kms.TagResourceResponse, err error)) <-chan int {
	res0 := w.obj.TagResourceWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) TagResourceWithChan(ctx context.Context, request *kms.TagResourceRequest) (<-chan *kms.TagResourceResponse, <-chan error) {
	res0, res1 := w.obj.TagResourceWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UntagResource(ctx context.Context, request *kms.UntagResourceRequest) (*kms.UntagResourceResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UntagResource")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.UntagResourceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.UntagResource"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.UntagResource", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.UntagResource", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UntagResource(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UntagResourceWithCallback(ctx context.Context, request *kms.UntagResourceRequest, callback func(response *kms.UntagResourceResponse, err error)) <-chan int {
	res0 := w.obj.UntagResourceWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UntagResourceWithChan(ctx context.Context, request *kms.UntagResourceRequest) (<-chan *kms.UntagResourceResponse, <-chan error) {
	res0, res1 := w.obj.UntagResourceWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateAlias(ctx context.Context, request *kms.UpdateAliasRequest) (*kms.UpdateAliasResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateAlias")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.UpdateAliasResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.UpdateAlias"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.UpdateAlias", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.UpdateAlias", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateAlias(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateAliasWithCallback(ctx context.Context, request *kms.UpdateAliasRequest, callback func(response *kms.UpdateAliasResponse, err error)) <-chan int {
	res0 := w.obj.UpdateAliasWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateAliasWithChan(ctx context.Context, request *kms.UpdateAliasRequest) (<-chan *kms.UpdateAliasResponse, <-chan error) {
	res0, res1 := w.obj.UpdateAliasWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateCertificateStatus(ctx context.Context, request *kms.UpdateCertificateStatusRequest) (*kms.UpdateCertificateStatusResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateCertificateStatus")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.UpdateCertificateStatusResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.UpdateCertificateStatus"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.UpdateCertificateStatus", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.UpdateCertificateStatus", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateCertificateStatus(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateCertificateStatusWithCallback(ctx context.Context, request *kms.UpdateCertificateStatusRequest, callback func(response *kms.UpdateCertificateStatusResponse, err error)) <-chan int {
	res0 := w.obj.UpdateCertificateStatusWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateCertificateStatusWithChan(ctx context.Context, request *kms.UpdateCertificateStatusRequest) (<-chan *kms.UpdateCertificateStatusResponse, <-chan error) {
	res0, res1 := w.obj.UpdateCertificateStatusWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateKeyDescription(ctx context.Context, request *kms.UpdateKeyDescriptionRequest) (*kms.UpdateKeyDescriptionResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateKeyDescription")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.UpdateKeyDescriptionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.UpdateKeyDescription"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.UpdateKeyDescription", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.UpdateKeyDescription", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateKeyDescription(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateKeyDescriptionWithCallback(ctx context.Context, request *kms.UpdateKeyDescriptionRequest, callback func(response *kms.UpdateKeyDescriptionResponse, err error)) <-chan int {
	res0 := w.obj.UpdateKeyDescriptionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateKeyDescriptionWithChan(ctx context.Context, request *kms.UpdateKeyDescriptionRequest) (<-chan *kms.UpdateKeyDescriptionResponse, <-chan error) {
	res0, res1 := w.obj.UpdateKeyDescriptionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateRotationPolicy(ctx context.Context, request *kms.UpdateRotationPolicyRequest) (*kms.UpdateRotationPolicyResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateRotationPolicy")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.UpdateRotationPolicyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.UpdateRotationPolicy"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.UpdateRotationPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.UpdateRotationPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateRotationPolicy(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateRotationPolicyWithCallback(ctx context.Context, request *kms.UpdateRotationPolicyRequest, callback func(response *kms.UpdateRotationPolicyResponse, err error)) <-chan int {
	res0 := w.obj.UpdateRotationPolicyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateRotationPolicyWithChan(ctx context.Context, request *kms.UpdateRotationPolicyRequest) (<-chan *kms.UpdateRotationPolicyResponse, <-chan error) {
	res0, res1 := w.obj.UpdateRotationPolicyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateSecret(ctx context.Context, request *kms.UpdateSecretRequest) (*kms.UpdateSecretResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateSecret")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.UpdateSecretResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.UpdateSecret"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.UpdateSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.UpdateSecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateSecret(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateSecretRotationPolicy(ctx context.Context, request *kms.UpdateSecretRotationPolicyRequest) (*kms.UpdateSecretRotationPolicyResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateSecretRotationPolicy")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.UpdateSecretRotationPolicyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.UpdateSecretRotationPolicy"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.UpdateSecretRotationPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.UpdateSecretRotationPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateSecretRotationPolicy(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateSecretRotationPolicyWithCallback(ctx context.Context, request *kms.UpdateSecretRotationPolicyRequest, callback func(response *kms.UpdateSecretRotationPolicyResponse, err error)) <-chan int {
	res0 := w.obj.UpdateSecretRotationPolicyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateSecretRotationPolicyWithChan(ctx context.Context, request *kms.UpdateSecretRotationPolicyRequest) (<-chan *kms.UpdateSecretRotationPolicyResponse, <-chan error) {
	res0, res1 := w.obj.UpdateSecretRotationPolicyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateSecretVersionStage(ctx context.Context, request *kms.UpdateSecretVersionStageRequest) (*kms.UpdateSecretVersionStageResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateSecretVersionStage")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.UpdateSecretVersionStageResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.UpdateSecretVersionStage"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.UpdateSecretVersionStage", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.UpdateSecretVersionStage", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateSecretVersionStage(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateSecretVersionStageWithCallback(ctx context.Context, request *kms.UpdateSecretVersionStageRequest, callback func(response *kms.UpdateSecretVersionStageResponse, err error)) <-chan int {
	res0 := w.obj.UpdateSecretVersionStageWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateSecretVersionStageWithChan(ctx context.Context, request *kms.UpdateSecretVersionStageRequest) (<-chan *kms.UpdateSecretVersionStageResponse, <-chan error) {
	res0, res1 := w.obj.UpdateSecretVersionStageWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateSecretWithCallback(ctx context.Context, request *kms.UpdateSecretRequest, callback func(response *kms.UpdateSecretResponse, err error)) <-chan int {
	res0 := w.obj.UpdateSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateSecretWithChan(ctx context.Context, request *kms.UpdateSecretRequest) (<-chan *kms.UpdateSecretResponse, <-chan error) {
	res0, res1 := w.obj.UpdateSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UploadCertificate(ctx context.Context, request *kms.UploadCertificateRequest) (*kms.UploadCertificateResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UploadCertificate")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var response *kms.UploadCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.UploadCertificate"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("kms.Client.UploadCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("kms.Client.UploadCertificate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UploadCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UploadCertificateWithCallback(ctx context.Context, request *kms.UploadCertificateRequest, callback func(response *kms.UploadCertificateResponse, err error)) <-chan int {
	res0 := w.obj.UploadCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UploadCertificateWithChan(ctx context.Context, request *kms.UploadCertificateRequest) (<-chan *kms.UploadCertificateResponse, <-chan error) {
	res0, res1 := w.obj.UploadCertificateWithChan(request)
	return res0, res1
}
