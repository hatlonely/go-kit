package config

import (
	"context"
	"errors"
	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestACMProvider(t *testing.T) {
	Convey("TestACMProvider", t, func() {
		Convey("test normal", func() {
			patches := ApplyMethod(reflect.TypeOf(&config_client.ConfigClient{}), "GetConfig", func(
				client *config_client.ConfigClient, param vo.ConfigParam) (string, error) {
				So(param.DataId, ShouldEqual, "test-data-id")
				So(param.Group, ShouldEqual, "test-group")
				return "hello world", nil
			}).ApplyMethod(reflect.TypeOf(&config_client.ConfigClient{}), "PublishConfig", func(
				client *config_client.ConfigClient, params vo.ConfigParam) (bool, error) {
				So(params.DataId, ShouldEqual, "test-data-id")
				So(params.Group, ShouldEqual, "test-group")
				return true, nil
			})
			defer patches.Reset()

			provider, err := NewACMProviderWithOptions(&ACMProviderOptions{
				ClientConfig: constant.ClientConfig{
					Endpoint:    "acm.aliyun.com:8080",
					NamespaceId: "test-namespace",
					AccessKey:   "test-access-key-id",
					SecretKey:   "test-access-key-secret",
					TimeoutMs:   5000,
				},
				DataID: "test-data-id",
				Group:  "test-group",
			})
			So(err, ShouldBeNil)

			Convey("Dump", func() {
				So(provider.Dump([]byte("hello world")), ShouldBeNil)
			})

			Convey("Load", func() {
				buf, err := provider.Load()
				So(err, ShouldBeNil)
				So(string(buf), ShouldEqual, "hello world")
			})
		})

		Convey("test get config failed", func() {
			patches := ApplyMethod(reflect.TypeOf(&config_client.ConfigClient{}), "GetConfig", func(
				client *config_client.ConfigClient, param vo.ConfigParam) (string, error) {
				return "", errors.New("timeout")
			})
			defer patches.Reset()

			provider, err := NewACMProviderWithOptions(&ACMProviderOptions{
				ClientConfig: constant.ClientConfig{
					Endpoint:    "acm.aliyun.com:8080",
					NamespaceId: "test-namespace",
					AccessKey:   "test-access-key-id",
					SecretKey:   "test-access-key-secret",
					TimeoutMs:   5000,
				},
				DataID: "test-data-id",
				Group:  "test-group",
			})
			So(err, ShouldNotBeNil)
			So(provider, ShouldBeNil)
		})

		Convey("test put config failed", func() {
			patches := ApplyMethod(reflect.TypeOf(&config_client.ConfigClient{}), "GetConfig", func(
				client *config_client.ConfigClient, param vo.ConfigParam) (string, error) {
				return "hello world", nil
			}).ApplyMethodSeq(reflect.TypeOf(&config_client.ConfigClient{}), "PublishConfig", []OutputCell{
				{Values: Params{true, errors.New("timeout")}},
				{Values: Params{false, nil}},
				{Values: Params{true, nil}},
			})
			defer patches.Reset()

			provider, err := NewACMProviderWithOptions(&ACMProviderOptions{
				ClientConfig: constant.ClientConfig{
					Endpoint:    "acm.aliyun.com:8080",
					NamespaceId: "test-namespace",
					AccessKey:   "test-access-key-id",
					SecretKey:   "test-access-key-secret",
					TimeoutMs:   5000,
				},
				DataID: "test-data-id",
				Group:  "test-group",
			})
			So(err, ShouldBeNil)
			So(provider.Dump([]byte("hello world")), ShouldNotBeNil)
			So(provider.Dump([]byte("hello world")), ShouldNotBeNil)
		})

		Convey("test listen", func() {
			patches := ApplyMethod(reflect.TypeOf(&config_client.ConfigClient{}), "GetConfig", func(
				client *config_client.ConfigClient, param vo.ConfigParam) (string, error) {
				So(param.DataId, ShouldEqual, "test-data-id")
				So(param.Group, ShouldEqual, "test-group")
				return "hello world", nil
			}).ApplyMethod(reflect.TypeOf(&config_client.ConfigClient{}), "PublishConfig", func(
				client *config_client.ConfigClient, params vo.ConfigParam) (bool, error) {
				So(params.DataId, ShouldEqual, "test-data-id")
				So(params.Group, ShouldEqual, "test-group")
				return true, nil
			}).ApplyMethod(reflect.TypeOf(&config_client.ConfigClient{}), "ListenConfig", func(
				client *config_client.ConfigClient, params vo.ConfigParam) (err error) {
				So(params.DataId, ShouldEqual, "test-data-id")
				So(params.Group, ShouldEqual, "test-group")
				go func() {
					params.OnChange("test-namespace", "test-group", "test-data-id", "hello world 1")
				}()
				return nil
			})
			defer patches.Reset()

			provider, err := NewACMProviderWithOptions(&ACMProviderOptions{
				ClientConfig: constant.ClientConfig{
					Endpoint:    "acm.aliyun.com:8080",
					NamespaceId: "test-namespace",
					AccessKey:   "test-access-key-id",
					SecretKey:   "test-access-key-secret",
					TimeoutMs:   5000,
				},
				DataID: "test-data-id",
				Group:  "test-group",
			})
			So(err, ShouldBeNil)

			So(provider.EventLoop(context.Background()), ShouldBeNil)

			<-provider.events
			buf, err := provider.Load()
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, "hello world 1")
		})
	})
}
