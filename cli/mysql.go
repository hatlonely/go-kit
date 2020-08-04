package cli

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/hatlonely/go-kit/config"
)

func NewMysql(opts ...MySQLOption) (*gorm.DB, error) {
	options := defaultMySQLOptions
	for _, opt := range opts {
		opt(&options)
	}
	return NewMysqlWithOptions(&options)
}

func NewMysqlWithConfig(conf *config.Config) (*gorm.DB, error) {
	options := defaultMySQLOptions
	if err := conf.Unmarshal(&options); err != nil {
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
	if options.ConnMaxLifetime != 0 {
		db.DB().SetConnMaxLifetime(options.ConnMaxLifetime)
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
	Username        string
	Password        string
	Database        string
	Host            string
	Port            int
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
	LogMode         bool
}

var defaultMySQLOptions = MySQLOptions{
	Host:            "localhost",
	Port:            3306,
	ConnMaxLifetime: 60 * time.Second,
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
		options.ConnMaxLifetime = connMaxLifeTime
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
