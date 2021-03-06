package bind

import (
	"os"
	"strings"

	"github.com/hatlonely/go-kit/strx"
)

func NewEnvGetter(opts ...EnvGetterOption) *EnvGetter {
	var options = defaultEnvGetterOptions
	for _, opt := range opts {
		opt(options)
	}
	return &EnvGetter{
		prefix:    options.Prefix,
		separator: options.Separator,
	}
}

type EnvGetterOptions struct {
	Prefix    string
	Separator string
}

type EnvGetterOption func(opt *EnvGetterOptions)

var defaultEnvGetterOptions = &EnvGetterOptions{
	Prefix:    "",
	Separator: "_",
}

func WithEnvPrefix(prefix string) EnvGetterOption {
	return func(opt *EnvGetterOptions) {
		opt.Prefix = prefix
	}
}

func WithEnvSeparator(separator string) EnvGetterOption {
	return func(opt *EnvGetterOptions) {
		opt.Separator = separator
	}
}

type EnvGetter struct {
	prefix    string
	separator string
}

func (g *EnvGetter) Get(key string) (interface{}, bool) {
	val := os.Getenv(g.TransformKey(key))
	if val == "" {
		return nil, false
	}
	return val, true
}

func (g *EnvGetter) TransformKey(key string) string {
	if g.prefix != "" {
		key = g.prefix + g.separator + key
	}
	key = strx.SnakeNameAllCaps(key)
	key = strings.Replace(key, ".", g.separator, -1)
	key = strings.Replace(key, "[", g.separator, -1)
	key = strings.Replace(key, "]", g.separator, -1)
	key = strings.Replace(key, g.separator+g.separator, g.separator, -1)
	return key
}
