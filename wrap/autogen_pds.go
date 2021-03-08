// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"fmt"
	"time"

	"github.com/alibabacloud-go/pds-sdk/client"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

func NewPDSClientWrapper(
	obj *client.Client,
	retry *micro.Retry,
	options *WrapperOptions,
	durationMetric *prometheus.HistogramVec,
	inflightMetric *prometheus.GaugeVec,
	rateLimiter micro.RateLimiter,
	parallelController micro.ParallelController) *PDSClientWrapper {
	return &PDSClientWrapper{
		obj:                obj,
		retry:              retry,
		options:            options,
		durationMetric:     durationMetric,
		inflightMetric:     inflightMetric,
		rateLimiter:        rateLimiter,
		parallelController: parallelController,
	}
}

type PDSClientWrapper struct {
	obj                *client.Client
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *PDSClientWrapper) Unwrap() *client.Client {
	return w.obj
}

func (w *PDSClientWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *PDSClientWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *PDSClientWrapper) OnRateLimiterChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *PDSClientWrapper) OnParallelControllerChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *PDSClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        fmt.Sprintf("%s_client_Client_durationMs", options.Name),
		Help:        "client Client response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        fmt.Sprintf("%s_client_Client_inflight", options.Name),
		Help:        "client Client inflight",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "custom"})
}

