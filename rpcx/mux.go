package rpcx

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func WithMuxMetadata() runtime.ServeMuxOption {
	return runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		requestID := req.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.NewV4().String()
			req.Header.Set("X-Request-Id", requestID)
			return metadata.Pairs("X-Remote-Addr", req.RemoteAddr, "X-Request-Id", requestID)
		}
		return metadata.Pairs("X-Remote-Addr", req.RemoteAddr)
	})
}

func WithMuxIncomingHeaderMatcher() runtime.ServeMuxOption {
	return runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		if strings.HasPrefix(key, "X-") || strings.HasPrefix(key, "x-") {
			return key, true
		}

		return runtime.DefaultHeaderMatcher(key)
	})
}

func WithMuxOutgoingHeaderMatcher() runtime.ServeMuxOption {
	return runtime.WithOutgoingHeaderMatcher(func(key string) (string, bool) {
		return key, true
	})
}

func WithMuxProtoErrorHandler() runtime.ServeMuxOption {
	return runtime.WithProtoErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
		writer.Header().Set("X-Request-Id", request.Header.Get("X-Request-Id"))
		writer.Header().Set("Content-Type", "application/json")

		s := status.Convert(err)
		if len(s.Details()) >= 1 {
			if e, ok := s.Details()[0].(*EInfo); ok {
				e.Status = int64(runtime.HTTPStatusFromCode(codes.Code(e.Status)))
				writer.WriteHeader(int(e.Status))
				buf, _ := json.Marshal(e)
				_, _ = writer.Write(buf)
				return
			}
		}
		e := NewInternalError(err, request.Header.Get("X-Request-Id")).Info
		e.Status = int64(runtime.HTTPStatusFromCode(s.Code()))
		e.Code = http.StatusText(int(e.Status))
		writer.WriteHeader(int(e.Status))
		buf, _ := json.Marshal(e)
		_, _ = writer.Write(buf)
	})
}
