package binging

import (
	"github.com/gin-gonic/gin"
)

type Getter interface {
	Get(key string) (value interface{}, exists bool)
}

type MapGetter map[string]interface{}

func (g MapGetter) Get(key string) (interface{}, bool) {
	v, ok := g[key]
	return v, ok
}

type GinGetter gin.Context

func (g *GinGetter) Get(key string) (interface{}, bool) {
	if val, ok := (*gin.Context)(g).GetQuery(key); ok {
		return val, true
	}
	if val, ok := (*gin.Context)(g).Get(key); ok {
		return val, true
	}
	return nil, false
}
