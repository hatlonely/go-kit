package service

import (
	"context"
	"time"

	"github.com/hatlonely/go-kit/rpcx"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"

	"github.com/hatlonely/go-kit/examples/rpcx/api/gen/go/api"
)

type Options struct {
	SleepTime time.Duration
}

func NewExampleServiceWithOptions(options *Options) (*ExampleService, error) {
	return &ExampleService{
		options: options,
	}, nil
}

type ExampleService struct {
	api.ExampleServiceServer

	options *Options
}

func (s *ExampleService) Echo(ctx context.Context, req *api.EchoReq) (*api.EchoRes, error) {
	time.Sleep(s.options.SleepTime)

	return &api.EchoRes{
		Message: req.Message,
	}, nil
}

func (s *ExampleService) Add(ctx context.Context, req *api.AddReq) (*api.AddRes, error) {
	time.Sleep(s.options.SleepTime)

	if req.I1 > 100 || req.I2 > 100 {
		return nil, rpcx.NewError(errors.Errorf("parameter too large"), codes.InvalidArgument, "InvalidArgument", "parameter too large")
	}

	return &api.AddRes{
		Val: req.I1 + req.I2,
	}, nil
}
