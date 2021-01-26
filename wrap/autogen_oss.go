// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type OSSBucketWrapper struct {
	obj              *oss.Bucket
	retry            *Retry
	options          *WrapperOptions
	durationMetric   *prometheus.HistogramVec
	totalMetric      *prometheus.CounterVec
	rateLimiterGroup RateLimiterGroup
}

func (w *OSSBucketWrapper) Unwrap() *oss.Bucket {
	return w.obj
}

type OSSClientWrapper struct {
	obj              *oss.Client
	retry            *Retry
	options          *WrapperOptions
	durationMetric   *prometheus.HistogramVec
	totalMetric      *prometheus.CounterVec
	rateLimiterGroup RateLimiterGroup
}

func (w *OSSClientWrapper) Unwrap() *oss.Client {
	return w.obj
}

func (w *OSSClientWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *OSSClientWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *OSSClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "oss_Client_durationMs",
		Help:        "oss Client response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "oss_Client_total",
		Help:        "oss Client request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
}

func (w *OSSBucketWrapper) AbortMultipartUpload(ctx context.Context, imur oss.InitiateMultipartUploadResult, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.AbortMultipartUpload")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.AbortMultipartUpload"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.AbortMultipartUpload", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.AbortMultipartUpload", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.AbortMultipartUpload(imur, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) AppendObject(ctx context.Context, objectKey string, reader io.Reader, appendPosition int64, options ...oss.Option) (int64, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.AppendObject")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.AppendObject"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.AppendObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.AppendObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.AppendObject(objectKey, reader, appendPosition, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) CompleteMultipartUpload(ctx context.Context, imur oss.InitiateMultipartUploadResult, parts []oss.UploadPart, options ...oss.Option) (oss.CompleteMultipartUploadResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.CompleteMultipartUpload")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.CompleteMultipartUploadResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.CompleteMultipartUpload"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.CompleteMultipartUpload", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.CompleteMultipartUpload", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.CompleteMultipartUpload(imur, parts, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) CopyFile(ctx context.Context, srcBucketName string, srcObjectKey string, destObjectKey string, partSize int64, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.CopyFile")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.CopyFile"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.CopyFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.CopyFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.CopyFile(srcBucketName, srcObjectKey, destObjectKey, partSize, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) CopyObject(ctx context.Context, srcObjectKey string, destObjectKey string, options ...oss.Option) (oss.CopyObjectResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.CopyObject")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.CopyObjectResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.CopyObject"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.CopyObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.CopyObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.CopyObject(srcObjectKey, destObjectKey, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) CopyObjectFrom(ctx context.Context, srcBucketName string, srcObjectKey string, destObjectKey string, options ...oss.Option) (oss.CopyObjectResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.CopyObjectFrom")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.CopyObjectResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.CopyObjectFrom"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.CopyObjectFrom", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.CopyObjectFrom", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.CopyObjectFrom(srcBucketName, srcObjectKey, destObjectKey, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) CopyObjectTo(ctx context.Context, destBucketName string, destObjectKey string, srcObjectKey string, options ...oss.Option) (oss.CopyObjectResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.CopyObjectTo")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.CopyObjectResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.CopyObjectTo"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.CopyObjectTo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.CopyObjectTo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.CopyObjectTo(destBucketName, destObjectKey, srcObjectKey, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) CreateLiveChannel(ctx context.Context, channelName string, config oss.LiveChannelConfiguration) (oss.CreateLiveChannelResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.CreateLiveChannel")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.CreateLiveChannelResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.CreateLiveChannel"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.CreateLiveChannel", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.CreateLiveChannel", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.CreateLiveChannel(channelName, config)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) CreateSelectCsvObjectMeta(ctx context.Context, key string, csvMeta oss.CsvMetaRequest, options ...oss.Option) (oss.MetaEndFrameCSV, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.CreateSelectCsvObjectMeta")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.MetaEndFrameCSV
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.CreateSelectCsvObjectMeta"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.CreateSelectCsvObjectMeta", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.CreateSelectCsvObjectMeta", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.CreateSelectCsvObjectMeta(key, csvMeta, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) CreateSelectJsonObjectMeta(ctx context.Context, key string, jsonMeta oss.JsonMetaRequest, options ...oss.Option) (oss.MetaEndFrameJSON, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.CreateSelectJsonObjectMeta")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.MetaEndFrameJSON
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.CreateSelectJsonObjectMeta"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.CreateSelectJsonObjectMeta", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.CreateSelectJsonObjectMeta", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.CreateSelectJsonObjectMeta(key, jsonMeta, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) DeleteLiveChannel(ctx context.Context, channelName string) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.DeleteLiveChannel")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.DeleteLiveChannel"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.DeleteLiveChannel", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.DeleteLiveChannel", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DeleteLiveChannel(channelName)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) DeleteObject(ctx context.Context, objectKey string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.DeleteObject")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.DeleteObject"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.DeleteObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.DeleteObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DeleteObject(objectKey, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) DeleteObjectTagging(ctx context.Context, objectKey string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.DeleteObjectTagging")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.DeleteObjectTagging"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.DeleteObjectTagging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.DeleteObjectTagging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DeleteObjectTagging(objectKey, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) DeleteObjectVersions(ctx context.Context, objectVersions []oss.DeleteObject, options ...oss.Option) (oss.DeleteObjectVersionsResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.DeleteObjectVersions")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.DeleteObjectVersionsResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.DeleteObjectVersions"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.DeleteObjectVersions", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.DeleteObjectVersions", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DeleteObjectVersions(objectVersions, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) DeleteObjects(ctx context.Context, objectKeys []string, options ...oss.Option) (oss.DeleteObjectsResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.DeleteObjects")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.DeleteObjectsResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.DeleteObjects"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.DeleteObjects", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.DeleteObjects", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DeleteObjects(objectKeys, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) Do(ctx context.Context, method string, objectName string, params map[string]interface{}, options []oss.Option, data io.Reader, listener oss.ProgressListener) (*oss.Response, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.Do")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 *oss.Response
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.Do"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(method, objectName, params, options, data, listener)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) DoAppendObject(ctx context.Context, request *oss.AppendObjectRequest, options []oss.Option) (*oss.AppendObjectResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.DoAppendObject")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 *oss.AppendObjectResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.DoAppendObject"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.DoAppendObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.DoAppendObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DoAppendObject(request, options)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) DoGetObject(ctx context.Context, request *oss.GetObjectRequest, options []oss.Option) (*oss.GetObjectResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.DoGetObject")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 *oss.GetObjectResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.DoGetObject"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.DoGetObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.DoGetObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DoGetObject(request, options)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) DoGetObjectWithURL(ctx context.Context, signedURL string, options []oss.Option) (*oss.GetObjectResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.DoGetObjectWithURL")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 *oss.GetObjectResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.DoGetObjectWithURL"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.DoGetObjectWithURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.DoGetObjectWithURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DoGetObjectWithURL(signedURL, options)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) DoPostSelectObject(ctx context.Context, key string, params map[string]interface{}, buf *bytes.Buffer, options ...oss.Option) (*oss.SelectObjectResponse, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.DoPostSelectObject")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 *oss.SelectObjectResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.DoPostSelectObject"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.DoPostSelectObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.DoPostSelectObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DoPostSelectObject(key, params, buf, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) DoPutObject(ctx context.Context, request *oss.PutObjectRequest, options []oss.Option) (*oss.Response, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.DoPutObject")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 *oss.Response
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.DoPutObject"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.DoPutObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.DoPutObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DoPutObject(request, options)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) DoPutObjectWithURL(ctx context.Context, signedURL string, reader io.Reader, options []oss.Option) (*oss.Response, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.DoPutObjectWithURL")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 *oss.Response
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.DoPutObjectWithURL"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.DoPutObjectWithURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.DoPutObjectWithURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DoPutObjectWithURL(signedURL, reader, options)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) DoUploadPart(ctx context.Context, request *oss.UploadPartRequest, options []oss.Option) (*oss.UploadPartResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.DoUploadPart")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 *oss.UploadPartResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.DoUploadPart"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.DoUploadPart", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.DoUploadPart", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DoUploadPart(request, options)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) DownloadFile(ctx context.Context, objectKey string, filePath string, partSize int64, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.DownloadFile")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.DownloadFile"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.DownloadFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.DownloadFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DownloadFile(objectKey, filePath, partSize, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) GetConfig(ctx context.Context) *oss.Config {
	res0 := w.obj.GetConfig()
	return res0
}

