// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

func NewOTSTableStoreClientWrapper(
	obj *tablestore.TableStoreClient,
	retry *micro.Retry,
	options *WrapperOptions,
	durationMetric *prometheus.HistogramVec,
	inflightMetric *prometheus.GaugeVec,
	rateLimiter micro.RateLimiter,
	parallelController micro.ParallelController) *OTSTableStoreClientWrapper {
	return &OTSTableStoreClientWrapper{
		obj:                obj,
		retry:              retry,
		options:            options,
		durationMetric:     durationMetric,
		inflightMetric:     inflightMetric,
		rateLimiter:        rateLimiter,
		parallelController: parallelController,
	}
}

type OTSTableStoreClientWrapper struct {
	obj                *tablestore.TableStoreClient
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
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

func (w *OTSTableStoreClientWrapper) OnRateLimiterChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *OTSTableStoreClientWrapper) OnParallelControllerChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *OTSTableStoreClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "tablestore_TableStoreClient_durationMs",
		Help:        "tablestore TableStoreClient response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "tablestore_TableStoreClient_inflight",
		Help:        "tablestore TableStoreClient inflight",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "custom"})
}

func (w *OTSTableStoreClientWrapper) AbortTransaction(ctx context.Context, request *tablestore.AbortTransactionRequest) (*tablestore.AbortTransactionResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.AbortTransactionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.AbortTransaction", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.AbortTransaction", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.AbortTransaction", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.AbortTransaction")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.AbortTransaction", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.AbortTransaction", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.AbortTransaction", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.AbortTransaction(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) BatchGetRow(ctx context.Context, request *tablestore.BatchGetRowRequest) (*tablestore.BatchGetRowResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.BatchGetRowResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.BatchGetRow", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.BatchGetRow", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.BatchGetRow", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.BatchGetRow")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.BatchGetRow", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.BatchGetRow", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.BatchGetRow", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.BatchGetRow(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) BatchWriteRow(ctx context.Context, request *tablestore.BatchWriteRowRequest) (*tablestore.BatchWriteRowResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.BatchWriteRowResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.BatchWriteRow", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.BatchWriteRow", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.BatchWriteRow", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.BatchWriteRow")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.BatchWriteRow", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.BatchWriteRow", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.BatchWriteRow", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.BatchWriteRow(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) CommitTransaction(ctx context.Context, request *tablestore.CommitTransactionRequest) (*tablestore.CommitTransactionResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.CommitTransactionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.CommitTransaction", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.CommitTransaction", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.CommitTransaction", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CommitTransaction")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.CommitTransaction", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.CommitTransaction", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.CommitTransaction", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.CommitTransaction(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) ComputeSplitPointsBySize(ctx context.Context, req *tablestore.ComputeSplitPointsBySizeRequest) (*tablestore.ComputeSplitPointsBySizeResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.ComputeSplitPointsBySizeResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.ComputeSplitPointsBySize", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.ComputeSplitPointsBySize", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.ComputeSplitPointsBySize", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ComputeSplitPointsBySize")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.ComputeSplitPointsBySize", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.ComputeSplitPointsBySize", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.ComputeSplitPointsBySize", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ComputeSplitPointsBySize(req)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) CreateIndex(ctx context.Context, request *tablestore.CreateIndexRequest) (*tablestore.CreateIndexResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.CreateIndexResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.CreateIndex", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.CreateIndex", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.CreateIndex", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CreateIndex")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.CreateIndex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.CreateIndex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.CreateIndex", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.CreateIndex(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) CreateSearchIndex(ctx context.Context, request *tablestore.CreateSearchIndexRequest) (*tablestore.CreateSearchIndexResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.CreateSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.CreateSearchIndex", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.CreateSearchIndex", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.CreateSearchIndex", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CreateSearchIndex")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.CreateSearchIndex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.CreateSearchIndex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.CreateSearchIndex", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.CreateSearchIndex(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) CreateTable(ctx context.Context, request *tablestore.CreateTableRequest) (*tablestore.CreateTableResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.CreateTableResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.CreateTable", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.CreateTable", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.CreateTable", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.CreateTable")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.CreateTable", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.CreateTable", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.CreateTable", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.CreateTable(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DeleteIndex(ctx context.Context, request *tablestore.DeleteIndexRequest) (*tablestore.DeleteIndexResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.DeleteIndexResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.DeleteIndex", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.DeleteIndex", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.DeleteIndex", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteIndex")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.DeleteIndex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.DeleteIndex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.DeleteIndex", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DeleteIndex(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DeleteRow(ctx context.Context, request *tablestore.DeleteRowRequest) (*tablestore.DeleteRowResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.DeleteRowResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.DeleteRow", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.DeleteRow", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.DeleteRow", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteRow")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.DeleteRow", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.DeleteRow", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.DeleteRow", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DeleteRow(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DeleteSearchIndex(ctx context.Context, request *tablestore.DeleteSearchIndexRequest) (*tablestore.DeleteSearchIndexResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.DeleteSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.DeleteSearchIndex", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.DeleteSearchIndex", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.DeleteSearchIndex", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteSearchIndex")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.DeleteSearchIndex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.DeleteSearchIndex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.DeleteSearchIndex", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DeleteSearchIndex(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DeleteTable(ctx context.Context, request *tablestore.DeleteTableRequest) (*tablestore.DeleteTableResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.DeleteTableResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.DeleteTable", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.DeleteTable", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.DeleteTable", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DeleteTable")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.DeleteTable", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.DeleteTable", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.DeleteTable", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DeleteTable(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DescribeSearchIndex(ctx context.Context, request *tablestore.DescribeSearchIndexRequest) (*tablestore.DescribeSearchIndexResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.DescribeSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.DescribeSearchIndex", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.DescribeSearchIndex", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.DescribeSearchIndex", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DescribeSearchIndex")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.DescribeSearchIndex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.DescribeSearchIndex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.DescribeSearchIndex", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DescribeSearchIndex(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DescribeStream(ctx context.Context, req *tablestore.DescribeStreamRequest) (*tablestore.DescribeStreamResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.DescribeStreamResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.DescribeStream", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.DescribeStream", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.DescribeStream", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DescribeStream")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.DescribeStream", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.DescribeStream", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.DescribeStream", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DescribeStream(req)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) DescribeTable(ctx context.Context, request *tablestore.DescribeTableRequest) (*tablestore.DescribeTableResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.DescribeTableResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.DescribeTable", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.DescribeTable", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.DescribeTable", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.DescribeTable")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.DescribeTable", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.DescribeTable", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.DescribeTable", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DescribeTable(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) GetRange(ctx context.Context, request *tablestore.GetRangeRequest) (*tablestore.GetRangeResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.GetRangeResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.GetRange", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.GetRange", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.GetRange", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetRange")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.GetRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.GetRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.GetRange", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetRange(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) GetRow(ctx context.Context, request *tablestore.GetRowRequest) (*tablestore.GetRowResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.GetRowResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.GetRow", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.GetRow", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.GetRow", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetRow")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.GetRow", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.GetRow", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.GetRow", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetRow(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) GetShardIterator(ctx context.Context, req *tablestore.GetShardIteratorRequest) (*tablestore.GetShardIteratorResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.GetShardIteratorResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.GetShardIterator", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.GetShardIterator", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.GetShardIterator", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetShardIterator")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.GetShardIterator", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.GetShardIterator", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.GetShardIterator", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetShardIterator(req)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) GetStreamRecord(ctx context.Context, req *tablestore.GetStreamRecordRequest) (*tablestore.GetStreamRecordResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.GetStreamRecordResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.GetStreamRecord", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.GetStreamRecord", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.GetStreamRecord", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetStreamRecord")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.GetStreamRecord", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.GetStreamRecord", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.GetStreamRecord", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetStreamRecord(req)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) ListSearchIndex(ctx context.Context, request *tablestore.ListSearchIndexRequest) (*tablestore.ListSearchIndexResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.ListSearchIndexResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.ListSearchIndex", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.ListSearchIndex", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.ListSearchIndex", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ListSearchIndex")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.ListSearchIndex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.ListSearchIndex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.ListSearchIndex", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListSearchIndex(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) ListStream(ctx context.Context, req *tablestore.ListStreamRequest) (*tablestore.ListStreamResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.ListStreamResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.ListStream", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.ListStream", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.ListStream", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ListStream")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.ListStream", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.ListStream", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.ListStream", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListStream(req)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) ListTable(ctx context.Context) (*tablestore.ListTableResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.ListTableResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.ListTable", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.ListTable", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.ListTable", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.ListTable")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.ListTable", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.ListTable", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.ListTable", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListTable()
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) PutRow(ctx context.Context, request *tablestore.PutRowRequest) (*tablestore.PutRowResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.PutRowResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.PutRow", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.PutRow", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.PutRow", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.PutRow")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.PutRow", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.PutRow", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.PutRow", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.PutRow(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) Search(ctx context.Context, request *tablestore.SearchRequest) (*tablestore.SearchResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.SearchResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.Search", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.Search", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.Search", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.Search")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.Search", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.Search", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.Search", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Search(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) StartLocalTransaction(ctx context.Context, request *tablestore.StartLocalTransactionRequest) (*tablestore.StartLocalTransactionResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.StartLocalTransactionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.StartLocalTransaction", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.StartLocalTransaction", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.StartLocalTransaction", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.StartLocalTransaction")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.StartLocalTransaction", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.StartLocalTransaction", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.StartLocalTransaction", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.StartLocalTransaction(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) UpdateRow(ctx context.Context, request *tablestore.UpdateRowRequest) (*tablestore.UpdateRowResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.UpdateRowResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.UpdateRow", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.UpdateRow", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.UpdateRow", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.UpdateRow")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.UpdateRow", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.UpdateRow", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.UpdateRow", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.UpdateRow(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *OTSTableStoreClientWrapper) UpdateTable(ctx context.Context, request *tablestore.UpdateTableRequest) (*tablestore.UpdateTableResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.UpdateTableResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.UpdateTable", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.UpdateTable", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.UpdateTable", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.UpdateTable")
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
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.UpdateTable", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.UpdateTable", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.UpdateTable", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.UpdateTable(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}
