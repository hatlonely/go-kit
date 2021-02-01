## Feature

1. 支持 `Get<T>` 风格 api
2. 支持 `Unmarshal` 到结构体，包括数组，map 的嵌套，支持默认值，validator，支持不同风格的 key
3. 支持将配置项动态绑定到原子变量
4. 支持监听配置变更事件
5. 支持子配置，子配置依然是一个完全的配置对象
6. 支持多种配置文件格式，包括 `Yaml`/`Json5`/`Toml`/`Ini`/`Properties`
7. 支持多种配置文件后端，包括本地文件/阿里云表格存储OTS/阿里云配置服务ACM
8. 支持配置项的自动加密解密
9. 支持多种加密方式 `AES`/阿里云密钥服务KMS，以及多种加密方式的组合
10. 提供配置管理工具 `cfg`
11. 支持用户注册自己的配置后端，加解密服务，以及配置格式解析器

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
	Host            string        `dft:"localhost"`
	Port            string        `dft:"3306"`
	ConnMaxLifeTime time.Duration `dft:"5m"`
	MaxIdleConns    int           `dft:"10" rule:"x <= 100 && x >= 0" validate:"gte=0,lte=100"`
	MaxOpenConns    int           `dft:"30"`
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
	if err := cfg.Unmarshal(&options, refx.WithCamelName(), refx.WithPlaygroundValidator()); err != nil {
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
	Host            string        `dft:"localhost"`
	Port            string        `dft:"3306"`
	ConnMaxLifeTime time.Duration `dft:"5m"`
	MaxIdleConns    int           `dft:"10" rule:"x <= 100 && x >= 0" validate:"gte=0,lte=100"`
	MaxOpenConns    int           `dft:"30"`
}

var options *MysqlOptions

func main() {
	cfg, err := config.NewConfigWithSimpleFile("test.json")
	if err != nil {
		panic(err)
	}

	cfg.AddOnItemChangeHandler("mysql", func(cfg *config.Config) {
		var opt MysqlOptions
		if err := cfg.Unmarshal(&opt, refx.WithCamelName(), refx.WithPlaygroundValidator()); err != nil {
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

## 配置文件

你可以通过[配置文件](docs/配置文件.md) (通常命名为 `base.json`) 来构造复杂的配置对象

```json
{
  "decoder": {
    "type": "Json"
  },
  "provider": {
    "type": "Local",
    "options": {
      "filename": "test.json"
    }
  },
  "cipher": {
    "type": "Group",
    "options": [{
      "type": "AES",
      "options": {
        "base64Key": "IrjXy4vx7iwgCLaUeu5TVUA9TkgMwSw3QWcgE/IW5W0="
      }
    }, {
      "type": "Base64"
    }]
  }
}
```

## 管理工具

你可以通过如下命令获取配置工具

```sh
go get -u github.com/hatlonely/go-kit/cmd/cfg
```

通过[配置工具](docs/配置工具.md)可以方便地管理你配置，包括配置的获取，更新，迁移，回滚，备份


```sh
cfg --camelName --inBaseFile base.json -a get
cfg --camelName --inBaseFile base.json -a get --key mysql
cfg --camelName --inBaseFile base.json -a get --key mysql.password
cfg --camelName --inBaseFile base.json -a diff --key mysql.username --val hatlonely
cfg --camelName --inBaseFile base.json -a diff --key mysql.password --val 12345678
cfg --camelName --inBaseFile base.json -a diff --key mysql --jsonVal '{
    "connMaxLifeTime": "60s",
    "database": "testdb2",
    "host": "127.0.0.1",
    "maxIdleConns": 10,
    "maxOpenConns": 20,
    "password": "",
    "port": 3306,
    "username": "hatlonely"
}'
cfg --camelName --inBaseFile base.json -a set --key mysql.username --val hatlonely
cfg --camelName --inBaseFile base.json -a set --key mysql.password --val 12345678
cfg --camelName --inBaseFile base.json -a set --key mysql --jsonVal '{
    "connMaxLifeTime": "60s",
    "database": "testdb2",
    "host": "127.0.0.1",
    "maxIdleConns": 10,
    "maxOpenConns": 20,
    "password": "",
    "port": 3306,
    "username": "hatlonely"
}'
cfg --camelName --inBaseFile base.json -a rollback --backupFile cfg.backup.json.20201113.224234
cfg --camelName --inBaseFile base.json -a put --setCipherKeys mysql.password,redis.password
cfg --camelName --inBaseFile base.json -a put --addCipherKeys email.password
cfg --camelName --inBaseFile base.json -a put --noCipher
cfg --camelName --inBaseFile base.json -a diff --outBaseFile base_ots.json
cfg --camelName --inBaseFile base.json -a put --outBaseFile base_ots.json
```
