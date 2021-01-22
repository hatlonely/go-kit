// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/opentracing/opentracing-go"
)

type OTSTableStoreClientWrapper struct {
	obj   *tablestore.TableStoreClient
	retry *Retry
}

func (w *OTSTableStoreClientWrapper) AbortTransaction(ctx context.Context, request *tablestore.AbortTransactionRequest) (*tablestore.AbortTransactionResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.AbortTransaction")
	defer span.Finish()

	var res0 *tablestore.AbortTransactionResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.AbortTransaction(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) BatchGetRow(ctx context.Context, request *tablestore.BatchGetRowRequest) (*tablestore.BatchGetRowResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.BatchGetRow")
	defer span.Finish()

	var res0 *tablestore.BatchGetRowResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.BatchGetRow(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) BatchWriteRow(ctx context.Context, request *tablestore.BatchWriteRowRequest) (*tablestore.BatchWriteRowResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.BatchWriteRow")
	defer span.Finish()

	var res0 *tablestore.BatchWriteRowResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.BatchWriteRow(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) CommitTransaction(ctx context.Context, request *tablestore.CommitTransactionRequest) (*tablestore.CommitTransactionResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CommitTransaction")
	defer span.Finish()

	var res0 *tablestore.CommitTransactionResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.CommitTransaction(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) ComputeSplitPointsBySize(ctx context.Context, req *tablestore.ComputeSplitPointsBySizeRequest) (*tablestore.ComputeSplitPointsBySizeResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ComputeSplitPointsBySize")
	defer span.Finish()

	var res0 *tablestore.ComputeSplitPointsBySizeResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.ComputeSplitPointsBySize(req)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) CreateIndex(ctx context.Context, request *tablestore.CreateIndexRequest) (*tablestore.CreateIndexResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CreateIndex")
	defer span.Finish()

	var res0 *tablestore.CreateIndexResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.CreateIndex(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) CreateSearchIndex(ctx context.Context, request *tablestore.CreateSearchIndexRequest) (*tablestore.CreateSearchIndexResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CreateSearchIndex")
	defer span.Finish()

	var res0 *tablestore.CreateSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.CreateSearchIndex(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) CreateTable(ctx context.Context, request *tablestore.CreateTableRequest) (*tablestore.CreateTableResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CreateTable")
	defer span.Finish()

	var res0 *tablestore.CreateTableResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.CreateTable(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DeleteIndex(ctx context.Context, request *tablestore.DeleteIndexRequest) (*tablestore.DeleteIndexResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteIndex")
	defer span.Finish()

	var res0 *tablestore.DeleteIndexResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.DeleteIndex(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DeleteRow(ctx context.Context, request *tablestore.DeleteRowRequest) (*tablestore.DeleteRowResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteRow")
	defer span.Finish()

	var res0 *tablestore.DeleteRowResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.DeleteRow(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DeleteSearchIndex(ctx context.Context, request *tablestore.DeleteSearchIndexRequest) (*tablestore.DeleteSearchIndexResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteSearchIndex")
	defer span.Finish()

	var res0 *tablestore.DeleteSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.DeleteSearchIndex(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DeleteTable(ctx context.Context, request *tablestore.DeleteTableRequest) (*tablestore.DeleteTableResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteTable")
	defer span.Finish()

	var res0 *tablestore.DeleteTableResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.DeleteTable(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DescribeSearchIndex(ctx context.Context, request *tablestore.DescribeSearchIndexRequest) (*tablestore.DescribeSearchIndexResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DescribeSearchIndex")
	defer span.Finish()

	var res0 *tablestore.DescribeSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.DescribeSearchIndex(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DescribeStream(ctx context.Context, req *tablestore.DescribeStreamRequest) (*tablestore.DescribeStreamResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DescribeStream")
	defer span.Finish()

	var res0 *tablestore.DescribeStreamResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.DescribeStream(req)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DescribeTable(ctx context.Context, request *tablestore.DescribeTableRequest) (*tablestore.DescribeTableResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DescribeTable")
	defer span.Finish()

	var res0 *tablestore.DescribeTableResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.DescribeTable(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) GetRange(ctx context.Context, request *tablestore.GetRangeRequest) (*tablestore.GetRangeResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetRange")
	defer span.Finish()

	var res0 *tablestore.GetRangeResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.GetRange(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) GetRow(ctx context.Context, request *tablestore.GetRowRequest) (*tablestore.GetRowResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetRow")
	defer span.Finish()

	var res0 *tablestore.GetRowResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.GetRow(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) GetShardIterator(ctx context.Context, req *tablestore.GetShardIteratorRequest) (*tablestore.GetShardIteratorResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetShardIterator")
	defer span.Finish()

	var res0 *tablestore.GetShardIteratorResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.GetShardIterator(req)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) GetStreamRecord(ctx context.Context, req *tablestore.GetStreamRecordRequest) (*tablestore.GetStreamRecordResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetStreamRecord")
	defer span.Finish()

	var res0 *tablestore.GetStreamRecordResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.GetStreamRecord(req)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) ListSearchIndex(ctx context.Context, request *tablestore.ListSearchIndexRequest) (*tablestore.ListSearchIndexResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ListSearchIndex")
	defer span.Finish()

	var res0 *tablestore.ListSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.ListSearchIndex(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) ListStream(ctx context.Context, req *tablestore.ListStreamRequest) (*tablestore.ListStreamResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ListStream")
	defer span.Finish()

	var res0 *tablestore.ListStreamResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.ListStream(req)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) ListTable(ctx context.Context) (*tablestore.ListTableResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ListTable")
	defer span.Finish()

	var res0 *tablestore.ListTableResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.ListTable()
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) PutRow(ctx context.Context, request *tablestore.PutRowRequest) (*tablestore.PutRowResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.PutRow")
	defer span.Finish()

	var res0 *tablestore.PutRowResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.PutRow(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) Search(ctx context.Context, request *tablestore.SearchRequest) (*tablestore.SearchResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.Search")
	defer span.Finish()

	var res0 *tablestore.SearchResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.Search(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) StartLocalTransaction(ctx context.Context, request *tablestore.StartLocalTransactionRequest) (*tablestore.StartLocalTransactionResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.StartLocalTransaction")
	defer span.Finish()

	var res0 *tablestore.StartLocalTransactionResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.StartLocalTransaction(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) UpdateRow(ctx context.Context, request *tablestore.UpdateRowRequest) (*tablestore.UpdateRowResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.UpdateRow")
	defer span.Finish()

	var res0 *tablestore.UpdateRowResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.UpdateRow(request)
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) UpdateTable(ctx context.Context, request *tablestore.UpdateTableRequest) (*tablestore.UpdateTableResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.UpdateTable")
	defer span.Finish()

	var res0 *tablestore.UpdateTableResponse
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.UpdateTable(request)
		return err
	})
	return res0, err
}
