package rpcx

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

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
	return &api.AddRes{
		Val: req.I1 + req.I2,
	}, nil
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
			Headers:          []string{"X-Request-Id"},
		})
		So(err, ShouldBeNil)

		api.RegisterExampleServiceServer(server.GRPCServer(), &ExampleService{})
		So(server.RegisterServiceHandlerFunc(api.RegisterExampleServiceHandlerFromEndpoint), ShouldBeNil)

		go server.Run()
		defer server.Stop()

		var res api.EchoRes
		client := NewHttpClient()
		So(client.Get("http://127.0.0.1/v1/echo", map[string]string{"message": "hello world"}, nil, nil, nil, &res), ShouldBeNil)
		So(res.Message, ShouldEqual, "hello world")
	})
}
