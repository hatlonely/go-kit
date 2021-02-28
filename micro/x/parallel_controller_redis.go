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

type RedisParallelController struct {
	client            *wrap.RedisClientWrapper
	options           *RedisParallelControllerOptions
	acquireScriptSha1 string
	releaseScriptSha1 string
}

func NewRedisParallelControllerWithOptions(options *RedisParallelControllerOptions) (*RedisParallelController, error) {
	if options.Interval == 0 {
		options.Interval = time.Microsecond
	}

	client, err := wrap.NewRedisClientWrapperWithOptions(&options.Redis)
	if err != nil {
		return nil, errors.WithMessage(err, "wrap.NewRedisClientWrapperWithOptions")
	}

	c := &RedisParallelController{client: client, options: options}
	{
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		res := client.ScriptLoad(ctx, `
local r = redis.call('GET', KEYS[1])
if not r or tonumber(r) < tonumber(ARGV[1]) then
  redis.call('INCR', KEYS[1])
  return 1
end
return 0
`)
		if res == nil {
			return nil, errors.New("redis return nil")
		}
		if res.Err() != nil {
			return nil, errors.Wrap(res.Err(), "client.ScriptLoad failed")
		}
		c.acquireScriptSha1 = res.Val()
	}
	{
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		res := client.ScriptLoad(ctx, `
local r = redis.call('GET', KEYS[1])
if not r then
  return 0
end
if tonumber(r) > 0 then
  redis.call('DECR', KEYS[1])
  return 1
end
return 0
`)
		if res == nil {
			return nil, errors.New("redis return nil")
		}
		if res.Err() != nil {
			return nil, errors.Wrap(res.Err(), "client.ScriptLoad failed")
		}
		c.releaseScriptSha1 = res.Val()
	}

	return c, nil
}

func (c *RedisParallelController) TryAcquire(ctx context.Context, key string) (int, error) {
	maxToken := c.calculateMaxToken(key)
	if maxToken == 0 {
		return 0, nil
	}

	key = c.generateKey(key)
	res := c.client.EvalSha(ctx, c.acquireScriptSha1, []string{key}, maxToken)
	if res == nil {
		return 0, errors.New("redis return nil")
	}
	if res.Err() != nil {
		return 0, errors.Wrap(res.Err(), "c.client.Incr failed")
	}
	if res.Val() == int64(1) {
		return 0, nil
	}
	return 0, micro.ErrParallelControl
}

func (c *RedisParallelController) Acquire(ctx context.Context, key string) (int, error) {
	maxToken := c.calculateMaxToken(key)
	if maxToken == 0 {
		return 0, nil
	}

	key = c.generateKey(key)
	for {
		res := c.client.EvalSha(ctx, c.acquireScriptSha1, []string{key}, maxToken)
		if res == nil {
			return 0, errors.New("redis return nil")
		}
		if res.Err() != nil {
			return 0, errors.Wrap(res.Err(), "c.client.Incr failed")
		}
		if res.Val() == int64(1) {
			return 0, nil
		}
		select {
		case <-ctx.Done():
			return 0, micro.ErrContextCancel
		case <-time.After(time.Duration(rand.Int63n(c.options.Interval.Microseconds())) * time.Microsecond):
		}
	}
}

func (c *RedisParallelController) Release(ctx context.Context, key string, token int) error {
	maxToken := c.calculateMaxToken(key)
	if maxToken == 0 {
		return nil
	}

	key = c.generateKey(key)
	res := c.client.EvalSha(ctx, c.releaseScriptSha1, []string{key})
	if res == nil {
		return errors.New("redis return nil")
	}
	if res.Err() != nil {
		return errors.Wrap(res.Err(), "c.client.Decr failed")
	}
	return nil
}

func (c *RedisParallelController) generateKey(key string) string {
	if c.options.Prefix == "" {
		return key
	}
	return fmt.Sprintf("%s_%s", c.options.Prefix, key)
}

func (c *RedisParallelController) calculateMaxToken(key string) int {
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
