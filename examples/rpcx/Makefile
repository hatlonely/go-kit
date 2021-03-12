binary=rpcx-example
repository=rpcx-example
version=$(shell git describe --tags | awk '{print(substr($$0,2,length($$0)))}')
export GOPROXY=https://goproxy.cn

define BUILD_VERSION
  version: $(shell git describe --tags)
gitremote: $(shell git remote -v | grep fetch | awk '{print $$2}')
   commit: $(shell git rev-parse HEAD)
 datetime: $(shell date '+%Y-%m-%d %H:%M:%S')
 hostname: $(shell hostname):$(shell pwd)
goversion: $(shell go version)
endef
export BUILD_VERSION

.PHONY: build
build: cmd/main.go $(wildcard internal/*/*.go) Makefile vendor
	mkdir -p build/bin && mkdir -p build/config
	go build -ldflags "-X 'main.Version=$$BUILD_VERSION'" -o build/bin/${binary} cmd/main.go
	cp config/app.json build/config/

clean:
	rm -rf build

vendor: go.mod go.sum
	go mod tidy
	go mod vendor

.PHONY: codegen
codegen: api/example.proto submodule
	if [ ! -z "$(shell docker ps --filter name=protobuf -q)" ]; then \
		docker stop protobuf; \
	fi
	docker run --name protobuf -d --rm hatlonely/protobuf:1.0.0 tail -f /dev/null
	docker exec protobuf mkdir -p api
	docker cp $< protobuf:/$<
	docker cp rpc-api protobuf:/
	docker exec protobuf bash -c "mkdir -p api/gen/go && mkdir -p api/gen/swagger"
	docker exec protobuf bash -c "protoc -Irpc-api -I. --go_out api/gen/go --go_opt paths=source_relative $<"
	docker exec protobuf bash -c "protoc -Irpc-api -I. --go-grpc_out api/gen/go --go-grpc_opt paths=source_relative $<"
	docker exec protobuf bash -c "protoc -Irpc-api -I. --grpc-gateway_out api/gen/go --grpc-gateway_opt logtostderr=true,paths=source_relative $<"
	docker exec protobuf bash -c "protoc -Irpc-api -I. --openapiv2_out api/gen/swagger --openapiv2_opt logtostderr=true $<"
	docker cp protobuf:api/gen api
	docker stop protobuf

.PHONY: submodule
submodule:
	git submodule init
	git submodule update

.PHONY: image
image:
	docker build --tag=hatlonely/${repository}:${version} .
