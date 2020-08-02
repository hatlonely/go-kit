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
	db.DB().SetConnMaxLifetime(options.ConnMaxLifetime)
	db.DB().SetMaxOpenConns(options.MaxOpenConns)
	db.DB().SetMaxIdleConns(options.MaxIdleConns)
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
}

var defaultMySQLOptions = MySQLOptions{
	Host:            "localhost",
	Port:            3306,
	ConnMaxLifetime: 60 * time.Second,
	MaxIdleConns:    10,
	MaxOpenConns:    20,
}

type MySQLOption func(options *MySQLOptions)
