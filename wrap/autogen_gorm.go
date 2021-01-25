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

func (w *GORMDBWrapper) NewMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "gorm_DB_durationMs",
		Help:        "gorm DB response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "gorm_DB_total",
		Help:        "gorm DB request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode"})
}

func (w GORMDBWrapper) AddError(ctx context.Context, err error) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.AddError")
		defer span.Finish()
	}

	res0 := w.obj.AddError(err)
	return res0
}

func (w GORMDBWrapper) AddForeignKey(ctx context.Context, field string, dest string, onDelete string, onUpdate string) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.AddForeignKey")
		defer span.Finish()
	}

	w.obj = w.obj.AddForeignKey(field, dest, onDelete, onUpdate)
	return w
}

func (w GORMDBWrapper) AddIndex(ctx context.Context, indexName string, columns ...string) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.AddIndex")
		defer span.Finish()
	}

	w.obj = w.obj.AddIndex(indexName, columns...)
	return w
}

func (w GORMDBWrapper) AddUniqueIndex(ctx context.Context, indexName string, columns ...string) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.AddUniqueIndex")
		defer span.Finish()
	}

	w.obj = w.obj.AddUniqueIndex(indexName, columns...)
	return w
}

func (w GORMDBWrapper) Assign(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Assign")
		defer span.Finish()
	}

	w.obj = w.obj.Assign(attrs...)
	return w
}

func (w GORMDBWrapper) Association(ctx context.Context, column string) *gorm.Association {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Association")
		defer span.Finish()
	}

	res0 := w.obj.Association(column)
	return res0
}

func (w GORMDBWrapper) Attrs(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Attrs")
		defer span.Finish()
	}

	w.obj = w.obj.Attrs(attrs...)
	return w
}

func (w GORMDBWrapper) AutoMigrate(ctx context.Context, values ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.AutoMigrate")
		defer span.Finish()
	}

	w.obj = w.obj.AutoMigrate(values...)
	return w
}

func (w GORMDBWrapper) Begin(ctx context.Context) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Begin")
		defer span.Finish()
	}

	w.obj = w.obj.Begin()
	return w
}

func (w GORMDBWrapper) BeginTx(ctx context.Context, opts *sql.TxOptions) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.BeginTx")
		defer span.Finish()
	}

	w.obj = w.obj.BeginTx(ctx, opts)
	return w
}

func (w GORMDBWrapper) BlockGlobalUpdate(ctx context.Context, enable bool) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.BlockGlobalUpdate")
		defer span.Finish()
	}

	w.obj = w.obj.BlockGlobalUpdate(enable)
	return w
}

func (w GORMDBWrapper) Callback(ctx context.Context) *gorm.Callback {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Callback")
		defer span.Finish()
	}

	res0 := w.obj.Callback()
	return res0
}

func (w GORMDBWrapper) Close(ctx context.Context) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Close")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.Close()
		return err
	})
	return err
}

func (w GORMDBWrapper) Commit(ctx context.Context) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Commit")
		defer span.Finish()
	}

	w.obj = w.obj.Commit()
	return w
}

func (w GORMDBWrapper) CommonDB(ctx context.Context) gorm.SQLCommon {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.CommonDB")
		defer span.Finish()
	}

	res0 := w.obj.CommonDB()
	return res0
}

func (w GORMDBWrapper) Count(ctx context.Context, value interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Count")
		defer span.Finish()
	}

	w.obj = w.obj.Count(value)
	return w
}

func (w GORMDBWrapper) Create(ctx context.Context, value interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Create")
		defer span.Finish()
	}

	w.obj = w.obj.Create(value)
	return w
}

func (w GORMDBWrapper) CreateTable(ctx context.Context, models ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.CreateTable")
		defer span.Finish()
	}

	w.obj = w.obj.CreateTable(models...)
	return w
}

func (w GORMDBWrapper) DB(ctx context.Context) *sql.DB {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.DB")
		defer span.Finish()
	}

	res0 := w.obj.DB()
	return res0
}

func (w GORMDBWrapper) Debug(ctx context.Context) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Debug")
		defer span.Finish()
	}

	w.obj = w.obj.Debug()
	return w
}

