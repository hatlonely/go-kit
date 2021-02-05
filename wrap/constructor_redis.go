package wrap

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

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

type RedisClientWrapperOptions struct {
	Retry       micro.RetryOptions
	Wrapper     WrapperOptions
	Redis       RedisOptions
	RateLimiter micro.RateLimiterOptions
}

func NewRedisClientWrapperWithOptions(options *RedisClientWrapperOptions, opts ...refx.Option) (*RedisClientWrapper, error) {
	var w RedisClientWrapper
	var err error

	w.options = &options.Wrapper
	w.retry, err = micro.NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}
	w.rateLimiter, err = micro.NewRateLimiterWithOptions(&options.RateLimiter, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "NewRateLimiter failed")
	}
	if w.options.EnableMetric {
		w.CreateMetric(w.options)
	}

	client := redis.NewClient(&redis.Options{
		Addr:         options.Redis.Addr,
		DialTimeout:  options.Redis.DialTimeout,
		ReadTimeout:  options.Redis.ReadTimeout,
		WriteTimeout: options.Redis.WriteTimeout,
		MaxRetries:   options.Redis.MaxRetries,
		PoolSize:     options.Redis.PoolSize,
		DB:           options.Redis.DB,
		Password:     options.Redis.Password,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, errors.Wrap(err, "redis.Client.Ping failed")
	}
	w.obj = client

	return &w, nil
}

func NewRedisClientWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*RedisClientWrapper, error) {
	var options RedisClientWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewRedisClientWrapperWithOptions(&options, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "NewRedisClusterClientWrapperWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Wrapper"), w.OnWrapperChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Retry"), w.OnRetryChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiter"), w.OnRateLimiterChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Redis"), func(cfg *config.Config) error {
		var options RedisOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
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
			return errors.Wrap(err, "redis.Client.Ping failed")
		}
		w.obj = client
		return nil
	})

	return w, err
}

type RedisClusterOptions struct {
	Addrs        []string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MaxRetries   int
	PoolSize     int
	Password     string
}

type RedisClusterClientWrapperOptions struct {
	Retry        micro.RetryOptions
	Wrapper      WrapperOptions
	RedisCluster RedisClusterOptions
	RateLimiter  micro.RateLimiterOptions
}

func NewRedisClusterClientWrapperWithOptions(options *RedisClusterClientWrapperOptions) (*RedisClusterClientWrapper, error) {
	var w RedisClusterClientWrapper
	var err error

	w.options = &options.Wrapper
	w.retry, err = micro.NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}
	w.rateLimiter, err = micro.NewRateLimiterWithOptions(&options.RateLimiter)
	if err != nil {
		return nil, errors.Wrap(err, "NewRateLimiterWithOptions failed")
	}
	if w.options.EnableMetric {
		w.CreateMetric(w.options)
	}

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        options.RedisCluster.Addrs,
		DialTimeout:  options.RedisCluster.DialTimeout,
		ReadTimeout:  options.RedisCluster.ReadTimeout,
		WriteTimeout: options.RedisCluster.WriteTimeout,
		MaxRetries:   options.RedisCluster.MaxRetries,
		PoolSize:     options.RedisCluster.PoolSize,
		Password:     options.RedisCluster.Password,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, errors.Wrap(err, "redis.ClusterClient.Ping failed")
	}
	w.obj = client

	return &w, nil
}

func NewRedisClusterClientWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*RedisClusterClientWrapper, error) {
	var options RedisClusterClientWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewRedisClusterClientWrapperWithOptions(&options)
	if err != nil {
		return nil, errors.Wrap(err, "NewRedisClusterClientWrapperWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Wrapper"), w.OnWrapperChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Retry"), w.OnRetryChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiter"), w.OnRateLimiterChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RedisCluster"), func(cfg *config.Config) error {
		var options RedisClusterOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		client := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        options.Addrs,
			DialTimeout:  options.DialTimeout,
			ReadTimeout:  options.ReadTimeout,
			WriteTimeout: options.WriteTimeout,
			MaxRetries:   options.MaxRetries,
			PoolSize:     options.PoolSize,
			Password:     options.Password,
		})
		if err := client.Ping().Err(); err != nil {
			return errors.Wrap(err, "redis.ClusterClient.Ping failed")
		}
		w.obj = client
		return nil
	})

	return w, err
}
