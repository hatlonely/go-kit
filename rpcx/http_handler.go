package rpcx

import (
	"encoding/json"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	uuid "github.com/satori/go.uuid"
)

type HttpHandlerOptions struct {
	EnableTrace      bool
	RequestIDMetaKey string
}

func NewHttpHandlerWithOptions(options *HttpHandlerOptions) *HttpHandler {
	return &HttpHandler{
		options:        options,
		httpHandlerMap: map[string]http.Handler{},
	}
}

type HttpHandler struct {
	options        *HttpHandlerOptions
	defaultHandler http.Handler
	httpHandlerMap map[string]http.Handler
	preHandlers    []HttpPreHandler
}

type HttpPreHandler func(r *http.Request) error

func (h *HttpHandler) SetDefaultHandler(handler http.Handler) {
	h.defaultHandler = handler
}

func (h *HttpHandler) AddHandler(path string, handler http.Handler) {
	h.httpHandlerMap[path] = handler
}

func (h *HttpHandler) AddPreHandler(handler HttpPreHandler) {
	h.preHandlers = append(h.preHandlers, handler)
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.options.EnableTrace {
		parentSpanContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err == nil || err == opentracing.ErrSpanContextNotFound {
			serverSpan := opentracing.GlobalTracer().StartSpan(
				"ServeHTTP",
				ext.RPCServerOption(parentSpanContext),
				opentracing.Tag{Key: string(ext.Component), Value: "grpc-gateway"},
			)
			r = r.WithContext(opentracing.ContextWithSpan(r.Context(), serverSpan))
			defer serverSpan.Finish()
		}
	}

	var requestID string
	if h.options.RequestIDMetaKey != "" {
		requestID = r.Header.Get(h.options.RequestIDMetaKey)
		if requestID == "" {
			requestID = uuid.NewV4().String()
			r.Header.Set(h.options.RequestIDMetaKey, requestID)
		}
	}

	for _, preHandler := range h.preHandlers {
		if err := preHandler(r); err != nil {
			var httpError *HttpError
			switch e := err.(type) {
			case *HttpError:
				httpError = e
			case *Error:
				httpError = (*HttpError)(e.Detail)
			default:
				httpError = (*HttpError)(NewInternalError(e).Detail)
			}
			httpError.RequestID = requestID
			buf, _ := json.Marshal(httpError)
			_, _ = w.Write(buf)
			return
		}
	}

	if handler, ok := h.httpHandlerMap[r.URL.Path]; ok {
		handler.ServeHTTP(w, r)
	} else {
		h.defaultHandler.ServeHTTP(w, r)
	}
}
