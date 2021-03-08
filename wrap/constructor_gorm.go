package wrap

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

type GormOptions struct {
	Username        string `dft:"root"`
	Password        string
	Database        string
	Host            string        `dft:"localhost"`
	Port            int           `dft:"3306"`
	ConnMaxLifeTime time.Duration `dft:"60s"`
	MaxIdleConns    int           `dft:"10"`
	MaxOpenConns    int           `dft:"20"`
	LogMode         bool
}

type GORMDBWrapperOptions struct {
	Retry              micro.RetryOptions
	Wrapper            WrapperOptions
	Gorm               GormOptions
	RateLimiter        micro.RateLimiterOptions
	ParallelController micro.ParallelControllerOptions
}

func NewGORMDBWithOptions(options *GormOptions) (*gorm.DB, error) {
	client, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		options.Username, options.Password, options.Host, options.Port, options.Database,
	))
	if err != nil {
		return nil, errors.Wrap(err, "gorm.Open failed")
	}
	if options.ConnMaxLifeTime != 0 {
		client.DB().SetConnMaxLifetime(options.ConnMaxLifeTime)
	}
	if options.MaxOpenConns != 0 {
		client.DB().SetMaxOpenConns(options.MaxOpenConns)
	}
	if options.MaxIdleConns != 0 {
		client.DB().SetMaxIdleConns(options.MaxIdleConns)
	}
	client.LogMode(options.LogMode)
	return client, nil
}

func NewGORMDBWrapperWithOptions(options *GORMDBWrapperOptions, opts ...refx.Option) (*GORMDBWrapper, error) {
	var w GORMDBWrapper
	var err error

	w.options = &options.Wrapper
	w.retry, err = micro.NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "micro.NewRetryWithOptions failed")
	}
	w.rateLimiter, err = micro.NewRateLimiterWithOptions(&options.RateLimiter, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "micro.NewRateLimiterWithOptions failed")
	}
	w.parallelController, err = micro.NewParallelControllerWithOptions(&options.ParallelController, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "micro.NewParallelControllerWithOptions failed")
	}
	if w.options.EnableMetric {
		w.CreateMetric(w.options)
	}
	client, err := NewGORMDBWithOptions(&options.Gorm)
	if err != nil {
		return nil, errors.WithMessage(err, "NewGORMDBWithOptions failed")
	}
	w.obj = client

	return &w, nil
}

func NewGORMDBWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*GORMDBWrapper, error) {
	var options GORMDBWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewGORMDBWrapperWithOptions(&options, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "NewGORMDBWrapperWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Wrapper"), w.OnWrapperChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Retry"), w.OnRetryChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiter"), w.OnRateLimiterChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("ParallelController"), w.OnParallelControllerChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Gorm"), func(cfg *config.Config) error {
		var options GormOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		client, err := NewGORMDBWithOptions(&options)
		if err != nil {
			return errors.WithMessage(err, "NewGORMDBWithOptions failed")
		}
		w.obj = client
		return nil
	})

	return w, err
}
