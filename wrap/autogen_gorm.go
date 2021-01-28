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
	"github.com/hatlonely/go-kit/refx"
)

type GORMDBWrapper struct {
	obj              *gorm.DB
	retry            *Retry
	options          *WrapperOptions
	durationMetric   *prometheus.HistogramVec
	inflightMetric   *prometheus.GaugeVec
	rateLimiterGroup RateLimiterGroup
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

func (w *GORMDBWrapper) OnRateLimiterGroupChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w GORMDBWrapper) AddError(ctx context.Context, err error) error {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.AddError")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.AddError")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.AddError", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.AddError", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.AddError", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.AddError(err)
	return res0
}

func (w GORMDBWrapper) AddForeignKey(ctx context.Context, field string, dest string, onDelete string, onUpdate string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.AddForeignKey"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.AddForeignKey")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.AddForeignKey", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.AddForeignKey(field, dest, onDelete, onUpdate)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) AddIndex(ctx context.Context, indexName string, columns ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.AddIndex"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.AddIndex")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.AddIndex", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.AddIndex(indexName, columns...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) AddUniqueIndex(ctx context.Context, indexName string, columns ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.AddUniqueIndex"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.AddUniqueIndex")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.AddUniqueIndex", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.AddUniqueIndex(indexName, columns...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Assign(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Assign"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Assign")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Assign", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Assign(attrs...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Association(ctx context.Context, column string) *gorm.Association {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.Association")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Association")
		for key, val := range w.options.Trace.ConstTags {
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
			w.durationMetric.WithLabelValues("gorm.DB.Association", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.Association(column)
	return res0
}

func (w GORMDBWrapper) Attrs(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Attrs"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Attrs")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Attrs", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Attrs(attrs...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) AutoMigrate(ctx context.Context, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.AutoMigrate"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.AutoMigrate")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.AutoMigrate", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.AutoMigrate(values...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Begin(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Begin"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Begin")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Begin", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Begin()
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) BeginTx(ctx context.Context, opts *sql.TxOptions) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.BeginTx"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.BeginTx")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.BeginTx", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BeginTx(ctx, opts)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) BlockGlobalUpdate(ctx context.Context, enable bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.BlockGlobalUpdate"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.BlockGlobalUpdate")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.BlockGlobalUpdate", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BlockGlobalUpdate(enable)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Callback(ctx context.Context) *gorm.Callback {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.Callback")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Callback")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.Callback", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.Callback", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.Callback", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.Callback()
	return res0
}

func (w GORMDBWrapper) Close(ctx context.Context) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Close"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Close")
			for key, val := range w.options.Trace.ConstTags {
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
		return err
	})
	return err
}

func (w GORMDBWrapper) Commit(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Commit"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Commit")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Commit", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Commit()
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) CommonDB(ctx context.Context) gorm.SQLCommon {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.CommonDB")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.CommonDB")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.CommonDB", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.CommonDB", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.CommonDB", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.CommonDB()
	return res0
}

func (w GORMDBWrapper) Count(ctx context.Context, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Count"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Count")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Count", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Count(value)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Create(ctx context.Context, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Create"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Create")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Create", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Create(value)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) CreateTable(ctx context.Context, models ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.CreateTable"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.CreateTable")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.CreateTable", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.CreateTable(models...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) DB(ctx context.Context) *sql.DB {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.DB")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.DB")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.DB", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.DB", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.DB", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.DB()
	return res0
}

func (w GORMDBWrapper) Debug(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Debug"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Debug")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Debug", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Debug()
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Delete(ctx context.Context, value interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Delete"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Delete")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Delete", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Delete(value, where...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Dialect(ctx context.Context) gorm.Dialect {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.Dialect")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Dialect")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.Dialect", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.Dialect", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.Dialect", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.Dialect()
	return res0
}

