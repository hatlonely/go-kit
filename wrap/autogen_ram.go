// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

func NewRAMClientWrapper(
	obj *ram.Client,
	retry *micro.Retry,
	options *WrapperOptions,
	durationMetric *prometheus.HistogramVec,
	inflightMetric *prometheus.GaugeVec,
	rateLimiter micro.RateLimiter,
	parallelController micro.ParallelController) *RAMClientWrapper {
	return &RAMClientWrapper{
		obj:                obj,
		retry:              retry,
		options:            options,
		durationMetric:     durationMetric,
		inflightMetric:     inflightMetric,
		rateLimiter:        rateLimiter,
		parallelController: parallelController,
	}
}

type RAMClientWrapper struct {
	obj                *ram.Client
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *RAMClientWrapper) Unwrap() *ram.Client {
	return w.obj
}

func (w *RAMClientWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *RAMClientWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options micro.RetryOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		retry, err := micro.NewRetryWithOptions(&options)
		if err != nil {
			return errors.Wrap(err, "NewRetryWithOptions failed")
		}
		w.retry = retry
		return nil
	}
}

func (w *RAMClientWrapper) OnRateLimiterChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options micro.RateLimiterOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		rateLimiter, err := micro.NewRateLimiterWithOptions(&options, opts...)
		if err != nil {
			return errors.Wrap(err, "NewRateLimiterWithOptions failed")
		}
		w.rateLimiter = rateLimiter
		return nil
	}
}

func (w *RAMClientWrapper) OnParallelControllerChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options micro.ParallelControllerOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		parallelController, err := micro.NewParallelControllerWithOptions(&options, opts...)
		if err != nil {
			return errors.Wrap(err, "NewParallelControllerWithOptions failed")
		}
		w.parallelController = parallelController
		return nil
	}
}

func (w *RAMClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        fmt.Sprintf("%s_ram_Client_durationMs", options.Name),
		Help:        "ram Client response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        fmt.Sprintf("%s_ram_Client_inflight", options.Name),
		Help:        "ram Client inflight",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "custom"})
}

