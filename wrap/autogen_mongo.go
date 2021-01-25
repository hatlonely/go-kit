// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"

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

type MongoCollectionWrapper struct {
	obj            *mongo.Collection
	retry          *Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	totalMetric    *prometheus.CounterVec
}

type MongoDatabaseWrapper struct {
	obj            *mongo.Database
	retry          *Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	totalMetric    *prometheus.CounterVec
}

func (w *MongoClientWrapper) Unwrap() *mongo.Client {
	return w.obj
}

func (w *MongoCollectionWrapper) Unwrap() *mongo.Collection {
	return w.obj
}

func (w *MongoDatabaseWrapper) Unwrap() *mongo.Database {
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
	}, []string{"method", "errCode"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "mongo_Client_total",
		Help:        "mongo Client request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode"})
}

func (w *MongoClientWrapper) Connect(ctx context.Context) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.Connect")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.Connect(ctx)
		return err
	})
	return err
}

func (w *MongoClientWrapper) Database(ctx context.Context, name string, opts ...*options.DatabaseOptions) *MongoDatabaseWrapper {
	res0 := w.obj.Database(name, opts...)
	return &MongoDatabaseWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric}
}

func (w *MongoClientWrapper) Disconnect(ctx context.Context) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.Disconnect")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.Disconnect(ctx)
		return err
	})
	return err
}

func (w *MongoClientWrapper) ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) ([]string, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.ListDatabaseNames")
		defer span.Finish()
	}

	var res0 []string
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.ListDatabaseNames(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoClientWrapper) ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (mongo.ListDatabasesResult, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.ListDatabases")
		defer span.Finish()
	}

	var res0 mongo.ListDatabasesResult
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.ListDatabases(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoClientWrapper) NumberSessionsInProgress(ctx context.Context) int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.NumberSessionsInProgress")
		defer span.Finish()
	}

	res0 := w.obj.NumberSessionsInProgress()
	return res0
}

func (w *MongoClientWrapper) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.Ping")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.Ping(ctx, rp)
		return err
	})
	return err
}

func (w *MongoClientWrapper) StartSession(ctx context.Context, opts ...*options.SessionOptions) (mongo.Session, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.StartSession")
		defer span.Finish()
	}

	var res0 mongo.Session
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.StartSession(opts...)
		return err
	})
	return res0, err
}

func (w *MongoClientWrapper) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.UseSession")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.UseSession(ctx, fn)
		return err
	})
	return err
}

func (w *MongoClientWrapper) UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(mongo.SessionContext) error) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.UseSessionWithOptions")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.UseSessionWithOptions(ctx, opts, fn)
		return err
	})
	return err
}

func (w *MongoClientWrapper) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Client.Watch")
		defer span.Finish()
	}

	var res0 *mongo.ChangeStream
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.Watch(ctx, pipeline, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Aggregate")
		defer span.Finish()
	}

	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.Aggregate(ctx, pipeline, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.BulkWrite")
		defer span.Finish()
	}

	var res0 *mongo.BulkWriteResult
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.BulkWrite(ctx, models, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Clone(ctx context.Context, opts ...*options.CollectionOptions) (*MongoCollectionWrapper, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Clone")
		defer span.Finish()
	}

	var res0 *mongo.Collection
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.Clone(opts...)
		return err
	})
	return &MongoCollectionWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric}, err
}

func (w *MongoCollectionWrapper) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.CountDocuments")
		defer span.Finish()
	}

	var res0 int64
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.CountDocuments(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Database(ctx context.Context) *MongoDatabaseWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Database")
		defer span.Finish()
	}

	res0 := w.obj.Database()
	return &MongoDatabaseWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric}
}

func (w *MongoCollectionWrapper) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.DeleteMany")
		defer span.Finish()
	}

	var res0 *mongo.DeleteResult
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.DeleteMany(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.DeleteOne")
		defer span.Finish()
	}

	var res0 *mongo.DeleteResult
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.DeleteOne(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Distinct")
		defer span.Finish()
	}

	var res0 []interface{}
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.Distinct(ctx, fieldName, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Drop(ctx context.Context) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Drop")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.Drop(ctx)
		return err
	})
	return err
}

func (w *MongoCollectionWrapper) EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.EstimatedDocumentCount")
		defer span.Finish()
	}

	var res0 int64
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.EstimatedDocumentCount(ctx, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Find")
		defer span.Finish()
	}

	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.Find(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.FindOne")
		defer span.Finish()
	}

	res0 := w.obj.FindOne(ctx, filter, opts...)
	return res0
}

func (w *MongoCollectionWrapper) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.FindOneAndDelete")
		defer span.Finish()
	}

	res0 := w.obj.FindOneAndDelete(ctx, filter, opts...)
	return res0
}

