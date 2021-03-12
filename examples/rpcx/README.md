# rpcx example

## 命令

- 添加依赖 `git submodule add git@github.com:hatlonely/rpc-api.git`
- 代码生成 `make codegen`
- 编译运行 `make build`

## 测试运行

启动服务

```shell
go run cmd/main.go -c config/app.json
```

发送请求

```shell
curl http://127.0.0.1/v1/echo?message=helloworld
curl -XPOST -d '{"i1":12, "i2": 34}' http://127.0.0.1/v1/add
curl -XPOST -d '{"i1":123, "i2": 456}' http://127.0.0.1/v1/add
```
