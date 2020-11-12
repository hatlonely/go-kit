## Quick Start

你可以直接使用 Get 方法获取配置项

```go
package main

import (
	"fmt"

	"github.com/hatlonely/go-kit/config"
)

func main() {
	cfg, err := config.NewConfigWithSimpleFile("test.json")
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.GetString("mysql.username"))
	fmt.Println(cfg.GetInt("mysql.port"))
}
```

使用 Unmarshal 直接将配置映射到结构体中

```go
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
```

类似 flag 的 bind 方式将配置项绑定到变量或者结构体上，这个变量会随着配置的变化而自动更新

```go
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
	if err := config.InitWithSimpleFile("test.json"); err != nil {
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
```

也可以通过监听配置项的变化，来实现动态的配置更新

```go
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
	cfg, err := config.NewConfigWithSimpleFile("test.json")
	if err != nil {
		panic(err)
	}

	cfg.AddOnItemChangeHandler("mysql", func(cfg *config.Config) {
		var opt MysqlOptions
		if err := cfg.Unmarshal(&opt, refx.WithCamelName()); err != nil {
			fmt.Println(err)
			return
		}
		options = &opt
	})

	if err := cfg.Watch(); err != nil {
		panic(err)
	}
	defer cfg.Stop()

	fmt.Println(options)
}
```
