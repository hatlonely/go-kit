package rpcx

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TraceWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		h.ServeHTTP(w, r)
	})
}

func MetricWrapper(h http.Handler) http.Handler {
	promHandler := promhttp.Handler()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			promHandler.ServeHTTP(w, r)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func MapHandlerWrapper(h http.Handler, handlerMap map[string]http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handler, ok := handlerMap[r.URL.Path]; ok {
			handler.ServeHTTP(w, r)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