func (w *OSSBucketWrapper) GetLiveChannelHistory(ctx context.Context, channelName string) (oss.LiveChannelHistory, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.GetLiveChannelHistory")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.LiveChannelHistory
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.GetLiveChannelHistory"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.GetLiveChannelHistory", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.GetLiveChannelHistory", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetLiveChannelHistory(channelName)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) GetLiveChannelInfo(ctx context.Context, channelName string) (oss.LiveChannelConfiguration, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.GetLiveChannelInfo")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.LiveChannelConfiguration
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.GetLiveChannelInfo"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.GetLiveChannelInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.GetLiveChannelInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetLiveChannelInfo(channelName)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) GetLiveChannelStat(ctx context.Context, channelName string) (oss.LiveChannelStat, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.GetLiveChannelStat")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.LiveChannelStat
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.GetLiveChannelStat"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.GetLiveChannelStat", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.GetLiveChannelStat", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetLiveChannelStat(channelName)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) GetObject(ctx context.Context, objectKey string, options ...oss.Option) (io.ReadCloser, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.GetObject")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 io.ReadCloser
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.GetObject"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.GetObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.GetObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetObject(objectKey, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) GetObjectACL(ctx context.Context, objectKey string, options ...oss.Option) (oss.GetObjectACLResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.GetObjectACL")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.GetObjectACLResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.GetObjectACL"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.GetObjectACL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.GetObjectACL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetObjectACL(objectKey, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) GetObjectDetailedMeta(ctx context.Context, objectKey string, options ...oss.Option) (http.Header, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.GetObjectDetailedMeta")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 http.Header
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.GetObjectDetailedMeta"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.GetObjectDetailedMeta", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.GetObjectDetailedMeta", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetObjectDetailedMeta(objectKey, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) GetObjectMeta(ctx context.Context, objectKey string, options ...oss.Option) (http.Header, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.GetObjectMeta")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 http.Header
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.GetObjectMeta"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.GetObjectMeta", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.GetObjectMeta", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetObjectMeta(objectKey, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) GetObjectTagging(ctx context.Context, objectKey string, options ...oss.Option) (oss.GetObjectTaggingResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.GetObjectTagging")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.GetObjectTaggingResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.GetObjectTagging"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.GetObjectTagging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.GetObjectTagging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetObjectTagging(objectKey, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) GetObjectToFile(ctx context.Context, objectKey string, filePath string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.GetObjectToFile")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.GetObjectToFile"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.GetObjectToFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.GetObjectToFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.GetObjectToFile(objectKey, filePath, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) GetObjectToFileWithURL(ctx context.Context, signedURL string, filePath string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.GetObjectToFileWithURL")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.GetObjectToFileWithURL"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.GetObjectToFileWithURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.GetObjectToFileWithURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.GetObjectToFileWithURL(signedURL, filePath, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) GetObjectWithURL(ctx context.Context, signedURL string, options ...oss.Option) (io.ReadCloser, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.GetObjectWithURL")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 io.ReadCloser
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.GetObjectWithURL"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.GetObjectWithURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.GetObjectWithURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetObjectWithURL(signedURL, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) GetSymlink(ctx context.Context, objectKey string, options ...oss.Option) (http.Header, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.GetSymlink")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 http.Header
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.GetSymlink"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.GetSymlink", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.GetSymlink", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetSymlink(objectKey, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) GetVodPlaylist(ctx context.Context, channelName string, startTime time.Time, endTime time.Time) (io.ReadCloser, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.GetVodPlaylist")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 io.ReadCloser
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.GetVodPlaylist"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.GetVodPlaylist", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.GetVodPlaylist", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetVodPlaylist(channelName, startTime, endTime)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) InitiateMultipartUpload(ctx context.Context, objectKey string, options ...oss.Option) (oss.InitiateMultipartUploadResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.InitiateMultipartUpload")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.InitiateMultipartUploadResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.InitiateMultipartUpload"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.InitiateMultipartUpload", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.InitiateMultipartUpload", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.InitiateMultipartUpload(objectKey, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) IsObjectExist(ctx context.Context, objectKey string, options ...oss.Option) (bool, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.IsObjectExist")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 bool
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.IsObjectExist"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.IsObjectExist", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.IsObjectExist", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.IsObjectExist(objectKey, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) ListLiveChannel(ctx context.Context, options ...oss.Option) (oss.ListLiveChannelResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.ListLiveChannel")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.ListLiveChannelResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.ListLiveChannel"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.ListLiveChannel", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.ListLiveChannel", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListLiveChannel(options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) ListMultipartUploads(ctx context.Context, options ...oss.Option) (oss.ListMultipartUploadResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.ListMultipartUploads")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.ListMultipartUploadResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.ListMultipartUploads"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.ListMultipartUploads", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.ListMultipartUploads", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListMultipartUploads(options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) ListObjectVersions(ctx context.Context, options ...oss.Option) (oss.ListObjectVersionsResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.ListObjectVersions")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.ListObjectVersionsResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.ListObjectVersions"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.ListObjectVersions", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.ListObjectVersions", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListObjectVersions(options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) ListObjects(ctx context.Context, options ...oss.Option) (oss.ListObjectsResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.ListObjects")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.ListObjectsResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.ListObjects"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.ListObjects", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.ListObjects", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListObjects(options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) ListObjectsV2(ctx context.Context, options ...oss.Option) (oss.ListObjectsResultV2, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.ListObjectsV2")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.ListObjectsResultV2
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.ListObjectsV2"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.ListObjectsV2", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.ListObjectsV2", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListObjectsV2(options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) ListUploadedParts(ctx context.Context, imur oss.InitiateMultipartUploadResult, options ...oss.Option) (oss.ListUploadedPartsResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.ListUploadedParts")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.ListUploadedPartsResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.ListUploadedParts"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.ListUploadedParts", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.ListUploadedParts", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListUploadedParts(imur, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) OptionsMethod(ctx context.Context, objectKey string, options ...oss.Option) (http.Header, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.OptionsMethod")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 http.Header
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.OptionsMethod"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.OptionsMethod", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.OptionsMethod", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.OptionsMethod(objectKey, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) PostVodPlaylist(ctx context.Context, channelName string, playlistName string, startTime time.Time, endTime time.Time) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.PostVodPlaylist")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.PostVodPlaylist"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.PostVodPlaylist", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.PostVodPlaylist", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.PostVodPlaylist(channelName, playlistName, startTime, endTime)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) ProcessObject(ctx context.Context, objectKey string, process string, options ...oss.Option) (oss.ProcessObjectResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.ProcessObject")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.ProcessObjectResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.ProcessObject"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.ProcessObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.ProcessObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ProcessObject(objectKey, process, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) PutLiveChannelStatus(ctx context.Context, channelName string, status string) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.PutLiveChannelStatus")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.PutLiveChannelStatus"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.PutLiveChannelStatus", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.PutLiveChannelStatus", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.PutLiveChannelStatus(channelName, status)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) PutObject(ctx context.Context, objectKey string, reader io.Reader, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.PutObject")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.PutObject"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.PutObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.PutObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.PutObject(objectKey, reader, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) PutObjectFromFile(ctx context.Context, objectKey string, filePath string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.PutObjectFromFile")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.PutObjectFromFile"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.PutObjectFromFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.PutObjectFromFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.PutObjectFromFile(objectKey, filePath, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) PutObjectFromFileWithURL(ctx context.Context, signedURL string, filePath string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.PutObjectFromFileWithURL")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.PutObjectFromFileWithURL"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.PutObjectFromFileWithURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.PutObjectFromFileWithURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.PutObjectFromFileWithURL(signedURL, filePath, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) PutObjectTagging(ctx context.Context, objectKey string, tagging oss.Tagging, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.PutObjectTagging")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.PutObjectTagging"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.PutObjectTagging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.PutObjectTagging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.PutObjectTagging(objectKey, tagging, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) PutObjectWithURL(ctx context.Context, signedURL string, reader io.Reader, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.PutObjectWithURL")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.PutObjectWithURL"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.PutObjectWithURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.PutObjectWithURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.PutObjectWithURL(signedURL, reader, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) PutSymlink(ctx context.Context, symObjectKey string, targetObjectKey string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.PutSymlink")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.PutSymlink"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.PutSymlink", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.PutSymlink", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.PutSymlink(symObjectKey, targetObjectKey, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) RestoreObject(ctx context.Context, objectKey string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.RestoreObject")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.RestoreObject"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.RestoreObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.RestoreObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.RestoreObject(objectKey, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) RestoreObjectDetail(ctx context.Context, objectKey string, restoreConfig oss.RestoreConfiguration, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.RestoreObjectDetail")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.RestoreObjectDetail"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.RestoreObjectDetail", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.RestoreObjectDetail", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.RestoreObjectDetail(objectKey, restoreConfig, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) RestoreObjectXML(ctx context.Context, objectKey string, configXML string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.RestoreObjectXML")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.RestoreObjectXML"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.RestoreObjectXML", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.RestoreObjectXML", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.RestoreObjectXML(objectKey, configXML, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) SelectObject(ctx context.Context, key string, selectReq oss.SelectRequest, options ...oss.Option) (io.ReadCloser, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.SelectObject")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 io.ReadCloser
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.SelectObject"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.SelectObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.SelectObject", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.SelectObject(key, selectReq, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) SelectObjectIntoFile(ctx context.Context, key string, fileName string, selectReq oss.SelectRequest, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.SelectObjectIntoFile")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.SelectObjectIntoFile"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.SelectObjectIntoFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.SelectObjectIntoFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SelectObjectIntoFile(key, fileName, selectReq, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) SetObjectACL(ctx context.Context, objectKey string, objectACL oss.ACLType, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.SetObjectACL")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.SetObjectACL"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.SetObjectACL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.SetObjectACL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetObjectACL(objectKey, objectACL, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) SetObjectMeta(ctx context.Context, objectKey string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.SetObjectMeta")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.SetObjectMeta"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.SetObjectMeta", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.SetObjectMeta", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetObjectMeta(objectKey, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) SignRtmpURL(ctx context.Context, channelName string, playlistName string, expires int64) (string, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.SignRtmpURL")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 string
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.SignRtmpURL"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.SignRtmpURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.SignRtmpURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.SignRtmpURL(channelName, playlistName, expires)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) SignURL(ctx context.Context, objectKey string, method oss.HTTPMethod, expiredInSec int64, options ...oss.Option) (string, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.SignURL")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 string
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.SignURL"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.SignURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.SignURL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.SignURL(objectKey, method, expiredInSec, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) UploadFile(ctx context.Context, objectKey string, filePath string, partSize int64, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.UploadFile")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.UploadFile"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.UploadFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.UploadFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.UploadFile(objectKey, filePath, partSize, options...)
		return err
	})
	return err
}

