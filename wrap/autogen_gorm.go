// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

type GORMAssociationWrapper struct {
	obj                *gorm.Association
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w GORMAssociationWrapper) Unwrap() *gorm.Association {
	return w.obj
}

func NewGORMDBWrapper(
	obj *gorm.DB,
	retry *micro.Retry,
	options *WrapperOptions,
	durationMetric *prometheus.HistogramVec,
	inflightMetric *prometheus.GaugeVec,
	rateLimiter micro.RateLimiter,
	parallelController micro.ParallelController) *GORMDBWrapper {
	return &GORMDBWrapper{
		obj:                obj,
		retry:              retry,
		options:            options,
		durationMetric:     durationMetric,
		inflightMetric:     inflightMetric,
		rateLimiter:        rateLimiter,
		parallelController: parallelController,
	}
}

type GORMDBWrapper struct {
	obj                *gorm.DB
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w GORMDBWrapper) Unwrap() *gorm.DB {
	return w.obj
}

func (w *GORMDBWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *GORMDBWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *GORMDBWrapper) OnRateLimiterChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *GORMDBWrapper) OnParallelControllerChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *GORMDBWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        fmt.Sprintf("%s_gorm_DB_durationMs", options.Name),
		Help:        "gorm DB response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        fmt.Sprintf("%s_gorm_DB_inflight", options.Name),
		Help:        "gorm DB inflight",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "custom"})
}

func (w GORMAssociationWrapper) Append(ctx context.Context, values ...interface{}) GORMAssociationWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.Association
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Association.Append", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Association.Append", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Association.Append", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.Association.Append", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.Association.Append", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.Association.Append", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.Association.Append", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Append(values...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMAssociationWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMAssociationWrapper) Clear(ctx context.Context) GORMAssociationWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.Association
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Association.Clear", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Association.Clear", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Association.Clear", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.Association.Clear", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.Association.Clear", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.Association.Clear", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.Association.Clear", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Clear()
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMAssociationWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMAssociationWrapper) Count() int {
	res0 := w.obj.Count()
	return res0
}

func (w GORMAssociationWrapper) Delete(ctx context.Context, values ...interface{}) GORMAssociationWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.Association
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Association.Delete", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Association.Delete", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Association.Delete", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.Association.Delete", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.Association.Delete", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.Association.Delete", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.Association.Delete", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Delete(values...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMAssociationWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMAssociationWrapper) Find(ctx context.Context, value interface{}) GORMAssociationWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.Association
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Association.Find", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Association.Find", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Association.Find", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.Association.Find", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.Association.Find", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.Association.Find", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.Association.Find", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Find(value)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMAssociationWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMAssociationWrapper) Replace(ctx context.Context, values ...interface{}) GORMAssociationWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.Association
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Association.Replace", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Association.Replace", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Association.Replace", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.Association.Replace", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.Association.Replace", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.Association.Replace", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.Association.Replace", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Replace(values...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMAssociationWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) AddError(err error) error {
	res0 := w.obj.AddError(err)
	return res0
}

