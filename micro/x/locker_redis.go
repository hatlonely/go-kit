package microx

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/hatlonely/go-kit/logger"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/wrap"
)

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

type RedisLocker struct {
	client  *wrap.RedisClientWrapper
	options *RedisLockerOptions
	uuidMap sync.Map

	delScriptSha1   string
	renewScriptSha1 string
	log             Logger
}

type redisLockerValue struct {
	value     string
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	updatedAt time.Time
}

func NewRedisLockerWithOptions(options *RedisLockerOptions, opts ...refx.Option) (*RedisLocker, error) {
	l := &RedisLocker{
		options: options,
		log:     logger.NewStdoutTextLogger(),
	}

	client, err := wrap.NewRedisClientWrapperWithOptions(&options.Redis, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "wrap.NewRedisClientWrapperWithOptions failed")
	}
	l.client = client

	{
		res := client.ScriptLoad(context.Background(), `
if redis.call('GET', KEYS[1]) == ARGV[1] then
  redis.call('DEL', KEYS[1])
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
		l.delScriptSha1 = res.Val()
	}
	{
		res := client.ScriptLoad(context.Background(), `
if redis.call('GET', KEYS[1]) == ARGV[1] then
  redis.call('EXPIRE', KEYS[1], ARGV[2])
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
		l.renewScriptSha1 = res.Val()
	}

	return l, nil
}

func (l *RedisLocker) SetLogger(log Logger) {
	l.log = log
}

func (l *RedisLocker) TryLock(ctx context.Context, key string) error {
	_, ok := l.uuidMap.Load(key)
	if ok {
		return micro.ErrLocked
	}
	uid := uuid.NewV4().String()
	res := l.client.SetNX(ctx, l.calculateKey(key), uid, l.options.Expiration)
	if res == nil {
		return errors.New("redis return nil")
	}
	if !res.Val() {
		return micro.ErrLocked
	}
	val := &redisLockerValue{
		value: uid,
	}
	val.ctx, val.cancel = context.WithCancel(context.Background())
	val.wg.Add(1)
	l.uuidMap.Store(key, val)
	go l.renew(key)
	return nil
}

func (l *RedisLocker) Lock(ctx context.Context, key string) error {
	for {
		_, ok := l.uuidMap.Load(key)
		if !ok {
			break
		}
		select {
		case <-ctx.Done():
			return micro.ErrContextCancel
		case <-time.After(time.Duration(rand.Int63n(l.options.Interval.Microseconds())) * time.Millisecond):
		}
	}

	uid := uuid.NewV4().String()
	for {
		res := l.client.SetNX(ctx, l.calculateKey(key), uid, l.options.Expiration)
		if res == nil {
			return errors.New("redis return nil")
		}
		if res.Val() {
			val := &redisLockerValue{
				value: uid,
			}
			val.ctx, val.cancel = context.WithCancel(context.Background())
			val.wg.Add(1)
			l.uuidMap.Store(key, val)
			go l.renew(key)
			break
		}
		select {
		case <-ctx.Done():
			return micro.ErrContextCancel
		case <-time.After(time.Duration(rand.Int63n(l.options.Interval.Microseconds())) * time.Millisecond):
		}
	}
	return nil
}

func (l *RedisLocker) Unlock(ctx context.Context, key string) error {
	v, ok := l.uuidMap.Load(key)
	if !ok {
		return errors.New("release unlocked key")
	}
	val := v.(*redisLockerValue)
	val.cancel()
	val.wg.Wait()
	l.uuidMap.Delete(key)

	res := l.client.EvalSha(ctx, l.delScriptSha1, []string{l.calculateKey(key)}, val.value)
	if res == nil {
		return errors.New("redis return nil")
	}
	if res.Err() != nil {
		return errors.Wrap(res.Err(), "client.EvalSha failed")
	}
	if res.Val() == 0 {
		return errors.New("delete others key is not allowed")
	}
	return nil
}

func (l *RedisLocker) renew(key string) {
	for {
		v, ok := l.uuidMap.Load(key)
		if !ok {
			break
		}
		val := v.(*redisLockerValue)
		sleepTime := l.options.Expiration - time.Now().Sub(val.updatedAt) - l.options.RenewTime

		res := l.client.EvalSha(context.Background(), l.renewScriptSha1, []string{l.calculateKey(key)}, val.value, l.options.Expiration.Seconds())
		if res == nil {
			l.log.Warnf("client.EvalSha redis return nil")
			continue
		}
		if res.Err() != nil {
			l.log.Warnf("client.EvalSha renew script failed. err: [%v]", res.Err())
			continue
		}
		// key 不存在
		if res.Val() == 0 {
			break
		}

		val.updatedAt = time.Now()
		l.uuidMap.Store(key, val)

		select {
		case <-val.ctx.Done():
			val.wg.Done()
			return
		case <-time.After(sleepTime):
		}
	}
}

func (l *RedisLocker) calculateKey(key string) string {
	if l.options.Prefix != "" {
		return fmt.Sprintf("%s_%s", l.options.Prefix, key)
	}
	return key
}