func (w *OSSBucketWrapper) UploadPart(ctx context.Context, imur oss.InitiateMultipartUploadResult, reader io.Reader, partSize int64, partNumber int, options ...oss.Option) (oss.UploadPart, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.UploadPart")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.UploadPart
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.UploadPart"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.UploadPart", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.UploadPart", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.UploadPart(imur, reader, partSize, partNumber, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) UploadPartCopy(ctx context.Context, imur oss.InitiateMultipartUploadResult, srcBucketName string, srcObjectKey string, startPosition int64, partSize int64, partNumber int, options ...oss.Option) (oss.UploadPart, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.UploadPartCopy")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.UploadPart
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.UploadPartCopy"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.UploadPartCopy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.UploadPartCopy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.UploadPartCopy(imur, srcBucketName, srcObjectKey, startPosition, partSize, partNumber, options...)
		return err
	})
	return res0, err
}

func (w *OSSBucketWrapper) UploadPartFromFile(ctx context.Context, imur oss.InitiateMultipartUploadResult, filePath string, startPosition int64, partSize int64, partNumber int, options ...oss.Option) (oss.UploadPart, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Bucket.UploadPartFromFile")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.UploadPart
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Bucket.UploadPartFromFile"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Bucket.UploadPartFromFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Bucket.UploadPartFromFile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.UploadPartFromFile(imur, filePath, startPosition, partSize, partNumber, options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) AbortBucketWorm(ctx context.Context, bucketName string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.AbortBucketWorm")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.AbortBucketWorm"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.AbortBucketWorm", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.AbortBucketWorm", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.AbortBucketWorm(bucketName, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) Bucket(bucketName string) (*OSSBucketWrapper, error) {
	var res0 *oss.Bucket
	var err error
	res0, err = w.obj.Bucket(bucketName)
	return &OSSBucketWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric, rateLimiterGroup: w.rateLimiterGroup}, err
}