func (w GORMDBWrapper) AddForeignKey(ctx context.Context, field string, dest string, onDelete string, onUpdate string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.AddForeignKey", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.AddForeignKey", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.AddForeignKey", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.AddForeignKey", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.AddForeignKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.AddForeignKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.AddForeignKey", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.AddForeignKey(field, dest, onDelete, onUpdate)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) AddIndex(ctx context.Context, indexName string, columns ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.AddIndex", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.AddIndex", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.AddIndex", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.AddIndex", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.AddIndex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.AddIndex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.AddIndex", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.AddIndex(indexName, columns...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) AddUniqueIndex(ctx context.Context, indexName string, columns ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.AddUniqueIndex", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.AddUniqueIndex", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.AddUniqueIndex", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.AddUniqueIndex", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.AddUniqueIndex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.AddUniqueIndex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.AddUniqueIndex", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.AddUniqueIndex(indexName, columns...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Assign(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Assign", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Assign", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Assign", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Assign", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Assign", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Assign", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Assign", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Assign(attrs...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Association(ctx context.Context, column string) GORMAssociationWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.Association
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Association", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Association", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Association", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Association", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Association", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Association", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Association", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Association(column)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMAssociationWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Attrs(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Attrs", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Attrs", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Attrs", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Attrs", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Attrs", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Attrs", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Attrs", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Attrs(attrs...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) AutoMigrate(ctx context.Context, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.AutoMigrate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.AutoMigrate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.AutoMigrate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.AutoMigrate", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.AutoMigrate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.AutoMigrate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.AutoMigrate", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.AutoMigrate(values...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Begin(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Begin", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Begin", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Begin", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Begin", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Begin", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Begin", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Begin", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Begin()
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) BeginTx(ctx context.Context, opts *sql.TxOptions) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.BeginTx", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.BeginTx", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.BeginTx", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.BeginTx", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.BeginTx", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.BeginTx", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.BeginTx", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BeginTx(ctx, opts)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) BlockGlobalUpdate(ctx context.Context, enable bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.BlockGlobalUpdate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.BlockGlobalUpdate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.BlockGlobalUpdate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.BlockGlobalUpdate", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.BlockGlobalUpdate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.BlockGlobalUpdate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.BlockGlobalUpdate", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BlockGlobalUpdate(enable)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Callback() *gorm.Callback {
	res0 := w.obj.Callback()
	return res0
}

func (w GORMDBWrapper) Close(ctx context.Context) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Close", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Close", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Close", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Close", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Close", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Close", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Close", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Close()
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w GORMDBWrapper) Commit(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Commit", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Commit", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Commit", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Commit", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Commit", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Commit", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Commit", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Commit()
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) CommonDB() gorm.SQLCommon {
	res0 := w.obj.CommonDB()
	return res0
}

func (w GORMDBWrapper) Count(ctx context.Context, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Count", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Count", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Count", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Count", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Count", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Count", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Count", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Count(value)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Create(ctx context.Context, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Create", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Create", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Create", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Create", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Create", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Create", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Create", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Create(value)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) CreateTable(ctx context.Context, models ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.CreateTable", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.CreateTable", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.CreateTable", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.CreateTable", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.CreateTable", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.CreateTable", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.CreateTable", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.CreateTable(models...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) DB() *sql.DB {
	res0 := w.obj.DB()
	return res0
}

func (w GORMDBWrapper) Debug(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Debug", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Debug", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Debug", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Debug", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Debug", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Debug", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Debug", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Debug()
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Delete(ctx context.Context, value interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Delete", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Delete", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Delete", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Delete", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Delete", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Delete", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Delete", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Delete(value, where...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Dialect() gorm.Dialect {
	res0 := w.obj.Dialect()
	return res0
}

func (w GORMDBWrapper) DropColumn(ctx context.Context, column string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.DropColumn", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.DropColumn", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.DropColumn", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.DropColumn", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.DropColumn", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.DropColumn", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.DropColumn", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.DropColumn(column)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) DropTable(ctx context.Context, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.DropTable", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.DropTable", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.DropTable", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.DropTable", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.DropTable", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.DropTable", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.DropTable", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.DropTable(values...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) DropTableIfExists(ctx context.Context, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.DropTableIfExists", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.DropTableIfExists", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.DropTableIfExists", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.DropTableIfExists", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.DropTableIfExists", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.DropTableIfExists", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.DropTableIfExists", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.DropTableIfExists(values...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Exec(ctx context.Context, sql string, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Exec", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Exec", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Exec", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Exec", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Exec", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Exec", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Exec", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Exec(sql, values...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Find(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Find", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Find", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Find", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Find", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Find", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Find", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Find", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Find(out, where...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) First(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.First", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.First", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.First", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.First", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.First", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.First", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.First", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.First(out, where...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) FirstOrCreate(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.FirstOrCreate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.FirstOrCreate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.FirstOrCreate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.FirstOrCreate", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.FirstOrCreate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.FirstOrCreate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.FirstOrCreate", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.FirstOrCreate(out, where...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) FirstOrInit(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.FirstOrInit", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.FirstOrInit", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.FirstOrInit", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.FirstOrInit", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.FirstOrInit", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.FirstOrInit", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.FirstOrInit", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.FirstOrInit(out, where...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Get(name string) (interface{}, bool) {
	value, ok := w.obj.Get(name)
	return value, ok
}

func (w GORMDBWrapper) GetErrors() []error {
	res0 := w.obj.GetErrors()
	return res0
}

func (w GORMDBWrapper) Group(ctx context.Context, query string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Group", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Group", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Group", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Group", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Group", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Group", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Group", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Group(query)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) HasBlockGlobalUpdate() bool {
	res0 := w.obj.HasBlockGlobalUpdate()
	return res0
}

func (w GORMDBWrapper) HasTable(value interface{}) bool {
	res0 := w.obj.HasTable(value)
	return res0
}

func (w GORMDBWrapper) Having(ctx context.Context, query interface{}, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Having", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Having", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Having", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Having", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Having", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Having", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Having", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Having(query, values...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) InstantSet(ctx context.Context, name string, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.InstantSet", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.InstantSet", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.InstantSet", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.InstantSet", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.InstantSet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.InstantSet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.InstantSet", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.InstantSet(name, value)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Joins(ctx context.Context, query string, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Joins", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Joins", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Joins", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Joins", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Joins", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Joins", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Joins", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Joins(query, args...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Last(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Last", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Last", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Last", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Last", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Last", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Last", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Last", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Last(out, where...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Limit(ctx context.Context, limit interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Limit", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Limit", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Limit", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Limit", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Limit", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Limit", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Limit", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Limit(limit)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) LogMode(ctx context.Context, enable bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.LogMode", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.LogMode", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.LogMode", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.LogMode", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.LogMode", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.LogMode", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.LogMode", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LogMode(enable)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Model(ctx context.Context, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Model", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Model", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Model", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Model", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Model", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Model", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Model", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Model(value)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) ModifyColumn(ctx context.Context, column string, typ string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.ModifyColumn", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.ModifyColumn", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.ModifyColumn", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.ModifyColumn", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.ModifyColumn", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.ModifyColumn", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.ModifyColumn", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ModifyColumn(column, typ)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) New(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.New", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.New", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.New", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.New", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.New", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.New", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.New", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.New()
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) NewRecord(value interface{}) bool {
	res0 := w.obj.NewRecord(value)
	return res0
}

func (w GORMDBWrapper) NewScope(value interface{}) *gorm.Scope {
	res0 := w.obj.NewScope(value)
	return res0
}

func (w GORMDBWrapper) Not(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Not", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Not", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Not", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Not", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Not", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Not", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Not", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Not(query, args...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Offset(ctx context.Context, offset interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Offset", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Offset", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Offset", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Offset", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Offset", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Offset", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Offset", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Offset(offset)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Omit(ctx context.Context, columns ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Omit", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Omit", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Omit", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Omit", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Omit", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Omit", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Omit", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Omit(columns...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Or(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Or", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Or", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Or", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Or", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Or", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Or", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Or", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Or(query, args...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Order(ctx context.Context, value interface{}, reorder ...bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Order", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Order", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Order", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Order", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Order", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Order", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Order", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Order(value, reorder...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Pluck(ctx context.Context, column string, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Pluck", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Pluck", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Pluck", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Pluck", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Pluck", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Pluck", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Pluck", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Pluck(column, value)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Preload(ctx context.Context, column string, conditions ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Preload", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Preload", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Preload", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Preload", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Preload", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Preload", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Preload", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Preload(column, conditions...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Preloads(ctx context.Context, out interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Preloads", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Preloads", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Preloads", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Preloads", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Preloads", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Preloads", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Preloads", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Preloads(out)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) QueryExpr() *gorm.SqlExpr {
	res0 := w.obj.QueryExpr()
	return res0
}

func (w GORMDBWrapper) Raw(ctx context.Context, sql string, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Raw", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Raw", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Raw", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Raw", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Raw", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Raw", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Raw", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Raw(sql, values...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) RecordNotFound() bool {
	res0 := w.obj.RecordNotFound()
	return res0
}

func (w GORMDBWrapper) Related(ctx context.Context, value interface{}, foreignKeys ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Related", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Related", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Related", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Related", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Related", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Related", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Related", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Related(value, foreignKeys...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) RemoveForeignKey(ctx context.Context, field string, dest string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.RemoveForeignKey", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.RemoveForeignKey", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.RemoveForeignKey", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.RemoveForeignKey", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.RemoveForeignKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.RemoveForeignKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.RemoveForeignKey", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RemoveForeignKey(field, dest)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) RemoveIndex(ctx context.Context, indexName string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.RemoveIndex", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.RemoveIndex", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.RemoveIndex", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.RemoveIndex", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.RemoveIndex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.RemoveIndex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.RemoveIndex", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RemoveIndex(indexName)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Rollback(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Rollback", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Rollback", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Rollback", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Rollback", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Rollback", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Rollback", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Rollback", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Rollback()
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) RollbackUnlessCommitted(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.RollbackUnlessCommitted", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.RollbackUnlessCommitted", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.RollbackUnlessCommitted", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.RollbackUnlessCommitted", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.RollbackUnlessCommitted", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.RollbackUnlessCommitted", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.RollbackUnlessCommitted", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RollbackUnlessCommitted()
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Row() *sql.Row {
	res0 := w.obj.Row()
	return res0
}

func (w GORMDBWrapper) Rows(ctx context.Context) (*sql.Rows, error) {
	ctxOptions := FromContext(ctx)
	var res0 *sql.Rows
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Rows", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Rows", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Rows", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Rows", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Rows", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Rows", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Rows", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Rows()
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w GORMDBWrapper) Save(ctx context.Context, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Save", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Save", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Save", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Save", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Save", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Save", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Save", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Save(value)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Scan(ctx context.Context, dest interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Scan", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Scan", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Scan", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Scan", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Scan", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Scan", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Scan", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Scan(dest)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) ScanRows(ctx context.Context, rows *sql.Rows, result interface{}) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.ScanRows", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.ScanRows", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.ScanRows", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.ScanRows", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.ScanRows", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.ScanRows", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.ScanRows", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.ScanRows(rows, result)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w GORMDBWrapper) Scopes(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Scopes", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Scopes", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Scopes", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Scopes", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Scopes", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Scopes", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Scopes", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Scopes(funcs...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Select(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Select", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Select", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Select", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Select", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Select", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Select", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Select", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Select(query, args...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Set(ctx context.Context, name string, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Set", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Set", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Set", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Set", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Set", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Set", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Set", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Set(name, value)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) SetJoinTableHandler(source interface{}, column string, handler gorm.JoinTableHandlerInterface) {
	w.obj.SetJoinTableHandler(source, column, handler)
}

func (w GORMDBWrapper) SetNowFuncOverride(ctx context.Context, nowFuncOverride func() time.Time) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.SetNowFuncOverride", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.SetNowFuncOverride", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.SetNowFuncOverride", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.SetNowFuncOverride", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.SetNowFuncOverride", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.SetNowFuncOverride", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.SetNowFuncOverride", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SetNowFuncOverride(nowFuncOverride)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) SingularTable(enable bool) {
	w.obj.SingularTable(enable)
}

func (w GORMDBWrapper) SubQuery() *gorm.SqlExpr {
	res0 := w.obj.SubQuery()
	return res0
}

func (w GORMDBWrapper) Table(ctx context.Context, name string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Table", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Table", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Table", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Table", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Table", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Table", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Table", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Table(name)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Take(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Take", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Take", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Take", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Take", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Take", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Take", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Take", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Take(out, where...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Transaction(ctx context.Context, fc func(tx *gorm.DB) error) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Transaction", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Transaction", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Transaction", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Transaction", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Transaction", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Transaction", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Transaction", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Transaction(fc)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w GORMDBWrapper) Unscoped(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Unscoped", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Unscoped", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Unscoped", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Unscoped", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Unscoped", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Unscoped", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Unscoped", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Unscoped()
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Update(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Update", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Update", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Update", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Update", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Update", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Update", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Update", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Update(attrs...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) UpdateColumn(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.UpdateColumn", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.UpdateColumn", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.UpdateColumn", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.UpdateColumn", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.UpdateColumn", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.UpdateColumn", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.UpdateColumn", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.UpdateColumn(attrs...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) UpdateColumns(ctx context.Context, values interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.UpdateColumns", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.UpdateColumns", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.UpdateColumns", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.UpdateColumns", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.UpdateColumns", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.UpdateColumns", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.UpdateColumns", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.UpdateColumns(values)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Updates(ctx context.Context, values interface{}, ignoreProtectedAttrs ...bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Updates", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Updates", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Updates", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Updates", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Updates", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Updates", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Updates", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Updates(values, ignoreProtectedAttrs...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w GORMDBWrapper) Where(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DB.Where", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DB.Where", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DB.Where", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Where", ctxOptions.startSpanOpts...)
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
			w.inflightMetric.WithLabelValues("gorm.DB.Where", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("gorm.DB.Where", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("gorm.DB.Where", ErrCode(res0.Error), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Where(query, args...)
		if res0.Error != nil && span != nil {
			span.SetTag("error", res0.Error.Error())
		}
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}
