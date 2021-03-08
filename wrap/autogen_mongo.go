// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

func NewMongoClientWrapper(
	obj *mongo.Client,
	retry *micro.Retry,
	options *WrapperOptions,
	durationMetric *prometheus.HistogramVec,
	inflightMetric *prometheus.GaugeVec,
	rateLimiter micro.RateLimiter,
	parallelController micro.ParallelController) *MongoClientWrapper {
	return &MongoClientWrapper{
		obj:                obj,
		retry:              retry,
		options:            options,
		durationMetric:     durationMetric,
		inflightMetric:     inflightMetric,
		rateLimiter:        rateLimiter,
		parallelController: parallelController,
	}
}

type MongoClientWrapper struct {
	obj                *mongo.Client
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *MongoClientWrapper) Unwrap() *mongo.Client {
	return w.obj
}

func (w *MongoClientWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *MongoClientWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *MongoClientWrapper) OnRateLimiterChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *MongoClientWrapper) OnParallelControllerChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *MongoClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        fmt.Sprintf("%s_mongo_Client_durationMs", options.Name),
		Help:        "mongo Client response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        fmt.Sprintf("%s_mongo_Client_inflight", options.Name),
		Help:        "mongo Client inflight",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "custom"})
}

type MongoCollectionWrapper struct {
	obj                *mongo.Collection
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *MongoCollectionWrapper) Unwrap() *mongo.Collection {
	return w.obj
}

type MongoDatabaseWrapper struct {
	obj                *mongo.Database
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *MongoDatabaseWrapper) Unwrap() *mongo.Database {
	return w.obj
}

func (w *MongoClientWrapper) Connect(ctx context.Context) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.Connect", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.Connect", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.Connect", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Client.Connect")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Client.Connect", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Client.Connect", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Client.Connect", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Connect(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *MongoClientWrapper) Database(name string, opts ...*options.DatabaseOptions) *MongoDatabaseWrapper {
	res0 := w.obj.Database(name, opts...)
	return &MongoDatabaseWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *MongoClientWrapper) Disconnect(ctx context.Context) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.Disconnect", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.Disconnect", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.Disconnect", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Client.Disconnect")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Client.Disconnect", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Client.Disconnect", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Client.Disconnect", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Disconnect(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *MongoClientWrapper) ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) ([]string, error) {
	ctxOptions := FromContext(ctx)
	var res0 []string
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListDatabaseNames", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListDatabaseNames", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListDatabaseNames", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Client.ListDatabaseNames")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Client.ListDatabaseNames", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Client.ListDatabaseNames", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Client.ListDatabaseNames", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListDatabaseNames(ctx, filter, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoClientWrapper) ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (mongo.ListDatabasesResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 mongo.ListDatabasesResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListDatabases", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListDatabases", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListDatabases", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Client.ListDatabases")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Client.ListDatabases", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Client.ListDatabases", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Client.ListDatabases", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListDatabases(ctx, filter, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoClientWrapper) NumberSessionsInProgress() int {
	res0 := w.obj.NumberSessionsInProgress()
	return res0
}

