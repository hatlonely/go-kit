## Feature

1. 基本兼容标准 `flag` 库的使用方法
2. 支持选项参数和位置参数
3. 支持 `json` 格式参数值
4. 支持从结构体构造 `flag` 对象
5. 支持自动生成帮助信息
6. 支持反序列化到结构体

## Quick Start

兼容标准 `flag` 的用法

```go
package main

import (
	"fmt"

	"github.com/hatlonely/go-kit/flag"
)

func main() {
	key1 := flag.Int("key1", 123, "key1 usage")
	key2 := flag.String("key2", "val2", "key2 usage")
	if err := flag.Parse(); err != nil {
		fmt.Println(flag.Usage())
		return
	}
	fmt.Println("key1 =>", *key1)
	fmt.Println("key2 =>", *key2)
}
```

add 风格 api

```go
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

```

options 风格用法

```go
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
```
