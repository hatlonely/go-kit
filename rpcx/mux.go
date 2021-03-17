package rpcx

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

type MuxInterceptorOptions struct {
	Headers                 []string `dft:"X-Request-Id"`
	RequestIDMetaKey        string   `dft:"X-Request-Id"`
	UsePascalNameErrKey     bool
	MarshalUseProtoNames    bool
	MarshalEmitUnpopulated  bool
	UnmarshalDiscardUnknown bool
}

func NewMuxInterceptorWithOptions(options *MuxInterceptorOptions) (*MuxInterceptor, error) {
	m := &MuxInterceptor{
		options:              options,
		errorDetailMarshaler: JsonMarshalErrorDetail,
	}
	if options.UsePascalNameErrKey {
		m.errorDetailMarshaler = JsonMarshalErrorDetailWithFieldKey
	}

	return m, nil
}

type MuxInterceptor struct {
	options *MuxInterceptorOptions

	errorDetailMarshaler ErrorDetailMarshaler
}

type ErrorDetailMarshaler func(detail *ErrorDetail) []byte

func (m *MuxInterceptor) SetErrorDetailMarshaler(errorDetailMarshaler ErrorDetailMarshaler) {
	m.errorDetailMarshaler = errorDetailMarshaler
}

func (m *MuxInterceptor) ServeMuxOptions() []runtime.ServeMuxOption {
	var opts []runtime.ServeMuxOption
	opts = append(opts, m.MuxMetaData())
	opts = append(opts, m.MuxIncomingHeaderMatcher())
	opts = append(opts, m.MuxOutgoingHeaderMatcher())
	opts = append(opts, m.MuxProtoErrorHandler())
	opts = append(opts, m.MuxMarshalerOption())
	return opts
}

func (m *MuxInterceptor) MuxMarshalerOption() runtime.ServeMuxOption {
	return runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   m.options.MarshalUseProtoNames,
			EmitUnpopulated: m.options.MarshalEmitUnpopulated,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: m.options.UnmarshalDiscardUnknown,
		},
	})
}

func (m *MuxInterceptor) MuxMetaData() runtime.ServeMuxOption {
	return runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		requestID := req.Header.Get(m.options.RequestIDMetaKey)
		if requestID == "" {
			requestID = uuid.NewV4().String()
			req.Header.Set(m.options.RequestIDMetaKey, requestID)
			return metadata.Pairs("X-Remote-Addr", req.RemoteAddr, m.options.RequestIDMetaKey, requestID)
		}
		return metadata.Pairs("X-Remote-Addr", req.RemoteAddr)
	})
}

func (m *MuxInterceptor) MuxIncomingHeaderMatcher() runtime.ServeMuxOption {
	return runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		if strings.HasPrefix(key, "X-") || strings.HasPrefix(key, "x-") {
			return key, true
		}
		return runtime.DefaultHeaderMatcher(key)
	})
}

func (m *MuxInterceptor) MuxOutgoingHeaderMatcher() runtime.ServeMuxOption {
	return runtime.WithOutgoingHeaderMatcher(func(key string) (string, bool) {
		return key, true
	})
}

func (m *MuxInterceptor) MuxProtoErrorHandler() runtime.ServeMuxOption {
	return runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, res http.ResponseWriter, req *http.Request, err error) {
		res.Header().Set("Content-Type", "application/json")
		for _, header := range m.options.Headers {
			res.Header().Set(header, req.Header.Get(header))
		}

		e := StatusErrorDetail(err, req.Header.Get(m.options.RequestIDMetaKey))
		res.WriteHeader(int(e.Status))
		_, _ = res.Write(m.errorDetailMarshaler(e))
	})
}

func JsonMarshalErrorDetail(detail *ErrorDetail) []byte {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(detail)
	return buf
}

func JsonMarshalErrorDetailWithFieldKey(detail *ErrorDetail) []byte {
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
	} else {
		// 未走到 grpc handler 中的错误处理
		return NewErrorf(nil, s.Code(), http.StatusText(runtime.HTTPStatusFromCode(s.Code())), s.Message()).SetRequestID(requestID).Detail
	}
	return NewInternalError(err).SetRequestID(requestID).Detail
}
