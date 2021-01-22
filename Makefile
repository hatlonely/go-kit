export GOPROXY=https://goproxy.cn
export GOPRIVATE=gitlab.alibaba-inc.com

define BUILD_VERSION
  version: $(shell git describe --tags)
gitremote: $(shell git remote -v | grep fetch | awk '{print $$2}')
   commit: $(shell git rev-parse HEAD)
 datetime: $(shell date '+%Y-%m-%d %H:%M:%S')
 hostname: $(shell hostname):$(shell pwd)
goversion: $(shell go version)
endef
export BUILD_VERSION

test: vendor
	go test -gcflags=all=-l -cover ./alics
	go test -gcflags=all=-l -cover ./bind
	go test -gcflags=all=-l -cover ./cast
	go test -gcflags=all=-l -cover ./config
	go test -gcflags=all=-l -cover ./flag
	go test -gcflags=all=-l -cover ./logger
	go test -gcflags=all=-l -cover ./refx
	go test -gcflags=all=-l -cover ./rpcx
	go test -gcflags=all=-l -cover ./strx
	go test -gcflags=all=-l -cover ./validator
	go test -gcflags=all=-l -cover ./ops

vendor:
	go mod tidy
	go mod vendor

build: build/bin/cfg build/bin/ops build/bin/gen

build/bin/cfg: cmd/cfg/main.go $(wildcard config/*.go) Makefile vendor
	mkdir -p build/bin
	go build -ldflags "-X 'main.Version=$$BUILD_VERSION'" -o $@ $<

build/bin/ops: cmd/ops/main.go $(wildcard ops/*.go) Makefile vendor
	mkdir -p build/bin
	go build -ldflags "-X 'main.Version=$$BUILD_VERSION'" -o $@ $<

build/bin/gen: cmd/gen/main.go $(wildcard astx/*.go) Makefile vendor
	mkdir -p build/bin
	go build -ldflags "-X 'main.Version=$$BUILD_VERSION'" -o $@ $<

wrap: wrap/autogen_ots.go wrap/autogen_kms.go wrap/autogen_acm.go wrap/autogen_oss.go

wrap/autogen_ots.go: build/bin/gen vendor $(wildcard astx/*.go)
	build/bin/gen --goPath vendor \
		--pkgPath "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore" \
		--package tablestore \
		--classPrefix OTS \
		--classes TableStoreClient \
		--output $@

wrap/autogen_kms.go: build/bin/gen vendor $(wildcard astx/*.go)
	build/bin/gen --goPath vendor \
		--pkgPath "github.com/aliyun/alibaba-cloud-sdk-go/services/kms" \
		--package kms \
		--classPrefix KMS \
		--classes Client \
		--output $@

wrap/autogen_acm.go: build/bin/gen vendor $(wildcard astx/*.go)
	build/bin/gen --goPath vendor \
		--pkgPath "github.com/nacos-group/nacos-sdk-go/clients/config_client" \
		--package config_client \
		--classPrefix ACM \
		--classes ConfigClient \
		--output $@

wrap/autogen_oss.go: build/bin/gen vendor $(wildcard astx/*.go)
	build/bin/gen --goPath ~/hatlonely/gitlab.alibaba-inc.com/hatlonely.hl/convert-wrapper/vendor \
		--pkgPath "github.com/aliyun/aliyun-oss-go-sdk/oss" \
		--package oss \
		--classPrefix OSS \
		--classes Client,Bucket \
		--rule.trace '{"Client": {"exclude": "^Bucket$$"}}' \
		--rule.retry '{"Client": {"exclude": "^Bucket$$"}}' \
		--output $@
