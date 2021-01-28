package wrap

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
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
	Retry            RetryOptions
	Wrapper          WrapperOptions
	Gorm             GormOptions
	RateLimiterGroup RateLimiterGroupOptions
}

func NewGORMDBWrapperWithOptions(options *GORMDBWrapperOptions, opts ...refx.Option) (*GORMDBWrapper, error) {
	var w GORMDBWrapper
	var err error

	w.options = &options.Wrapper
	w.retry, err = NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}
	w.rateLimiterGroup, err = NewRateLimiterGroupWithOptions(&options.RateLimiterGroup, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "NewRateLimiterGroupWithOptions failed")
	}
	if w.options.EnableMetric {
		w.CreateMetric(w.options)
	}

	client, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		options.Gorm.Username, options.Gorm.Password, options.Gorm.Host, options.Gorm.Port, options.Gorm.Database,
	))
	if err != nil {
		return nil, err
	}
	if options.Gorm.ConnMaxLifeTime != 0 {
		client.DB().SetConnMaxLifetime(options.Gorm.ConnMaxLifeTime)
	}
	if options.Gorm.MaxOpenConns != 0 {
		client.DB().SetMaxOpenConns(options.Gorm.MaxOpenConns)
	}
	if options.Gorm.MaxIdleConns != 0 {
		client.DB().SetMaxIdleConns(options.Gorm.MaxIdleConns)
	}
	client.LogMode(options.Gorm.LogMode)
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
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiterGroup"), w.OnRateLimiterGroupChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Gorm"), func(cfg *config.Config) error {
		var options GormOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		client, err := gorm.Open("mysql", fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			options.Username, options.Password, options.Host, options.Port, options.Database,
		))
		if err != nil {
			return errors.Wrap(err, "gorm.Open failed")
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
		w.obj = client
		return nil
	})

	return w, err
}
