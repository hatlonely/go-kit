package cli

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

func NewMysql(opts ...MySQLOption) (*gorm.DB, error) {
	options := defaultMySQLOptions
	for _, opt := range opts {
		opt(&options)
	}
	return NewMysqlWithOptions(&options)
}

func NewMysqlWithConfig(cfg *config.Config, opts ...refx.Option) (*gorm.DB, error) {
	options := defaultMySQLOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, err
	}
	return NewMysqlWithOptions(&options)
}

func NewMysqlWithOptions(options *MySQLOptions) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		options.Username, options.Password, options.Host, options.Port, options.Database,
	))
	if err != nil {
		return nil, err
	}
	if options.ConnMaxLifeTime != 0 {
		db.DB().SetConnMaxLifetime(options.ConnMaxLifeTime)
	}
	if options.MaxOpenConns != 0 {
		db.DB().SetMaxOpenConns(options.MaxOpenConns)
	}
	if options.MaxIdleConns != 0 {
		db.DB().SetMaxIdleConns(options.MaxIdleConns)
	}
	db.LogMode(options.LogMode)
	return db, nil
}

type MySQLOptions struct {
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

var defaultMySQLOptions = MySQLOptions{
	Host:            "localhost",
	Port:            3306,
	ConnMaxLifeTime: 60 * time.Second,
	MaxIdleConns:    10,
	MaxOpenConns:    20,
	LogMode:         false,
}

type MySQLOption func(options *MySQLOptions)

func WithMysqlAuth(username, password string) MySQLOption {
	return func(options *MySQLOptions) {
		options.Username = username
		options.Password = password
	}
}

func WithMysqlAddr(host string, port int) MySQLOption {
	return func(options *MySQLOptions) {
		options.Host = host
		options.Port = port
	}
}

func WithMysqlDatabase(database string) MySQLOption {
	return func(options *MySQLOptions) {
		options.Database = database
	}
}

func WithMysqlConnMaxLifeTime(connMaxLifeTime time.Duration) MySQLOption {
	return func(options *MySQLOptions) {
		options.ConnMaxLifeTime = connMaxLifeTime
	}
}

func WithMysqlConns(maxIdle int, maxOpen int) MySQLOption {
	return func(options *MySQLOptions) {
		options.MaxIdleConns = maxIdle
		options.MaxOpenConns = maxOpen
	}
}

func WithLogMode() MySQLOption {
	return func(options *MySQLOptions) {
		options.LogMode = true
	}
}
