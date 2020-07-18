package config

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOTSProvider(t *testing.T) {
	Convey("TestOTSProvider", t, func() {
		otsCli := tablestore.NewClient(
			"https://imm-dev-hl.cn-shanghai.ots.aliyuncs.com",
			"imm-dev-hl",
			"xx",
			"xx",
		)
		provider, err := NewOTSProvider(otsCli, "IMMConfigV2", "test", 100*time.Millisecond)
		So(err, ShouldBeNil)
		{
			So(provider.Dump([]byte("hello world")), ShouldBeNil)
			buf, _, err := OTSGetRow(otsCli, "IMMConfigV2", "test")
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, "hello world")
		}
		{
			ctx, cancel := context.WithCancel(context.Background())
			So(provider.EventLoop(ctx), ShouldBeNil)

			for len(provider.Events()) != 0 {
				<-provider.Events()
			}

			for i := 0; i < 5; i++ {
				So(OTSPutRow(otsCli, "IMMConfigV2", "test", []byte(fmt.Sprintf("hello world %v", i))), ShouldBeNil)
				<-provider.Events()
				buf, err := provider.Load()
				So(err, ShouldBeNil)
				fmt.Println(string(buf))
				So(string(buf), ShouldEqual, fmt.Sprintf("hello world %v", i))
				time.Sleep(time.Second)
			}

			cancel()
		}
	})
}