func (w *OSSClientWrapper) CompleteBucketWorm(ctx context.Context, bucketName string, wormID string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.CompleteBucketWorm")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.CompleteBucketWorm"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.CompleteBucketWorm", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.CompleteBucketWorm", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.CompleteBucketWorm(bucketName, wormID, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) CreateBucket(ctx context.Context, bucketName string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.CreateBucket")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.CreateBucket"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.CreateBucket", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.CreateBucket", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.CreateBucket(bucketName, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) DeleteBucket(ctx context.Context, bucketName string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.DeleteBucket")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DeleteBucket"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.DeleteBucket", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.DeleteBucket", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DeleteBucket(bucketName, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) DeleteBucketCORS(ctx context.Context, bucketName string) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.DeleteBucketCORS")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DeleteBucketCORS"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.DeleteBucketCORS", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.DeleteBucketCORS", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DeleteBucketCORS(bucketName)
		return err
	})
	return err
}

func (w *OSSClientWrapper) DeleteBucketEncryption(ctx context.Context, bucketName string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.DeleteBucketEncryption")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DeleteBucketEncryption"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.DeleteBucketEncryption", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.DeleteBucketEncryption", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DeleteBucketEncryption(bucketName, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) DeleteBucketInventory(ctx context.Context, bucketName string, strInventoryId string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.DeleteBucketInventory")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DeleteBucketInventory"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.DeleteBucketInventory", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.DeleteBucketInventory", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DeleteBucketInventory(bucketName, strInventoryId, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) DeleteBucketLifecycle(ctx context.Context, bucketName string) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.DeleteBucketLifecycle")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DeleteBucketLifecycle"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.DeleteBucketLifecycle", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.DeleteBucketLifecycle", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DeleteBucketLifecycle(bucketName)
		return err
	})
	return err
}