func (w *PDSClientWrapper) AccountRevoke(ctx context.Context, request *client.RevokeRequest) (*client.AccountRevokeModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.AccountRevokeModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AccountRevoke", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AccountRevoke", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AccountRevoke", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.AccountRevoke")
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
			w.inflightMetric.WithLabelValues("client.Client.AccountRevoke", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.AccountRevoke", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.AccountRevoke", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.AccountRevoke(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) AccountRevokeEx(ctx context.Context, request *client.RevokeRequest, runtime *client.RuntimeOptions) (*client.AccountRevokeModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.AccountRevokeModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AccountRevokeEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AccountRevokeEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AccountRevokeEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.AccountRevokeEx")
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
			w.inflightMetric.WithLabelValues("client.Client.AccountRevokeEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.AccountRevokeEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.AccountRevokeEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.AccountRevokeEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) AccountToken(ctx context.Context, request *client.AccountTokenRequest) (*client.AccountTokenModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.AccountTokenModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AccountToken", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AccountToken", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AccountToken", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.AccountToken")
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
			w.inflightMetric.WithLabelValues("client.Client.AccountToken", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.AccountToken", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.AccountToken", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.AccountToken(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) AccountTokenEx(ctx context.Context, request *client.AccountTokenRequest, runtime *client.RuntimeOptions) (*client.AccountTokenModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.AccountTokenModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AccountTokenEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AccountTokenEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AccountTokenEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.AccountTokenEx")
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
			w.inflightMetric.WithLabelValues("client.Client.AccountTokenEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.AccountTokenEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.AccountTokenEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.AccountTokenEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) AdminListStores(ctx context.Context, request *client.AdminListStoresRequest) (*client.AdminListStoresModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.AdminListStoresModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AdminListStores", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AdminListStores", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AdminListStores", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.AdminListStores")
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
			w.inflightMetric.WithLabelValues("client.Client.AdminListStores", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.AdminListStores", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.AdminListStores", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.AdminListStores(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) AdminListStoresEx(ctx context.Context, request *client.AdminListStoresRequest, runtime *client.RuntimeOptions) (*client.AdminListStoresModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.AdminListStoresModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AdminListStoresEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AdminListStoresEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AdminListStoresEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.AdminListStoresEx")
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
			w.inflightMetric.WithLabelValues("client.Client.AdminListStoresEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.AdminListStoresEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.AdminListStoresEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.AdminListStoresEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) AppendUserAgent(userAgent *string) {
	w.obj.AppendUserAgent(userAgent)
}

func (w *PDSClientWrapper) BatchOperation(ctx context.Context, request *client.BatchRequest) (*client.BatchOperationModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.BatchOperationModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.BatchOperation", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.BatchOperation", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.BatchOperation", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.BatchOperation")
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
			w.inflightMetric.WithLabelValues("client.Client.BatchOperation", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.BatchOperation", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.BatchOperation", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.BatchOperation(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) BatchOperationEx(ctx context.Context, request *client.BatchRequest, runtime *client.RuntimeOptions) (*client.BatchOperationModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.BatchOperationModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.BatchOperationEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.BatchOperationEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.BatchOperationEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.BatchOperationEx")
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
			w.inflightMetric.WithLabelValues("client.Client.BatchOperationEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.BatchOperationEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.BatchOperationEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.BatchOperationEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CancelLink(ctx context.Context, request *client.CancelLinkRequest) (*client.CancelLinkModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CancelLinkModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CancelLink", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CancelLink", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CancelLink", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CancelLink")
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
			w.inflightMetric.WithLabelValues("client.Client.CancelLink", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CancelLink", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CancelLink", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CancelLink(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CancelLinkEx(ctx context.Context, request *client.CancelLinkRequest, runtime *client.RuntimeOptions) (*client.CancelLinkModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CancelLinkModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CancelLinkEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CancelLinkEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CancelLinkEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CancelLinkEx")
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
			w.inflightMetric.WithLabelValues("client.Client.CancelLinkEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CancelLinkEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CancelLinkEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CancelLinkEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CancelShareLink(ctx context.Context, request *client.CancelShareLinkRequest) (*client.CancelShareLinkModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CancelShareLinkModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CancelShareLink", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CancelShareLink", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CancelShareLink", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CancelShareLink")
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
			w.inflightMetric.WithLabelValues("client.Client.CancelShareLink", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CancelShareLink", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CancelShareLink", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CancelShareLink(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CancelShareLinkEx(ctx context.Context, request *client.CancelShareLinkRequest, runtime *client.RuntimeOptions) (*client.CancelShareLinkModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CancelShareLinkModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CancelShareLinkEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CancelShareLinkEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CancelShareLinkEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CancelShareLinkEx")
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
			w.inflightMetric.WithLabelValues("client.Client.CancelShareLinkEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CancelShareLinkEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CancelShareLinkEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CancelShareLinkEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ChangePassword(ctx context.Context, request *client.DefaultChangePasswordRequest) (*client.ChangePasswordModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ChangePasswordModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ChangePassword", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ChangePassword", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ChangePassword", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ChangePassword")
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
			w.inflightMetric.WithLabelValues("client.Client.ChangePassword", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ChangePassword", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ChangePassword", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ChangePassword(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ChangePasswordEx(ctx context.Context, request *client.DefaultChangePasswordRequest, runtime *client.RuntimeOptions) (*client.ChangePasswordModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ChangePasswordModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ChangePasswordEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ChangePasswordEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ChangePasswordEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ChangePasswordEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ChangePasswordEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ChangePasswordEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ChangePasswordEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ChangePasswordEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CheckExist(ctx context.Context, request *client.MobileCheckExistRequest) (*client.CheckExistModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CheckExistModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CheckExist", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CheckExist", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CheckExist", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CheckExist")
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
			w.inflightMetric.WithLabelValues("client.Client.CheckExist", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CheckExist", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CheckExist", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CheckExist(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CheckExistEx(ctx context.Context, request *client.MobileCheckExistRequest, runtime *client.RuntimeOptions) (*client.CheckExistModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CheckExistModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CheckExistEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CheckExistEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CheckExistEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CheckExistEx")
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
			w.inflightMetric.WithLabelValues("client.Client.CheckExistEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CheckExistEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CheckExistEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CheckExistEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CompleteFile(ctx context.Context, request *client.CompleteFileRequest) (*client.CompleteFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CompleteFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CompleteFile", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CompleteFile", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CompleteFile", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CompleteFile")
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
			w.inflightMetric.WithLabelValues("client.Client.CompleteFile", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CompleteFile", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CompleteFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CompleteFile(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CompleteFileEx(ctx context.Context, request *client.CompleteFileRequest, runtime *client.RuntimeOptions) (*client.CompleteFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CompleteFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CompleteFileEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CompleteFileEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CompleteFileEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CompleteFileEx")
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
			w.inflightMetric.WithLabelValues("client.Client.CompleteFileEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CompleteFileEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CompleteFileEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CompleteFileEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ConfirmLink(ctx context.Context, request *client.ConfirmLinkRequest) (*client.ConfirmLinkModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ConfirmLinkModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ConfirmLink", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ConfirmLink", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ConfirmLink", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ConfirmLink")
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
			w.inflightMetric.WithLabelValues("client.Client.ConfirmLink", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ConfirmLink", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ConfirmLink", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ConfirmLink(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ConfirmLinkEx(ctx context.Context, request *client.ConfirmLinkRequest, runtime *client.RuntimeOptions) (*client.ConfirmLinkModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ConfirmLinkModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ConfirmLinkEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ConfirmLinkEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ConfirmLinkEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ConfirmLinkEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ConfirmLinkEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ConfirmLinkEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ConfirmLinkEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ConfirmLinkEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CopyFile(ctx context.Context, request *client.CopyFileRequest) (*client.CopyFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CopyFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CopyFile", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CopyFile", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CopyFile", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CopyFile")
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
			w.inflightMetric.WithLabelValues("client.Client.CopyFile", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CopyFile", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CopyFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CopyFile(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CopyFileEx(ctx context.Context, request *client.CopyFileRequest, runtime *client.RuntimeOptions) (*client.CopyFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CopyFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CopyFileEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CopyFileEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CopyFileEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CopyFileEx")
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
			w.inflightMetric.WithLabelValues("client.Client.CopyFileEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CopyFileEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CopyFileEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CopyFileEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CreateDrive(ctx context.Context, request *client.CreateDriveRequest) (*client.CreateDriveModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CreateDriveModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateDrive", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateDrive", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateDrive", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CreateDrive")
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
			w.inflightMetric.WithLabelValues("client.Client.CreateDrive", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CreateDrive", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CreateDrive", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CreateDrive(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CreateDriveEx(ctx context.Context, request *client.CreateDriveRequest, runtime *client.RuntimeOptions) (*client.CreateDriveModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CreateDriveModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateDriveEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateDriveEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateDriveEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CreateDriveEx")
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
			w.inflightMetric.WithLabelValues("client.Client.CreateDriveEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CreateDriveEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CreateDriveEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CreateDriveEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CreateFile(ctx context.Context, request *client.CreateFileRequest) (*client.CreateFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CreateFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateFile", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateFile", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateFile", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CreateFile")
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
			w.inflightMetric.WithLabelValues("client.Client.CreateFile", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CreateFile", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CreateFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CreateFile(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CreateFileEx(ctx context.Context, request *client.CreateFileRequest, runtime *client.RuntimeOptions) (*client.CreateFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CreateFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateFileEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateFileEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateFileEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CreateFileEx")
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
			w.inflightMetric.WithLabelValues("client.Client.CreateFileEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CreateFileEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CreateFileEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CreateFileEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CreateShare(ctx context.Context, request *client.CreateShareRequest) (*client.CreateShareModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CreateShareModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateShare", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateShare", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateShare", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CreateShare")
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
			w.inflightMetric.WithLabelValues("client.Client.CreateShare", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CreateShare", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CreateShare", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CreateShare(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CreateShareEx(ctx context.Context, request *client.CreateShareRequest, runtime *client.RuntimeOptions) (*client.CreateShareModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CreateShareModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateShareEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateShareEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateShareEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CreateShareEx")
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
			w.inflightMetric.WithLabelValues("client.Client.CreateShareEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CreateShareEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CreateShareEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CreateShareEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CreateShareLink(ctx context.Context, request *client.CreateShareLinkRequest) (*client.CreateShareLinkModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CreateShareLinkModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateShareLink", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateShareLink", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateShareLink", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CreateShareLink")
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
			w.inflightMetric.WithLabelValues("client.Client.CreateShareLink", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CreateShareLink", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CreateShareLink", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CreateShareLink(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CreateShareLinkEx(ctx context.Context, request *client.CreateShareLinkRequest, runtime *client.RuntimeOptions) (*client.CreateShareLinkModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CreateShareLinkModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateShareLinkEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateShareLinkEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateShareLinkEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CreateShareLinkEx")
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
			w.inflightMetric.WithLabelValues("client.Client.CreateShareLinkEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CreateShareLinkEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CreateShareLinkEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CreateShareLinkEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CreateStory(ctx context.Context, request *client.CreateStoryRequest) (*client.CreateStoryModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CreateStoryModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateStory", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateStory", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateStory", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CreateStory")
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
			w.inflightMetric.WithLabelValues("client.Client.CreateStory", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CreateStory", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CreateStory", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CreateStory(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CreateStoryEx(ctx context.Context, request *client.CreateStoryRequest, runtime *client.RuntimeOptions) (*client.CreateStoryModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CreateStoryModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateStoryEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateStoryEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateStoryEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CreateStoryEx")
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
			w.inflightMetric.WithLabelValues("client.Client.CreateStoryEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CreateStoryEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CreateStoryEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CreateStoryEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CreateUser(ctx context.Context, request *client.CreateUserRequest) (*client.CreateUserModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CreateUserModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateUser", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateUser", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateUser", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CreateUser")
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
			w.inflightMetric.WithLabelValues("client.Client.CreateUser", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CreateUser", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CreateUser", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CreateUser(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) CreateUserEx(ctx context.Context, request *client.CreateUserRequest, runtime *client.RuntimeOptions) (*client.CreateUserModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.CreateUserModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateUserEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateUserEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateUserEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.CreateUserEx")
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
			w.inflightMetric.WithLabelValues("client.Client.CreateUserEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.CreateUserEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.CreateUserEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.CreateUserEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) DeleteDrive(ctx context.Context, request *client.DeleteDriveRequest) (*client.DeleteDriveModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.DeleteDriveModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteDrive", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteDrive", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteDrive", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.DeleteDrive")
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
			w.inflightMetric.WithLabelValues("client.Client.DeleteDrive", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.DeleteDrive", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.DeleteDrive", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.DeleteDrive(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) DeleteDriveEx(ctx context.Context, request *client.DeleteDriveRequest, runtime *client.RuntimeOptions) (*client.DeleteDriveModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.DeleteDriveModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteDriveEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteDriveEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteDriveEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.DeleteDriveEx")
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
			w.inflightMetric.WithLabelValues("client.Client.DeleteDriveEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.DeleteDriveEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.DeleteDriveEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.DeleteDriveEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) DeleteFile(ctx context.Context, request *client.DeleteFileRequest) (*client.DeleteFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.DeleteFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteFile", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteFile", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteFile", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.DeleteFile")
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
			w.inflightMetric.WithLabelValues("client.Client.DeleteFile", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.DeleteFile", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.DeleteFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.DeleteFile(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) DeleteFileEx(ctx context.Context, request *client.DeleteFileRequest, runtime *client.RuntimeOptions) (*client.DeleteFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.DeleteFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteFileEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteFileEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteFileEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.DeleteFileEx")
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
			w.inflightMetric.WithLabelValues("client.Client.DeleteFileEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.DeleteFileEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.DeleteFileEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.DeleteFileEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) DeleteShare(ctx context.Context, request *client.DeleteShareRequest) (*client.DeleteShareModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.DeleteShareModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteShare", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteShare", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteShare", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.DeleteShare")
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
			w.inflightMetric.WithLabelValues("client.Client.DeleteShare", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.DeleteShare", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.DeleteShare", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.DeleteShare(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) DeleteShareEx(ctx context.Context, request *client.DeleteShareRequest, runtime *client.RuntimeOptions) (*client.DeleteShareModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.DeleteShareModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteShareEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteShareEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteShareEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.DeleteShareEx")
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
			w.inflightMetric.WithLabelValues("client.Client.DeleteShareEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.DeleteShareEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.DeleteShareEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.DeleteShareEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) DeleteUser(ctx context.Context, request *client.DeleteUserRequest) (*client.DeleteUserModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.DeleteUserModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteUser", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteUser", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteUser", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.DeleteUser")
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
			w.inflightMetric.WithLabelValues("client.Client.DeleteUser", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.DeleteUser", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.DeleteUser", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.DeleteUser(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) DeleteUserEx(ctx context.Context, request *client.DeleteUserRequest, runtime *client.RuntimeOptions) (*client.DeleteUserModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.DeleteUserModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteUserEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteUserEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteUserEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.DeleteUserEx")
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
			w.inflightMetric.WithLabelValues("client.Client.DeleteUserEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.DeleteUserEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.DeleteUserEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.DeleteUserEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetAccessKeyId(ctx context.Context) (*string, error) {
	ctxOptions := FromContext(ctx)
	var _result *string
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetAccessKeyId", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetAccessKeyId", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetAccessKeyId", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetAccessKeyId")
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
			w.inflightMetric.WithLabelValues("client.Client.GetAccessKeyId", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetAccessKeyId", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetAccessKeyId", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetAccessKeyId()
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetAccessKeySecret(ctx context.Context) (*string, error) {
	ctxOptions := FromContext(ctx)
	var _result *string
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetAccessKeySecret", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetAccessKeySecret", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetAccessKeySecret", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetAccessKeySecret")
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
			w.inflightMetric.WithLabelValues("client.Client.GetAccessKeySecret", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetAccessKeySecret", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetAccessKeySecret", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetAccessKeySecret()
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetAccessToken(ctx context.Context) (*string, error) {
	ctxOptions := FromContext(ctx)
	var _result *string
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetAccessToken", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetAccessToken", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetAccessToken", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetAccessToken")
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
			w.inflightMetric.WithLabelValues("client.Client.GetAccessToken", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetAccessToken", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetAccessToken", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetAccessToken()
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetAccessTokenByLinkInfo(ctx context.Context, request *client.GetAccessTokenByLinkInfoRequest) (*client.GetAccessTokenByLinkInfoModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetAccessTokenByLinkInfoModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetAccessTokenByLinkInfo", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetAccessTokenByLinkInfo", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetAccessTokenByLinkInfo", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetAccessTokenByLinkInfo")
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
			w.inflightMetric.WithLabelValues("client.Client.GetAccessTokenByLinkInfo", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetAccessTokenByLinkInfo", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetAccessTokenByLinkInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetAccessTokenByLinkInfo(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetAccessTokenByLinkInfoEx(ctx context.Context, request *client.GetAccessTokenByLinkInfoRequest, runtime *client.RuntimeOptions) (*client.GetAccessTokenByLinkInfoModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetAccessTokenByLinkInfoModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetAccessTokenByLinkInfoEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetAccessTokenByLinkInfoEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetAccessTokenByLinkInfoEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetAccessTokenByLinkInfoEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetAccessTokenByLinkInfoEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetAccessTokenByLinkInfoEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetAccessTokenByLinkInfoEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetAccessTokenByLinkInfoEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetAsyncTaskInfo(ctx context.Context, request *client.GetAsyncTaskRequest) (*client.GetAsyncTaskInfoModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetAsyncTaskInfoModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetAsyncTaskInfo", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetAsyncTaskInfo", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetAsyncTaskInfo", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetAsyncTaskInfo")
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
			w.inflightMetric.WithLabelValues("client.Client.GetAsyncTaskInfo", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetAsyncTaskInfo", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetAsyncTaskInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetAsyncTaskInfo(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetAsyncTaskInfoEx(ctx context.Context, request *client.GetAsyncTaskRequest, runtime *client.RuntimeOptions) (*client.GetAsyncTaskInfoModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetAsyncTaskInfoModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetAsyncTaskInfoEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetAsyncTaskInfoEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetAsyncTaskInfoEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetAsyncTaskInfoEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetAsyncTaskInfoEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetAsyncTaskInfoEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetAsyncTaskInfoEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetAsyncTaskInfoEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetCaptcha(ctx context.Context, request *client.GetCaptchaRequest) (*client.GetCaptchaModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetCaptchaModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetCaptcha", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetCaptcha", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetCaptcha", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetCaptcha")
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
			w.inflightMetric.WithLabelValues("client.Client.GetCaptcha", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetCaptcha", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetCaptcha", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetCaptcha(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetCaptchaEx(ctx context.Context, request *client.GetCaptchaRequest, runtime *client.RuntimeOptions) (*client.GetCaptchaModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetCaptchaModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetCaptchaEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetCaptchaEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetCaptchaEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetCaptchaEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetCaptchaEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetCaptchaEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetCaptchaEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetCaptchaEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetDefaultDrive(ctx context.Context, request *client.GetDefaultDriveRequest) (*client.GetDefaultDriveModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetDefaultDriveModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetDefaultDrive", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetDefaultDrive", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetDefaultDrive", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetDefaultDrive")
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
			w.inflightMetric.WithLabelValues("client.Client.GetDefaultDrive", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetDefaultDrive", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetDefaultDrive", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetDefaultDrive(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetDefaultDriveEx(ctx context.Context, request *client.GetDefaultDriveRequest, runtime *client.RuntimeOptions) (*client.GetDefaultDriveModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetDefaultDriveModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetDefaultDriveEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetDefaultDriveEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetDefaultDriveEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetDefaultDriveEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetDefaultDriveEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetDefaultDriveEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetDefaultDriveEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetDefaultDriveEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetDownloadUrl(ctx context.Context, request *client.GetDownloadUrlRequest) (*client.GetDownloadUrlModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetDownloadUrlModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetDownloadUrl", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetDownloadUrl", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetDownloadUrl", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetDownloadUrl")
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
			w.inflightMetric.WithLabelValues("client.Client.GetDownloadUrl", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetDownloadUrl", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetDownloadUrl", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetDownloadUrl(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetDownloadUrlEx(ctx context.Context, request *client.GetDownloadUrlRequest, runtime *client.RuntimeOptions) (*client.GetDownloadUrlModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetDownloadUrlModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetDownloadUrlEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetDownloadUrlEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetDownloadUrlEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetDownloadUrlEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetDownloadUrlEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetDownloadUrlEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetDownloadUrlEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetDownloadUrlEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetDrive(ctx context.Context, request *client.GetDriveRequest) (*client.GetDriveModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetDriveModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetDrive", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetDrive", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetDrive", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetDrive")
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
			w.inflightMetric.WithLabelValues("client.Client.GetDrive", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetDrive", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetDrive", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetDrive(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetDriveEx(ctx context.Context, request *client.GetDriveRequest, runtime *client.RuntimeOptions) (*client.GetDriveModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetDriveModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetDriveEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetDriveEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetDriveEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetDriveEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetDriveEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetDriveEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetDriveEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetDriveEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetExpireTime() *string {
	_result := w.obj.GetExpireTime()
	return _result
}

func (w *PDSClientWrapper) GetFile(ctx context.Context, request *client.GetFileRequest) (*client.GetFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetFile", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetFile", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetFile", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetFile")
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
			w.inflightMetric.WithLabelValues("client.Client.GetFile", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetFile", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetFile(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetFileByPath(ctx context.Context, request *client.GetFileByPathRequest) (*client.GetFileByPathModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetFileByPathModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetFileByPath", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetFileByPath", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetFileByPath", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetFileByPath")
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
			w.inflightMetric.WithLabelValues("client.Client.GetFileByPath", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetFileByPath", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetFileByPath", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetFileByPath(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetFileByPathEx(ctx context.Context, request *client.GetFileByPathRequest, runtime *client.RuntimeOptions) (*client.GetFileByPathModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetFileByPathModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetFileByPathEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetFileByPathEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetFileByPathEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetFileByPathEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetFileByPathEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetFileByPathEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetFileByPathEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetFileByPathEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetFileEx(ctx context.Context, request *client.GetFileRequest, runtime *client.RuntimeOptions) (*client.GetFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetFileEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetFileEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetFileEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetFileEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetFileEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetFileEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetFileEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetFileEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetLastCursor(ctx context.Context, request *client.GetLastCursorRequest) (*client.GetLastCursorModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetLastCursorModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetLastCursor", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetLastCursor", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetLastCursor", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetLastCursor")
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
			w.inflightMetric.WithLabelValues("client.Client.GetLastCursor", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetLastCursor", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetLastCursor", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetLastCursor(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetLastCursorEx(ctx context.Context, request *client.GetLastCursorRequest, runtime *client.RuntimeOptions) (*client.GetLastCursorModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetLastCursorModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetLastCursorEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetLastCursorEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetLastCursorEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetLastCursorEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetLastCursorEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetLastCursorEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetLastCursorEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetLastCursorEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetLinkInfo(ctx context.Context, request *client.GetByLinkInfoRequest) (*client.GetLinkInfoModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetLinkInfoModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetLinkInfo", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetLinkInfo", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetLinkInfo", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetLinkInfo")
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
			w.inflightMetric.WithLabelValues("client.Client.GetLinkInfo", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetLinkInfo", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetLinkInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetLinkInfo(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetLinkInfoByUserId(ctx context.Context, request *client.GetLinkInfoByUserIDRequest) (*client.GetLinkInfoByUserIdModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetLinkInfoByUserIdModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetLinkInfoByUserId", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetLinkInfoByUserId", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetLinkInfoByUserId", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetLinkInfoByUserId")
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
			w.inflightMetric.WithLabelValues("client.Client.GetLinkInfoByUserId", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetLinkInfoByUserId", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetLinkInfoByUserId", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetLinkInfoByUserId(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetLinkInfoByUserIdEx(ctx context.Context, request *client.GetLinkInfoByUserIDRequest, runtime *client.RuntimeOptions) (*client.GetLinkInfoByUserIdModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetLinkInfoByUserIdModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetLinkInfoByUserIdEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetLinkInfoByUserIdEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetLinkInfoByUserIdEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetLinkInfoByUserIdEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetLinkInfoByUserIdEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetLinkInfoByUserIdEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetLinkInfoByUserIdEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetLinkInfoByUserIdEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetLinkInfoEx(ctx context.Context, request *client.GetByLinkInfoRequest, runtime *client.RuntimeOptions) (*client.GetLinkInfoModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetLinkInfoModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetLinkInfoEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetLinkInfoEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetLinkInfoEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetLinkInfoEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetLinkInfoEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetLinkInfoEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetLinkInfoEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetLinkInfoEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetMediaPlayUrl(ctx context.Context, request *client.GetMediaPlayURLRequest) (*client.GetMediaPlayUrlModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetMediaPlayUrlModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetMediaPlayUrl", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetMediaPlayUrl", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetMediaPlayUrl", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetMediaPlayUrl")
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
			w.inflightMetric.WithLabelValues("client.Client.GetMediaPlayUrl", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetMediaPlayUrl", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetMediaPlayUrl", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetMediaPlayUrl(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetMediaPlayUrlEx(ctx context.Context, request *client.GetMediaPlayURLRequest, runtime *client.RuntimeOptions) (*client.GetMediaPlayUrlModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetMediaPlayUrlModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetMediaPlayUrlEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetMediaPlayUrlEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetMediaPlayUrlEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetMediaPlayUrlEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetMediaPlayUrlEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetMediaPlayUrlEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetMediaPlayUrlEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetMediaPlayUrlEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetOfficeEditUrl(ctx context.Context, request *client.GetOfficeEditUrlRequest) (*client.GetOfficeEditUrlModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetOfficeEditUrlModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetOfficeEditUrl", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetOfficeEditUrl", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetOfficeEditUrl", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetOfficeEditUrl")
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
			w.inflightMetric.WithLabelValues("client.Client.GetOfficeEditUrl", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetOfficeEditUrl", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetOfficeEditUrl", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetOfficeEditUrl(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetOfficeEditUrlEx(ctx context.Context, request *client.GetOfficeEditUrlRequest, runtime *client.RuntimeOptions) (*client.GetOfficeEditUrlModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetOfficeEditUrlModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetOfficeEditUrlEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetOfficeEditUrlEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetOfficeEditUrlEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetOfficeEditUrlEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetOfficeEditUrlEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetOfficeEditUrlEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetOfficeEditUrlEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetOfficeEditUrlEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetOfficePreviewUrl(ctx context.Context, request *client.GetOfficePreviewUrlRequest) (*client.GetOfficePreviewUrlModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetOfficePreviewUrlModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetOfficePreviewUrl", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetOfficePreviewUrl", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetOfficePreviewUrl", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetOfficePreviewUrl")
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
			w.inflightMetric.WithLabelValues("client.Client.GetOfficePreviewUrl", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetOfficePreviewUrl", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetOfficePreviewUrl", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetOfficePreviewUrl(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetOfficePreviewUrlEx(ctx context.Context, request *client.GetOfficePreviewUrlRequest, runtime *client.RuntimeOptions) (*client.GetOfficePreviewUrlModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetOfficePreviewUrlModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetOfficePreviewUrlEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetOfficePreviewUrlEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetOfficePreviewUrlEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetOfficePreviewUrlEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetOfficePreviewUrlEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetOfficePreviewUrlEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetOfficePreviewUrlEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetOfficePreviewUrlEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetPathname(nickname *string, path *string) *string {
	_result := w.obj.GetPathname(nickname, path)
	return _result
}

func (w *PDSClientWrapper) GetPhotoCount(ctx context.Context, request *client.GetImageCountRequest) (*client.GetPhotoCountModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetPhotoCountModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetPhotoCount", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetPhotoCount", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetPhotoCount", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetPhotoCount")
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
			w.inflightMetric.WithLabelValues("client.Client.GetPhotoCount", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetPhotoCount", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetPhotoCount", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetPhotoCount(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetPhotoCountEx(ctx context.Context, request *client.GetImageCountRequest, runtime *client.RuntimeOptions) (*client.GetPhotoCountModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetPhotoCountModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetPhotoCountEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetPhotoCountEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetPhotoCountEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetPhotoCountEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetPhotoCountEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetPhotoCountEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetPhotoCountEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetPhotoCountEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetPublicKey(ctx context.Context, request *client.GetPublicKeyRequest) (*client.GetPublicKeyModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetPublicKeyModel
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
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetPublicKey")
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
			w.inflightMetric.WithLabelValues("client.Client.GetPublicKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetPublicKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetPublicKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetPublicKey(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetPublicKeyEx(ctx context.Context, request *client.GetPublicKeyRequest, runtime *client.RuntimeOptions) (*client.GetPublicKeyModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetPublicKeyModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetPublicKeyEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetPublicKeyEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetPublicKeyEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetPublicKeyEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetPublicKeyEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetPublicKeyEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetPublicKeyEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetPublicKeyEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetRefreshToken() *string {
	_result := w.obj.GetRefreshToken()
	return _result
}

func (w *PDSClientWrapper) GetSecurityToken(ctx context.Context) (*string, error) {
	ctxOptions := FromContext(ctx)
	var _result *string
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetSecurityToken", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetSecurityToken", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetSecurityToken", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetSecurityToken")
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
			w.inflightMetric.WithLabelValues("client.Client.GetSecurityToken", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetSecurityToken", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetSecurityToken", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetSecurityToken()
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetShare(ctx context.Context, request *client.GetShareRequest) (*client.GetShareModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetShareModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetShare", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetShare", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetShare", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetShare")
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
			w.inflightMetric.WithLabelValues("client.Client.GetShare", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetShare", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetShare", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetShare(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetShareByAnonymous(ctx context.Context, request *client.GetShareLinkByAnonymousRequest) (*client.GetShareByAnonymousModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetShareByAnonymousModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetShareByAnonymous", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetShareByAnonymous", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetShareByAnonymous", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetShareByAnonymous")
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
			w.inflightMetric.WithLabelValues("client.Client.GetShareByAnonymous", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetShareByAnonymous", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetShareByAnonymous", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetShareByAnonymous(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetShareByAnonymousEx(ctx context.Context, request *client.GetShareLinkByAnonymousRequest, runtime *client.RuntimeOptions) (*client.GetShareByAnonymousModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetShareByAnonymousModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetShareByAnonymousEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetShareByAnonymousEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetShareByAnonymousEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetShareByAnonymousEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetShareByAnonymousEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetShareByAnonymousEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetShareByAnonymousEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetShareByAnonymousEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetShareEx(ctx context.Context, request *client.GetShareRequest, runtime *client.RuntimeOptions) (*client.GetShareModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetShareModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetShareEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetShareEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetShareEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetShareEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetShareEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetShareEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetShareEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetShareEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetShareId(ctx context.Context, request *client.GetShareLinkIDRequest) (*client.GetShareIdModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetShareIdModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetShareId", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetShareId", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetShareId", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetShareId")
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
			w.inflightMetric.WithLabelValues("client.Client.GetShareId", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetShareId", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetShareId", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetShareId(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetShareIdEx(ctx context.Context, request *client.GetShareLinkIDRequest, runtime *client.RuntimeOptions) (*client.GetShareIdModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetShareIdModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetShareIdEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetShareIdEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetShareIdEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetShareIdEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetShareIdEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetShareIdEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetShareIdEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetShareIdEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetShareToken(ctx context.Context, request *client.GetShareLinkTokenRequest) (*client.GetShareTokenModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetShareTokenModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetShareToken", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetShareToken", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetShareToken", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetShareToken")
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
			w.inflightMetric.WithLabelValues("client.Client.GetShareToken", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetShareToken", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetShareToken", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetShareToken(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetShareTokenEx(ctx context.Context, request *client.GetShareLinkTokenRequest, runtime *client.RuntimeOptions) (*client.GetShareTokenModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetShareTokenModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetShareTokenEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetShareTokenEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetShareTokenEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetShareTokenEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetShareTokenEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetShareTokenEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetShareTokenEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetShareTokenEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetStoryDetail(ctx context.Context, request *client.GetStoryDetailRequest) (*client.GetStoryDetailModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetStoryDetailModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetStoryDetail", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetStoryDetail", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetStoryDetail", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetStoryDetail")
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
			w.inflightMetric.WithLabelValues("client.Client.GetStoryDetail", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetStoryDetail", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetStoryDetail", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetStoryDetail(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetStoryDetailEx(ctx context.Context, request *client.GetStoryDetailRequest, runtime *client.RuntimeOptions) (*client.GetStoryDetailModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetStoryDetailModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetStoryDetailEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetStoryDetailEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetStoryDetailEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetStoryDetailEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetStoryDetailEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetStoryDetailEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetStoryDetailEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetStoryDetailEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetUploadUrl(ctx context.Context, request *client.GetUploadUrlRequest) (*client.GetUploadUrlModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetUploadUrlModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetUploadUrl", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetUploadUrl", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetUploadUrl", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetUploadUrl")
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
			w.inflightMetric.WithLabelValues("client.Client.GetUploadUrl", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetUploadUrl", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetUploadUrl", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetUploadUrl(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetUploadUrlEx(ctx context.Context, request *client.GetUploadUrlRequest, runtime *client.RuntimeOptions) (*client.GetUploadUrlModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetUploadUrlModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetUploadUrlEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetUploadUrlEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetUploadUrlEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetUploadUrlEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetUploadUrlEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetUploadUrlEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetUploadUrlEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetUploadUrlEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetUser(ctx context.Context, request *client.GetUserRequest) (*client.GetUserModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetUserModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetUser", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetUser", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetUser", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetUser")
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
			w.inflightMetric.WithLabelValues("client.Client.GetUser", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetUser", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetUser", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetUser(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetUserAccessToken(ctx context.Context, request *client.GetUserAccessTokenRequest) (*client.GetUserAccessTokenModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetUserAccessTokenModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetUserAccessToken", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetUserAccessToken", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetUserAccessToken", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetUserAccessToken")
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
			w.inflightMetric.WithLabelValues("client.Client.GetUserAccessToken", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetUserAccessToken", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetUserAccessToken", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetUserAccessToken(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetUserAccessTokenEx(ctx context.Context, request *client.GetUserAccessTokenRequest, runtime *client.RuntimeOptions) (*client.GetUserAccessTokenModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetUserAccessTokenModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetUserAccessTokenEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetUserAccessTokenEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetUserAccessTokenEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetUserAccessTokenEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetUserAccessTokenEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetUserAccessTokenEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetUserAccessTokenEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetUserAccessTokenEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetUserAgent() *string {
	_result := w.obj.GetUserAgent()
	return _result
}

func (w *PDSClientWrapper) GetUserEx(ctx context.Context, request *client.GetUserRequest, runtime *client.RuntimeOptions) (*client.GetUserModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetUserModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetUserEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetUserEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetUserEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetUserEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetUserEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetUserEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetUserEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetUserEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetVideoPreviewSpriteUrl(ctx context.Context, request *client.GetVideoPreviewSpriteURLRequest) (*client.GetVideoPreviewSpriteUrlModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetVideoPreviewSpriteUrlModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetVideoPreviewSpriteUrl", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetVideoPreviewSpriteUrl", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetVideoPreviewSpriteUrl", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetVideoPreviewSpriteUrl")
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
			w.inflightMetric.WithLabelValues("client.Client.GetVideoPreviewSpriteUrl", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetVideoPreviewSpriteUrl", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetVideoPreviewSpriteUrl", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetVideoPreviewSpriteUrl(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetVideoPreviewSpriteUrlEx(ctx context.Context, request *client.GetVideoPreviewSpriteURLRequest, runtime *client.RuntimeOptions) (*client.GetVideoPreviewSpriteUrlModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetVideoPreviewSpriteUrlModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetVideoPreviewSpriteUrlEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetVideoPreviewSpriteUrlEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetVideoPreviewSpriteUrlEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetVideoPreviewSpriteUrlEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetVideoPreviewSpriteUrlEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetVideoPreviewSpriteUrlEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetVideoPreviewSpriteUrlEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetVideoPreviewSpriteUrlEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetVideoPreviewUrl(ctx context.Context, request *client.GetVideoPreviewURLRequest) (*client.GetVideoPreviewUrlModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetVideoPreviewUrlModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetVideoPreviewUrl", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetVideoPreviewUrl", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetVideoPreviewUrl", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetVideoPreviewUrl")
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
			w.inflightMetric.WithLabelValues("client.Client.GetVideoPreviewUrl", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetVideoPreviewUrl", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetVideoPreviewUrl", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetVideoPreviewUrl(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) GetVideoPreviewUrlEx(ctx context.Context, request *client.GetVideoPreviewURLRequest, runtime *client.RuntimeOptions) (*client.GetVideoPreviewUrlModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.GetVideoPreviewUrlModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetVideoPreviewUrlEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetVideoPreviewUrlEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetVideoPreviewUrlEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.GetVideoPreviewUrlEx")
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
			w.inflightMetric.WithLabelValues("client.Client.GetVideoPreviewUrlEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.GetVideoPreviewUrlEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.GetVideoPreviewUrlEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.GetVideoPreviewUrlEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) Init(ctx context.Context, config *client.Config) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.Init", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.Init", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.Init", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.Init")
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
			w.inflightMetric.WithLabelValues("client.Client.Init", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.Init", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.Init", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Init(config)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *PDSClientWrapper) Link(ctx context.Context, request *client.AccountLinkRequest) (*client.LinkModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.LinkModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.Link", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.Link", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.Link", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.Link")
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
			w.inflightMetric.WithLabelValues("client.Client.Link", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.Link", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.Link", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.Link(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) LinkEx(ctx context.Context, request *client.AccountLinkRequest, runtime *client.RuntimeOptions) (*client.LinkModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.LinkModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.LinkEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.LinkEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.LinkEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.LinkEx")
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
			w.inflightMetric.WithLabelValues("client.Client.LinkEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.LinkEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.LinkEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.LinkEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListAddressGroups(ctx context.Context, request *client.ListImageAddressGroupsRequest) (*client.ListAddressGroupsModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListAddressGroupsModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListAddressGroups", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListAddressGroups", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListAddressGroups", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListAddressGroups")
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
			w.inflightMetric.WithLabelValues("client.Client.ListAddressGroups", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListAddressGroups", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListAddressGroups", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListAddressGroups(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListAddressGroupsEx(ctx context.Context, request *client.ListImageAddressGroupsRequest, runtime *client.RuntimeOptions) (*client.ListAddressGroupsModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListAddressGroupsModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListAddressGroupsEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListAddressGroupsEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListAddressGroupsEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListAddressGroupsEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ListAddressGroupsEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListAddressGroupsEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListAddressGroupsEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListAddressGroupsEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListDrives(ctx context.Context, request *client.ListDriveRequest) (*client.ListDrivesModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListDrivesModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListDrives", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListDrives", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListDrives", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListDrives")
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
			w.inflightMetric.WithLabelValues("client.Client.ListDrives", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListDrives", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListDrives", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListDrives(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListDrivesEx(ctx context.Context, request *client.ListDriveRequest, runtime *client.RuntimeOptions) (*client.ListDrivesModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListDrivesModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListDrivesEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListDrivesEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListDrivesEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListDrivesEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ListDrivesEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListDrivesEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListDrivesEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListDrivesEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListFacegroups(ctx context.Context, request *client.ListImageFaceGroupsRequest) (*client.ListFacegroupsModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListFacegroupsModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListFacegroups", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListFacegroups", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListFacegroups", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListFacegroups")
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
			w.inflightMetric.WithLabelValues("client.Client.ListFacegroups", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListFacegroups", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListFacegroups", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListFacegroups(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListFacegroupsEx(ctx context.Context, request *client.ListImageFaceGroupsRequest, runtime *client.RuntimeOptions) (*client.ListFacegroupsModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListFacegroupsModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListFacegroupsEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListFacegroupsEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListFacegroupsEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListFacegroupsEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ListFacegroupsEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListFacegroupsEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListFacegroupsEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListFacegroupsEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListFile(ctx context.Context, request *client.ListFileRequest) (*client.ListFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListFile", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListFile", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListFile", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListFile")
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
			w.inflightMetric.WithLabelValues("client.Client.ListFile", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListFile", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListFile(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListFileByAnonymous(ctx context.Context, request *client.ListByAnonymousRequest) (*client.ListFileByAnonymousModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListFileByAnonymousModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListFileByAnonymous", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListFileByAnonymous", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListFileByAnonymous", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListFileByAnonymous")
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
			w.inflightMetric.WithLabelValues("client.Client.ListFileByAnonymous", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListFileByAnonymous", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListFileByAnonymous", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListFileByAnonymous(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListFileByAnonymousEx(ctx context.Context, request *client.ListByAnonymousRequest, runtime *client.RuntimeOptions) (*client.ListFileByAnonymousModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListFileByAnonymousModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListFileByAnonymousEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListFileByAnonymousEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListFileByAnonymousEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListFileByAnonymousEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ListFileByAnonymousEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListFileByAnonymousEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListFileByAnonymousEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListFileByAnonymousEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListFileByCustomIndexKey(ctx context.Context, request *client.ListFileByCustomIndexKeyRequest) (*client.ListFileByCustomIndexKeyModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListFileByCustomIndexKeyModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListFileByCustomIndexKey", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListFileByCustomIndexKey", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListFileByCustomIndexKey", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListFileByCustomIndexKey")
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
			w.inflightMetric.WithLabelValues("client.Client.ListFileByCustomIndexKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListFileByCustomIndexKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListFileByCustomIndexKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListFileByCustomIndexKey(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListFileByCustomIndexKeyEx(ctx context.Context, request *client.ListFileByCustomIndexKeyRequest, runtime *client.RuntimeOptions) (*client.ListFileByCustomIndexKeyModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListFileByCustomIndexKeyModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListFileByCustomIndexKeyEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListFileByCustomIndexKeyEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListFileByCustomIndexKeyEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListFileByCustomIndexKeyEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ListFileByCustomIndexKeyEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListFileByCustomIndexKeyEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListFileByCustomIndexKeyEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListFileByCustomIndexKeyEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListFileDelta(ctx context.Context, request *client.ListFileDeltaRequest) (*client.ListFileDeltaModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListFileDeltaModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListFileDelta", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListFileDelta", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListFileDelta", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListFileDelta")
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
			w.inflightMetric.WithLabelValues("client.Client.ListFileDelta", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListFileDelta", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListFileDelta", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListFileDelta(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListFileDeltaEx(ctx context.Context, request *client.ListFileDeltaRequest, runtime *client.RuntimeOptions) (*client.ListFileDeltaModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListFileDeltaModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListFileDeltaEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListFileDeltaEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListFileDeltaEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListFileDeltaEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ListFileDeltaEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListFileDeltaEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListFileDeltaEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListFileDeltaEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListFileEx(ctx context.Context, request *client.ListFileRequest, runtime *client.RuntimeOptions) (*client.ListFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListFileEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListFileEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListFileEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListFileEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ListFileEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListFileEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListFileEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListFileEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListMyDrives(ctx context.Context, request *client.ListMyDriveRequest) (*client.ListMyDrivesModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListMyDrivesModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListMyDrives", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListMyDrives", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListMyDrives", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListMyDrives")
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
			w.inflightMetric.WithLabelValues("client.Client.ListMyDrives", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListMyDrives", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListMyDrives", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListMyDrives(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListMyDrivesEx(ctx context.Context, request *client.ListMyDriveRequest, runtime *client.RuntimeOptions) (*client.ListMyDrivesModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListMyDrivesModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListMyDrivesEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListMyDrivesEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListMyDrivesEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListMyDrivesEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ListMyDrivesEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListMyDrivesEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListMyDrivesEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListMyDrivesEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListShare(ctx context.Context, request *client.ListShareRequest) (*client.ListShareModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListShareModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListShare", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListShare", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListShare", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListShare")
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
			w.inflightMetric.WithLabelValues("client.Client.ListShare", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListShare", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListShare", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListShare(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListShareEx(ctx context.Context, request *client.ListShareRequest, runtime *client.RuntimeOptions) (*client.ListShareModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListShareModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListShareEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListShareEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListShareEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListShareEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ListShareEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListShareEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListShareEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListShareEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListShareLink(ctx context.Context, request *client.ListShareLinkRequest) (*client.ListShareLinkModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListShareLinkModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListShareLink", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListShareLink", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListShareLink", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListShareLink")
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
			w.inflightMetric.WithLabelValues("client.Client.ListShareLink", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListShareLink", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListShareLink", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListShareLink(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListShareLinkEx(ctx context.Context, request *client.ListShareLinkRequest, runtime *client.RuntimeOptions) (*client.ListShareLinkModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListShareLinkModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListShareLinkEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListShareLinkEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListShareLinkEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListShareLinkEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ListShareLinkEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListShareLinkEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListShareLinkEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListShareLinkEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListStory(ctx context.Context, request *client.ListStoryRequest) (*client.ListStoryModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListStoryModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListStory", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListStory", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListStory", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListStory")
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
			w.inflightMetric.WithLabelValues("client.Client.ListStory", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListStory", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListStory", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListStory(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListStoryEx(ctx context.Context, request *client.ListStoryRequest, runtime *client.RuntimeOptions) (*client.ListStoryModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListStoryModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListStoryEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListStoryEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListStoryEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListStoryEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ListStoryEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListStoryEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListStoryEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListStoryEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListTags(ctx context.Context, request *client.ListImageTagsRequest) (*client.ListTagsModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListTagsModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListTags", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListTags", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListTags", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListTags")
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
			w.inflightMetric.WithLabelValues("client.Client.ListTags", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListTags", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListTags", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListTags(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListTagsEx(ctx context.Context, request *client.ListImageTagsRequest, runtime *client.RuntimeOptions) (*client.ListTagsModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListTagsModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListTagsEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListTagsEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListTagsEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListTagsEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ListTagsEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListTagsEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListTagsEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListTagsEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListUploadedParts(ctx context.Context, request *client.ListUploadedPartRequest) (*client.ListUploadedPartsModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListUploadedPartsModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListUploadedParts", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListUploadedParts", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListUploadedParts", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListUploadedParts")
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
			w.inflightMetric.WithLabelValues("client.Client.ListUploadedParts", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListUploadedParts", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListUploadedParts", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListUploadedParts(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListUploadedPartsEx(ctx context.Context, request *client.ListUploadedPartRequest, runtime *client.RuntimeOptions) (*client.ListUploadedPartsModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListUploadedPartsModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListUploadedPartsEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListUploadedPartsEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListUploadedPartsEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListUploadedPartsEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ListUploadedPartsEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListUploadedPartsEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListUploadedPartsEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListUploadedPartsEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListUsers(ctx context.Context, request *client.ListUserRequest) (*client.ListUsersModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListUsersModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListUsers", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListUsers", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListUsers", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListUsers")
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
			w.inflightMetric.WithLabelValues("client.Client.ListUsers", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListUsers", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListUsers", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListUsers(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ListUsersEx(ctx context.Context, request *client.ListUserRequest, runtime *client.RuntimeOptions) (*client.ListUsersModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ListUsersModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListUsersEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListUsersEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListUsersEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ListUsersEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ListUsersEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ListUsersEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ListUsersEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ListUsersEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) Login(ctx context.Context, request *client.MobileLoginRequest) (*client.LoginModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.LoginModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.Login", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.Login", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.Login", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.Login")
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
			w.inflightMetric.WithLabelValues("client.Client.Login", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.Login", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.Login", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.Login(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) LoginEx(ctx context.Context, request *client.MobileLoginRequest, runtime *client.RuntimeOptions) (*client.LoginModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.LoginModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.LoginEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.LoginEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.LoginEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.LoginEx")
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
			w.inflightMetric.WithLabelValues("client.Client.LoginEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.LoginEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.LoginEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.LoginEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) MobileSendSmsCode(ctx context.Context, request *client.MobileSendSmsCodeRequest) (*client.MobileSendSmsCodeModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.MobileSendSmsCodeModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.MobileSendSmsCode", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.MobileSendSmsCode", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.MobileSendSmsCode", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.MobileSendSmsCode")
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
			w.inflightMetric.WithLabelValues("client.Client.MobileSendSmsCode", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.MobileSendSmsCode", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.MobileSendSmsCode", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.MobileSendSmsCode(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) MobileSendSmsCodeEx(ctx context.Context, request *client.MobileSendSmsCodeRequest, runtime *client.RuntimeOptions) (*client.MobileSendSmsCodeModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.MobileSendSmsCodeModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.MobileSendSmsCodeEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.MobileSendSmsCodeEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.MobileSendSmsCodeEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.MobileSendSmsCodeEx")
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
			w.inflightMetric.WithLabelValues("client.Client.MobileSendSmsCodeEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.MobileSendSmsCodeEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.MobileSendSmsCodeEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.MobileSendSmsCodeEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) MoveFile(ctx context.Context, request *client.MoveFileRequest) (*client.MoveFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.MoveFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.MoveFile", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.MoveFile", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.MoveFile", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.MoveFile")
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
			w.inflightMetric.WithLabelValues("client.Client.MoveFile", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.MoveFile", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.MoveFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.MoveFile(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) MoveFileEx(ctx context.Context, request *client.MoveFileRequest, runtime *client.RuntimeOptions) (*client.MoveFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.MoveFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.MoveFileEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.MoveFileEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.MoveFileEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.MoveFileEx")
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
			w.inflightMetric.WithLabelValues("client.Client.MoveFileEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.MoveFileEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.MoveFileEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.MoveFileEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ParseKeywords(ctx context.Context, request *client.ParseKeywordsRequest) (*client.ParseKeywordsModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ParseKeywordsModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ParseKeywords", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ParseKeywords", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ParseKeywords", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ParseKeywords")
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
			w.inflightMetric.WithLabelValues("client.Client.ParseKeywords", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ParseKeywords", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ParseKeywords", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ParseKeywords(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ParseKeywordsEx(ctx context.Context, request *client.ParseKeywordsRequest, runtime *client.RuntimeOptions) (*client.ParseKeywordsModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ParseKeywordsModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ParseKeywordsEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ParseKeywordsEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ParseKeywordsEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ParseKeywordsEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ParseKeywordsEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ParseKeywordsEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ParseKeywordsEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ParseKeywordsEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) Register(ctx context.Context, request *client.MobileRegisterRequest) (*client.RegisterModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.RegisterModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.Register", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.Register", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.Register", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.Register")
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
			w.inflightMetric.WithLabelValues("client.Client.Register", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.Register", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.Register", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.Register(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) RegisterEx(ctx context.Context, request *client.MobileRegisterRequest, runtime *client.RuntimeOptions) (*client.RegisterModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.RegisterModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.RegisterEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.RegisterEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.RegisterEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.RegisterEx")
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
			w.inflightMetric.WithLabelValues("client.Client.RegisterEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.RegisterEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.RegisterEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.RegisterEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) RemoveStoryImages(ctx context.Context, request *client.RemoveStoryImagesRequest) (*client.RemoveStoryImagesModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.RemoveStoryImagesModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.RemoveStoryImages", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.RemoveStoryImages", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.RemoveStoryImages", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.RemoveStoryImages")
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
			w.inflightMetric.WithLabelValues("client.Client.RemoveStoryImages", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.RemoveStoryImages", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.RemoveStoryImages", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.RemoveStoryImages(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) RemoveStoryImagesEx(ctx context.Context, request *client.RemoveStoryImagesRequest, runtime *client.RuntimeOptions) (*client.RemoveStoryImagesModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.RemoveStoryImagesModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.RemoveStoryImagesEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.RemoveStoryImagesEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.RemoveStoryImagesEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.RemoveStoryImagesEx")
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
			w.inflightMetric.WithLabelValues("client.Client.RemoveStoryImagesEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.RemoveStoryImagesEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.RemoveStoryImagesEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.RemoveStoryImagesEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ScanFileMeta(ctx context.Context, request *client.ScanFileMetaRequest) (*client.ScanFileMetaModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ScanFileMetaModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ScanFileMeta", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ScanFileMeta", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ScanFileMeta", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ScanFileMeta")
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
			w.inflightMetric.WithLabelValues("client.Client.ScanFileMeta", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ScanFileMeta", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ScanFileMeta", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ScanFileMeta(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) ScanFileMetaEx(ctx context.Context, request *client.ScanFileMetaRequest, runtime *client.RuntimeOptions) (*client.ScanFileMetaModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.ScanFileMetaModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ScanFileMetaEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ScanFileMetaEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ScanFileMetaEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.ScanFileMetaEx")
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
			w.inflightMetric.WithLabelValues("client.Client.ScanFileMetaEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.ScanFileMetaEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.ScanFileMetaEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.ScanFileMetaEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) SearchAddressGroups(ctx context.Context, request *client.SearchImageAddressGroupsRequest) (*client.SearchAddressGroupsModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.SearchAddressGroupsModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SearchAddressGroups", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SearchAddressGroups", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SearchAddressGroups", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.SearchAddressGroups")
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
			w.inflightMetric.WithLabelValues("client.Client.SearchAddressGroups", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.SearchAddressGroups", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.SearchAddressGroups", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.SearchAddressGroups(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) SearchAddressGroupsEx(ctx context.Context, request *client.SearchImageAddressGroupsRequest, runtime *client.RuntimeOptions) (*client.SearchAddressGroupsModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.SearchAddressGroupsModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SearchAddressGroupsEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SearchAddressGroupsEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SearchAddressGroupsEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.SearchAddressGroupsEx")
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
			w.inflightMetric.WithLabelValues("client.Client.SearchAddressGroupsEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.SearchAddressGroupsEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.SearchAddressGroupsEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.SearchAddressGroupsEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) SearchFile(ctx context.Context, request *client.SearchFileRequest) (*client.SearchFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.SearchFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SearchFile", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SearchFile", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SearchFile", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.SearchFile")
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
			w.inflightMetric.WithLabelValues("client.Client.SearchFile", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.SearchFile", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.SearchFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.SearchFile(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) SearchFileEx(ctx context.Context, request *client.SearchFileRequest, runtime *client.RuntimeOptions) (*client.SearchFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.SearchFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SearchFileEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SearchFileEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SearchFileEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.SearchFileEx")
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
			w.inflightMetric.WithLabelValues("client.Client.SearchFileEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.SearchFileEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.SearchFileEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.SearchFileEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) SearchUser(ctx context.Context, request *client.SearchUserRequest) (*client.SearchUserModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.SearchUserModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SearchUser", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SearchUser", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SearchUser", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.SearchUser")
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
			w.inflightMetric.WithLabelValues("client.Client.SearchUser", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.SearchUser", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.SearchUser", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.SearchUser(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) SearchUserEx(ctx context.Context, request *client.SearchUserRequest, runtime *client.RuntimeOptions) (*client.SearchUserModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.SearchUserModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SearchUserEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SearchUserEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SearchUserEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.SearchUserEx")
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
			w.inflightMetric.WithLabelValues("client.Client.SearchUserEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.SearchUserEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.SearchUserEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.SearchUserEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) SetAccessToken(token *string) {
	w.obj.SetAccessToken(token)
}

func (w *PDSClientWrapper) SetExpireTime(ctx context.Context, expireTime *string) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SetExpireTime", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SetExpireTime", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SetExpireTime", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.SetExpireTime")
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
			w.inflightMetric.WithLabelValues("client.Client.SetExpireTime", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.SetExpireTime", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.SetExpireTime", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetExpireTime(expireTime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *PDSClientWrapper) SetPassword(ctx context.Context, request *client.DefaultSetPasswordRequest) (*client.SetPasswordModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.SetPasswordModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SetPassword", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SetPassword", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SetPassword", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.SetPassword")
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
			w.inflightMetric.WithLabelValues("client.Client.SetPassword", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.SetPassword", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.SetPassword", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.SetPassword(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) SetPasswordEx(ctx context.Context, request *client.DefaultSetPasswordRequest, runtime *client.RuntimeOptions) (*client.SetPasswordModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.SetPasswordModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SetPasswordEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SetPasswordEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SetPasswordEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.SetPasswordEx")
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
			w.inflightMetric.WithLabelValues("client.Client.SetPasswordEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.SetPasswordEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.SetPasswordEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.SetPasswordEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) SetRefreshToken(token *string) {
	w.obj.SetRefreshToken(token)
}

func (w *PDSClientWrapper) SetUserAgent(userAgent *string) {
	w.obj.SetUserAgent(userAgent)
}

func (w *PDSClientWrapper) Token(ctx context.Context, request *client.RefreshOfficeEditTokenRequest) (*client.TokenModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.TokenModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.Token", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.Token", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.Token", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.Token")
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
			w.inflightMetric.WithLabelValues("client.Client.Token", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.Token", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.Token", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.Token(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) TokenEx(ctx context.Context, request *client.RefreshOfficeEditTokenRequest, runtime *client.RuntimeOptions) (*client.TokenModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.TokenModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.TokenEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.TokenEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.TokenEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.TokenEx")
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
			w.inflightMetric.WithLabelValues("client.Client.TokenEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.TokenEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.TokenEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.TokenEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) UpdateDrive(ctx context.Context, request *client.UpdateDriveRequest) (*client.UpdateDriveModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.UpdateDriveModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateDrive", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateDrive", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateDrive", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.UpdateDrive")
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
			w.inflightMetric.WithLabelValues("client.Client.UpdateDrive", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.UpdateDrive", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.UpdateDrive", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.UpdateDrive(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) UpdateDriveEx(ctx context.Context, request *client.UpdateDriveRequest, runtime *client.RuntimeOptions) (*client.UpdateDriveModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.UpdateDriveModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateDriveEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateDriveEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateDriveEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.UpdateDriveEx")
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
			w.inflightMetric.WithLabelValues("client.Client.UpdateDriveEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.UpdateDriveEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.UpdateDriveEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.UpdateDriveEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) UpdateFacegroupInfo(ctx context.Context, request *client.UpdateFaceGroupInfoRequest) (*client.UpdateFacegroupInfoModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.UpdateFacegroupInfoModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateFacegroupInfo", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateFacegroupInfo", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateFacegroupInfo", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.UpdateFacegroupInfo")
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
			w.inflightMetric.WithLabelValues("client.Client.UpdateFacegroupInfo", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.UpdateFacegroupInfo", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.UpdateFacegroupInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.UpdateFacegroupInfo(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) UpdateFacegroupInfoEx(ctx context.Context, request *client.UpdateFaceGroupInfoRequest, runtime *client.RuntimeOptions) (*client.UpdateFacegroupInfoModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.UpdateFacegroupInfoModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateFacegroupInfoEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateFacegroupInfoEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateFacegroupInfoEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.UpdateFacegroupInfoEx")
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
			w.inflightMetric.WithLabelValues("client.Client.UpdateFacegroupInfoEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.UpdateFacegroupInfoEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.UpdateFacegroupInfoEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.UpdateFacegroupInfoEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) UpdateFile(ctx context.Context, request *client.UpdateFileMetaRequest) (*client.UpdateFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.UpdateFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateFile", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateFile", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateFile", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.UpdateFile")
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
			w.inflightMetric.WithLabelValues("client.Client.UpdateFile", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.UpdateFile", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.UpdateFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.UpdateFile(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) UpdateFileEx(ctx context.Context, request *client.UpdateFileMetaRequest, runtime *client.RuntimeOptions) (*client.UpdateFileModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.UpdateFileModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateFileEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateFileEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateFileEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.UpdateFileEx")
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
			w.inflightMetric.WithLabelValues("client.Client.UpdateFileEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.UpdateFileEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.UpdateFileEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.UpdateFileEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) UpdateShare(ctx context.Context, request *client.UpdateShareRequest) (*client.UpdateShareModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.UpdateShareModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateShare", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateShare", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateShare", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.UpdateShare")
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
			w.inflightMetric.WithLabelValues("client.Client.UpdateShare", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.UpdateShare", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.UpdateShare", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.UpdateShare(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) UpdateShareEx(ctx context.Context, request *client.UpdateShareRequest, runtime *client.RuntimeOptions) (*client.UpdateShareModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.UpdateShareModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateShareEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateShareEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateShareEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.UpdateShareEx")
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
			w.inflightMetric.WithLabelValues("client.Client.UpdateShareEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.UpdateShareEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.UpdateShareEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.UpdateShareEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) UpdateUser(ctx context.Context, request *client.UpdateUserRequest) (*client.UpdateUserModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.UpdateUserModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateUser", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateUser", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateUser", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.UpdateUser")
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
			w.inflightMetric.WithLabelValues("client.Client.UpdateUser", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.UpdateUser", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.UpdateUser", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.UpdateUser(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) UpdateUserEx(ctx context.Context, request *client.UpdateUserRequest, runtime *client.RuntimeOptions) (*client.UpdateUserModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.UpdateUserModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateUserEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateUserEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateUserEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.UpdateUserEx")
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
			w.inflightMetric.WithLabelValues("client.Client.UpdateUserEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.UpdateUserEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.UpdateUserEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.UpdateUserEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) VerifyCode(ctx context.Context, request *client.VerifyCodeRequest) (*client.VerifyCodeModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.VerifyCodeModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.VerifyCode", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.VerifyCode", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.VerifyCode", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.VerifyCode")
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
			w.inflightMetric.WithLabelValues("client.Client.VerifyCode", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.VerifyCode", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.VerifyCode", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.VerifyCode(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}

func (w *PDSClientWrapper) VerifyCodeEx(ctx context.Context, request *client.VerifyCodeRequest, runtime *client.RuntimeOptions) (*client.VerifyCodeModel, error) {
	ctxOptions := FromContext(ctx)
	var _result *client.VerifyCodeModel
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.VerifyCodeEx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.VerifyCodeEx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.VerifyCodeEx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "client.Client.VerifyCodeEx")
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
			w.inflightMetric.WithLabelValues("client.Client.VerifyCodeEx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("client.Client.VerifyCodeEx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("client.Client.VerifyCodeEx", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		_result, err = w.obj.VerifyCodeEx(request, runtime)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return _result, err
}
