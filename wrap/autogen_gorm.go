// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"database/sql"
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
	obj            *gorm.Association
	retry          *micro.Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	inflightMetric *prometheus.GaugeVec
	rateLimiter    micro.RateLimiter
}

func (w GORMAssociationWrapper) Unwrap() *gorm.Association {
	return w.obj
}

type GORMDBWrapper struct {
	obj            *gorm.DB
	retry          *micro.Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	inflightMetric *prometheus.GaugeVec
	rateLimiter    micro.RateLimiter
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

func (w *GORMDBWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "gorm_DB_durationMs",
		Help:        "gorm DB response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "gorm_DB_inflight",
		Help:        "gorm DB inflight",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "custom"})
}

func (w GORMAssociationWrapper) Append(ctx context.Context, values ...interface{}) GORMAssociationWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.Association
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Association.Append"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.Association.Append")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMAssociationWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMAssociationWrapper) Clear(ctx context.Context) GORMAssociationWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.Association
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Association.Clear"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.Association.Clear")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMAssociationWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
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
			if err := w.rateLimiter.Wait(ctx, "Association.Delete"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.Association.Delete")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMAssociationWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMAssociationWrapper) Find(ctx context.Context, value interface{}) GORMAssociationWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.Association
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Association.Find"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.Association.Find")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMAssociationWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMAssociationWrapper) Replace(ctx context.Context, values ...interface{}) GORMAssociationWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.Association
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Association.Replace"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.Association.Replace")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMAssociationWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
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
			if err := w.rateLimiter.Wait(ctx, "DB.AddForeignKey"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.AddForeignKey")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) AddIndex(ctx context.Context, indexName string, columns ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.AddIndex"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.AddIndex")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) AddUniqueIndex(ctx context.Context, indexName string, columns ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.AddUniqueIndex"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.AddUniqueIndex")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Assign(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Assign"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Assign")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Association(ctx context.Context, column string) GORMAssociationWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.Association
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Association"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Association")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMAssociationWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Attrs(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Attrs"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Attrs")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) AutoMigrate(ctx context.Context, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.AutoMigrate"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.AutoMigrate")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Begin(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Begin"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Begin")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) BeginTx(ctx context.Context, opts *sql.TxOptions) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.BeginTx"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.BeginTx")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) BlockGlobalUpdate(ctx context.Context, enable bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.BlockGlobalUpdate"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.BlockGlobalUpdate")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
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
			if err := w.rateLimiter.Wait(ctx, "DB.Close"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Close")
			for key, val := range w.options.Trace.ConstTags {
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
			if err := w.rateLimiter.Wait(ctx, "DB.Commit"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Commit")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
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
			if err := w.rateLimiter.Wait(ctx, "DB.Count"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Count")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Create(ctx context.Context, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Create"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Create")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) CreateTable(ctx context.Context, models ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.CreateTable"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.CreateTable")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
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
			if err := w.rateLimiter.Wait(ctx, "DB.Debug"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Debug")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Delete(ctx context.Context, value interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Delete"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Delete")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
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
			if err := w.rateLimiter.Wait(ctx, "DB.DropColumn"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.DropColumn")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) DropTable(ctx context.Context, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.DropTable"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.DropTable")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) DropTableIfExists(ctx context.Context, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.DropTableIfExists"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.DropTableIfExists")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Exec(ctx context.Context, sql string, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Exec"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Exec")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Find(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Find"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Find")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) First(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.First"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.First")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) FirstOrCreate(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.FirstOrCreate"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.FirstOrCreate")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) FirstOrInit(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.FirstOrInit"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.FirstOrInit")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
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
			if err := w.rateLimiter.Wait(ctx, "DB.Group"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Group")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
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
			if err := w.rateLimiter.Wait(ctx, "DB.Having"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Having")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) InstantSet(ctx context.Context, name string, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.InstantSet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.InstantSet")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Joins(ctx context.Context, query string, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Joins"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Joins")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Last(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Last"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Last")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Limit(ctx context.Context, limit interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Limit"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Limit")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) LogMode(ctx context.Context, enable bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.LogMode"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.LogMode")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Model(ctx context.Context, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Model"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Model")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) ModifyColumn(ctx context.Context, column string, typ string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.ModifyColumn"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.ModifyColumn")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) New(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.New"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.New")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
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
			if err := w.rateLimiter.Wait(ctx, "DB.Not"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Not")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Offset(ctx context.Context, offset interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Offset"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Offset")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Omit(ctx context.Context, columns ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Omit"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Omit")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Or(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Or"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Or")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Order(ctx context.Context, value interface{}, reorder ...bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Order"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Order")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Pluck(ctx context.Context, column string, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Pluck"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Pluck")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Preload(ctx context.Context, column string, conditions ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Preload"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Preload")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Preloads(ctx context.Context, out interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Preloads"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Preloads")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
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
			if err := w.rateLimiter.Wait(ctx, "DB.Raw"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Raw")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
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
			if err := w.rateLimiter.Wait(ctx, "DB.Related"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Related")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) RemoveForeignKey(ctx context.Context, field string, dest string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.RemoveForeignKey"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.RemoveForeignKey")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) RemoveIndex(ctx context.Context, indexName string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.RemoveIndex"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.RemoveIndex")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Rollback(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Rollback"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Rollback")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) RollbackUnlessCommitted(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.RollbackUnlessCommitted"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.RollbackUnlessCommitted")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
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
			if err := w.rateLimiter.Wait(ctx, "DB.Rows"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Rows")
			for key, val := range w.options.Trace.ConstTags {
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
			if err := w.rateLimiter.Wait(ctx, "DB.Save"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Save")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Scan(ctx context.Context, dest interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Scan"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Scan")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) ScanRows(ctx context.Context, rows *sql.Rows, result interface{}) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.ScanRows"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.ScanRows")
			for key, val := range w.options.Trace.ConstTags {
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
			if err := w.rateLimiter.Wait(ctx, "DB.Scopes"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Scopes")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Select(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Select"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Select")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Set(ctx context.Context, name string, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Set"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Set")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) SetJoinTableHandler(source interface{}, column string, handler gorm.JoinTableHandlerInterface) {
	w.obj.SetJoinTableHandler(source, column, handler)
}

func (w GORMDBWrapper) SetNowFuncOverride(ctx context.Context, nowFuncOverride func() time.Time) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.SetNowFuncOverride"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.SetNowFuncOverride")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
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
			if err := w.rateLimiter.Wait(ctx, "DB.Table"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Table")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Take(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Take"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Take")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Transaction(ctx context.Context, fc func(tx *gorm.DB) error) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Transaction"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Transaction")
			for key, val := range w.options.Trace.ConstTags {
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
			if err := w.rateLimiter.Wait(ctx, "DB.Unscoped"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Unscoped")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Update(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Update"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Update")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) UpdateColumn(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.UpdateColumn"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.UpdateColumn")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) UpdateColumns(ctx context.Context, values interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.UpdateColumns"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.UpdateColumns")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Updates(ctx context.Context, values interface{}, ignoreProtectedAttrs ...bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Updates"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Updates")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}

func (w GORMDBWrapper) Where(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "DB.Where"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "gorm.DB.Where")
			for key, val := range w.options.Trace.ConstTags {
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
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter}
}
