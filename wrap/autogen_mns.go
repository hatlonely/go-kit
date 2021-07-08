// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"fmt"
	"time"

	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

func NewMNSClientWrapper(
	obj ali_mns.MNSClient,
	retry *micro.Retry,
	options *WrapperOptions,
	durationMetric *prometheus.HistogramVec,
	inflightMetric *prometheus.GaugeVec,
	rateLimiter micro.RateLimiter,
	parallelController micro.ParallelController) *MNSClientWrapper {
	return &MNSClientWrapper{
		obj:                obj,
		retry:              retry,
		options:            options,
		durationMetric:     durationMetric,
		inflightMetric:     inflightMetric,
		rateLimiter:        rateLimiter,
		parallelController: parallelController,
	}
}

type MNSClientWrapper struct {
	obj                ali_mns.MNSClient
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *MNSClientWrapper) Unwrap() ali_mns.MNSClient {
	return w.obj
}

func (w *MNSClientWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *MNSClientWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *MNSClientWrapper) OnRateLimiterChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *MNSClientWrapper) OnParallelControllerChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *MNSClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        fmt.Sprintf("%s_ali_mns_MNSClient_durationMs", options.Name),
		Help:        "ali_mns MNSClient response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        fmt.Sprintf("%s_ali_mns_MNSClient_inflight", options.Name),
		Help:        "ali_mns MNSClient inflight",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "custom"})
}

type MNSQueueWrapper struct {
	obj                *ali_mns.MNSQueue
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *MNSQueueWrapper) Unwrap() *ali_mns.MNSQueue {
	return w.obj
}

type MNSTopicWrapper struct {
	obj                *ali_mns.MNSTopic
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *MNSTopicWrapper) Unwrap() *ali_mns.MNSTopic {
	return w.obj
}