func (w GORMDBWrapper) Delete(ctx context.Context, value interface{}, where ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Delete")
		defer span.Finish()
	}

	w.obj = w.obj.Delete(value, where...)
	return w
}

func (w GORMDBWrapper) Dialect(ctx context.Context) gorm.Dialect {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Dialect")
		defer span.Finish()
	}

	res0 := w.obj.Dialect()
	return res0
}

func (w GORMDBWrapper) DropColumn(ctx context.Context, column string) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.DropColumn")
		defer span.Finish()
	}

	w.obj = w.obj.DropColumn(column)
	return w
}

func (w GORMDBWrapper) DropTable(ctx context.Context, values ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.DropTable")
		defer span.Finish()
	}

	w.obj = w.obj.DropTable(values...)
	return w
}

func (w GORMDBWrapper) DropTableIfExists(ctx context.Context, values ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.DropTableIfExists")
		defer span.Finish()
	}

	w.obj = w.obj.DropTableIfExists(values...)
	return w
}

func (w GORMDBWrapper) Exec(ctx context.Context, sql string, values ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Exec")
		defer span.Finish()
	}

	w.obj = w.obj.Exec(sql, values...)
	return w
}

func (w GORMDBWrapper) Find(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Find")
		defer span.Finish()
	}

	w.obj = w.obj.Find(out, where...)
	return w
}

func (w GORMDBWrapper) First(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.First")
		defer span.Finish()
	}

	w.obj = w.obj.First(out, where...)
	return w
}

func (w GORMDBWrapper) FirstOrCreate(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.FirstOrCreate")
		defer span.Finish()
	}

	w.obj = w.obj.FirstOrCreate(out, where...)
	return w
}

func (w GORMDBWrapper) FirstOrInit(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.FirstOrInit")
		defer span.Finish()
	}

	w.obj = w.obj.FirstOrInit(out, where...)
	return w
}

func (w GORMDBWrapper) Get(ctx context.Context, name string) (interface{}, bool) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Get")
		defer span.Finish()
	}

	value, ok := w.obj.Get(name)
	return value, ok
}

func (w GORMDBWrapper) GetErrors(ctx context.Context) []error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.GetErrors")
		defer span.Finish()
	}

	res0 := w.obj.GetErrors()
	return res0
}

func (w GORMDBWrapper) Group(ctx context.Context, query string) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Group")
		defer span.Finish()
	}

	w.obj = w.obj.Group(query)
	return w
}

func (w GORMDBWrapper) HasBlockGlobalUpdate(ctx context.Context) bool {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.HasBlockGlobalUpdate")
		defer span.Finish()
	}

	res0 := w.obj.HasBlockGlobalUpdate()
	return res0
}

func (w GORMDBWrapper) HasTable(ctx context.Context, value interface{}) bool {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.HasTable")
		defer span.Finish()
	}

	res0 := w.obj.HasTable(value)
	return res0
}

func (w GORMDBWrapper) Having(ctx context.Context, query interface{}, values ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Having")
		defer span.Finish()
	}

	w.obj = w.obj.Having(query, values...)
	return w
}

func (w GORMDBWrapper) InstantSet(ctx context.Context, name string, value interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.InstantSet")
		defer span.Finish()
	}

	w.obj = w.obj.InstantSet(name, value)
	return w
}

func (w GORMDBWrapper) Joins(ctx context.Context, query string, args ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Joins")
		defer span.Finish()
	}

	w.obj = w.obj.Joins(query, args...)
	return w
}

func (w GORMDBWrapper) Last(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Last")
		defer span.Finish()
	}

	w.obj = w.obj.Last(out, where...)
	return w
}

func (w GORMDBWrapper) Limit(ctx context.Context, limit interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Limit")
		defer span.Finish()
	}

	w.obj = w.obj.Limit(limit)
	return w
}

func (w GORMDBWrapper) LogMode(ctx context.Context, enable bool) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.LogMode")
		defer span.Finish()
	}

	w.obj = w.obj.LogMode(enable)
	return w
}

func (w GORMDBWrapper) Model(ctx context.Context, value interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Model")
		defer span.Finish()
	}

	w.obj = w.obj.Model(value)
	return w
}

