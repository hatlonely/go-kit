package rpcx

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/hatlonely/go-kit/refx"
)

type CORSOption func(options *CORSOptions)

func WithCORSAllowAll() CORSOption {
	return func(options *CORSOptions) {
		options.AllowAll = true
	}
}

func WithAllowRegex(allowRegex ...string) CORSOption {
	return func(options *CORSOptions) {
		options.AllowRegex = allowRegex
	}
}

func WithAllowOrigin(allowOrigin ...string) CORSOption {
	return func(options *CORSOptions) {
		options.AllowOrigin = allowOrigin
	}
}

func WithAllowMethod(allowMethod ...string) CORSOption {
	return func(options *CORSOptions) {
		options.AllowMethod = allowMethod
	}
}

func WithAllowHeader(allowHeader ...string) CORSOption {
	return func(options *CORSOptions) {
		options.AllowHeader = allowHeader
	}
}

func CORS(h http.Handler, opts ...CORSOption) http.Handler {
	var options = CORSOptions{}
	_ = refx.SetDefaultValue(&options)

	for _, opt := range opts {
		opt(&options)
	}

	return CORSWithOptions(h, &options)
}

// 参考 https://github.com/grpc-ecosystem/grpc-gateway/issues/544
func CORSWithOptions(h http.Handler, options *CORSOptions) http.Handler {
	var allowRegex []*regexp.Regexp
	for _, v := range options.AllowRegex {
		allowRegex = append(allowRegex, regexp.MustCompile(v))
	}
	allowOrigin := map[string]bool{}
	for _, v := range options.AllowOrigin {
		allowOrigin[v] = true
	}
	allowMethod := strings.Join(options.AllowMethod, ",")
	allowHeader := strings.Join(options.AllowHeader, ",")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
		out:
			for {
				if options.AllowAll {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					break
				}
				if allowOrigin[origin] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
				for _, re := range allowRegex {
					if re.MatchString(origin) {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break out
					}
				}
			}

			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Headers", allowHeader)
				w.Header().Set("Access-Control-Allow-Methods", allowMethod)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
