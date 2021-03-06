package microx

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/wrap"
)

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

type RedisTimedParallelController struct {
	client  *wrap.RedisClientWrapper
	options *RedisTimedParallelControllerOptions

	acquireScriptSha1 string
}

func NewRedisTimedParallelControllerWithOptions(options *RedisTimedParallelControllerOptions) (*RedisTimedParallelController, error) {
	if options.Interval == 0 {
		options.Interval = time.Microsecond
	}

	c := &RedisTimedParallelController{
		options: options,
	}

	client, err := wrap.NewRedisClientWrapperWithOptions(&options.Redis)
	if err != nil {
		return nil, errors.Wrap(err, "wrap.NewRedisClientWrapperWithOptions failed")
	}

	{
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		// key maxToken expiration
		res := client.ScriptLoad(ctx, `
local time = redis.call('TIME')
local timestampMicroseconds = tonumber(time[1]) * 1000000 + tonumber(time[2])
local maxToken = tonumber(ARGV[1])
local expiration = tonumber(ARGV[2])
while true do
	local vs = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', '+inf', 'limit', '0', '1')
	if not vs[1] then
		break
	end
	if timestampMicroseconds - expiration < tonumber(vs[1]) then
		break
	end
	redis.call('ZREM', KEYS[1], vs[1])
end
if redis.call('ZCOUNT', KEYS[1], '-inf', '+inf') < maxToken then
	redis.call('ZADD', KEYS[1], timestampMicroseconds, timestampMicroseconds)
	return timestampMicroseconds
else
	return 0
end
`)
		if res == nil {
			return nil, errors.New("redis return nil")
		}
		if res.Err() != nil {
			return nil, errors.Wrap(res.Err(), "client.ScriptLoad failed")
		}
		c.acquireScriptSha1 = res.Val()
	}

	c.client = client
	return c, nil
}

func (c *RedisTimedParallelController) TryAcquire(ctx context.Context, key string) (int, error) {
	maxToken := c.calculateMaxToken(key)
	if maxToken == 0 {
		return 0, nil
	}

	key = c.generateKey(key)
	res := c.client.EvalSha(ctx, c.acquireScriptSha1, []string{key}, maxToken, c.options.Expiration.Microseconds())
	if res == nil {
		return 0, errors.New("redis return nil")
	}
	if res.Err() != nil {
		return 0, errors.Wrap(res.Err(), "c.client.EvalSha acquireScriptSha1 failed")
	}
	if res.Val() != int64(0) {
		return int(res.Val().(int64)), nil
	}
	return 0, micro.ErrParallelControl
}

func (c *RedisTimedParallelController) Acquire(ctx context.Context, key string) (int, error) {
	maxToken := c.calculateMaxToken(key)
	if maxToken == 0 {
		return 0, nil
	}

	key = c.generateKey(key)
	for {
		res := c.client.EvalSha(ctx, c.acquireScriptSha1, []string{key}, maxToken, c.options.Expiration.Microseconds())
		if res == nil {
			return 0, errors.New("redis return nil")
		}
		if res.Err() != nil {
			return 0, errors.Wrap(res.Err(), "c.client.EvalSha acquireScriptSha1 failed")
		}
		if res.Val() != int64(0) {
			return int(res.Val().(int64)), nil
		}
		select {
		case <-ctx.Done():
			return 0, micro.ErrContextCancel
		case <-time.After(time.Duration(rand.Int63n(c.options.Interval.Microseconds())) * time.Microsecond):
		}
	}
}

func (c *RedisTimedParallelController) Release(ctx context.Context, key string, token int) error {
	maxToken := c.calculateMaxToken(key)
	if maxToken == 0 {
		return nil
	}

	key = c.generateKey(key)
	res := c.client.ZRem(ctx, key, token)
	if res == nil {
		return errors.New("redis return nil")
	}
	if res.Err() != nil {
		return errors.Wrap(res.Err(), "c.client.Decr failed")
	}
	return nil
}

func (c *RedisTimedParallelController) generateKey(key string) string {
	if c.options.Prefix == "" {
		return key
	}
	return fmt.Sprintf("%s_%s", c.options.Prefix, key)
}

func (c *RedisTimedParallelController) calculateMaxToken(key string) int {
	if maxToken, ok := c.options.MaxToken[key]; ok {
		return maxToken
	}

	if idx := strings.Index(key, "|"); idx != -1 {
		if maxToken, ok := c.options.MaxToken[key[:idx]]; ok {
			return maxToken
		}
	}

	return c.options.DefaultMaxToken
}