func (w GORMDBWrapper) ModifyColumn(ctx context.Context, column string, typ string) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.ModifyColumn")
		defer span.Finish()
	}

	w.obj = w.obj.ModifyColumn(column, typ)
	return w
}

func (w GORMDBWrapper) New(ctx context.Context) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.New")
		defer span.Finish()
	}

	w.obj = w.obj.New()
	return w
}

func (w GORMDBWrapper) NewRecord(ctx context.Context, value interface{}) bool {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.NewRecord")
		defer span.Finish()
	}

	res0 := w.obj.NewRecord(value)
	return res0
}

func (w GORMDBWrapper) NewScope(ctx context.Context, value interface{}) *gorm.Scope {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.NewScope")
		defer span.Finish()
	}

	res0 := w.obj.NewScope(value)
	return res0
}

func (w GORMDBWrapper) Not(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Not")
		defer span.Finish()
	}

	w.obj = w.obj.Not(query, args...)
	return w
}

func (w GORMDBWrapper) Offset(ctx context.Context, offset interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Offset")
		defer span.Finish()
	}

	w.obj = w.obj.Offset(offset)
	return w
}

func (w GORMDBWrapper) Omit(ctx context.Context, columns ...string) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Omit")
		defer span.Finish()
	}

	w.obj = w.obj.Omit(columns...)
	return w
}

func (w GORMDBWrapper) Or(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Or")
		defer span.Finish()
	}

	w.obj = w.obj.Or(query, args...)
	return w
}

func (w GORMDBWrapper) Order(ctx context.Context, value interface{}, reorder ...bool) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Order")
		defer span.Finish()
	}

	w.obj = w.obj.Order(value, reorder...)
	return w
}

func (w GORMDBWrapper) Pluck(ctx context.Context, column string, value interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Pluck")
		defer span.Finish()
	}

	w.obj = w.obj.Pluck(column, value)
	return w
}

func (w GORMDBWrapper) Preload(ctx context.Context, column string, conditions ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Preload")
		defer span.Finish()
	}

	w.obj = w.obj.Preload(column, conditions...)
	return w
}

func (w GORMDBWrapper) Preloads(ctx context.Context, out interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Preloads")
		defer span.Finish()
	}

	w.obj = w.obj.Preloads(out)
	return w
}

func (w GORMDBWrapper) QueryExpr(ctx context.Context) *gorm.SqlExpr {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.QueryExpr")
		defer span.Finish()
	}

	res0 := w.obj.QueryExpr()
	return res0
}

func (w GORMDBWrapper) Raw(ctx context.Context, sql string, values ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Raw")
		defer span.Finish()
	}

	w.obj = w.obj.Raw(sql, values...)
	return w
}

func (w GORMDBWrapper) RecordNotFound(ctx context.Context) bool {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.RecordNotFound")
		defer span.Finish()
	}

	res0 := w.obj.RecordNotFound()
	return res0
}

func (w GORMDBWrapper) Related(ctx context.Context, value interface{}, foreignKeys ...string) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Related")
		defer span.Finish()
	}

	w.obj = w.obj.Related(value, foreignKeys...)
	return w
}

func (w GORMDBWrapper) RemoveForeignKey(ctx context.Context, field string, dest string) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.RemoveForeignKey")
		defer span.Finish()
	}

	w.obj = w.obj.RemoveForeignKey(field, dest)
	return w
}

func (w GORMDBWrapper) RemoveIndex(ctx context.Context, indexName string) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.RemoveIndex")
		defer span.Finish()
	}

	w.obj = w.obj.RemoveIndex(indexName)
	return w
}

func (w GORMDBWrapper) Rollback(ctx context.Context) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Rollback")
		defer span.Finish()
	}

	w.obj = w.obj.Rollback()
	return w
}

func (w GORMDBWrapper) RollbackUnlessCommitted(ctx context.Context) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.RollbackUnlessCommitted")
		defer span.Finish()
	}

	w.obj = w.obj.RollbackUnlessCommitted()
	return w
}

