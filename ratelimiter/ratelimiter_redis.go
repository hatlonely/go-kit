package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/wrap"
)

type RedisRateLimiterOptions struct {
	Redis wrap.RedisClientWrapperOptions
	QPS   map[string]int
}

type RedisRateLimiter struct {
	client  *wrap.RedisClientWrapper
	options *RedisRateLimiterOptions
}

// https://redislabs.com/redis-best-practices/basic-rate-limiting/
// https://redis.io/commands/INCR#pattern-rate-limiter-1

func NewRedisRateLimiterWithOptions(options *RedisRateLimiterOptions) (*RedisRateLimiter, error) {
	client, err := wrap.NewRedisClientWrapperWithOptions(&options.Redis)
	if err != nil {
		return nil, errors.Wrap(err, "NewRedisClientWrapperWithOptions failed")
	}

	return &RedisRateLimiter{client: client, options: options}, nil
}

func (r *RedisRateLimiter) Allow(ctx context.Context, key string) bool {
	qps, ok := r.options.QPS[key]
	if !ok {
		return true
	}

	now := time.Now()
	tsKey := fmt.Sprintf("%s_%d", key, now.Unix())
	var res *redis.IntCmd
	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		res = pipe.Incr(tsKey)
		pipe.Expire(tsKey, 3*time.Second)
		return nil
	})
	if err != nil {
		return true
	}

	return int(res.Val()) <= qps
}

func (r *RedisRateLimiter) Wait(ctx context.Context, key string) error {
	return r.WaitN(ctx, key, 1)
}

func (r *RedisRateLimiter) WaitN(ctx context.Context, key string, n int) error {
	qps, ok := r.options.QPS[key]
	if !ok {
		return nil
	}

	for {
		now := time.Now()
		tsKey := fmt.Sprintf("%s_%d", key, now.Unix())
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
