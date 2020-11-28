package rpcx

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/hatlonely/go-kit/refx"
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

func WithMuxProtoErrorHandler(opts ...MuxOption) runtime.ServeMuxOption {
	var options MuxOptions
	_ = refx.SetDefaultValue(&options)
	for _, opt := range opts {
		opt(&options)
	}
	detailMarshal := jsonMarshalErrorDetail
	if options.UseFieldKey {
		detailMarshal = jsonMarshalErrorDetailWithFieldKey
	}

	return runtime.WithProtoErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, res http.ResponseWriter, req *http.Request, err error) {
		res.Header().Set("Content-Type", options.ContentType)
		for _, header := range options.Headers {
			res.Header().Set(header, req.Header.Get(header))
		}

		e := StatusErrorDetail(err, req.Header.Get("X-Request-Id"))
		res.WriteHeader(int(e.Status))
		_, _ = res.Write(detailMarshal(e))
	})
}

func jsonMarshalErrorDetail(detail *ErrorDetail) []byte {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(detail)
	return buf
}

func jsonMarshalErrorDetailWithFieldKey(detail *ErrorDetail) []byte {
	m := map[string]interface{}{}
	if detail.Status != 0 {
		m["Status"] = detail.Status
	}
	if detail.RequestID != "" {
		m["RequestID"] = detail.RequestID
	}
	if detail.Code != "" {
		m["Code"] = detail.Code
	}
	if detail.Message != "" {
		m["Message"] = detail.Message
	}
	if detail.Refer != "" {
		m["Refer"] = detail.Refer
	}
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(m)
	return buf
}

func StatusErrorDetail(err error, requestID string) *ErrorDetail {
	s := status.Convert(err)
	if len(s.Details()) >= 1 {
		if e, ok := s.Details()[0].(*ErrorDetail); ok {
			return e
		}
	}
	return NewInternalError(err).SetRequestID(requestID).Detail
}

type MuxOptions struct {
	Headers     []string `dft:"X-Request-Id"`
	ContentType string   `dft:"application/json"`
	UseFieldKey bool
}

type MuxOption func(options *MuxOptions)

func WithMuxHeaders(headers ...string) MuxOption {
	return func(options *MuxOptions) {
		options.Headers = headers
	}
}

func WithMuxContentType(contentType string) MuxOption {
	return func(options *MuxOptions) {
		options.ContentType = contentType
	}
}

func WithMuxUseFieldKey() MuxOption {
	return func(options *MuxOptions) {
		options.UseFieldKey = true
	}
}
