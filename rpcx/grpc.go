package rpcx

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/hatlonely/go-kit/logger"
	"github.com/hatlonely/go-kit/validator"
)

func MetaDataGetRequestID(ctx context.Context) string {
	return MetaDataIncomingGet(ctx, "x-request-id")
}

func MetaDataIncomingGet(ctx context.Context, key string) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals, ok := md[key]; ok {
			return strings.Join(vals, ",")
		}
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
	options := defaultGRPCOptions
	for _, opt := range opts {
		opt(&options)
	}

	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var requestID, remoteIP string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals, ok := md["x-request-id"]; ok {
				requestID = strings.Join(vals, ",")
			}
			if requestID == "" {
				requestID = uuid.NewV4().String()
				md.Set("x-request-id", requestID)
			}
			if vals, ok := md["x-remote-addr"]; ok {
				remoteIP = strings.Split(strings.Join(vals, ","), ":")[0]
			}
		}

		ctx = NewRPCXContext(ctx)

		var res interface{}
		var err error
		ts := time.Now()
		defer func() {
			if perr := recover(); perr != nil {
				err = errors.Wrap(fmt.Errorf("%v\n%v", string(debug.Stack()), perr), "panic")
			}
			p, ok := peer.FromContext(ctx)
			clientIP := ""
			if ok && p != nil {
				clientIP = p.Addr.String()
			}

			log.Info(map[string]interface{}{
				"requestID": requestID,
				"remoteIP":  remoteIP,
				"clientIP":  clientIP,
				"method":    info.FullMethod,
				"req":       req,
				"ctx":       ctx.Value(grpcCtxKey{}),
				"res":       res,
				"err":       err,
				"errStack":  fmt.Sprintf("%+v", err),
				"resTimeMs": time.Now().Sub(ts).Milliseconds(),
			})

			headers := map[string]string{}
			for _, header := range options.Headers {
				headers[header] = MetaDataIncomingGet(ctx, header)
			}
			_ = grpc.SendHeader(ctx, metadata.New(headers))
		}()

		if err = validator.Validate(req); err != nil {
			err = NewErrorWithoutRefer(err, codes.InvalidArgument, requestID, "InvalidArgument", err.Error())
		} else {
			res, err = handler(ctx, req)
		}

		if err != nil {
			switch e := err.(type) {
			case *Error:
				return res, e.ToStatus().Err()
			}
			return res, NewInternalError(err, requestID).ToStatus().Err()
		}
		return res, nil
	})
}

type GRPCOptions struct {
	Headers []string
}

var defaultGRPCOptions = GRPCOptions{
	Headers: []string{"X-Request-Id"},
}

type GRPCOption func(options *GRPCOptions)

func WithGRPCHeaders(headers ...string) GRPCOption {
	return func(options *GRPCOptions) {
		options.Headers = headers
	}
}