func (w *RAMClientWrapper) AddUserToGroup(ctx context.Context, request *ram.AddUserToGroupRequest) (*ram.AddUserToGroupResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.AddUserToGroupResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AddUserToGroup", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AddUserToGroup", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AddUserToGroup", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.AddUserToGroup")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.AddUserToGroup", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.AddUserToGroup", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.AddUserToGroup", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AddUserToGroup(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) AddUserToGroupWithCallback(request *ram.AddUserToGroupRequest, callback func(response *ram.AddUserToGroupResponse, err error)) <-chan int {
	res0 := w.obj.AddUserToGroupWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) AddUserToGroupWithChan(request *ram.AddUserToGroupRequest) (<-chan *ram.AddUserToGroupResponse, <-chan error) {
	res0, res1 := w.obj.AddUserToGroupWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) AttachPolicyToGroup(ctx context.Context, request *ram.AttachPolicyToGroupRequest) (*ram.AttachPolicyToGroupResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.AttachPolicyToGroupResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AttachPolicyToGroup", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AttachPolicyToGroup", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AttachPolicyToGroup", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.AttachPolicyToGroup")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.AttachPolicyToGroup", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.AttachPolicyToGroup", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.AttachPolicyToGroup", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AttachPolicyToGroup(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) AttachPolicyToGroupWithCallback(request *ram.AttachPolicyToGroupRequest, callback func(response *ram.AttachPolicyToGroupResponse, err error)) <-chan int {
	res0 := w.obj.AttachPolicyToGroupWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) AttachPolicyToGroupWithChan(request *ram.AttachPolicyToGroupRequest) (<-chan *ram.AttachPolicyToGroupResponse, <-chan error) {
	res0, res1 := w.obj.AttachPolicyToGroupWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) AttachPolicyToRole(ctx context.Context, request *ram.AttachPolicyToRoleRequest) (*ram.AttachPolicyToRoleResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.AttachPolicyToRoleResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AttachPolicyToRole", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AttachPolicyToRole", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AttachPolicyToRole", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.AttachPolicyToRole")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.AttachPolicyToRole", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.AttachPolicyToRole", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.AttachPolicyToRole", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AttachPolicyToRole(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) AttachPolicyToRoleWithCallback(request *ram.AttachPolicyToRoleRequest, callback func(response *ram.AttachPolicyToRoleResponse, err error)) <-chan int {
	res0 := w.obj.AttachPolicyToRoleWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) AttachPolicyToRoleWithChan(request *ram.AttachPolicyToRoleRequest) (<-chan *ram.AttachPolicyToRoleResponse, <-chan error) {
	res0, res1 := w.obj.AttachPolicyToRoleWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) AttachPolicyToUser(ctx context.Context, request *ram.AttachPolicyToUserRequest) (*ram.AttachPolicyToUserResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.AttachPolicyToUserResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.AttachPolicyToUser", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.AttachPolicyToUser", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.AttachPolicyToUser", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.AttachPolicyToUser")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.AttachPolicyToUser", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.AttachPolicyToUser", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.AttachPolicyToUser", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.AttachPolicyToUser(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) AttachPolicyToUserWithCallback(request *ram.AttachPolicyToUserRequest, callback func(response *ram.AttachPolicyToUserResponse, err error)) <-chan int {
	res0 := w.obj.AttachPolicyToUserWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) AttachPolicyToUserWithChan(request *ram.AttachPolicyToUserRequest) (<-chan *ram.AttachPolicyToUserResponse, <-chan error) {
	res0, res1 := w.obj.AttachPolicyToUserWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) BindMFADevice(ctx context.Context, request *ram.BindMFADeviceRequest) (*ram.BindMFADeviceResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.BindMFADeviceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.BindMFADevice", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.BindMFADevice", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.BindMFADevice", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.BindMFADevice")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.BindMFADevice", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.BindMFADevice", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.BindMFADevice", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.BindMFADevice(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) BindMFADeviceWithCallback(request *ram.BindMFADeviceRequest, callback func(response *ram.BindMFADeviceResponse, err error)) <-chan int {
	res0 := w.obj.BindMFADeviceWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) BindMFADeviceWithChan(request *ram.BindMFADeviceRequest) (<-chan *ram.BindMFADeviceResponse, <-chan error) {
	res0, res1 := w.obj.BindMFADeviceWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ChangePassword(ctx context.Context, request *ram.ChangePasswordRequest) (*ram.ChangePasswordResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ChangePasswordResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ChangePassword", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ChangePassword", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ChangePassword", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ChangePassword")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ChangePassword", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ChangePassword", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ChangePassword", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ChangePassword(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ChangePasswordWithCallback(request *ram.ChangePasswordRequest, callback func(response *ram.ChangePasswordResponse, err error)) <-chan int {
	res0 := w.obj.ChangePasswordWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ChangePasswordWithChan(request *ram.ChangePasswordRequest) (<-chan *ram.ChangePasswordResponse, <-chan error) {
	res0, res1 := w.obj.ChangePasswordWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ClearAccountAlias(ctx context.Context, request *ram.ClearAccountAliasRequest) (*ram.ClearAccountAliasResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ClearAccountAliasResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ClearAccountAlias", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ClearAccountAlias", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ClearAccountAlias", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ClearAccountAlias")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ClearAccountAlias", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ClearAccountAlias", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ClearAccountAlias", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ClearAccountAlias(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ClearAccountAliasWithCallback(request *ram.ClearAccountAliasRequest, callback func(response *ram.ClearAccountAliasResponse, err error)) <-chan int {
	res0 := w.obj.ClearAccountAliasWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ClearAccountAliasWithChan(request *ram.ClearAccountAliasRequest) (<-chan *ram.ClearAccountAliasResponse, <-chan error) {
	res0, res1 := w.obj.ClearAccountAliasWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) CreateAccessKey(ctx context.Context, request *ram.CreateAccessKeyRequest) (*ram.CreateAccessKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.CreateAccessKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateAccessKey", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateAccessKey", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateAccessKey", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.CreateAccessKey")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.CreateAccessKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.CreateAccessKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.CreateAccessKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateAccessKey(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) CreateAccessKeyWithCallback(request *ram.CreateAccessKeyRequest, callback func(response *ram.CreateAccessKeyResponse, err error)) <-chan int {
	res0 := w.obj.CreateAccessKeyWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) CreateAccessKeyWithChan(request *ram.CreateAccessKeyRequest) (<-chan *ram.CreateAccessKeyResponse, <-chan error) {
	res0, res1 := w.obj.CreateAccessKeyWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) CreateGroup(ctx context.Context, request *ram.CreateGroupRequest) (*ram.CreateGroupResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.CreateGroupResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateGroup", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateGroup", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateGroup", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.CreateGroup")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.CreateGroup", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.CreateGroup", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.CreateGroup", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateGroup(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) CreateGroupWithCallback(request *ram.CreateGroupRequest, callback func(response *ram.CreateGroupResponse, err error)) <-chan int {
	res0 := w.obj.CreateGroupWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) CreateGroupWithChan(request *ram.CreateGroupRequest) (<-chan *ram.CreateGroupResponse, <-chan error) {
	res0, res1 := w.obj.CreateGroupWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) CreateLoginProfile(ctx context.Context, request *ram.CreateLoginProfileRequest) (*ram.CreateLoginProfileResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.CreateLoginProfileResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateLoginProfile", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateLoginProfile", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateLoginProfile", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.CreateLoginProfile")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.CreateLoginProfile", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.CreateLoginProfile", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.CreateLoginProfile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateLoginProfile(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) CreateLoginProfileWithCallback(request *ram.CreateLoginProfileRequest, callback func(response *ram.CreateLoginProfileResponse, err error)) <-chan int {
	res0 := w.obj.CreateLoginProfileWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) CreateLoginProfileWithChan(request *ram.CreateLoginProfileRequest) (<-chan *ram.CreateLoginProfileResponse, <-chan error) {
	res0, res1 := w.obj.CreateLoginProfileWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) CreatePolicy(ctx context.Context, request *ram.CreatePolicyRequest) (*ram.CreatePolicyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.CreatePolicyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreatePolicy", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreatePolicy", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreatePolicy", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.CreatePolicy")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.CreatePolicy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.CreatePolicy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.CreatePolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreatePolicy(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) CreatePolicyVersion(ctx context.Context, request *ram.CreatePolicyVersionRequest) (*ram.CreatePolicyVersionResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.CreatePolicyVersionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreatePolicyVersion", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreatePolicyVersion", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreatePolicyVersion", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.CreatePolicyVersion")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.CreatePolicyVersion", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.CreatePolicyVersion", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.CreatePolicyVersion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreatePolicyVersion(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) CreatePolicyVersionWithCallback(request *ram.CreatePolicyVersionRequest, callback func(response *ram.CreatePolicyVersionResponse, err error)) <-chan int {
	res0 := w.obj.CreatePolicyVersionWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) CreatePolicyVersionWithChan(request *ram.CreatePolicyVersionRequest) (<-chan *ram.CreatePolicyVersionResponse, <-chan error) {
	res0, res1 := w.obj.CreatePolicyVersionWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) CreatePolicyWithCallback(request *ram.CreatePolicyRequest, callback func(response *ram.CreatePolicyResponse, err error)) <-chan int {
	res0 := w.obj.CreatePolicyWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) CreatePolicyWithChan(request *ram.CreatePolicyRequest) (<-chan *ram.CreatePolicyResponse, <-chan error) {
	res0, res1 := w.obj.CreatePolicyWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) CreateRole(ctx context.Context, request *ram.CreateRoleRequest) (*ram.CreateRoleResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.CreateRoleResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateRole", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateRole", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateRole", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.CreateRole")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.CreateRole", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.CreateRole", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.CreateRole", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateRole(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) CreateRoleWithCallback(request *ram.CreateRoleRequest, callback func(response *ram.CreateRoleResponse, err error)) <-chan int {
	res0 := w.obj.CreateRoleWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) CreateRoleWithChan(request *ram.CreateRoleRequest) (<-chan *ram.CreateRoleResponse, <-chan error) {
	res0, res1 := w.obj.CreateRoleWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) CreateUser(ctx context.Context, request *ram.CreateUserRequest) (*ram.CreateUserResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.CreateUserResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateUser", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateUser", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateUser", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.CreateUser")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.CreateUser", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.CreateUser", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.CreateUser", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateUser(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) CreateUserWithCallback(request *ram.CreateUserRequest, callback func(response *ram.CreateUserResponse, err error)) <-chan int {
	res0 := w.obj.CreateUserWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) CreateUserWithChan(request *ram.CreateUserRequest) (<-chan *ram.CreateUserResponse, <-chan error) {
	res0, res1 := w.obj.CreateUserWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) CreateVirtualMFADevice(ctx context.Context, request *ram.CreateVirtualMFADeviceRequest) (*ram.CreateVirtualMFADeviceResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.CreateVirtualMFADeviceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.CreateVirtualMFADevice", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.CreateVirtualMFADevice", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.CreateVirtualMFADevice", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.CreateVirtualMFADevice")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.CreateVirtualMFADevice", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.CreateVirtualMFADevice", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.CreateVirtualMFADevice", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.CreateVirtualMFADevice(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) CreateVirtualMFADeviceWithCallback(request *ram.CreateVirtualMFADeviceRequest, callback func(response *ram.CreateVirtualMFADeviceResponse, err error)) <-chan int {
	res0 := w.obj.CreateVirtualMFADeviceWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) CreateVirtualMFADeviceWithChan(request *ram.CreateVirtualMFADeviceRequest) (<-chan *ram.CreateVirtualMFADeviceResponse, <-chan error) {
	res0, res1 := w.obj.CreateVirtualMFADeviceWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) DeleteAccessKey(ctx context.Context, request *ram.DeleteAccessKeyRequest) (*ram.DeleteAccessKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.DeleteAccessKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteAccessKey", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteAccessKey", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteAccessKey", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.DeleteAccessKey")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.DeleteAccessKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.DeleteAccessKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.DeleteAccessKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteAccessKey(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) DeleteAccessKeyWithCallback(request *ram.DeleteAccessKeyRequest, callback func(response *ram.DeleteAccessKeyResponse, err error)) <-chan int {
	res0 := w.obj.DeleteAccessKeyWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) DeleteAccessKeyWithChan(request *ram.DeleteAccessKeyRequest) (<-chan *ram.DeleteAccessKeyResponse, <-chan error) {
	res0, res1 := w.obj.DeleteAccessKeyWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) DeleteGroup(ctx context.Context, request *ram.DeleteGroupRequest) (*ram.DeleteGroupResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.DeleteGroupResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteGroup", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteGroup", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteGroup", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.DeleteGroup")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.DeleteGroup", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.DeleteGroup", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.DeleteGroup", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteGroup(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) DeleteGroupWithCallback(request *ram.DeleteGroupRequest, callback func(response *ram.DeleteGroupResponse, err error)) <-chan int {
	res0 := w.obj.DeleteGroupWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) DeleteGroupWithChan(request *ram.DeleteGroupRequest) (<-chan *ram.DeleteGroupResponse, <-chan error) {
	res0, res1 := w.obj.DeleteGroupWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) DeleteLoginProfile(ctx context.Context, request *ram.DeleteLoginProfileRequest) (*ram.DeleteLoginProfileResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.DeleteLoginProfileResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteLoginProfile", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteLoginProfile", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteLoginProfile", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.DeleteLoginProfile")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.DeleteLoginProfile", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.DeleteLoginProfile", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.DeleteLoginProfile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteLoginProfile(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) DeleteLoginProfileWithCallback(request *ram.DeleteLoginProfileRequest, callback func(response *ram.DeleteLoginProfileResponse, err error)) <-chan int {
	res0 := w.obj.DeleteLoginProfileWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) DeleteLoginProfileWithChan(request *ram.DeleteLoginProfileRequest) (<-chan *ram.DeleteLoginProfileResponse, <-chan error) {
	res0, res1 := w.obj.DeleteLoginProfileWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) DeletePolicy(ctx context.Context, request *ram.DeletePolicyRequest) (*ram.DeletePolicyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.DeletePolicyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeletePolicy", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeletePolicy", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeletePolicy", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.DeletePolicy")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.DeletePolicy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.DeletePolicy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.DeletePolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeletePolicy(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) DeletePolicyVersion(ctx context.Context, request *ram.DeletePolicyVersionRequest) (*ram.DeletePolicyVersionResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.DeletePolicyVersionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeletePolicyVersion", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeletePolicyVersion", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeletePolicyVersion", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.DeletePolicyVersion")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.DeletePolicyVersion", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.DeletePolicyVersion", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.DeletePolicyVersion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeletePolicyVersion(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) DeletePolicyVersionWithCallback(request *ram.DeletePolicyVersionRequest, callback func(response *ram.DeletePolicyVersionResponse, err error)) <-chan int {
	res0 := w.obj.DeletePolicyVersionWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) DeletePolicyVersionWithChan(request *ram.DeletePolicyVersionRequest) (<-chan *ram.DeletePolicyVersionResponse, <-chan error) {
	res0, res1 := w.obj.DeletePolicyVersionWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) DeletePolicyWithCallback(request *ram.DeletePolicyRequest, callback func(response *ram.DeletePolicyResponse, err error)) <-chan int {
	res0 := w.obj.DeletePolicyWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) DeletePolicyWithChan(request *ram.DeletePolicyRequest) (<-chan *ram.DeletePolicyResponse, <-chan error) {
	res0, res1 := w.obj.DeletePolicyWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) DeleteRole(ctx context.Context, request *ram.DeleteRoleRequest) (*ram.DeleteRoleResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.DeleteRoleResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteRole", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteRole", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteRole", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.DeleteRole")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.DeleteRole", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.DeleteRole", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.DeleteRole", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteRole(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) DeleteRoleWithCallback(request *ram.DeleteRoleRequest, callback func(response *ram.DeleteRoleResponse, err error)) <-chan int {
	res0 := w.obj.DeleteRoleWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) DeleteRoleWithChan(request *ram.DeleteRoleRequest) (<-chan *ram.DeleteRoleResponse, <-chan error) {
	res0, res1 := w.obj.DeleteRoleWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) DeleteUser(ctx context.Context, request *ram.DeleteUserRequest) (*ram.DeleteUserResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.DeleteUserResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteUser", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteUser", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteUser", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.DeleteUser")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.DeleteUser", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.DeleteUser", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.DeleteUser", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteUser(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) DeleteUserWithCallback(request *ram.DeleteUserRequest, callback func(response *ram.DeleteUserResponse, err error)) <-chan int {
	res0 := w.obj.DeleteUserWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) DeleteUserWithChan(request *ram.DeleteUserRequest) (<-chan *ram.DeleteUserResponse, <-chan error) {
	res0, res1 := w.obj.DeleteUserWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) DeleteVirtualMFADevice(ctx context.Context, request *ram.DeleteVirtualMFADeviceRequest) (*ram.DeleteVirtualMFADeviceResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.DeleteVirtualMFADeviceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DeleteVirtualMFADevice", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DeleteVirtualMFADevice", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DeleteVirtualMFADevice", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.DeleteVirtualMFADevice")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.DeleteVirtualMFADevice", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.DeleteVirtualMFADevice", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.DeleteVirtualMFADevice", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DeleteVirtualMFADevice(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) DeleteVirtualMFADeviceWithCallback(request *ram.DeleteVirtualMFADeviceRequest, callback func(response *ram.DeleteVirtualMFADeviceResponse, err error)) <-chan int {
	res0 := w.obj.DeleteVirtualMFADeviceWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) DeleteVirtualMFADeviceWithChan(request *ram.DeleteVirtualMFADeviceRequest) (<-chan *ram.DeleteVirtualMFADeviceResponse, <-chan error) {
	res0, res1 := w.obj.DeleteVirtualMFADeviceWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) DetachPolicyFromGroup(ctx context.Context, request *ram.DetachPolicyFromGroupRequest) (*ram.DetachPolicyFromGroupResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.DetachPolicyFromGroupResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DetachPolicyFromGroup", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DetachPolicyFromGroup", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DetachPolicyFromGroup", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.DetachPolicyFromGroup")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.DetachPolicyFromGroup", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.DetachPolicyFromGroup", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.DetachPolicyFromGroup", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DetachPolicyFromGroup(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) DetachPolicyFromGroupWithCallback(request *ram.DetachPolicyFromGroupRequest, callback func(response *ram.DetachPolicyFromGroupResponse, err error)) <-chan int {
	res0 := w.obj.DetachPolicyFromGroupWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) DetachPolicyFromGroupWithChan(request *ram.DetachPolicyFromGroupRequest) (<-chan *ram.DetachPolicyFromGroupResponse, <-chan error) {
	res0, res1 := w.obj.DetachPolicyFromGroupWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) DetachPolicyFromRole(ctx context.Context, request *ram.DetachPolicyFromRoleRequest) (*ram.DetachPolicyFromRoleResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.DetachPolicyFromRoleResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DetachPolicyFromRole", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DetachPolicyFromRole", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DetachPolicyFromRole", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.DetachPolicyFromRole")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.DetachPolicyFromRole", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.DetachPolicyFromRole", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.DetachPolicyFromRole", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DetachPolicyFromRole(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) DetachPolicyFromRoleWithCallback(request *ram.DetachPolicyFromRoleRequest, callback func(response *ram.DetachPolicyFromRoleResponse, err error)) <-chan int {
	res0 := w.obj.DetachPolicyFromRoleWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) DetachPolicyFromRoleWithChan(request *ram.DetachPolicyFromRoleRequest) (<-chan *ram.DetachPolicyFromRoleResponse, <-chan error) {
	res0, res1 := w.obj.DetachPolicyFromRoleWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) DetachPolicyFromUser(ctx context.Context, request *ram.DetachPolicyFromUserRequest) (*ram.DetachPolicyFromUserResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.DetachPolicyFromUserResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.DetachPolicyFromUser", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.DetachPolicyFromUser", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.DetachPolicyFromUser", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.DetachPolicyFromUser")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.DetachPolicyFromUser", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.DetachPolicyFromUser", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.DetachPolicyFromUser", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.DetachPolicyFromUser(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) DetachPolicyFromUserWithCallback(request *ram.DetachPolicyFromUserRequest, callback func(response *ram.DetachPolicyFromUserResponse, err error)) <-chan int {
	res0 := w.obj.DetachPolicyFromUserWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) DetachPolicyFromUserWithChan(request *ram.DetachPolicyFromUserRequest) (<-chan *ram.DetachPolicyFromUserResponse, <-chan error) {
	res0, res1 := w.obj.DetachPolicyFromUserWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) GetAccessKeyLastUsed(ctx context.Context, request *ram.GetAccessKeyLastUsedRequest) (*ram.GetAccessKeyLastUsedResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.GetAccessKeyLastUsedResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetAccessKeyLastUsed", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetAccessKeyLastUsed", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetAccessKeyLastUsed", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.GetAccessKeyLastUsed")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.GetAccessKeyLastUsed", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.GetAccessKeyLastUsed", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.GetAccessKeyLastUsed", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetAccessKeyLastUsed(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) GetAccessKeyLastUsedWithCallback(request *ram.GetAccessKeyLastUsedRequest, callback func(response *ram.GetAccessKeyLastUsedResponse, err error)) <-chan int {
	res0 := w.obj.GetAccessKeyLastUsedWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) GetAccessKeyLastUsedWithChan(request *ram.GetAccessKeyLastUsedRequest) (<-chan *ram.GetAccessKeyLastUsedResponse, <-chan error) {
	res0, res1 := w.obj.GetAccessKeyLastUsedWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) GetAccountAlias(ctx context.Context, request *ram.GetAccountAliasRequest) (*ram.GetAccountAliasResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.GetAccountAliasResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetAccountAlias", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetAccountAlias", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetAccountAlias", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.GetAccountAlias")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.GetAccountAlias", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.GetAccountAlias", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.GetAccountAlias", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetAccountAlias(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) GetAccountAliasWithCallback(request *ram.GetAccountAliasRequest, callback func(response *ram.GetAccountAliasResponse, err error)) <-chan int {
	res0 := w.obj.GetAccountAliasWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) GetAccountAliasWithChan(request *ram.GetAccountAliasRequest) (<-chan *ram.GetAccountAliasResponse, <-chan error) {
	res0, res1 := w.obj.GetAccountAliasWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) GetGroup(ctx context.Context, request *ram.GetGroupRequest) (*ram.GetGroupResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.GetGroupResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetGroup", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetGroup", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetGroup", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.GetGroup")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.GetGroup", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.GetGroup", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.GetGroup", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetGroup(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) GetGroupWithCallback(request *ram.GetGroupRequest, callback func(response *ram.GetGroupResponse, err error)) <-chan int {
	res0 := w.obj.GetGroupWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) GetGroupWithChan(request *ram.GetGroupRequest) (<-chan *ram.GetGroupResponse, <-chan error) {
	res0, res1 := w.obj.GetGroupWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) GetLoginProfile(ctx context.Context, request *ram.GetLoginProfileRequest) (*ram.GetLoginProfileResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.GetLoginProfileResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetLoginProfile", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetLoginProfile", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetLoginProfile", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.GetLoginProfile")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.GetLoginProfile", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.GetLoginProfile", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.GetLoginProfile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetLoginProfile(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) GetLoginProfileWithCallback(request *ram.GetLoginProfileRequest, callback func(response *ram.GetLoginProfileResponse, err error)) <-chan int {
	res0 := w.obj.GetLoginProfileWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) GetLoginProfileWithChan(request *ram.GetLoginProfileRequest) (<-chan *ram.GetLoginProfileResponse, <-chan error) {
	res0, res1 := w.obj.GetLoginProfileWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) GetPasswordPolicy(ctx context.Context, request *ram.GetPasswordPolicyRequest) (*ram.GetPasswordPolicyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.GetPasswordPolicyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetPasswordPolicy", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetPasswordPolicy", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetPasswordPolicy", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.GetPasswordPolicy")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.GetPasswordPolicy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.GetPasswordPolicy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.GetPasswordPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetPasswordPolicy(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) GetPasswordPolicyWithCallback(request *ram.GetPasswordPolicyRequest, callback func(response *ram.GetPasswordPolicyResponse, err error)) <-chan int {
	res0 := w.obj.GetPasswordPolicyWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) GetPasswordPolicyWithChan(request *ram.GetPasswordPolicyRequest) (<-chan *ram.GetPasswordPolicyResponse, <-chan error) {
	res0, res1 := w.obj.GetPasswordPolicyWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) GetPolicy(ctx context.Context, request *ram.GetPolicyRequest) (*ram.GetPolicyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.GetPolicyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetPolicy", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetPolicy", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetPolicy", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.GetPolicy")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.GetPolicy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.GetPolicy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.GetPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetPolicy(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) GetPolicyVersion(ctx context.Context, request *ram.GetPolicyVersionRequest) (*ram.GetPolicyVersionResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.GetPolicyVersionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetPolicyVersion", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetPolicyVersion", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetPolicyVersion", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.GetPolicyVersion")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.GetPolicyVersion", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.GetPolicyVersion", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.GetPolicyVersion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetPolicyVersion(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) GetPolicyVersionWithCallback(request *ram.GetPolicyVersionRequest, callback func(response *ram.GetPolicyVersionResponse, err error)) <-chan int {
	res0 := w.obj.GetPolicyVersionWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) GetPolicyVersionWithChan(request *ram.GetPolicyVersionRequest) (<-chan *ram.GetPolicyVersionResponse, <-chan error) {
	res0, res1 := w.obj.GetPolicyVersionWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) GetPolicyWithCallback(request *ram.GetPolicyRequest, callback func(response *ram.GetPolicyResponse, err error)) <-chan int {
	res0 := w.obj.GetPolicyWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) GetPolicyWithChan(request *ram.GetPolicyRequest) (<-chan *ram.GetPolicyResponse, <-chan error) {
	res0, res1 := w.obj.GetPolicyWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) GetRole(ctx context.Context, request *ram.GetRoleRequest) (*ram.GetRoleResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.GetRoleResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetRole", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetRole", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetRole", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.GetRole")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.GetRole", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.GetRole", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.GetRole", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetRole(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) GetRoleWithCallback(request *ram.GetRoleRequest, callback func(response *ram.GetRoleResponse, err error)) <-chan int {
	res0 := w.obj.GetRoleWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) GetRoleWithChan(request *ram.GetRoleRequest) (<-chan *ram.GetRoleResponse, <-chan error) {
	res0, res1 := w.obj.GetRoleWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) GetSecurityPreference(ctx context.Context, request *ram.GetSecurityPreferenceRequest) (*ram.GetSecurityPreferenceResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.GetSecurityPreferenceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetSecurityPreference", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetSecurityPreference", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetSecurityPreference", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.GetSecurityPreference")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.GetSecurityPreference", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.GetSecurityPreference", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.GetSecurityPreference", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetSecurityPreference(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) GetSecurityPreferenceWithCallback(request *ram.GetSecurityPreferenceRequest, callback func(response *ram.GetSecurityPreferenceResponse, err error)) <-chan int {
	res0 := w.obj.GetSecurityPreferenceWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) GetSecurityPreferenceWithChan(request *ram.GetSecurityPreferenceRequest) (<-chan *ram.GetSecurityPreferenceResponse, <-chan error) {
	res0, res1 := w.obj.GetSecurityPreferenceWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) GetUser(ctx context.Context, request *ram.GetUserRequest) (*ram.GetUserResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.GetUserResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetUser", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetUser", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetUser", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.GetUser")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.GetUser", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.GetUser", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.GetUser", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetUser(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) GetUserMFAInfo(ctx context.Context, request *ram.GetUserMFAInfoRequest) (*ram.GetUserMFAInfoResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.GetUserMFAInfoResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.GetUserMFAInfo", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.GetUserMFAInfo", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.GetUserMFAInfo", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.GetUserMFAInfo")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.GetUserMFAInfo", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.GetUserMFAInfo", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.GetUserMFAInfo", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.GetUserMFAInfo(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) GetUserMFAInfoWithCallback(request *ram.GetUserMFAInfoRequest, callback func(response *ram.GetUserMFAInfoResponse, err error)) <-chan int {
	res0 := w.obj.GetUserMFAInfoWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) GetUserMFAInfoWithChan(request *ram.GetUserMFAInfoRequest) (<-chan *ram.GetUserMFAInfoResponse, <-chan error) {
	res0, res1 := w.obj.GetUserMFAInfoWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) GetUserWithCallback(request *ram.GetUserRequest, callback func(response *ram.GetUserResponse, err error)) <-chan int {
	res0 := w.obj.GetUserWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) GetUserWithChan(request *ram.GetUserRequest) (<-chan *ram.GetUserResponse, <-chan error) {
	res0, res1 := w.obj.GetUserWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ListAccessKeys(ctx context.Context, request *ram.ListAccessKeysRequest) (*ram.ListAccessKeysResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ListAccessKeysResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListAccessKeys", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListAccessKeys", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListAccessKeys", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ListAccessKeys")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ListAccessKeys", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ListAccessKeys", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ListAccessKeys", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListAccessKeys(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ListAccessKeysWithCallback(request *ram.ListAccessKeysRequest, callback func(response *ram.ListAccessKeysResponse, err error)) <-chan int {
	res0 := w.obj.ListAccessKeysWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ListAccessKeysWithChan(request *ram.ListAccessKeysRequest) (<-chan *ram.ListAccessKeysResponse, <-chan error) {
	res0, res1 := w.obj.ListAccessKeysWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ListEntitiesForPolicy(ctx context.Context, request *ram.ListEntitiesForPolicyRequest) (*ram.ListEntitiesForPolicyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ListEntitiesForPolicyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListEntitiesForPolicy", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListEntitiesForPolicy", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListEntitiesForPolicy", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ListEntitiesForPolicy")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ListEntitiesForPolicy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ListEntitiesForPolicy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ListEntitiesForPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListEntitiesForPolicy(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ListEntitiesForPolicyWithCallback(request *ram.ListEntitiesForPolicyRequest, callback func(response *ram.ListEntitiesForPolicyResponse, err error)) <-chan int {
	res0 := w.obj.ListEntitiesForPolicyWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ListEntitiesForPolicyWithChan(request *ram.ListEntitiesForPolicyRequest) (<-chan *ram.ListEntitiesForPolicyResponse, <-chan error) {
	res0, res1 := w.obj.ListEntitiesForPolicyWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ListGroups(ctx context.Context, request *ram.ListGroupsRequest) (*ram.ListGroupsResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ListGroupsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListGroups", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListGroups", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListGroups", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ListGroups")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ListGroups", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ListGroups", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ListGroups", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListGroups(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ListGroupsForUser(ctx context.Context, request *ram.ListGroupsForUserRequest) (*ram.ListGroupsForUserResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ListGroupsForUserResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListGroupsForUser", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListGroupsForUser", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListGroupsForUser", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ListGroupsForUser")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ListGroupsForUser", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ListGroupsForUser", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ListGroupsForUser", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListGroupsForUser(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ListGroupsForUserWithCallback(request *ram.ListGroupsForUserRequest, callback func(response *ram.ListGroupsForUserResponse, err error)) <-chan int {
	res0 := w.obj.ListGroupsForUserWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ListGroupsForUserWithChan(request *ram.ListGroupsForUserRequest) (<-chan *ram.ListGroupsForUserResponse, <-chan error) {
	res0, res1 := w.obj.ListGroupsForUserWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ListGroupsWithCallback(request *ram.ListGroupsRequest, callback func(response *ram.ListGroupsResponse, err error)) <-chan int {
	res0 := w.obj.ListGroupsWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ListGroupsWithChan(request *ram.ListGroupsRequest) (<-chan *ram.ListGroupsResponse, <-chan error) {
	res0, res1 := w.obj.ListGroupsWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ListPolicies(ctx context.Context, request *ram.ListPoliciesRequest) (*ram.ListPoliciesResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ListPoliciesResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListPolicies", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListPolicies", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListPolicies", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ListPolicies")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ListPolicies", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ListPolicies", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ListPolicies", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListPolicies(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ListPoliciesForGroup(ctx context.Context, request *ram.ListPoliciesForGroupRequest) (*ram.ListPoliciesForGroupResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ListPoliciesForGroupResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListPoliciesForGroup", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListPoliciesForGroup", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListPoliciesForGroup", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ListPoliciesForGroup")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ListPoliciesForGroup", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ListPoliciesForGroup", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ListPoliciesForGroup", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListPoliciesForGroup(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ListPoliciesForGroupWithCallback(request *ram.ListPoliciesForGroupRequest, callback func(response *ram.ListPoliciesForGroupResponse, err error)) <-chan int {
	res0 := w.obj.ListPoliciesForGroupWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ListPoliciesForGroupWithChan(request *ram.ListPoliciesForGroupRequest) (<-chan *ram.ListPoliciesForGroupResponse, <-chan error) {
	res0, res1 := w.obj.ListPoliciesForGroupWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ListPoliciesForRole(ctx context.Context, request *ram.ListPoliciesForRoleRequest) (*ram.ListPoliciesForRoleResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ListPoliciesForRoleResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListPoliciesForRole", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListPoliciesForRole", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListPoliciesForRole", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ListPoliciesForRole")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ListPoliciesForRole", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ListPoliciesForRole", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ListPoliciesForRole", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListPoliciesForRole(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ListPoliciesForRoleWithCallback(request *ram.ListPoliciesForRoleRequest, callback func(response *ram.ListPoliciesForRoleResponse, err error)) <-chan int {
	res0 := w.obj.ListPoliciesForRoleWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ListPoliciesForRoleWithChan(request *ram.ListPoliciesForRoleRequest) (<-chan *ram.ListPoliciesForRoleResponse, <-chan error) {
	res0, res1 := w.obj.ListPoliciesForRoleWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ListPoliciesForUser(ctx context.Context, request *ram.ListPoliciesForUserRequest) (*ram.ListPoliciesForUserResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ListPoliciesForUserResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListPoliciesForUser", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListPoliciesForUser", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListPoliciesForUser", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ListPoliciesForUser")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ListPoliciesForUser", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ListPoliciesForUser", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ListPoliciesForUser", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListPoliciesForUser(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ListPoliciesForUserWithCallback(request *ram.ListPoliciesForUserRequest, callback func(response *ram.ListPoliciesForUserResponse, err error)) <-chan int {
	res0 := w.obj.ListPoliciesForUserWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ListPoliciesForUserWithChan(request *ram.ListPoliciesForUserRequest) (<-chan *ram.ListPoliciesForUserResponse, <-chan error) {
	res0, res1 := w.obj.ListPoliciesForUserWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ListPoliciesWithCallback(request *ram.ListPoliciesRequest, callback func(response *ram.ListPoliciesResponse, err error)) <-chan int {
	res0 := w.obj.ListPoliciesWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ListPoliciesWithChan(request *ram.ListPoliciesRequest) (<-chan *ram.ListPoliciesResponse, <-chan error) {
	res0, res1 := w.obj.ListPoliciesWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ListPolicyVersions(ctx context.Context, request *ram.ListPolicyVersionsRequest) (*ram.ListPolicyVersionsResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ListPolicyVersionsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListPolicyVersions", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListPolicyVersions", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListPolicyVersions", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ListPolicyVersions")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ListPolicyVersions", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ListPolicyVersions", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ListPolicyVersions", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListPolicyVersions(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ListPolicyVersionsWithCallback(request *ram.ListPolicyVersionsRequest, callback func(response *ram.ListPolicyVersionsResponse, err error)) <-chan int {
	res0 := w.obj.ListPolicyVersionsWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ListPolicyVersionsWithChan(request *ram.ListPolicyVersionsRequest) (<-chan *ram.ListPolicyVersionsResponse, <-chan error) {
	res0, res1 := w.obj.ListPolicyVersionsWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ListRoles(ctx context.Context, request *ram.ListRolesRequest) (*ram.ListRolesResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ListRolesResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListRoles", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListRoles", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListRoles", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ListRoles")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ListRoles", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ListRoles", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ListRoles", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListRoles(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ListRolesWithCallback(request *ram.ListRolesRequest, callback func(response *ram.ListRolesResponse, err error)) <-chan int {
	res0 := w.obj.ListRolesWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ListRolesWithChan(request *ram.ListRolesRequest) (<-chan *ram.ListRolesResponse, <-chan error) {
	res0, res1 := w.obj.ListRolesWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ListUsers(ctx context.Context, request *ram.ListUsersRequest) (*ram.ListUsersResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ListUsersResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListUsers", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListUsers", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListUsers", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ListUsers")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ListUsers", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ListUsers", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ListUsers", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListUsers(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ListUsersForGroup(ctx context.Context, request *ram.ListUsersForGroupRequest) (*ram.ListUsersForGroupResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ListUsersForGroupResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListUsersForGroup", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListUsersForGroup", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListUsersForGroup", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ListUsersForGroup")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ListUsersForGroup", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ListUsersForGroup", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ListUsersForGroup", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListUsersForGroup(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ListUsersForGroupWithCallback(request *ram.ListUsersForGroupRequest, callback func(response *ram.ListUsersForGroupResponse, err error)) <-chan int {
	res0 := w.obj.ListUsersForGroupWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ListUsersForGroupWithChan(request *ram.ListUsersForGroupRequest) (<-chan *ram.ListUsersForGroupResponse, <-chan error) {
	res0, res1 := w.obj.ListUsersForGroupWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ListUsersWithCallback(request *ram.ListUsersRequest, callback func(response *ram.ListUsersResponse, err error)) <-chan int {
	res0 := w.obj.ListUsersWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ListUsersWithChan(request *ram.ListUsersRequest) (<-chan *ram.ListUsersResponse, <-chan error) {
	res0, res1 := w.obj.ListUsersWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) ListVirtualMFADevices(ctx context.Context, request *ram.ListVirtualMFADevicesRequest) (*ram.ListVirtualMFADevicesResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.ListVirtualMFADevicesResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.ListVirtualMFADevices", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.ListVirtualMFADevices", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.ListVirtualMFADevices", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.ListVirtualMFADevices")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.ListVirtualMFADevices", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.ListVirtualMFADevices", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.ListVirtualMFADevices", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.ListVirtualMFADevices(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) ListVirtualMFADevicesWithCallback(request *ram.ListVirtualMFADevicesRequest, callback func(response *ram.ListVirtualMFADevicesResponse, err error)) <-chan int {
	res0 := w.obj.ListVirtualMFADevicesWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) ListVirtualMFADevicesWithChan(request *ram.ListVirtualMFADevicesRequest) (<-chan *ram.ListVirtualMFADevicesResponse, <-chan error) {
	res0, res1 := w.obj.ListVirtualMFADevicesWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) RemoveUserFromGroup(ctx context.Context, request *ram.RemoveUserFromGroupRequest) (*ram.RemoveUserFromGroupResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.RemoveUserFromGroupResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.RemoveUserFromGroup", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.RemoveUserFromGroup", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.RemoveUserFromGroup", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.RemoveUserFromGroup")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.RemoveUserFromGroup", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.RemoveUserFromGroup", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.RemoveUserFromGroup", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.RemoveUserFromGroup(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) RemoveUserFromGroupWithCallback(request *ram.RemoveUserFromGroupRequest, callback func(response *ram.RemoveUserFromGroupResponse, err error)) <-chan int {
	res0 := w.obj.RemoveUserFromGroupWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) RemoveUserFromGroupWithChan(request *ram.RemoveUserFromGroupRequest) (<-chan *ram.RemoveUserFromGroupResponse, <-chan error) {
	res0, res1 := w.obj.RemoveUserFromGroupWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) SetAccountAlias(ctx context.Context, request *ram.SetAccountAliasRequest) (*ram.SetAccountAliasResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.SetAccountAliasResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SetAccountAlias", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SetAccountAlias", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SetAccountAlias", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.SetAccountAlias")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.SetAccountAlias", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.SetAccountAlias", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.SetAccountAlias", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.SetAccountAlias(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) SetAccountAliasWithCallback(request *ram.SetAccountAliasRequest, callback func(response *ram.SetAccountAliasResponse, err error)) <-chan int {
	res0 := w.obj.SetAccountAliasWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) SetAccountAliasWithChan(request *ram.SetAccountAliasRequest) (<-chan *ram.SetAccountAliasResponse, <-chan error) {
	res0, res1 := w.obj.SetAccountAliasWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) SetDefaultPolicyVersion(ctx context.Context, request *ram.SetDefaultPolicyVersionRequest) (*ram.SetDefaultPolicyVersionResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.SetDefaultPolicyVersionResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SetDefaultPolicyVersion", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SetDefaultPolicyVersion", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SetDefaultPolicyVersion", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.SetDefaultPolicyVersion")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.SetDefaultPolicyVersion", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.SetDefaultPolicyVersion", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.SetDefaultPolicyVersion", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.SetDefaultPolicyVersion(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) SetDefaultPolicyVersionWithCallback(request *ram.SetDefaultPolicyVersionRequest, callback func(response *ram.SetDefaultPolicyVersionResponse, err error)) <-chan int {
	res0 := w.obj.SetDefaultPolicyVersionWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) SetDefaultPolicyVersionWithChan(request *ram.SetDefaultPolicyVersionRequest) (<-chan *ram.SetDefaultPolicyVersionResponse, <-chan error) {
	res0, res1 := w.obj.SetDefaultPolicyVersionWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) SetPasswordPolicy(ctx context.Context, request *ram.SetPasswordPolicyRequest) (*ram.SetPasswordPolicyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.SetPasswordPolicyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SetPasswordPolicy", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SetPasswordPolicy", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SetPasswordPolicy", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.SetPasswordPolicy")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.SetPasswordPolicy", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.SetPasswordPolicy", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.SetPasswordPolicy", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.SetPasswordPolicy(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) SetPasswordPolicyWithCallback(request *ram.SetPasswordPolicyRequest, callback func(response *ram.SetPasswordPolicyResponse, err error)) <-chan int {
	res0 := w.obj.SetPasswordPolicyWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) SetPasswordPolicyWithChan(request *ram.SetPasswordPolicyRequest) (<-chan *ram.SetPasswordPolicyResponse, <-chan error) {
	res0, res1 := w.obj.SetPasswordPolicyWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) SetSecurityPreference(ctx context.Context, request *ram.SetSecurityPreferenceRequest) (*ram.SetSecurityPreferenceResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.SetSecurityPreferenceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.SetSecurityPreference", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.SetSecurityPreference", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.SetSecurityPreference", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.SetSecurityPreference")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.SetSecurityPreference", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.SetSecurityPreference", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.SetSecurityPreference", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.SetSecurityPreference(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) SetSecurityPreferenceWithCallback(request *ram.SetSecurityPreferenceRequest, callback func(response *ram.SetSecurityPreferenceResponse, err error)) <-chan int {
	res0 := w.obj.SetSecurityPreferenceWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) SetSecurityPreferenceWithChan(request *ram.SetSecurityPreferenceRequest) (<-chan *ram.SetSecurityPreferenceResponse, <-chan error) {
	res0, res1 := w.obj.SetSecurityPreferenceWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) UnbindMFADevice(ctx context.Context, request *ram.UnbindMFADeviceRequest) (*ram.UnbindMFADeviceResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.UnbindMFADeviceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UnbindMFADevice", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UnbindMFADevice", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UnbindMFADevice", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.UnbindMFADevice")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.UnbindMFADevice", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.UnbindMFADevice", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.UnbindMFADevice", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UnbindMFADevice(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) UnbindMFADeviceWithCallback(request *ram.UnbindMFADeviceRequest, callback func(response *ram.UnbindMFADeviceResponse, err error)) <-chan int {
	res0 := w.obj.UnbindMFADeviceWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) UnbindMFADeviceWithChan(request *ram.UnbindMFADeviceRequest) (<-chan *ram.UnbindMFADeviceResponse, <-chan error) {
	res0, res1 := w.obj.UnbindMFADeviceWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) UpdateAccessKey(ctx context.Context, request *ram.UpdateAccessKeyRequest) (*ram.UpdateAccessKeyResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.UpdateAccessKeyResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateAccessKey", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateAccessKey", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateAccessKey", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.UpdateAccessKey")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.UpdateAccessKey", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.UpdateAccessKey", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.UpdateAccessKey", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateAccessKey(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) UpdateAccessKeyWithCallback(request *ram.UpdateAccessKeyRequest, callback func(response *ram.UpdateAccessKeyResponse, err error)) <-chan int {
	res0 := w.obj.UpdateAccessKeyWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) UpdateAccessKeyWithChan(request *ram.UpdateAccessKeyRequest) (<-chan *ram.UpdateAccessKeyResponse, <-chan error) {
	res0, res1 := w.obj.UpdateAccessKeyWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) UpdateGroup(ctx context.Context, request *ram.UpdateGroupRequest) (*ram.UpdateGroupResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.UpdateGroupResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateGroup", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateGroup", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateGroup", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.UpdateGroup")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.UpdateGroup", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.UpdateGroup", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.UpdateGroup", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateGroup(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) UpdateGroupWithCallback(request *ram.UpdateGroupRequest, callback func(response *ram.UpdateGroupResponse, err error)) <-chan int {
	res0 := w.obj.UpdateGroupWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) UpdateGroupWithChan(request *ram.UpdateGroupRequest) (<-chan *ram.UpdateGroupResponse, <-chan error) {
	res0, res1 := w.obj.UpdateGroupWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) UpdateLoginProfile(ctx context.Context, request *ram.UpdateLoginProfileRequest) (*ram.UpdateLoginProfileResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.UpdateLoginProfileResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateLoginProfile", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateLoginProfile", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateLoginProfile", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.UpdateLoginProfile")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.UpdateLoginProfile", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.UpdateLoginProfile", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.UpdateLoginProfile", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateLoginProfile(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) UpdateLoginProfileWithCallback(request *ram.UpdateLoginProfileRequest, callback func(response *ram.UpdateLoginProfileResponse, err error)) <-chan int {
	res0 := w.obj.UpdateLoginProfileWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) UpdateLoginProfileWithChan(request *ram.UpdateLoginProfileRequest) (<-chan *ram.UpdateLoginProfileResponse, <-chan error) {
	res0, res1 := w.obj.UpdateLoginProfileWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) UpdateRole(ctx context.Context, request *ram.UpdateRoleRequest) (*ram.UpdateRoleResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.UpdateRoleResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateRole", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateRole", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateRole", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.UpdateRole")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.UpdateRole", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.UpdateRole", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.UpdateRole", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateRole(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) UpdateRoleWithCallback(request *ram.UpdateRoleRequest, callback func(response *ram.UpdateRoleResponse, err error)) <-chan int {
	res0 := w.obj.UpdateRoleWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) UpdateRoleWithChan(request *ram.UpdateRoleRequest) (<-chan *ram.UpdateRoleResponse, <-chan error) {
	res0, res1 := w.obj.UpdateRoleWithChan(request)
	return res0, res1
}

func (w *RAMClientWrapper) UpdateUser(ctx context.Context, request *ram.UpdateUserRequest) (*ram.UpdateUserResponse, error) {
	ctxOptions := FromContext(ctx)
	var response *ram.UpdateUserResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.Client.UpdateUser", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.Client.UpdateUser", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.Client.UpdateUser", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "ram.Client.UpdateUser")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("ram.Client.UpdateUser", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("ram.Client.UpdateUser", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("ram.Client.UpdateUser", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		response, err = w.obj.UpdateUser(request)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return response, err
}

func (w *RAMClientWrapper) UpdateUserWithCallback(request *ram.UpdateUserRequest, callback func(response *ram.UpdateUserResponse, err error)) <-chan int {
	res0 := w.obj.UpdateUserWithCallback(request, callback)
	return res0
}

func (w *RAMClientWrapper) UpdateUserWithChan(request *ram.UpdateUserRequest) (<-chan *ram.UpdateUserResponse, <-chan error) {
	res0, res1 := w.obj.UpdateUserWithChan(request)
	return res0, res1
}
