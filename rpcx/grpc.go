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
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	playgroundValidator "gopkg.in/go-playground/validator.v9"

	"github.com/hatlonely/go-kit/logger"
	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/validator"
)

func MetaDataGetRequestID(ctx context.Context) string {
	return MetaDataIncomingGet(ctx, "x-request-id")
}

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
	m := ctx.Value(grpcCtxKey{}).(map[string]interface{})
	m[key] = val
}

func CtxGet(ctx context.Context, key string) interface{} {
	m := ctx.Value(grpcCtxKey{}).(map[string]interface{})
	return m[key]
}

func WithGRPCDecorator(log *logger.Logger, opts ...GRPCOption) grpc.ServerOption {
	var options GRPCOptions
	_ = refx.SetDefaultValue(&options)
	for _, opt := range opts {
		opt(&options)
	}
	return WithGRPCDecoratorWithOptions(log, &options)
}

func WithGRPCDecoratorWithOptions(log *logger.Logger, options *GRPCOptions) grpc.ServerOption {
	if options.Hostname == "" {
		options.Hostname = hostname()
	}
	if options.PrivateIP == "" {
		options.PrivateIP = privateIP()
	}
	if options.validator == nil {
		switch options.Validator {
		case "playground":
			WithPlaygroundValidator()(options)
		case "default":
			WithDefaultValidator()(options)
		case "":
			options.validator = func(i interface{}) error {
				return nil
			}
		default:
			panic(fmt.Sprintf("invalid validator [%v], should be one of [playground, default]", options.Validator))
		}
	}

	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var requestID, remoteIP string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			requestID = strings.Join(md.Get("x-request-id"), ",")
			if requestID == "" {
				requestID = uuid.NewV4().String()
				md.Set("x-request-id", requestID)
			}
			remoteIP = strings.Split(strings.Join(md.Get("x-remote-addr"), ","), ":")[0]
		}

		ctx = NewRPCXContext(ctx)

		var res interface{}
		var err error
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
				rpcCode = e.code.String()
				errCode = e.Detail.Code
				status = int(e.Detail.Status)
			}

			md, _ := metadata.FromIncomingContext(ctx)
			meta := map[string]string{}
			for key, val := range md {
				meta[key] = strings.Join(val, ",")
			}

			log.Info(map[string]interface{}{
				"requestID": requestID,
				"hostname":  options.Hostname,
				"privateIP": options.PrivateIP,
				"remoteIP":  remoteIP,
				"clientIP":  clientIP,
				"method":    info.FullMethod,
				"rpcCode":   rpcCode,
				"errCode":   errCode,
				"status":    status,
				"meta":      meta,
				"req":       req,
				"ctx":       ctx.Value(grpcCtxKey{}),
				"res":       res,
				"err":       err,
				"errStack":  fmt.Sprintf("%+v", err),
				"resTimeMs": time.Now().Sub(ts).Milliseconds(),
			})

			headers := map[string]string{}
			for _, header := range options.Headers {
				headers[header] = MetaDataIncomingGet(ctx, strings.ToLower(header))
			}
			_ = grpc.SendHeader(ctx, metadata.New(headers))
		}()

		if err = options.validator(req); err != nil {
			err = NewError(codes.InvalidArgument, "InvalidArgument", err.Error(), err)
		} else {
			res, err = handler(ctx, req)
		}

		if err != nil {
			switch e := errors.Cause(err).(type) {
			case *Error:
				if e.Detail.Status == 0 {
					e.Detail.Status = int32(runtime.HTTPStatusFromCode(e.code))
				}
				err = NewError(e.code, e.Detail.Code, e.Detail.Message, err).SetStatus(int(e.Detail.Status)).SetRequestID(e.Detail.RequestID).SetRefer(e.Detail.Refer)
			default:
				err = NewInternalError(err)
			}
			return res, err.(*Error).SetRequestID(requestID).ToStatus().Err()
		}
		return res, nil
	})
}

type GRPCOptions struct {
	Headers   []string `dft:"X-Request-Id"`
	PrivateIP string
	Hostname  string
	Validator string

	validator func(interface{}) error
}

type GRPCOption func(options *GRPCOptions)

func WithGRPCHeaders(headers ...string) GRPCOption {
	return func(options *GRPCOptions) {
		options.Headers = headers
	}
}

func WithPrivateIP(privateIP string) GRPCOption {
	return func(options *GRPCOptions) {
		options.PrivateIP = privateIP
	}
}

func WithHostname(hostname string) GRPCOption {
	return func(options *GRPCOptions) {
		options.Hostname = hostname
	}
}

func WithPlaygroundValidator() GRPCOption {
	validate := playgroundValidator.New()
	return func(options *GRPCOptions) {
		options.validator = validate.Struct
	}
}

func WithDefaultValidator() GRPCOption {
	return func(options *GRPCOptions) {
		options.validator = validator.Validate
	}
}

func WithValidator(fun func(interface{}) error) GRPCOption {
	return func(options *GRPCOptions) {
		options.validator = fun
	}
}

func privateIP() string {
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

func hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return ""
	}
	return name
}
