package config_test

import (
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/logger"
	"github.com/hatlonely/go-kit/refx"
)

func CreateTestFile() {
	fp, _ := os.Create("test.json")
	_, _ = fp.WriteString(`{
  "http": {
    "port": 80
  },
  "grpc": {
    "port": 6080
  },
  "service": {
    "accountExpiration": "5m",
    "captchaExpiration": "30m"
  },
  "redis": {
    "addr": "127.0.0.1:6379",
    "dialTimeout": "200ms",
    "readTimeout": "200ms",
    "writeTimeout": "200ms",
    "maxRetries": 3,
    "poolSize": 20,
    "db": 0,
    "password": ""
  },
  "mysql": {
    "username": "root",
    "password": "",
    "database": "account",
    "host": "127.0.0.1",
    "port": 3306,
    "connMaxLifeTime": "60s",
    "maxIdleConns": 10,
    "maxOpenConns": 20
  },
  "email": {
    "from": "hatlonely@foxmail.com",
    "password": "123456",
    "server": "smtp.qq.com",
    "port": 25
  },
  "logger": {
    "grpc": {
      "level": "Info",
      "writers": [{
        "type": "RotateFile",
        "rotateFileWriter": {
          "filename": "log/account.rpc",
          "maxAge": "24h",
          "formatter": {
            "type": "Json"
          }
        }
      }]
    },
    "info": {
      "level": "Info",
      "writers": [{
        "type": "RotateFile",
        "rotateFileWriter": {
          "filename": "log/account.log",
          "maxAge": "48h",
          "formatter": {
            "type": "Json"
          }
        }
      }]
    }
  }
}`)
	_ = fp.Close()
}

func CreateBaseFile() {
	fp, _ := os.Create("base.json")
	_, _ = fp.WriteString(`{
  "decoder": {
    "type": "Json"
  },
  "provider": {
    "type": "Local",
	"localProvider": {
      "filename": "test.json"
    }
  },
  "cipher": {
    "type": "Group",
    "cipherGroup": [{
      "type": "AES",
      "aesCipher": {
        "base64Key": "IrjXy4vx7iwgCLaUeu5TVUA9TkgMwSw3QWcgE/IW5W0="
      }
    }, {
      "type": "Base64",
    }]
  }
}`)
}

func DeleteTestFile() {
	_ = os.Remove("test.json")
}

func DeleteBaseFile() {
	_ = os.Remove("base.json")
}

// 场景一：直接使用 config 的全局 Get 方法
// 可读性：代码最简单，无需提前声明，写起来最方便，可读性也还可以
// 复用性：配置项是写死的，这样的模块很难复用
// 可测试性：需要调用 config 的 set 方法，mock 使用到的配置项，mock 的代价较高
// 可维护性：配置项散落在代码中，新增和修改都不太方便
// 安全性：可使用 GetD，GetE 之类的方法来做一些错误处理，关联的多个配置项无法保证原子性（这种场景触发几率较低，目前还未碰到过）

func TestExample1(t *testing.T) {
	CreateTestFile()

	// package main
	if err := config.InitWithSimpleFile("test.json"); err != nil {
		panic(err)
	}
	config.SetLogger(logger.NewStdoutLogger())
	if err := config.Watch(); err != nil {
		panic(err)
	}
	defer config.Stop()

	// package module
	fmt.Println(config.GetString("redis.addr"))
	fmt.Println(config.GetString("mysql.host"))
	fmt.Println(config.GetString("logger.info.writers[0].type"))

	// package test
	_ = config.UnsafeSet("redis.addr", "127.0.0.1:6379")
	_ = config.UnsafeSet("mysql.host", "localhost")
	_ = config.UnsafeSet("logger.info.writers[0].type", "Stdout")

	DeleteTestFile()
}