func (w *MongoClientWrapper) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.Ping", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.Ping", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.Ping", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Client.Ping")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Client.Ping", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Client.Ping", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Client.Ping", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Ping(ctx, rp)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *MongoClientWrapper) StartSession(ctx context.Context, opts ...*options.SessionOptions) (mongo.Session, error) {
	ctxOptions := FromContext(ctx)
	var res0 mongo.Session
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.StartSession", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.StartSession", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.StartSession", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Client.StartSession")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Client.StartSession", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Client.StartSession", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Client.StartSession", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.StartSession(opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoClientWrapper) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UseSession", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UseSession", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UseSession", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Client.UseSession")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Client.UseSession", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Client.UseSession", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Client.UseSession", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.UseSession(ctx, fn)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *MongoClientWrapper) UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(mongo.SessionContext) error) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UseSessionWithOptions", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UseSessionWithOptions", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UseSessionWithOptions", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Client.UseSessionWithOptions")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Client.UseSessionWithOptions", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Client.UseSessionWithOptions", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Client.UseSessionWithOptions", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.UseSessionWithOptions(ctx, opts, fn)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *MongoClientWrapper) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.ChangeStream
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.Watch", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.Watch", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.Watch", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Client.Watch")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Client.Watch", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Client.Watch", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Client.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Watch(ctx, pipeline, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.Aggregate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.Aggregate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.Aggregate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.Aggregate")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.Aggregate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.Aggregate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.Aggregate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Aggregate(ctx, pipeline, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.BulkWriteResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.BulkWrite", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.BulkWrite", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.BulkWrite", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.BulkWrite")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.BulkWrite", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.BulkWrite", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.BulkWrite", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.BulkWrite(ctx, models, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Clone(ctx context.Context, opts ...*options.CollectionOptions) (*MongoCollectionWrapper, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.Collection
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.Clone", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.Clone", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.Clone", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.Clone")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.Clone", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.Clone", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.Clone", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Clone(opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return &MongoCollectionWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}, err
}

func (w *MongoCollectionWrapper) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	ctxOptions := FromContext(ctx)
	var res0 int64
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.CountDocuments", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.CountDocuments", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.CountDocuments", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.CountDocuments")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.CountDocuments", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.CountDocuments", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.CountDocuments", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.CountDocuments(ctx, filter, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Database() *MongoDatabaseWrapper {
	res0 := w.obj.Database()
	return &MongoDatabaseWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *MongoCollectionWrapper) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.DeleteResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.DeleteMany", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.DeleteMany", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.DeleteMany", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.DeleteMany")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.DeleteMany", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.DeleteMany", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.DeleteMany", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DeleteMany(ctx, filter, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.DeleteResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.DeleteOne", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.DeleteOne", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.DeleteOne", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.DeleteOne")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.DeleteOne", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.DeleteOne", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.DeleteOne", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DeleteOne(ctx, filter, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	ctxOptions := FromContext(ctx)
	var res0 []interface{}
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.Distinct", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.Distinct", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.Distinct", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.Distinct")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.Distinct", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.Distinct", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.Distinct", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Distinct(ctx, fieldName, filter, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Drop(ctx context.Context) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.Drop", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.Drop", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.Drop", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.Drop")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.Drop", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.Drop", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.Drop", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Drop(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *MongoCollectionWrapper) EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	ctxOptions := FromContext(ctx)
	var res0 int64
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.EstimatedDocumentCount", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.EstimatedDocumentCount", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.EstimatedDocumentCount", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.EstimatedDocumentCount")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.EstimatedDocumentCount", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.EstimatedDocumentCount", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.EstimatedDocumentCount", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.EstimatedDocumentCount(ctx, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.Find", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.Find", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.Find", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.Find")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.Find", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.Find", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.Find", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Find(ctx, filter, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	res0 := w.obj.FindOne(ctx, filter, opts...)
	return res0
}

func (w *MongoCollectionWrapper) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	res0 := w.obj.FindOneAndDelete(ctx, filter, opts...)
	return res0
}

func (w *MongoCollectionWrapper) FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	res0 := w.obj.FindOneAndReplace(ctx, filter, replacement, opts...)
	return res0
}

func (w *MongoCollectionWrapper) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	res0 := w.obj.FindOneAndUpdate(ctx, filter, update, opts...)
	return res0
}

func (w *MongoCollectionWrapper) Indexes() mongo.IndexView {
	res0 := w.obj.Indexes()
	return res0
}

