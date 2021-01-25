package wrap

import (
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/pkg/errors"
)

var retryDelayTypeMap = map[string]retry.DelayTypeFunc{
	"BackOff": retry.BackOffDelay,
	"Fixed":   retry.FixedDelay,
	"Random":  retry.RandomDelay,
}

var retryRetryIfMap = map[string]retry.RetryIfFunc{}

type RetryOptions struct {
	Attempts      uint
	Delay         time.Duration
	MaxDelay      time.Duration
	MaxJitter     time.Duration
	LastErrorOnly bool
	DelayType     string `dft:"BackOff"`
	RetryIf       string
}

type Retry struct {
	OnRetry   retry.OnRetryFunc
	RetryIf   retry.RetryIfFunc
	DelayType retry.DelayTypeFunc

	options *RetryOptions
}

func NewRetryWithOptions(options *RetryOptions) (*Retry, error) {
	delayType, err := parseDelayType(options.DelayType)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetry failed")
	}

	r := &Retry{
		options:   options,
		DelayType: delayType,
	}

	if options.RetryIf != "" {
		if retryIf, ok := retryRetryIfMap[options.RetryIf]; ok {
			r.RetryIf = retryIf
		} else {
			return nil, errors.Errorf("unsupported retryIf [%v] func", options.RetryIf)
		}
	}

	return r, nil
}

func parseDelayType(dt string) (retry.DelayTypeFunc, error) {
	if dt == "" {
		return nil, nil
	}

	vs := strings.Split(dt, ",")
	if len(vs) == 1 {
		if delayType, ok := retryDelayTypeMap[dt]; ok {
			return delayType, nil
		}
		return nil, errors.Errorf("unsupported delayType [%v]", dt)
	}

	var funs []retry.DelayTypeFunc
	for _, v := range vs {
		fun, err := parseDelayType(strings.TrimSpace(v))
		if err != nil {
			return nil, err
		}
		funs = append(funs, fun)
	}
	return retry.CombineDelay(funs...), nil
}

func (r *Retry) Do(fun func() error) error {
	var opts []retry.Option

	if r.options.Attempts != 0 {
		opts = append(opts, retry.Attempts(r.options.Attempts))
	}
	if r.options.Delay != 0 {
		opts = append(opts, retry.Delay(r.options.Delay))
	}
	if r.options.MaxDelay != 0 {
		opts = append(opts, retry.MaxDelay(r.options.MaxDelay))
	}
	if r.options.MaxJitter != 0 {
		opts = append(opts, retry.MaxJitter(r.options.MaxJitter))
	}
	if r.options.LastErrorOnly {
		opts = append(opts, retry.LastErrorOnly(true))
	}
	if r.DelayType != nil {
		opts = append(opts, retry.DelayType(r.DelayType))
	}
	if r.OnRetry != nil {
		opts = append(opts, retry.OnRetry(r.OnRetry))
	}
	if r.RetryIf != nil {
		opts = append(opts, retry.RetryIf(r.RetryIf))
	}

	return retry.Do(fun, opts...)
}