func (w *MNSQueueWrapper) BatchDeleteMessage(ctx context.Context, receiptHandles ...string) (ali_mns.BatchMessageDeleteErrorResponse, error) {
	ctxOptions := FromContext(ctx)
	var resp ali_mns.BatchMessageDeleteErrorResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MNSQueue.BatchDeleteMessage", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MNSQueue.BatchDeleteMessage", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MNSQueue.BatchDeleteMessage", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ali_mns.MNSQueue.BatchDeleteMessage", ctxOptions.StartSpanOpts...)
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
			w.inflightMetric.WithLabelValues("ali_mns.MNSQueue.BatchDeleteMessage", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ali_mns.MNSQueue.BatchDeleteMessage", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ali_mns.MNSQueue.BatchDeleteMessage", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		resp, err = w.obj.BatchDeleteMessage(receiptHandles...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return resp, err
}

func (w *MNSQueueWrapper) BatchPeekMessage(respChan chan ali_mns.BatchMessageReceiveResponse, errChan chan error, numOfMessages int32) {
	w.obj.BatchPeekMessage(respChan, errChan, numOfMessages)
}

func (w *MNSQueueWrapper) BatchReceiveMessage(respChan chan ali_mns.BatchMessageReceiveResponse, errChan chan error, numOfMessages int32, waitseconds ...int64) {
	w.obj.BatchReceiveMessage(respChan, errChan, numOfMessages, waitseconds...)
}

func (w *MNSQueueWrapper) BatchSendMessage(ctx context.Context, messages ...ali_mns.MessageSendRequest) (ali_mns.BatchMessageSendResponse, error) {
	ctxOptions := FromContext(ctx)
	var resp ali_mns.BatchMessageSendResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MNSQueue.BatchSendMessage", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MNSQueue.BatchSendMessage", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MNSQueue.BatchSendMessage", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ali_mns.MNSQueue.BatchSendMessage", ctxOptions.StartSpanOpts...)
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
			w.inflightMetric.WithLabelValues("ali_mns.MNSQueue.BatchSendMessage", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ali_mns.MNSQueue.BatchSendMessage", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ali_mns.MNSQueue.BatchSendMessage", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		resp, err = w.obj.BatchSendMessage(messages...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return resp, err
}

func (w *MNSQueueWrapper) ChangeMessageVisibility(ctx context.Context, receiptHandle string, visibilityTimeout int64) (ali_mns.MessageVisibilityChangeResponse, error) {
	ctxOptions := FromContext(ctx)
	var resp ali_mns.MessageVisibilityChangeResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MNSQueue.ChangeMessageVisibility", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MNSQueue.ChangeMessageVisibility", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MNSQueue.ChangeMessageVisibility", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ali_mns.MNSQueue.ChangeMessageVisibility", ctxOptions.StartSpanOpts...)
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
			w.inflightMetric.WithLabelValues("ali_mns.MNSQueue.ChangeMessageVisibility", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ali_mns.MNSQueue.ChangeMessageVisibility", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ali_mns.MNSQueue.ChangeMessageVisibility", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		resp, err = w.obj.ChangeMessageVisibility(receiptHandle, visibilityTimeout)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return resp, err
}

func (w *MNSQueueWrapper) DeleteMessage(ctx context.Context, receiptHandle string) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MNSQueue.DeleteMessage", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MNSQueue.DeleteMessage", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MNSQueue.DeleteMessage", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ali_mns.MNSQueue.DeleteMessage", ctxOptions.StartSpanOpts...)
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
			w.inflightMetric.WithLabelValues("ali_mns.MNSQueue.DeleteMessage", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ali_mns.MNSQueue.DeleteMessage", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ali_mns.MNSQueue.DeleteMessage", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.DeleteMessage(receiptHandle)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *MNSQueueWrapper) Name() string {
	res0 := w.obj.Name()
	return res0
}

func (w *MNSQueueWrapper) PeekMessage(respChan chan ali_mns.MessageReceiveResponse, errChan chan error) {
	w.obj.PeekMessage(respChan, errChan)
}

func (w *MNSQueueWrapper) QPSMonitor() *ali_mns.QPSMonitor {
	res0 := w.obj.QPSMonitor()
	return res0
}

func (w *MNSQueueWrapper) ReceiveMessage(respChan chan ali_mns.MessageReceiveResponse, errChan chan error, waitseconds ...int64) {
	w.obj.ReceiveMessage(respChan, errChan, waitseconds...)
}

func (w *MNSQueueWrapper) SendMessage(ctx context.Context, message ali_mns.MessageSendRequest) (ali_mns.MessageSendResponse, error) {
	ctxOptions := FromContext(ctx)
	var resp ali_mns.MessageSendResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MNSQueue.SendMessage", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MNSQueue.SendMessage", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MNSQueue.SendMessage", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ali_mns.MNSQueue.SendMessage", ctxOptions.StartSpanOpts...)
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
			w.inflightMetric.WithLabelValues("ali_mns.MNSQueue.SendMessage", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ali_mns.MNSQueue.SendMessage", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ali_mns.MNSQueue.SendMessage", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		resp, err = w.obj.SendMessage(message)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return resp, err
}

func (w *MNSTopicWrapper) GenerateMailEndpoint(mailAddress string) string {
	res0 := w.obj.GenerateMailEndpoint(mailAddress)
	return res0
}

func (w *MNSTopicWrapper) GenerateQueueEndpoint(queueName string) string {
	res0 := w.obj.GenerateQueueEndpoint(queueName)
	return res0
}

func (w *MNSTopicWrapper) GetSubscriptionAttributes(ctx context.Context, subscriptionName string) (ali_mns.SubscriptionAttribute, error) {
	ctxOptions := FromContext(ctx)
	var attr ali_mns.SubscriptionAttribute
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MNSTopic.GetSubscriptionAttributes", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MNSTopic.GetSubscriptionAttributes", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MNSTopic.GetSubscriptionAttributes", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ali_mns.MNSTopic.GetSubscriptionAttributes", ctxOptions.StartSpanOpts...)
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
			w.inflightMetric.WithLabelValues("ali_mns.MNSTopic.GetSubscriptionAttributes", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ali_mns.MNSTopic.GetSubscriptionAttributes", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ali_mns.MNSTopic.GetSubscriptionAttributes", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		attr, err = w.obj.GetSubscriptionAttributes(subscriptionName)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return attr, err
}

func (w *MNSTopicWrapper) ListSubscriptionByTopic(ctx context.Context, nextMarker string, retNumber int32, prefix string) (ali_mns.Subscriptions, error) {
	ctxOptions := FromContext(ctx)
	var subscriptions ali_mns.Subscriptions
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MNSTopic.ListSubscriptionByTopic", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MNSTopic.ListSubscriptionByTopic", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MNSTopic.ListSubscriptionByTopic", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ali_mns.MNSTopic.ListSubscriptionByTopic", ctxOptions.StartSpanOpts...)
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
			w.inflightMetric.WithLabelValues("ali_mns.MNSTopic.ListSubscriptionByTopic", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ali_mns.MNSTopic.ListSubscriptionByTopic", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ali_mns.MNSTopic.ListSubscriptionByTopic", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		subscriptions, err = w.obj.ListSubscriptionByTopic(nextMarker, retNumber, prefix)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return subscriptions, err
}

func (w *MNSTopicWrapper) ListSubscriptionDetailByTopic(ctx context.Context, nextMarker string, retNumber int32, prefix string) (ali_mns.SubscriptionDetails, error) {
	ctxOptions := FromContext(ctx)
	var subscriptionDetails ali_mns.SubscriptionDetails
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MNSTopic.ListSubscriptionDetailByTopic", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MNSTopic.ListSubscriptionDetailByTopic", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MNSTopic.ListSubscriptionDetailByTopic", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ali_mns.MNSTopic.ListSubscriptionDetailByTopic", ctxOptions.StartSpanOpts...)
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
			w.inflightMetric.WithLabelValues("ali_mns.MNSTopic.ListSubscriptionDetailByTopic", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ali_mns.MNSTopic.ListSubscriptionDetailByTopic", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ali_mns.MNSTopic.ListSubscriptionDetailByTopic", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		subscriptionDetails, err = w.obj.ListSubscriptionDetailByTopic(nextMarker, retNumber, prefix)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return subscriptionDetails, err
}

func (w *MNSTopicWrapper) Name() string {
	res0 := w.obj.Name()
	return res0
}

func (w *MNSTopicWrapper) PublishMessage(ctx context.Context, message ali_mns.MessagePublishRequest) (ali_mns.MessageSendResponse, error) {
	ctxOptions := FromContext(ctx)
	var resp ali_mns.MessageSendResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MNSTopic.PublishMessage", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MNSTopic.PublishMessage", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MNSTopic.PublishMessage", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ali_mns.MNSTopic.PublishMessage", ctxOptions.StartSpanOpts...)
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
			w.inflightMetric.WithLabelValues("ali_mns.MNSTopic.PublishMessage", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ali_mns.MNSTopic.PublishMessage", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ali_mns.MNSTopic.PublishMessage", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		resp, err = w.obj.PublishMessage(message)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return resp, err
}

func (w *MNSTopicWrapper) SetSubscriptionAttributes(ctx context.Context, subscriptionName string, notifyStrategy ali_mns.NotifyStrategyType) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MNSTopic.SetSubscriptionAttributes", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MNSTopic.SetSubscriptionAttributes", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MNSTopic.SetSubscriptionAttributes", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ali_mns.MNSTopic.SetSubscriptionAttributes", ctxOptions.StartSpanOpts...)
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
			w.inflightMetric.WithLabelValues("ali_mns.MNSTopic.SetSubscriptionAttributes", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ali_mns.MNSTopic.SetSubscriptionAttributes", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ali_mns.MNSTopic.SetSubscriptionAttributes", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.SetSubscriptionAttributes(subscriptionName, notifyStrategy)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *MNSTopicWrapper) Subscribe(ctx context.Context, subscriptionName string, message ali_mns.MessageSubsribeRequest) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MNSTopic.Subscribe", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MNSTopic.Subscribe", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MNSTopic.Subscribe", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ali_mns.MNSTopic.Subscribe", ctxOptions.StartSpanOpts...)
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
			w.inflightMetric.WithLabelValues("ali_mns.MNSTopic.Subscribe", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ali_mns.MNSTopic.Subscribe", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ali_mns.MNSTopic.Subscribe", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Subscribe(subscriptionName, message)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *MNSTopicWrapper) Unsubscribe(ctx context.Context, subscriptionName string) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MNSTopic.Unsubscribe", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MNSTopic.Unsubscribe", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MNSTopic.Unsubscribe", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ali_mns.MNSTopic.Unsubscribe", ctxOptions.StartSpanOpts...)
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
			w.inflightMetric.WithLabelValues("ali_mns.MNSTopic.Unsubscribe", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ali_mns.MNSTopic.Unsubscribe", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ali_mns.MNSTopic.Unsubscribe", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Unsubscribe(subscriptionName)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}
