package main

import (
	"fmt"

	"github.com/hatlonely/go-kit/flag"
	"github.com/hatlonely/go-kit/refx"
)

type Options struct {
	Key1 int    `flag:"--key1; usage: key1 usage; default: 123"`
	Key2 string `flag:"--key2; usage: key2 usage; default: val2"`
	Key3 string `flag:"-k, --key3; usage: key3 usage"`
	Arg1 string `flag:"arg; usage: arg1 usage"`
}

func main() {
	if err := flag.Struct(&Options{}, refx.WithCamelName()); err != nil {
		panic(err)
	}
	if err := flag.Parse(); err != nil {
		fmt.Println(flag.Usage())
		return
	}
	fmt.Println("key1 =>", flag.GetInt("key1"))
	fmt.Println("key2 =>", flag.GetString("key2"))
	fmt.Println("key3 =>", flag.GetString("key3"))
	fmt.Println("arg1 =>", flag.GetString("arg1"))
}
