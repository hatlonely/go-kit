// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type RedisClientWrapper struct {
	obj              *redis.Client
	retry            *Retry
	options          *WrapperOptions
	durationMetric   *prometheus.HistogramVec
	totalMetric      *prometheus.CounterVec
	rateLimiterGroup RateLimiterGroup
}

func (w *RedisClientWrapper) Unwrap() *redis.Client {
	return w.obj
}

func (w *RedisClientWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *RedisClientWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *RedisClientWrapper) OnRateLimiterGroupChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options RateLimiterGroupOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		rateLimiterGroup, err := NewRateLimiterGroup(&options)
		if err != nil {
			return errors.Wrap(err, "NewRateLimiterGroup failed")
		}
		w.rateLimiterGroup = rateLimiterGroup
		return nil
	}
}

func (w *RedisClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "redis_Client_durationMs",
		Help:        "redis Client response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "redis_Client_total",
		Help:        "redis Client request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
}

type RedisClusterClientWrapper struct {
	obj              *redis.ClusterClient
	retry            *Retry
	options          *WrapperOptions
	durationMetric   *prometheus.HistogramVec
	totalMetric      *prometheus.CounterVec
	rateLimiterGroup RateLimiterGroup
}

func (w *RedisClusterClientWrapper) Unwrap() *redis.ClusterClient {
	return w.obj
}

func (w *RedisClusterClientWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *RedisClusterClientWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *RedisClusterClientWrapper) OnRateLimiterGroupChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options RateLimiterGroupOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		rateLimiterGroup, err := NewRateLimiterGroup(&options)
		if err != nil {
			return errors.Wrap(err, "NewRateLimiterGroup failed")
		}
		w.rateLimiterGroup = rateLimiterGroup
		return nil
	}
}

func (w *RedisClusterClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "redis_ClusterClient_durationMs",
		Help:        "redis ClusterClient response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "redis_ClusterClient_total",
		Help:        "redis ClusterClient request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
}

func (w *RedisClientWrapper) Context() context.Context {
	res0 := w.obj.Context()
	return res0
}

func (w *RedisClientWrapper) Options() *redis.Options {
	res0 := w.obj.Options()
	return res0
}

func (w *RedisClientWrapper) PSubscribe(channels ...string) *redis.PubSub {
	res0 := w.obj.PSubscribe(channels...)
	return res0
}

func (w *RedisClientWrapper) Pipeline() redis.Pipeliner {
	res0 := w.obj.Pipeline()
	return res0
}

func (w *RedisClientWrapper) Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.Client.Pipelined")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 []redis.Cmder
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.Pipelined"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("redis.Client.Pipelined", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("redis.Client.Pipelined", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Pipelined(fn)
		return err
	})
	return res0, err
}

func (w *RedisClientWrapper) PoolStats() *redis.PoolStats {
	res0 := w.obj.PoolStats()
	return res0
}

func (w *RedisClientWrapper) SetLimiter(l redis.Limiter) *RedisClientWrapper {
	w.obj = w.obj.SetLimiter(l)
	return w
}

func (w *RedisClientWrapper) Subscribe(channels ...string) *redis.PubSub {
	res0 := w.obj.Subscribe(channels...)
	return res0
}

func (w *RedisClientWrapper) TxPipeline() redis.Pipeliner {
	res0 := w.obj.TxPipeline()
	return res0
}

func (w *RedisClientWrapper) TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.Client.TxPipelined")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 []redis.Cmder
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "Client.TxPipelined"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("redis.Client.TxPipelined", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("redis.Client.TxPipelined", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.TxPipelined(fn)
		return err
	})
	return res0, err
}

func (w *RedisClientWrapper) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.Client.Watch")
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
			if err := w.rateLimiterGroup.Wait(ctx, "Client.Watch"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("redis.Client.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("redis.Client.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Watch(fn, keys...)
		return err
	})
	return err
}

func (w *RedisClientWrapper) WithContext(ctx context.Context) *RedisClientWrapper {
	w.obj = w.obj.WithContext(ctx)
	return w
}

