package wrap

import (
	"context"
	"time"

	"github.com/go-errors/errors"
	"golang.org/x/time/rate"
)

type RateLimiterGroup interface {
	Wait(ctx context.Context, key string) error
}

func NewRateLimiterGroup(options *RateLimiterGroupOptions) (RateLimiterGroup, error) {
	switch options.Type {
	case "":
		return nil, nil
	case "LocalGroup":
		return NewLocalGroupRateLimiter(&options.LocalGroupRateLimiter)
	case "LocalShare":
		return NewLocalShareRateLimiter(&options.LocalShareRateLimiter)
	}
	return nil, errors.Errorf("unsupported rate limiter type: [%v]", options.Type)
}

type RateLimiterGroupOptions struct {
	Type                  string
	LocalGroupRateLimiter LocalGroupRateLimiterOptions
	LocalShareRateLimiter LocalShareRateLimiterOptions
}

type LocalGroupRateLimiter struct {
	limiters map[string]*rate.Limiter
}

type LocalGroupRateLimiterOptions map[string]struct {
	Interval time.Duration `rule:"x > 0"`
	Burst    int           `rule:"x > 0"`
}

func NewLocalGroupRateLimiter(options *LocalGroupRateLimiterOptions) (*LocalGroupRateLimiter, error) {
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

func NewLocalShareRateLimiter(options *LocalShareRateLimiterOptions) (*LocalShareRateLimiter, error) {
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
