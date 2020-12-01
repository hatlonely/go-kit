package main

import (
	"fmt"

	"github.com/hatlonely/go-kit/bind"
	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/flag"
	"github.com/hatlonely/go-kit/refx"
)

var Version = "unknown"

type Options struct {
	flag.Options

	Key1 int    `dft:"123"`
	Key2 string `dft:"val"`
	Key3 struct {
		Key4 int `dft:"456"`
		Key5 string
	}
}

func main() {
	var options Options
	if err := flag.Struct(&options, refx.WithCamelName()); err != nil {
		panic(err)
	}
	if err := flag.Parse(flag.WithJsonVal()); err != nil {
		fmt.Println(flag.Usage())
		return
	}
	if options.Help {
		fmt.Println(flag.Usage())
		return
	}
	if options.Version {
		fmt.Println(Version)
		return
	}
	if options.ConfigPath == "" {
		options.ConfigPath = "test.json"
	}
	cfg, err := config.NewConfigWithSimpleFile(options.ConfigPath)
	if err != nil {
		panic(err)
	}

	if err := bind.Bind(&options, []bind.Getter{flag.Instance(), bind.NewEnvGetter(bind.WithEnvPrefix("TEST")), cfg}, refx.WithCamelName()); err != nil {
		panic(err)
	}

	fmt.Println(options)
}
