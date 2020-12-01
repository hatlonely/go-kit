bind 的设计初衷是统一 flag/config/env 的配置读取, 通过反射机制提供对任意实现 `Get(key string) (value interface{}, exists bool)` 接口对象，Unmarshal 到特定对象上，

flag/config 均实现了 `Get` 接口，可以直接当做 `Getter` 使用，此外 bind 提供 EnvGetter 从环境变量中获取值

## Quick Start

```go
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
```