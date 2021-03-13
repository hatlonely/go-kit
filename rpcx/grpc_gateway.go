package rpcx

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/logger"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

type GrpcGatewayOptions struct {
	HttpPort    int
	GrpcPort    int
	ExitTimeout time.Duration `dft:"10s"` // 服务退出最长等待时间

	Headers                 []string `dft:"X-Request-Id"` // 请求透传包头，默认透传 x- 开头包头
	Validators              []string // 请求包体校验，默认不校验
	PrivateIP               string
	Hostname                string
	RequestIDMetaKey        string `dft:"x-request-id"` // request id 包头 key
	UsePascalNameLogKey     bool   // rpc 日志使用大写驼峰风格，默认小写驼峰
	UsePascalNameErrKey     bool   // 错误返回包体使用大写驼峰风格，默认小写驼峰
	MarshalUseProtoNames    bool   // 序列化返回时，使用 proto 文件中的名字，默认使用小写驼峰
	MarshalEmitUnpopulated  bool   // 序列化返回时，返回空字段，默认不返回
	UnmarshalDiscardUnknown bool   // 反序列化请求时，丢弃未知的字段，默认返回请求包体错误

	Name         string
	EnableTrace  bool
	EnableMetric bool
	EnablePprof  bool
	Trace        struct {
		ConstTags map[string]string
	}
	Metric struct {
		Buckets     []float64
		ConstLabels map[string]string
	}

	RateLimiterHeader        string
	ParallelControllerHeader string
	RateLimiter              micro.RateLimiterOptions
	ParallelController       micro.ParallelControllerOptions
	Jaeger                   jaegercfg.Configuration
}

type GrpcGateway struct {
	grpcInterceptor *GrpcInterceptor
	muxInterceptor  *MuxInterceptor

	grpcServer  *grpc.Server
	muxServer   *runtime.ServeMux
	httpServer  *http.Server
	traceCloser io.Closer

	httpHandlerMap map[string]http.Handler

	options *GrpcGatewayOptions

	appLog Logger
	appRpc Logger
}

func NewGrpcGatewayWithConfig(cfg *config.Config, opts ...refx.Option) (*GrpcGateway, error) {
	var options GrpcGatewayOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.WithMessage(err, "cfg.Unmarshal failed")
	}
	g, err := NewGrpcGatewayWithOptions(&options)
	if err != nil {
		return nil, errors.WithMessage(err, "NewGrpcGatewayWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiter"), func(cfg *config.Config) error {
		var options micro.RateLimiterOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.WithMessage(err, "cfg.Unmarshal failed")
		}
		rateLimiter, err := micro.NewRateLimiterWithOptions(&options, opts...)
		if err != nil {
			return errors.WithMessage(err, "micro.NewRateLimiterWithOptions failed")
		}
		g.grpcInterceptor.rateLimiter = rateLimiter
		return nil
	})
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("ParallelController"), func(cfg *config.Config) error {
		var options micro.ParallelControllerOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.WithMessage(err, "cfg.Unmarshal failed")
		}
		parallelController, err := micro.NewParallelControllerWithOptions(&options, opts...)
		if err != nil {
			return errors.WithMessage(err, "micro.NewParallelControllerWithOptions failed")
		}
		g.grpcInterceptor.parallelController = parallelController
		return nil
	})

	return g, nil
}

func NewGrpcGatewayWithOptions(options *GrpcGatewayOptions, opts ...refx.Option) (*GrpcGateway, error) {
	grpcInterceptor, err := NewGrpcInterceptorWithOptions(&GrpcInterceptorOptions{
		Headers:                  options.Headers,
		PrivateIP:                options.PrivateIP,
		Hostname:                 options.Hostname,
		Validators:               options.Validators,
		UsePascalNameLogKey:      options.UsePascalNameLogKey,
		RequestIDMetaKey:         options.RequestIDMetaKey,
		RateLimiter:              options.RateLimiter,
		ParallelController:       options.ParallelController,
		RateLimiterHeader:        options.RateLimiterHeader,
		ParallelControllerHeader: options.ParallelControllerHeader,
		Name:                     options.Name,
		EnableTrace:              options.EnableTrace,
		EnableMetric:             options.EnableMetric,
		Trace:                    options.Trace,
		Metric:                   options.Metric,
	}, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "NewGrpcGatewayWithOptions failed")
	}

	muxInterceptor, _ := NewMuxInterceptorWithOptions(&MuxInterceptorOptions{
		Headers:                 options.Headers,
		UsePascalNameErrKey:     options.UsePascalNameErrKey,
		RequestIDMetaKey:        options.RequestIDMetaKey,
		MarshalUseProtoNames:    options.MarshalUseProtoNames,
		MarshalEmitUnpopulated:  options.MarshalEmitUnpopulated,
		UnmarshalDiscardUnknown: options.UnmarshalDiscardUnknown,
	})

	g := &GrpcGateway{
		grpcInterceptor: grpcInterceptor,
		muxInterceptor:  muxInterceptor,
		options:         options,
		httpHandlerMap:  map[string]http.Handler{},
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

func (g *GrpcGateway) SetLogger(log Logger, rpc Logger) {
	g.grpcInterceptor.SetLogger(log, rpc)
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

func (g *GrpcGateway) AddGrpcPreHandler(handler GrpcPreHandler) {
	g.grpcInterceptor.AddPreHandler(handler)
}

func (g *GrpcGateway) AddHttpHandler(path string, handler http.Handler) {
	g.httpHandlerMap[path] = handler
}

func (g *GrpcGateway) HandleHttp(method string, path string, handleFunc runtime.HandlerFunc) error {
	return g.muxServer.HandlePath(method, path, handleFunc)
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
		g.AddHttpHandler("/metrics", promhttp.Handler())
	}
	if g.options.EnablePprof {
		g.AddHttpHandler("/debug/pprof/", http.HandlerFunc(pprof.Index))
		g.AddHttpHandler("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		g.AddHttpHandler("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		g.AddHttpHandler("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		g.AddHttpHandler("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	}
	if g.options.EnableTrace {
		handler = TraceWrapper(handler)
	}
	if len(g.httpHandlerMap) != 0 {
		handler = MapHandlerWrapper(handler, g.httpHandlerMap)
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
