# wrap

对一些常见库进行 wrap，并保留原生库的所有功能，wrap 代码通过工具自动生成

## Feature

1. 支持 retry
2. 支持 trace
3. 支持 metric
4. 支持 rateLimiter
5. 支持 parallelController
6. 支持配置热加载
7. 阿里 client 支持从 ecs role 中获取授权信息
8. wrap 代码由自动生成

## Quick Start

```shell
client, err := wrap.NewOTSTableStoreClientWrapperWithOptions(&OTSTableStoreClientWrapperOptions{...})
refx.Muxt(err)

res, err := client.GetRow(ctx, &tablestore.GetRowRequest{...})
if err != nil {
  ...
}
```

## Context

在 Context 中可以传递一些动态的选项，用来改变函数执行的逻辑

```go
type CtxOptions struct {
	DisableTrace           bool
	DisableMetric          bool
	MetricCustomLabelValue string
	TraceTags              map[string]string
	startSpanOpts          []opentracing.StartSpanOption
}
```

用户可以通过 `NewContext` 来创建新的 context

```go
res, err := w.GetRow(wrap.NewContext(ctx, wrap.WithCtxDisableTrace(), wrap.WithMetricCustomLabelValue("myCustomVal")), &tablestore.GetRowRequest{...})
```

## wrap 函数

```go
func (w *OTSTableStoreClientWrapper) GetRow(ctx context.Context, request *tablestore.GetRowRequest) (*tablestore.GetRowResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *tablestore.GetRowResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TableStoreClient.GetRow", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TableStoreClient.GetRow", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TableStoreClient.GetRow", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetRow", ctxOptions.startSpanOpts...)
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.GetRow", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("tablestore.TableStoreClient.GetRow", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("tablestore.TableStoreClient.GetRow", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.GetRow(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}
```

## 微服务组件

参考 [micro](../micro)

## 热加载

Wrapper 创建的代码逻辑是手动编写的，提供两种创建方式，`NewWithOptions` 以及 `NewWithConfig`，`NewWithConfig` 使用 config 模块的动态加载机制实现了配置的热加载，
热加载包括三个部分，Wrapper、Retry 以及被 Wrap 的类，其中 Wrapper 和 Retry 的热加载逻辑可以自动生成，但需要在构造方法中调用， 此外，Metric 需要提前定义，是无法热加载的


```go
func NewOTSTableStoreClientWrapperWithOptions(options *OTSTableStoreClientWrapperOptions, opts ...refx.Option) (*OTSTableStoreClientWrapper, error) {
	var w OTSTableStoreClientWrapper
	var err error

	w.options = &options.Wrapper
	w.retry, err = micro.NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "micro.NewRetryWithOptions failed")
	}
	w.rateLimiter, err = micro.NewRateLimiterWithOptions(&options.RateLimiter, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "micro.NewRateLimiterWithOptions failed")
	}
	w.parallelController, err = micro.NewParallelControllerWithOptions(&options.ParallelController, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "micro.NewParallelControllerWithOptions failed")
	}
	if w.options.EnableMetric {
		w.CreateMetric(w.options)
	}

	if options.OTS.AccessKeyID != "" {
		client := tablestore.NewClient(options.OTS.Endpoint, options.OTS.InstanceName, options.OTS.AccessKeyID, options.OTS.AccessKeySecret)
		if _, err := client.ListTable(); err != nil {
			return nil, errors.Wrap(err, "tablestore.TableStoreClient.ListTable failed")
		}
		w.obj = client
	} else {
		res, err := alics.ECSMetaDataRamSecurityCredentials()
		if err != nil {
			return nil, errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentials failed")
		}
		client := tablestore.NewClientWithConfig(options.OTS.Endpoint, options.OTS.InstanceName, res.AccessKeyID, res.AccessKeySecret, res.SecurityToken, nil)
		if _, err := client.ListTable(); err != nil {
			return nil, errors.Wrap(err, "tablestore.TableStoreClient.ListTable failed")
		}
		w.obj = client
		go w.UpdateCredentialByECSRole(res, &options.OTS)
	}

	return &w, nil
}

func NewOTSTableStoreClientWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*OTSTableStoreClientWrapper, error) {
	var options OTSTableStoreClientWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewOTSTableStoreClientWrapperWithOptions(&options, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "NewOTSTableStoreClientWrapperWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Wrapper"), w.OnWrapperChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Retry"), w.OnRetryChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiter"), w.OnRateLimiterChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("ParallelController"), w.OnParallelControllerChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("OTS"), func(cfg *config.Config) error {
		var options OTSOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}

		if options.AccessKeyID != "" {
			client := tablestore.NewClient(options.Endpoint, options.InstanceName, options.AccessKeyID, options.AccessKeySecret)
			if _, err := client.ListTable(); err != nil {
				return errors.Wrap(err, "tablestore.TableStoreClient.ListTable failed")
			}
			w.obj = client
			return nil
		}

		res, err := alics.ECSMetaDataRamSecurityCredentials()
		if err != nil {
			return errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentials failed")
		}
		client := tablestore.NewClientWithConfig(options.Endpoint, options.InstanceName, res.AccessKeyID, res.AccessKeySecret, res.SecurityToken, nil)
		if _, err := client.ListTable(); err != nil {
			return errors.Wrap(err, "tablestore.TableStoreClient.ListTable failed")
		}
		w.obj = client
		go w.UpdateCredentialByECSRole(res, &options)
		return nil
	})

	return w, err
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
