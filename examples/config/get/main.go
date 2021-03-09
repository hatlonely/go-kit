package main

import (
	"fmt"

	"github.com/hatlonely/go-kit/config"
)

func main() {
	cfg, err := config.NewConfigWithSimpleFile("../test.json")
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.GetString("mysql.username"))
	fmt.Println(cfg.GetInt("mysql.port"))
}
