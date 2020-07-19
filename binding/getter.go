package binding

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hatlonely/go-kit/strex"
)

type Getter interface {
	Get(key string) (value interface{}, exists bool)
}

type MapGetter map[string]interface{}

func (g MapGetter) Get(key string) (interface{}, bool) {
	v, ok := g[key]
	return v, ok
}

type GinCtxGetter gin.Context

func (g *GinCtxGetter) Get(key string) (interface{}, bool) {
	if val, ok := (*gin.Context)(g).GetQuery(key); ok {
		return val, true
	}
	if val, ok := (*gin.Context)(g).Get(key); ok {
		return val, true
	}
	return nil, false
}

func NewEnvGetter(prefix string, separator string) *EnvGetter {
	return &EnvGetter{
		prefix:    prefix,
		separator: separator,
	}
}

type EnvGetter struct {
	prefix    string
	separator string
}

func (g *EnvGetter) Get(key string) (interface{}, bool) {
	if g.prefix != "" {
		key = g.prefix + "_" + key
	}
	key = strex.SnakeNameAllCaps(key)
	key = strings.Replace(key, ".", "_", -1)
	key = strings.Replace(key, "[", "_", -1)
	key = strings.Replace(key, "]", "_", -1)
	key = strings.Replace(key, "__", "_", -1)
	str := os.Getenv(key)
	if str == "" {
		return nil, false
	}
	return str, true
}