func (w *OSSClientWrapper) DeleteBucketLogging(ctx context.Context, bucketName string) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.DeleteBucketLogging")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DeleteBucketLogging"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.DeleteBucketLogging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.DeleteBucketLogging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DeleteBucketLogging(bucketName)
		return err
	})
	return err
}

func (w *OSSClientWrapper) DeleteBucketPolicy(ctx context.Context, bucketName string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.DeleteBucketPolicy")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DeleteBucketPolicy"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.DeleteBucketPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.DeleteBucketPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DeleteBucketPolicy(bucketName, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) DeleteBucketQosInfo(ctx context.Context, bucketName string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.DeleteBucketQosInfo")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DeleteBucketQosInfo"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.DeleteBucketQosInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.DeleteBucketQosInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DeleteBucketQosInfo(bucketName, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) DeleteBucketTagging(ctx context.Context, bucketName string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.DeleteBucketTagging")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DeleteBucketTagging"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.DeleteBucketTagging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.DeleteBucketTagging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DeleteBucketTagging(bucketName, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) DeleteBucketWebsite(ctx context.Context, bucketName string) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.DeleteBucketWebsite")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.DeleteBucketWebsite"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.DeleteBucketWebsite", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.DeleteBucketWebsite", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DeleteBucketWebsite(bucketName)
		return err
	})
	return err
}