func (w GORMDBWrapper) Row(ctx context.Context) *sql.Row {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Row")
		defer span.Finish()
	}

	res0 := w.obj.Row()
	return res0
}

func (w GORMDBWrapper) Rows(ctx context.Context) (*sql.Rows, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Rows")
		defer span.Finish()
	}

	var res0 *sql.Rows
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.Rows()
		return err
	})
	return res0, err
}

func (w GORMDBWrapper) Save(ctx context.Context, value interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Save")
		defer span.Finish()
	}

	w.obj = w.obj.Save(value)
	return w
}

func (w GORMDBWrapper) Scan(ctx context.Context, dest interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Scan")
		defer span.Finish()
	}

	w.obj = w.obj.Scan(dest)
	return w
}

func (w GORMDBWrapper) ScanRows(ctx context.Context, rows *sql.Rows, result interface{}) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.ScanRows")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.ScanRows(rows, result)
		return err
	})
	return err
}

func (w GORMDBWrapper) Scopes(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Scopes")
		defer span.Finish()
	}

	w.obj = w.obj.Scopes(funcs...)
	return w
}

func (w GORMDBWrapper) Select(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Select")
		defer span.Finish()
	}

	w.obj = w.obj.Select(query, args...)
	return w
}

func (w GORMDBWrapper) Set(ctx context.Context, name string, value interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Set")
		defer span.Finish()
	}

	w.obj = w.obj.Set(name, value)
	return w
}

func (w GORMDBWrapper) SetJoinTableHandler(ctx context.Context, source interface{}, column string, handler gorm.JoinTableHandlerInterface) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.SetJoinTableHandler")
		defer span.Finish()
	}
	w.obj.SetJoinTableHandler(source, column, handler)
}

func (w GORMDBWrapper) SetNowFuncOverride(ctx context.Context, nowFuncOverride func() time.Time) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.SetNowFuncOverride")
		defer span.Finish()
	}

	w.obj = w.obj.SetNowFuncOverride(nowFuncOverride)
	return w
}

func (w GORMDBWrapper) SingularTable(ctx context.Context, enable bool) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.SingularTable")
		defer span.Finish()
	}
	w.obj.SingularTable(enable)
}

func (w GORMDBWrapper) SubQuery(ctx context.Context) *gorm.SqlExpr {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.SubQuery")
		defer span.Finish()
	}

	res0 := w.obj.SubQuery()
	return res0
}

func (w GORMDBWrapper) Table(ctx context.Context, name string) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Table")
		defer span.Finish()
	}

	w.obj = w.obj.Table(name)
	return w
}

func (w GORMDBWrapper) Take(ctx context.Context, out interface{}, where ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Take")
		defer span.Finish()
	}

	w.obj = w.obj.Take(out, where...)
	return w
}

func (w GORMDBWrapper) Transaction(ctx context.Context, fc func(tx *gorm.DB) error) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Transaction")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.Transaction(fc)
		return err
	})
	return err
}

func (w GORMDBWrapper) Unscoped(ctx context.Context) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Unscoped")
		defer span.Finish()
	}

	w.obj = w.obj.Unscoped()
	return w
}

func (w GORMDBWrapper) Update(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Update")
		defer span.Finish()
	}

	w.obj = w.obj.Update(attrs...)
	return w
}

func (w GORMDBWrapper) UpdateColumn(ctx context.Context, attrs ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.UpdateColumn")
		defer span.Finish()
	}

	w.obj = w.obj.UpdateColumn(attrs...)
	return w
}

func (w GORMDBWrapper) UpdateColumns(ctx context.Context, values interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.UpdateColumns")
		defer span.Finish()
	}

	w.obj = w.obj.UpdateColumns(values)
	return w
}

func (w GORMDBWrapper) Updates(ctx context.Context, values interface{}, ignoreProtectedAttrs ...bool) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Updates")
		defer span.Finish()
	}

	w.obj = w.obj.Updates(values, ignoreProtectedAttrs...)
	return w
}

func (w GORMDBWrapper) Where(ctx context.Context, query interface{}, args ...interface{}) GORMDBWrapper {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "gorm.DB.Where")
		defer span.Finish()
	}

	w.obj = w.obj.Where(query, args...)
	return w
}