// 场景二：使用全局 config 的 bind 类型方法，类 flag 的使用方式
// 可读性：代码比较简单，类似 flag 的使用方法，需提前将变量绑定一个 key 上，使用时直接使用变量，可读性较好
// 复用性：配置项依然是写死的，复用性较差
// 可测试性：可重新赋值配置项，相比较调用 config 的 set 方法，稍微好一点点
// 可维护性：配置项可集中在模块中声明，维护起来稍微方便一些
// 安全性：可以在检查声明中增加失败回调保证安全，关联的多个配置项同样无法保证原子性
func TestExample2(t *testing.T) {
	CreateTestFile()

	// package main
	if err := config.InitWithSimpleFile("test.json"); err != nil {
		panic(err)
	}
	config.SetLogger(logger.NewStdoutLogger())
	if err := config.Watch(); err != nil {
		panic(err)
	}
	defer config.Stop()

	// package module
	var redisAddr = config.String("redis.addr", config.OnFail(func(err error) {
		fmt.Println(err)
	}))
	var mysqlHost = config.String("mysql.host")

	fmt.Println(redisAddr.Get())
	fmt.Println(mysqlHost.Get())

	// package test
	redisAddr = config.NewAtomicString("127.0.0.1:6379")
	mysqlHost = config.NewAtomicString("localhost")

	DeleteTestFile()
}

// 场景三: 将配置项定义成一个结构体，使用全局 config 的 bind 对象方法
// 可读性：可读性较好
// 复用性：配置项依然是写死的，复用性较差
// 可测试性：需要手动设置 option 中的各个配置项，和场景二差不多
// 可维护性：配置项可集中在一个结构体中，和场景二差不多
// 安全性：可保证关联配置项的原子性
func TestExample3(t *testing.T) {
	CreateTestFile()

	// package main
	if err := config.InitWithSimpleFile("test.json"); err != nil {
		panic(err)
	}
	config.SetLogger(logger.NewStdoutLogger())
	if err := config.Watch(); err != nil {
		panic(err)
	}
	defer config.Stop()

	// package module
	type Options struct {
		Addr         string        `dft:"127.0.0.1:6379"`
		DialTimeout  time.Duration `dft:"300ms"`
		ReadTimeout  time.Duration `dft:"300ms"`
		WriteTimeout time.Duration `dft:"300ms"`
		MaxRetries   int           `dft:"3"`
		PoolSize     int           `dft:"20"`
		DB           int
		Password     string
	}

	var option = config.Bind("redis", Options{}, config.WithUnmarshalOptions(refx.WithCamelName()))
	fmt.Println(option.Load().(Options).Addr)
	fmt.Println(option.Load().(Options).DialTimeout)
	fmt.Println(option.Load().(Options).ReadTimeout)

	// package test
	option = &atomic.Value{}
	option.Store(Options{
		Addr:        "127.0.0.1:6379",
		DialTimeout: 300 * time.Millisecond,
		ReadTimeout: 300 * time.Millisecond,
	})

	DeleteTestFile()
}

// 场景四: 在 OnChangeHandler 中初始化变量，模块中仅声明变量，不绑定 key
// 可读性：模块内部的可读性较好，main 会随着配置项增多而变差
// 复用性：配置项不再和 key 绑定，模块是可复用的
// 可测试性：和场景二一样
// 可维护性：和场景二一样
// 安全性：和场景二一样
func TestExample4(t *testing.T) {
	CreateTestFile()

	// package module
	var RedisAddr config.AtomicString
	var MysqlHost config.AtomicString

	// package main
	cfg, err := config.NewSimpleFileConfig("test.json")
	if err != nil {
		panic(err)
	}
	cfg.SetLogger(logger.NewStdoutLogger())
	cfg.AddOnChangeHandler(func(conf *config.Config) {
		RedisAddr.Set(conf.GetString("redis.addr"))
		MysqlHost.Set(conf.GetString("mysql.host"))
	})
	if err := cfg.Watch(); err != nil {
		panic(err)
	}
	defer cfg.Stop()

	// package module
	fmt.Println(RedisAddr.Get())
	fmt.Println(MysqlHost.Get())

	// package test
	RedisAddr = *config.NewAtomicString("127.0.0.1:6379")
	MysqlHost = *config.NewAtomicString("localhost")

	DeleteTestFile()
}