func (w *OSSClientWrapper) ExtendBucketWorm(ctx context.Context, bucketName string, retentionDays int, wormID string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.ExtendBucketWorm")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ExtendBucketWorm"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.ExtendBucketWorm", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.ExtendBucketWorm", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.ExtendBucketWorm(bucketName, retentionDays, wormID, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) GetBucketACL(ctx context.Context, bucketName string) (oss.GetBucketACLResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketACL")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.GetBucketACLResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketACL"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketACL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketACL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketACL(bucketName)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketAsyncTask(ctx context.Context, bucketName string, taskID string, options ...oss.Option) (oss.AsynFetchTaskInfo, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketAsyncTask")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.AsynFetchTaskInfo
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketAsyncTask"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketAsyncTask", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketAsyncTask", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketAsyncTask(bucketName, taskID, options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketCORS(ctx context.Context, bucketName string) (oss.GetBucketCORSResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketCORS")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.GetBucketCORSResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketCORS"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketCORS", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketCORS", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketCORS(bucketName)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketEncryption(ctx context.Context, bucketName string, options ...oss.Option) (oss.GetBucketEncryptionResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketEncryption")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.GetBucketEncryptionResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketEncryption"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketEncryption", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketEncryption", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketEncryption(bucketName, options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketInfo(ctx context.Context, bucketName string, options ...oss.Option) (oss.GetBucketInfoResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketInfo")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.GetBucketInfoResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketInfo"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketInfo(bucketName, options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketInventory(ctx context.Context, bucketName string, strInventoryId string, options ...oss.Option) (oss.InventoryConfiguration, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketInventory")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.InventoryConfiguration
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketInventory"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketInventory", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketInventory", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketInventory(bucketName, strInventoryId, options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketLifecycle(ctx context.Context, bucketName string) (oss.GetBucketLifecycleResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketLifecycle")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.GetBucketLifecycleResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketLifecycle"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketLifecycle", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketLifecycle", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketLifecycle(bucketName)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketLocation(ctx context.Context, bucketName string) (string, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketLocation")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 string
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketLocation"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketLocation", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketLocation", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketLocation(bucketName)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketLogging(ctx context.Context, bucketName string) (oss.GetBucketLoggingResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketLogging")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.GetBucketLoggingResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketLogging"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketLogging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketLogging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketLogging(bucketName)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketPolicy(ctx context.Context, bucketName string, options ...oss.Option) (string, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketPolicy")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 string
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketPolicy"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketPolicy(bucketName, options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketQosInfo(ctx context.Context, bucketName string, options ...oss.Option) (oss.BucketQoSConfiguration, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketQosInfo")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.BucketQoSConfiguration
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketQosInfo"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketQosInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketQosInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketQosInfo(bucketName, options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketReferer(ctx context.Context, bucketName string) (oss.GetBucketRefererResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketReferer")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.GetBucketRefererResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketReferer"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketReferer", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketReferer", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketReferer(bucketName)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketRequestPayment(ctx context.Context, bucketName string, options ...oss.Option) (oss.RequestPaymentConfiguration, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketRequestPayment")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.RequestPaymentConfiguration
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketRequestPayment"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketRequestPayment", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketRequestPayment", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketRequestPayment(bucketName, options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketStat(ctx context.Context, bucketName string) (oss.GetBucketStatResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketStat")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.GetBucketStatResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketStat"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketStat", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketStat", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketStat(bucketName)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketTagging(ctx context.Context, bucketName string, options ...oss.Option) (oss.GetBucketTaggingResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketTagging")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.GetBucketTaggingResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketTagging"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketTagging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketTagging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketTagging(bucketName, options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketVersioning(ctx context.Context, bucketName string, options ...oss.Option) (oss.GetBucketVersioningResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketVersioning")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.GetBucketVersioningResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketVersioning"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketVersioning", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketVersioning", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketVersioning(bucketName, options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketWebsite(ctx context.Context, bucketName string) (oss.GetBucketWebsiteResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketWebsite")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.GetBucketWebsiteResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketWebsite"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketWebsite", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketWebsite", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketWebsite(bucketName)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetBucketWorm(ctx context.Context, bucketName string, options ...oss.Option) (oss.WormConfiguration, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetBucketWorm")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.WormConfiguration
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetBucketWorm"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetBucketWorm", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetBucketWorm", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetBucketWorm(bucketName, options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) GetUserQoSInfo(ctx context.Context, options ...oss.Option) (oss.UserQoSConfiguration, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.GetUserQoSInfo")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.UserQoSConfiguration
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.GetUserQoSInfo"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.GetUserQoSInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.GetUserQoSInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetUserQoSInfo(options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) InitiateBucketWorm(ctx context.Context, bucketName string, retentionDays int, options ...oss.Option) (string, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.InitiateBucketWorm")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 string
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.InitiateBucketWorm"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.InitiateBucketWorm", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.InitiateBucketWorm", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.InitiateBucketWorm(bucketName, retentionDays, options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) IsBucketExist(ctx context.Context, bucketName string) (bool, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.IsBucketExist")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 bool
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.IsBucketExist"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.IsBucketExist", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.IsBucketExist", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.IsBucketExist(bucketName)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) LimitUploadSpeed(ctx context.Context, upSpeed int) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.LimitUploadSpeed")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.LimitUploadSpeed"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.LimitUploadSpeed", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.LimitUploadSpeed", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.LimitUploadSpeed(upSpeed)
		return err
	})
	return err
}

