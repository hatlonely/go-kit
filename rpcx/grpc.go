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

	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/strx"
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

func GRPCUnaryInterceptor(log Logger, opts ...GRPCOption) grpc.ServerOption {
	var options GRPCOptions
	_ = refx.SetDefaultValue(&options)
	for _, opt := range opts {
		opt(&options)
	}
	return GRPCUnaryInterceptorWithOptions(log, &options)
}

func GRPCUnaryInterceptorWithOptions(log Logger, options *GRPCOptions) grpc.ServerOption {
	if options.Hostname == "" {
		options.Hostname = Hostname()
	}
	if options.PrivateIP == "" {
		options.PrivateIP = PrivateIP()
	}
	for _, validate := range options.Validators {
		switch validate {
		case "playground":
			WithPlaygroundValidator()(options)
		case "default":
			WithDefaultValidator()(options)
		default:
			panic(fmt.Sprintf("invalid validator [%v], should be one of [playground, default]", validate))
		}
	}

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

type GRPCOptions struct {
	Headers       []string `dft:"X-Request-Id"`
	PrivateIP     string
	Hostname      string
	Validators    []string
	PascalNameKey bool

	validators  []func(interface{}) error
	preHandlers []func(ctx context.Context, req interface{}) error
}

type GRPCOption func(options *GRPCOptions)

func WithGRPCPreHandlers(handlers ...func(ctx context.Context, req interface{}) error) GRPCOption {
	return func(options *GRPCOptions) {
		options.preHandlers = append(options.preHandlers, handlers...)
	}
}

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
		options.validators = append(options.validators, validate.Struct)
	}
}

func WithDefaultValidator() GRPCOption {
	return func(options *GRPCOptions) {
		options.validators = append(options.validators, validator.Validate)
	}
}

func WithValidators(fun ...func(interface{}) error) GRPCOption {
	return func(options *GRPCOptions) {
		options.validators = append(options.validators, fun...)
	}
}

func WithPascalNameKey() GRPCOption {
	return func(options *GRPCOptions) {
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
