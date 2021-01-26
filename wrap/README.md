# wrap

wrap 对一些常见库进行 wrap 提供，保留原生库的所有功能，并提供通用的 trace/metric/retry 功能

## Feature

1. 支持 trace/metric/retry
2. 支持自动从 ecs role 中获取 ak/sk/token 更新 client
3. 支持配置热加载
4. wrap 代码自动生成

## 代码结构

代码分为三个部分

1. 所有库公用的结构，封装了一些通用的逻辑和结构，包括 `retry/error/logger/wrapper`
2. 不同库逻辑相似性不高，不可自动生成的代码，目前主要是各个库的构造方式 `constructor_<T>.go`
3. 可自动生成的代码，包括各个函数接口，以及其他一些相似性较高的逻辑，比如 `OnRetryChange`

## Retry

封装了`github.com/avast/retry-go` 相关的逻辑，支持 retry 的配置化，提供常见的通用 retryIf 方法，支持用户可以拓展自己的 retryIf

```go
type RetryOptions struct {
	Attempts uint `dft:"3"`
	// 重试间隔，当 DelayType == Fixed || DelayType == BackOff 时有效
	Delay time.Duration `dft:"1s"`
	// 最大重试间隔时间
	MaxDelay time.Duration `dft:"3m"`
	// 随机重试间隔，当 DelayType == Random 时有效
	MaxJitter     time.Duration `dft:"1s"`
	LastErrorOnly bool
	// Fixed: 固定重试间隔为 Delay
	// BackOff: 以 Delay 为基础，重试间隔指数增长，Delay, Delay >> 1, Delay >> 2
	// Random: 重试间隔为 (0, MaxJitter) 中的一个随机值
	// Fixed,Random: 重试间隔 (Delay, Delay + MaxJitter) 中的一个随机值
	DelayType string `dft:"BackOff"`
	// 使用定义在 retryRetryIfMap 中的 retryIf 方法
	// 默认重试非 retry.Unrecoverable 的所有错误
	RetryIf string
}
```

### retryIf 的拓展

在 `retryRetryIfMap` 中定义了默认的 `retryIf` 处理，用户可以通过 `RegisterRetryRetryIf` 方法来定义自己的 `retryIf` 方法

```go
func init() {
    wrap.RegisterRetryRetryIf("MyRetryIf", func(err error) bool {
    	switch e := err.(type) {
        case *MyError:
            if e.HttpStatusCode >= http.StatusInternalServerError {
                return true
            }
            if strings.Contains(e.Error(), "timeout") {
                return true
            }
        }
        return false
    })
}
```

### Retry 函数代码模板

```go
err = w.retry.Do(func() error {
    {{.function.resultList}} = w.obj.{{.function.name}}({{.function.paramList}})
    return {{.function.lastResult}}
})
```

## Trace

trace 目前支持 opentracing，用 GlobalTracer 从 context 中获取 trace

### Trace 函数代码模板

```go
if w.options.EnableTrace {
    span, _ := opentracing.StartSpanFromContext(ctx, "{{.package}}.{{.class}}.{{.function.name}}")
    defer span.Finish()
}
```

## Metric

有两个 Metric

- durationMetric 使用 HistogramVec 用于统计平均响应时间和分位数
- totalMetric 使用 CounterVec 用于统计 QPS 和错误分布

### Metric 初始化模板

```go
func (w *{{.wrapClass}}) CreateMetric(options *{{.wrapPackagePrefix}}WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "{{.package}}_{{.class}}_durationMs",
		Help:        "{{.package}} {{.class}} response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "{{.package}}_{{.class}}_total",
		Help:        "{{.package}} {{.class}} request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
}
```

### Metric 函数代码模板

```go
if w.options.EnableMetric {
    ts := time.Now()
    defer func() {
        w.totalMetric.WithLabelValues("{{.package}}.{{.class}}.{{.function.name}}", {{.function.errCode}}).Inc()
        w.durationMetric.WithLabelValues("{{.package}}.{{.class}}.{{.function.name}}", {{.function.errCode}}).Observe(float64(time.Now().Sub(ts).Milliseconds()))
    }()
}
```

### Label

- `method`: 方法名
- `errCode`: 错误码
- `custom`: 用户自定义 Label，常见的使用场景是一个 client 有多种不同业务在使用，业务字段可以设置在这个字段上

## Context

在 Context 中可以传递一些动态的选项，用来改变函数执行的逻辑

```go
type CtxOptions struct {
	DisableTrace           bool
	DisableMetric          bool
	MetricCustomLabelValue string
}
```

用户可以通过 `NewContext` 来创建新的 context，调用方式类似

```go
res, err := w.GetRow(wrap.NewContext(ctx, wrap.WithCtxDisableTrace(), wrap.WithMetricCustomLabelValue("myCustomVal")), &tablestore.GetRowRequest{...})
```

## 热加载

Wrapper 创建的代码逻辑是手动编写的，提供两种创建方式，`NewWithOptions` 以及 `NewWithConfig`，`NewWithConfig` 使用 config 模块的动态加载机制实现了配置的热加载，
热加载包括三个部分，Wrapper、Retry 以及被 Wrap 的类，其中 Wrapper 和 Retry 的热加载逻辑可以自动生成，但需要在构造方法中调用，
此外，Metric 需要提前定义，是无法热加载的


```go
func (w *{{.wrapClass}}) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options {{.wrapPackagePrefix}}WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *{{.wrapClass}}) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options {{.wrapPackagePrefix}}RetryOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		retry, err := {{.wrapPackagePrefix}}NewRetryWithOptions(&options)
		if err != nil {
			return errors.Wrap(err, "NewRetryWithOptions failed")
		}
		w.retry = retry
		return nil
	}
}
```

## Unwrap

有几种场景下，Wrap 无法完全覆盖原生类的所有功能，提供 Unwrap 方法来获取原生类对象

1. 原生类方法使用了私有的类作为参数或者返回，比如 `gorm.SetLogger(log logger)`，这里的 logger 为一个私有的接口，无法 wrap
2. 原生类成员无法 wrap，但是又需要直接访问，比如 `gorm.DB.Error`

## ECS RAM Role 支持

对于阿里云相关的客户端，如果配置中没有提供 ak/sk 信息，wrap 模块会自动尝试从 esc 中获取 ram role 并创建 client，并启动一个后台协程，
在过期之前更新 client

### client 更新策略

ECS 的 Role 会在过期前 30 分钟更新 ak/sk/token，后台协程会在过期前 25 分钟尝试获取新的 ak/sk/token，如果获取失败，30s 之后重试直到成功

### 与配置热加载的冲突

由于配置文件更新，client 被重新创建，此时的老的后台协程可能依旧存在，导致配置更新失效，这里对后台协程引入一个简单的机制解决这个问题，

1. 每次循环开始时保存当前的 client
2. 从等待更新时间中唤醒后再次检查当前保存的 client 和 wrapper 中的 client 是否一致，如果不一致说明已有别的协程更新过 client 直接退出当前协程
3. 获取新的 ecs role，并且创建新的 client
4. 再次检查 client 和 wrapper 中的 client 是否一致，不一致直接退出
5. 更新 client，进入下一次循环

这个方案并不严谨，极端情况下，配置更新会失效，后台协程会依然存在
