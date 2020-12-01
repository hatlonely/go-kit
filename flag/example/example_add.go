package main

import (
	"fmt"

	"github.com/hatlonely/go-kit/flag"
)

func main() {
	flag.AddOption("key1", "key1 usage", flag.DefaultValue("123"))
	flag.AddOption("key2", "key2 usage", flag.DefaultValue("val2"))
	flag.AddOption("key3", "key3 usage", flag.Shorthand("k"))
	flag.AddArgument("arg1", "arg1 usage")
	if err := flag.Parse(); err != nil {
		fmt.Println(flag.Usage())
		return
	}
	fmt.Println("key1 =>", flag.GetInt("key1"))
	fmt.Println("key2 =>", flag.GetString("key2"))
	fmt.Println("key3 =>", flag.GetString("key3"))
	fmt.Println("arg1 =>", flag.GetString("arg1"))
}
