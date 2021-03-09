# astx

使用 `go/ast` 相关库的语法解析功能，做一些代码自动生成工具

## astx

语法解析模块

- `ParseFunction`:从目录中解析出其中的函数

## wrap

[wrap](../wrap) 代码自动生成工具

- 支持 retry
- 支持 trace
- 支持 metric
- 支持 rateLimiter
- 支持 parallelController
- 支持错误在返回体中（`github.com/go-redis/redis`/`github.com/jinzhu/gorm`）
- 支持继承（`github.com/go-redis/redis`）

### wrap 工具

可通过如下方式工具获取

```shell
go get github.com/hatlonely/go-kit/cmd/gen
```

查看帮助

```shell
gen wrap -h
```

代码生成示例

```shell
gen wrap --sourcePath vendor \
  --packagePath "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore" \
  --packageName tablestore \
  --classPrefix OTS \
  --starClasses TableStoreClient \
  --rule.mainClass.include "^TableStoreClient$$" \
  --output $@
  
gen wrap --sourcePath vendor \
  --packagePath "github.com/go-redis/redis" \
  --packageName redis \
  --classPrefix Redis \
  --errorField "Err()" \
  --inherit '{"Client": ["cmdable", "baseClient"], "ClusterClient": ["cmdable"]}' \
  --starClasses Client,ClusterClient \
  --rule.mainClass.include "^(?i:(Client)|(ClusterClient))$$" \
  --rule.errorInResult.include "^\*redis\..*Cmd$$" \
  --output $@
```

### wrap 类模板

```go
type OTSTableStoreClientWrapper struct {
	obj                *tablestore.TableStoreClient
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}
```

### wrap 方法代码模板

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
			span, _ = opentracing.StartSpanFromContext(ctx, "tablestore.TableStoreClient.GetRow")
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

### wrap 方法

- wrap 相关类的所有公开方法
- 和原生类保持一致的方法接口
    * 如果方法含有 `context.Context` 参数，接口和原生接口一致
    * 如果方法没有 `context.Context` 参数，增加 `context.Context` 作为 wrap 方法的第一个参数
  
```go
func (tableStoreClient *TableStoreClient) GetRow(request *GetRowRequest) (*GetRowResponse, error)
func (w *OTSTableStoreClientWrapper) GetRow(ctx context.Context, request *tablestore.GetRowRequest) (*tablestore.GetRowResponse, error)


func (coll *Collection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*Cursor, error)
func (w *MongoCollectionWrapper) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
```

### 返回 wrap 类

原生类方法返回另一个原生类的场景，返回对应的 wrap 类

比如 `github.com/aliyun/aliyun-oss-go-sdk/oss` 中如下函数，同时 wrap `Client` 和 `Bucket` 类

```go
func (client Client) Bucket(bucketName string) (*Bucket, error)
func (w *OSSClientWrapper) Bucket(bucketName string) (*OSSBucketWrapper, error)
```

### 返回体错误

比如 `github.com/go-redis/redis` 中如下函数的错误包含在 `Status.Cmd` 中

```go
func (c *cmdable) Set(key string, value interface{}, expiration time.Duration) *StatusCmd
```

可通过设置 `ErrorField` 以及 `Rule.ErrorInResult` 指定错误字段，用于 retry

### 继承

`github.com/go-redis/redis` 中 `redis.Client` 继承了 `cmdable`

```go
func (c *cmdable) Set(key string, value interface{}, expiration time.Duration) *StatusCmd
```

可通过设置 `Inherit` 字段，生成继承类的 wrap 方法

### 跳过私有参数方法

`github.com/jinzhu/gorm` 中如下方法的参数为私有，无法 wrap

```go
func (s *DB) SetLogger(log logger)
```

可通过设置 Rule.Function 设置为 `'{"DB": {"exclude": "^SetLogger$$"}}'` 跳过该函数
