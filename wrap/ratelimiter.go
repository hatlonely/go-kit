package wrap

import (
	"context"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"

	"github.com/hatlonely/go-kit/refx"
)

//func MustConstructorWithOption

func RegisterRateLimiterGroup(key string, fun interface{}) {
	rt := reflect.TypeOf(fun)
	if rt.Kind() != reflect.Func {
		panic("fun should be a function type")
	}
	if rt.NumIn() > 1 {
		panic("fun parameters number should not greater than 1")
	}
	if rt.NumOut() > 2 {
		panic("fun return number should not greater than 2")
	}
	if rt.NumOut() == 0 || !rt.Out(0).Implements(reflect.TypeOf((*RateLimiterGroup)(nil)).Elem()) {
		panic("fun should return a RateLimiterGroup")
	}
	if rt.NumOut() == 2 && rt.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		panic("fun should return an error")
	}

	rateLimiterGroupMap[key] = &RateLimiterGroupConstructor{
		HasParam:    rt.NumIn() == 1,
		ReturnError: rt.NumOut() == 2,
		ParamType:   rt.In(0),
		FuncValue:   reflect.ValueOf(fun),
	}
}

type RateLimiterGroupConstructor struct {
	HasParam    bool
	ReturnError bool
	ParamType   reflect.Type
	FuncValue   reflect.Value
}

var rateLimiterGroupMap = map[string]*RateLimiterGroupConstructor{}

type RateLimiterGroup interface {
	Wait(ctx context.Context, key string) error
}

func NewRateLimiterGroupWithOptions(options *RateLimiterGroupOptions, opts ...refx.Option) (RateLimiterGroup, error) {
	constructor, ok := rateLimiterGroupMap[options.Type]
	if !ok {
		return nil, errors.Errorf("unsupported rate limiter type: [%v]", options.Type)
	}

	var result []reflect.Value
	if constructor.HasParam {
		params := reflect.New(reflect.TypeOf(&LocalGroupRateLimiterOptions{}))
		if err := refx.InterfaceToStruct(options.V, params.Interface(), opts...); err != nil {
			return nil, errors.Wrap(err, "refx.InterfaceToStruct failed")
		}
		result = constructor.FuncValue.Call([]reflect.Value{params.Elem()})
	} else {
		result = constructor.FuncValue.Call(nil)
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, result[1].Interface().(error)
		}
		return result[0].Interface().(RateLimiterGroup), nil
	}

	return result[0].Interface().(RateLimiterGroup), nil
}

type RateLimiterGroupOptions struct {
	Type                  string
	LocalGroupRateLimiter LocalGroupRateLimiterOptions
	LocalShareRateLimiter LocalShareRateLimiterOptions
	V                     interface{}
}

func init() {
	RegisterRateLimiterGroup("LocalGroup", NewLocalGroupRateLimiterWithOptions)
	RegisterRateLimiterGroup("ShareGroup", NewLocalShareRateLimiterWithOptions)
}

type LocalGroupRateLimiter struct {
	limiters map[string]*rate.Limiter
}

type LocalGroupRateLimiterOptions map[string]struct {
	Interval time.Duration `rule:"x > 0"`
	Burst    int           `rule:"x > 0"`
}

func NewLocalGroupRateLimiterWithOptions(options *LocalGroupRateLimiterOptions) (*LocalGroupRateLimiter, error) {
	limiters := map[string]*rate.Limiter{}

	for key, val := range *options {
		if val.Interval <= 0 {
			return nil, errors.Errorf("interval should be positive")
		}
		if val.Burst <= 0 {
			return nil, errors.Errorf("bucket should be positive")
		}
		limiters[key] = rate.NewLimiter(rate.Every(val.Interval), val.Burst)
	}

	return &LocalGroupRateLimiter{limiters: limiters}, nil
}

func (r *LocalGroupRateLimiter) Wait(ctx context.Context, key string) error {
	limiter, ok := r.limiters[key]
	if !ok {
		return nil
	}

	return limiter.Wait(ctx)
}

type LocalShareRateLimiterOptions struct {
	Interval time.Duration `rule:"x > 0"`
	Burst    int           `rule:"x > 0"`
	Cost     map[string]int
}

type LocalShareRateLimiter struct {
	limiter *rate.Limiter

	cost map[string]int
}

func NewLocalShareRateLimiterWithOptions(options *LocalShareRateLimiterOptions) (*LocalShareRateLimiter, error) {
	limiter := rate.NewLimiter(rate.Every(options.Interval), options.Burst)

	for _, cost := range options.Cost {
		if cost <= 0 {
			return nil, errors.Errorf("cost should be positive")
		}
	}

	return &LocalShareRateLimiter{
		limiter: limiter,
		cost:    options.Cost,
	}, nil
}

func (r *LocalShareRateLimiter) Wait(ctx context.Context, key string) error {
	cost, ok := r.cost[key]
	if !ok {
		return nil
	}

	return r.limiter.WaitN(ctx, cost)
}
