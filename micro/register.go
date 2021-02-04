package micro

func init() {
	RegisterRateLimiter("LocalGroup", NewLocalGroupRateLimiterWithOptions)
	RegisterRateLimiter("LocalShare", NewLocalShareRateLimiterWithOptions)

	RegisterParallelControllerGroup("LocalGroup", NewLocalParallelControllerGroupWithOptions)
}