// 场景五: 在 OnChangeHandler 中初始化变量，模块中仅声明变量，不绑定 key
// 可读性：模块内部的可读性较好，main 中每个模块都只有一次 Unmarshal 操作，不会随着配置项增多而增加复杂性，可读性较好
// 复用性：配置项不和 key 绑定，模块是可复用的
// 可测试性：和场景三一样
// 可维护性：和场景三一样
// 安全性：和场景三一样
func TestExample5(t *testing.T) {
	CreateTestFile()

	// package module
	type Options struct {
		Addr         string        `dft:"127.0.0.1:6379"`
		DialTimeout  time.Duration `dft:"300ms"`
		ReadTimeout  time.Duration `dft:"300ms"`
		WriteTimeout time.Duration `dft:"300ms"`
		MaxRetries   int           `dft:"3"`
		PoolSize     int           `dft:"20"`
		DB           int
		Password     string
	}
	var options atomic.Value

	// package main
	cfg, err := config.NewSimpleFileConfig("test.json")
	if err != nil {
		panic(err)
	}
	cfg.SetLogger(logger.NewStdoutLogger())
	cfg.AddOnItemChangeHandler("redis", func(conf *config.Config) {
		var opt Options
		if err := conf.Unmarshal(&opt, refx.WithCamelName()); err != nil {
			return
		}
		options.Store(opt)
	})
	if err := cfg.Watch(); err != nil {
		panic(err)
	}
	defer cfg.Stop()

	// package module
	fmt.Println(options.Load().(Options).Addr)
	fmt.Println(options.Load().(Options).DialTimeout)
	fmt.Println(options.Load().(Options).ReadTimeout)

	// package test
	options.Store(Options{
		Addr:        "127.0.0.1:6379",
		DialTimeout: 200 * time.Millisecond,
		ReadTimeout: 200 * time.Millisecond,
	})

	DeleteTestFile()
}

// 场景六: 使用绑定变量作为参数传递给构造函数
// 可读性：模块内部的可读性较好，对象在 main 中构造，main 会随着对象增多而变复杂
// 复用性：可以做到对象级别的复用，复用性最好
// 可测试性：直接通过构造函数构造测试对象，测试很方便
// 可维护性：每个对象维护自己的动态参数，维护比较方便
// 安全性：无法保证关联配置的原子性
func TestExample6(t *testing.T) {
	CreateTestFile()

	// package main
	cfg, err := config.NewSimpleFileConfig("test.json")
	if err != nil {
		panic(err)
	}
	cfg.SetLogger(logger.NewStdoutLogger())
	if err := cfg.Watch(); err != nil {
		panic(err)
	}
	defer cfg.Stop()

	myType := NewMyType6(cfg.String("redis.addr"), cfg.Duration("redis.dialTimeout"), cfg.Duration("mysql.readTimeout"))

	// package module
	myType.DoSomething()

	// package test
	testType := NewMyType6(config.NewAtomicString("127.0.0.1:6378"), config.NewAtomicDuration(100*time.Millisecond), config.NewAtomicDuration(100*time.Millisecond))
	testType.DoSomething()

	DeleteTestFile()
}

// package module
type MyType6 struct {
	Addr        *config.AtomicString
	DialTimeout *config.AtomicDuration
	ReadTimeout *config.AtomicDuration
}

func NewMyType6(addr *config.AtomicString, dialTimeout *config.AtomicDuration, readTimeout *config.AtomicDuration) *MyType6 {
	return &MyType6{
		Addr:        addr,
		DialTimeout: dialTimeout,
		ReadTimeout: readTimeout,
	}
}

