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
	obj            *gorm.DB
	retry          *Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	totalMetric    *prometheus.CounterVec
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

func (w *GORMDBWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "gorm_DB_durationMs",
		Help:        "gorm DB response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "gorm_DB_total",
		Help:        "gorm DB request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
}

func (w GORMDBWrapper) AddError(ctx context.Context, err error) error {
	ctxOptions := FromContext(ctx)

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

	var res0 error
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.AddError", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.AddError", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.AddError(err)
	return res0
}

func (w GORMDBWrapper) AddForeignKey(ctx context.Context, field string, dest string, onDelete string, onUpdate string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.AddForeignKey", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.AddForeignKey", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.AddForeignKey(field, dest, onDelete, onUpdate)
	return w
}

func (w GORMDBWrapper) AddIndex(ctx context.Context, indexName string, columns ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.AddIndex", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.AddIndex", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.AddIndex(indexName, columns...)
	return w
}

func (w GORMDBWrapper) AddUniqueIndex(ctx context.Context, indexName string, columns ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.AddUniqueIndex", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.AddUniqueIndex", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.AddUniqueIndex(indexName, columns...)
	return w
}

func (w GORMDBWrapper) Assign(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Assign", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Assign", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Assign(attrs...)
	return w
}

func (w GORMDBWrapper) Association(ctx context.Context, column string) *gorm.Association {
	ctxOptions := FromContext(ctx)

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

	var res0 *gorm.Association
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Association", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Association", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.Association(column)
	return res0
}

func (w GORMDBWrapper) Attrs(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Attrs", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Attrs", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Attrs(attrs...)
	return w
}

func (w GORMDBWrapper) AutoMigrate(ctx context.Context, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.AutoMigrate", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.AutoMigrate", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.AutoMigrate(values...)
	return w
}

func (w GORMDBWrapper) Begin(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Begin", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Begin", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Begin()
	return w
}

func (w GORMDBWrapper) BeginTx(ctx context.Context, opts *sql.TxOptions) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.BeginTx", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.BeginTx", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.BeginTx(ctx, opts)
	return w
}

func (w GORMDBWrapper) BlockGlobalUpdate(ctx context.Context, enable bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.BlockGlobalUpdate", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.BlockGlobalUpdate", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.BlockGlobalUpdate(enable)
	return w
}

func (w GORMDBWrapper) Callback(ctx context.Context) *gorm.Callback {
	ctxOptions := FromContext(ctx)

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

	var res0 *gorm.Callback
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Callback", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Callback", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.Callback()
	return res0
}

