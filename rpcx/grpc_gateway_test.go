package rpcx

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"

	"github.com/hatlonely/go-kit/examples/rpcx/api/gen/go/api"
)

type ExampleService struct {
	api.ExampleServiceServer
}

func (s *ExampleService) Echo(ctx context.Context, req *api.EchoReq) (*api.EchoRes, error) {
	return &api.EchoRes{
		Message: req.Message,
	}, nil
}

func (s *ExampleService) Add(ctx context.Context, req *api.AddReq) (*api.AddRes, error) {
	if req.I1 < 0 {
		return nil, NewError(errors.New("i1 should be positive"), codes.InvalidArgument, "InvalidArgument", "i1 should be positive")
	}
	if req.I2 < 0 {
		return nil, errors.Wrap(NewError(errors.New("i2 should be positive"), codes.InvalidArgument, "BadRequest", "i2 should be positive"), "wrap error")
	}
	if req.I1 == 101 {
		return nil, errors.New("hello world")
	}
	if req.I1 == 102 {
		panic("panic")
	}

	return &api.AddRes{
		Val: req.I1 + req.I2,
	}, nil
}

func TestGrpcGateway_AddGrpcPreHandler(t *testing.T) {
	Convey("TestGrpcGateway", t, func() {
		server, err := NewGrpcGatewayWithOptions(&GrpcGatewayOptions{
			HttpPort:         80,
			GrpcPort:         6080,
			EnableTrace:      false,
			EnableMetric:     false,
			EnablePprof:      false,
			Validators:       []string{"Default"},
			RequestIDMetaKey: "x-request-id",
			Headers:          []string{"X-Request-Id", "X-User-Id"},
		})
		So(err, ShouldBeNil)

		api.RegisterExampleServiceServer(server.GRPCServer(), &ExampleService{})
		So(server.RegisterServiceHandlerFunc(api.RegisterExampleServiceHandlerFromEndpoint), ShouldBeNil)

		server.AddGrpcPreHandler(func(ctx context.Context, req interface{}) error {
			if MetaDataIncomingGet(ctx, "x-user-id") == "" {
				return NewError(errors.New("permission deny"), codes.PermissionDenied, "Forbidden", "permission deny")
			}
			return nil
		})

		go server.Run()
		defer server.Stop()

		time.Sleep(100 * time.Millisecond)

		client := NewHttpClient()

		Convey("pass", func() {
			var res api.EchoRes
			resMeta := map[string]string{}
			So(client.Get(
				"http://127.0.0.1/v1/echo",
				map[string]string{"message": "hello world"},
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				nil,
				&resMeta,
				&res,
			), ShouldBeNil)
			So(resMeta["X-Request-Id"], ShouldResemble, "test-request-id")
			So(resMeta["X-User-Id"], ShouldResemble, "121231")
			So(res.Message, ShouldEqual, "hello world")
		})

		Convey("permission deny", func() {
			var res api.AddRes
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id"},
				&api.AddReq{I1: -12, I2: 34},
				&resMeta,
				&res,
			)
			e := err.(*HttpError)
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Status, ShouldEqual, 403)
			So(e.Code, ShouldEqual, "Forbidden")
			So(e.Message, ShouldEqual, "permission deny")
		})
	})
}

func TestGrpcGateway(t *testing.T) {
	Convey("TestGrpcGateway", t, func() {
		server, err := NewGrpcGatewayWithOptions(&GrpcGatewayOptions{
			HttpPort:         80,
			GrpcPort:         6080,
			EnableTrace:      false,
			EnableMetric:     false,
			EnablePprof:      false,
			Validators:       []string{"Default"},
			RequestIDMetaKey: "x-request-id",
			Headers:          []string{"X-Request-Id", "X-User-Id"},
		})
		So(err, ShouldBeNil)

		api.RegisterExampleServiceServer(server.GRPCServer(), &ExampleService{})
		So(server.RegisterServiceHandlerFunc(api.RegisterExampleServiceHandlerFromEndpoint), ShouldBeNil)

		go server.Run()
		defer server.Stop()

		time.Sleep(100 * time.Millisecond)

		client := NewHttpClient()

		Convey("test echo", func() {
			var res api.EchoRes
			resMeta := map[string]string{}
			So(client.Get(
				"http://127.0.0.1/v1/echo",
				map[string]string{"message": "hello world"},
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				nil,
				&resMeta,
				&res,
			), ShouldBeNil)
			So(resMeta["X-Request-Id"], ShouldResemble, "test-request-id")
			So(resMeta["X-User-Id"], ShouldResemble, "121231")
			So(res.Message, ShouldEqual, "hello world")
		})

		Convey("test add", func() {
			var res api.AddRes
			resMeta := map[string]string{}
			So(client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				&api.AddReq{I1: 12, I2: 34},
				&resMeta,
				&res,
			), ShouldBeNil)
			So(resMeta["X-Request-Id"], ShouldResemble, "test-request-id")
			So(resMeta["X-User-Id"], ShouldResemble, "121231")
			So(res.Val, ShouldEqual, 46)
		})

		Convey("test error", func() {
			var res api.AddRes
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				&api.AddReq{I1: -12, I2: 34},
				&resMeta,
				&res,
			)
			e := err.(*HttpError)
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Status, ShouldEqual, 400)
			So(e.Code, ShouldEqual, "InvalidArgument")
			So(e.Message, ShouldEqual, "i1 should be positive")
		})

		Convey("test wrap error", func() {
			var res api.AddRes
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				&api.AddReq{I1: 12, I2: -34},
				&resMeta,
				&res,
			)
			e := err.(*HttpError)
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Status, ShouldEqual, 400)
			So(e.Code, ShouldEqual, "BadRequest")
			So(e.Message, ShouldEqual, "i2 should be positive")
		})

		Convey("test internal error", func() {
			var res api.AddRes
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				&api.AddReq{I1: 101, I2: 34},
				&resMeta,
				&res,
			)
			e := err.(*HttpError)
			fmt.Println(e)
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Status, ShouldEqual, 500)
			So(e.Code, ShouldEqual, "InternalError")
			So(e.Message, ShouldEqual, "unknown error")
		})

		Convey("test panic", func() {
			var res api.AddRes
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				&api.AddReq{I1: 102, I2: 34},
				&resMeta,
				&res,
			)
			e := err.(*HttpError)
			fmt.Println(e)
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Status, ShouldEqual, 500)
			So(e.Code, ShouldEqual, "InternalError")
			So(e.Message, ShouldEqual, "unknown error")
		})

		Convey("test not found", func() {
			var res interface{}
			err := client.Get(
				"http://127.0.0.1/v1/found",
				map[string]string{"message": "hello world"},
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				nil,
				nil,
				&res,
			)
			So(reflect.TypeOf(err), ShouldEqual, reflect.TypeOf(&HttpError{}))
			e := err.(*HttpError)
			So(e.Status, ShouldEqual, 404)
			So(e.Code, ShouldEqual, http.StatusText(http.StatusNotFound))
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Message, ShouldEqual, http.StatusText(http.StatusNotFound))
		})

		Convey("test not implement", func() {
			var res api.EchoRes
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/echo",
				map[string]string{"message": "hello world"},
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				nil,
				&resMeta,
				&res,
			)
			e := err.(*HttpError)
			So(e.Status, ShouldEqual, 501)
			So(e.Code, ShouldEqual, http.StatusText(http.StatusNotImplemented))
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Message, ShouldEqual, "Method Not Allowed")
		})
	})
}
