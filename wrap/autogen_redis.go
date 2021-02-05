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
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

type RedisClientWrapper struct {
	obj            *redis.Client
	retry          *micro.Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	inflightMetric *prometheus.GaugeVec
	rateLimiter    micro.RateLimiter
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

func (w *RedisClientWrapper) OnRateLimiterChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *RedisClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "redis_Client_durationMs",
		Help:        "redis Client response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "redis_Client_inflight",
		Help:        "redis Client inflight",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "custom"})
}

type RedisClusterClientWrapper struct {
	obj            *redis.ClusterClient
	retry          *micro.Retry
	options        *WrapperOptions
	durationMetric *prometheus.HistogramVec
	inflightMetric *prometheus.GaugeVec
	rateLimiter    micro.RateLimiter
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

func (w *RedisClusterClientWrapper) OnRateLimiterChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *RedisClusterClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "redis_ClusterClient_durationMs",
		Help:        "redis ClusterClient response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "redis_ClusterClient_inflight",
		Help:        "redis ClusterClient inflight",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "custom"})
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
	var res0 []redis.Cmder
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Pipelined"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Pipelined")
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
			w.inflightMetric.WithLabelValues("redis.Client.Pipelined", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Pipelined", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Pipelined", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Pipelined(fn)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
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
	var res0 []redis.Cmder
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.TxPipelined"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.TxPipelined")
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
			w.inflightMetric.WithLabelValues("redis.Client.TxPipelined", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.TxPipelined", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.TxPipelined", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.TxPipelined(fn)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *RedisClientWrapper) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Watch"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Watch")
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
			w.inflightMetric.WithLabelValues("redis.Client.Watch", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Watch", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Watch(fn, keys...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *RedisClientWrapper) WithContext(ctx context.Context) *RedisClientWrapper {
	w.obj = w.obj.WithContext(ctx)
	return w
}

