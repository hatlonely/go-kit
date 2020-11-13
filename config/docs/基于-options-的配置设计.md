基于 options 的配置设计是指通过在模块内部将需要通过配置传入的数据定义成一个 `Options` 结构，使用 `config` 提供的 `Unmarshal` 机制直接获取配置，`config` 自身的配置就采用了这种设计

## 静态配置

以一个 MySQL client 为例，我们可以定义如下选项:

```go
type MySQLOptions struct {
	Username        string
	Password        string
	Database        string
	Host            string
	Port            int
	ConnMaxLifeTime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
	LogMode         bool
}
```

然后提供一个基于 Options 的构造函数

```go
func NewMySQLWithOptions(options *MySQLOptions) (*MysqlCli, error)
```

在 main 中定义一个全局的 Options

```go
type Options struct {
    MySQL cli.MySQLOptions
    Redis cli.RedisOptions
}
```

调用 Unmarshal 配置到 Options，直接用 Options 构造对象

```go
var options Options
if err := cfg.Unmarshal(&options); err != nil {
    panic(err)
}
mysqlCli, err := cli.NewMySQLWithOptions(&options.MySQL)
if err != nil {
    panic(err)
}
redisCli, err := cli.NewRedisWithOptions(&options.Redis)
if err != nil {
    panic(err)
}
```

## 动态配置

和静态配置不同的是，动态配置需要在每次配置更新的时候调用 Unmarshal，因此需要用到 config 的时间处理机制，在 Key 上绑定事件处理 Unmarshal，这里提供两种实现

1. 提供额外的 `OnChangeHandler(cfg *config.Config)`

```go
func (c *MySQLCli)OnChangeHandler(cfg *config.Config) {
    var options MysqlOptions
    if err := cfg.Unmarshal(&options); err != nil {
        fmt.Println(err)
        return
    }
    c.options = &options
}
```

再在 main 中统一注册

```go
cfg.AddOnItemChangeHandler("mysql", mysqlCli.OnChangeHandler)
cfg.AddOnItemChangeHandler("redis", redisCli.OnChangeHandler)
```

2. 提供基于 config 的构造函数

```go
func NewMySQLWithConfig(cfg *config.Config) (*MySQLCli, error) {
    var options MySQLOptions
    if err := cfg.Unmarshal(&options); err != nil {
        return nil, err
    }
    cli, err := NewMySQLWithOptions(&options)
    if err != nil {
        return nil, err
    }
    cfg.AddOnChangeHandler(func (cfg *config.Config) {
        var options MySQLOptions
        if err := cfg.Unmarshal(&options); err != nil {
            return nil, err
        }
        cli.options = &options
    })
    return cli, nil
}
```

再在 main 中调用注册

```go
mysqlCli, err := NewMySQLWithConfig(cfg.Sub("mysql"))
if err != nil {
    panic(err)
}
redisCli, err := NewMySQLWithConfig(cfg.Sub("redis"))
if err != nil {
    panic(err)
}
```
