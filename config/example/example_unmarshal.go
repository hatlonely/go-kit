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

type Options struct {
	Mysql MysqlOptions
}

func main() {
	cfg, err := config.NewConfigWithSimpleFile("test.json")
	if err != nil {
		panic(err)
	}

	var options Options
	if err := cfg.Unmarshal(&options, refx.WithCamelName()); err != nil {
		panic(err)
	}

	fmt.Println(options)
}