func (w *RedisClientWrapper) Close(ctx context.Context) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Close"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Close")
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
			w.inflightMetric.WithLabelValues("redis.Client.Close", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Close", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Close", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Close()
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *RedisClientWrapper) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.Cmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Do"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Do")
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
			w.inflightMetric.WithLabelValues("redis.Client.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Do", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Do(args...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Process(ctx context.Context, cmd redis.Cmder) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Process"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Process")
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
			w.inflightMetric.WithLabelValues("redis.Client.Process", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Process", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Process", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Process(cmd)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *RedisClientWrapper) String() string {
	res0 := w.obj.String()
	return res0
}

func (w *RedisClientWrapper) WrapProcess(fn func(oldProcess func(cmd redis.Cmder) error) func(cmd redis.Cmder) error) {
	w.obj.WrapProcess(fn)
}

func (w *RedisClientWrapper) WrapProcessPipeline(fn func(oldProcess func([]redis.Cmder) error) func([]redis.Cmder) error) {
	w.obj.WrapProcessPipeline(fn)
}

func (w *RedisClientWrapper) Append(ctx context.Context, key string, value string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Append"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Append")
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
			w.inflightMetric.WithLabelValues("redis.Client.Append", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Append", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Append", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Append(key, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.BLPop"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.BLPop")
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
			w.inflightMetric.WithLabelValues("redis.Client.BLPop", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.BLPop", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.BLPop", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BLPop(timeout, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.BRPop"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.BRPop")
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
			w.inflightMetric.WithLabelValues("redis.Client.BRPop", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.BRPop", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.BRPop", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BRPop(timeout, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) BRPopLPush(ctx context.Context, source string, destination string, timeout time.Duration) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.BRPopLPush"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.BRPopLPush")
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
			w.inflightMetric.WithLabelValues("redis.Client.BRPopLPush", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.BRPopLPush", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.BRPopLPush", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BRPopLPush(source, destination, timeout)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZWithKeyCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.BZPopMax"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.BZPopMax")
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
			w.inflightMetric.WithLabelValues("redis.Client.BZPopMax", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.BZPopMax", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.BZPopMax", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BZPopMax(timeout, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZWithKeyCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.BZPopMin"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.BZPopMin")
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
			w.inflightMetric.WithLabelValues("redis.Client.BZPopMin", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.BZPopMin", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.BZPopMin", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BZPopMin(timeout, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) BgRewriteAOF(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.BgRewriteAOF"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.BgRewriteAOF")
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
			w.inflightMetric.WithLabelValues("redis.Client.BgRewriteAOF", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.BgRewriteAOF", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.BgRewriteAOF", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BgRewriteAOF()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) BgSave(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.BgSave"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.BgSave")
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
			w.inflightMetric.WithLabelValues("redis.Client.BgSave", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.BgSave", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.BgSave", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BgSave()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.BitCount"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.BitCount")
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
			w.inflightMetric.WithLabelValues("redis.Client.BitCount", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.BitCount", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.BitCount", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BitCount(key, bitCount)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) BitOpAnd(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.BitOpAnd"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.BitOpAnd")
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
			w.inflightMetric.WithLabelValues("redis.Client.BitOpAnd", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.BitOpAnd", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.BitOpAnd", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BitOpAnd(destKey, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) BitOpNot(ctx context.Context, destKey string, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.BitOpNot"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.BitOpNot")
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
			w.inflightMetric.WithLabelValues("redis.Client.BitOpNot", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.BitOpNot", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.BitOpNot", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BitOpNot(destKey, key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) BitOpOr(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.BitOpOr"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.BitOpOr")
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
			w.inflightMetric.WithLabelValues("redis.Client.BitOpOr", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.BitOpOr", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.BitOpOr", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BitOpOr(destKey, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) BitOpXor(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.BitOpXor"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.BitOpXor")
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
			w.inflightMetric.WithLabelValues("redis.Client.BitOpXor", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.BitOpXor", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.BitOpXor", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BitOpXor(destKey, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.BitPos"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.BitPos")
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
			w.inflightMetric.WithLabelValues("redis.Client.BitPos", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.BitPos", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.BitPos", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BitPos(key, bit, pos...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClientGetName(ctx context.Context) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClientGetName"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClientGetName")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClientGetName", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClientGetName", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClientGetName", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientGetName()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClientID(ctx context.Context) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClientID"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClientID")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClientID", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClientID", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClientID", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientID()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClientKill(ctx context.Context, ipPort string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClientKill"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClientKill")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClientKill", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClientKill", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClientKill", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientKill(ipPort)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClientKillByFilter(ctx context.Context, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClientKillByFilter"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClientKillByFilter")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClientKillByFilter", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClientKillByFilter", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClientKillByFilter", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientKillByFilter(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClientList(ctx context.Context) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClientList"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClientList")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClientList", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClientList", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClientList", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientList()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClientPause(ctx context.Context, dur time.Duration) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClientPause"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClientPause")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClientPause", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClientPause", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClientPause", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientPause(dur)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClientUnblock(ctx context.Context, id int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClientUnblock"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClientUnblock")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClientUnblock", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClientUnblock", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClientUnblock", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientUnblock(id)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClientUnblockWithError(ctx context.Context, id int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClientUnblockWithError"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClientUnblockWithError")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClientUnblockWithError", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClientUnblockWithError", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClientUnblockWithError", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientUnblockWithError(id)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterAddSlots(ctx context.Context, slots ...int) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterAddSlots"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterAddSlots")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterAddSlots", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterAddSlots", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterAddSlots", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterAddSlots(slots...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterAddSlotsRange(ctx context.Context, min int, max int) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterAddSlotsRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterAddSlotsRange")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterAddSlotsRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterAddSlotsRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterAddSlotsRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterAddSlotsRange(min, max)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterCountFailureReports(ctx context.Context, nodeID string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterCountFailureReports"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterCountFailureReports")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterCountFailureReports", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterCountFailureReports", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterCountFailureReports", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterCountFailureReports(nodeID)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterCountKeysInSlot(ctx context.Context, slot int) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterCountKeysInSlot"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterCountKeysInSlot")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterCountKeysInSlot", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterCountKeysInSlot", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterCountKeysInSlot", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterCountKeysInSlot(slot)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterDelSlots(ctx context.Context, slots ...int) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterDelSlots"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterDelSlots")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterDelSlots", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterDelSlots", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterDelSlots", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterDelSlots(slots...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterDelSlotsRange(ctx context.Context, min int, max int) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterDelSlotsRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterDelSlotsRange")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterDelSlotsRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterDelSlotsRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterDelSlotsRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterDelSlotsRange(min, max)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterFailover(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterFailover"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterFailover")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterFailover", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterFailover", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterFailover", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterFailover()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterForget(ctx context.Context, nodeID string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterForget"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterForget")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterForget", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterForget", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterForget", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterForget(nodeID)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterGetKeysInSlot"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterGetKeysInSlot")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterGetKeysInSlot", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterGetKeysInSlot", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterGetKeysInSlot", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterGetKeysInSlot(slot, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterInfo(ctx context.Context) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterInfo"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterInfo")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterInfo", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterInfo", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterInfo", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterInfo()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterKeySlot(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterKeySlot"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterKeySlot")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterKeySlot", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterKeySlot", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterKeySlot", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterKeySlot(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterMeet(ctx context.Context, host string, port string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterMeet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterMeet")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterMeet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterMeet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterMeet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterMeet(host, port)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterNodes(ctx context.Context) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterNodes"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterNodes")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterNodes", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterNodes", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterNodes", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterNodes()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterReplicate(ctx context.Context, nodeID string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterReplicate"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterReplicate")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterReplicate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterReplicate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterReplicate", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterReplicate(nodeID)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterResetHard(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterResetHard"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterResetHard")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterResetHard", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterResetHard", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterResetHard", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterResetHard()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterResetSoft(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterResetSoft"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterResetSoft")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterResetSoft", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterResetSoft", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterResetSoft", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterResetSoft()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterSaveConfig(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterSaveConfig"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterSaveConfig")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterSaveConfig", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterSaveConfig", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterSaveConfig", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterSaveConfig()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterSlaves(ctx context.Context, nodeID string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterSlaves"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterSlaves")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterSlaves", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterSlaves", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterSlaves", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterSlaves(nodeID)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ClusterSlots(ctx context.Context) *redis.ClusterSlotsCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ClusterSlotsCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ClusterSlots"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ClusterSlots")
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
			w.inflightMetric.WithLabelValues("redis.Client.ClusterSlots", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ClusterSlots", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ClusterSlots", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterSlots()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Command(ctx context.Context) *redis.CommandsInfoCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.CommandsInfoCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Command"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Command")
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
			w.inflightMetric.WithLabelValues("redis.Client.Command", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Command", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Command", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Command()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ConfigGet(ctx context.Context, parameter string) *redis.SliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.SliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ConfigGet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ConfigGet")
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
			w.inflightMetric.WithLabelValues("redis.Client.ConfigGet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ConfigGet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ConfigGet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ConfigGet(parameter)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ConfigResetStat(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ConfigResetStat"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ConfigResetStat")
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
			w.inflightMetric.WithLabelValues("redis.Client.ConfigResetStat", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ConfigResetStat", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ConfigResetStat", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ConfigResetStat()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ConfigRewrite(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ConfigRewrite"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ConfigRewrite")
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
			w.inflightMetric.WithLabelValues("redis.Client.ConfigRewrite", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ConfigRewrite", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ConfigRewrite", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ConfigRewrite()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ConfigSet(ctx context.Context, parameter string, value string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ConfigSet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ConfigSet")
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
			w.inflightMetric.WithLabelValues("redis.Client.ConfigSet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ConfigSet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ConfigSet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ConfigSet(parameter, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) DBSize(ctx context.Context) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.DBSize"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.DBSize")
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
			w.inflightMetric.WithLabelValues("redis.Client.DBSize", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.DBSize", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.DBSize", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.DBSize()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) DbSize(ctx context.Context) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.DbSize"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.DbSize")
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
			w.inflightMetric.WithLabelValues("redis.Client.DbSize", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.DbSize", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.DbSize", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.DbSize()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) DebugObject(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.DebugObject"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.DebugObject")
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
			w.inflightMetric.WithLabelValues("redis.Client.DebugObject", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.DebugObject", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.DebugObject", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.DebugObject(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Decr(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Decr"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Decr")
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
			w.inflightMetric.WithLabelValues("redis.Client.Decr", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Decr", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Decr", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Decr(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.DecrBy"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.DecrBy")
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
			w.inflightMetric.WithLabelValues("redis.Client.DecrBy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.DecrBy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.DecrBy", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.DecrBy(key, decrement)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Del"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Del")
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
			w.inflightMetric.WithLabelValues("redis.Client.Del", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Del", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Del", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Del(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Dump(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Dump"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Dump")
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
			w.inflightMetric.WithLabelValues("redis.Client.Dump", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Dump", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Dump", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Dump(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Echo(ctx context.Context, message interface{}) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Echo"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Echo")
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
			w.inflightMetric.WithLabelValues("redis.Client.Echo", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Echo", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Echo", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Echo(message)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.Cmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Eval"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Eval")
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
			w.inflightMetric.WithLabelValues("redis.Client.Eval", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Eval", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Eval", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Eval(script, keys, args...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.Cmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.EvalSha"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.EvalSha")
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
			w.inflightMetric.WithLabelValues("redis.Client.EvalSha", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.EvalSha", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.EvalSha", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.EvalSha(sha1, keys, args...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Exists"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Exists")
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
			w.inflightMetric.WithLabelValues("redis.Client.Exists", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Exists", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Exists", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Exists(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Expire"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Expire")
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
			w.inflightMetric.WithLabelValues("redis.Client.Expire", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Expire", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Expire", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Expire(key, expiration)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ExpireAt"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ExpireAt")
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
			w.inflightMetric.WithLabelValues("redis.Client.ExpireAt", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ExpireAt", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ExpireAt", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ExpireAt(key, tm)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) FlushAll(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.FlushAll"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.FlushAll")
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
			w.inflightMetric.WithLabelValues("redis.Client.FlushAll", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.FlushAll", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.FlushAll", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.FlushAll()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) FlushAllAsync(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.FlushAllAsync"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.FlushAllAsync")
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
			w.inflightMetric.WithLabelValues("redis.Client.FlushAllAsync", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.FlushAllAsync", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.FlushAllAsync", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.FlushAllAsync()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) FlushDB(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.FlushDB"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.FlushDB")
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
			w.inflightMetric.WithLabelValues("redis.Client.FlushDB", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.FlushDB", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.FlushDB", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.FlushDB()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) FlushDBAsync(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.FlushDBAsync"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.FlushDBAsync")
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
			w.inflightMetric.WithLabelValues("redis.Client.FlushDBAsync", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.FlushDBAsync", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.FlushDBAsync", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.FlushDBAsync()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) FlushDb(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.FlushDb"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.FlushDb")
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
			w.inflightMetric.WithLabelValues("redis.Client.FlushDb", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.FlushDb", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.FlushDb", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.FlushDb()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.GeoAdd"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.GeoAdd")
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
			w.inflightMetric.WithLabelValues("redis.Client.GeoAdd", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.GeoAdd", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.GeoAdd", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoAdd(key, geoLocation...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) GeoDist(ctx context.Context, key string, member1 string, member2 string, unit string) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.GeoDist"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.GeoDist")
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
			w.inflightMetric.WithLabelValues("redis.Client.GeoDist", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.GeoDist", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.GeoDist", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoDist(key, member1, member2, unit)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) GeoHash(ctx context.Context, key string, members ...string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.GeoHash"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.GeoHash")
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
			w.inflightMetric.WithLabelValues("redis.Client.GeoHash", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.GeoHash", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.GeoHash", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoHash(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) GeoPos(ctx context.Context, key string, members ...string) *redis.GeoPosCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.GeoPosCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.GeoPos"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.GeoPos")
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
			w.inflightMetric.WithLabelValues("redis.Client.GeoPos", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.GeoPos", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.GeoPos", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoPos(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) GeoRadius(ctx context.Context, key string, longitude float64, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.GeoLocationCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.GeoRadius"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.GeoRadius")
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
			w.inflightMetric.WithLabelValues("redis.Client.GeoRadius", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.GeoRadius", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.GeoRadius", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoRadius(key, longitude, latitude, query)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) GeoRadiusByMember(ctx context.Context, key string, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.GeoLocationCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.GeoRadiusByMember"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.GeoRadiusByMember")
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
			w.inflightMetric.WithLabelValues("redis.Client.GeoRadiusByMember", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.GeoRadiusByMember", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.GeoRadiusByMember", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoRadiusByMember(key, member, query)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) GeoRadiusByMemberRO(ctx context.Context, key string, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.GeoLocationCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.GeoRadiusByMemberRO"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.GeoRadiusByMemberRO")
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
			w.inflightMetric.WithLabelValues("redis.Client.GeoRadiusByMemberRO", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.GeoRadiusByMemberRO", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.GeoRadiusByMemberRO", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoRadiusByMemberRO(key, member, query)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) GeoRadiusRO(ctx context.Context, key string, longitude float64, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.GeoLocationCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.GeoRadiusRO"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.GeoRadiusRO")
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
			w.inflightMetric.WithLabelValues("redis.Client.GeoRadiusRO", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.GeoRadiusRO", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.GeoRadiusRO", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoRadiusRO(key, longitude, latitude, query)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Get(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Get"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Get")
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
			w.inflightMetric.WithLabelValues("redis.Client.Get", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Get", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Get", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Get(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.GetBit"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.GetBit")
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
			w.inflightMetric.WithLabelValues("redis.Client.GetBit", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.GetBit", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.GetBit", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GetBit(key, offset)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) GetRange(ctx context.Context, key string, start int64, end int64) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.GetRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.GetRange")
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
			w.inflightMetric.WithLabelValues("redis.Client.GetRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.GetRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.GetRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GetRange(key, start, end)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.GetSet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.GetSet")
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
			w.inflightMetric.WithLabelValues("redis.Client.GetSet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.GetSet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.GetSet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GetSet(key, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.HDel"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.HDel")
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
			w.inflightMetric.WithLabelValues("redis.Client.HDel", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.HDel", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.HDel", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HDel(key, fields...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) HExists(ctx context.Context, key string, field string) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.HExists"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.HExists")
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
			w.inflightMetric.WithLabelValues("redis.Client.HExists", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.HExists", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.HExists", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HExists(key, field)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) HGet(ctx context.Context, key string, field string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.HGet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.HGet")
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
			w.inflightMetric.WithLabelValues("redis.Client.HGet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.HGet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.HGet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HGet(key, field)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringStringMapCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.HGetAll"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.HGetAll")
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
			w.inflightMetric.WithLabelValues("redis.Client.HGetAll", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.HGetAll", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.HGetAll", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HGetAll(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) HIncrBy(ctx context.Context, key string, field string, incr int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.HIncrBy"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.HIncrBy")
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
			w.inflightMetric.WithLabelValues("redis.Client.HIncrBy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.HIncrBy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.HIncrBy", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HIncrBy(key, field, incr)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) HIncrByFloat(ctx context.Context, key string, field string, incr float64) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.HIncrByFloat"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.HIncrByFloat")
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
			w.inflightMetric.WithLabelValues("redis.Client.HIncrByFloat", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.HIncrByFloat", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.HIncrByFloat", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HIncrByFloat(key, field, incr)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.HKeys"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.HKeys")
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
			w.inflightMetric.WithLabelValues("redis.Client.HKeys", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.HKeys", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.HKeys", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HKeys(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) HLen(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.HLen"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.HLen")
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
			w.inflightMetric.WithLabelValues("redis.Client.HLen", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.HLen", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.HLen", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HLen(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.SliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.HMGet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.HMGet")
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
			w.inflightMetric.WithLabelValues("redis.Client.HMGet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.HMGet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.HMGet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HMGet(key, fields...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) HMSet(ctx context.Context, key string, fields map[string]interface{}) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.HMSet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.HMSet")
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
			w.inflightMetric.WithLabelValues("redis.Client.HMSet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.HMSet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.HMSet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HMSet(key, fields)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ScanCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.HScan"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.HScan")
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
			w.inflightMetric.WithLabelValues("redis.Client.HScan", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.HScan", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.HScan", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HScan(key, cursor, match, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) HSet(ctx context.Context, key string, field string, value interface{}) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.HSet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.HSet")
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
			w.inflightMetric.WithLabelValues("redis.Client.HSet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.HSet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.HSet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HSet(key, field, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) HSetNX(ctx context.Context, key string, field string, value interface{}) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.HSetNX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.HSetNX")
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
			w.inflightMetric.WithLabelValues("redis.Client.HSetNX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.HSetNX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.HSetNX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HSetNX(key, field, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.HVals"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.HVals")
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
			w.inflightMetric.WithLabelValues("redis.Client.HVals", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.HVals", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.HVals", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HVals(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Incr(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Incr"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Incr")
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
			w.inflightMetric.WithLabelValues("redis.Client.Incr", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Incr", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Incr", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Incr(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.IncrBy"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.IncrBy")
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
			w.inflightMetric.WithLabelValues("redis.Client.IncrBy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.IncrBy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.IncrBy", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.IncrBy(key, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.IncrByFloat"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.IncrByFloat")
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
			w.inflightMetric.WithLabelValues("redis.Client.IncrByFloat", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.IncrByFloat", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.IncrByFloat", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.IncrByFloat(key, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Info(ctx context.Context, section ...string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Info"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Info")
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
			w.inflightMetric.WithLabelValues("redis.Client.Info", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Info", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Info", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Info(section...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Keys"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Keys")
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
			w.inflightMetric.WithLabelValues("redis.Client.Keys", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Keys", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Keys", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Keys(pattern)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.LIndex"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.LIndex")
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
			w.inflightMetric.WithLabelValues("redis.Client.LIndex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.LIndex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.LIndex", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LIndex(key, index)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) LInsert(ctx context.Context, key string, op string, pivot interface{}, value interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.LInsert"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.LInsert")
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
			w.inflightMetric.WithLabelValues("redis.Client.LInsert", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.LInsert", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.LInsert", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LInsert(key, op, pivot, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) LInsertAfter(ctx context.Context, key string, pivot interface{}, value interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.LInsertAfter"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.LInsertAfter")
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
			w.inflightMetric.WithLabelValues("redis.Client.LInsertAfter", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.LInsertAfter", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.LInsertAfter", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LInsertAfter(key, pivot, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) LInsertBefore(ctx context.Context, key string, pivot interface{}, value interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.LInsertBefore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.LInsertBefore")
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
			w.inflightMetric.WithLabelValues("redis.Client.LInsertBefore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.LInsertBefore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.LInsertBefore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LInsertBefore(key, pivot, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) LLen(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.LLen"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.LLen")
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
			w.inflightMetric.WithLabelValues("redis.Client.LLen", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.LLen", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.LLen", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LLen(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) LPop(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.LPop"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.LPop")
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
			w.inflightMetric.WithLabelValues("redis.Client.LPop", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.LPop", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.LPop", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LPop(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.LPush"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.LPush")
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
			w.inflightMetric.WithLabelValues("redis.Client.LPush", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.LPush", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.LPush", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LPush(key, values...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) LPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.LPushX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.LPushX")
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
			w.inflightMetric.WithLabelValues("redis.Client.LPushX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.LPushX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.LPushX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LPushX(key, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) LRange(ctx context.Context, key string, start int64, stop int64) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.LRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.LRange")
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
			w.inflightMetric.WithLabelValues("redis.Client.LRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.LRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.LRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LRange(key, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.LRem"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.LRem")
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
			w.inflightMetric.WithLabelValues("redis.Client.LRem", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.LRem", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.LRem", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LRem(key, count, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.LSet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.LSet")
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
			w.inflightMetric.WithLabelValues("redis.Client.LSet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.LSet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.LSet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LSet(key, index, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) LTrim(ctx context.Context, key string, start int64, stop int64) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.LTrim"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.LTrim")
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
			w.inflightMetric.WithLabelValues("redis.Client.LTrim", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.LTrim", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.LTrim", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LTrim(key, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) LastSave(ctx context.Context) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.LastSave"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.LastSave")
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
			w.inflightMetric.WithLabelValues("redis.Client.LastSave", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.LastSave", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.LastSave", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LastSave()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.SliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.MGet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.MGet")
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
			w.inflightMetric.WithLabelValues("redis.Client.MGet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.MGet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.MGet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.MGet(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) MSet(ctx context.Context, pairs ...interface{}) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.MSet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.MSet")
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
			w.inflightMetric.WithLabelValues("redis.Client.MSet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.MSet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.MSet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.MSet(pairs...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) MSetNX(ctx context.Context, pairs ...interface{}) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.MSetNX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.MSetNX")
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
			w.inflightMetric.WithLabelValues("redis.Client.MSetNX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.MSetNX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.MSetNX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.MSetNX(pairs...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) MemoryUsage(ctx context.Context, key string, samples ...int) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.MemoryUsage"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.MemoryUsage")
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
			w.inflightMetric.WithLabelValues("redis.Client.MemoryUsage", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.MemoryUsage", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.MemoryUsage", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.MemoryUsage(key, samples...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Migrate(ctx context.Context, host string, port string, key string, db int64, timeout time.Duration) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Migrate"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Migrate")
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
			w.inflightMetric.WithLabelValues("redis.Client.Migrate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Migrate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Migrate", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Migrate(host, port, key, db, timeout)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Move(ctx context.Context, key string, db int64) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Move"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Move")
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
			w.inflightMetric.WithLabelValues("redis.Client.Move", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Move", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Move", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Move(key, db)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ObjectEncoding(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ObjectEncoding"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ObjectEncoding")
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
			w.inflightMetric.WithLabelValues("redis.Client.ObjectEncoding", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ObjectEncoding", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ObjectEncoding", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ObjectEncoding(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ObjectIdleTime(ctx context.Context, key string) *redis.DurationCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.DurationCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ObjectIdleTime"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ObjectIdleTime")
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
			w.inflightMetric.WithLabelValues("redis.Client.ObjectIdleTime", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ObjectIdleTime", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ObjectIdleTime", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ObjectIdleTime(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ObjectRefCount(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ObjectRefCount"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ObjectRefCount")
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
			w.inflightMetric.WithLabelValues("redis.Client.ObjectRefCount", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ObjectRefCount", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ObjectRefCount", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ObjectRefCount(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) PExpire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.PExpire"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.PExpire")
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
			w.inflightMetric.WithLabelValues("redis.Client.PExpire", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.PExpire", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.PExpire", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PExpire(key, expiration)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) PExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.PExpireAt"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.PExpireAt")
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
			w.inflightMetric.WithLabelValues("redis.Client.PExpireAt", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.PExpireAt", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.PExpireAt", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PExpireAt(key, tm)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) PFAdd(ctx context.Context, key string, els ...interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.PFAdd"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.PFAdd")
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
			w.inflightMetric.WithLabelValues("redis.Client.PFAdd", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.PFAdd", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.PFAdd", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PFAdd(key, els...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) PFCount(ctx context.Context, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.PFCount"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.PFCount")
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
			w.inflightMetric.WithLabelValues("redis.Client.PFCount", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.PFCount", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.PFCount", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PFCount(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) PFMerge(ctx context.Context, dest string, keys ...string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.PFMerge"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.PFMerge")
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
			w.inflightMetric.WithLabelValues("redis.Client.PFMerge", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.PFMerge", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.PFMerge", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PFMerge(dest, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) PTTL(ctx context.Context, key string) *redis.DurationCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.DurationCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.PTTL"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.PTTL")
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
			w.inflightMetric.WithLabelValues("redis.Client.PTTL", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.PTTL", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.PTTL", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PTTL(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Persist(ctx context.Context, key string) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Persist"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Persist")
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
			w.inflightMetric.WithLabelValues("redis.Client.Persist", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Persist", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Persist", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Persist(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Ping(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Ping"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Ping")
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
			w.inflightMetric.WithLabelValues("redis.Client.Ping", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Ping", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Ping", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Ping()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) PubSubChannels(ctx context.Context, pattern string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.PubSubChannels"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.PubSubChannels")
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
			w.inflightMetric.WithLabelValues("redis.Client.PubSubChannels", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.PubSubChannels", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.PubSubChannels", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PubSubChannels(pattern)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) PubSubNumPat(ctx context.Context) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.PubSubNumPat"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.PubSubNumPat")
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
			w.inflightMetric.WithLabelValues("redis.Client.PubSubNumPat", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.PubSubNumPat", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.PubSubNumPat", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PubSubNumPat()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) PubSubNumSub(ctx context.Context, channels ...string) *redis.StringIntMapCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringIntMapCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.PubSubNumSub"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.PubSubNumSub")
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
			w.inflightMetric.WithLabelValues("redis.Client.PubSubNumSub", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.PubSubNumSub", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.PubSubNumSub", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PubSubNumSub(channels...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Publish"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Publish")
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
			w.inflightMetric.WithLabelValues("redis.Client.Publish", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Publish", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Publish", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Publish(channel, message)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Quit(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Quit"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Quit")
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
			w.inflightMetric.WithLabelValues("redis.Client.Quit", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Quit", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Quit", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Quit()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) RPop(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.RPop"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.RPop")
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
			w.inflightMetric.WithLabelValues("redis.Client.RPop", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.RPop", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.RPop", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RPop(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) RPopLPush(ctx context.Context, source string, destination string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.RPopLPush"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.RPopLPush")
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
			w.inflightMetric.WithLabelValues("redis.Client.RPopLPush", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.RPopLPush", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.RPopLPush", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RPopLPush(source, destination)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.RPush"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.RPush")
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
			w.inflightMetric.WithLabelValues("redis.Client.RPush", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.RPush", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.RPush", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RPush(key, values...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) RPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.RPushX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.RPushX")
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
			w.inflightMetric.WithLabelValues("redis.Client.RPushX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.RPushX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.RPushX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RPushX(key, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) RandomKey(ctx context.Context) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.RandomKey"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.RandomKey")
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
			w.inflightMetric.WithLabelValues("redis.Client.RandomKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.RandomKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.RandomKey", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RandomKey()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ReadOnly(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ReadOnly"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ReadOnly")
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
			w.inflightMetric.WithLabelValues("redis.Client.ReadOnly", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ReadOnly", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ReadOnly", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ReadOnly()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ReadWrite(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ReadWrite"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ReadWrite")
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
			w.inflightMetric.WithLabelValues("redis.Client.ReadWrite", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ReadWrite", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ReadWrite", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ReadWrite()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Rename(ctx context.Context, key string, newkey string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Rename"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Rename")
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
			w.inflightMetric.WithLabelValues("redis.Client.Rename", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Rename", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Rename", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Rename(key, newkey)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) RenameNX(ctx context.Context, key string, newkey string) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.RenameNX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.RenameNX")
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
			w.inflightMetric.WithLabelValues("redis.Client.RenameNX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.RenameNX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.RenameNX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RenameNX(key, newkey)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Restore(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Restore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Restore")
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
			w.inflightMetric.WithLabelValues("redis.Client.Restore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Restore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Restore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Restore(key, ttl, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.RestoreReplace"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.RestoreReplace")
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
			w.inflightMetric.WithLabelValues("redis.Client.RestoreReplace", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.RestoreReplace", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.RestoreReplace", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RestoreReplace(key, ttl, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SAdd"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SAdd")
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
			w.inflightMetric.WithLabelValues("redis.Client.SAdd", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SAdd", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SAdd", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SAdd(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SCard(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SCard"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SCard")
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
			w.inflightMetric.WithLabelValues("redis.Client.SCard", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SCard", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SCard", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SCard(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SDiff"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SDiff")
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
			w.inflightMetric.WithLabelValues("redis.Client.SDiff", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SDiff", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SDiff", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SDiff(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SDiffStore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SDiffStore")
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
			w.inflightMetric.WithLabelValues("redis.Client.SDiffStore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SDiffStore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SDiffStore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SDiffStore(destination, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SInter(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SInter"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SInter")
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
			w.inflightMetric.WithLabelValues("redis.Client.SInter", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SInter", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SInter", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SInter(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SInterStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SInterStore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SInterStore")
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
			w.inflightMetric.WithLabelValues("redis.Client.SInterStore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SInterStore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SInterStore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SInterStore(destination, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SIsMember"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SIsMember")
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
			w.inflightMetric.WithLabelValues("redis.Client.SIsMember", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SIsMember", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SIsMember", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SIsMember(key, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SMembers"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SMembers")
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
			w.inflightMetric.WithLabelValues("redis.Client.SMembers", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SMembers", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SMembers", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SMembers(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SMembersMap(ctx context.Context, key string) *redis.StringStructMapCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringStructMapCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SMembersMap"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SMembersMap")
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
			w.inflightMetric.WithLabelValues("redis.Client.SMembersMap", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SMembersMap", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SMembersMap", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SMembersMap(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SMove(ctx context.Context, source string, destination string, member interface{}) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SMove"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SMove")
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
			w.inflightMetric.WithLabelValues("redis.Client.SMove", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SMove", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SMove", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SMove(source, destination, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SPop(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SPop"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SPop")
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
			w.inflightMetric.WithLabelValues("redis.Client.SPop", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SPop", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SPop", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SPop(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SPopN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SPopN"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SPopN")
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
			w.inflightMetric.WithLabelValues("redis.Client.SPopN", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SPopN", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SPopN", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SPopN(key, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SRandMember(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SRandMember"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SRandMember")
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
			w.inflightMetric.WithLabelValues("redis.Client.SRandMember", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SRandMember", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SRandMember", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SRandMember(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SRandMemberN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SRandMemberN"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SRandMemberN")
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
			w.inflightMetric.WithLabelValues("redis.Client.SRandMemberN", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SRandMemberN", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SRandMemberN", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SRandMemberN(key, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SRem"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SRem")
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
			w.inflightMetric.WithLabelValues("redis.Client.SRem", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SRem", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SRem", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SRem(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ScanCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SScan"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SScan")
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
			w.inflightMetric.WithLabelValues("redis.Client.SScan", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SScan", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SScan", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SScan(key, cursor, match, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SUnion(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SUnion"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SUnion")
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
			w.inflightMetric.WithLabelValues("redis.Client.SUnion", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SUnion", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SUnion", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SUnion(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SUnionStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SUnionStore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SUnionStore")
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
			w.inflightMetric.WithLabelValues("redis.Client.SUnionStore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SUnionStore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SUnionStore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SUnionStore(destination, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Save(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Save"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Save")
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
			w.inflightMetric.WithLabelValues("redis.Client.Save", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Save", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Save", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Save()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ScanCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Scan"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Scan")
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
			w.inflightMetric.WithLabelValues("redis.Client.Scan", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Scan", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Scan", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Scan(cursor, match, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ScriptExists"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ScriptExists")
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
			w.inflightMetric.WithLabelValues("redis.Client.ScriptExists", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ScriptExists", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ScriptExists", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ScriptExists(hashes...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ScriptFlush(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ScriptFlush"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ScriptFlush")
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
			w.inflightMetric.WithLabelValues("redis.Client.ScriptFlush", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ScriptFlush", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ScriptFlush", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ScriptFlush()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ScriptKill(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ScriptKill"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ScriptKill")
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
			w.inflightMetric.WithLabelValues("redis.Client.ScriptKill", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ScriptKill", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ScriptKill", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ScriptKill()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ScriptLoad(ctx context.Context, script string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ScriptLoad"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ScriptLoad")
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
			w.inflightMetric.WithLabelValues("redis.Client.ScriptLoad", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ScriptLoad", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ScriptLoad", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ScriptLoad(script)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Set"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Set")
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
			w.inflightMetric.WithLabelValues("redis.Client.Set", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Set", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Set", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Set(key, value, expiration)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SetBit"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SetBit")
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
			w.inflightMetric.WithLabelValues("redis.Client.SetBit", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SetBit", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SetBit", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SetBit(key, offset, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SetNX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SetNX")
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
			w.inflightMetric.WithLabelValues("redis.Client.SetNX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SetNX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SetNX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SetNX(key, value, expiration)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SetRange(ctx context.Context, key string, offset int64, value string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SetRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SetRange")
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
			w.inflightMetric.WithLabelValues("redis.Client.SetRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SetRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SetRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SetRange(key, offset, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SetXX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SetXX")
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
			w.inflightMetric.WithLabelValues("redis.Client.SetXX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SetXX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SetXX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SetXX(key, value, expiration)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Shutdown(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Shutdown"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Shutdown")
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
			w.inflightMetric.WithLabelValues("redis.Client.Shutdown", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Shutdown", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Shutdown", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Shutdown()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ShutdownNoSave(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ShutdownNoSave"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ShutdownNoSave")
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
			w.inflightMetric.WithLabelValues("redis.Client.ShutdownNoSave", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ShutdownNoSave", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ShutdownNoSave", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ShutdownNoSave()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ShutdownSave(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ShutdownSave"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ShutdownSave")
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
			w.inflightMetric.WithLabelValues("redis.Client.ShutdownSave", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ShutdownSave", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ShutdownSave", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ShutdownSave()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SlaveOf(ctx context.Context, host string, port string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SlaveOf"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SlaveOf")
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
			w.inflightMetric.WithLabelValues("redis.Client.SlaveOf", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SlaveOf", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SlaveOf", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SlaveOf(host, port)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SlowLog() {
	w.obj.SlowLog()
}

func (w *RedisClientWrapper) Sort(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Sort"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Sort")
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
			w.inflightMetric.WithLabelValues("redis.Client.Sort", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Sort", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Sort", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Sort(key, sort)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) *redis.SliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.SliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SortInterfaces"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SortInterfaces")
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
			w.inflightMetric.WithLabelValues("redis.Client.SortInterfaces", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SortInterfaces", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SortInterfaces", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SortInterfaces(key, sort)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) SortStore(ctx context.Context, key string, store string, sort *redis.Sort) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.SortStore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.SortStore")
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
			w.inflightMetric.WithLabelValues("redis.Client.SortStore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.SortStore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.SortStore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SortStore(key, store, sort)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) StrLen(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.StrLen"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.StrLen")
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
			w.inflightMetric.WithLabelValues("redis.Client.StrLen", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.StrLen", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.StrLen", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.StrLen(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Sync() {
	w.obj.Sync()
}

func (w *RedisClientWrapper) TTL(ctx context.Context, key string) *redis.DurationCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.DurationCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.TTL"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.TTL")
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
			w.inflightMetric.WithLabelValues("redis.Client.TTL", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.TTL", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.TTL", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.TTL(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Time(ctx context.Context) *redis.TimeCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.TimeCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Time"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Time")
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
			w.inflightMetric.WithLabelValues("redis.Client.Time", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Time", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Time", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Time()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Touch(ctx context.Context, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Touch"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Touch")
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
			w.inflightMetric.WithLabelValues("redis.Client.Touch", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Touch", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Touch", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Touch(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Type(ctx context.Context, key string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Type"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Type")
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
			w.inflightMetric.WithLabelValues("redis.Client.Type", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Type", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Type", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Type(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Unlink(ctx context.Context, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Unlink"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Unlink")
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
			w.inflightMetric.WithLabelValues("redis.Client.Unlink", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Unlink", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Unlink", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Unlink(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) Wait(ctx context.Context, numSlaves int, timeout time.Duration) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.Wait"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.Wait")
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
			w.inflightMetric.WithLabelValues("redis.Client.Wait", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.Wait", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.Wait", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Wait(numSlaves, timeout)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XAck(ctx context.Context, stream string, group string, ids ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XAck"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XAck")
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
			w.inflightMetric.WithLabelValues("redis.Client.XAck", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XAck", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XAck", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XAck(stream, group, ids...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XAdd"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XAdd")
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
			w.inflightMetric.WithLabelValues("redis.Client.XAdd", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XAdd", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XAdd", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XAdd(a)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XClaim(ctx context.Context, a *redis.XClaimArgs) *redis.XMessageSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XMessageSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XClaim"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XClaim")
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
			w.inflightMetric.WithLabelValues("redis.Client.XClaim", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XClaim", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XClaim", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XClaim(a)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XClaimJustID(ctx context.Context, a *redis.XClaimArgs) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XClaimJustID"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XClaimJustID")
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
			w.inflightMetric.WithLabelValues("redis.Client.XClaimJustID", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XClaimJustID", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XClaimJustID", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XClaimJustID(a)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XDel(ctx context.Context, stream string, ids ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XDel"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XDel")
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
			w.inflightMetric.WithLabelValues("redis.Client.XDel", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XDel", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XDel", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XDel(stream, ids...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XGroupCreate(ctx context.Context, stream string, group string, start string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XGroupCreate"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XGroupCreate")
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
			w.inflightMetric.WithLabelValues("redis.Client.XGroupCreate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XGroupCreate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XGroupCreate", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XGroupCreate(stream, group, start)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XGroupCreateMkStream(ctx context.Context, stream string, group string, start string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XGroupCreateMkStream"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XGroupCreateMkStream")
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
			w.inflightMetric.WithLabelValues("redis.Client.XGroupCreateMkStream", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XGroupCreateMkStream", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XGroupCreateMkStream", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XGroupCreateMkStream(stream, group, start)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XGroupDelConsumer(ctx context.Context, stream string, group string, consumer string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XGroupDelConsumer"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XGroupDelConsumer")
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
			w.inflightMetric.WithLabelValues("redis.Client.XGroupDelConsumer", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XGroupDelConsumer", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XGroupDelConsumer", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XGroupDelConsumer(stream, group, consumer)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XGroupDestroy(ctx context.Context, stream string, group string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XGroupDestroy"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XGroupDestroy")
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
			w.inflightMetric.WithLabelValues("redis.Client.XGroupDestroy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XGroupDestroy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XGroupDestroy", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XGroupDestroy(stream, group)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XGroupSetID(ctx context.Context, stream string, group string, start string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XGroupSetID"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XGroupSetID")
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
			w.inflightMetric.WithLabelValues("redis.Client.XGroupSetID", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XGroupSetID", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XGroupSetID", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XGroupSetID(stream, group, start)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XLen(ctx context.Context, stream string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XLen"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XLen")
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
			w.inflightMetric.WithLabelValues("redis.Client.XLen", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XLen", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XLen", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XLen(stream)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XPending(ctx context.Context, stream string, group string) *redis.XPendingCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XPendingCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XPending"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XPending")
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
			w.inflightMetric.WithLabelValues("redis.Client.XPending", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XPending", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XPending", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XPending(stream, group)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XPendingExt(ctx context.Context, a *redis.XPendingExtArgs) *redis.XPendingExtCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XPendingExtCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XPendingExt"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XPendingExt")
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
			w.inflightMetric.WithLabelValues("redis.Client.XPendingExt", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XPendingExt", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XPendingExt", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XPendingExt(a)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XRange(ctx context.Context, stream string, start string, stop string) *redis.XMessageSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XMessageSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XRange")
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
			w.inflightMetric.WithLabelValues("redis.Client.XRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XRange(stream, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XRangeN(ctx context.Context, stream string, start string, stop string, count int64) *redis.XMessageSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XMessageSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XRangeN"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XRangeN")
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
			w.inflightMetric.WithLabelValues("redis.Client.XRangeN", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XRangeN", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XRangeN", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XRangeN(stream, start, stop, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XRead(ctx context.Context, a *redis.XReadArgs) *redis.XStreamSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XStreamSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XRead"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XRead")
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
			w.inflightMetric.WithLabelValues("redis.Client.XRead", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XRead", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XRead", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XRead(a)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XReadGroup(ctx context.Context, a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XStreamSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XReadGroup"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XReadGroup")
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
			w.inflightMetric.WithLabelValues("redis.Client.XReadGroup", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XReadGroup", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XReadGroup", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XReadGroup(a)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XReadStreams(ctx context.Context, streams ...string) *redis.XStreamSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XStreamSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XReadStreams"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XReadStreams")
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
			w.inflightMetric.WithLabelValues("redis.Client.XReadStreams", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XReadStreams", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XReadStreams", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XReadStreams(streams...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XRevRange(ctx context.Context, stream string, start string, stop string) *redis.XMessageSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XMessageSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XRevRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XRevRange")
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
			w.inflightMetric.WithLabelValues("redis.Client.XRevRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XRevRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XRevRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XRevRange(stream, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XRevRangeN(ctx context.Context, stream string, start string, stop string, count int64) *redis.XMessageSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XMessageSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XRevRangeN"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XRevRangeN")
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
			w.inflightMetric.WithLabelValues("redis.Client.XRevRangeN", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XRevRangeN", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XRevRangeN", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XRevRangeN(stream, start, stop, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XTrim(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XTrim"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XTrim")
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
			w.inflightMetric.WithLabelValues("redis.Client.XTrim", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XTrim", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XTrim", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XTrim(key, maxLen)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) XTrimApprox(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.XTrimApprox"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.XTrimApprox")
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
			w.inflightMetric.WithLabelValues("redis.Client.XTrimApprox", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.XTrimApprox", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.XTrimApprox", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XTrimApprox(key, maxLen)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZAdd"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZAdd")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZAdd", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZAdd", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZAdd", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZAdd(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZAddCh(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZAddCh"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZAddCh")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZAddCh", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZAddCh", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZAddCh", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZAddCh(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZAddNX(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZAddNX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZAddNX")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZAddNX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZAddNX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZAddNX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZAddNX(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZAddNXCh(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZAddNXCh"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZAddNXCh")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZAddNXCh", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZAddNXCh", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZAddNXCh", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZAddNXCh(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZAddXX(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZAddXX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZAddXX")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZAddXX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZAddXX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZAddXX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZAddXX(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZAddXXCh(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZAddXXCh"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZAddXXCh")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZAddXXCh", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZAddXXCh", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZAddXXCh", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZAddXXCh(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZCard(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZCard"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZCard")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZCard", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZCard", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZCard", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZCard(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZCount(ctx context.Context, key string, min string, max string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZCount"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZCount")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZCount", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZCount", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZCount", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZCount(key, min, max)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZIncr(ctx context.Context, key string, member redis.Z) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZIncr"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZIncr")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZIncr", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZIncr", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZIncr", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZIncr(key, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZIncrBy"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZIncrBy")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZIncrBy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZIncrBy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZIncrBy", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZIncrBy(key, increment, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZIncrNX(ctx context.Context, key string, member redis.Z) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZIncrNX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZIncrNX")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZIncrNX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZIncrNX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZIncrNX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZIncrNX(key, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZIncrXX(ctx context.Context, key string, member redis.Z) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZIncrXX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZIncrXX")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZIncrXX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZIncrXX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZIncrXX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZIncrXX(key, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZInterStore(ctx context.Context, destination string, store redis.ZStore, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZInterStore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZInterStore")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZInterStore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZInterStore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZInterStore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZInterStore(destination, store, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZLexCount(ctx context.Context, key string, min string, max string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZLexCount"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZLexCount")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZLexCount", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZLexCount", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZLexCount", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZLexCount(key, min, max)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZPopMax(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZPopMax"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZPopMax")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZPopMax", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZPopMax", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZPopMax", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZPopMax(key, count...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZPopMin(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZPopMin"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZPopMin")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZPopMin", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZPopMin", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZPopMin", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZPopMin(key, count...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRange(ctx context.Context, key string, start int64, stop int64) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRange")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRange(key, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRangeByLex(ctx context.Context, key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRangeByLex"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRangeByLex")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRangeByLex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRangeByLex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRangeByLex", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRangeByLex(key, opt)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRangeByScore(ctx context.Context, key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRangeByScore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRangeByScore")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRangeByScore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRangeByScore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRangeByScore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRangeByScore(key, opt)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRangeByScoreWithScores(ctx context.Context, key string, opt redis.ZRangeBy) *redis.ZSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRangeByScoreWithScores"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRangeByScoreWithScores")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRangeByScoreWithScores", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRangeByScoreWithScores", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRangeByScoreWithScores", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRangeByScoreWithScores(key, opt)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRangeWithScores(ctx context.Context, key string, start int64, stop int64) *redis.ZSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRangeWithScores"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRangeWithScores")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRangeWithScores", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRangeWithScores", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRangeWithScores", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRangeWithScores(key, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRank(ctx context.Context, key string, member string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRank"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRank")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRank", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRank", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRank", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRank(key, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRem"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRem")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRem", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRem", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRem", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRem(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRemRangeByLex(ctx context.Context, key string, min string, max string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRemRangeByLex"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRemRangeByLex")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRemRangeByLex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRemRangeByLex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRemRangeByLex", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRemRangeByLex(key, min, max)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRemRangeByRank(ctx context.Context, key string, start int64, stop int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRemRangeByRank"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRemRangeByRank")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRemRangeByRank", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRemRangeByRank", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRemRangeByRank", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRemRangeByRank(key, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRemRangeByScore(ctx context.Context, key string, min string, max string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRemRangeByScore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRemRangeByScore")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRemRangeByScore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRemRangeByScore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRemRangeByScore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRemRangeByScore(key, min, max)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRevRange(ctx context.Context, key string, start int64, stop int64) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRevRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRevRange")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRevRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRevRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRevRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRevRange(key, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRevRangeByLex(ctx context.Context, key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRevRangeByLex"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRevRangeByLex")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRevRangeByLex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRevRangeByLex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRevRangeByLex", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRevRangeByLex(key, opt)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRevRangeByScore(ctx context.Context, key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRevRangeByScore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRevRangeByScore")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRevRangeByScore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRevRangeByScore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRevRangeByScore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRevRangeByScore(key, opt)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt redis.ZRangeBy) *redis.ZSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRevRangeByScoreWithScores"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRevRangeByScoreWithScores")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRevRangeByScoreWithScores", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRevRangeByScoreWithScores", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRevRangeByScoreWithScores", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRevRangeByScoreWithScores(key, opt)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRevRangeWithScores(ctx context.Context, key string, start int64, stop int64) *redis.ZSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRevRangeWithScores"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRevRangeWithScores")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRevRangeWithScores", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRevRangeWithScores", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRevRangeWithScores", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRevRangeWithScores(key, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZRevRank(ctx context.Context, key string, member string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZRevRank"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZRevRank")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZRevRank", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZRevRank", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZRevRank", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRevRank(key, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ScanCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZScan"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZScan")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZScan", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZScan", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZScan", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZScan(key, cursor, match, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZScore(ctx context.Context, key string, member string) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZScore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZScore")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZScore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZScore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZScore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZScore(key, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClientWrapper) ZUnionStore(ctx context.Context, dest string, store redis.ZStore, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "Client.ZUnionStore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.Client.ZUnionStore")
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
			w.inflightMetric.WithLabelValues("redis.Client.ZUnionStore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.Client.ZUnionStore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.Client.ZUnionStore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZUnionStore(dest, store, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Close(ctx context.Context) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Close"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Close")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Close", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Close", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Close", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Close()
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) Context() context.Context {
	res0 := w.obj.Context()
	return res0
}

func (w *RedisClusterClientWrapper) DBSize(ctx context.Context) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.DBSize"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.DBSize")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.DBSize", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.DBSize", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.DBSize", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.DBSize()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.Cmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Do"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Do")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Do", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Do(args...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ForEachMaster(ctx context.Context, fn func(client *redis.Client) error) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ForEachMaster"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ForEachMaster")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ForEachMaster", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ForEachMaster", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ForEachMaster", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.ForEachMaster(fn)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) ForEachNode(ctx context.Context, fn func(client *redis.Client) error) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ForEachNode"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ForEachNode")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ForEachNode", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ForEachNode", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ForEachNode", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.ForEachNode(fn)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) ForEachSlave(ctx context.Context, fn func(client *redis.Client) error) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ForEachSlave"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ForEachSlave")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ForEachSlave", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ForEachSlave", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ForEachSlave", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.ForEachSlave(fn)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
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
	var res0 []redis.Cmder
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Pipelined"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Pipelined")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Pipelined", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Pipelined", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Pipelined", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Pipelined(fn)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
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
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Process"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Process")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Process", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Process", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Process", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Process(cmd)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return err
}

func (w *RedisClusterClientWrapper) ReloadState(ctx context.Context) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ReloadState"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ReloadState")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ReloadState", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ReloadState", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ReloadState", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.ReloadState()
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
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
	var res0 []redis.Cmder
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.TxPipelined"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.TxPipelined")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.TxPipelined", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.TxPipelined", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.TxPipelined", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.TxPipelined(fn)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *RedisClusterClientWrapper) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	ctxOptions := FromContext(ctx)
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Watch"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Watch")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Watch", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Watch", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Watch", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		err = w.obj.Watch(fn, keys...)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
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

func (w *RedisClusterClientWrapper) Append(ctx context.Context, key string, value string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Append"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Append")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Append", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Append", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Append", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Append(key, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.BLPop"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.BLPop")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.BLPop", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.BLPop", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.BLPop", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BLPop(timeout, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.BRPop"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.BRPop")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.BRPop", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.BRPop", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.BRPop", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BRPop(timeout, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) BRPopLPush(ctx context.Context, source string, destination string, timeout time.Duration) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.BRPopLPush"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.BRPopLPush")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.BRPopLPush", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.BRPopLPush", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.BRPopLPush", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BRPopLPush(source, destination, timeout)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZWithKeyCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.BZPopMax"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.BZPopMax")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.BZPopMax", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.BZPopMax", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.BZPopMax", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BZPopMax(timeout, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZWithKeyCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.BZPopMin"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.BZPopMin")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.BZPopMin", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.BZPopMin", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.BZPopMin", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BZPopMin(timeout, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) BgRewriteAOF(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.BgRewriteAOF"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.BgRewriteAOF")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.BgRewriteAOF", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.BgRewriteAOF", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.BgRewriteAOF", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BgRewriteAOF()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) BgSave(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.BgSave"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.BgSave")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.BgSave", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.BgSave", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.BgSave", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BgSave()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.BitCount"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.BitCount")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.BitCount", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.BitCount", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.BitCount", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BitCount(key, bitCount)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) BitOpAnd(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.BitOpAnd"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.BitOpAnd")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.BitOpAnd", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.BitOpAnd", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.BitOpAnd", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BitOpAnd(destKey, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) BitOpNot(ctx context.Context, destKey string, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.BitOpNot"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.BitOpNot")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.BitOpNot", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.BitOpNot", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.BitOpNot", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BitOpNot(destKey, key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) BitOpOr(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.BitOpOr"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.BitOpOr")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.BitOpOr", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.BitOpOr", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.BitOpOr", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BitOpOr(destKey, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) BitOpXor(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.BitOpXor"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.BitOpXor")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.BitOpXor", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.BitOpXor", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.BitOpXor", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BitOpXor(destKey, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.BitPos"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.BitPos")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.BitPos", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.BitPos", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.BitPos", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.BitPos(key, bit, pos...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClientGetName(ctx context.Context) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClientGetName"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClientGetName")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientGetName", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientGetName", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClientGetName", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientGetName()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClientID(ctx context.Context) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClientID"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClientID")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientID", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientID", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClientID", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientID()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClientKill(ctx context.Context, ipPort string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClientKill"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClientKill")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientKill", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientKill", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClientKill", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientKill(ipPort)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClientKillByFilter(ctx context.Context, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClientKillByFilter"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClientKillByFilter")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientKillByFilter", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientKillByFilter", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClientKillByFilter", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientKillByFilter(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClientList(ctx context.Context) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClientList"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClientList")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientList", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientList", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClientList", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientList()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClientPause(ctx context.Context, dur time.Duration) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClientPause"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClientPause")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientPause", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientPause", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClientPause", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientPause(dur)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClientUnblock(ctx context.Context, id int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClientUnblock"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClientUnblock")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientUnblock", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientUnblock", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClientUnblock", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientUnblock(id)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClientUnblockWithError(ctx context.Context, id int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClientUnblockWithError"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClientUnblockWithError")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientUnblockWithError", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClientUnblockWithError", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClientUnblockWithError", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClientUnblockWithError(id)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterAddSlots(ctx context.Context, slots ...int) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterAddSlots"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterAddSlots")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterAddSlots", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterAddSlots", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterAddSlots", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterAddSlots(slots...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterAddSlotsRange(ctx context.Context, min int, max int) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterAddSlotsRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterAddSlotsRange")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterAddSlotsRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterAddSlotsRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterAddSlotsRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterAddSlotsRange(min, max)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterCountFailureReports(ctx context.Context, nodeID string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterCountFailureReports"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterCountFailureReports")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterCountFailureReports", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterCountFailureReports", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterCountFailureReports", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterCountFailureReports(nodeID)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterCountKeysInSlot(ctx context.Context, slot int) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterCountKeysInSlot"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterCountKeysInSlot")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterCountKeysInSlot", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterCountKeysInSlot", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterCountKeysInSlot", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterCountKeysInSlot(slot)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterDelSlots(ctx context.Context, slots ...int) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterDelSlots"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterDelSlots")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterDelSlots", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterDelSlots", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterDelSlots", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterDelSlots(slots...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterDelSlotsRange(ctx context.Context, min int, max int) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterDelSlotsRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterDelSlotsRange")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterDelSlotsRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterDelSlotsRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterDelSlotsRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterDelSlotsRange(min, max)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterFailover(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterFailover"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterFailover")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterFailover", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterFailover", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterFailover", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterFailover()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterForget(ctx context.Context, nodeID string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterForget"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterForget")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterForget", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterForget", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterForget", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterForget(nodeID)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterGetKeysInSlot"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterGetKeysInSlot")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterGetKeysInSlot", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterGetKeysInSlot", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterGetKeysInSlot", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterGetKeysInSlot(slot, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterInfo(ctx context.Context) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterInfo"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterInfo")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterInfo", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterInfo", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterInfo", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterInfo()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterKeySlot(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterKeySlot"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterKeySlot")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterKeySlot", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterKeySlot", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterKeySlot", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterKeySlot(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterMeet(ctx context.Context, host string, port string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterMeet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterMeet")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterMeet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterMeet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterMeet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterMeet(host, port)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterNodes(ctx context.Context) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterNodes"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterNodes")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterNodes", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterNodes", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterNodes", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterNodes()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterReplicate(ctx context.Context, nodeID string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterReplicate"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterReplicate")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterReplicate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterReplicate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterReplicate", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterReplicate(nodeID)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterResetHard(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterResetHard"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterResetHard")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterResetHard", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterResetHard", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterResetHard", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterResetHard()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterResetSoft(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterResetSoft"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterResetSoft")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterResetSoft", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterResetSoft", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterResetSoft", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterResetSoft()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterSaveConfig(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterSaveConfig"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterSaveConfig")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterSaveConfig", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterSaveConfig", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterSaveConfig", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterSaveConfig()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterSlaves(ctx context.Context, nodeID string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterSlaves"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterSlaves")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterSlaves", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterSlaves", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterSlaves", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterSlaves(nodeID)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ClusterSlots(ctx context.Context) *redis.ClusterSlotsCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ClusterSlotsCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ClusterSlots"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ClusterSlots")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterSlots", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ClusterSlots", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ClusterSlots", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ClusterSlots()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Command(ctx context.Context) *redis.CommandsInfoCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.CommandsInfoCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Command"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Command")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Command", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Command", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Command", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Command()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ConfigGet(ctx context.Context, parameter string) *redis.SliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.SliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ConfigGet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ConfigGet")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ConfigGet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ConfigGet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ConfigGet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ConfigGet(parameter)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ConfigResetStat(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ConfigResetStat"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ConfigResetStat")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ConfigResetStat", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ConfigResetStat", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ConfigResetStat", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ConfigResetStat()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ConfigRewrite(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ConfigRewrite"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ConfigRewrite")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ConfigRewrite", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ConfigRewrite", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ConfigRewrite", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ConfigRewrite()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ConfigSet(ctx context.Context, parameter string, value string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ConfigSet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ConfigSet")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ConfigSet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ConfigSet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ConfigSet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ConfigSet(parameter, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) DbSize(ctx context.Context) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.DbSize"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.DbSize")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.DbSize", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.DbSize", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.DbSize", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.DbSize()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) DebugObject(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.DebugObject"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.DebugObject")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.DebugObject", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.DebugObject", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.DebugObject", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.DebugObject(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Decr(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Decr"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Decr")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Decr", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Decr", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Decr", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Decr(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.DecrBy"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.DecrBy")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.DecrBy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.DecrBy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.DecrBy", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.DecrBy(key, decrement)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Del"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Del")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Del", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Del", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Del", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Del(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Dump(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Dump"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Dump")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Dump", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Dump", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Dump", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Dump(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Echo(ctx context.Context, message interface{}) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Echo"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Echo")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Echo", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Echo", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Echo", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Echo(message)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.Cmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Eval"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Eval")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Eval", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Eval", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Eval", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Eval(script, keys, args...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.Cmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.EvalSha"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.EvalSha")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.EvalSha", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.EvalSha", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.EvalSha", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.EvalSha(sha1, keys, args...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Exists"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Exists")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Exists", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Exists", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Exists", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Exists(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Expire"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Expire")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Expire", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Expire", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Expire", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Expire(key, expiration)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ExpireAt"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ExpireAt")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ExpireAt", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ExpireAt", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ExpireAt", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ExpireAt(key, tm)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) FlushAll(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.FlushAll"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.FlushAll")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.FlushAll", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.FlushAll", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.FlushAll", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.FlushAll()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) FlushAllAsync(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.FlushAllAsync"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.FlushAllAsync")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.FlushAllAsync", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.FlushAllAsync", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.FlushAllAsync", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.FlushAllAsync()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) FlushDB(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.FlushDB"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.FlushDB")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.FlushDB", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.FlushDB", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.FlushDB", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.FlushDB()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) FlushDBAsync(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.FlushDBAsync"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.FlushDBAsync")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.FlushDBAsync", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.FlushDBAsync", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.FlushDBAsync", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.FlushDBAsync()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) FlushDb(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.FlushDb"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.FlushDb")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.FlushDb", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.FlushDb", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.FlushDb", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.FlushDb()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.GeoAdd"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.GeoAdd")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoAdd", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoAdd", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.GeoAdd", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoAdd(key, geoLocation...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) GeoDist(ctx context.Context, key string, member1 string, member2 string, unit string) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.GeoDist"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.GeoDist")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoDist", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoDist", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.GeoDist", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoDist(key, member1, member2, unit)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) GeoHash(ctx context.Context, key string, members ...string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.GeoHash"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.GeoHash")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoHash", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoHash", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.GeoHash", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoHash(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) GeoPos(ctx context.Context, key string, members ...string) *redis.GeoPosCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.GeoPosCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.GeoPos"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.GeoPos")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoPos", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoPos", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.GeoPos", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoPos(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) GeoRadius(ctx context.Context, key string, longitude float64, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.GeoLocationCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.GeoRadius"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.GeoRadius")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoRadius", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoRadius", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.GeoRadius", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoRadius(key, longitude, latitude, query)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) GeoRadiusByMember(ctx context.Context, key string, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.GeoLocationCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.GeoRadiusByMember"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.GeoRadiusByMember")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoRadiusByMember", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoRadiusByMember", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.GeoRadiusByMember", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoRadiusByMember(key, member, query)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) GeoRadiusByMemberRO(ctx context.Context, key string, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.GeoLocationCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.GeoRadiusByMemberRO"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.GeoRadiusByMemberRO")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoRadiusByMemberRO", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoRadiusByMemberRO", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.GeoRadiusByMemberRO", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoRadiusByMemberRO(key, member, query)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) GeoRadiusRO(ctx context.Context, key string, longitude float64, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.GeoLocationCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.GeoRadiusRO"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.GeoRadiusRO")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoRadiusRO", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.GeoRadiusRO", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.GeoRadiusRO", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GeoRadiusRO(key, longitude, latitude, query)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Get(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Get"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Get")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Get", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Get", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Get", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Get(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.GetBit"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.GetBit")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.GetBit", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.GetBit", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.GetBit", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GetBit(key, offset)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) GetRange(ctx context.Context, key string, start int64, end int64) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.GetRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.GetRange")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.GetRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.GetRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.GetRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GetRange(key, start, end)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.GetSet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.GetSet")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.GetSet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.GetSet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.GetSet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.GetSet(key, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.HDel"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.HDel")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.HDel", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.HDel", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.HDel", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HDel(key, fields...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) HExists(ctx context.Context, key string, field string) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.HExists"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.HExists")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.HExists", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.HExists", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.HExists", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HExists(key, field)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) HGet(ctx context.Context, key string, field string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.HGet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.HGet")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.HGet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.HGet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.HGet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HGet(key, field)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringStringMapCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.HGetAll"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.HGetAll")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.HGetAll", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.HGetAll", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.HGetAll", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HGetAll(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) HIncrBy(ctx context.Context, key string, field string, incr int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.HIncrBy"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.HIncrBy")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.HIncrBy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.HIncrBy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.HIncrBy", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HIncrBy(key, field, incr)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) HIncrByFloat(ctx context.Context, key string, field string, incr float64) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.HIncrByFloat"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.HIncrByFloat")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.HIncrByFloat", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.HIncrByFloat", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.HIncrByFloat", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HIncrByFloat(key, field, incr)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.HKeys"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.HKeys")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.HKeys", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.HKeys", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.HKeys", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HKeys(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) HLen(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.HLen"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.HLen")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.HLen", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.HLen", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.HLen", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HLen(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.SliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.HMGet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.HMGet")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.HMGet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.HMGet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.HMGet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HMGet(key, fields...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) HMSet(ctx context.Context, key string, fields map[string]interface{}) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.HMSet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.HMSet")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.HMSet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.HMSet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.HMSet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HMSet(key, fields)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ScanCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.HScan"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.HScan")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.HScan", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.HScan", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.HScan", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HScan(key, cursor, match, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) HSet(ctx context.Context, key string, field string, value interface{}) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.HSet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.HSet")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.HSet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.HSet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.HSet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HSet(key, field, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) HSetNX(ctx context.Context, key string, field string, value interface{}) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.HSetNX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.HSetNX")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.HSetNX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.HSetNX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.HSetNX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HSetNX(key, field, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.HVals"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.HVals")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.HVals", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.HVals", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.HVals", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.HVals(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Incr(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Incr"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Incr")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Incr", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Incr", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Incr", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Incr(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.IncrBy"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.IncrBy")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.IncrBy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.IncrBy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.IncrBy", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.IncrBy(key, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.IncrByFloat"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.IncrByFloat")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.IncrByFloat", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.IncrByFloat", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.IncrByFloat", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.IncrByFloat(key, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Info(ctx context.Context, section ...string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Info"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Info")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Info", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Info", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Info", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Info(section...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Keys"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Keys")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Keys", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Keys", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Keys", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Keys(pattern)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.LIndex"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.LIndex")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.LIndex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.LIndex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.LIndex", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LIndex(key, index)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) LInsert(ctx context.Context, key string, op string, pivot interface{}, value interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.LInsert"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.LInsert")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.LInsert", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.LInsert", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.LInsert", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LInsert(key, op, pivot, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) LInsertAfter(ctx context.Context, key string, pivot interface{}, value interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.LInsertAfter"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.LInsertAfter")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.LInsertAfter", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.LInsertAfter", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.LInsertAfter", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LInsertAfter(key, pivot, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) LInsertBefore(ctx context.Context, key string, pivot interface{}, value interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.LInsertBefore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.LInsertBefore")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.LInsertBefore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.LInsertBefore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.LInsertBefore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LInsertBefore(key, pivot, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) LLen(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.LLen"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.LLen")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.LLen", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.LLen", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.LLen", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LLen(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) LPop(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.LPop"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.LPop")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.LPop", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.LPop", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.LPop", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LPop(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.LPush"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.LPush")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.LPush", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.LPush", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.LPush", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LPush(key, values...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) LPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.LPushX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.LPushX")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.LPushX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.LPushX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.LPushX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LPushX(key, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) LRange(ctx context.Context, key string, start int64, stop int64) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.LRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.LRange")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.LRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.LRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.LRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LRange(key, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.LRem"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.LRem")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.LRem", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.LRem", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.LRem", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LRem(key, count, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.LSet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.LSet")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.LSet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.LSet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.LSet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LSet(key, index, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) LTrim(ctx context.Context, key string, start int64, stop int64) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.LTrim"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.LTrim")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.LTrim", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.LTrim", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.LTrim", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LTrim(key, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) LastSave(ctx context.Context) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.LastSave"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.LastSave")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.LastSave", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.LastSave", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.LastSave", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.LastSave()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.SliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.MGet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.MGet")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.MGet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.MGet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.MGet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.MGet(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) MSet(ctx context.Context, pairs ...interface{}) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.MSet"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.MSet")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.MSet", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.MSet", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.MSet", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.MSet(pairs...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) MSetNX(ctx context.Context, pairs ...interface{}) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.MSetNX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.MSetNX")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.MSetNX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.MSetNX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.MSetNX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.MSetNX(pairs...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) MemoryUsage(ctx context.Context, key string, samples ...int) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.MemoryUsage"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.MemoryUsage")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.MemoryUsage", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.MemoryUsage", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.MemoryUsage", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.MemoryUsage(key, samples...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Migrate(ctx context.Context, host string, port string, key string, db int64, timeout time.Duration) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Migrate"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Migrate")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Migrate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Migrate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Migrate", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Migrate(host, port, key, db, timeout)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Move(ctx context.Context, key string, db int64) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Move"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Move")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Move", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Move", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Move", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Move(key, db)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ObjectEncoding(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ObjectEncoding"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ObjectEncoding")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ObjectEncoding", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ObjectEncoding", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ObjectEncoding", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ObjectEncoding(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ObjectIdleTime(ctx context.Context, key string) *redis.DurationCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.DurationCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ObjectIdleTime"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ObjectIdleTime")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ObjectIdleTime", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ObjectIdleTime", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ObjectIdleTime", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ObjectIdleTime(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ObjectRefCount(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ObjectRefCount"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ObjectRefCount")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ObjectRefCount", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ObjectRefCount", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ObjectRefCount", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ObjectRefCount(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) PExpire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.PExpire"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.PExpire")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.PExpire", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.PExpire", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.PExpire", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PExpire(key, expiration)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) PExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.PExpireAt"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.PExpireAt")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.PExpireAt", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.PExpireAt", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.PExpireAt", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PExpireAt(key, tm)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) PFAdd(ctx context.Context, key string, els ...interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.PFAdd"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.PFAdd")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.PFAdd", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.PFAdd", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.PFAdd", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PFAdd(key, els...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) PFCount(ctx context.Context, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.PFCount"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.PFCount")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.PFCount", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.PFCount", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.PFCount", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PFCount(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) PFMerge(ctx context.Context, dest string, keys ...string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.PFMerge"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.PFMerge")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.PFMerge", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.PFMerge", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.PFMerge", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PFMerge(dest, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) PTTL(ctx context.Context, key string) *redis.DurationCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.DurationCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.PTTL"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.PTTL")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.PTTL", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.PTTL", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.PTTL", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PTTL(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Persist(ctx context.Context, key string) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Persist"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Persist")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Persist", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Persist", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Persist", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Persist(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Ping(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Ping"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Ping")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Ping", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Ping", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Ping", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Ping()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) PubSubChannels(ctx context.Context, pattern string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.PubSubChannels"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.PubSubChannels")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.PubSubChannels", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.PubSubChannels", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.PubSubChannels", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PubSubChannels(pattern)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) PubSubNumPat(ctx context.Context) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.PubSubNumPat"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.PubSubNumPat")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.PubSubNumPat", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.PubSubNumPat", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.PubSubNumPat", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PubSubNumPat()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) PubSubNumSub(ctx context.Context, channels ...string) *redis.StringIntMapCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringIntMapCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.PubSubNumSub"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.PubSubNumSub")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.PubSubNumSub", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.PubSubNumSub", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.PubSubNumSub", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.PubSubNumSub(channels...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Publish"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Publish")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Publish", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Publish", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Publish", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Publish(channel, message)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Quit(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Quit"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Quit")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Quit", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Quit", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Quit", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Quit()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) RPop(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.RPop"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.RPop")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.RPop", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.RPop", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.RPop", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RPop(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) RPopLPush(ctx context.Context, source string, destination string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.RPopLPush"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.RPopLPush")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.RPopLPush", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.RPopLPush", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.RPopLPush", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RPopLPush(source, destination)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.RPush"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.RPush")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.RPush", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.RPush", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.RPush", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RPush(key, values...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) RPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.RPushX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.RPushX")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.RPushX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.RPushX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.RPushX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RPushX(key, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) RandomKey(ctx context.Context) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.RandomKey"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.RandomKey")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.RandomKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.RandomKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.RandomKey", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RandomKey()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ReadOnly(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ReadOnly"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ReadOnly")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ReadOnly", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ReadOnly", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ReadOnly", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ReadOnly()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ReadWrite(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ReadWrite"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ReadWrite")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ReadWrite", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ReadWrite", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ReadWrite", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ReadWrite()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Rename(ctx context.Context, key string, newkey string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Rename"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Rename")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Rename", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Rename", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Rename", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Rename(key, newkey)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) RenameNX(ctx context.Context, key string, newkey string) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.RenameNX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.RenameNX")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.RenameNX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.RenameNX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.RenameNX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RenameNX(key, newkey)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Restore(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Restore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Restore")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Restore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Restore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Restore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Restore(key, ttl, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.RestoreReplace"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.RestoreReplace")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.RestoreReplace", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.RestoreReplace", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.RestoreReplace", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.RestoreReplace(key, ttl, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SAdd"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SAdd")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SAdd", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SAdd", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SAdd", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SAdd(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SCard(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SCard"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SCard")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SCard", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SCard", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SCard", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SCard(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SDiff"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SDiff")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SDiff", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SDiff", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SDiff", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SDiff(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SDiffStore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SDiffStore")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SDiffStore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SDiffStore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SDiffStore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SDiffStore(destination, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SInter(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SInter"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SInter")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SInter", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SInter", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SInter", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SInter(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SInterStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SInterStore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SInterStore")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SInterStore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SInterStore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SInterStore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SInterStore(destination, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SIsMember"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SIsMember")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SIsMember", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SIsMember", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SIsMember", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SIsMember(key, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SMembers"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SMembers")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SMembers", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SMembers", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SMembers", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SMembers(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SMembersMap(ctx context.Context, key string) *redis.StringStructMapCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringStructMapCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SMembersMap"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SMembersMap")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SMembersMap", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SMembersMap", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SMembersMap", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SMembersMap(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SMove(ctx context.Context, source string, destination string, member interface{}) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SMove"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SMove")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SMove", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SMove", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SMove", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SMove(source, destination, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SPop(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SPop"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SPop")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SPop", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SPop", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SPop", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SPop(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SPopN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SPopN"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SPopN")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SPopN", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SPopN", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SPopN", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SPopN(key, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SRandMember(ctx context.Context, key string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SRandMember"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SRandMember")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SRandMember", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SRandMember", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SRandMember", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SRandMember(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SRandMemberN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SRandMemberN"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SRandMemberN")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SRandMemberN", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SRandMemberN", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SRandMemberN", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SRandMemberN(key, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SRem"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SRem")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SRem", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SRem", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SRem", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SRem(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ScanCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SScan"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SScan")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SScan", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SScan", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SScan", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SScan(key, cursor, match, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SUnion(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SUnion"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SUnion")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SUnion", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SUnion", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SUnion", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SUnion(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SUnionStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SUnionStore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SUnionStore")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SUnionStore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SUnionStore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SUnionStore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SUnionStore(destination, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Save(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Save"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Save")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Save", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Save", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Save", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Save()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ScanCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Scan"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Scan")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Scan", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Scan", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Scan", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Scan(cursor, match, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ScriptExists"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ScriptExists")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ScriptExists", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ScriptExists", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ScriptExists", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ScriptExists(hashes...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ScriptFlush(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ScriptFlush"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ScriptFlush")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ScriptFlush", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ScriptFlush", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ScriptFlush", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ScriptFlush()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ScriptKill(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ScriptKill"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ScriptKill")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ScriptKill", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ScriptKill", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ScriptKill", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ScriptKill()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ScriptLoad(ctx context.Context, script string) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ScriptLoad"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ScriptLoad")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ScriptLoad", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ScriptLoad", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ScriptLoad", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ScriptLoad(script)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Set"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Set")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Set", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Set", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Set", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Set(key, value, expiration)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SetBit"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SetBit")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SetBit", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SetBit", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SetBit", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SetBit(key, offset, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SetNX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SetNX")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SetNX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SetNX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SetNX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SetNX(key, value, expiration)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SetRange(ctx context.Context, key string, offset int64, value string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SetRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SetRange")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SetRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SetRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SetRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SetRange(key, offset, value)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.BoolCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SetXX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SetXX")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SetXX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SetXX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SetXX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SetXX(key, value, expiration)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Shutdown(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Shutdown"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Shutdown")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Shutdown", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Shutdown", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Shutdown", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Shutdown()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ShutdownNoSave(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ShutdownNoSave"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ShutdownNoSave")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ShutdownNoSave", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ShutdownNoSave", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ShutdownNoSave", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ShutdownNoSave()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ShutdownSave(ctx context.Context) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ShutdownSave"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ShutdownSave")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ShutdownSave", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ShutdownSave", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ShutdownSave", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ShutdownSave()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SlaveOf(ctx context.Context, host string, port string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SlaveOf"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SlaveOf")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SlaveOf", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SlaveOf", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SlaveOf", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SlaveOf(host, port)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SlowLog() {
	w.obj.SlowLog()
}

func (w *RedisClusterClientWrapper) Sort(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Sort"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Sort")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Sort", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Sort", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Sort", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Sort(key, sort)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) *redis.SliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.SliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SortInterfaces"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SortInterfaces")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SortInterfaces", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SortInterfaces", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SortInterfaces", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SortInterfaces(key, sort)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) SortStore(ctx context.Context, key string, store string, sort *redis.Sort) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.SortStore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.SortStore")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.SortStore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.SortStore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.SortStore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.SortStore(key, store, sort)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) StrLen(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.StrLen"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.StrLen")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.StrLen", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.StrLen", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.StrLen", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.StrLen(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Sync() {
	w.obj.Sync()
}

func (w *RedisClusterClientWrapper) TTL(ctx context.Context, key string) *redis.DurationCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.DurationCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.TTL"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.TTL")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.TTL", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.TTL", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.TTL", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.TTL(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Time(ctx context.Context) *redis.TimeCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.TimeCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Time"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Time")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Time", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Time", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Time", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Time()
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Touch(ctx context.Context, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Touch"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Touch")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Touch", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Touch", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Touch", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Touch(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Type(ctx context.Context, key string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Type"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Type")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Type", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Type", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Type", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Type(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Unlink(ctx context.Context, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Unlink"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Unlink")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Unlink", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Unlink", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Unlink", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Unlink(keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) Wait(ctx context.Context, numSlaves int, timeout time.Duration) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.Wait"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.Wait")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.Wait", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.Wait", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.Wait", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.Wait(numSlaves, timeout)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XAck(ctx context.Context, stream string, group string, ids ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XAck"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XAck")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XAck", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XAck", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XAck", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XAck(stream, group, ids...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XAdd"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XAdd")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XAdd", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XAdd", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XAdd", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XAdd(a)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XClaim(ctx context.Context, a *redis.XClaimArgs) *redis.XMessageSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XMessageSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XClaim"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XClaim")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XClaim", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XClaim", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XClaim", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XClaim(a)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XClaimJustID(ctx context.Context, a *redis.XClaimArgs) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XClaimJustID"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XClaimJustID")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XClaimJustID", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XClaimJustID", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XClaimJustID", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XClaimJustID(a)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XDel(ctx context.Context, stream string, ids ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XDel"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XDel")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XDel", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XDel", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XDel", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XDel(stream, ids...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XGroupCreate(ctx context.Context, stream string, group string, start string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XGroupCreate"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XGroupCreate")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XGroupCreate", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XGroupCreate", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XGroupCreate", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XGroupCreate(stream, group, start)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XGroupCreateMkStream(ctx context.Context, stream string, group string, start string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XGroupCreateMkStream"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XGroupCreateMkStream")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XGroupCreateMkStream", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XGroupCreateMkStream", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XGroupCreateMkStream", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XGroupCreateMkStream(stream, group, start)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XGroupDelConsumer(ctx context.Context, stream string, group string, consumer string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XGroupDelConsumer"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XGroupDelConsumer")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XGroupDelConsumer", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XGroupDelConsumer", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XGroupDelConsumer", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XGroupDelConsumer(stream, group, consumer)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XGroupDestroy(ctx context.Context, stream string, group string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XGroupDestroy"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XGroupDestroy")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XGroupDestroy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XGroupDestroy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XGroupDestroy", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XGroupDestroy(stream, group)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XGroupSetID(ctx context.Context, stream string, group string, start string) *redis.StatusCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StatusCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XGroupSetID"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XGroupSetID")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XGroupSetID", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XGroupSetID", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XGroupSetID", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XGroupSetID(stream, group, start)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XLen(ctx context.Context, stream string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XLen"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XLen")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XLen", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XLen", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XLen", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XLen(stream)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XPending(ctx context.Context, stream string, group string) *redis.XPendingCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XPendingCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XPending"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XPending")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XPending", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XPending", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XPending", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XPending(stream, group)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XPendingExt(ctx context.Context, a *redis.XPendingExtArgs) *redis.XPendingExtCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XPendingExtCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XPendingExt"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XPendingExt")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XPendingExt", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XPendingExt", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XPendingExt", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XPendingExt(a)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XRange(ctx context.Context, stream string, start string, stop string) *redis.XMessageSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XMessageSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XRange")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XRange(stream, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XRangeN(ctx context.Context, stream string, start string, stop string, count int64) *redis.XMessageSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XMessageSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XRangeN"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XRangeN")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XRangeN", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XRangeN", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XRangeN", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XRangeN(stream, start, stop, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XRead(ctx context.Context, a *redis.XReadArgs) *redis.XStreamSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XStreamSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XRead"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XRead")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XRead", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XRead", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XRead", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XRead(a)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XReadGroup(ctx context.Context, a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XStreamSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XReadGroup"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XReadGroup")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XReadGroup", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XReadGroup", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XReadGroup", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XReadGroup(a)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XReadStreams(ctx context.Context, streams ...string) *redis.XStreamSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XStreamSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XReadStreams"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XReadStreams")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XReadStreams", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XReadStreams", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XReadStreams", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XReadStreams(streams...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XRevRange(ctx context.Context, stream string, start string, stop string) *redis.XMessageSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XMessageSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XRevRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XRevRange")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XRevRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XRevRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XRevRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XRevRange(stream, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XRevRangeN(ctx context.Context, stream string, start string, stop string, count int64) *redis.XMessageSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.XMessageSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XRevRangeN"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XRevRangeN")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XRevRangeN", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XRevRangeN", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XRevRangeN", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XRevRangeN(stream, start, stop, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XTrim(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XTrim"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XTrim")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XTrim", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XTrim", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XTrim", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XTrim(key, maxLen)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) XTrimApprox(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.XTrimApprox"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.XTrimApprox")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.XTrimApprox", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.XTrimApprox", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.XTrimApprox", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.XTrimApprox(key, maxLen)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZAdd"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZAdd")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZAdd", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZAdd", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZAdd", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZAdd(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZAddCh(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZAddCh"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZAddCh")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZAddCh", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZAddCh", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZAddCh", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZAddCh(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZAddNX(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZAddNX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZAddNX")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZAddNX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZAddNX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZAddNX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZAddNX(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZAddNXCh(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZAddNXCh"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZAddNXCh")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZAddNXCh", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZAddNXCh", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZAddNXCh", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZAddNXCh(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZAddXX(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZAddXX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZAddXX")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZAddXX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZAddXX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZAddXX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZAddXX(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZAddXXCh(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZAddXXCh"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZAddXXCh")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZAddXXCh", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZAddXXCh", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZAddXXCh", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZAddXXCh(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZCard(ctx context.Context, key string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZCard"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZCard")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZCard", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZCard", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZCard", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZCard(key)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZCount(ctx context.Context, key string, min string, max string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZCount"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZCount")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZCount", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZCount", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZCount", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZCount(key, min, max)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZIncr(ctx context.Context, key string, member redis.Z) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZIncr"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZIncr")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZIncr", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZIncr", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZIncr", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZIncr(key, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZIncrBy"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZIncrBy")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZIncrBy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZIncrBy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZIncrBy", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZIncrBy(key, increment, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZIncrNX(ctx context.Context, key string, member redis.Z) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZIncrNX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZIncrNX")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZIncrNX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZIncrNX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZIncrNX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZIncrNX(key, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZIncrXX(ctx context.Context, key string, member redis.Z) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZIncrXX"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZIncrXX")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZIncrXX", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZIncrXX", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZIncrXX", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZIncrXX(key, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZInterStore(ctx context.Context, destination string, store redis.ZStore, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZInterStore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZInterStore")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZInterStore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZInterStore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZInterStore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZInterStore(destination, store, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZLexCount(ctx context.Context, key string, min string, max string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZLexCount"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZLexCount")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZLexCount", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZLexCount", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZLexCount", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZLexCount(key, min, max)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZPopMax(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZPopMax"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZPopMax")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZPopMax", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZPopMax", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZPopMax", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZPopMax(key, count...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZPopMin(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZPopMin"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZPopMin")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZPopMin", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZPopMin", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZPopMin", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZPopMin(key, count...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRange(ctx context.Context, key string, start int64, stop int64) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRange")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRange(key, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRangeByLex(ctx context.Context, key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRangeByLex"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRangeByLex")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRangeByLex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRangeByLex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRangeByLex", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRangeByLex(key, opt)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRangeByScore(ctx context.Context, key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRangeByScore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRangeByScore")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRangeByScore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRangeByScore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRangeByScore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRangeByScore(key, opt)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRangeByScoreWithScores(ctx context.Context, key string, opt redis.ZRangeBy) *redis.ZSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRangeByScoreWithScores"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRangeByScoreWithScores")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRangeByScoreWithScores", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRangeByScoreWithScores", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRangeByScoreWithScores", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRangeByScoreWithScores(key, opt)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRangeWithScores(ctx context.Context, key string, start int64, stop int64) *redis.ZSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRangeWithScores"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRangeWithScores")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRangeWithScores", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRangeWithScores", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRangeWithScores", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRangeWithScores(key, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRank(ctx context.Context, key string, member string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRank"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRank")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRank", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRank", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRank", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRank(key, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRem"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRem")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRem", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRem", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRem", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRem(key, members...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRemRangeByLex(ctx context.Context, key string, min string, max string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRemRangeByLex"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRemRangeByLex")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRemRangeByLex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRemRangeByLex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRemRangeByLex", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRemRangeByLex(key, min, max)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRemRangeByRank(ctx context.Context, key string, start int64, stop int64) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRemRangeByRank"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRemRangeByRank")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRemRangeByRank", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRemRangeByRank", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRemRangeByRank", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRemRangeByRank(key, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRemRangeByScore(ctx context.Context, key string, min string, max string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRemRangeByScore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRemRangeByScore")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRemRangeByScore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRemRangeByScore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRemRangeByScore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRemRangeByScore(key, min, max)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRevRange(ctx context.Context, key string, start int64, stop int64) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRevRange"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRevRange")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRevRange", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRevRange", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRevRange", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRevRange(key, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRevRangeByLex(ctx context.Context, key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRevRangeByLex"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRevRangeByLex")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRevRangeByLex", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRevRangeByLex", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRevRangeByLex", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRevRangeByLex(key, opt)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRevRangeByScore(ctx context.Context, key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.StringSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRevRangeByScore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRevRangeByScore")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRevRangeByScore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRevRangeByScore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRevRangeByScore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRevRangeByScore(key, opt)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt redis.ZRangeBy) *redis.ZSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRevRangeByScoreWithScores"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRevRangeByScoreWithScores")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRevRangeByScoreWithScores", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRevRangeByScoreWithScores", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRevRangeByScoreWithScores", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRevRangeByScoreWithScores(key, opt)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRevRangeWithScores(ctx context.Context, key string, start int64, stop int64) *redis.ZSliceCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ZSliceCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRevRangeWithScores"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRevRangeWithScores")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRevRangeWithScores", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRevRangeWithScores", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRevRangeWithScores", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRevRangeWithScores(key, start, stop)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZRevRank(ctx context.Context, key string, member string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZRevRank"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZRevRank")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRevRank", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZRevRank", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZRevRank", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZRevRank(key, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.ScanCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZScan"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZScan")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZScan", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZScan", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZScan", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZScan(key, cursor, match, count)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZScore(ctx context.Context, key string, member string) *redis.FloatCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.FloatCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZScore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZScore")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZScore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZScore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZScore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZScore(key, member)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}

func (w *RedisClusterClientWrapper) ZUnionStore(ctx context.Context, dest string, store redis.ZStore, keys ...string) *redis.IntCmd {
	ctxOptions := FromContext(ctx)
	var res0 *redis.IntCmd
	_ = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, "ClusterClient.ZUnionStore"); err != nil {
				return err
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "redis.ClusterClient.ZUnionStore")
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
			w.inflightMetric.WithLabelValues("redis.ClusterClient.ZUnionStore", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("redis.ClusterClient.ZUnionStore", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("redis.ClusterClient.ZUnionStore", ErrCode(res0.Err()), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0 = w.obj.ZUnionStore(dest, store, keys...)
		if res0.Err() != nil && span != nil {
			span.SetTag("error", res0.Err().Error())
		}
		return res0.Err()
	})
	return res0
}
