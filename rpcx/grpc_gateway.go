package rpcx

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type GRPCGatewayInterceptorOptions struct {
	Headers          []string `dft:"X-Request-Id"`
	Validators       []string
	PrivateIP        string
	Hostname         string
	RequestIDMetaKey string `dft:"x-request-id"`
	ContentType      string `dft:"application/json"`
	PascalNameKey    bool
	EnableTracing    bool
	UseFieldKey      bool
	MarshalOmitempty bool
}

type GRPCGatewayInterceptor struct {
	grpcInterceptor *GRPCInterceptor
	muxInterceptor  *MuxInterceptor
}

func NewGRPCGatewayInterceptorWithOptions(options *GRPCGatewayInterceptorOptions) (*GRPCGatewayInterceptor, error) {
	grpcInterceptor, err := NewGRPCInterceptorWithOptions(&GRPCInterceptorOptions{
		Headers:          options.Headers,
		PrivateIP:        options.PrivateIP,
		Hostname:         options.Hostname,
		Validators:       options.Validators,
		PascalNameKey:    options.PascalNameKey,
		RequestIDMetaKey: options.RequestIDMetaKey,
		EnableTracing:    options.EnableTracing,
	})
	if err != nil {
		return nil, errors.Wrap(err, "NewGRPCInterceptorWithOptions failed")
	}

	muxInterceptor, err := NewMuxInterceptorWithOptions(&MuxInterceptorOptions{
		Headers:          options.Headers,
		ContentType:      options.ContentType,
		UseFieldKey:      options.UseFieldKey,
		RequestIDMetaKey: options.RequestIDMetaKey,
		MarshalOmitempty: options.MarshalOmitempty,
	})

	return &GRPCGatewayInterceptor{
		grpcInterceptor: grpcInterceptor,
		muxInterceptor:  muxInterceptor,
	}, nil
}

func (g *GRPCGatewayInterceptor) SetLogger(logger Logger) {
	g.grpcInterceptor.SetLogger(logger)
}

func (g *GRPCGatewayInterceptor) ServerOption() grpc.ServerOption {
	return g.grpcInterceptor.ServerOption()
}

func (g *GRPCGatewayInterceptor) DialOptions() []grpc.DialOption {
	return g.grpcInterceptor.DialOptions()
}

func (g *GRPCGatewayInterceptor) ServeMuxOptions() []runtime.ServeMuxOption {
	return g.muxInterceptor.ServeMuxOptions()
}
