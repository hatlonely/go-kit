# micro/x

`micro/x` 提供了 micro 组件的拓展实现

## RedisRateLimiter

基于 Redis 的 `RateLimiter` 实现，方案参考 <https://redislabs.com/redis-best-practices/basic-rate-limiting/>

使用时间戳作为 key，获取 token 使用 `incrby` 命令

这个方案中时间戳来自于客户端，需要客户端设置时钟同步，否则在压力较大时，时钟较慢的机器会一直获取不到 token

```go
type RedisRateLimiterOptions struct {
	Redis wrap.RedisClientWrapperOptions
	// 限流窗口长度
	Window time.Duration `dft:"1s"`
	// key 前缀，可当成命名空间使用
	Prefix string
	// QPS 计算规则
	// 1. key 在 map 中，直接返回 key 对应的 qps
	// 2. key 按 '|' 分割，第 0 个字符串作为 key，如果在 map 中，返回 qps
	// 3. 返回 DefaultQPS
	QPS map[string]int
	// QPS 中未匹配到，使用 DefaultQPS，默认为 0，不限流
	DefaultQPS int
}
```

## OTSRateLimiter

基于 OTS 的 `RateLimiter` 实现，和 `RedisRateLimiter` 类似，使用 `UpdateRow` 中的条件增量更新

```go
type OTSRateLimiterOptions struct {
	OTS wrap.OTSTableStoreClientWrapperOptions
	// OTS 表名
	Table string
	// 限流窗口长度
	Window time.Duration `dft:"1s"`
	// key 前缀，可当成命名空间使用
	Prefix string
	// QPS 计算规则
	// 1. key 在 map 中，直接返回 key 对应的 qps
	// 2. key 按 '|' 分割，第 0 个字符串作为 key，如果在 map 中，返回 qps
	// 3. 返回 DefaultQPS
	QPS map[string]int
	// QPS 中未匹配到，使用 DefaultQPS，默认为 0，不限流
	DefaultQPS int
}
```

## RedisParallelController

基于 Redis 的 `ParallelController` 实现

- `Acquire` 使用 `incr` 命令实现，`Release` 使用 `decr` 命令实现
- `incr/decr` 前需要先判断当前的 token 数量是否已经达到最大值，且这个操作需要保证原子性，使用 lua 脚本实现

这种方案，没有过期设置，如果 token 发生泄漏无法自动恢复

```go
type RedisParallelControllerOptions struct {
	Redis wrap.RedisClientWrapperOptions
	// key 前缀，可当成命名空间使用
	Prefix string
	// MaxToken 计算规则
	// 1. key 在 map 中，直接返回 key 对应的 qps
	// 2. key 按 '|' 分割，第 0 个字符串作为 key，如果在 map 中，返回 qps
	// 3. 返回 DefaultQPS
	MaxToken map[string]int
	// MaxToken 中未匹配到，使用 DefaultMaxToken，默认为 0，不限流
	DefaultMaxToken int
	// 获取 token 失败时重试时间间隔最大值
	Interval time.Duration
}
```

## OTSParallelController

基于 OTS 的 `ParallelController` 实现，和 `OTSRateLimiter` 的实现类似，使用 `UpdateRow` 中的条件增量更新

```go
type OTSParallelControllerOptions struct {
	OTS wrap.OTSTableStoreClientWrapperOptions
	// OTS 表名
	Table string
	// key 前缀，可当成命名空间使用
	Prefix string
	// MaxToken 计算规则
	// 1. key 在 map 中，直接返回 key 对应的 qps
	// 2. key 按 '|' 分割，第 0 个字符串作为 key，如果在 map 中，返回 qps
	// 3. 返回 DefaultQPS
	MaxToken map[string]int
	// MaxToken 中未匹配到，使用 DefaultMaxToken，默认为 0，不限流
	DefaultMaxToken int
	// 获取 token 失败时重试时间间隔最大值
	Interval time.Duration
}
```

## RedisTimedParallelController

基于 Redis 的 `ParallelController` 实现，增加超时自动释放机制

- 基于 redis 的 order set 实现
- 使用 microsecond 时间戳作为 token，保存在 redis 的 order set 中
- `Acquire` 时，将当前时间戳添加到 order set 中，并且释放已经过期的 token，这个操作比较复杂，需要使用 lua 脚本实现
- `Release` 直接使用 zrem 即可

```go
type RedisTimedParallelControllerOptions struct {
	Redis wrap.RedisClientWrapperOptions
	// key 前缀，可当成命名空间使用
	Prefix string
	// MaxToken 计算规则
	// 1. key 在 map 中，直接返回 key 对应的 qps
	// 2. key 按 '|' 分割，第 0 个字符串作为 key，如果在 map 中，返回 qps
	// 3. 返回 DefaultQPS
	MaxToken map[string]int
	// MaxToken 中未匹配到，使用 DefaultMaxToken，默认为 0，不限流
	DefaultMaxToken int
	// 获取 token 失败时重试时间间隔最大值
	Interval time.Duration
	// token 过期时间
	Expiration time.Duration
}
```

## RedisLocker

基于 Redis 的 `Locker` 实现

- `Lock` 使用 `setnx` 命令实现
- `Unlock` 使用 `del` 命令实现，但是需要判断锁的 val 是否为当前的 client，防止 Unlock 不属于自己持有的锁
- 每把锁启动额外的协程自动更新 token，使用 `expire` 命令刷新锁时间

```go
type RedisLockerOptions struct {
	Redis wrap.RedisClientWrapperOptions
	// key 前缀，可当成命名空间使用
	Prefix string
	// 锁过期时间
	Expiration time.Duration
	// 锁刷新时间，锁刷新间隔 = Expiration - RenewTime
	RenewTime time.Duration
	// 所获取失败最大重试间隔
	Interval time.Duration
}
```