func (w GORMDBWrapper) DropColumn(ctx context.Context, column string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.DropColumn"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.DropColumn")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.DropColumn", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.DropColumn(column)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) DropTable(ctx context.Context, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.DropTable"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.DropTable")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.DropTable", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.DropTable(values...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) DropTableIfExists(ctx context.Context, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.DropTableIfExists"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.DropTableIfExists")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.DropTableIfExists", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.DropTableIfExists(values...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Exec(ctx context.Context, sql string, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Exec"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Exec")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Exec", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Exec(sql, values...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Find(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Find"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Find")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Find", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Find(out, where...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) First(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.First"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.First")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.First", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.First(out, where...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) FirstOrCreate(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.FirstOrCreate"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.FirstOrCreate")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.FirstOrCreate", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.FirstOrCreate(out, where...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) FirstOrInit(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.FirstOrInit"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.FirstOrInit")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.FirstOrInit", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.FirstOrInit(out, where...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Get(ctx context.Context, name string) (interface{}, bool) {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.Get")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Get")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.Get", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.Get", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.Get", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	value, ok := w.obj.Get(name)
	return value, ok
}

func (w GORMDBWrapper) GetErrors(ctx context.Context) []error {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.GetErrors")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.GetErrors")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.GetErrors", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.GetErrors", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.GetErrors", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.GetErrors()
	return res0
}

func (w GORMDBWrapper) Group(ctx context.Context, query string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Group"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Group")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Group", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Group(query)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) HasBlockGlobalUpdate(ctx context.Context) bool {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.HasBlockGlobalUpdate")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.HasBlockGlobalUpdate")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.HasBlockGlobalUpdate", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.HasBlockGlobalUpdate", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.HasBlockGlobalUpdate", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.HasBlockGlobalUpdate()
	return res0
}

func (w GORMDBWrapper) HasTable(ctx context.Context, value interface{}) bool {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.HasTable")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.HasTable")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.HasTable", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.HasTable", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.HasTable", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.HasTable(value)
	return res0
}

func (w GORMDBWrapper) Having(ctx context.Context, query interface{}, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Having"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Having")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Having", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Having(query, values...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) InstantSet(ctx context.Context, name string, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.InstantSet"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.InstantSet")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.InstantSet", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.InstantSet(name, value)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Joins(ctx context.Context, query string, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Joins"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Joins")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Joins", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Joins(query, args...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Last(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Last"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Last")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Last", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Last(out, where...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Limit(ctx context.Context, limit interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Limit"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Limit")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Limit", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Limit(limit)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) LogMode(ctx context.Context, enable bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.LogMode"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.LogMode")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.LogMode", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LogMode(enable)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Model(ctx context.Context, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Model"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Model")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Model", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Model(value)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) ModifyColumn(ctx context.Context, column string, typ string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.ModifyColumn"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.ModifyColumn")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.ModifyColumn", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ModifyColumn(column, typ)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) New(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.New"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.New")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.New", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.New()
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) NewRecord(ctx context.Context, value interface{}) bool {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.NewRecord")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.NewRecord")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.NewRecord", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.NewRecord", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.NewRecord", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.NewRecord(value)
	return res0
}

func (w GORMDBWrapper) NewScope(ctx context.Context, value interface{}) *gorm.Scope {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.NewScope")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.NewScope")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.NewScope", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.NewScope", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.NewScope", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.NewScope(value)
	return res0
}

func (w GORMDBWrapper) Not(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Not"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Not")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Not", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Not(query, args...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Offset(ctx context.Context, offset interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Offset"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Offset")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Offset", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Offset(offset)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Omit(ctx context.Context, columns ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Omit"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Omit")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Omit", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Omit(columns...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Or(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Or"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Or")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Or", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Or(query, args...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Order(ctx context.Context, value interface{}, reorder ...bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Order"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Order")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Order", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Order(value, reorder...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Pluck(ctx context.Context, column string, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Pluck"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Pluck")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Pluck", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Pluck(column, value)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Preload(ctx context.Context, column string, conditions ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Preload"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Preload")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Preload", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Preload(column, conditions...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Preloads(ctx context.Context, out interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Preloads"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Preloads")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Preloads", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Preloads(out)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) QueryExpr(ctx context.Context) *gorm.SqlExpr {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.QueryExpr")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.QueryExpr")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.QueryExpr", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.QueryExpr", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.QueryExpr", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.QueryExpr()
	return res0
}

func (w GORMDBWrapper) Raw(ctx context.Context, sql string, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Raw"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Raw")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Raw", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Raw(sql, values...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) RecordNotFound(ctx context.Context) bool {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.RecordNotFound")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.RecordNotFound")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.RecordNotFound", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.RecordNotFound", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.RecordNotFound", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.RecordNotFound()
	return res0
}

