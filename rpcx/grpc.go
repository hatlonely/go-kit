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

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	playgroundValidator "gopkg.in/go-playground/validator.v9"

	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/strx"
	"github.com/hatlonely/go-kit/validator"
)

func MetaDataIncomingGet(ctx context.Context, key string) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get(key), ",")
	}
	return ""
}

func MetaDataIncomingSet(ctx context.Context, key string, val string) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		md.Set(key, val)
	}
}

type grpcCtxKey struct{}

func NewRPCXContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, grpcCtxKey{}, map[string]interface{}{})
}

func CtxSet(ctx context.Context, key string, val interface{}) {
	m := ctx.Value(grpcCtxKey{})
	if m == nil {
		return
	}
	m.(map[string]interface{})[key] = val
}

func CtxGet(ctx context.Context, key string) interface{} {
	m := ctx.Value(grpcCtxKey{})
	if m == nil {
		return nil
	}
	return m.(map[string]interface{})[key]
}

func FromRPCXContext(ctx context.Context) map[string]interface{} {
	m := ctx.Value(grpcCtxKey{})
	if m == nil {
		return nil
	}
	return m.(map[string]interface{})
}

type Logger interface {
	Info(v interface{})
}

func GRPCUnaryInterceptor(log Logger, opts ...GRPCUnaryInterceptorOption) grpc.ServerOption {
	var options GRPCUnaryInterceptorOptions
	refx.SetDefaultValueP(&options)
	for _, opt := range opts {
		opt(&options)
	}
	return GRPCUnaryInterceptorWithOptions(log, &options)
}

func GRPCUnaryInterceptorWithOptions(log Logger, options *GRPCUnaryInterceptorOptions) grpc.ServerOption {
	if options.Hostname == "" {
		options.Hostname = Hostname()
	}
	if options.PrivateIP == "" {
		options.PrivateIP = PrivateIP()
	}
	for _, validate := range options.Validators {
		switch validate {
		case "playground":
			WithGRPCUnaryInterceptorPlaygroundValidator()(options)
		case "default":
			WithGRPCUnaryInterceptorDefaultValidator()(options)
		default:
			panic(fmt.Sprintf("invalid validator [%v], should be one of [playground, default]", validate))
		}
	}
	options.RequestIDMetaKey = strings.ToLower(options.RequestIDMetaKey)

	requestIDKey := "requestID"
	hostnameKey := "hostname"
	privateIPKey := "privateIP"
	remoteIPKey := "remoteIP"
	clientIPKey := "clientIP"
	methodKey := "method"
	rpcCodeKey := "rpcCode"
	errCodeKey := "errCode"
	statusKey := "status"
	metaKey := "meta"
	reqKey := "req"
	ctxKey := "ctx"
	resKey := "res"
	errKey := "err"
	errStackKey := "errStack"
	resTimeMsKey := "resTimeMs"
	if options.PascalNameKey {
		for _, key := range []*string{
			&requestIDKey, &hostnameKey, &privateIPKey, &remoteIPKey, &clientIPKey, &methodKey, &rpcCodeKey, &errCodeKey,
			&statusKey, &metaKey, &reqKey, &ctxKey, &resKey, &errKey, &errStackKey, &resTimeMsKey,
		} {
			*key = strx.PascalName(*key)
		}
	}

	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		span, ctx := opentracing.StartSpanFromContext(ctx, "grpc")
		span.SetTag(methodKey, info.FullMethod)
		defer span.Finish()

		var requestID, remoteIP string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			requestID = strings.Join(md.Get(options.RequestIDMetaKey), ",")
			if requestID == "" {
				requestID = uuid.NewV4().String()
				md.Set(options.RequestIDMetaKey, requestID)
			}
			remoteIP = strings.Split(strings.Join(md.Get("x-remote-addr"), ","), ":")[0]
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

			log.Info(map[string]interface{}{
				requestIDKey: requestID,
				hostnameKey:  options.Hostname,
				privateIPKey: options.PrivateIP,
				remoteIPKey:  remoteIP,
				clientIPKey:  clientIP,
				methodKey:    info.FullMethod,
				rpcCodeKey:   rpcCode,
				errCodeKey:   errCode,
				statusKey:    status,
				metaKey:      meta,
				reqKey:       req,
				ctxKey:       ctx.Value(grpcCtxKey{}),
				resKey:       res,
				errKey:       err,
				errStackKey:  fmt.Sprintf("%+v", err),
				resTimeMsKey: time.Now().Sub(ts).Milliseconds(),
			})

			headers := map[string]string{}
			for _, header := range options.Headers {
				headers[header] = MetaDataIncomingGet(ctx, strings.ToLower(header))
			}
			_ = grpc.SendHeader(ctx, metadata.New(headers))
			if err != nil {
				err = err.(*Error).SetRequestID(requestID).ToStatus().Err()
			}
		}()

		if err == nil {
			for _, h := range options.preHandlers {
				if err = h(ctx, req); err != nil {
					break
				}
			}
		}

		if err == nil {
			for _, validate := range options.validators {
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

type GRPCUnaryInterceptorOptions struct {
	Headers          []string `dft:"X-Request-Id"`
	PrivateIP        string
	Hostname         string
	Validators       []string
	PascalNameKey    bool
	RequestIDMetaKey string `dft:"x-request-id"`

	validators  []func(interface{}) error
	preHandlers []func(ctx context.Context, req interface{}) error
}

type GRPCUnaryInterceptorOption func(options *GRPCUnaryInterceptorOptions)

func WithGRPCUnaryInterceptorPreHandlers(handlers ...func(ctx context.Context, req interface{}) error) GRPCUnaryInterceptorOption {
	return func(options *GRPCUnaryInterceptorOptions) {
		options.preHandlers = append(options.preHandlers, handlers...)
	}
}

func WithGRPCUnaryInterceptorHeaders(headers ...string) GRPCUnaryInterceptorOption {
	return func(options *GRPCUnaryInterceptorOptions) {
		options.Headers = headers
	}
}

func WithGRPCUnaryInterceptorPrivateIP(privateIP string) GRPCUnaryInterceptorOption {
	return func(options *GRPCUnaryInterceptorOptions) {
		options.PrivateIP = privateIP
	}
}

func WithGRPCUnaryInterceptorHostname(hostname string) GRPCUnaryInterceptorOption {
	return func(options *GRPCUnaryInterceptorOptions) {
		options.Hostname = hostname
	}
}

func WithGRPCUnaryInterceptorRequestIDMetaKey(requestIDMetaKey string) GRPCUnaryInterceptorOption {
	return func(options *GRPCUnaryInterceptorOptions) {
		options.RequestIDMetaKey = requestIDMetaKey
	}
}

func WithGRPCUnaryInterceptorPlaygroundValidator() GRPCUnaryInterceptorOption {
	validate := playgroundValidator.New()
	return func(options *GRPCUnaryInterceptorOptions) {
		options.validators = append(options.validators, validate.Struct)
	}
}

func WithGRPCUnaryInterceptorDefaultValidator() GRPCUnaryInterceptorOption {
	return func(options *GRPCUnaryInterceptorOptions) {
		options.validators = append(options.validators, validator.Validate)
	}
}

func WithGRPCUnaryInterceptorValidators(fun ...func(interface{}) error) GRPCUnaryInterceptorOption {
	return func(options *GRPCUnaryInterceptorOptions) {
		options.validators = append(options.validators, fun...)
	}
}

func WithGRPCUnaryInterceptorPascalNameKey() GRPCUnaryInterceptorOption {
	return func(options *GRPCUnaryInterceptorOptions) {
		options.PascalNameKey = true
	}
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
