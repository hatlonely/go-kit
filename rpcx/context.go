package rpcx

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

func MetaDataIncomingGet(ctx context.Context, key string) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get(key), ",")
	}
	return ""
}

func MetaDataIncomingSet(ctx context.Context, key string, val string) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		md.Set(key, val)
	}
}

type rpcxCtxKey struct{}

func NewRPCXContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, rpcxCtxKey{}, map[string]interface{}{})
}

func CtxSet(ctx context.Context, key string, val interface{}) {
	m := ctx.Value(rpcxCtxKey{})
	if m == nil {
		return
	}
	m.(map[string]interface{})[key] = val
}

func CtxGet(ctx context.Context, key string) interface{} {
	m := ctx.Value(rpcxCtxKey{})
	if m == nil {
		return nil
	}
	return m.(map[string]interface{})[key]
}

func FromRPCXContext(ctx context.Context) map[string]interface{} {
	m := ctx.Value(rpcxCtxKey{})
	if m == nil {
		return nil
	}
	return m.(map[string]interface{})
}
