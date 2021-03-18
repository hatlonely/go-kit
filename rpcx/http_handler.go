package rpcx

import (
	"bytes"
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

type CORSOptions struct {
	AllowAll    bool
	AllowRegex  []string
	AllowOrigin []string
	AllowMethod []string `dft:"GET,HEAD,POST,PUT,DELETE"`
	AllowHeader []string `dft:"Content-Type,Accept"`
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
	postHandlers   []HttpPostHandler

	allowRegex  []*regexp.Regexp
	allowOrigin map[string]bool
	allowMethod string
	allowHeader string
}

type HttpPreHandler func(w http.ResponseWriter, r *http.Request) error
type HttpPostHandler func(w *BufferedHttpResponseWriter, r *http.Request) (*BufferedHttpResponseWriter, error)

func (h *HttpHandler) SetDefaultHandler(handler http.Handler) {
	h.defaultHandler = handler
}

func (h *HttpHandler) AddHandler(path string, handler http.Handler) {
	h.httpHandlerMap[path] = handler
}

func (h *HttpHandler) AddPreHandler(handler HttpPreHandler) {
	h.preHandlers = append(h.preHandlers, handler)
}

func (h *HttpHandler) AddPostHandler(handler HttpPostHandler) {
	h.postHandlers = append(h.postHandlers, handler)
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bw := NewBufferedHttpResponseWriter()

	defer func() {
		for key, vals := range bw.Header() {
			for _, val := range vals {
				w.Header().Set(key, val)
			}
		}
		w.WriteHeader(bw.Status())
		_, _ = w.Write(bw.Body())
	}()

	// 参考 https://github.com/grpc-ecosystem/grpc-gateway/issues/544
	if h.options.EnableCors {
		if origin := r.Header.Get("Origin"); origin != "" {
			allow := true
		out:
			for {
				if h.options.Cors.AllowAll {
					bw.Header().Set("Access-Control-Allow-Origin", "*")
					break
				}
				if h.allowOrigin[origin] {
					bw.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
				for _, re := range h.allowRegex {
					if re.MatchString(origin) {
						bw.Header().Set("Access-Control-Allow-Origin", origin)
						break out
					}
				}
				allow = false
				break
			}
			if !allow {
				bw.WriteHeader(http.StatusForbidden)
				return
			}

			if r.Method == "OPTIONS" {
				bw.Header().Set("Access-Control-Allow-Headers", h.allowHeader)
				bw.Header().Set("Access-Control-Allow-Methods", h.allowMethod)
				bw.WriteHeader(http.StatusOK)
				return
			}
		} else {
			bw.WriteHeader(http.StatusForbidden)
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
		if err := preHandler(bw, r); err != nil {
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
			bw.WriteHeader(int(httpError.Status))
			_, _ = bw.Write(buf)
			return
		}
	}

	if handler, ok := h.httpHandlerMap[r.URL.Path]; ok {
		handler.ServeHTTP(bw, r)
	} else {
		h.defaultHandler.ServeHTTP(bw, r)
	}

	for _, postHandler := range h.postHandlers {
		var err error
		bw, err = postHandler(bw, r)
		if err != nil {
			return
		}
	}
}

type BufferedHttpResponseWriter struct {
	status int
	body   bytes.Buffer
	header http.Header
}

func NewBufferedHttpResponseWriter() *BufferedHttpResponseWriter {
	return &BufferedHttpResponseWriter{
		status: http.StatusOK,
		header: map[string][]string{},
	}
}

func (w *BufferedHttpResponseWriter) Header() http.Header {
	return w.header
}

func (w *BufferedHttpResponseWriter) Write(buf []byte) (int, error) {
	return w.body.Write(buf)
}

func (w *BufferedHttpResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

func (w *BufferedHttpResponseWriter) Status() int {
	return w.status
}

func (w *BufferedHttpResponseWriter) Body() []byte {
	return w.body.Bytes()
}

func (w *BufferedHttpResponseWriter) SetHeader(header http.Header) {
	w.header = header
}
