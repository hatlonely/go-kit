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

func (w *OTSTableStoreClientWrapper) CreateSearchIndex(ctx context.Context, request *tablestore.CreateSearchIndexRequest) (*tablestore.CreateSearchIndexResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CreateSearchIndex")
	defer span.Finish()

	var createSearchIndexResponse *tablestore.CreateSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		createSearchIndexResponse, err = w.obj.CreateSearchIndex(request)
		return err
	})
	return createSearchIndexResponse, err
}

func (w *OTSTableStoreClientWrapper) DeleteSearchIndex(ctx context.Context, request *tablestore.DeleteSearchIndexRequest) (*tablestore.DeleteSearchIndexResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteSearchIndex")
	defer span.Finish()

	var deleteSearchIndexResponse *tablestore.DeleteSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		deleteSearchIndexResponse, err = w.obj.DeleteSearchIndex(request)
		return err
	})
	return deleteSearchIndexResponse, err
}

func (w *OTSTableStoreClientWrapper) ListSearchIndex(ctx context.Context, request *tablestore.ListSearchIndexRequest) (*tablestore.ListSearchIndexResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ListSearchIndex")
	defer span.Finish()

	var listSearchIndexResponse *tablestore.ListSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		listSearchIndexResponse, err = w.obj.ListSearchIndex(request)
		return err
	})
	return listSearchIndexResponse, err
}

func (w *OTSTableStoreClientWrapper) DescribeSearchIndex(ctx context.Context, request *tablestore.DescribeSearchIndexRequest) (*tablestore.DescribeSearchIndexResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DescribeSearchIndex")
	defer span.Finish()

	var describeSearchIndexResponse *tablestore.DescribeSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		describeSearchIndexResponse, err = w.obj.DescribeSearchIndex(request)
		return err
	})
	return describeSearchIndexResponse, err
}

func (w *OTSTableStoreClientWrapper) Search(ctx context.Context, request *tablestore.SearchRequest) (*tablestore.SearchResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.Search")
	defer span.Finish()

	var searchResponse *tablestore.SearchResponse
	var err error
	err = w.retry.Do(func() error {
		searchResponse, err = w.obj.Search(request)
		return err
	})
	return searchResponse, err
}

func (w *OTSTableStoreClientWrapper) CreateTable(ctx context.Context, request *tablestore.CreateTableRequest) (*tablestore.CreateTableResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CreateTable")
	defer span.Finish()

	var createTableResponse *tablestore.CreateTableResponse
	var err error
	err = w.retry.Do(func() error {
		createTableResponse, err = w.obj.CreateTable(request)
		return err
	})
	return createTableResponse, err
}

func (w *OTSTableStoreClientWrapper) CreateIndex(ctx context.Context, request *tablestore.CreateIndexRequest) (*tablestore.CreateIndexResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CreateIndex")
	defer span.Finish()

	var createIndexResponse *tablestore.CreateIndexResponse
	var err error
	err = w.retry.Do(func() error {
		createIndexResponse, err = w.obj.CreateIndex(request)
		return err
	})
	return createIndexResponse, err
}

func (w *OTSTableStoreClientWrapper) DeleteIndex(ctx context.Context, request *tablestore.DeleteIndexRequest) (*tablestore.DeleteIndexResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteIndex")
	defer span.Finish()

	var deleteIndexResponse *tablestore.DeleteIndexResponse
	var err error
	err = w.retry.Do(func() error {
		deleteIndexResponse, err = w.obj.DeleteIndex(request)
		return err
	})
	return deleteIndexResponse, err
}

func (w *OTSTableStoreClientWrapper) ListTable(ctx context.Context) (*tablestore.ListTableResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ListTable")
	defer span.Finish()

	var listTableResponse *tablestore.ListTableResponse
	var err error
	err = w.retry.Do(func() error {
		listTableResponse, err = w.obj.ListTable()
		return err
	})
	return listTableResponse, err
}

