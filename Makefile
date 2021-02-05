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

wrap: wrap/autogen_ots.go \
	wrap/autogen_kms.go \
	wrap/autogen_acm.go \
	wrap/autogen_oss.go \
	wrap/autogen_gorm.go \
	wrap/autogen_elasticsearch.go \
	wrap/autogen_mongo.go \
	wrap/autogen_redis.go \
	wrap/autogen_sts.go

wrap/autogen_ots.go: build/bin/gen vendor $(wildcard astx/*.go)
	build/bin/gen wrap --sourcePath vendor \
		--packagePath "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore" \
		--packageName tablestore \
		--classPrefix OTS \
		--starClasses TableStoreClient \
		--rule.createMetric.include "^TableStoreClient$$" \
		--rule.onWrapperChange.include "^TableStoreClient$$" \
		--rule.onRetryChange.include "^TableStoreClient$$" \
		--rule.onRateLimiterChange.include "^TableStoreClient$$" \
		--rule.onParallelControllerChange.include "^TableStoreClient$$" \
		--output $@

wrap/autogen_kms.go: build/bin/gen vendor $(wildcard astx/*.go)
	build/bin/gen wrap --sourcePath vendor \
		--packagePath "github.com/aliyun/alibaba-cloud-sdk-go/services/kms" \
		--packageName kms \
		--classPrefix KMS \
		--starClasses Client \
		--rule.createMetric.include "^Client$$" \
		--rule.onWrapperChange.include "^Client$$" \
		--rule.onRetryChange.include "^Client$$" \
		--rule.onRateLimiterChange.include "^Client$$" \
		--rule.onParallelControllerChange.include "^Client$$" \
		--output $@

wrap/autogen_sts.go: build/bin/gen vendor $(wildcard astx/*.go)
	build/bin/gen wrap --sourcePath vendor \
		--packagePath "github.com/aliyun/alibaba-cloud-sdk-go/services/sts" \
		--packageName sts \
		--classPrefix STS \
		--starClasses Client \
		--rule.createMetric.include "^Client$$" \
		--rule.onWrapperChange.include "^Client$$" \
		--rule.onRetryChange.include "^Client$$" \
		--rule.onRateLimiterChange.include "^Client$$" \
		--rule.onParallelControllerChange.include "^Client$$" \
		--output $@

wrap/autogen_acm.go: build/bin/gen vendor $(wildcard astx/*.go)
	build/bin/gen wrap --sourcePath vendor \
		--packagePath "github.com/nacos-group/nacos-sdk-go/clients/config_client" \
		--packageName config_client \
		--classPrefix ACM \
		--starClasses ConfigClient \
		--rule.createMetric.include "^ConfigClient$$" \
		--rule.onWrapperChange.include "^ConfigClient$$" \
		--rule.onRetryChange.include "^ConfigClient$$" \
		--rule.onRateLimiterChange.include "^ConfigClient$$" \
		--rule.onParallelControllerChange.include "^ConfigClient$$" \
		--output $@

wrap/autogen_oss.go: build/bin/gen vendor $(wildcard astx/*.go)
	build/bin/gen wrap --sourcePath vendor \
		--packagePath "github.com/aliyun/aliyun-oss-go-sdk/oss" \
		--packageName oss \
		--classPrefix OSS \
		--starClasses Client,Bucket \
		--rule.createMetric.include "^Client$$" \
		--rule.onWrapperChange.include "^Client$$" \
		--rule.onRetryChange.include "^Client$$" \
		--rule.onRateLimiterChange.include "^Client$$" \
		--rule.onParallelControllerChange.include "^Client$$" \
		--rule.trace '{"Client": {"exclude": "^Bucket$$"}}' \
		--rule.retry '{"Client": {"exclude": "^Bucket$$"}}' \
		--rule.metric '{"Client": {"exclude": "^Bucket$$"}}' \
		--rule.rateLimiter '{"Client": {"exclude": "^Bucket$$"}}' \
		--rule.parallelController '{"Client": {"exclude": "^Bucket$$"}}' \
		--output $@

wrap/autogen_gorm.go: build/bin/gen vendor $(wildcard astx/*.go)
	build/bin/gen wrap --sourcePath vendor \
		--packagePath "github.com/jinzhu/gorm" \
		--packageName gorm \
		--classPrefix GORM \
		--classes DB,Association \
		--errorField Error \
		--rule.createMetric.include "^DB$$" \
		--rule.onWrapperChange.include "^DB$$" \
		--rule.onRetryChange.include "^DB$$" \
		--rule.onRateLimiterChange.include "^DB$$" \
		--rule.onParallelControllerChange.include "^DB$$" \
		--rule.errorInResult.include "^(?i:(\*gorm.DB)|(\*gorm.Association))$$" \
		--rule.function '{"DB": {"exclude": "^SetLogger$$"}}' \
		--output $@

wrap/autogen_elasticsearch.go: build/bin/gen vendor $(wildcard astx/*.go)
	build/bin/gen wrap --sourcePath vendor \
		--packagePath "github.com/olivere/elastic/v7" \
		--packageName elastic \
		--classPrefix ES \
		--rule.starClass '{"include": "^(?i:(Client)|(.*Service))$$", "exclude": ".*"}' \
		--rule.createMetric.include "^Client$$" \
		--rule.onWrapperChange.include "^Client$$" \
		--rule.onRetryChange.include "^Client$$" \
		--rule.onRateLimiterChange.include "^Client$$" \
		--rule.onParallelControllerChange.include "^Client$$" \
		--rule.trace '{"default": {"exclude": ".*", "include": "^(Do)|(DoAsync)$$"}, "Client": {"exclude": ".*"}}' \
		--rule.retry '{"default": {"exclude": ".*", "include": "^(Do)|(DoAsync)$$"}, "Client": {"exclude": ".*"}}' \
		--rule.metric '{"default": {"exclude": ".*", "include": "^(Do)|(DoAsync)$$"}, "Client": {"exclude": ".*"}}' \
		--rule.rateLimiter '{"default": {"exclude": ".*", "include": "^(Do)|(DoAsync)$$"}, "Client": {"exclude": ".*"}}' \
		--rule.parallelController '{"default": {"exclude": ".*", "include": "^(Do)|(DoAsync)$$"}, "Client": {"exclude": ".*"}}' \
		--output $@

wrap/autogen_mongo.go: build/bin/gen vendor $(wildcard astx/*.go)
	build/bin/gen wrap --sourcePath vendor \
		--packagePath "go.mongodb.org/mongo-driver/mongo" \
		--packageName mongo \
		--classPrefix Mongo \
		--rule.starClass '{"include": "^(?i:(Client)|(Database)|(Collection))$$", "exclude": ".*"}' \
		--rule.createMetric.include "^Client$$" \
		--rule.onWrapperChange.include "^Client$$" \
		--rule.onRetryChange.include "^Client$$" \
		--rule.onRateLimiterChange.include "^Client$$" \
		--rule.onParallelControllerChange.include "^Client$$" \
		--rule.trace '{"Client": {"exclude": "^Database$$"}, "Database": {"exclude": "^Collection$$"}}' \
		--rule.metric '{"Client": {"exclude": "^Database$$"}, "Database": {"exclude": "^Collection$$"}}' \
		--rule.rateLimiter '{"Client": {"exclude": "^Database$$"}, "Database": {"exclude": "^Collection$$"}}' \
		--rule.parallelController '{"Client": {"exclude": "^Database$$"}, "Database": {"exclude": "^Collection$$"}}' \
		--output $@

wrap/autogen_redis.go: build/bin/gen vendor $(wildcard astx/*.go)
	build/bin/gen wrap --sourcePath vendor \
		--packagePath "github.com/go-redis/redis" \
		--packageName redis \
		--classPrefix Redis \
		--errorField "Err()" \
		--inherit '{"Client": ["cmdable", "baseClient"], "ClusterClient": ["cmdable"]}' \
		--starClasses Client,ClusterClient \
		--rule.createMetric.include "^(?i:(Client)|(ClusterClient))$$" \
		--rule.onWrapperChange.include "^(?i:(Client)|(ClusterClient))$$" \
		--rule.onRetryChange.include "^(?i:(Client)|(ClusterClient))$$" \
		--rule.onRateLimiterChange.include "^(?i:(Client)|(ClusterClient))$$" \
		--rule.onParallelControllerChange.include "^(?i:(Client)|(ClusterClient))$$" \
		--rule.errorInResult.include "^\*redis\..*Cmd$$" \
		--output $@
