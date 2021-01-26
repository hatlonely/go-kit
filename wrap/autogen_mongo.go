// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
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
	"github.com/hatlonely/go-kit/refx"
)

type MongoClientWrapper struct {
	obj            *mongo.Client
	retry          *Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	totalMetric    *prometheus.CounterVec
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

func (w *MongoClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "mongo_Client_durationMs",
		Help:        "mongo Client response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "mongo_Client_total",
		Help:        "mongo Client request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
}

type MongoCollectionWrapper struct {
	obj            *mongo.Collection
	retry          *Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	totalMetric    *prometheus.CounterVec
}

func (w *MongoCollectionWrapper) Unwrap() *mongo.Collection {
	return w.obj
}

type MongoDatabaseWrapper struct {
	obj            *mongo.Database
	retry          *Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	totalMetric    *prometheus.CounterVec
}

func (w *MongoDatabaseWrapper) Unwrap() *mongo.Database {
	return w.obj
}

func (w *MongoClientWrapper) Connect(ctx context.Context) error {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.Connect")
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
				w.totalMetric.WithLabelValues("mongo.Client.Connect", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Client.Connect", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		err = w.obj.Connect(ctx)
		return err
	})
	return err
}

func (w *MongoClientWrapper) Database(ctx context.Context, name string, opts ...*options.DatabaseOptions) *MongoDatabaseWrapper {
	var res0 *mongo.Database
	res0 = w.obj.Database(name, opts...)
	return &MongoDatabaseWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric}
}

func (w *MongoClientWrapper) Disconnect(ctx context.Context) error {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.Disconnect")
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
				w.totalMetric.WithLabelValues("mongo.Client.Disconnect", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Client.Disconnect", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		err = w.obj.Disconnect(ctx)
		return err
	})
	return err
}

func (w *MongoClientWrapper) ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) ([]string, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.ListDatabaseNames")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 []string
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Client.ListDatabaseNames", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Client.ListDatabaseNames", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.ListDatabaseNames(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoClientWrapper) ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (mongo.ListDatabasesResult, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.ListDatabases")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 mongo.ListDatabasesResult
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Client.ListDatabases", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Client.ListDatabases", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.ListDatabases(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoClientWrapper) NumberSessionsInProgress(ctx context.Context) int {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.NumberSessionsInProgress")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 int
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("mongo.Client.NumberSessionsInProgress", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("mongo.Client.NumberSessionsInProgress", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.NumberSessionsInProgress()
	return res0
}

func (w *MongoClientWrapper) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.Ping")
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
				w.totalMetric.WithLabelValues("mongo.Client.Ping", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Client.Ping", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		err = w.obj.Ping(ctx, rp)
		return err
	})
	return err
}

func (w *MongoClientWrapper) StartSession(ctx context.Context, opts ...*options.SessionOptions) (mongo.Session, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.StartSession")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 mongo.Session
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Client.StartSession", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Client.StartSession", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.StartSession(opts...)
		return err
	})
	return res0, err
}

func (w *MongoClientWrapper) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.UseSession")
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
				w.totalMetric.WithLabelValues("mongo.Client.UseSession", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Client.UseSession", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		err = w.obj.UseSession(ctx, fn)
		return err
	})
	return err
}

func (w *MongoClientWrapper) UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(mongo.SessionContext) error) error {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.UseSessionWithOptions")
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
				w.totalMetric.WithLabelValues("mongo.Client.UseSessionWithOptions", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Client.UseSessionWithOptions", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		err = w.obj.UseSessionWithOptions(ctx, opts, fn)
		return err
	})
	return err
}

func (w *MongoClientWrapper) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.Watch")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.ChangeStream
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Client.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Client.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.Watch(ctx, pipeline, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Aggregate")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.Aggregate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.Aggregate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.Aggregate(ctx, pipeline, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.BulkWrite")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.BulkWriteResult
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.BulkWrite", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.BulkWrite", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.BulkWrite(ctx, models, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Clone(ctx context.Context, opts ...*options.CollectionOptions) (*MongoCollectionWrapper, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Clone")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.Collection
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.Clone", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.Clone", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.Clone(opts...)
		return err
	})
	return &MongoCollectionWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric}, err
}