func (w *OTSTableStoreClientWrapper) DeleteTable(ctx context.Context, request *tablestore.DeleteTableRequest) (*tablestore.DeleteTableResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteTable")
	defer span.Finish()

	var deleteTableResponse *tablestore.DeleteTableResponse
	var err error
	err = w.retry.Do(func() error {
		deleteTableResponse, err = w.obj.DeleteTable(request)
		return err
	})
	return deleteTableResponse, err
}

func (w *OTSTableStoreClientWrapper) DescribeTable(ctx context.Context, request *tablestore.DescribeTableRequest) (*tablestore.DescribeTableResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DescribeTable")
	defer span.Finish()

	var describeTableResponse *tablestore.DescribeTableResponse
	var err error
	err = w.retry.Do(func() error {
		describeTableResponse, err = w.obj.DescribeTable(request)
		return err
	})
	return describeTableResponse, err
}

func (w *OTSTableStoreClientWrapper) UpdateTable(ctx context.Context, request *tablestore.UpdateTableRequest) (*tablestore.UpdateTableResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.UpdateTable")
	defer span.Finish()

	var updateTableResponse *tablestore.UpdateTableResponse
	var err error
	err = w.retry.Do(func() error {
		updateTableResponse, err = w.obj.UpdateTable(request)
		return err
	})
	return updateTableResponse, err
}

func (w *OTSTableStoreClientWrapper) PutRow(ctx context.Context, request *tablestore.PutRowRequest) (*tablestore.PutRowResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.PutRow")
	defer span.Finish()

	var putRowResponse *tablestore.PutRowResponse
	var err error
	err = w.retry.Do(func() error {
		putRowResponse, err = w.obj.PutRow(request)
		return err
	})
	return putRowResponse, err
}

func (w *OTSTableStoreClientWrapper) DeleteRow(ctx context.Context, request *tablestore.DeleteRowRequest) (*tablestore.DeleteRowResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteRow")
	defer span.Finish()

	var deleteRowResponse *tablestore.DeleteRowResponse
	var err error
	err = w.retry.Do(func() error {
		deleteRowResponse, err = w.obj.DeleteRow(request)
		return err
	})
	return deleteRowResponse, err
}

func (w *OTSTableStoreClientWrapper) GetRow(ctx context.Context, request *tablestore.GetRowRequest) (*tablestore.GetRowResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetRow")
	defer span.Finish()

	var getRowResponse *tablestore.GetRowResponse
	var err error
	err = w.retry.Do(func() error {
		getRowResponse, err = w.obj.GetRow(request)
		return err
	})
	return getRowResponse, err
}

func (w *OTSTableStoreClientWrapper) UpdateRow(ctx context.Context, request *tablestore.UpdateRowRequest) (*tablestore.UpdateRowResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.UpdateRow")
	defer span.Finish()

	var updateRowResponse *tablestore.UpdateRowResponse
	var err error
	err = w.retry.Do(func() error {
		updateRowResponse, err = w.obj.UpdateRow(request)
		return err
	})
	return updateRowResponse, err
}

func (w *OTSTableStoreClientWrapper) BatchGetRow(ctx context.Context, request *tablestore.BatchGetRowRequest) (*tablestore.BatchGetRowResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.BatchGetRow")
	defer span.Finish()

	var batchGetRowResponse *tablestore.BatchGetRowResponse
	var err error
	err = w.retry.Do(func() error {
		batchGetRowResponse, err = w.obj.BatchGetRow(request)
		return err
	})
	return batchGetRowResponse, err
}

func (w *OTSTableStoreClientWrapper) BatchWriteRow(ctx context.Context, request *tablestore.BatchWriteRowRequest) (*tablestore.BatchWriteRowResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.BatchWriteRow")
	defer span.Finish()

	var batchWriteRowResponse *tablestore.BatchWriteRowResponse
	var err error
	err = w.retry.Do(func() error {
		batchWriteRowResponse, err = w.obj.BatchWriteRow(request)
		return err
	})
	return batchWriteRowResponse, err
}

func (w *OTSTableStoreClientWrapper) GetRange(ctx context.Context, request *tablestore.GetRangeRequest) (*tablestore.GetRangeResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetRange")
	defer span.Finish()

	var getRangeResponse *tablestore.GetRangeResponse
	var err error
	err = w.retry.Do(func() error {
		getRangeResponse, err = w.obj.GetRange(request)
		return err
	})
	return getRangeResponse, err
}

