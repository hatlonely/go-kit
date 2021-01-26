package wrap

import (
	"context"
)

type WrapperOptions struct {
	EnableTrace  bool
	EnableMetric bool
	Metric       struct {
		Buckets     []float64
		ConstLabels map[string]string
	}
}

type CtxOptions struct {
	DisableTrace           bool
	DisableMetric          bool
	MetricCustomLabelValue string
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

func WithMetricCustomLabelValue(val string) CtxOption {
	return func(options *CtxOptions) {
		options.MetricCustomLabelValue = val
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
