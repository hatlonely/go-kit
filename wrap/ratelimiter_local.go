package wrap

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

func init() {
	RegisterRateLimiterGroup("LocalGroup", NewLocalGroupRateLimiterWithOptions)
	RegisterRateLimiterGroup("ShareGroup", NewLocalShareRateLimiterWithOptions)
}

type LocalGroupRateLimiter struct {
	limiters map[string]*rate.Limiter
}

type LocalGroupRateLimiterOptions map[string]struct {
	Interval time.Duration `rule:"x > 0"`
	Burst    int           `rule:"x > 0"`
}

func NewLocalGroupRateLimiterWithOptions(options *LocalGroupRateLimiterOptions) (*LocalGroupRateLimiter, error) {
	limiters := map[string]*rate.Limiter{}

	for key, val := range *options {
		if val.Interval <= 0 {
			return nil, errors.Errorf("interval should be positive")
		}
		if val.Burst <= 0 {
			return nil, errors.Errorf("bucket should be positive")
		}
		limiters[key] = rate.NewLimiter(rate.Every(val.Interval), val.Burst)
	}

	return &LocalGroupRateLimiter{limiters: limiters}, nil
}

func (r *LocalGroupRateLimiter) Wait(ctx context.Context, key string) error {
	limiter, ok := r.limiters[key]
	if !ok {
		return nil
	}

	return limiter.Wait(ctx)
}

type LocalShareRateLimiterOptions struct {
	Interval time.Duration `rule:"x > 0"`
	Burst    int           `rule:"x > 0"`
	Cost     map[string]int
}

type LocalShareRateLimiter struct {
	limiter *rate.Limiter

	cost map[string]int
}

func NewLocalShareRateLimiterWithOptions(options *LocalShareRateLimiterOptions) (*LocalShareRateLimiter, error) {
	limiter := rate.NewLimiter(rate.Every(options.Interval), options.Burst)

	for _, cost := range options.Cost {
		if cost <= 0 {
			return nil, errors.Errorf("cost should be positive")
		}
	}

	return &LocalShareRateLimiter{
		limiter: limiter,
		cost:    options.Cost,
	}, nil
}

func (r *LocalShareRateLimiter) Wait(ctx context.Context, key string) error {
	cost, ok := r.cost[key]
	if !ok {
		return nil
	}

	return r.limiter.WaitN(ctx, cost)
}