func (t MyType6) DoSomething() {
	fmt.Println(t.Addr.Get())
	fmt.Println(t.DialTimeout.Get())
	fmt.Println(t.ReadTimeout.Get())
}

// 场景七: 使用绑定的结构体作为参数传递给构造函数
// 可读性：模块内部的可读性较好，对象在 main 中构造，main 会随着对象增多而变复杂
// 复用性：可以做到对象级别的复用，复用性好
// 可测试性：直接通过构造函数构造测试对象，测试很方便
// 可维护性：每个对象维护自己的动态参数，维护比较方便
// 安全性：可以保证关联配置项的原子性
func TestExample7(t *testing.T) {
	// package main
	conf, err := config.NewConfigWithBaseFile("testfile/base.json")
	if err != nil {
		panic(err)
	}
	if err := conf.Watch(); err != nil {
		panic(err)
	}
	defer conf.Stop()

	myType := NewMyType7(conf.Bind("OSS", MyType2Option{}))
	myType.DoSomething()

	// package test
	var option atomic.Value
	option.Store(MyType2Option{
		AccessKeyID:     "test-ak",
		AccessKeySecret: "test-sk",
		Endpoint:        "endpoint",
	})
	testType := NewMyType7(&option)
	testType.DoSomething()
}

// package module
type MyType7 struct {
	option *atomic.Value
}

type MyType2Option struct {
	AccessKeyID     string
	AccessKeySecret string
	Endpoint        string
}

func NewMyType7(option *atomic.Value) *MyType7 {
	return &MyType7{
		option: option,
	}
}

func (t MyType7) DoSomething() {
	fmt.Println(t.option.Load().(MyType2Option).AccessKeyID)
	fmt.Println(t.option.Load().(MyType2Option).AccessKeySecret)
	fmt.Println(t.option.Load().(MyType2Option).Endpoint)
}

// 场景八: 类似缓存的使用方式，用配置的 key 来初始化对象
// 可读性：模块内部的可读性较好，对象在 main 中构造，main 会随着对象增多而变复杂
// 复用性：不依赖固定的配置项，对象是可复用的
// 可测试性：需要 mock 全局的配置对象，mock 代价较高，和场景一类似
// 可维护性：和场景一类似
// 安全性：和场景一类似
func TestExample8(t *testing.T) {
	// package main
	if err := config.Init("testfile/base.json"); err != nil {
		panic(err)
	}
	if err := config.Watch(); err != nil {
		panic(err)
	}
	defer config.Stop()

	myType := NewMyType8("OSS.AccessKeyID", "OSS.AccessKeySecret", "OSS.Endpoint")

	// package module
	myType.DoSomething()

	// package test
	testType := NewMyType8("OSS.AccessKeyID", "OSS.AccessKeySecret", "OSS.Endpoint")
	_ = config.UnsafeSet("OSS.AccessKeyID", "test-ak")
	_ = config.UnsafeSet("OSS.AccessKeyID", "test-sk")
	_ = config.UnsafeSet("OSS.AccessKeyID", "endpoint")
	testType.DoSomething()
}

// package module
type MyType8 struct {
	AccessKeyIDConfigKey     string
	AccessKeySecretConfigKey string
	EndpointConfigKey        string
}

func NewMyType8(AccessKeyIDConfigKey string, AccessKeySecretConfigKey string, EndpointConfigKey string) *MyType8 {
	return &MyType8{
		AccessKeyIDConfigKey:     AccessKeyIDConfigKey,
		AccessKeySecretConfigKey: AccessKeySecretConfigKey,
		EndpointConfigKey:        EndpointConfigKey,
	}
}

func (t MyType8) DoSomething() {
	fmt.Println(config.GetString(t.AccessKeyIDConfigKey))
	fmt.Println(config.GetString(t.AccessKeySecretConfigKey))
	fmt.Println(config.GetString(t.EndpointConfigKey))
}
