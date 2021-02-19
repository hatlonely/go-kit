# micro

micro 定义了通用的微服务组件接口，提供了一些单机版的实现，`x` 中提供了一些拓展的实现

## Retry

提供功能丰富的重试功能，对外仅提供一个接口 `Do(fun func() error) error`

基于 `github.com/avast/retry-go`，封装成 `Retry` 对象，支持从配置中加载，可作为函数参数传递

```go
type RetryOptions struct {
	Attempts uint `dft:"3"`
	// 重试间隔，当 DelayType == Fixed || DelayType == BackOff 时有效
	Delay time.Duration `dft:"1s"`
	// 最大重试间隔时间
	MaxDelay time.Duration `dft:"3m"`
	// 随机重试间隔，当 DelayType == Random 时有效
	MaxJitter     time.Duration `dft:"1s"`
	LastErrorOnly bool
	// Fixed: 固定重试间隔为 Delay
	// BackOff: 以 Delay 为基础，重试间隔指数增长，Delay, Delay >> 1, Delay >> 2
	// Random: 重试间隔为 (0, MaxJitter) 中的一个随机值
	// Fixed,Random: 重试间隔 (Delay, Delay + MaxJitter) 中的一个随机值
	DelayType string `dft:"BackOff"`
	// 使用定义在 retryRetryIfMap 中的 retryIf 方法
	// 默认重试非 retry.Unrecoverable 的所有错误
	RetryIf string
}
```

使用示例

```go
retry, err := micro.NewRetryWithOptions(&RetryOptions{
	Attempts: 3,
})
if err != nil {
	panic(err)
}

var res *Response
retry.Do(func() error {
	res, err = doSomething()
	if err != nil {
		return err
	}
	return nil
})
```

## RateLimiter

限流器，用来控制一段时间内的请求数量，主要用于以下几种使用场景

- 客户端使用，控制访问后端服务的频次，以保证流量异常增长时不会对后端服务造成影响
- 服务端使用，控制客户端访问自身的频次，以此来保护自身服务，
- 服务端使用，限制某个用户的流量

```go
type RateLimiter interface {
	// 立即返回是否限流
	// 1. 成功返回 nil
	// 2. 限流返回 ErrFlowControl
	// 3. 未知错误返回 err
	Allow(ctx context.Context, key string) error
	// 阻塞等待限流结果返回，或者立即返回错误
	// 1. 成功返回 nil
	// 2. 未知错误返回 err，比如网络错误，不可能返回 ErrFlowControl
	Wait(ctx context.Context, key string) error
	// 阻塞等待 n 个限流结果，或者立即返回错误
	// 1. 成功返回 nil
	// 2. 未知错误返回 err，比如网络错误，不可能返回 ErrFlowControl
	WaitN(ctx context.Context, key string, n int) (err error)
}
```

基于 `golang.org/x/time/rate` 默认提供两种本地限流策略，

- `LocalGroupRateLimiter`: 每个方法都有单独的限流器，互不干扰
- `LocalShareRateLimiter`: 所有方法共享同一个限流器，不同的方法可以设置每次请求消耗的 token 数量

使用示例

```go
r, err := NewRateLimiterWithOptions(&RateLimiterOptions{
Type: "LocalGroup",
	Options: &LocalGroupRateLimiterOptions{
		"key": {
			Interval: time.Second,
			Burst:    2,
		},
	},
})

_ = r.Wait(context.Background(), "key")
doSomething()
```

可以通过 `RegisterRateLimiter` 拓展自己的 `RateLimiter` 实现

## ParallelController

并发控制器，和限流器类似，但是没有时间概念，控制的是正在处理的请求并发数量，token 需要通过 `Release` 方法返还，也有几种使用场景

- 客户端使用，控制访问后端的并发，以保证不会因为并发太大将后端服务压垮
- 服务端使用，用来控制客户端访问自身的并发数量
- 服务端使用，用来控制某个用户的并发数量

```go
type ParallelController interface {
	// 尝试获取 token，如果已空，返回 ErrParallelControl
	// 1. 成功返回 nil
	// 2. 已空，返回 ErrParallelControl
	// 3. 其他未知错误返回 err
	TryAcquire(ctx context.Context, key string) (int, error)
	// 获取 token，如果已空，阻塞等待
	// 1. 成功返回 nil
	// 2. 未知错误返回 err
	Acquire(ctx context.Context, key string) (int, error)
	// 返还 token，如果已满，直接丢弃
	// 1. 成功返回 nil
	// 2. 未知错误返回 err
	Release(ctx context.Context, key string, token int) error
}
```

默认提供两种本地的并发控制器

- `LocalSemaphoreParallelController`: 基于信号量实现的并发控制器
- `LocalTimedSemaphoreParallelController`：使用 treemap 实现带超时机制的并发控制器

使用示例

```go
r, err := NewParallelControllerWithOptions(&ParallelControllerOptions{
	Type: "LocalSemaphore",
	Options: &LocalSemaphoreParallelControllerOptions{
		"key": 2,
	},
})

token, err := r.Acquire(context.Background(), "key")
if err != nil {
    panic(err)
}
doSomething()
r.Release(context.Background(), "key", token)
```

可以通过 `RegisterParallelController` 拓展自己的 `ParallelController` 实现

## Locker

锁

```go
type Locker interface {
	// 尝试获取锁
	// 1. 获取成功，返回 nil
	// 2. 获取失败，返回 ErrLocked
	// 3. 其他错误，返回 err
	TryLock(ctx context.Context, key string) error
	// 获取锁
	// 1. 获取成功，返回 nil
	// 2. 其他错误，返回 err
	Lock(ctx context.Context, key string) error
	// 释放锁
	// 1. 释放成功，返回 nil
	// 2. 其他错误，返回 err
	Unlock(ctx context.Context, key string) error
}
```

默认提供基于信号量的本地实现，`LocalLocker`
