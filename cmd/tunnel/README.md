# tunnel

tunnel 是一个反向访问工具，在服务端无法提供固定可访问地址的场景下，可将 tunnel-agent 部署在服务端，
tunnel-agent 会与 tunnel-server 建立一条通讯隧道，并由 tunnel-server 对外提供和原服务一致的服务

## 工作模式

```
|--------------|       |--------------|
|    server    | <-X-- |    client    |
|--------------|       |--------------|
```

```
|--------------|       |---------------|
|    server    |       |               |
|--------------|       |               |
       ^               |               |       |--------------|
       |               | tunnel-server | <---- |    client    |
       |               |               |       |--------------|
|--------------|       |               |
| tunnel-agent | ----> |               |
|--------------|       |---------------|
```

## 工具获取

**go get 获取**

```shell
go get -u github.com/hatlonely/go-kit/cmd/tunnel/tunnel-server
go get -u github.com/hatlonely/go-kit/cmd/tunnel/tunnel-agent
export PATH=$PATH:$GOPATH/bin
```

**源码编译**

```shell
git clone https://github.com/hatlonely/go-kit.git
make build
```

## Quick Start

```shell
# 在 8000 端口启动一个测试的 http 服务
python -m SimpleHTTPServer 8000

# 在 8888 启动 tunnel-server 服务，tunnel 端口 5080
tunnel-server --server.tunnelPort 5080 --server.serverPort 8888

# 启动 agent，将 tunnel 请求转发到 http 服务
tunnel-agent --agent.tunnelAddr 127.0.0.1:5080 --agent.serverAddr 127.0.0.1:8000

# 发送测试请求
curl 127.0.0.1:8888
```
