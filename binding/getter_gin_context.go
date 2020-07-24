package binding

import (
	"github.com/gin-gonic/gin"
)

func NewGinCtxGetter(ctx *gin.Context) *GinCtxGetter {
	return (*GinCtxGetter)(ctx)
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