func (w GORMDBWrapper) Close(ctx context.Context) error {
	ctxOptions := FromContext(ctx)

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

	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("gorm.DB.Close", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Commit", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Commit", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Commit()
	return w
}

func (w GORMDBWrapper) CommonDB(ctx context.Context) gorm.SQLCommon {
	ctxOptions := FromContext(ctx)

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

	var res0 gorm.SQLCommon
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.CommonDB", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.CommonDB", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.CommonDB()
	return res0
}

func (w GORMDBWrapper) Count(ctx context.Context, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Count", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Count", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Count(value)
	return w
}

func (w GORMDBWrapper) Create(ctx context.Context, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Create", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Create", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Create(value)
	return w
}

func (w GORMDBWrapper) CreateTable(ctx context.Context, models ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.CreateTable", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.CreateTable", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.CreateTable(models...)
	return w
}

func (w GORMDBWrapper) DB(ctx context.Context) *sql.DB {
	ctxOptions := FromContext(ctx)

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

	var res0 *sql.DB
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.DB", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.DB", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.DB()
	return res0
}

func (w GORMDBWrapper) Debug(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Debug", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Debug", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Debug()
	return w
}

func (w GORMDBWrapper) Delete(ctx context.Context, value interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Delete", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Delete", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Delete(value, where...)
	return w
}

func (w GORMDBWrapper) Dialect(ctx context.Context) gorm.Dialect {
	ctxOptions := FromContext(ctx)

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

	var res0 gorm.Dialect
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Dialect", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Dialect", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.Dialect()
	return res0
}

func (w GORMDBWrapper) DropColumn(ctx context.Context, column string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.DropColumn", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.DropColumn", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.DropColumn(column)
	return w
}

func (w GORMDBWrapper) DropTable(ctx context.Context, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.DropTable", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.DropTable", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.DropTable(values...)
	return w
}

func (w GORMDBWrapper) DropTableIfExists(ctx context.Context, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.DropTableIfExists", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.DropTableIfExists", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.DropTableIfExists(values...)
	return w
}

func (w GORMDBWrapper) Exec(ctx context.Context, sql string, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Exec", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Exec", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Exec(sql, values...)
	return w
}

func (w GORMDBWrapper) Find(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Find", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Find", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Find(out, where...)
	return w
}

func (w GORMDBWrapper) First(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.First", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.First", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.First(out, where...)
	return w
}

func (w GORMDBWrapper) FirstOrCreate(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.FirstOrCreate", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.FirstOrCreate", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.FirstOrCreate(out, where...)
	return w
}

func (w GORMDBWrapper) FirstOrInit(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.FirstOrInit", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.FirstOrInit", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.FirstOrInit(out, where...)
	return w
}

func (w GORMDBWrapper) Get(ctx context.Context, name string) (interface{}, bool) {
	ctxOptions := FromContext(ctx)

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

	var value interface{}
	var ok bool
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Get", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Get", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	value, ok = w.obj.Get(name)
	return value, ok
}

func (w GORMDBWrapper) GetErrors(ctx context.Context) []error {
	ctxOptions := FromContext(ctx)

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

	var res0 []error
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.GetErrors", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.GetErrors", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.GetErrors()
	return res0
}

func (w GORMDBWrapper) Group(ctx context.Context, query string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Group", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Group", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Group(query)
	return w
}

func (w GORMDBWrapper) HasBlockGlobalUpdate(ctx context.Context) bool {
	ctxOptions := FromContext(ctx)

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

	var res0 bool
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.HasBlockGlobalUpdate", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.HasBlockGlobalUpdate", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.HasBlockGlobalUpdate()
	return res0
}

func (w GORMDBWrapper) HasTable(ctx context.Context, value interface{}) bool {
	ctxOptions := FromContext(ctx)

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

	var res0 bool
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.HasTable", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.HasTable", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.HasTable(value)
	return res0
}

func (w GORMDBWrapper) Having(ctx context.Context, query interface{}, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Having", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Having", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Having(query, values...)
	return w
}

func (w GORMDBWrapper) InstantSet(ctx context.Context, name string, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.InstantSet", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.InstantSet", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.InstantSet(name, value)
	return w
}

func (w GORMDBWrapper) Joins(ctx context.Context, query string, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Joins", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Joins", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Joins(query, args...)
	return w
}

func (w GORMDBWrapper) Last(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Last", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Last", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Last(out, where...)
	return w
}

func (w GORMDBWrapper) Limit(ctx context.Context, limit interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Limit", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Limit", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Limit(limit)
	return w
}

func (w GORMDBWrapper) LogMode(ctx context.Context, enable bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.LogMode", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.LogMode", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.LogMode(enable)
	return w
}

func (w GORMDBWrapper) Model(ctx context.Context, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Model", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Model", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Model(value)
	return w
}

func (w GORMDBWrapper) ModifyColumn(ctx context.Context, column string, typ string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.ModifyColumn", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.ModifyColumn", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.ModifyColumn(column, typ)
	return w
}

func (w GORMDBWrapper) New(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.New", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.New", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.New()
	return w
}

func (w GORMDBWrapper) NewRecord(ctx context.Context, value interface{}) bool {
	ctxOptions := FromContext(ctx)

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

	var res0 bool
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.NewRecord", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.NewRecord", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.NewRecord(value)
	return res0
}

func (w GORMDBWrapper) NewScope(ctx context.Context, value interface{}) *gorm.Scope {
	ctxOptions := FromContext(ctx)

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

	var res0 *gorm.Scope
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.NewScope", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.NewScope", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.NewScope(value)
	return res0
}

func (w GORMDBWrapper) Not(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Not", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Not", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Not(query, args...)
	return w
}

func (w GORMDBWrapper) Offset(ctx context.Context, offset interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Offset", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Offset", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Offset(offset)
	return w
}

func (w GORMDBWrapper) Omit(ctx context.Context, columns ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Omit", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Omit", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Omit(columns...)
	return w
}

func (w GORMDBWrapper) Or(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Or", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Or", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Or(query, args...)
	return w
}

func (w GORMDBWrapper) Order(ctx context.Context, value interface{}, reorder ...bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Order", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Order", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Order(value, reorder...)
	return w
}

func (w GORMDBWrapper) Pluck(ctx context.Context, column string, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Pluck", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Pluck", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Pluck(column, value)
	return w
}

func (w GORMDBWrapper) Preload(ctx context.Context, column string, conditions ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Preload", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Preload", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Preload(column, conditions...)
	return w
}

func (w GORMDBWrapper) Preloads(ctx context.Context, out interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Preloads", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Preloads", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Preloads(out)
	return w
}

func (w GORMDBWrapper) QueryExpr(ctx context.Context) *gorm.SqlExpr {
	ctxOptions := FromContext(ctx)

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

	var res0 *gorm.SqlExpr
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.QueryExpr", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.QueryExpr", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.QueryExpr()
	return res0
}

func (w GORMDBWrapper) Raw(ctx context.Context, sql string, values ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Raw", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Raw", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Raw(sql, values...)
	return w
}

func (w GORMDBWrapper) RecordNotFound(ctx context.Context) bool {
	ctxOptions := FromContext(ctx)

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

	var res0 bool
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.RecordNotFound", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.RecordNotFound", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.RecordNotFound()
	return res0
}

func (w GORMDBWrapper) Related(ctx context.Context, value interface{}, foreignKeys ...string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Related", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Related", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Related(value, foreignKeys...)
	return w
}

func (w GORMDBWrapper) RemoveForeignKey(ctx context.Context, field string, dest string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.RemoveForeignKey", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.RemoveForeignKey", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.RemoveForeignKey(field, dest)
	return w
}

func (w GORMDBWrapper) RemoveIndex(ctx context.Context, indexName string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.RemoveIndex", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.RemoveIndex", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.RemoveIndex(indexName)
	return w
}

func (w GORMDBWrapper) Rollback(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Rollback", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Rollback", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Rollback()
	return w
}

func (w GORMDBWrapper) RollbackUnlessCommitted(ctx context.Context) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.RollbackUnlessCommitted", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.RollbackUnlessCommitted", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.RollbackUnlessCommitted()
	return w
}

