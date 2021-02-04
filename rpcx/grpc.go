package rpcx

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	playgroundValidator "gopkg.in/go-playground/validator.v9"

	"github.com/hatlonely/go-kit/logger"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/strx"
	"github.com/hatlonely/go-kit/validator"
)

func NewGrpcInterceptorWithOptions(options *GrpcInterceptorOptions, opts ...refx.Option) (*GrpcInterceptor, error) {
	if options.Hostname == "" {
		options.Hostname = Hostname()
	}
	if options.PrivateIP == "" {
		options.PrivateIP = PrivateIP()
	}
	options.RequestIDMetaKey = strings.ToLower(options.RequestIDMetaKey)

	g := &GrpcInterceptor{
		options: options,
		appRpc:  logger.NewStdoutJsonLogger(),
	}
	g.requestIDKey = "requestID"
	g.hostnameKey = "hostname"
	g.privateIPKey = "privateIP"
	g.remoteIPKey = "remoteIP"
	g.clientIPKey = "clientIP"
	g.methodKey = "method"
	g.rpcCodeKey = "rpcCode"
	g.errCodeKey = "errCode"
	g.statusKey = "status"
	g.metaKey = "meta"
	g.reqKey = "req"
	g.ctxKey = "ctx"
	g.resKey = "res"
	g.errKey = "err"
	g.errStackKey = "errStack"
	g.resTimeMsKey = "resTimeMs"
	if options.PascalNameKey {
		for _, key := range []*string{
			&g.requestIDKey, &g.hostnameKey, &g.privateIPKey, &g.remoteIPKey, &g.clientIPKey, &g.methodKey, &g.rpcCodeKey, &g.errCodeKey,
			&g.statusKey, &g.metaKey, &g.reqKey, &g.ctxKey, &g.resKey, &g.errKey, &g.errStackKey, &g.resTimeMsKey,
		} {
			*key = strx.PascalName(*key)
		}
	}

	for _, v := range options.Validators {
		switch v {
		case "Playground":
			validate := playgroundValidator.New()
			g.validators = append(g.validators, validate.Struct)
		case "Default":
			g.validators = append(g.validators, validator.Validate)
		default:
			return nil, errors.Errorf("invalid validator type [%v]", v)
		}
	}

	rateLimiter, err := micro.NewRateLimiterWithOptions(&options.RateLimiter, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "micro.NewRateLimiterWithOptions failed")
	}
	g.rateLimiter = rateLimiter
	parallelCtl, err := micro.NewParallelControllerGroupWithOptions(&options.ParallelController, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "micro.NewParallelControllerGroupWithOptions failed")
	}
	g.parallelCtl = parallelCtl

	return g, nil
}

type GrpcInterceptor struct {
	options *GrpcInterceptorOptions

	rateLimiter micro.RateLimiter
	parallelCtl micro.ParallelControllerGroup

	validators  []func(interface{}) error
	preHandlers []func(ctx context.Context, req interface{}) error

	requestIDKey string
	hostnameKey  string
	privateIPKey string
	remoteIPKey  string
	clientIPKey  string
	methodKey    string
	rpcCodeKey   string
	errCodeKey   string
	statusKey    string
	metaKey      string
	reqKey       string
	ctxKey       string
	resKey       string
	errKey       string
	errStackKey  string
	resTimeMsKey string

	appRpc Logger
	appLog Logger
}

type GrpcPreHandler func(ctx context.Context, req interface{}) error

func (g *GrpcInterceptor) AddPreHandler(handler GrpcPreHandler) {
	g.preHandlers = append(g.preHandlers, handler)
}

func (g *GrpcInterceptor) SetLogger(log, rpc Logger) {
	g.appLog = log
	g.appRpc = rpc
}

func (g *GrpcInterceptor) DialOptions() []grpc.DialOption {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	if g.options.EnableTrace {
		opts = append(opts, grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
			),
		))
	}
	return opts
}

