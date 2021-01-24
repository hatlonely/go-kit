package wrap

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type GORMDBWrapperOptions struct {
	Retry   RetryOptions
	Wrapper WrapperOptions
	Gorm    struct {
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
}

func NewGORMDBWrapperWithOptions(options *GORMDBWrapperOptions) (*GORMDBWrapper, error) {
	retry, err := NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}

	cli, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		options.Gorm.Username, options.Gorm.Password, options.Gorm.Host, options.Gorm.Port, options.Gorm.Database,
	))
	if err != nil {
		return nil, err
	}
	if options.Gorm.ConnMaxLifeTime != 0 {
		cli.DB().SetConnMaxLifetime(options.Gorm.ConnMaxLifeTime)
	}
	if options.Gorm.MaxOpenConns != 0 {
		cli.DB().SetMaxOpenConns(options.Gorm.MaxOpenConns)
	}
	if options.Gorm.MaxIdleConns != 0 {
		cli.DB().SetMaxIdleConns(options.Gorm.MaxIdleConns)
	}
	cli.LogMode(options.Gorm.LogMode)
	return &GORMDBWrapper{
		obj:     cli,
		retry:   retry,
		options: &options.Wrapper,
	}, nil
}
