package micro

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

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

func (r *LocalShareRateLimiter) Allow(ctx context.Context, key string) error {
	cost, ok := r.cost[key]
	if !ok {
		return nil
	}

	if r.limiter.AllowN(time.Now(), cost) {
		return nil
	}

	return ErrFlowControl
}

func (r *LocalShareRateLimiter) Wait(ctx context.Context, key string) error {
	return r.WaitN(ctx, key, 1)
}

func (r *LocalShareRateLimiter) WaitN(ctx context.Context, key string, n int) error {
	cost, ok := r.cost[key]
	if !ok {
		return nil
	}

	return r.limiter.WaitN(ctx, cost*n)
}