func (g *GrpcInterceptor) ServerOption() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		var requestID, remoteIP string
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			requestID = strings.Join(md.Get(g.options.RequestIDMetaKey), ",")
			if requestID == "" {
				requestID = uuid.NewV4().String()
				md.Set(g.options.RequestIDMetaKey, requestID)
			}
			remoteIP = strings.Split(strings.Join(md.Get("x-remote-addr"), ","), ":")[0]
		}

		if g.options.EnableTrace {
			spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(md))
			if err == nil || err == opentracing.ErrSpanContextNotFound {
				span := opentracing.GlobalTracer().StartSpan("GrpcInterceptor", ext.RPCServerOption(spanCtx))
				span.SetTag(g.methodKey, info.FullMethod)
				ctx = opentracing.ContextWithSpan(ctx, span)
				defer span.Finish()
			}
		}

		ctx = NewRPCXContext(ctx)

		ts := time.Now()
		defer func() {
			if perr := recover(); perr != nil {
				err = NewInternalError(errors.Wrap(fmt.Errorf("%v\n%v", string(debug.Stack()), perr), "panic"))
			}

			clientIP := ""
			if p, ok := peer.FromContext(ctx); ok && p != nil {
				clientIP = p.Addr.String()
			}

			rpcCode := codes.OK.String()
			errCode := "OK"
			status := http.StatusOK
			if err != nil {
				e := err.(*Error)
				rpcCode = e.Code.String()
				errCode = e.Detail.Code
				status = int(e.Detail.Status)
			}

			md, _ := metadata.FromIncomingContext(ctx)
			meta := map[string]string{}
			for key, val := range md {
				meta[key] = strings.Join(val, ",")
			}

			g.appRpc.Info(map[string]interface{}{
				g.requestIDKey: requestID,
				g.hostnameKey:  g.options.Hostname,
				g.privateIPKey: g.options.PrivateIP,
				g.remoteIPKey:  remoteIP,
				g.clientIPKey:  clientIP,
				g.methodKey:    info.FullMethod,
				g.rpcCodeKey:   rpcCode,
				g.errCodeKey:   errCode,
				g.statusKey:    status,
				g.metaKey:      meta,
				g.reqKey:       req,
				g.ctxKey:       ctx.Value(rpcxCtxKey{}),
				g.resKey:       res,
				g.errKey:       err,
				g.errStackKey:  fmt.Sprintf("%+v", err),
				g.resTimeMsKey: time.Now().Sub(ts).Milliseconds(),
			})

			headers := map[string]string{}
			for _, header := range g.options.Headers {
				headers[header] = MetaDataIncomingGet(ctx, strings.ToLower(header))
			}
			_ = grpc.SendHeader(ctx, metadata.New(headers))
			if err != nil {
				err = err.(*Error).SetRequestID(requestID).ToStatus().Err()
			}
		}()

		if err == nil {
			if g.rateLimiter != nil {
				key := info.FullMethod
				if g.options.RateLimiterHeader != "" {
					key = fmt.Sprintf("%s|%s", strings.Join(md.Get(g.options.RateLimiterHeader), ","), info.FullMethod)
				}
				if err = g.rateLimiter.Allow(ctx, key); err != nil {
					if err == micro.ErrFlowControl {
						err = NewError(codes.ResourceExhausted, "ResourceExhausted", err.Error(), err)
					}
				}
			}
		}

		if err == nil {
			if g.parallelCtl != nil {
				key := info.FullMethod
				if g.options.ParallelControllerHeader != "" {
					key = fmt.Sprintf("%s|%s", strings.Join(md.Get(g.options.ParallelControllerHeader), ","), info.FullMethod)
				}
				if err = g.parallelCtl.GetToken(ctx, key); err != nil {
					if err == micro.ErrFlowControl {
						err = NewError(codes.ResourceExhausted, "ResourceExhausted", err.Error(), err)
					} else {
						g.appLog.Warnf("g.parallelCtl.GetToken failed. err: [%+v]", err)
					}
				}
				defer func() {
					if err := g.parallelCtl.PutToken(ctx, key); err != nil {
						g.appLog.Warnf("g.parallelCtl.PutToken failed. err: [%+v]", err)
					}
				}()
			}
		}

		if err == nil {
			for _, h := range g.preHandlers {
				if err = h(ctx, req); err != nil {
					break
				}
			}
		}

		if err == nil {
			for _, validate := range g.validators {
				if err = validate(req); err != nil {
					err = NewError(codes.InvalidArgument, "InvalidArgument", err.Error(), err)
					break
				}
			}
		}

		if err == nil {
			res, err = handler(ctx, req)
		}

		if err != nil {
			switch e := errors.Cause(err).(type) {
			case *Error:
				if e.Detail.Status == 0 {
					e.Detail.Status = int32(runtime.HTTPStatusFromCode(e.Code))
				}
				err = NewError(e.Code, e.Detail.Code, e.Detail.Message, err).SetStatus(int(e.Detail.Status)).SetRequestID(e.Detail.RequestID).SetRefer(e.Detail.Refer)
			default:
				err = NewInternalError(err)
			}
			return res, err
		}

		return res, nil
	})
}

type GrpcInterceptorOptions struct {
	Headers          []string `dft:"X-Request-Id"`
	PrivateIP        string
	Hostname         string
	Validators       []string
	PascalNameKey    bool
	RequestIDMetaKey string `dft:"x-request-id"`
	EnableTrace      bool

	RateLimiterHeader        string
	ParallelControllerHeader string
	RateLimiter              micro.RateLimiterOptions
	ParallelController       micro.ParallelControllerGroupOptions
}

func PrivateIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}
	for _, a := range addrs {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return "unknown"
}

func Hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return ""
	}
	return name
}
