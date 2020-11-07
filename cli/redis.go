package cli

import (
	"time"

	"github.com/go-redis/redis"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

func NewRedis(opts ...RedisOption) (*redis.Client, error) {
	options := defaultRedisOptions
	for _, opt := range opts {
		opt(&options)
	}
	return NewRedisWithOptions(&options)
}

func NewRedisWithConfig(cfg *config.Config, opts ...refx.Option) (*redis.Client, error) {
	options := defaultRedisOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, err
	}
	return NewRedisWithOptions(&options)
}

func NewRedisWithOptions(options *RedisOptions) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         options.Addr,
		DialTimeout:  options.DialTimeout,
		ReadTimeout:  options.ReadTimeout,
		WriteTimeout: options.WriteTimeout,
		MaxRetries:   options.MaxRetries,
		PoolSize:     options.PoolSize,
		DB:           options.DB,
		Password:     options.Password,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return client, nil
}

type RedisOptions struct {
	Addr         string        `dft:"127.0.0.1:6379"`
	DialTimeout  time.Duration `dft:"300ms"`
	ReadTimeout  time.Duration `dft:"300ms"`
	WriteTimeout time.Duration `dft:"300ms"`
	MaxRetries   int           `dft:"3"`
	PoolSize     int           `dft:"20"`
	DB           int
	Password     string
}

var defaultRedisOptions = RedisOptions{
	Addr:         "127.0.0.1:6379",
	DialTimeout:  300 * time.Millisecond,
	ReadTimeout:  300 * time.Millisecond,
	WriteTimeout: 300 * time.Millisecond,
	MaxRetries:   3,
	PoolSize:     20,
	DB:           0,
}

type RedisOption func(options *RedisOptions)

func WithRedisAddr(addr string) RedisOption {
	return func(options *RedisOptions) {
		options.Addr = addr
	}
}

func WithRedisTimeout(dialTimeout, readTimeout, writeTimeout time.Duration) RedisOption {
	return func(options *RedisOptions) {
		options.DialTimeout = dialTimeout
		options.ReadTimeout = readTimeout
		options.WriteTimeout = writeTimeout
	}
}

func WithRedisMaxRetries(maxRetries int) RedisOption {
	return func(options *RedisOptions) {
		options.MaxRetries = maxRetries
	}
}

func WithRedisPoolSize(poolSize int) RedisOption {
	return func(options *RedisOptions) {
		options.PoolSize = poolSize
	}
}

func WithRedisPassword(password string) RedisOption {
	return func(options *RedisOptions) {
		options.Password = password
	}
}

func WithRedisDB(db int) RedisOption {
	return func(options *RedisOptions) {
		options.DB = db
	}
}
