// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"

	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/opentracing/opentracing-go"
)

type ACMConfigClientWrapper struct {
	obj     *config_client.ConfigClient
	retry   *Retry
	options *WrapperOptions
}

func (w *ACMConfigClientWrapper) Unwrap() *config_client.ConfigClient {
	return w.obj
}

func (w *ACMConfigClientWrapper) CancelListenConfig(ctx context.Context, param vo.ConfigParam) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.CancelListenConfig")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.CancelListenConfig(param)
		return err
	})
	return err
}

func (w *ACMConfigClientWrapper) DeleteConfig(ctx context.Context, param vo.ConfigParam) (bool, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.DeleteConfig")
		defer span.Finish()
	}

	var deleted bool
	var err error
	err = w.retry.Do(func() error {
		deleted, err = w.obj.DeleteConfig(param)
		return err
	})
	return deleted, err
}

func (w *ACMConfigClientWrapper) GetConfig(ctx context.Context, param vo.ConfigParam) (string, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.GetConfig")
		defer span.Finish()
	}

	var content string
	var err error
	err = w.retry.Do(func() error {
		content, err = w.obj.GetConfig(param)
		return err
	})
	return content, err
}

func (w *ACMConfigClientWrapper) ListenConfig(ctx context.Context, param vo.ConfigParam) error {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.ListenConfig")
		defer span.Finish()
	}

	var err error
	err = w.retry.Do(func() error {
		err = w.obj.ListenConfig(param)
		return err
	})
	return err
}

func (w *ACMConfigClientWrapper) PublishConfig(ctx context.Context, param vo.ConfigParam) (bool, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.PublishConfig")
		defer span.Finish()
	}

	var published bool
	var err error
	err = w.retry.Do(func() error {
		published, err = w.obj.PublishConfig(param)
		return err
	})
	return published, err
}

func (w *ACMConfigClientWrapper) SearchConfig(ctx context.Context, param vo.SearchConfigParm) (*model.ConfigPage, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "config_client.ConfigClient.SearchConfig")
		defer span.Finish()
	}

	var res0 *model.ConfigPage
	var err error
	err = w.retry.Do(func() error {
		res0, err = w.obj.SearchConfig(param)
		return err
	})
	return res0, err
}
