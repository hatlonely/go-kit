package rpcx

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	uuid "github.com/satori/go.uuid"
)

type HttpHandlerOptions struct {
	EnableTrace      bool
	EnableCors       bool
	RequestIDMetaKey string
	Cors             CORSOptions
}

func NewHttpHandlerWithOptions(options *HttpHandlerOptions) *HttpHandler {
	var allowRegex []*regexp.Regexp
	for _, v := range options.Cors.AllowRegex {
		allowRegex = append(allowRegex, regexp.MustCompile(v))
	}
	allowOrigin := map[string]bool{}
	for _, v := range options.Cors.AllowOrigin {
		allowOrigin[v] = true
	}

	return &HttpHandler{
		options:        options,
		httpHandlerMap: map[string]http.Handler{},
		allowRegex:     allowRegex,
		allowOrigin:    allowOrigin,
		allowMethod:    strings.Join(options.Cors.AllowMethod, ","),
		allowHeader:    strings.Join(options.Cors.AllowHeader, ","),
	}
}

type HttpHandler struct {
	options        *HttpHandlerOptions
	defaultHandler http.Handler
	httpHandlerMap map[string]http.Handler
	preHandlers    []HttpPreHandler

	allowRegex  []*regexp.Regexp
	allowOrigin map[string]bool
	allowMethod string
	allowHeader string
}

type HttpPreHandler func(w http.ResponseWriter, r *http.Request) error

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
	// 参考 https://github.com/grpc-ecosystem/grpc-gateway/issues/544
	if h.options.EnableCors {
		if origin := r.Header.Get("Origin"); origin != "" {
			allow := true
		out:
			for {
				if h.options.Cors.AllowAll {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					break
				}
				if h.allowOrigin[origin] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
				for _, re := range h.allowRegex {
					if re.MatchString(origin) {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break out
					}
				}
				allow = false
				break
			}
			if !allow {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Headers", h.allowHeader)
				w.Header().Set("Access-Control-Allow-Methods", h.allowMethod)
				w.WriteHeader(http.StatusOK)
				return
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

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
		if err := preHandler(w, r); err != nil {
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
			w.WriteHeader(int(httpError.Status))
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
