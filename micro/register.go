package micro

func init() {
	RegisterRateLimiter("LocalGroup", NewLocalGroupRateLimiterWithOptions)
	RegisterRateLimiter("LocalShare", NewLocalShareRateLimiterWithOptions)

	RegisterParallelController("Local", NewLocalParallelControllerWithOptions)
}