func (w *MongoCollectionWrapper) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.CountDocuments")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 int64
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.CountDocuments", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.CountDocuments", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.CountDocuments(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Database(ctx context.Context) *MongoDatabaseWrapper {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Database")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.Database
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("mongo.Collection.Database", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("mongo.Collection.Database", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.Database()
	return &MongoDatabaseWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric}
}

func (w *MongoCollectionWrapper) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.DeleteMany")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.DeleteResult
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.DeleteMany", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.DeleteMany", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.DeleteMany(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.DeleteOne")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.DeleteResult
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.DeleteOne", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.DeleteOne", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.DeleteOne(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Distinct")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 []interface{}
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.Distinct", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.Distinct", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.Distinct(ctx, fieldName, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Drop(ctx context.Context) error {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Drop")
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
				w.totalMetric.WithLabelValues("mongo.Collection.Drop", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.Drop", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		err = w.obj.Drop(ctx)
		return err
	})
	return err
}

func (w *MongoCollectionWrapper) EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.EstimatedDocumentCount")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 int64
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.EstimatedDocumentCount", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.EstimatedDocumentCount", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.EstimatedDocumentCount(ctx, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Find")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.Find", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.Find", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.Find(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.FindOne")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.SingleResult
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("mongo.Collection.FindOne", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("mongo.Collection.FindOne", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.FindOne(ctx, filter, opts...)
	return res0
}

func (w *MongoCollectionWrapper) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.FindOneAndDelete")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.SingleResult
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("mongo.Collection.FindOneAndDelete", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("mongo.Collection.FindOneAndDelete", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.FindOneAndDelete(ctx, filter, opts...)
	return res0
}

func (w *MongoCollectionWrapper) FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.FindOneAndReplace")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.SingleResult
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("mongo.Collection.FindOneAndReplace", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("mongo.Collection.FindOneAndReplace", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.FindOneAndReplace(ctx, filter, replacement, opts...)
	return res0
}

func (w *MongoCollectionWrapper) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.FindOneAndUpdate")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.SingleResult
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("mongo.Collection.FindOneAndUpdate", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("mongo.Collection.FindOneAndUpdate", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.FindOneAndUpdate(ctx, filter, update, opts...)
	return res0
}

func (w *MongoCollectionWrapper) Indexes(ctx context.Context) mongo.IndexView {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Indexes")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 mongo.IndexView
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("mongo.Collection.Indexes", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("mongo.Collection.Indexes", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.Indexes()
	return res0
}

func (w *MongoCollectionWrapper) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.InsertMany")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.InsertManyResult
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.InsertMany", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.InsertMany", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.InsertMany(ctx, documents, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.InsertOne")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.InsertOneResult
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.InsertOne", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.InsertOne", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.InsertOne(ctx, document, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Name(ctx context.Context) string {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Name")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 string
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("mongo.Collection.Name", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("mongo.Collection.Name", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.Name()
	return res0
}

func (w *MongoCollectionWrapper) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.ReplaceOne")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.UpdateResult
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.ReplaceOne", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.ReplaceOne", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.ReplaceOne(ctx, filter, replacement, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.UpdateMany")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.UpdateResult
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.UpdateMany", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.UpdateMany", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.UpdateMany(ctx, filter, update, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.UpdateOne")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.UpdateResult
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.UpdateOne", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.UpdateOne", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.UpdateOne(ctx, filter, update, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Watch")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.ChangeStream
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Collection.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Collection.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.Watch(ctx, pipeline, opts...)
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.Aggregate")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Database.Aggregate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Database.Aggregate", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.Aggregate(ctx, pipeline, opts...)
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) Client(ctx context.Context) *MongoClientWrapper {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.Client")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.Client
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("mongo.Database.Client", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("mongo.Database.Client", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.Client()
	return &MongoClientWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric}
}

func (w *MongoDatabaseWrapper) Collection(ctx context.Context, name string, opts ...*options.CollectionOptions) *MongoCollectionWrapper {
	var res0 *mongo.Collection
	res0 = w.obj.Collection(name, opts...)
	return &MongoCollectionWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric}
}

func (w *MongoDatabaseWrapper) CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.CreateCollection")
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
				w.totalMetric.WithLabelValues("mongo.Database.CreateCollection", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Database.CreateCollection", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		err = w.obj.CreateCollection(ctx, name, opts...)
		return err
	})
	return err
}

func (w *MongoDatabaseWrapper) CreateView(ctx context.Context, viewName string, viewOn string, pipeline interface{}, opts ...*options.CreateViewOptions) error {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.CreateView")
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
				w.totalMetric.WithLabelValues("mongo.Database.CreateView", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Database.CreateView", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		err = w.obj.CreateView(ctx, viewName, viewOn, pipeline, opts...)
		return err
	})
	return err
}

func (w *MongoDatabaseWrapper) Drop(ctx context.Context) error {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.Drop")
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
				w.totalMetric.WithLabelValues("mongo.Database.Drop", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Database.Drop", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		err = w.obj.Drop(ctx)
		return err
	})
	return err
}

func (w *MongoDatabaseWrapper) ListCollectionNames(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) ([]string, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.ListCollectionNames")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 []string
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Database.ListCollectionNames", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Database.ListCollectionNames", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.ListCollectionNames(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) ListCollections(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) (*mongo.Cursor, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.ListCollections")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Database.ListCollections", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Database.ListCollections", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.ListCollections(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) Name(ctx context.Context) string {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.Name")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 string
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("mongo.Database.Name", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("mongo.Database.Name", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.Name()
	return res0
}

func (w *MongoDatabaseWrapper) ReadConcern(ctx context.Context) *readconcern.ReadConcern {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.ReadConcern")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *readconcern.ReadConcern
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("mongo.Database.ReadConcern", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("mongo.Database.ReadConcern", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.ReadConcern()
	return res0
}

func (w *MongoDatabaseWrapper) ReadPreference(ctx context.Context) *readpref.ReadPref {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.ReadPreference")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *readpref.ReadPref
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("mongo.Database.ReadPreference", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("mongo.Database.ReadPreference", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.ReadPreference()
	return res0
}

func (w *MongoDatabaseWrapper) RunCommand(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *mongo.SingleResult {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.RunCommand")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.SingleResult
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("mongo.Database.RunCommand", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("mongo.Database.RunCommand", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.RunCommand(ctx, runCommand, opts...)
	return res0
}

func (w *MongoDatabaseWrapper) RunCommandCursor(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) (*mongo.Cursor, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.RunCommandCursor")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Database.RunCommandCursor", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Database.RunCommandCursor", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.RunCommandCursor(ctx, runCommand, opts...)
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.Watch")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *mongo.ChangeStream
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("mongo.Database.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("mongo.Database.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.Watch(ctx, pipeline, opts...)
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) WriteConcern(ctx context.Context) *writeconcern.WriteConcern {
	ctxOptions := FromContext(ctx)

	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.WriteConcern")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}

	var res0 *writeconcern.WriteConcern
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("mongo.Database.WriteConcern", "OK", ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("mongo.Database.WriteConcern", "OK", ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}

	res0 = w.obj.WriteConcern()
	return res0
}