func (w *MongoCollectionWrapper) FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.FindOneAndReplace")
		defer span.Finish()
	}

	res0 := w.obj.FindOneAndReplace(ctx, filter, replacement, opts...)
	return res0
}

func (w *MongoCollectionWrapper) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.FindOneAndUpdate")
		defer span.Finish()
	}

	res0 := w.obj.FindOneAndUpdate(ctx, filter, update, opts...)
	return res0
}

func (w *MongoCollectionWrapper) Indexes(ctx context.Context) mongo.IndexView {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Indexes")
		defer span.Finish()
	}

	res0 := w.obj.Indexes()
	return res0
}

func (w *MongoCollectionWrapper) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.InsertMany")
		defer span.Finish()
	}

	var res0 *mongo.InsertManyResult
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.InsertMany(ctx, documents, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.InsertOne")
		defer span.Finish()
	}

	var res0 *mongo.InsertOneResult
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.InsertOne(ctx, document, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Name(ctx context.Context) string {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Name")
		defer span.Finish()
	}

	res0 := w.obj.Name()
	return res0
}

func (w *MongoCollectionWrapper) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.ReplaceOne")
		defer span.Finish()
	}

	var res0 *mongo.UpdateResult
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.ReplaceOne(ctx, filter, replacement, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.UpdateMany")
		defer span.Finish()
	}

	var res0 *mongo.UpdateResult
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.UpdateMany(ctx, filter, update, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.UpdateOne")
		defer span.Finish()
	}

	var res0 *mongo.UpdateResult
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.UpdateOne(ctx, filter, update, opts...)
		return err
	})
	return res0, err
}

func (w *MongoCollectionWrapper) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Collection.Watch")
		defer span.Finish()
	}

	var res0 *mongo.ChangeStream
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.Watch(ctx, pipeline, opts...)
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.Aggregate")
		defer span.Finish()
	}

	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.Aggregate(ctx, pipeline, opts...)
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) Client(ctx context.Context) *MongoClientWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.Client")
		defer span.Finish()
	}

	res0 := w.obj.Client()
	return &MongoClientWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric}
}

func (w *MongoDatabaseWrapper) Collection(ctx context.Context, name string, opts ...*options.CollectionOptions) *MongoCollectionWrapper {
	res0 := w.obj.Collection(name, opts...)
	return &MongoCollectionWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric}
}

func (w *MongoDatabaseWrapper) CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.CreateCollection")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.CreateCollection(ctx, name, opts...)
		return err
	})
	return err
}

func (w *MongoDatabaseWrapper) CreateView(ctx context.Context, viewName string, viewOn string, pipeline interface{}, opts ...*options.CreateViewOptions) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.CreateView")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.CreateView(ctx, viewName, viewOn, pipeline, opts...)
		return err
	})
	return err
}

func (w *MongoDatabaseWrapper) Drop(ctx context.Context) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.Drop")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.Drop(ctx)
		return err
	})
	return err
}

func (w *MongoDatabaseWrapper) ListCollectionNames(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) ([]string, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.ListCollectionNames")
		defer span.Finish()
	}

	var res0 []string
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.ListCollectionNames(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) ListCollections(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) (*mongo.Cursor, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.ListCollections")
		defer span.Finish()
	}

	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.ListCollections(ctx, filter, opts...)
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) Name(ctx context.Context) string {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.Name")
		defer span.Finish()
	}

	res0 := w.obj.Name()
	return res0
}

func (w *MongoDatabaseWrapper) ReadConcern(ctx context.Context) *readconcern.ReadConcern {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.ReadConcern")
		defer span.Finish()
	}

	res0 := w.obj.ReadConcern()
	return res0
}

func (w *MongoDatabaseWrapper) ReadPreference(ctx context.Context) *readpref.ReadPref {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.ReadPreference")
		defer span.Finish()
	}

	res0 := w.obj.ReadPreference()
	return res0
}

func (w *MongoDatabaseWrapper) RunCommand(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *mongo.SingleResult {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.RunCommand")
		defer span.Finish()
	}

	res0 := w.obj.RunCommand(ctx, runCommand, opts...)
	return res0
}

func (w *MongoDatabaseWrapper) RunCommandCursor(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) (*mongo.Cursor, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.RunCommandCursor")
		defer span.Finish()
	}

	var res0 *mongo.Cursor
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.RunCommandCursor(ctx, runCommand, opts...)
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.Watch")
		defer span.Finish()
	}

	var res0 *mongo.ChangeStream
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.Watch(ctx, pipeline, opts...)
		return err
	})
	return res0, err
}

func (w *MongoDatabaseWrapper) WriteConcern(ctx context.Context) *writeconcern.WriteConcern {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "mongo.Database.WriteConcern")
		defer span.Finish()
	}

	res0 := w.obj.WriteConcern()
	return res0
}
