package main

import (
	"fmt"
	"time"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type ServiceOptions struct {
	AccountExpiration time.Duration
	CaptchaExpiration time.Duration
}

var mysqlUsername = config.String("mysql.username")
var mysqlPort = config.Int("mysql.port", config.OnSucc(func(c *config.Config) {
	fmt.Println(c.Get(""))
}))
var serviceOptions = config.Bind("service", ServiceOptions{}, config.OnFail(func(err error) {
	fmt.Println(err)
}), config.WithUnmarshalOptions(refx.WithCamelName()))

func main() {
	if err := config.InitWithSimpleFile("../test.json"); err != nil {
		panic(err)
	}
	if err := config.Watch(); err != nil {
		panic(err)
	}
	defer config.Stop()

	fmt.Println(mysqlUsername.Get())
	fmt.Println(mysqlPort.Get())
	fmt.Println(serviceOptions.Load().(ServiceOptions))
}