func (w GORMDBWrapper) Related(ctx context.Context, value interface{}, foreignKeys ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Related"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Related")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Related", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Related(value, foreignKeys...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) RemoveForeignKey(ctx context.Context, field string, dest string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.RemoveForeignKey"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.RemoveForeignKey")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.RemoveForeignKey", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RemoveForeignKey(field, dest)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) RemoveIndex(ctx context.Context, indexName string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.RemoveIndex"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.RemoveIndex")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.RemoveIndex", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RemoveIndex(indexName)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Rollback(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Rollback"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Rollback")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Rollback", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Rollback()
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) RollbackUnlessCommitted(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.RollbackUnlessCommitted"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.RollbackUnlessCommitted")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.RollbackUnlessCommitted", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RollbackUnlessCommitted()
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Row(ctx context.Context) *sql.Row {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.Row")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Row")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.Row", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.Row", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.Row", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.Row()
	return res0
}

func (w GORMDBWrapper) Rows(ctx context.Context) (*sql.Rows, error) {
	ctxOptions := FromContext(ctx)
	var res0 *sql.Rows
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Rows"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Rows")
			for key, val := range w.options.Trace.ConstTags {
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
		return err
	})
	return res0, err
}

func (w GORMDBWrapper) Save(ctx context.Context, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Save"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Save")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Save", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Save(value)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Scan(ctx context.Context, dest interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Scan"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Scan")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Scan", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Scan(dest)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) ScanRows(ctx context.Context, rows *sql.Rows, result interface{}) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.ScanRows"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.ScanRows")
			for key, val := range w.options.Trace.ConstTags {
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
		return err
	})
	return err
}

func (w GORMDBWrapper) Scopes(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Scopes"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Scopes")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Scopes", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Scopes(funcs...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Select(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Select"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Select")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Select", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Select(query, args...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Set(ctx context.Context, name string, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Set"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Set")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Set", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Set(name, value)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) SetJoinTableHandler(ctx context.Context, source interface{}, column string, handler gorm.JoinTableHandlerInterface) {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.SetJoinTableHandler")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.SetJoinTableHandler")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.SetJoinTableHandler", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.SetJoinTableHandler", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.SetJoinTableHandler", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	w.obj.SetJoinTableHandler(source, column, handler)
}

func (w GORMDBWrapper) SetNowFuncOverride(ctx context.Context, nowFuncOverride func() time.Time) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.SetNowFuncOverride"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.SetNowFuncOverride")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.SetNowFuncOverride", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SetNowFuncOverride(nowFuncOverride)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) SingularTable(ctx context.Context, enable bool) {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.SingularTable")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.SingularTable")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.SingularTable", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.SingularTable", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.SingularTable", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	w.obj.SingularTable(enable)
}

func (w GORMDBWrapper) SubQuery(ctx context.Context) *gorm.SqlExpr {
	ctxOptions := FromContext(ctx)
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "DB.SubQuery")
	}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.SubQuery")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("gorm.DB.SubQuery", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("gorm.DB.SubQuery", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("gorm.DB.SubQuery", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
	res0 := w.obj.SubQuery()
	return res0
}

func (w GORMDBWrapper) Table(ctx context.Context, name string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Table"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Table")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Table", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Table(name)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Take(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Take"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Take")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Take", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Take(out, where...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Transaction(ctx context.Context, fc func(tx *gorm.DB) error) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Transaction"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Transaction")
			for key, val := range w.options.Trace.ConstTags {
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
		return err
	})
	return err
}

func (w GORMDBWrapper) Unscoped(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Unscoped"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Unscoped")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Unscoped", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Unscoped()
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Update(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Update"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Update")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Update", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Update(attrs...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) UpdateColumn(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.UpdateColumn"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.UpdateColumn")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.UpdateColumn", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.UpdateColumn(attrs...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) UpdateColumns(ctx context.Context, values interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.UpdateColumns"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.UpdateColumns")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.UpdateColumns", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.UpdateColumns(values)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Updates(ctx context.Context, values interface{}, ignoreProtectedAttrs ...bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Updates"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Updates")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Updates", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Updates(values, ignoreProtectedAttrs...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}

func (w GORMDBWrapper) Where(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)
	var res0 *gorm.DB
	_ = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "DB.Where"); err != nil {
				return err
			}
		}
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Where")
			for key, val := range w.options.Trace.ConstTags {
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
				w.durationMetric.WithLabelValues("gorm.DB.Where", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Where(query, args...)
		return res0.Error
	})
	return GORMDBWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}
}