func (w GORMDBWrapper) Row(ctx context.Context) *sql.Row {
	ctxOptions := FromContext(ctx)

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

	var res0 *sql.Row
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Row", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Row", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.Row()
	return res0
}

func (w GORMDBWrapper) Rows(ctx context.Context) (*sql.Rows, error) {
	ctxOptions := FromContext(ctx)

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

	var res0 *sql.Rows
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("gorm.DB.Rows", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Save", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Save", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Save(value)
	return w
}

func (w GORMDBWrapper) Scan(ctx context.Context, dest interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Scan", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Scan", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Scan(dest)
	return w
}

func (w GORMDBWrapper) ScanRows(ctx context.Context, rows *sql.Rows, result interface{}) error {
	ctxOptions := FromContext(ctx)

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

	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("gorm.DB.ScanRows", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Scopes", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Scopes", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Scopes(funcs...)
	return w
}

func (w GORMDBWrapper) Select(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Select", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Select", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Select(query, args...)
	return w
}

func (w GORMDBWrapper) Set(ctx context.Context, name string, value interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Set", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Set", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Set(name, value)
	return w
}

func (w GORMDBWrapper) SetJoinTableHandler(ctx context.Context, source interface{}, column string, handler gorm.JoinTableHandlerInterface) {
	ctxOptions := FromContext(ctx)

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

	w.obj.SetJoinTableHandler(source, column, handler)
}

func (w GORMDBWrapper) SetNowFuncOverride(ctx context.Context, nowFuncOverride func() time.Time) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.SetNowFuncOverride", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.SetNowFuncOverride", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.SetNowFuncOverride(nowFuncOverride)
	return w
}

func (w GORMDBWrapper) SingularTable(ctx context.Context, enable bool) {
	ctxOptions := FromContext(ctx)

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

	w.obj.SingularTable(enable)
}

func (w GORMDBWrapper) SubQuery(ctx context.Context) *gorm.SqlExpr {
	ctxOptions := FromContext(ctx)

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

	var res0 *gorm.SqlExpr
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.SubQuery", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.SubQuery", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.SubQuery()
	return res0
}

func (w GORMDBWrapper) Table(ctx context.Context, name string) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Table", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Table", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Table(name)
	return w
}

func (w GORMDBWrapper) Take(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Take", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Take", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Take(out, where...)
	return w
}

func (w GORMDBWrapper) Transaction(ctx context.Context, fc func(tx *gorm.DB) error) error {
	ctxOptions := FromContext(ctx)

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

	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("gorm.DB.Transaction", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Unscoped", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Unscoped", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Unscoped()
	return w
}

func (w GORMDBWrapper) Update(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Update", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Update", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Update(attrs...)
	return w
}

func (w GORMDBWrapper) UpdateColumn(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.UpdateColumn", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.UpdateColumn", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.UpdateColumn(attrs...)
	return w
}

func (w GORMDBWrapper) UpdateColumns(ctx context.Context, values interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.UpdateColumns", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.UpdateColumns", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.UpdateColumns(values)
	return w
}

func (w GORMDBWrapper) Updates(ctx context.Context, values interface{}, ignoreProtectedAttrs ...bool) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Updates", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Updates", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Updates(values, ignoreProtectedAttrs...)
	return w
}

func (w GORMDBWrapper) Where(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	ctxOptions := FromContext(ctx)

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
		defer func() {
			w.totalMetric.WithLabelValues("gorm.DB.Where", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("gorm.DB.Where", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	w.obj = w.obj.Where(query, args...)
	return w
}