func (w *OSSClientWrapper) ListBucketInventory(ctx context.Context, bucketName string, continuationToken string, options ...oss.Option) (oss.ListInventoryConfigurationsResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.ListBucketInventory")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.ListInventoryConfigurationsResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ListBucketInventory"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.ListBucketInventory", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.ListBucketInventory", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListBucketInventory(bucketName, continuationToken, options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) ListBuckets(ctx context.Context, options ...oss.Option) (oss.ListBucketsResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.ListBuckets")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.ListBucketsResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.ListBuckets"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.ListBuckets", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.ListBuckets", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.ListBuckets(options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) SetBucketACL(ctx context.Context, bucketName string, bucketACL oss.ACLType) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketACL")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketACL"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketACL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketACL", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketACL(bucketName, bucketACL)
		return err
	})
	return err
}

func (w *OSSClientWrapper) SetBucketAsyncTask(ctx context.Context, bucketName string, asynConf oss.AsyncFetchTaskConfiguration, options ...oss.Option) (oss.AsyncFetchTaskResult, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketAsyncTask")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 oss.AsyncFetchTaskResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketAsyncTask"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketAsyncTask", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketAsyncTask", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.SetBucketAsyncTask(bucketName, asynConf, options...)
		return err
	})
	return res0, err
}