func (w *OTSTableStoreClientWrapper) ListStream(ctx context.Context, req *tablestore.ListStreamRequest) (*tablestore.ListStreamResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ListStream")
	defer span.Finish()

	var listStreamResponse *tablestore.ListStreamResponse
	var err error
	err = w.retry.Do(func() error {
		listStreamResponse, err = w.obj.ListStream(req)
		return err
	})
	return listStreamResponse, err
}

func (w *OTSTableStoreClientWrapper) DescribeStream(ctx context.Context, req *tablestore.DescribeStreamRequest) (*tablestore.DescribeStreamResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DescribeStream")
	defer span.Finish()

	var describeStreamResponse *tablestore.DescribeStreamResponse
	var err error
	err = w.retry.Do(func() error {
		describeStreamResponse, err = w.obj.DescribeStream(req)
		return err
	})
	return describeStreamResponse, err
}

func (w *OTSTableStoreClientWrapper) GetShardIterator(ctx context.Context, req *tablestore.GetShardIteratorRequest) (*tablestore.GetShardIteratorResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetShardIterator")
	defer span.Finish()

	var getShardIteratorResponse *tablestore.GetShardIteratorResponse
	var err error
	err = w.retry.Do(func() error {
		getShardIteratorResponse, err = w.obj.GetShardIterator(req)
		return err
	})
	return getShardIteratorResponse, err
}

func (w *OTSTableStoreClientWrapper) GetStreamRecord(ctx context.Context, req *tablestore.GetStreamRecordRequest) (*tablestore.GetStreamRecordResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetStreamRecord")
	defer span.Finish()

	var getStreamRecordResponse *tablestore.GetStreamRecordResponse
	var err error
	err = w.retry.Do(func() error {
		getStreamRecordResponse, err = w.obj.GetStreamRecord(req)
		return err
	})
	return getStreamRecordResponse, err
}

func (w *OTSTableStoreClientWrapper) ComputeSplitPointsBySize(ctx context.Context, req *tablestore.ComputeSplitPointsBySizeRequest) (*tablestore.ComputeSplitPointsBySizeResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ComputeSplitPointsBySize")
	defer span.Finish()

	var computeSplitPointsBySizeResponse *tablestore.ComputeSplitPointsBySizeResponse
	var err error
	err = w.retry.Do(func() error {
		computeSplitPointsBySizeResponse, err = w.obj.ComputeSplitPointsBySize(req)
		return err
	})
	return computeSplitPointsBySizeResponse, err
}

func (w *OTSTableStoreClientWrapper) StartLocalTransaction(ctx context.Context, request *tablestore.StartLocalTransactionRequest) (*tablestore.StartLocalTransactionResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.StartLocalTransaction")
	defer span.Finish()

	var startLocalTransactionResponse *tablestore.StartLocalTransactionResponse
	var err error
	err = w.retry.Do(func() error {
		startLocalTransactionResponse, err = w.obj.StartLocalTransaction(request)
		return err
	})
	return startLocalTransactionResponse, err
}

func (w *OTSTableStoreClientWrapper) CommitTransaction(ctx context.Context, request *tablestore.CommitTransactionRequest) (*tablestore.CommitTransactionResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CommitTransaction")
	defer span.Finish()

	var commitTransactionResponse *tablestore.CommitTransactionResponse
	var err error
	err = w.retry.Do(func() error {
		commitTransactionResponse, err = w.obj.CommitTransaction(request)
		return err
	})
	return commitTransactionResponse, err
}

func (w *OTSTableStoreClientWrapper) AbortTransaction(ctx context.Context, request *tablestore.AbortTransactionRequest) (*tablestore.AbortTransactionResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.AbortTransaction")
	defer span.Finish()

	var abortTransactionResponse *tablestore.AbortTransactionResponse
	var err error
	err = w.retry.Do(func() error {
		abortTransactionResponse, err = w.obj.AbortTransaction(request)
		return err
	})
	return abortTransactionResponse, err
}
