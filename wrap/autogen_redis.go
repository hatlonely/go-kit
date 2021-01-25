// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type RedisClientWrapper struct {
	obj            *redis.Client
	retry          *Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	totalMetric    *prometheus.CounterVec
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

func (w *RedisClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "redis_Client_durationMs",
		Help:        "redis Client response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "redis_Client_total",
		Help:        "redis Client request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode"})
}

type RedisClusterClientWrapper struct {
	obj            *redis.ClusterClient
	retry          *Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	totalMetric    *prometheus.CounterVec
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

func (w *RedisClusterClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "redis_ClusterClient_durationMs",
		Help:        "redis ClusterClient response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "redis_ClusterClient_total",
		Help:        "redis ClusterClient request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode"})
}

func (w *RedisClientWrapper) Context(ctx context.Context) context.Context {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.Client.Context")
		defer span.Finish()
	}

	res0 := w.obj.Context()
	return res0
}

func (w *RedisClientWrapper) Options(ctx context.Context) *redis.Options {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.Client.Options")
		defer span.Finish()
	}

	res0 := w.obj.Options()
	return res0
}

func (w *RedisClientWrapper) PSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.Client.PSubscribe")
		defer span.Finish()
	}

	res0 := w.obj.PSubscribe(channels...)
	return res0
}

func (w *RedisClientWrapper) Pipeline(ctx context.Context) redis.Pipeliner {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.Client.Pipeline")
		defer span.Finish()
	}

	res0 := w.obj.Pipeline()
	return res0
}

func (w *RedisClientWrapper) Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.Client.Pipelined")
		defer span.Finish()
	}

	var res0 []redis.Cmder
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.Pipelined(fn)
		return err
	})
	return res0, err
}

func (w *RedisClientWrapper) PoolStats(ctx context.Context) *redis.PoolStats {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.Client.PoolStats")
		defer span.Finish()
	}

	res0 := w.obj.PoolStats()
	return res0
}

func (w *RedisClientWrapper) SetLimiter(ctx context.Context, l redis.Limiter) *RedisClientWrapper {
	w.obj = w.obj.SetLimiter(l)
	return w
}

func (w *RedisClientWrapper) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.Client.Subscribe")
		defer span.Finish()
	}

	res0 := w.obj.Subscribe(channels...)
	return res0
}

func (w *RedisClientWrapper) TxPipeline(ctx context.Context) redis.Pipeliner {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.Client.TxPipeline")
		defer span.Finish()
	}

	res0 := w.obj.TxPipeline()
	return res0
}

func (w *RedisClientWrapper) TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.Client.TxPipelined")
		defer span.Finish()
	}

	var res0 []redis.Cmder
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.TxPipelined(fn)
		return err
	})
	return res0, err
}

func (w *RedisClientWrapper) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.Client.Watch")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
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
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Close")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.Close()
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) Context(ctx context.Context) context.Context {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Context")
		defer span.Finish()
	}

	res0 := w.obj.Context()
	return res0
}

func (w *RedisClusterClientWrapper) DBSize(ctx context.Context) *redis.IntCmd {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.DBSize")
		defer span.Finish()
	}

	res0 := w.obj.DBSize()
	return res0
}

func (w *RedisClusterClientWrapper) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Do")
		defer span.Finish()
	}

	res0 := w.obj.Do(args...)
	return res0
}

func (w *RedisClusterClientWrapper) ForEachMaster(ctx context.Context, fn func(client *redis.Client) error) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ForEachMaster")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.ForEachMaster(fn)
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) ForEachNode(ctx context.Context, fn func(client *redis.Client) error) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ForEachNode")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.ForEachNode(fn)
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) ForEachSlave(ctx context.Context, fn func(client *redis.Client) error) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ForEachSlave")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.ForEachSlave(fn)
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) Options(ctx context.Context) *redis.ClusterOptions {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Options")
		defer span.Finish()
	}

	res0 := w.obj.Options()
	return res0
}

func (w *RedisClusterClientWrapper) PSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.PSubscribe")
		defer span.Finish()
	}

	res0 := w.obj.PSubscribe(channels...)
	return res0
}

func (w *RedisClusterClientWrapper) Pipeline(ctx context.Context) redis.Pipeliner {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Pipeline")
		defer span.Finish()
	}

	res0 := w.obj.Pipeline()
	return res0
}

func (w *RedisClusterClientWrapper) Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Pipelined")
		defer span.Finish()
	}

	var res0 []redis.Cmder
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.Pipelined(fn)
		return err
	})
	return res0, err
}

func (w *RedisClusterClientWrapper) PoolStats(ctx context.Context) *redis.PoolStats {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.PoolStats")
		defer span.Finish()
	}

	res0 := w.obj.PoolStats()
	return res0
}

func (w *RedisClusterClientWrapper) Process(ctx context.Context, cmd redis.Cmder) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Process")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.Process(cmd)
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) ReloadState(ctx context.Context) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ReloadState")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.ReloadState()
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Subscribe")
		defer span.Finish()
	}

	res0 := w.obj.Subscribe(channels...)
	return res0
}

func (w *RedisClusterClientWrapper) TxPipeline(ctx context.Context) redis.Pipeliner {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.TxPipeline")
		defer span.Finish()
	}

	res0 := w.obj.TxPipeline()
	return res0
}

func (w *RedisClusterClientWrapper) TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.TxPipelined")
		defer span.Finish()
	}

	var res0 []redis.Cmder
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.TxPipelined(fn)
		return err
	})
	return res0, err
}

func (w *RedisClusterClientWrapper) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Watch")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.Watch(fn, keys...)
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) WithContext(ctx context.Context) *RedisClusterClientWrapper {
	w.obj = w.obj.WithContext(ctx)
	return w
}

func (w *RedisClusterClientWrapper) WrapProcess(ctx context.Context, fn func(oldProcess func(redis.Cmder) error) func(redis.Cmder) error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.WrapProcess")
		defer span.Finish()
	}

	w.obj.WrapProcess(fn)
}

func (w *RedisClusterClientWrapper) WrapProcessPipeline(ctx context.Context, fn func(oldProcess func([]redis.Cmder) error) func([]redis.Cmder) error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.WrapProcessPipeline")
		defer span.Finish()
	}

	w.obj.WrapProcessPipeline(fn)
}
