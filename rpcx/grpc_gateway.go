package rpcx

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"

	"github.com/hatlonely/go-kit/logger"
	"github.com/hatlonely/go-kit/refx"
)

type GrpcGatewayOptions struct {
	HttpPort    int
	GrpcPort    int
	ExitTimeout time.Duration

	Headers          []string `dft:"X-Request-Id"`
	Validators       []string
	PrivateIP        string
	Hostname         string
	RequestIDMetaKey string `dft:"x-request-id"`
	ContentType      string `dft:"application/json"`
	PascalNameKey    bool
	EnableTrace      bool
	EnableMetric     bool
	UseFieldKey      bool
	MarshalOmitempty bool

	Jaeger jaegercfg.Configuration
}

type GrpcGateway struct {
	grpcInterceptor *GrpcInterceptor
	muxInterceptor  *MuxInterceptor

	grpcServer  *grpc.Server
	muxServer   *runtime.ServeMux
	httpServer  *http.Server
	traceCloser io.Closer

	options *GrpcGatewayOptions

	appLog Logger
	appRpc Logger
}

func NewGrpcGatewayWithOptions(options *GrpcGatewayOptions) (*GrpcGateway, error) {
	grpcInterceptor, err := NewGrpcInterceptorWithOptions(&GrpcInterceptorOptions{
		Headers:          options.Headers,
		PrivateIP:        options.PrivateIP,
		Hostname:         options.Hostname,
		Validators:       options.Validators,
		PascalNameKey:    options.PascalNameKey,
		RequestIDMetaKey: options.RequestIDMetaKey,
		EnableTrace:      options.EnableTrace,
	})
	if err != nil {
		return nil, errors.Wrap(err, "NewGrpcGatewayWithOptions failed")
	}

	muxInterceptor, _ := NewMuxInterceptorWithOptions(&MuxInterceptorOptions{
		Headers:          options.Headers,
		ContentType:      options.ContentType,
		UseFieldKey:      options.UseFieldKey,
		RequestIDMetaKey: options.RequestIDMetaKey,
		MarshalOmitempty: options.MarshalOmitempty,
	})

	g := &GrpcGateway{
		grpcInterceptor: grpcInterceptor,
		muxInterceptor:  muxInterceptor,
		options:         options,
		appLog:          logger.NewStdoutTextLogger(),
		appRpc:          logger.NewStdoutJsonLogger(),
	}

	g.grpcServer = grpc.NewServer(grpcInterceptor.ServerOption())
	g.muxServer = runtime.NewServeMux(muxInterceptor.ServeMuxOptions()...)

	if options.EnableTrace {
		tracer, closer, err := options.Jaeger.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			return nil, errors.Wrap(err, "Jaeger.NewTracer failed")
		}
		opentracing.SetGlobalTracer(tracer)
		g.traceCloser = closer
	}

	return g, nil
}

func (g *GrpcGateway) SetLogger(log, rpc Logger) {
	g.grpcInterceptor.SetLogger(rpc)
	g.appLog = log
	g.appRpc = rpc
}

func (g *GrpcGateway) GRPCServer() *grpc.Server {
	return g.grpcServer
}

func (g *GrpcGateway) MuxServer() *runtime.ServeMux {
	return g.muxServer
}

func (g *GrpcGateway) RegisterServiceHandlerFunc(fun func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)) error {
	return fun(context.Background(), g.muxServer, fmt.Sprintf("0.0.0.0:%v", g.options.GrpcPort), g.grpcInterceptor.DialOptions())
}

func (g *GrpcGateway) Run() {
	go func() {
		address, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", g.options.GrpcPort))
		refx.Must(err)
		refx.Must(g.grpcServer.Serve(address))
	}()

	var handler http.Handler
	handler = g.muxServer
	if g.options.EnableMetric {
		handler = MetricWrapper(handler)
	}
	if g.options.EnableTrace {
		handler = TraceWrapper(handler)
	}
	g.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%v", g.options.HttpPort),
		Handler: handler,
	}
	go func() {
		if err := g.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			g.appLog.Warnf("httpServer.ListenAndServe failed. err: [%v]", err)
		}
	}()
	g.appLog.Infof("server start success. httpPort: [%v], grpcPort: [%v]", g.options.HttpPort, g.options.GrpcPort)
}

func (g *GrpcGateway) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), g.options.ExitTimeout)
	defer cancel()
	if err := g.httpServer.Shutdown(ctx); err != nil {
		g.appLog.Warnf("httServer.Shutdown failed, err: [%v]", err)
	}
	g.grpcServer.Stop()
	if g.traceCloser != nil {
		_ = g.traceCloser.Close()
	}
}
