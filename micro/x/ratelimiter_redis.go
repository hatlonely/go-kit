package microx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/wrap"
)

type RedisRateLimiterOptions struct {
	Redis wrap.RedisClientWrapperOptions
	// key 前缀，可当成命名空间使用
	Prefix string
	// QPS 计算规则
	// 1. key 在 map 中，直接返回 key 对应的 qps
	// 2. key 按 '|' 分割，第 0 个字符串作为 key，如果在 map 中，返回 qps
	// 3. 返回 DefaultQPS
	QPS map[string]int
	// QPS 中未匹配到，使用 DefaultQPS，默认为 0，不限流
	DefaultQPS int
}

type RedisRateLimiter struct {
	client  *wrap.RedisClientWrapper
	options *RedisRateLimiterOptions
}

// https://redislabs.com/redis-best-practices/basic-rate-limiting/
// https://redis.io/commands/INCR#pattern-rate-limiter-1

func NewRedisRateLimiterWithConfig(cfg *config.Config, opts ...refx.Option) (*RedisRateLimiter, error) {
	var options RedisRateLimiterOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.WithMessage(err, "cfg.Unmarshal failed")
	}

	refxOptions := refx.NewOptions(opts...)
	client, err := wrap.NewRedisClientWrapperWithConfig(cfg.Sub(refxOptions.FormatKey("Redis")), opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "NewRedisRateLimiterWithConfig failed")
	}
	r := &RedisRateLimiter{client: client, options: &options}

	cfg.AddOnChangeHandler(func(cfg *config.Config) error {
		var options RedisRateLimiterOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.WithMessage(err, "cfg.Unmarshal failed")
		}
		r.options = &options
		return nil
	})

	return r, nil
}

func NewRedisRateLimiterWithOptions(options *RedisRateLimiterOptions) (*RedisRateLimiter, error) {
	client, err := wrap.NewRedisClientWrapperWithOptions(&options.Redis)
	if err != nil {
		return nil, errors.Wrap(err, "NewRedisClientWrapperWithOptions failed")
	}

	return &RedisRateLimiter{client: client, options: options}, nil
}

func (r *RedisRateLimiter) Allow(ctx context.Context, key string) error {
	qps := r.calculateQPS(key)
	if qps == 0 {
		return nil
	}

	now := time.Now()
	tsKey := fmt.Sprintf("%s_%d", r.generateKey(key), now.Unix())
	var res *redis.IntCmd
	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		res = pipe.Incr(tsKey)
		pipe.Expire(tsKey, 3*time.Second)
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "query redis failed")
	}

	if int(res.Val()) > qps {
		return micro.ErrFlowControl
	}

	return nil
}

func (r *RedisRateLimiter) Wait(ctx context.Context, key string) error {
	return r.WaitN(ctx, key, 1)
}

func (r *RedisRateLimiter) WaitN(ctx context.Context, key string, n int) error {
	qps := r.calculateQPS(key)
	if qps == 0 {
		return nil
	}

	for {
		now := time.Now()
		tsKey := fmt.Sprintf("%s_%d", r.generateKey(key), now.Unix())
		var res *redis.IntCmd
		_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			res = pipe.IncrBy(tsKey, int64(n))
			pipe.Expire(tsKey, 3*time.Second)
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "r.client.TxPipelined failed")
		}
		val := int(res.Val())
		if val <= qps {
			break
		}
		if val-qps <= n {
			n = val - qps
		}

		d := now.Add(time.Second).Sub(time.Now())
		if d > 0 {
			select {
			case <-time.After(d):
			case <-ctx.Done():
				return errors.New("cancel by ctx.Done")
			}
		}
	}
	return nil
}

func (r *RedisRateLimiter) generateKey(key string) string {
	if r.options.Prefix == "" {
		return key
	}
	return fmt.Sprintf("%s_%s", r.options.Prefix, key)
}

func (r *RedisRateLimiter) calculateQPS(key string) int {
	if qps, ok := r.options.QPS[key]; ok {
		return qps
	}

	if idx := strings.Index(key, "|"); idx != -1 {
		if qps, ok := r.options.QPS[key[:idx]]; ok {
			return qps
		}
	}

	return r.options.DefaultQPS
}