func (w *MongoCollectionWrapper) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.InsertManyResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.InsertMany", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.InsertMany", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.InsertMany", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.InsertMany")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.InsertMany", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.InsertMany", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.InsertMany", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.InsertMany(ctx, documents, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.InsertOneResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.InsertOne", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.InsertOne", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.InsertOne", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.InsertOne")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.InsertOne", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.InsertOne", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.InsertOne", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.InsertOne(ctx, document, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Name() string {
	res0 := w.obj.Name()
	return res0
}

func (w *MongoCollectionWrapper) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.UpdateResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.ReplaceOne", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.ReplaceOne", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.ReplaceOne", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.ReplaceOne")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.ReplaceOne", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.ReplaceOne", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.ReplaceOne", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ReplaceOne(ctx, filter, replacement, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.UpdateResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.UpdateMany", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.UpdateMany", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.UpdateMany", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.UpdateMany")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.UpdateMany", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.UpdateMany", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.UpdateMany", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.UpdateMany(ctx, filter, update, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.UpdateResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.UpdateOne", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.UpdateOne", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.UpdateOne", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.UpdateOne")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.UpdateOne", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.UpdateOne", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.UpdateOne", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.UpdateOne(ctx, filter, update, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.ChangeStream
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Collection.Watch", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Collection.Watch", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Collection.Watch", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Collection.Watch")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Collection.Watch", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Collection.Watch", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Collection.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Watch(ctx, pipeline, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Database.Aggregate", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Database.Aggregate", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Database.Aggregate", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Database.Aggregate")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Database.Aggregate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Database.Aggregate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Database.Aggregate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Aggregate(ctx, pipeline, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) Client() *MongoClientWrapper {
	res0 := w.obj.Client()
	return &MongoClientWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *MongoDatabaseWrapper) Collection(name string, opts ...*options.CollectionOptions) *MongoCollectionWrapper {
	res0 := w.obj.Collection(name, opts...)
	return &MongoCollectionWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *MongoDatabaseWrapper) CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Database.CreateCollection", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Database.CreateCollection", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Database.CreateCollection", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Database.CreateCollection")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Database.CreateCollection", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Database.CreateCollection", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Database.CreateCollection", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.CreateCollection(ctx, name, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *MongoDatabaseWrapper) CreateView(ctx context.Context, viewName string, viewOn string, pipeline interface{}, opts ...*options.CreateViewOptions) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Database.CreateView", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Database.CreateView", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Database.CreateView", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Database.CreateView")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Database.CreateView", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Database.CreateView", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Database.CreateView", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.CreateView(ctx, viewName, viewOn, pipeline, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *MongoDatabaseWrapper) Drop(ctx context.Context) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Database.Drop", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Database.Drop", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Database.Drop", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Database.Drop")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Database.Drop", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Database.Drop", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Database.Drop", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Drop(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *MongoDatabaseWrapper) ListCollectionNames(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) ([]string, error) {
	ctxOptions := FromContext(ctx)
	var res0 []string
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Database.ListCollectionNames", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Database.ListCollectionNames", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Database.ListCollectionNames", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Database.ListCollectionNames")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Database.ListCollectionNames", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Database.ListCollectionNames", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Database.ListCollectionNames", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListCollectionNames(ctx, filter, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) ListCollections(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) (*mongo.Cursor, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Database.ListCollections", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Database.ListCollections", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Database.ListCollections", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Database.ListCollections")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Database.ListCollections", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Database.ListCollections", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Database.ListCollections", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListCollections(ctx, filter, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) Name() string {
	res0 := w.obj.Name()
	return res0
}

func (w *MongoDatabaseWrapper) ReadConcern() *readconcern.ReadConcern {
	res0 := w.obj.ReadConcern()
	return res0
}

func (w *MongoDatabaseWrapper) ReadPreference() *readpref.ReadPref {
	res0 := w.obj.ReadPreference()
	return res0
}

func (w *MongoDatabaseWrapper) RunCommand(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *mongo.SingleResult {
	res0 := w.obj.RunCommand(ctx, runCommand, opts...)
	return res0
}

func (w *MongoDatabaseWrapper) RunCommandCursor(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) (*mongo.Cursor, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Database.RunCommandCursor", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Database.RunCommandCursor", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Database.RunCommandCursor", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Database.RunCommandCursor")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Database.RunCommandCursor", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Database.RunCommandCursor", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Database.RunCommandCursor", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.RunCommandCursor(ctx, runCommand, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	ctxOptions := FromContext(ctx)
	var res0 *mongo.ChangeStream
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Database.Watch", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Database.Watch", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Database.Watch", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "mongo.Database.Watch")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("mongo.Database.Watch", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("mongo.Database.Watch", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("mongo.Database.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Watch(ctx, pipeline, opts...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) WriteConcern() *writeconcern.WriteConcern {
	res0 := w.obj.WriteConcern()
	return res0
}
