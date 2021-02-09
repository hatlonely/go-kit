package microx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/wrap"
)

type RedisIncrParallelControllerOptions struct {
	Redis           wrap.RedisClientWrapperOptions
	Prefix          string
	MaxToken        map[string]int
	DefaultMaxToken int
	Interval        time.Duration
}

type RedisIncrParallelController struct {
	client            *wrap.RedisClientWrapper
	options           *RedisIncrParallelControllerOptions
	getTokenScriptSha string
	putTokenScriptSha string
}

func NewRedisIncrParallelControllerWithOptions(options *RedisIncrParallelControllerOptions) (*RedisIncrParallelController, error) {
	client, err := wrap.NewRedisClientWrapperWithOptions(&options.Redis)
	if err != nil {
		return nil, errors.WithMessage(err, "wrap.NewRedisClientWrapperWithOptions")
	}

	c := &RedisIncrParallelController{client: client, options: options}
	{
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		res := client.ScriptLoad(ctx, `
local r = redis.call('GET', KEYS[1])
if r < ARGV[1] then
  redis.call('INCR', KEYS[1])
  return 1
end
return 0
`)
		if res == nil {
			return nil, errors.New("redis return nil")
		}
		if res.Err() != nil {
			return nil, errors.Wrap(err, "client.ScriptLoad failed")
		}
		c.getTokenScriptSha = res.Val()
	}
	{

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		res := client.ScriptLoad(ctx, `
local r = redis.call('GET', KEYS[1])
if r > 0 then
  redis.call('DECR', KEYS[1])
  return 1
end
return 0
`)
		if res == nil {
			return nil, errors.New("redis return nil")
		}
		if res.Err() != nil {
			return nil, errors.Wrap(err, "client.ScriptLoad failed")
		}
		c.putTokenScriptSha = res.Val()
	}

	return c, nil
}

func (c *RedisIncrParallelController) TryAcquire(ctx context.Context, key string) (int, error) {
	maxToken := c.calculateMaxToken(key)
	if maxToken == 0 {
		return 0, nil
	}

	key = c.generateKey(key)
	res := c.client.EvalSha(ctx, c.getTokenScriptSha, []string{key}, maxToken)
	if res == nil {
		return 0, errors.New("redis return nil")
	}
	if res.Err() != nil {
		return 0, errors.Wrap(res.Err(), "c.client.Incr failed")
	}
	if res.Val() == 1 {
		return 0, nil
	}
	return 0, micro.ErrParallelControl
}

func (c *RedisIncrParallelController) Acquire(ctx context.Context, key string) (int, error) {
	maxToken := c.calculateMaxToken(key)
	if maxToken == 0 {
		return 0, nil
	}

	key = c.generateKey(key)
	for {
		res := c.client.EvalSha(ctx, c.getTokenScriptSha, []string{key}, maxToken)
		if res == nil {
			return 0, errors.New("redis return nil")
		}
		if res.Err() != nil {
			return 0, errors.Wrap(res.Err(), "c.client.Incr failed")
		}
		if res.Val() == 1 {
			return 0, nil
		}
		select {
		case <-ctx.Done():
			return 0, errors.New("context cancel")
		case <-time.After(c.options.Interval):
		}
	}
}

func (c *RedisIncrParallelController) Release(ctx context.Context, key string, token int) error {
	maxToken := c.calculateMaxToken(key)
	if maxToken == 0 {
		return nil
	}

	key = c.generateKey(key)
	res := c.client.EvalSha(ctx, c.putTokenScriptSha, []string{key})
	if res == nil {
		return errors.New("redis return nil")
	}
	if res.Err() != nil {
		return errors.Wrap(res.Err(), "c.client.Decr failed")
	}
	return nil
}

func (c *RedisIncrParallelController) generateKey(key string) string {
	if c.options.Prefix == "" {
		return key
	}
	return fmt.Sprintf("%s_%s", c.options.Prefix, key)
}

func (c *RedisIncrParallelController) calculateMaxToken(key string) int {
	if qps, ok := c.options.MaxToken[key]; ok {
		return qps
	}

	if idx := strings.Index(key, "|"); idx != -1 {
		if qps, ok := c.options.MaxToken[key[:idx]]; ok {
			return qps
		}
	}

	return c.options.DefaultMaxToken
}