func (w *OSSClientWrapper) SetBucketCORS(ctx context.Context, bucketName string, corsRules []oss.CORSRule) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketCORS")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketCORS"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketCORS", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketCORS", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketCORS(bucketName, corsRules)
		return err
	})
	return err
}

func (w *OSSClientWrapper) SetBucketEncryption(ctx context.Context, bucketName string, encryptionRule oss.ServerEncryptionRule, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketEncryption")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketEncryption"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketEncryption", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketEncryption", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketEncryption(bucketName, encryptionRule, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) SetBucketInventory(ctx context.Context, bucketName string, inventoryConfig oss.InventoryConfiguration, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketInventory")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketInventory"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketInventory", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketInventory", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketInventory(bucketName, inventoryConfig, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) SetBucketLifecycle(ctx context.Context, bucketName string, rules []oss.LifecycleRule) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketLifecycle")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketLifecycle"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketLifecycle", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketLifecycle", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketLifecycle(bucketName, rules)
		return err
	})
	return err
}

func (w *OSSClientWrapper) SetBucketLogging(ctx context.Context, bucketName string, targetBucket string, targetPrefix string, isEnable bool) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketLogging")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketLogging"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketLogging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketLogging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketLogging(bucketName, targetBucket, targetPrefix, isEnable)
		return err
	})
	return err
}

func (w *OSSClientWrapper) SetBucketPolicy(ctx context.Context, bucketName string, policy string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketPolicy")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketPolicy"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketPolicy(bucketName, policy, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) SetBucketQoSInfo(ctx context.Context, bucketName string, qosConf oss.BucketQoSConfiguration, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketQoSInfo")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketQoSInfo"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketQoSInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketQoSInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketQoSInfo(bucketName, qosConf, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) SetBucketReferer(ctx context.Context, bucketName string, referers []string, allowEmptyReferer bool) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketReferer")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketReferer"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketReferer", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketReferer", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketReferer(bucketName, referers, allowEmptyReferer)
		return err
	})
	return err
}

func (w *OSSClientWrapper) SetBucketRequestPayment(ctx context.Context, bucketName string, paymentConfig oss.RequestPaymentConfiguration, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketRequestPayment")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketRequestPayment"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketRequestPayment", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketRequestPayment", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketRequestPayment(bucketName, paymentConfig, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) SetBucketTagging(ctx context.Context, bucketName string, tagging oss.Tagging, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketTagging")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketTagging"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketTagging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketTagging", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketTagging(bucketName, tagging, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) SetBucketVersioning(ctx context.Context, bucketName string, versioningConfig oss.VersioningConfig, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketVersioning")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketVersioning"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketVersioning", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketVersioning", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketVersioning(bucketName, versioningConfig, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) SetBucketWebsite(ctx context.Context, bucketName string, indexDocument string, errorDocument string) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketWebsite")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketWebsite"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketWebsite", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketWebsite", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketWebsite(bucketName, indexDocument, errorDocument)
		return err
	})
	return err
}

func (w *OSSClientWrapper) SetBucketWebsiteDetail(ctx context.Context, bucketName string, wxml oss.WebsiteXML, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketWebsiteDetail")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketWebsiteDetail"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketWebsiteDetail", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketWebsiteDetail", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketWebsiteDetail(bucketName, wxml, options...)
		return err
	})
	return err
}

func (w *OSSClientWrapper) SetBucketWebsiteXml(ctx context.Context, bucketName string, webXml string, options ...oss.Option) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "oss.Client.SetBucketWebsiteXml")
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
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.SetBucketWebsiteXml"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("oss.Client.SetBucketWebsiteXml", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("oss.Client.SetBucketWebsiteXml", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetBucketWebsiteXml(bucketName, webXml, options...)
		return err
	})
	return err
}