func (w *RedisClusterClientWrapper) Close(ctx context.Context) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Close")
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
			if err := w.rateLimiterGroup.Wait(ctx, "ClusterClient.Close"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("redis.ClusterClient.Close", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Close", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Close()
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) Context() context.Context {
	res0 := w.obj.Context()
	return res0
}

func (w *RedisClusterClientWrapper) DBSize() *redis.IntCmd {
	res0 := w.obj.DBSize()
	return res0
}

func (w *RedisClusterClientWrapper) Do(args ...interface{}) *redis.Cmd {
	res0 := w.obj.Do(args...)
	return res0
}

func (w *RedisClusterClientWrapper) ForEachMaster(ctx context.Context, fn func(client *redis.Client) error) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ForEachMaster")
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
			if err := w.rateLimiterGroup.Wait(ctx, "ClusterClient.ForEachMaster"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("redis.ClusterClient.ForEachMaster", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ForEachMaster", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.ForEachMaster(fn)
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) ForEachNode(ctx context.Context, fn func(client *redis.Client) error) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ForEachNode")
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
			if err := w.rateLimiterGroup.Wait(ctx, "ClusterClient.ForEachNode"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("redis.ClusterClient.ForEachNode", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ForEachNode", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.ForEachNode(fn)
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) ForEachSlave(ctx context.Context, fn func(client *redis.Client) error) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ForEachSlave")
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
			if err := w.rateLimiterGroup.Wait(ctx, "ClusterClient.ForEachSlave"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("redis.ClusterClient.ForEachSlave", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ForEachSlave", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.ForEachSlave(fn)
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) Options() *redis.ClusterOptions {
	res0 := w.obj.Options()
	return res0
}

func (w *RedisClusterClientWrapper) PSubscribe(channels ...string) *redis.PubSub {
	res0 := w.obj.PSubscribe(channels...)
	return res0
}

func (w *RedisClusterClientWrapper) Pipeline() redis.Pipeliner {
	res0 := w.obj.Pipeline()
	return res0
}

func (w *RedisClusterClientWrapper) Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Pipelined")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 []redis.Cmder
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "ClusterClient.Pipelined"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("redis.ClusterClient.Pipelined", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Pipelined", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Pipelined(fn)
		return err
	})
	return res0, err
}

func (w *RedisClusterClientWrapper) PoolStats() *redis.PoolStats {
	res0 := w.obj.PoolStats()
	return res0
}

func (w *RedisClusterClientWrapper) Process(ctx context.Context, cmd redis.Cmder) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Process")
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
			if err := w.rateLimiterGroup.Wait(ctx, "ClusterClient.Process"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("redis.ClusterClient.Process", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Process", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Process(cmd)
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) ReloadState(ctx context.Context) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ReloadState")
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
			if err := w.rateLimiterGroup.Wait(ctx, "ClusterClient.ReloadState"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("redis.ClusterClient.ReloadState", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ReloadState", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.ReloadState()
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) Subscribe(channels ...string) *redis.PubSub {
	res0 := w.obj.Subscribe(channels...)
	return res0
}

func (w *RedisClusterClientWrapper) TxPipeline() redis.Pipeliner {
	res0 := w.obj.TxPipeline()
	return res0
}

func (w *RedisClusterClientWrapper) TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.TxPipelined")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
	var res0 []redis.Cmder
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "ClusterClient.TxPipelined"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("redis.ClusterClient.TxPipelined", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("redis.ClusterClient.TxPipelined", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.TxPipelined(fn)
		return err
	})
	return res0, err
}

func (w *RedisClusterClientWrapper) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	ctxOptions := FromContext(ctx)
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Watch")
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
			if err := w.rateLimiterGroup.Wait(ctx, "ClusterClient.Watch"); err != nil {
				return err
			}
		}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("redis.ClusterClient.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Watch(fn, keys...)
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) WithContext(ctx context.Context) *RedisClusterClientWrapper {
	w.obj = w.obj.WithContext(ctx)
	return w
}

func (w *RedisClusterClientWrapper) WrapProcess(fn func(oldProcess func(redis.Cmder) error) func(redis.Cmder) error) {
	w.obj.WrapProcess(fn)
}

func (w *RedisClusterClientWrapper) WrapProcessPipeline(fn func(oldProcess func([]redis.Cmder) error) func([]redis.Cmder) error) {
	w.obj.WrapProcessPipeline(fn)
}
