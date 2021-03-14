# rpcx


基于 grpc-gateway 的封装了 trace/metric/logger/validator/rateLimiter/parallelController/cors 等相关功能

## Feature

- 支持 trace
- 支持 metric
- 自动 logger
- 支持 validator
- 支持 rateLimiter
- 支持 parallelController
- 支持 cors

## Quick Start

[example](../examples/rpcx)

```go
func main() {
  	grpcGateway, err := rpcx.NewGrpcGatewayWithOptions(&options.GrpcGateway, refx.WithCamelName())
	refx.Must(err)
	
	api.RegisterExampleServiceServer(grpcGateway.GRPCServer(), svc)
	refx.Must(grpcGateway.RegisterServiceHandlerFunc(api.RegisterExampleServiceHandlerFromEndpoint))
	grpcGateway.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	grpcGateway.Stop()
	infoLog.Info("server exit properly")
}

```

配置文件

```js
{
  "name": "example",
  "httpPort": 80,
  "grpcPort": 6080,
  "exitTimeout": "10s", // 服务退出最长等待时间，默认 10s
  "headers": ["X-Request-Id"], // 请求透传包头，默认透传 x- 开头包头
  "validators": ["Default"], // 请求包体校验，默认不校验
  "privateIP": "",
  "hostname": "",
  "requestIDMetaKey": "x-request-id", // request id 包头 key，默认 x-request-id
  "usePascalNameLogKey": false, // rpc 日志使用大写驼峰风格，默认小写驼峰
  "usePascalNameErrKey": false, // 错误返回包体使用大写驼峰风格，默认小写驼峰
  "marshalUseProtoNames": true, // 序列化返回时，使用 proto 文件中的名字，默认使用小写驼峰
  "marshalEmitUnpopulated": false, // 序列化返回时，返回空字段，默认不返回
  "unmarshalDiscardUnknown": true, // 反序列化请求时，丢弃未知的字段，默认返回请求包体错误
  "enableTrace": true,
  "enableMetric": true,
  "enablePprof": true,
  "trace": {
    "constTags": {
      "key1": "val1"
    }
  },
  "metric": {
    "buckets": [0.1, 1, 10, 100],
    "constLabels": {
      "key1": "val1"
    }
  },
  "enableCors": false,
  "cors": {
    "allowAll": "false",
    "allowRegex": [],
    "allowOrigin": [],
    "allowMethod": ["GET", "HEAD", "POST", "PUT", "DELETE"],
    "allowHeader": ["Content-Type", "Accept"]
  },
  "rateLimiterHeader": "x-user-id",
  "rateLimiter": {
    "type": "RedisRateLimiterInstance"
  },
  "parallelControllerHeader": "x-user-id",
  "parallelController": {
    "type": "RedisTimedParallelControllerInstance"
  },
  "jaeger": {
    "serviceName": "example",
    "sampler": {
      "type": "const",
      "param": 1
    },
    "reporter": {
      "logSpans": false
    }
  }
}
```