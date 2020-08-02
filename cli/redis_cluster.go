package cli

import (
	"time"

	"github.com/go-redis/redis"

	"github.com/hatlonely/go-kit/config"
)

func NewRedisCluster(opts ...RedisClusterOption) (*redis.ClusterClient, error) {
	options := defaultRedisClusterOptions
	for _, opt := range opts {
		opt(&options)
	}
	return NewRedisClusterWithOptions(&options)
}

func NewRedisClusterWithConfig(conf *config.Config) (*redis.ClusterClient, error) {
	options := defaultRedisClusterOptions
	if err := conf.Unmarshal(&options); err != nil {
		return nil, err
	}
	return NewRedisClusterWithOptions(&options)
}

func NewRedisClusterWithOptions(options *RedisClusterOptions) (*redis.ClusterClient, error) {
	client := redis.NewClusterClient((*redis.ClusterOptions)(options))
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return client, nil
}

type RedisClusterOptions redis.ClusterOptions

var defaultRedisClusterOptions = RedisClusterOptions{
	Addrs:        []string{"127.0.0.1:6379"},
	DialTimeout:  300 * time.Millisecond,
	ReadTimeout:  300 * time.Millisecond,
	WriteTimeout: 300 * time.Millisecond,
	MaxRetries:   3,
	PoolSize:     20,
}

type RedisClusterOption func(options *RedisClusterOptions)

func WithRedisClusterAddrs(addr ...string) RedisClusterOption {
	return func(options *RedisClusterOptions) {
		options.Addrs = addr
	}
}

func WithRedisClusterTimeout(dialTimeout, readTimeout, writeTimeout time.Duration) RedisClusterOption {
	return func(options *RedisClusterOptions) {
		options.DialTimeout = dialTimeout
		options.ReadTimeout = readTimeout
		options.WriteTimeout = writeTimeout
	}
}

func WithRedisClusterMaxRetries(maxRetries int) RedisClusterOption {
	return func(options *RedisClusterOptions) {
		options.MaxRetries = maxRetries
	}
}

func WithRedisClusterPoolSize(poolSize int) RedisClusterOption {
	return func(options *RedisClusterOptions) {
		options.PoolSize = poolSize
	}
}

func WithRedisClusterPassword(password string) RedisClusterOption {
	return func(options *RedisClusterOptions) {
		options.Password = password
	}
}
