package wrap

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

type WrapperOptions struct {
	Name         string
	EnableTrace  bool
	EnableMetric bool
	Trace        struct {
		ConstTags map[string]string
	}
	Metric struct {
		Buckets     []float64
		ConstLabels map[string]string
	}
}

type CtxOptions struct {
	DisableTrace           bool
	DisableMetric          bool
	MetricCustomLabelValue string
	TraceTags              map[string]string
	startSpanOpts          []opentracing.StartSpanOption
}

type CtxOption func(options *CtxOptions)

func WithCtxDisableTrace() CtxOption {
	return func(options *CtxOptions) {
		options.DisableTrace = true
	}
}

func WithCtxDisableMetric() CtxOption {
	return func(options *CtxOptions) {
		options.DisableMetric = true
	}
}

func WithCtxMetricCustomLabelValue(val string) CtxOption {
	return func(options *CtxOptions) {
		options.MetricCustomLabelValue = val
	}
}

func WithCtxTraceTag(key string, val string) CtxOption {
	return func(options *CtxOptions) {
		if options.TraceTags == nil {
			options.TraceTags = map[string]string{}
		}
		options.TraceTags[key] = val
	}
}

func WithCtxStartSpanOpts(opts ...opentracing.StartSpanOption) CtxOption {
	return func(options *CtxOptions) {
		options.startSpanOpts = opts
	}
}

type ctxKey struct{}

func NewContext(ctx context.Context, opts ...CtxOption) context.Context {
	options := CtxOptions{
		MetricCustomLabelValue: "default",
	}
	for _, opt := range opts {
		opt(&options)
	}

	return context.WithValue(ctx, &ctxKey{}, &options)
}

func FromContext(ctx context.Context) *CtxOptions {
	val := ctx.Value(&ctxKey{})
	if val == nil {
		return &CtxOptions{
			MetricCustomLabelValue: "default",
		}
	}
	return val.(*CtxOptions)
}
