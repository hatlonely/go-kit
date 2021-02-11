package microx

import (
	"github.com/hatlonely/go-kit/micro"
)

func init() {
	micro.RegisterRateLimiter("Redis", NewRedisRateLimiterWithOptions)

	micro.RegisterParallelController("RedisIncr", NewRedisParallelControllerWithOptions)

	micro.RegisterLocker("Redis", NewRedisLockerWithOptions)
}
