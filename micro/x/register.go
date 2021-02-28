package microx

import (
	"github.com/hatlonely/go-kit/micro"
)

func init() {
	micro.RegisterRateLimiter("Redis", NewRedisRateLimiterWithOptions)
	micro.RegisterRateLimiter("OTS", NewOTSRateLimiterWithOptions)

	micro.RegisterParallelController("Redis", NewRedisParallelControllerWithOptions)
	micro.RegisterParallelController("RedisTimed", NewRedisTimedParallelControllerWithOptions)
	micro.RegisterParallelController("OTS", NewOTSParallelControllerWithOptions)

	micro.RegisterLocker("Redis", NewRedisLockerWithOptions)
}
