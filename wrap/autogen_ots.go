// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type OTSTableStoreClientWrapper struct {
	obj            *tablestore.TableStoreClient
	retry          *Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	totalMetric    *prometheus.CounterVec
}

func (w *OTSTableStoreClientWrapper) Unwrap() *tablestore.TableStoreClient {
	return w.obj
}

func (w *OTSTableStoreClientWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *OTSTableStoreClientWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *OTSTableStoreClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "tablestore_TableStoreClient_durationMs",
		Help:        "tablestore TableStoreClient response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "tablestore_TableStoreClient_total",
		Help:        "tablestore TableStoreClient request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode"})
}

func (w *OTSTableStoreClientWrapper) AbortTransaction(ctx context.Context, request *tablestore.AbortTransactionRequest) (*tablestore.AbortTransactionResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.AbortTransaction")
		defer span.Finish()
	}

	var res0 *tablestore.AbortTransactionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.AbortTransaction", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.AbortTransaction", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.AbortTransaction(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) BatchGetRow(ctx context.Context, request *tablestore.BatchGetRowRequest) (*tablestore.BatchGetRowResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.BatchGetRow")
		defer span.Finish()
	}

	var res0 *tablestore.BatchGetRowResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.BatchGetRow", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.BatchGetRow", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.BatchGetRow(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) BatchWriteRow(ctx context.Context, request *tablestore.BatchWriteRowRequest) (*tablestore.BatchWriteRowResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.BatchWriteRow")
		defer span.Finish()
	}

	var res0 *tablestore.BatchWriteRowResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.BatchWriteRow", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.BatchWriteRow", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.BatchWriteRow(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) CommitTransaction(ctx context.Context, request *tablestore.CommitTransactionRequest) (*tablestore.CommitTransactionResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CommitTransaction")
		defer span.Finish()
	}

	var res0 *tablestore.CommitTransactionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.CommitTransaction", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.CommitTransaction", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.CommitTransaction(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) ComputeSplitPointsBySize(ctx context.Context, req *tablestore.ComputeSplitPointsBySizeRequest) (*tablestore.ComputeSplitPointsBySizeResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ComputeSplitPointsBySize")
		defer span.Finish()
	}

	var res0 *tablestore.ComputeSplitPointsBySizeResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.ComputeSplitPointsBySize", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.ComputeSplitPointsBySize", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.ComputeSplitPointsBySize(req)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) CreateIndex(ctx context.Context, request *tablestore.CreateIndexRequest) (*tablestore.CreateIndexResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CreateIndex")
		defer span.Finish()
	}

	var res0 *tablestore.CreateIndexResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.CreateIndex", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.CreateIndex", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.CreateIndex(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) CreateSearchIndex(ctx context.Context, request *tablestore.CreateSearchIndexRequest) (*tablestore.CreateSearchIndexResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CreateSearchIndex")
		defer span.Finish()
	}

	var res0 *tablestore.CreateSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.CreateSearchIndex", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.CreateSearchIndex", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.CreateSearchIndex(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) CreateTable(ctx context.Context, request *tablestore.CreateTableRequest) (*tablestore.CreateTableResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CreateTable")
		defer span.Finish()
	}

	var res0 *tablestore.CreateTableResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.CreateTable", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.CreateTable", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.CreateTable(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DeleteIndex(ctx context.Context, request *tablestore.DeleteIndexRequest) (*tablestore.DeleteIndexResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteIndex")
		defer span.Finish()
	}

	var res0 *tablestore.DeleteIndexResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.DeleteIndex", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.DeleteIndex", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.DeleteIndex(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DeleteRow(ctx context.Context, request *tablestore.DeleteRowRequest) (*tablestore.DeleteRowResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteRow")
		defer span.Finish()
	}

	var res0 *tablestore.DeleteRowResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.DeleteRow", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.DeleteRow", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.DeleteRow(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DeleteSearchIndex(ctx context.Context, request *tablestore.DeleteSearchIndexRequest) (*tablestore.DeleteSearchIndexResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteSearchIndex")
		defer span.Finish()
	}

	var res0 *tablestore.DeleteSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.DeleteSearchIndex", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.DeleteSearchIndex", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.DeleteSearchIndex(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DeleteTable(ctx context.Context, request *tablestore.DeleteTableRequest) (*tablestore.DeleteTableResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteTable")
		defer span.Finish()
	}

	var res0 *tablestore.DeleteTableResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.DeleteTable", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.DeleteTable", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.DeleteTable(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DescribeSearchIndex(ctx context.Context, request *tablestore.DescribeSearchIndexRequest) (*tablestore.DescribeSearchIndexResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DescribeSearchIndex")
		defer span.Finish()
	}

	var res0 *tablestore.DescribeSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.DescribeSearchIndex", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.DescribeSearchIndex", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.DescribeSearchIndex(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DescribeStream(ctx context.Context, req *tablestore.DescribeStreamRequest) (*tablestore.DescribeStreamResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DescribeStream")
		defer span.Finish()
	}

	var res0 *tablestore.DescribeStreamResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.DescribeStream", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.DescribeStream", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.DescribeStream(req)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DescribeTable(ctx context.Context, request *tablestore.DescribeTableRequest) (*tablestore.DescribeTableResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DescribeTable")
		defer span.Finish()
	}

	var res0 *tablestore.DescribeTableResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.DescribeTable", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.DescribeTable", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.DescribeTable(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) GetRange(ctx context.Context, request *tablestore.GetRangeRequest) (*tablestore.GetRangeResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetRange")
		defer span.Finish()
	}

	var res0 *tablestore.GetRangeResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.GetRange", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.GetRange", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.GetRange(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) GetRow(ctx context.Context, request *tablestore.GetRowRequest) (*tablestore.GetRowResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetRow")
		defer span.Finish()
	}

	var res0 *tablestore.GetRowResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.GetRow", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.GetRow", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.GetRow(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) GetShardIterator(ctx context.Context, req *tablestore.GetShardIteratorRequest) (*tablestore.GetShardIteratorResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetShardIterator")
		defer span.Finish()
	}

	var res0 *tablestore.GetShardIteratorResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.GetShardIterator", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.GetShardIterator", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.GetShardIterator(req)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) GetStreamRecord(ctx context.Context, req *tablestore.GetStreamRecordRequest) (*tablestore.GetStreamRecordResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetStreamRecord")
		defer span.Finish()
	}

	var res0 *tablestore.GetStreamRecordResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.GetStreamRecord", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.GetStreamRecord", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.GetStreamRecord(req)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) ListSearchIndex(ctx context.Context, request *tablestore.ListSearchIndexRequest) (*tablestore.ListSearchIndexResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ListSearchIndex")
		defer span.Finish()
	}

	var res0 *tablestore.ListSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.ListSearchIndex", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.ListSearchIndex", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.ListSearchIndex(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) ListStream(ctx context.Context, req *tablestore.ListStreamRequest) (*tablestore.ListStreamResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ListStream")
		defer span.Finish()
	}

	var res0 *tablestore.ListStreamResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.ListStream", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.ListStream", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.ListStream(req)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) ListTable(ctx context.Context) (*tablestore.ListTableResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ListTable")
		defer span.Finish()
	}

	var res0 *tablestore.ListTableResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.ListTable", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.ListTable", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.ListTable()
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) PutRow(ctx context.Context, request *tablestore.PutRowRequest) (*tablestore.PutRowResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.PutRow")
		defer span.Finish()
	}

	var res0 *tablestore.PutRowResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.PutRow", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.PutRow", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.PutRow(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) Search(ctx context.Context, request *tablestore.SearchRequest) (*tablestore.SearchResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.Search")
		defer span.Finish()
	}

	var res0 *tablestore.SearchResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.Search", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.Search", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.Search(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) StartLocalTransaction(ctx context.Context, request *tablestore.StartLocalTransactionRequest) (*tablestore.StartLocalTransactionResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.StartLocalTransaction")
		defer span.Finish()
	}

	var res0 *tablestore.StartLocalTransactionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.StartLocalTransaction", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.StartLocalTransaction", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.StartLocalTransaction(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) UpdateRow(ctx context.Context, request *tablestore.UpdateRowRequest) (*tablestore.UpdateRowResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.UpdateRow")
		defer span.Finish()
	}

	var res0 *tablestore.UpdateRowResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.UpdateRow", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.UpdateRow", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.UpdateRow(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) UpdateTable(ctx context.Context, request *tablestore.UpdateTableRequest) (*tablestore.UpdateTableResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.UpdateTable")
		defer span.Finish()
	}

	var res0 *tablestore.UpdateTableResponse
	var err error
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("tablestore.TableStoreClient.UpdateTable", ErrCode(err)).Inc()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.UpdateTable", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		res0, err = w.obj.UpdateTable(request)
		return err
	})
	return res0, err
}
