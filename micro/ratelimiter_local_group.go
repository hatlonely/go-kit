package micro

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

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

func (r *LocalGroupRateLimiter) Allow(ctx context.Context, key string) error {
	limiter, ok := r.limiters[key]
	if !ok {
		return nil
	}
	if !limiter.Allow() {
		return ErrFlowControl
	}
	return nil
}

func (r *LocalGroupRateLimiter) Wait(ctx context.Context, key string) error {
	return r.WaitN(ctx, key, 1)
}

func (r *LocalGroupRateLimiter) WaitN(ctx context.Context, key string, n int) error {
	limiter, ok := r.limiters[key]
	if !ok {
		return nil
	}
	return limiter.WaitN(ctx, n)
}
