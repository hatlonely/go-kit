package rpcx

import (
	"context"
	"fmt"
	"runtime/debug"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/metadata"
)

func TestMetaDataIncomingGet(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"key1": "val1",
		"key2": "val2",
	}))

	Convey("TestMetaDataIncomingGet", t, func() {
		So(MetaDataIncomingGet(ctx, "key1"), ShouldEqual, "val1")
		So(MetaDataIncomingGet(ctx, "key2"), ShouldEqual, "val2")
	})
}

func TestMetaDataIncomingSet(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{}))

	Convey("TestMetaDataIncomingSet", t, func() {
		MetaDataIncomingSet(ctx, "key1", "val1")
		MetaDataIncomingSet(ctx, "key2", "val2")
		So(MetaDataIncomingGet(ctx, "key1"), ShouldEqual, "val1")
		So(MetaDataIncomingGet(ctx, "key2"), ShouldEqual, "val2")
	})
}

func TestCtxGetSet(t *testing.T) {
	ctx := NewRPCXContext(context.Background())

	Convey("TestCtxSet", t, func() {
		CtxSet(ctx, "key1", "val1")
		CtxSet(ctx, "key2", 2)
		So(CtxGet(ctx, "key1"), ShouldEqual, "val1")
		So(CtxGet(ctx, "key2"), ShouldEqual, 2)
	})
}

func TestPrivateIP(t *testing.T) {
	Convey("TestPrivateIP", t, func() {
		fmt.Println(PrivateIP())
		fmt.Println(Hostname())
	})
}

func TestGRPCUnaryInterceptor(t *testing.T) {
	Convey("TestGRPCUnaryInterceptor", t, func() {
		fun := func() (err error) {
			defer func() {
				if perr := recover(); perr != nil {
					err = NewInternalError(errors.Wrap(fmt.Errorf("%v\n%v", string(debug.Stack()), perr), "panic"))
				}
			}()
			panic("hello")
		}
		So(fun(), ShouldNotBeNil)
	})
}
