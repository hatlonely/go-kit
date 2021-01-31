package ratelimiter

import (
	"context"
)

// https://konghq.com/blog/how-to-design-a-scalable-rate-limiting-algorithm/

type RateLimiter interface {
	Wait(ctx context.Context) error
	WaitN(ctx context.Context, n int) (err error)
	Allow() bool
}
