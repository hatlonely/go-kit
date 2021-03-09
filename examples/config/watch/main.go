package main

import (
	"fmt"
	"time"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type MysqlOptions struct {
	Username        string
	Password        string
	Database        string
	Host            string
	Port            string
	ConnMaxLifeTime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

var options *MysqlOptions

func main() {
	cfg, err := config.NewConfigWithSimpleFile("../test.json")
	if err != nil {
		panic(err)
	}

	cfg.AddOnItemChangeHandler("mysql", func(cfg *config.Config) error {
		var opt MysqlOptions
		if err := cfg.Unmarshal(&opt, refx.WithCamelName()); err != nil {
			return err
		}
		options = &opt
		return nil
	})

	if err := cfg.Watch(); err != nil {
		panic(err)
	}
	defer cfg.Stop()

	fmt.Println(options)
}
