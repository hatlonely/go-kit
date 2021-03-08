// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
)

type ESAliasServiceWrapper struct {
	obj                *elastic.AliasService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESAliasServiceWrapper) Unwrap() *elastic.AliasService {
	return w.obj
}

type ESAliasesServiceWrapper struct {
	obj                *elastic.AliasesService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESAliasesServiceWrapper) Unwrap() *elastic.AliasesService {
	return w.obj
}

type ESBulkProcessorServiceWrapper struct {
	obj                *elastic.BulkProcessorService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESBulkProcessorServiceWrapper) Unwrap() *elastic.BulkProcessorService {
	return w.obj
}

type ESBulkServiceWrapper struct {
	obj                *elastic.BulkService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESBulkServiceWrapper) Unwrap() *elastic.BulkService {
	return w.obj
}

type ESCatAliasesServiceWrapper struct {
	obj                *elastic.CatAliasesService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESCatAliasesServiceWrapper) Unwrap() *elastic.CatAliasesService {
	return w.obj
}

type ESCatAllocationServiceWrapper struct {
	obj                *elastic.CatAllocationService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESCatAllocationServiceWrapper) Unwrap() *elastic.CatAllocationService {
	return w.obj
}

type ESCatCountServiceWrapper struct {
	obj                *elastic.CatCountService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESCatCountServiceWrapper) Unwrap() *elastic.CatCountService {
	return w.obj
}

type ESCatHealthServiceWrapper struct {
	obj                *elastic.CatHealthService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESCatHealthServiceWrapper) Unwrap() *elastic.CatHealthService {
	return w.obj
}

type ESCatIndicesServiceWrapper struct {
	obj                *elastic.CatIndicesService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESCatIndicesServiceWrapper) Unwrap() *elastic.CatIndicesService {
	return w.obj
}

type ESCatShardsServiceWrapper struct {
	obj                *elastic.CatShardsService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESCatShardsServiceWrapper) Unwrap() *elastic.CatShardsService {
	return w.obj
}

type ESClearScrollServiceWrapper struct {
	obj                *elastic.ClearScrollService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESClearScrollServiceWrapper) Unwrap() *elastic.ClearScrollService {
	return w.obj
}

func NewESClientWrapper(
	obj *elastic.Client,
	retry *micro.Retry,
	options *WrapperOptions,
	durationMetric *prometheus.HistogramVec,
	inflightMetric *prometheus.GaugeVec,
	rateLimiter micro.RateLimiter,
	parallelController micro.ParallelController) *ESClientWrapper {
	return &ESClientWrapper{
		obj:                obj,
		retry:              retry,
		options:            options,
		durationMetric:     durationMetric,
		inflightMetric:     inflightMetric,
		rateLimiter:        rateLimiter,
		parallelController: parallelController,
	}
}

type ESClientWrapper struct {
	obj                *elastic.Client
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESClientWrapper) Unwrap() *elastic.Client {
	return w.obj
}

func (w *ESClientWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *ESClientWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *ESClientWrapper) OnRateLimiterChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *ESClientWrapper) OnParallelControllerChange(opts ...refx.Option) config.OnChangeHandler {
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

func (w *ESClientWrapper) CreateMetric(options *WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        fmt.Sprintf("%s_elastic_Client_durationMs", options.Name),
		Help:        "elastic Client response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        fmt.Sprintf("%s_elastic_Client_inflight", options.Name),
		Help:        "elastic Client inflight",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "custom"})
}

type ESClusterHealthServiceWrapper struct {
	obj                *elastic.ClusterHealthService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESClusterHealthServiceWrapper) Unwrap() *elastic.ClusterHealthService {
	return w.obj
}

type ESClusterRerouteServiceWrapper struct {
	obj                *elastic.ClusterRerouteService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESClusterRerouteServiceWrapper) Unwrap() *elastic.ClusterRerouteService {
	return w.obj
}

type ESClusterStateServiceWrapper struct {
	obj                *elastic.ClusterStateService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESClusterStateServiceWrapper) Unwrap() *elastic.ClusterStateService {
	return w.obj
}

type ESClusterStatsServiceWrapper struct {
	obj                *elastic.ClusterStatsService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESClusterStatsServiceWrapper) Unwrap() *elastic.ClusterStatsService {
	return w.obj
}

type ESCountServiceWrapper struct {
	obj                *elastic.CountService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESCountServiceWrapper) Unwrap() *elastic.CountService {
	return w.obj
}

type ESDeleteByQueryServiceWrapper struct {
	obj                *elastic.DeleteByQueryService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESDeleteByQueryServiceWrapper) Unwrap() *elastic.DeleteByQueryService {
	return w.obj
}

type ESDeleteScriptServiceWrapper struct {
	obj                *elastic.DeleteScriptService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESDeleteScriptServiceWrapper) Unwrap() *elastic.DeleteScriptService {
	return w.obj
}

type ESDeleteServiceWrapper struct {
	obj                *elastic.DeleteService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESDeleteServiceWrapper) Unwrap() *elastic.DeleteService {
	return w.obj
}

type ESExistsServiceWrapper struct {
	obj                *elastic.ExistsService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESExistsServiceWrapper) Unwrap() *elastic.ExistsService {
	return w.obj
}

type ESExplainServiceWrapper struct {
	obj                *elastic.ExplainService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESExplainServiceWrapper) Unwrap() *elastic.ExplainService {
	return w.obj
}

type ESFieldCapsServiceWrapper struct {
	obj                *elastic.FieldCapsService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESFieldCapsServiceWrapper) Unwrap() *elastic.FieldCapsService {
	return w.obj
}

type ESGetScriptServiceWrapper struct {
	obj                *elastic.GetScriptService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESGetScriptServiceWrapper) Unwrap() *elastic.GetScriptService {
	return w.obj
}

type ESGetServiceWrapper struct {
	obj                *elastic.GetService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESGetServiceWrapper) Unwrap() *elastic.GetService {
	return w.obj
}

type ESIndexServiceWrapper struct {
	obj                *elastic.IndexService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndexServiceWrapper) Unwrap() *elastic.IndexService {
	return w.obj
}

type ESIndicesAnalyzeServiceWrapper struct {
	obj                *elastic.IndicesAnalyzeService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesAnalyzeServiceWrapper) Unwrap() *elastic.IndicesAnalyzeService {
	return w.obj
}

type ESIndicesClearCacheServiceWrapper struct {
	obj                *elastic.IndicesClearCacheService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesClearCacheServiceWrapper) Unwrap() *elastic.IndicesClearCacheService {
	return w.obj
}

type ESIndicesCloseServiceWrapper struct {
	obj                *elastic.IndicesCloseService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesCloseServiceWrapper) Unwrap() *elastic.IndicesCloseService {
	return w.obj
}

type ESIndicesCreateServiceWrapper struct {
	obj                *elastic.IndicesCreateService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesCreateServiceWrapper) Unwrap() *elastic.IndicesCreateService {
	return w.obj
}

type ESIndicesDeleteIndexTemplateServiceWrapper struct {
	obj                *elastic.IndicesDeleteIndexTemplateService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesDeleteIndexTemplateServiceWrapper) Unwrap() *elastic.IndicesDeleteIndexTemplateService {
	return w.obj
}

type ESIndicesDeleteServiceWrapper struct {
	obj                *elastic.IndicesDeleteService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesDeleteServiceWrapper) Unwrap() *elastic.IndicesDeleteService {
	return w.obj
}

type ESIndicesDeleteTemplateServiceWrapper struct {
	obj                *elastic.IndicesDeleteTemplateService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesDeleteTemplateServiceWrapper) Unwrap() *elastic.IndicesDeleteTemplateService {
	return w.obj
}

type ESIndicesExistsServiceWrapper struct {
	obj                *elastic.IndicesExistsService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesExistsServiceWrapper) Unwrap() *elastic.IndicesExistsService {
	return w.obj
}

type ESIndicesExistsTemplateServiceWrapper struct {
	obj                *elastic.IndicesExistsTemplateService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesExistsTemplateServiceWrapper) Unwrap() *elastic.IndicesExistsTemplateService {
	return w.obj
}

type ESIndicesFlushServiceWrapper struct {
	obj                *elastic.IndicesFlushService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesFlushServiceWrapper) Unwrap() *elastic.IndicesFlushService {
	return w.obj
}

type ESIndicesForcemergeServiceWrapper struct {
	obj                *elastic.IndicesForcemergeService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesForcemergeServiceWrapper) Unwrap() *elastic.IndicesForcemergeService {
	return w.obj
}

type ESIndicesFreezeServiceWrapper struct {
	obj                *elastic.IndicesFreezeService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesFreezeServiceWrapper) Unwrap() *elastic.IndicesFreezeService {
	return w.obj
}

type ESIndicesGetFieldMappingServiceWrapper struct {
	obj                *elastic.IndicesGetFieldMappingService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesGetFieldMappingServiceWrapper) Unwrap() *elastic.IndicesGetFieldMappingService {
	return w.obj
}

type ESIndicesGetIndexTemplateServiceWrapper struct {
	obj                *elastic.IndicesGetIndexTemplateService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesGetIndexTemplateServiceWrapper) Unwrap() *elastic.IndicesGetIndexTemplateService {
	return w.obj
}

type ESIndicesGetMappingServiceWrapper struct {
	obj                *elastic.IndicesGetMappingService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesGetMappingServiceWrapper) Unwrap() *elastic.IndicesGetMappingService {
	return w.obj
}

type ESIndicesGetServiceWrapper struct {
	obj                *elastic.IndicesGetService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesGetServiceWrapper) Unwrap() *elastic.IndicesGetService {
	return w.obj
}

type ESIndicesGetSettingsServiceWrapper struct {
	obj                *elastic.IndicesGetSettingsService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesGetSettingsServiceWrapper) Unwrap() *elastic.IndicesGetSettingsService {
	return w.obj
}

type ESIndicesGetTemplateServiceWrapper struct {
	obj                *elastic.IndicesGetTemplateService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesGetTemplateServiceWrapper) Unwrap() *elastic.IndicesGetTemplateService {
	return w.obj
}

type ESIndicesOpenServiceWrapper struct {
	obj                *elastic.IndicesOpenService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesOpenServiceWrapper) Unwrap() *elastic.IndicesOpenService {
	return w.obj
}

type ESIndicesPutIndexTemplateServiceWrapper struct {
	obj                *elastic.IndicesPutIndexTemplateService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) Unwrap() *elastic.IndicesPutIndexTemplateService {
	return w.obj
}

type ESIndicesPutMappingServiceWrapper struct {
	obj                *elastic.IndicesPutMappingService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesPutMappingServiceWrapper) Unwrap() *elastic.IndicesPutMappingService {
	return w.obj
}

type ESIndicesPutSettingsServiceWrapper struct {
	obj                *elastic.IndicesPutSettingsService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesPutSettingsServiceWrapper) Unwrap() *elastic.IndicesPutSettingsService {
	return w.obj
}

type ESIndicesPutTemplateServiceWrapper struct {
	obj                *elastic.IndicesPutTemplateService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesPutTemplateServiceWrapper) Unwrap() *elastic.IndicesPutTemplateService {
	return w.obj
}

type ESIndicesRolloverServiceWrapper struct {
	obj                *elastic.IndicesRolloverService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesRolloverServiceWrapper) Unwrap() *elastic.IndicesRolloverService {
	return w.obj
}

type ESIndicesSegmentsServiceWrapper struct {
	obj                *elastic.IndicesSegmentsService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesSegmentsServiceWrapper) Unwrap() *elastic.IndicesSegmentsService {
	return w.obj
}

type ESIndicesShrinkServiceWrapper struct {
	obj                *elastic.IndicesShrinkService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesShrinkServiceWrapper) Unwrap() *elastic.IndicesShrinkService {
	return w.obj
}

type ESIndicesStatsServiceWrapper struct {
	obj                *elastic.IndicesStatsService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesStatsServiceWrapper) Unwrap() *elastic.IndicesStatsService {
	return w.obj
}

type ESIndicesSyncedFlushServiceWrapper struct {
	obj                *elastic.IndicesSyncedFlushService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesSyncedFlushServiceWrapper) Unwrap() *elastic.IndicesSyncedFlushService {
	return w.obj
}

type ESIndicesUnfreezeServiceWrapper struct {
	obj                *elastic.IndicesUnfreezeService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIndicesUnfreezeServiceWrapper) Unwrap() *elastic.IndicesUnfreezeService {
	return w.obj
}

type ESIngestDeletePipelineServiceWrapper struct {
	obj                *elastic.IngestDeletePipelineService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIngestDeletePipelineServiceWrapper) Unwrap() *elastic.IngestDeletePipelineService {
	return w.obj
}

type ESIngestGetPipelineServiceWrapper struct {
	obj                *elastic.IngestGetPipelineService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIngestGetPipelineServiceWrapper) Unwrap() *elastic.IngestGetPipelineService {
	return w.obj
}

type ESIngestPutPipelineServiceWrapper struct {
	obj                *elastic.IngestPutPipelineService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIngestPutPipelineServiceWrapper) Unwrap() *elastic.IngestPutPipelineService {
	return w.obj
}

type ESIngestSimulatePipelineServiceWrapper struct {
	obj                *elastic.IngestSimulatePipelineService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESIngestSimulatePipelineServiceWrapper) Unwrap() *elastic.IngestSimulatePipelineService {
	return w.obj
}

type ESMgetServiceWrapper struct {
	obj                *elastic.MgetService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESMgetServiceWrapper) Unwrap() *elastic.MgetService {
	return w.obj
}

type ESMultiSearchServiceWrapper struct {
	obj                *elastic.MultiSearchService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESMultiSearchServiceWrapper) Unwrap() *elastic.MultiSearchService {
	return w.obj
}

type ESMultiTermvectorServiceWrapper struct {
	obj                *elastic.MultiTermvectorService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESMultiTermvectorServiceWrapper) Unwrap() *elastic.MultiTermvectorService {
	return w.obj
}

type ESNodesInfoServiceWrapper struct {
	obj                *elastic.NodesInfoService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESNodesInfoServiceWrapper) Unwrap() *elastic.NodesInfoService {
	return w.obj
}

type ESNodesStatsServiceWrapper struct {
	obj                *elastic.NodesStatsService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESNodesStatsServiceWrapper) Unwrap() *elastic.NodesStatsService {
	return w.obj
}

type ESPingServiceWrapper struct {
	obj                *elastic.PingService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESPingServiceWrapper) Unwrap() *elastic.PingService {
	return w.obj
}

type ESPutScriptServiceWrapper struct {
	obj                *elastic.PutScriptService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESPutScriptServiceWrapper) Unwrap() *elastic.PutScriptService {
	return w.obj
}

type ESRefreshServiceWrapper struct {
	obj                *elastic.RefreshService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESRefreshServiceWrapper) Unwrap() *elastic.RefreshService {
	return w.obj
}

type ESReindexServiceWrapper struct {
	obj                *elastic.ReindexService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESReindexServiceWrapper) Unwrap() *elastic.ReindexService {
	return w.obj
}

type ESScrollServiceWrapper struct {
	obj                *elastic.ScrollService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESScrollServiceWrapper) Unwrap() *elastic.ScrollService {
	return w.obj
}

type ESSearchServiceWrapper struct {
	obj                *elastic.SearchService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESSearchServiceWrapper) Unwrap() *elastic.SearchService {
	return w.obj
}

type ESSearchShardsServiceWrapper struct {
	obj                *elastic.SearchShardsService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESSearchShardsServiceWrapper) Unwrap() *elastic.SearchShardsService {
	return w.obj
}

type ESSnapshotCreateRepositoryServiceWrapper struct {
	obj                *elastic.SnapshotCreateRepositoryService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) Unwrap() *elastic.SnapshotCreateRepositoryService {
	return w.obj
}

type ESSnapshotCreateServiceWrapper struct {
	obj                *elastic.SnapshotCreateService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESSnapshotCreateServiceWrapper) Unwrap() *elastic.SnapshotCreateService {
	return w.obj
}

type ESSnapshotDeleteRepositoryServiceWrapper struct {
	obj                *elastic.SnapshotDeleteRepositoryService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESSnapshotDeleteRepositoryServiceWrapper) Unwrap() *elastic.SnapshotDeleteRepositoryService {
	return w.obj
}

type ESSnapshotDeleteServiceWrapper struct {
	obj                *elastic.SnapshotDeleteService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESSnapshotDeleteServiceWrapper) Unwrap() *elastic.SnapshotDeleteService {
	return w.obj
}

type ESSnapshotGetRepositoryServiceWrapper struct {
	obj                *elastic.SnapshotGetRepositoryService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESSnapshotGetRepositoryServiceWrapper) Unwrap() *elastic.SnapshotGetRepositoryService {
	return w.obj
}

type ESSnapshotGetServiceWrapper struct {
	obj                *elastic.SnapshotGetService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESSnapshotGetServiceWrapper) Unwrap() *elastic.SnapshotGetService {
	return w.obj
}

type ESSnapshotRestoreServiceWrapper struct {
	obj                *elastic.SnapshotRestoreService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESSnapshotRestoreServiceWrapper) Unwrap() *elastic.SnapshotRestoreService {
	return w.obj
}

type ESSnapshotStatusServiceWrapper struct {
	obj                *elastic.SnapshotStatusService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESSnapshotStatusServiceWrapper) Unwrap() *elastic.SnapshotStatusService {
	return w.obj
}

type ESSnapshotVerifyRepositoryServiceWrapper struct {
	obj                *elastic.SnapshotVerifyRepositoryService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESSnapshotVerifyRepositoryServiceWrapper) Unwrap() *elastic.SnapshotVerifyRepositoryService {
	return w.obj
}

type ESTasksCancelServiceWrapper struct {
	obj                *elastic.TasksCancelService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESTasksCancelServiceWrapper) Unwrap() *elastic.TasksCancelService {
	return w.obj
}

type ESTasksGetTaskServiceWrapper struct {
	obj                *elastic.TasksGetTaskService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESTasksGetTaskServiceWrapper) Unwrap() *elastic.TasksGetTaskService {
	return w.obj
}

type ESTasksListServiceWrapper struct {
	obj                *elastic.TasksListService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESTasksListServiceWrapper) Unwrap() *elastic.TasksListService {
	return w.obj
}

type ESTermvectorsServiceWrapper struct {
	obj                *elastic.TermvectorsService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESTermvectorsServiceWrapper) Unwrap() *elastic.TermvectorsService {
	return w.obj
}

type ESUpdateByQueryServiceWrapper struct {
	obj                *elastic.UpdateByQueryService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESUpdateByQueryServiceWrapper) Unwrap() *elastic.UpdateByQueryService {
	return w.obj
}

type ESUpdateServiceWrapper struct {
	obj                *elastic.UpdateService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESUpdateServiceWrapper) Unwrap() *elastic.UpdateService {
	return w.obj
}

type ESValidateServiceWrapper struct {
	obj                *elastic.ValidateService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESValidateServiceWrapper) Unwrap() *elastic.ValidateService {
	return w.obj
}

type ESXPackIlmDeleteLifecycleServiceWrapper struct {
	obj                *elastic.XPackIlmDeleteLifecycleService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackIlmDeleteLifecycleServiceWrapper) Unwrap() *elastic.XPackIlmDeleteLifecycleService {
	return w.obj
}

type ESXPackIlmGetLifecycleServiceWrapper struct {
	obj                *elastic.XPackIlmGetLifecycleService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackIlmGetLifecycleServiceWrapper) Unwrap() *elastic.XPackIlmGetLifecycleService {
	return w.obj
}

type ESXPackIlmPutLifecycleServiceWrapper struct {
	obj                *elastic.XPackIlmPutLifecycleService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) Unwrap() *elastic.XPackIlmPutLifecycleService {
	return w.obj
}

type ESXPackInfoServiceWrapper struct {
	obj                *elastic.XPackInfoService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackInfoServiceWrapper) Unwrap() *elastic.XPackInfoService {
	return w.obj
}

type ESXPackSecurityChangePasswordServiceWrapper struct {
	obj                *elastic.XPackSecurityChangePasswordService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackSecurityChangePasswordServiceWrapper) Unwrap() *elastic.XPackSecurityChangePasswordService {
	return w.obj
}

type ESXPackSecurityDeleteRoleMappingServiceWrapper struct {
	obj                *elastic.XPackSecurityDeleteRoleMappingService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackSecurityDeleteRoleMappingServiceWrapper) Unwrap() *elastic.XPackSecurityDeleteRoleMappingService {
	return w.obj
}

type ESXPackSecurityDeleteRoleServiceWrapper struct {
	obj                *elastic.XPackSecurityDeleteRoleService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackSecurityDeleteRoleServiceWrapper) Unwrap() *elastic.XPackSecurityDeleteRoleService {
	return w.obj
}

type ESXPackSecurityDeleteUserServiceWrapper struct {
	obj                *elastic.XPackSecurityDeleteUserService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackSecurityDeleteUserServiceWrapper) Unwrap() *elastic.XPackSecurityDeleteUserService {
	return w.obj
}

type ESXPackSecurityDisableUserServiceWrapper struct {
	obj                *elastic.XPackSecurityDisableUserService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackSecurityDisableUserServiceWrapper) Unwrap() *elastic.XPackSecurityDisableUserService {
	return w.obj
}

type ESXPackSecurityEnableUserServiceWrapper struct {
	obj                *elastic.XPackSecurityEnableUserService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackSecurityEnableUserServiceWrapper) Unwrap() *elastic.XPackSecurityEnableUserService {
	return w.obj
}

type ESXPackSecurityGetRoleMappingServiceWrapper struct {
	obj                *elastic.XPackSecurityGetRoleMappingService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackSecurityGetRoleMappingServiceWrapper) Unwrap() *elastic.XPackSecurityGetRoleMappingService {
	return w.obj
}

type ESXPackSecurityGetRoleServiceWrapper struct {
	obj                *elastic.XPackSecurityGetRoleService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackSecurityGetRoleServiceWrapper) Unwrap() *elastic.XPackSecurityGetRoleService {
	return w.obj
}

type ESXPackSecurityGetUserServiceWrapper struct {
	obj                *elastic.XPackSecurityGetUserService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackSecurityGetUserServiceWrapper) Unwrap() *elastic.XPackSecurityGetUserService {
	return w.obj
}

type ESXPackSecurityPutRoleMappingServiceWrapper struct {
	obj                *elastic.XPackSecurityPutRoleMappingService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackSecurityPutRoleMappingServiceWrapper) Unwrap() *elastic.XPackSecurityPutRoleMappingService {
	return w.obj
}

type ESXPackSecurityPutRoleServiceWrapper struct {
	obj                *elastic.XPackSecurityPutRoleService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackSecurityPutRoleServiceWrapper) Unwrap() *elastic.XPackSecurityPutRoleService {
	return w.obj
}

type ESXPackSecurityPutUserServiceWrapper struct {
	obj                *elastic.XPackSecurityPutUserService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackSecurityPutUserServiceWrapper) Unwrap() *elastic.XPackSecurityPutUserService {
	return w.obj
}

type ESXPackWatcherAckWatchServiceWrapper struct {
	obj                *elastic.XPackWatcherAckWatchService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackWatcherAckWatchServiceWrapper) Unwrap() *elastic.XPackWatcherAckWatchService {
	return w.obj
}

type ESXPackWatcherActivateWatchServiceWrapper struct {
	obj                *elastic.XPackWatcherActivateWatchService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackWatcherActivateWatchServiceWrapper) Unwrap() *elastic.XPackWatcherActivateWatchService {
	return w.obj
}

type ESXPackWatcherDeactivateWatchServiceWrapper struct {
	obj                *elastic.XPackWatcherDeactivateWatchService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackWatcherDeactivateWatchServiceWrapper) Unwrap() *elastic.XPackWatcherDeactivateWatchService {
	return w.obj
}

type ESXPackWatcherDeleteWatchServiceWrapper struct {
	obj                *elastic.XPackWatcherDeleteWatchService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackWatcherDeleteWatchServiceWrapper) Unwrap() *elastic.XPackWatcherDeleteWatchService {
	return w.obj
}

type ESXPackWatcherExecuteWatchServiceWrapper struct {
	obj                *elastic.XPackWatcherExecuteWatchService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackWatcherExecuteWatchServiceWrapper) Unwrap() *elastic.XPackWatcherExecuteWatchService {
	return w.obj
}

type ESXPackWatcherGetWatchServiceWrapper struct {
	obj                *elastic.XPackWatcherGetWatchService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackWatcherGetWatchServiceWrapper) Unwrap() *elastic.XPackWatcherGetWatchService {
	return w.obj
}

type ESXPackWatcherPutWatchServiceWrapper struct {
	obj                *elastic.XPackWatcherPutWatchService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackWatcherPutWatchServiceWrapper) Unwrap() *elastic.XPackWatcherPutWatchService {
	return w.obj
}

type ESXPackWatcherStartServiceWrapper struct {
	obj                *elastic.XPackWatcherStartService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackWatcherStartServiceWrapper) Unwrap() *elastic.XPackWatcherStartService {
	return w.obj
}

type ESXPackWatcherStatsServiceWrapper struct {
	obj                *elastic.XPackWatcherStatsService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackWatcherStatsServiceWrapper) Unwrap() *elastic.XPackWatcherStatsService {
	return w.obj
}

type ESXPackWatcherStopServiceWrapper struct {
	obj                *elastic.XPackWatcherStopService
	retry              *micro.Retry
	options            *WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

func (w *ESXPackWatcherStopServiceWrapper) Unwrap() *elastic.XPackWatcherStopService {
	return w.obj
}

func (w *ESAliasServiceWrapper) Action(action ...elastic.AliasAction) *ESAliasServiceWrapper {
	w.obj = w.obj.Action(action...)
	return w
}

func (w *ESAliasServiceWrapper) Add(indexName string, aliasName string) *ESAliasServiceWrapper {
	w.obj = w.obj.Add(indexName, aliasName)
	return w
}

func (w *ESAliasServiceWrapper) AddWithFilter(indexName string, aliasName string, filter elastic.Query) *ESAliasServiceWrapper {
	w.obj = w.obj.AddWithFilter(indexName, aliasName, filter)
	return w
}

func (w *ESAliasServiceWrapper) Do(ctx context.Context) (*elastic.AliasResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.AliasResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.AliasService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.AliasService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.AliasService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.AliasService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.AliasService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.AliasService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.AliasService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESAliasServiceWrapper) ErrorTrace(errorTrace bool) *ESAliasServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESAliasServiceWrapper) FilterPath(filterPath ...string) *ESAliasServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESAliasServiceWrapper) Header(name string, value string) *ESAliasServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESAliasServiceWrapper) Headers(headers http.Header) *ESAliasServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESAliasServiceWrapper) Human(human bool) *ESAliasServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESAliasServiceWrapper) Pretty(pretty bool) *ESAliasServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESAliasServiceWrapper) Remove(indexName string, aliasName string) *ESAliasServiceWrapper {
	w.obj = w.obj.Remove(indexName, aliasName)
	return w
}

func (w *ESAliasesServiceWrapper) Alias(alias ...string) *ESAliasesServiceWrapper {
	w.obj = w.obj.Alias(alias...)
	return w
}

func (w *ESAliasesServiceWrapper) Do(ctx context.Context) (*elastic.AliasesResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.AliasesResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.AliasesService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.AliasesService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.AliasesService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.AliasesService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.AliasesService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.AliasesService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.AliasesService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESAliasesServiceWrapper) ErrorTrace(errorTrace bool) *ESAliasesServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESAliasesServiceWrapper) FilterPath(filterPath ...string) *ESAliasesServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESAliasesServiceWrapper) Header(name string, value string) *ESAliasesServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESAliasesServiceWrapper) Headers(headers http.Header) *ESAliasesServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESAliasesServiceWrapper) Human(human bool) *ESAliasesServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESAliasesServiceWrapper) Index(index ...string) *ESAliasesServiceWrapper {
	w.obj = w.obj.Index(index...)
	return w
}

func (w *ESAliasesServiceWrapper) Pretty(pretty bool) *ESAliasesServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESBulkProcessorServiceWrapper) After(fn elastic.BulkAfterFunc) *ESBulkProcessorServiceWrapper {
	w.obj = w.obj.After(fn)
	return w
}

func (w *ESBulkProcessorServiceWrapper) Backoff(backoff elastic.Backoff) *ESBulkProcessorServiceWrapper {
	w.obj = w.obj.Backoff(backoff)
	return w
}

func (w *ESBulkProcessorServiceWrapper) Before(fn elastic.BulkBeforeFunc) *ESBulkProcessorServiceWrapper {
	w.obj = w.obj.Before(fn)
	return w
}

func (w *ESBulkProcessorServiceWrapper) BulkActions(bulkActions int) *ESBulkProcessorServiceWrapper {
	w.obj = w.obj.BulkActions(bulkActions)
	return w
}

func (w *ESBulkProcessorServiceWrapper) BulkSize(bulkSize int) *ESBulkProcessorServiceWrapper {
	w.obj = w.obj.BulkSize(bulkSize)
	return w
}

func (w *ESBulkProcessorServiceWrapper) Do(ctx context.Context) (*elastic.BulkProcessor, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.BulkProcessor
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.BulkProcessorService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.BulkProcessorService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.BulkProcessorService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.BulkProcessorService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.BulkProcessorService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.BulkProcessorService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.BulkProcessorService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESBulkProcessorServiceWrapper) FlushInterval(interval time.Duration) *ESBulkProcessorServiceWrapper {
	w.obj = w.obj.FlushInterval(interval)
	return w
}

func (w *ESBulkProcessorServiceWrapper) Name(name string) *ESBulkProcessorServiceWrapper {
	w.obj = w.obj.Name(name)
	return w
}

func (w *ESBulkProcessorServiceWrapper) RetryItemStatusCodes(retryItemStatusCodes ...int) *ESBulkProcessorServiceWrapper {
	w.obj = w.obj.RetryItemStatusCodes(retryItemStatusCodes...)
	return w
}

func (w *ESBulkProcessorServiceWrapper) Stats(wantStats bool) *ESBulkProcessorServiceWrapper {
	w.obj = w.obj.Stats(wantStats)
	return w
}

func (w *ESBulkProcessorServiceWrapper) Workers(num int) *ESBulkProcessorServiceWrapper {
	w.obj = w.obj.Workers(num)
	return w
}

func (w *ESBulkServiceWrapper) Add(requests ...elastic.BulkableRequest) *ESBulkServiceWrapper {
	w.obj = w.obj.Add(requests...)
	return w
}

func (w *ESBulkServiceWrapper) Do(ctx context.Context) (*elastic.BulkResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.BulkResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.BulkService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.BulkService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.BulkService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.BulkService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.BulkService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.BulkService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.BulkService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESBulkServiceWrapper) ErrorTrace(errorTrace bool) *ESBulkServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESBulkServiceWrapper) EstimatedSizeInBytes() int64 {
	res0 := w.obj.EstimatedSizeInBytes()
	return res0
}

func (w *ESBulkServiceWrapper) FilterPath(filterPath ...string) *ESBulkServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESBulkServiceWrapper) Header(name string, value string) *ESBulkServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESBulkServiceWrapper) Headers(headers http.Header) *ESBulkServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESBulkServiceWrapper) Human(human bool) *ESBulkServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESBulkServiceWrapper) Index(index string) *ESBulkServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESBulkServiceWrapper) NumberOfActions() int {
	res0 := w.obj.NumberOfActions()
	return res0
}

func (w *ESBulkServiceWrapper) Pipeline(pipeline string) *ESBulkServiceWrapper {
	w.obj = w.obj.Pipeline(pipeline)
	return w
}

func (w *ESBulkServiceWrapper) Pretty(pretty bool) *ESBulkServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESBulkServiceWrapper) Refresh(refresh string) *ESBulkServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESBulkServiceWrapper) Reset() {
	w.obj.Reset()
}

func (w *ESBulkServiceWrapper) Retrier(retrier elastic.Retrier) *ESBulkServiceWrapper {
	w.obj = w.obj.Retrier(retrier)
	return w
}

func (w *ESBulkServiceWrapper) Routing(routing string) *ESBulkServiceWrapper {
	w.obj = w.obj.Routing(routing)
	return w
}

func (w *ESBulkServiceWrapper) Timeout(timeout string) *ESBulkServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESBulkServiceWrapper) Type(typ string) *ESBulkServiceWrapper {
	w.obj = w.obj.Type(typ)
	return w
}

func (w *ESBulkServiceWrapper) WaitForActiveShards(waitForActiveShards string) *ESBulkServiceWrapper {
	w.obj = w.obj.WaitForActiveShards(waitForActiveShards)
	return w
}

func (w *ESCatAliasesServiceWrapper) Alias(alias ...string) *ESCatAliasesServiceWrapper {
	w.obj = w.obj.Alias(alias...)
	return w
}

func (w *ESCatAliasesServiceWrapper) Columns(columns ...string) *ESCatAliasesServiceWrapper {
	w.obj = w.obj.Columns(columns...)
	return w
}

func (w *ESCatAliasesServiceWrapper) Do(ctx context.Context) (elastic.CatAliasesResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 elastic.CatAliasesResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.CatAliasesService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.CatAliasesService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.CatAliasesService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.CatAliasesService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.CatAliasesService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.CatAliasesService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.CatAliasesService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESCatAliasesServiceWrapper) ErrorTrace(errorTrace bool) *ESCatAliasesServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESCatAliasesServiceWrapper) FilterPath(filterPath ...string) *ESCatAliasesServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESCatAliasesServiceWrapper) Header(name string, value string) *ESCatAliasesServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESCatAliasesServiceWrapper) Headers(headers http.Header) *ESCatAliasesServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESCatAliasesServiceWrapper) Human(human bool) *ESCatAliasesServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESCatAliasesServiceWrapper) Local(local bool) *ESCatAliasesServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESCatAliasesServiceWrapper) MasterTimeout(masterTimeout string) *ESCatAliasesServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESCatAliasesServiceWrapper) Pretty(pretty bool) *ESCatAliasesServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESCatAliasesServiceWrapper) Sort(fields ...string) *ESCatAliasesServiceWrapper {
	w.obj = w.obj.Sort(fields...)
	return w
}

func (w *ESCatAllocationServiceWrapper) Bytes(bytes string) *ESCatAllocationServiceWrapper {
	w.obj = w.obj.Bytes(bytes)
	return w
}

func (w *ESCatAllocationServiceWrapper) Columns(columns ...string) *ESCatAllocationServiceWrapper {
	w.obj = w.obj.Columns(columns...)
	return w
}

func (w *ESCatAllocationServiceWrapper) Do(ctx context.Context) (elastic.CatAllocationResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 elastic.CatAllocationResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.CatAllocationService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.CatAllocationService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.CatAllocationService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.CatAllocationService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.CatAllocationService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.CatAllocationService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.CatAllocationService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESCatAllocationServiceWrapper) ErrorTrace(errorTrace bool) *ESCatAllocationServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESCatAllocationServiceWrapper) FilterPath(filterPath ...string) *ESCatAllocationServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESCatAllocationServiceWrapper) Header(name string, value string) *ESCatAllocationServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESCatAllocationServiceWrapper) Headers(headers http.Header) *ESCatAllocationServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESCatAllocationServiceWrapper) Human(human bool) *ESCatAllocationServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESCatAllocationServiceWrapper) Local(local bool) *ESCatAllocationServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESCatAllocationServiceWrapper) MasterTimeout(masterTimeout string) *ESCatAllocationServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESCatAllocationServiceWrapper) NodeID(nodes ...string) *ESCatAllocationServiceWrapper {
	w.obj = w.obj.NodeID(nodes...)
	return w
}

func (w *ESCatAllocationServiceWrapper) Pretty(pretty bool) *ESCatAllocationServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESCatAllocationServiceWrapper) Sort(fields ...string) *ESCatAllocationServiceWrapper {
	w.obj = w.obj.Sort(fields...)
	return w
}

func (w *ESCatCountServiceWrapper) Columns(columns ...string) *ESCatCountServiceWrapper {
	w.obj = w.obj.Columns(columns...)
	return w
}

func (w *ESCatCountServiceWrapper) Do(ctx context.Context) (elastic.CatCountResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 elastic.CatCountResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.CatCountService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.CatCountService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.CatCountService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.CatCountService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.CatCountService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.CatCountService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.CatCountService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESCatCountServiceWrapper) ErrorTrace(errorTrace bool) *ESCatCountServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESCatCountServiceWrapper) FilterPath(filterPath ...string) *ESCatCountServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESCatCountServiceWrapper) Header(name string, value string) *ESCatCountServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESCatCountServiceWrapper) Headers(headers http.Header) *ESCatCountServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESCatCountServiceWrapper) Human(human bool) *ESCatCountServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESCatCountServiceWrapper) Index(index ...string) *ESCatCountServiceWrapper {
	w.obj = w.obj.Index(index...)
	return w
}

func (w *ESCatCountServiceWrapper) Local(local bool) *ESCatCountServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESCatCountServiceWrapper) MasterTimeout(masterTimeout string) *ESCatCountServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESCatCountServiceWrapper) Pretty(pretty bool) *ESCatCountServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESCatCountServiceWrapper) Sort(fields ...string) *ESCatCountServiceWrapper {
	w.obj = w.obj.Sort(fields...)
	return w
}

func (w *ESCatHealthServiceWrapper) Columns(columns ...string) *ESCatHealthServiceWrapper {
	w.obj = w.obj.Columns(columns...)
	return w
}

func (w *ESCatHealthServiceWrapper) DisableTimestamping(disable bool) *ESCatHealthServiceWrapper {
	w.obj = w.obj.DisableTimestamping(disable)
	return w
}

func (w *ESCatHealthServiceWrapper) Do(ctx context.Context) (elastic.CatHealthResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 elastic.CatHealthResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.CatHealthService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.CatHealthService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.CatHealthService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.CatHealthService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.CatHealthService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.CatHealthService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.CatHealthService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESCatHealthServiceWrapper) ErrorTrace(errorTrace bool) *ESCatHealthServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESCatHealthServiceWrapper) FilterPath(filterPath ...string) *ESCatHealthServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESCatHealthServiceWrapper) Header(name string, value string) *ESCatHealthServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESCatHealthServiceWrapper) Headers(headers http.Header) *ESCatHealthServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESCatHealthServiceWrapper) Human(human bool) *ESCatHealthServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESCatHealthServiceWrapper) Local(local bool) *ESCatHealthServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESCatHealthServiceWrapper) MasterTimeout(masterTimeout string) *ESCatHealthServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESCatHealthServiceWrapper) Pretty(pretty bool) *ESCatHealthServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESCatHealthServiceWrapper) Sort(fields ...string) *ESCatHealthServiceWrapper {
	w.obj = w.obj.Sort(fields...)
	return w
}

func (w *ESCatIndicesServiceWrapper) Bytes(bytes string) *ESCatIndicesServiceWrapper {
	w.obj = w.obj.Bytes(bytes)
	return w
}

func (w *ESCatIndicesServiceWrapper) Columns(columns ...string) *ESCatIndicesServiceWrapper {
	w.obj = w.obj.Columns(columns...)
	return w
}

func (w *ESCatIndicesServiceWrapper) Do(ctx context.Context) (elastic.CatIndicesResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 elastic.CatIndicesResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.CatIndicesService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.CatIndicesService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.CatIndicesService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.CatIndicesService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.CatIndicesService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.CatIndicesService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.CatIndicesService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESCatIndicesServiceWrapper) ErrorTrace(errorTrace bool) *ESCatIndicesServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESCatIndicesServiceWrapper) FilterPath(filterPath ...string) *ESCatIndicesServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESCatIndicesServiceWrapper) Header(name string, value string) *ESCatIndicesServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESCatIndicesServiceWrapper) Headers(headers http.Header) *ESCatIndicesServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESCatIndicesServiceWrapper) Health(healthState string) *ESCatIndicesServiceWrapper {
	w.obj = w.obj.Health(healthState)
	return w
}

func (w *ESCatIndicesServiceWrapper) Human(human bool) *ESCatIndicesServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESCatIndicesServiceWrapper) Index(index string) *ESCatIndicesServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESCatIndicesServiceWrapper) Local(local bool) *ESCatIndicesServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESCatIndicesServiceWrapper) MasterTimeout(masterTimeout string) *ESCatIndicesServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESCatIndicesServiceWrapper) Pretty(pretty bool) *ESCatIndicesServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESCatIndicesServiceWrapper) PrimaryOnly(primaryOnly bool) *ESCatIndicesServiceWrapper {
	w.obj = w.obj.PrimaryOnly(primaryOnly)
	return w
}

func (w *ESCatIndicesServiceWrapper) Sort(fields ...string) *ESCatIndicesServiceWrapper {
	w.obj = w.obj.Sort(fields...)
	return w
}

func (w *ESCatShardsServiceWrapper) Bytes(bytes string) *ESCatShardsServiceWrapper {
	w.obj = w.obj.Bytes(bytes)
	return w
}

func (w *ESCatShardsServiceWrapper) Columns(columns ...string) *ESCatShardsServiceWrapper {
	w.obj = w.obj.Columns(columns...)
	return w
}

func (w *ESCatShardsServiceWrapper) Do(ctx context.Context) (elastic.CatShardsResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 elastic.CatShardsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.CatShardsService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.CatShardsService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.CatShardsService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.CatShardsService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.CatShardsService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.CatShardsService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.CatShardsService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESCatShardsServiceWrapper) ErrorTrace(errorTrace bool) *ESCatShardsServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESCatShardsServiceWrapper) FilterPath(filterPath ...string) *ESCatShardsServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESCatShardsServiceWrapper) Header(name string, value string) *ESCatShardsServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESCatShardsServiceWrapper) Headers(headers http.Header) *ESCatShardsServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESCatShardsServiceWrapper) Human(human bool) *ESCatShardsServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESCatShardsServiceWrapper) Index(index ...string) *ESCatShardsServiceWrapper {
	w.obj = w.obj.Index(index...)
	return w
}

func (w *ESCatShardsServiceWrapper) Local(local bool) *ESCatShardsServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESCatShardsServiceWrapper) MasterTimeout(masterTimeout string) *ESCatShardsServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESCatShardsServiceWrapper) Pretty(pretty bool) *ESCatShardsServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESCatShardsServiceWrapper) Sort(fields ...string) *ESCatShardsServiceWrapper {
	w.obj = w.obj.Sort(fields...)
	return w
}

func (w *ESCatShardsServiceWrapper) Time(time string) *ESCatShardsServiceWrapper {
	w.obj = w.obj.Time(time)
	return w
}

func (w *ESClearScrollServiceWrapper) Do(ctx context.Context) (*elastic.ClearScrollResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.ClearScrollResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.ClearScrollService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.ClearScrollService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.ClearScrollService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.ClearScrollService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.ClearScrollService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.ClearScrollService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.ClearScrollService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESClearScrollServiceWrapper) ErrorTrace(errorTrace bool) *ESClearScrollServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESClearScrollServiceWrapper) FilterPath(filterPath ...string) *ESClearScrollServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESClearScrollServiceWrapper) Header(name string, value string) *ESClearScrollServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESClearScrollServiceWrapper) Headers(headers http.Header) *ESClearScrollServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESClearScrollServiceWrapper) Human(human bool) *ESClearScrollServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESClearScrollServiceWrapper) Pretty(pretty bool) *ESClearScrollServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESClearScrollServiceWrapper) ScrollId(scrollIds ...string) *ESClearScrollServiceWrapper {
	w.obj = w.obj.ScrollId(scrollIds...)
	return w
}

func (w *ESClearScrollServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESClientWrapper) Alias() *ESAliasServiceWrapper {
	res0 := w.obj.Alias()
	return &ESAliasServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Aliases() *ESAliasesServiceWrapper {
	res0 := w.obj.Aliases()
	return &ESAliasesServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Bulk() *ESBulkServiceWrapper {
	res0 := w.obj.Bulk()
	return &ESBulkServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) BulkProcessor() *ESBulkProcessorServiceWrapper {
	res0 := w.obj.BulkProcessor()
	return &ESBulkProcessorServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) CatAliases() *ESCatAliasesServiceWrapper {
	res0 := w.obj.CatAliases()
	return &ESCatAliasesServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) CatAllocation() *ESCatAllocationServiceWrapper {
	res0 := w.obj.CatAllocation()
	return &ESCatAllocationServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) CatCount() *ESCatCountServiceWrapper {
	res0 := w.obj.CatCount()
	return &ESCatCountServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) CatHealth() *ESCatHealthServiceWrapper {
	res0 := w.obj.CatHealth()
	return &ESCatHealthServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) CatIndices() *ESCatIndicesServiceWrapper {
	res0 := w.obj.CatIndices()
	return &ESCatIndicesServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) CatShards() *ESCatShardsServiceWrapper {
	res0 := w.obj.CatShards()
	return &ESCatShardsServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) ClearCache(indices ...string) *ESIndicesClearCacheServiceWrapper {
	res0 := w.obj.ClearCache(indices...)
	return &ESIndicesClearCacheServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) ClearScroll(scrollIds ...string) *ESClearScrollServiceWrapper {
	res0 := w.obj.ClearScroll(scrollIds...)
	return &ESClearScrollServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) CloseIndex(name string) *ESIndicesCloseServiceWrapper {
	res0 := w.obj.CloseIndex(name)
	return &ESIndicesCloseServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) ClusterHealth() *ESClusterHealthServiceWrapper {
	res0 := w.obj.ClusterHealth()
	return &ESClusterHealthServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) ClusterReroute() *ESClusterRerouteServiceWrapper {
	res0 := w.obj.ClusterReroute()
	return &ESClusterRerouteServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) ClusterState() *ESClusterStateServiceWrapper {
	res0 := w.obj.ClusterState()
	return &ESClusterStateServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) ClusterStats() *ESClusterStatsServiceWrapper {
	res0 := w.obj.ClusterStats()
	return &ESClusterStatsServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Count(indices ...string) *ESCountServiceWrapper {
	res0 := w.obj.Count(indices...)
	return &ESCountServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) CreateIndex(name string) *ESIndicesCreateServiceWrapper {
	res0 := w.obj.CreateIndex(name)
	return &ESIndicesCreateServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Delete() *ESDeleteServiceWrapper {
	res0 := w.obj.Delete()
	return &ESDeleteServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) DeleteByQuery(indices ...string) *ESDeleteByQueryServiceWrapper {
	res0 := w.obj.DeleteByQuery(indices...)
	return &ESDeleteByQueryServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) DeleteIndex(indices ...string) *ESIndicesDeleteServiceWrapper {
	res0 := w.obj.DeleteIndex(indices...)
	return &ESIndicesDeleteServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) DeleteScript() *ESDeleteScriptServiceWrapper {
	res0 := w.obj.DeleteScript()
	return &ESDeleteScriptServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) ElasticsearchVersion(url string) (string, error) {
	var res0 string
	var err error
	res0, err = w.obj.ElasticsearchVersion(url)
	return res0, err
}

func (w *ESClientWrapper) Exists() *ESExistsServiceWrapper {
	res0 := w.obj.Exists()
	return &ESExistsServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Explain(index string, typ string, id string) *ESExplainServiceWrapper {
	res0 := w.obj.Explain(index, typ, id)
	return &ESExplainServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) FieldCaps(indices ...string) *ESFieldCapsServiceWrapper {
	res0 := w.obj.FieldCaps(indices...)
	return &ESFieldCapsServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Flush(indices ...string) *ESIndicesFlushServiceWrapper {
	res0 := w.obj.Flush(indices...)
	return &ESIndicesFlushServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Forcemerge(indices ...string) *ESIndicesForcemergeServiceWrapper {
	res0 := w.obj.Forcemerge(indices...)
	return &ESIndicesForcemergeServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) FreezeIndex(name string) *ESIndicesFreezeServiceWrapper {
	res0 := w.obj.FreezeIndex(name)
	return &ESIndicesFreezeServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Get() *ESGetServiceWrapper {
	res0 := w.obj.Get()
	return &ESGetServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) GetFieldMapping() *ESIndicesGetFieldMappingServiceWrapper {
	res0 := w.obj.GetFieldMapping()
	return &ESIndicesGetFieldMappingServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) GetMapping() *ESIndicesGetMappingServiceWrapper {
	res0 := w.obj.GetMapping()
	return &ESIndicesGetMappingServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) GetScript() *ESGetScriptServiceWrapper {
	res0 := w.obj.GetScript()
	return &ESGetScriptServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) HasPlugin(name string) (bool, error) {
	var res0 bool
	var err error
	res0, err = w.obj.HasPlugin(name)
	return res0, err
}

func (w *ESClientWrapper) Index() *ESIndexServiceWrapper {
	res0 := w.obj.Index()
	return &ESIndexServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IndexAnalyze() *ESIndicesAnalyzeServiceWrapper {
	res0 := w.obj.IndexAnalyze()
	return &ESIndicesAnalyzeServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IndexDeleteIndexTemplate(name string) *ESIndicesDeleteIndexTemplateServiceWrapper {
	res0 := w.obj.IndexDeleteIndexTemplate(name)
	return &ESIndicesDeleteIndexTemplateServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IndexDeleteTemplate(name string) *ESIndicesDeleteTemplateServiceWrapper {
	res0 := w.obj.IndexDeleteTemplate(name)
	return &ESIndicesDeleteTemplateServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IndexExists(indices ...string) *ESIndicesExistsServiceWrapper {
	res0 := w.obj.IndexExists(indices...)
	return &ESIndicesExistsServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IndexGet(indices ...string) *ESIndicesGetServiceWrapper {
	res0 := w.obj.IndexGet(indices...)
	return &ESIndicesGetServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IndexGetIndexTemplate(name string) *ESIndicesGetIndexTemplateServiceWrapper {
	res0 := w.obj.IndexGetIndexTemplate(name)
	return &ESIndicesGetIndexTemplateServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IndexGetSettings(indices ...string) *ESIndicesGetSettingsServiceWrapper {
	res0 := w.obj.IndexGetSettings(indices...)
	return &ESIndicesGetSettingsServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IndexGetTemplate(names ...string) *ESIndicesGetTemplateServiceWrapper {
	res0 := w.obj.IndexGetTemplate(names...)
	return &ESIndicesGetTemplateServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IndexNames() ([]string, error) {
	var res0 []string
	var err error
	res0, err = w.obj.IndexNames()
	return res0, err
}

func (w *ESClientWrapper) IndexPutIndexTemplate(name string) *ESIndicesPutIndexTemplateServiceWrapper {
	res0 := w.obj.IndexPutIndexTemplate(name)
	return &ESIndicesPutIndexTemplateServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IndexPutSettings(indices ...string) *ESIndicesPutSettingsServiceWrapper {
	res0 := w.obj.IndexPutSettings(indices...)
	return &ESIndicesPutSettingsServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IndexPutTemplate(name string) *ESIndicesPutTemplateServiceWrapper {
	res0 := w.obj.IndexPutTemplate(name)
	return &ESIndicesPutTemplateServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IndexSegments(indices ...string) *ESIndicesSegmentsServiceWrapper {
	res0 := w.obj.IndexSegments(indices...)
	return &ESIndicesSegmentsServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IndexStats(indices ...string) *ESIndicesStatsServiceWrapper {
	res0 := w.obj.IndexStats(indices...)
	return &ESIndicesStatsServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IndexTemplateExists(name string) *ESIndicesExistsTemplateServiceWrapper {
	res0 := w.obj.IndexTemplateExists(name)
	return &ESIndicesExistsTemplateServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IngestDeletePipeline(id string) *ESIngestDeletePipelineServiceWrapper {
	res0 := w.obj.IngestDeletePipeline(id)
	return &ESIngestDeletePipelineServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IngestGetPipeline(ids ...string) *ESIngestGetPipelineServiceWrapper {
	res0 := w.obj.IngestGetPipeline(ids...)
	return &ESIngestGetPipelineServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IngestPutPipeline(id string) *ESIngestPutPipelineServiceWrapper {
	res0 := w.obj.IngestPutPipeline(id)
	return &ESIngestPutPipelineServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IngestSimulatePipeline() *ESIngestSimulatePipelineServiceWrapper {
	res0 := w.obj.IngestSimulatePipeline()
	return &ESIngestSimulatePipelineServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) IsRunning() bool {
	res0 := w.obj.IsRunning()
	return res0
}

func (w *ESClientWrapper) Mget() *ESMgetServiceWrapper {
	res0 := w.obj.Mget()
	return &ESMgetServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) MultiGet() *ESMgetServiceWrapper {
	res0 := w.obj.MultiGet()
	return &ESMgetServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) MultiSearch() *ESMultiSearchServiceWrapper {
	res0 := w.obj.MultiSearch()
	return &ESMultiSearchServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) MultiTermVectors() *ESMultiTermvectorServiceWrapper {
	res0 := w.obj.MultiTermVectors()
	return &ESMultiTermvectorServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) NodesInfo() *ESNodesInfoServiceWrapper {
	res0 := w.obj.NodesInfo()
	return &ESNodesInfoServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) NodesStats() *ESNodesStatsServiceWrapper {
	res0 := w.obj.NodesStats()
	return &ESNodesStatsServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) OpenIndex(name string) *ESIndicesOpenServiceWrapper {
	res0 := w.obj.OpenIndex(name)
	return &ESIndicesOpenServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) PerformRequest(ctx context.Context, opt elastic.PerformRequestOptions) (*elastic.Response, error) {
	var res0 *elastic.Response
	var err error
	res0, err = w.obj.PerformRequest(ctx, opt)
	return res0, err
}

func (w *ESClientWrapper) Ping(url string) *ESPingServiceWrapper {
	res0 := w.obj.Ping(url)
	return &ESPingServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Plugins() ([]string, error) {
	var res0 []string
	var err error
	res0, err = w.obj.Plugins()
	return res0, err
}

func (w *ESClientWrapper) PutMapping() *ESIndicesPutMappingServiceWrapper {
	res0 := w.obj.PutMapping()
	return &ESIndicesPutMappingServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) PutScript() *ESPutScriptServiceWrapper {
	res0 := w.obj.PutScript()
	return &ESPutScriptServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Refresh(indices ...string) *ESRefreshServiceWrapper {
	res0 := w.obj.Refresh(indices...)
	return &ESRefreshServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Reindex() *ESReindexServiceWrapper {
	res0 := w.obj.Reindex()
	return &ESReindexServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) RolloverIndex(alias string) *ESIndicesRolloverServiceWrapper {
	res0 := w.obj.RolloverIndex(alias)
	return &ESIndicesRolloverServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Scroll(indices ...string) *ESScrollServiceWrapper {
	res0 := w.obj.Scroll(indices...)
	return &ESScrollServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Search(indices ...string) *ESSearchServiceWrapper {
	res0 := w.obj.Search(indices...)
	return &ESSearchServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) SearchShards(indices ...string) *ESSearchShardsServiceWrapper {
	res0 := w.obj.SearchShards(indices...)
	return &ESSearchShardsServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) ShrinkIndex(source string, target string) *ESIndicesShrinkServiceWrapper {
	res0 := w.obj.ShrinkIndex(source, target)
	return &ESIndicesShrinkServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) SnapshotCreate(repository string, snapshot string) *ESSnapshotCreateServiceWrapper {
	res0 := w.obj.SnapshotCreate(repository, snapshot)
	return &ESSnapshotCreateServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) SnapshotCreateRepository(repository string) *ESSnapshotCreateRepositoryServiceWrapper {
	res0 := w.obj.SnapshotCreateRepository(repository)
	return &ESSnapshotCreateRepositoryServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) SnapshotDelete(repository string, snapshot string) *ESSnapshotDeleteServiceWrapper {
	res0 := w.obj.SnapshotDelete(repository, snapshot)
	return &ESSnapshotDeleteServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) SnapshotDeleteRepository(repositories ...string) *ESSnapshotDeleteRepositoryServiceWrapper {
	res0 := w.obj.SnapshotDeleteRepository(repositories...)
	return &ESSnapshotDeleteRepositoryServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) SnapshotGet(repository string) *ESSnapshotGetServiceWrapper {
	res0 := w.obj.SnapshotGet(repository)
	return &ESSnapshotGetServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) SnapshotGetRepository(repositories ...string) *ESSnapshotGetRepositoryServiceWrapper {
	res0 := w.obj.SnapshotGetRepository(repositories...)
	return &ESSnapshotGetRepositoryServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) SnapshotRestore(repository string, snapshot string) *ESSnapshotRestoreServiceWrapper {
	res0 := w.obj.SnapshotRestore(repository, snapshot)
	return &ESSnapshotRestoreServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) SnapshotStatus() *ESSnapshotStatusServiceWrapper {
	res0 := w.obj.SnapshotStatus()
	return &ESSnapshotStatusServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) SnapshotVerifyRepository(repository string) *ESSnapshotVerifyRepositoryServiceWrapper {
	res0 := w.obj.SnapshotVerifyRepository(repository)
	return &ESSnapshotVerifyRepositoryServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Start() {
	w.obj.Start()
}

func (w *ESClientWrapper) Stop() {
	w.obj.Stop()
}

func (w *ESClientWrapper) String() string {
	res0 := w.obj.String()
	return res0
}

func (w *ESClientWrapper) SyncedFlush(indices ...string) *ESIndicesSyncedFlushServiceWrapper {
	res0 := w.obj.SyncedFlush(indices...)
	return &ESIndicesSyncedFlushServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) TasksCancel() *ESTasksCancelServiceWrapper {
	res0 := w.obj.TasksCancel()
	return &ESTasksCancelServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) TasksGetTask() *ESTasksGetTaskServiceWrapper {
	res0 := w.obj.TasksGetTask()
	return &ESTasksGetTaskServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) TasksList() *ESTasksListServiceWrapper {
	res0 := w.obj.TasksList()
	return &ESTasksListServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) TermVectors(index string) *ESTermvectorsServiceWrapper {
	res0 := w.obj.TermVectors(index)
	return &ESTermvectorsServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) UnfreezeIndex(name string) *ESIndicesUnfreezeServiceWrapper {
	res0 := w.obj.UnfreezeIndex(name)
	return &ESIndicesUnfreezeServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Update() *ESUpdateServiceWrapper {
	res0 := w.obj.Update()
	return &ESUpdateServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) UpdateByQuery(indices ...string) *ESUpdateByQueryServiceWrapper {
	res0 := w.obj.UpdateByQuery(indices...)
	return &ESUpdateByQueryServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) Validate(indices ...string) *ESValidateServiceWrapper {
	res0 := w.obj.Validate(indices...)
	return &ESValidateServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) WaitForGreenStatus(timeout string) error {
	var err error
	err = w.obj.WaitForGreenStatus(timeout)
	return err
}

func (w *ESClientWrapper) WaitForStatus(status string, timeout string) error {
	var err error
	err = w.obj.WaitForStatus(status, timeout)
	return err
}

func (w *ESClientWrapper) WaitForYellowStatus(timeout string) error {
	var err error
	err = w.obj.WaitForYellowStatus(timeout)
	return err
}

func (w *ESClientWrapper) XPackAsyncSearchDelete() *elastic.XPackAsyncSearchDelete {
	res0 := w.obj.XPackAsyncSearchDelete()
	return res0
}

func (w *ESClientWrapper) XPackAsyncSearchGet() *elastic.XPackAsyncSearchGet {
	res0 := w.obj.XPackAsyncSearchGet()
	return res0
}

func (w *ESClientWrapper) XPackAsyncSearchSubmit() *elastic.XPackAsyncSearchSubmit {
	res0 := w.obj.XPackAsyncSearchSubmit()
	return res0
}

func (w *ESClientWrapper) XPackIlmDeleteLifecycle() *ESXPackIlmDeleteLifecycleServiceWrapper {
	res0 := w.obj.XPackIlmDeleteLifecycle()
	return &ESXPackIlmDeleteLifecycleServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackIlmGetLifecycle() *ESXPackIlmGetLifecycleServiceWrapper {
	res0 := w.obj.XPackIlmGetLifecycle()
	return &ESXPackIlmGetLifecycleServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackIlmPutLifecycle() *ESXPackIlmPutLifecycleServiceWrapper {
	res0 := w.obj.XPackIlmPutLifecycle()
	return &ESXPackIlmPutLifecycleServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackInfo() *ESXPackInfoServiceWrapper {
	res0 := w.obj.XPackInfo()
	return &ESXPackInfoServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackSecurityChangePassword(username string) *ESXPackSecurityChangePasswordServiceWrapper {
	res0 := w.obj.XPackSecurityChangePassword(username)
	return &ESXPackSecurityChangePasswordServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackSecurityDeleteRole(roleName string) *ESXPackSecurityDeleteRoleServiceWrapper {
	res0 := w.obj.XPackSecurityDeleteRole(roleName)
	return &ESXPackSecurityDeleteRoleServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackSecurityDeleteRoleMapping(roleMappingName string) *ESXPackSecurityDeleteRoleMappingServiceWrapper {
	res0 := w.obj.XPackSecurityDeleteRoleMapping(roleMappingName)
	return &ESXPackSecurityDeleteRoleMappingServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackSecurityDeleteUser(username string) *ESXPackSecurityDeleteUserServiceWrapper {
	res0 := w.obj.XPackSecurityDeleteUser(username)
	return &ESXPackSecurityDeleteUserServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackSecurityDisableUser(username string) *ESXPackSecurityDisableUserServiceWrapper {
	res0 := w.obj.XPackSecurityDisableUser(username)
	return &ESXPackSecurityDisableUserServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackSecurityEnableUser(username string) *ESXPackSecurityEnableUserServiceWrapper {
	res0 := w.obj.XPackSecurityEnableUser(username)
	return &ESXPackSecurityEnableUserServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackSecurityGetRole(roleName string) *ESXPackSecurityGetRoleServiceWrapper {
	res0 := w.obj.XPackSecurityGetRole(roleName)
	return &ESXPackSecurityGetRoleServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackSecurityGetRoleMapping(roleMappingName string) *ESXPackSecurityGetRoleMappingServiceWrapper {
	res0 := w.obj.XPackSecurityGetRoleMapping(roleMappingName)
	return &ESXPackSecurityGetRoleMappingServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackSecurityGetUser(usernames ...string) *ESXPackSecurityGetUserServiceWrapper {
	res0 := w.obj.XPackSecurityGetUser(usernames...)
	return &ESXPackSecurityGetUserServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackSecurityPutRole(roleName string) *ESXPackSecurityPutRoleServiceWrapper {
	res0 := w.obj.XPackSecurityPutRole(roleName)
	return &ESXPackSecurityPutRoleServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackSecurityPutRoleMapping(roleMappingName string) *ESXPackSecurityPutRoleMappingServiceWrapper {
	res0 := w.obj.XPackSecurityPutRoleMapping(roleMappingName)
	return &ESXPackSecurityPutRoleMappingServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackSecurityPutUser(username string) *ESXPackSecurityPutUserServiceWrapper {
	res0 := w.obj.XPackSecurityPutUser(username)
	return &ESXPackSecurityPutUserServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackWatchAck(watchId string) *ESXPackWatcherAckWatchServiceWrapper {
	res0 := w.obj.XPackWatchAck(watchId)
	return &ESXPackWatcherAckWatchServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackWatchActivate(watchId string) *ESXPackWatcherActivateWatchServiceWrapper {
	res0 := w.obj.XPackWatchActivate(watchId)
	return &ESXPackWatcherActivateWatchServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackWatchDeactivate(watchId string) *ESXPackWatcherDeactivateWatchServiceWrapper {
	res0 := w.obj.XPackWatchDeactivate(watchId)
	return &ESXPackWatcherDeactivateWatchServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackWatchDelete(watchId string) *ESXPackWatcherDeleteWatchServiceWrapper {
	res0 := w.obj.XPackWatchDelete(watchId)
	return &ESXPackWatcherDeleteWatchServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackWatchExecute() *ESXPackWatcherExecuteWatchServiceWrapper {
	res0 := w.obj.XPackWatchExecute()
	return &ESXPackWatcherExecuteWatchServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackWatchGet(watchId string) *ESXPackWatcherGetWatchServiceWrapper {
	res0 := w.obj.XPackWatchGet(watchId)
	return &ESXPackWatcherGetWatchServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackWatchPut(watchId string) *ESXPackWatcherPutWatchServiceWrapper {
	res0 := w.obj.XPackWatchPut(watchId)
	return &ESXPackWatcherPutWatchServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackWatchStart() *ESXPackWatcherStartServiceWrapper {
	res0 := w.obj.XPackWatchStart()
	return &ESXPackWatcherStartServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackWatchStats() *ESXPackWatcherStatsServiceWrapper {
	res0 := w.obj.XPackWatchStats()
	return &ESXPackWatcherStatsServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClientWrapper) XPackWatchStop() *ESXPackWatcherStopServiceWrapper {
	res0 := w.obj.XPackWatchStop()
	return &ESXPackWatcherStopServiceWrapper{obj: res0, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}
}

func (w *ESClusterHealthServiceWrapper) Do(ctx context.Context) (*elastic.ClusterHealthResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.ClusterHealthResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.ClusterHealthService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.ClusterHealthService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.ClusterHealthService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.ClusterHealthService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.ClusterHealthService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.ClusterHealthService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.ClusterHealthService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESClusterHealthServiceWrapper) ErrorTrace(errorTrace bool) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESClusterHealthServiceWrapper) FilterPath(filterPath ...string) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESClusterHealthServiceWrapper) Header(name string, value string) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESClusterHealthServiceWrapper) Headers(headers http.Header) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESClusterHealthServiceWrapper) Human(human bool) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESClusterHealthServiceWrapper) Index(indices ...string) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESClusterHealthServiceWrapper) Level(level string) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.Level(level)
	return w
}

func (w *ESClusterHealthServiceWrapper) Local(local bool) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESClusterHealthServiceWrapper) MasterTimeout(masterTimeout string) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESClusterHealthServiceWrapper) Pretty(pretty bool) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESClusterHealthServiceWrapper) Timeout(timeout string) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESClusterHealthServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESClusterHealthServiceWrapper) WaitForActiveShards(waitForActiveShards int) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.WaitForActiveShards(waitForActiveShards)
	return w
}

func (w *ESClusterHealthServiceWrapper) WaitForGreenStatus() *ESClusterHealthServiceWrapper {
	w.obj = w.obj.WaitForGreenStatus()
	return w
}

func (w *ESClusterHealthServiceWrapper) WaitForNoRelocatingShards(waitForNoRelocatingShards bool) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.WaitForNoRelocatingShards(waitForNoRelocatingShards)
	return w
}

func (w *ESClusterHealthServiceWrapper) WaitForNodes(waitForNodes string) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.WaitForNodes(waitForNodes)
	return w
}

func (w *ESClusterHealthServiceWrapper) WaitForStatus(waitForStatus string) *ESClusterHealthServiceWrapper {
	w.obj = w.obj.WaitForStatus(waitForStatus)
	return w
}

func (w *ESClusterHealthServiceWrapper) WaitForYellowStatus() *ESClusterHealthServiceWrapper {
	w.obj = w.obj.WaitForYellowStatus()
	return w
}

func (w *ESClusterRerouteServiceWrapper) Add(commands ...elastic.AllocationCommand) *ESClusterRerouteServiceWrapper {
	w.obj = w.obj.Add(commands...)
	return w
}

func (w *ESClusterRerouteServiceWrapper) Body(body interface{}) *ESClusterRerouteServiceWrapper {
	w.obj = w.obj.Body(body)
	return w
}

func (w *ESClusterRerouteServiceWrapper) Do(ctx context.Context) (*elastic.ClusterRerouteResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.ClusterRerouteResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.ClusterRerouteService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.ClusterRerouteService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.ClusterRerouteService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.ClusterRerouteService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.ClusterRerouteService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.ClusterRerouteService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.ClusterRerouteService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESClusterRerouteServiceWrapper) DryRun(dryRun bool) *ESClusterRerouteServiceWrapper {
	w.obj = w.obj.DryRun(dryRun)
	return w
}

func (w *ESClusterRerouteServiceWrapper) ErrorTrace(errorTrace bool) *ESClusterRerouteServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESClusterRerouteServiceWrapper) Explain(explain bool) *ESClusterRerouteServiceWrapper {
	w.obj = w.obj.Explain(explain)
	return w
}

func (w *ESClusterRerouteServiceWrapper) FilterPath(filterPath ...string) *ESClusterRerouteServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESClusterRerouteServiceWrapper) Header(name string, value string) *ESClusterRerouteServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESClusterRerouteServiceWrapper) Headers(headers http.Header) *ESClusterRerouteServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESClusterRerouteServiceWrapper) Human(human bool) *ESClusterRerouteServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESClusterRerouteServiceWrapper) MasterTimeout(masterTimeout string) *ESClusterRerouteServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESClusterRerouteServiceWrapper) Metric(metrics ...string) *ESClusterRerouteServiceWrapper {
	w.obj = w.obj.Metric(metrics...)
	return w
}

func (w *ESClusterRerouteServiceWrapper) Pretty(pretty bool) *ESClusterRerouteServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESClusterRerouteServiceWrapper) RetryFailed(retryFailed bool) *ESClusterRerouteServiceWrapper {
	w.obj = w.obj.RetryFailed(retryFailed)
	return w
}

func (w *ESClusterRerouteServiceWrapper) Timeout(timeout string) *ESClusterRerouteServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESClusterRerouteServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESClusterStateServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESClusterStateServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESClusterStateServiceWrapper) Do(ctx context.Context) (*elastic.ClusterStateResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.ClusterStateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.ClusterStateService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.ClusterStateService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.ClusterStateService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.ClusterStateService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.ClusterStateService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.ClusterStateService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.ClusterStateService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESClusterStateServiceWrapper) ErrorTrace(errorTrace bool) *ESClusterStateServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESClusterStateServiceWrapper) ExpandWildcards(expandWildcards string) *ESClusterStateServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESClusterStateServiceWrapper) FilterPath(filterPath ...string) *ESClusterStateServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESClusterStateServiceWrapper) FlatSettings(flatSettings bool) *ESClusterStateServiceWrapper {
	w.obj = w.obj.FlatSettings(flatSettings)
	return w
}

func (w *ESClusterStateServiceWrapper) Header(name string, value string) *ESClusterStateServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESClusterStateServiceWrapper) Headers(headers http.Header) *ESClusterStateServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESClusterStateServiceWrapper) Human(human bool) *ESClusterStateServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESClusterStateServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESClusterStateServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESClusterStateServiceWrapper) Index(indices ...string) *ESClusterStateServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESClusterStateServiceWrapper) Local(local bool) *ESClusterStateServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESClusterStateServiceWrapper) MasterTimeout(masterTimeout string) *ESClusterStateServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESClusterStateServiceWrapper) Metric(metrics ...string) *ESClusterStateServiceWrapper {
	w.obj = w.obj.Metric(metrics...)
	return w
}

func (w *ESClusterStateServiceWrapper) Pretty(pretty bool) *ESClusterStateServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESClusterStateServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESClusterStatsServiceWrapper) Do(ctx context.Context) (*elastic.ClusterStatsResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.ClusterStatsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.ClusterStatsService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.ClusterStatsService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.ClusterStatsService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.ClusterStatsService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.ClusterStatsService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.ClusterStatsService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.ClusterStatsService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESClusterStatsServiceWrapper) ErrorTrace(errorTrace bool) *ESClusterStatsServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESClusterStatsServiceWrapper) FilterPath(filterPath ...string) *ESClusterStatsServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESClusterStatsServiceWrapper) FlatSettings(flatSettings bool) *ESClusterStatsServiceWrapper {
	w.obj = w.obj.FlatSettings(flatSettings)
	return w
}

func (w *ESClusterStatsServiceWrapper) Header(name string, value string) *ESClusterStatsServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESClusterStatsServiceWrapper) Headers(headers http.Header) *ESClusterStatsServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESClusterStatsServiceWrapper) Human(human bool) *ESClusterStatsServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESClusterStatsServiceWrapper) NodeId(nodeId []string) *ESClusterStatsServiceWrapper {
	w.obj = w.obj.NodeId(nodeId)
	return w
}

func (w *ESClusterStatsServiceWrapper) Pretty(pretty bool) *ESClusterStatsServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESClusterStatsServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESCountServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESCountServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESCountServiceWrapper) AnalyzeWildcard(analyzeWildcard bool) *ESCountServiceWrapper {
	w.obj = w.obj.AnalyzeWildcard(analyzeWildcard)
	return w
}

func (w *ESCountServiceWrapper) Analyzer(analyzer string) *ESCountServiceWrapper {
	w.obj = w.obj.Analyzer(analyzer)
	return w
}

func (w *ESCountServiceWrapper) BodyJson(body interface{}) *ESCountServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESCountServiceWrapper) BodyString(body string) *ESCountServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESCountServiceWrapper) DefaultOperator(defaultOperator string) *ESCountServiceWrapper {
	w.obj = w.obj.DefaultOperator(defaultOperator)
	return w
}

func (w *ESCountServiceWrapper) Df(df string) *ESCountServiceWrapper {
	w.obj = w.obj.Df(df)
	return w
}

func (w *ESCountServiceWrapper) Do(ctx context.Context) (int64, error) {
	ctxOptions := FromContext(ctx)
	var res0 int64
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.CountService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.CountService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.CountService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.CountService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.CountService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.CountService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.CountService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESCountServiceWrapper) ErrorTrace(errorTrace bool) *ESCountServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESCountServiceWrapper) ExpandWildcards(expandWildcards string) *ESCountServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESCountServiceWrapper) FilterPath(filterPath ...string) *ESCountServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESCountServiceWrapper) Header(name string, value string) *ESCountServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESCountServiceWrapper) Headers(headers http.Header) *ESCountServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESCountServiceWrapper) Human(human bool) *ESCountServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESCountServiceWrapper) IgnoreThrottled(ignoreThrottled bool) *ESCountServiceWrapper {
	w.obj = w.obj.IgnoreThrottled(ignoreThrottled)
	return w
}

func (w *ESCountServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESCountServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESCountServiceWrapper) Index(index ...string) *ESCountServiceWrapper {
	w.obj = w.obj.Index(index...)
	return w
}

func (w *ESCountServiceWrapper) Lenient(lenient bool) *ESCountServiceWrapper {
	w.obj = w.obj.Lenient(lenient)
	return w
}

func (w *ESCountServiceWrapper) LowercaseExpandedTerms(lowercaseExpandedTerms bool) *ESCountServiceWrapper {
	w.obj = w.obj.LowercaseExpandedTerms(lowercaseExpandedTerms)
	return w
}

func (w *ESCountServiceWrapper) MinScore(minScore interface{}) *ESCountServiceWrapper {
	w.obj = w.obj.MinScore(minScore)
	return w
}

func (w *ESCountServiceWrapper) Preference(preference string) *ESCountServiceWrapper {
	w.obj = w.obj.Preference(preference)
	return w
}

func (w *ESCountServiceWrapper) Pretty(pretty bool) *ESCountServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESCountServiceWrapper) Q(q string) *ESCountServiceWrapper {
	w.obj = w.obj.Q(q)
	return w
}

func (w *ESCountServiceWrapper) Query(query elastic.Query) *ESCountServiceWrapper {
	w.obj = w.obj.Query(query)
	return w
}

func (w *ESCountServiceWrapper) Routing(routing string) *ESCountServiceWrapper {
	w.obj = w.obj.Routing(routing)
	return w
}

func (w *ESCountServiceWrapper) TerminateAfter(terminateAfter int) *ESCountServiceWrapper {
	w.obj = w.obj.TerminateAfter(terminateAfter)
	return w
}

func (w *ESCountServiceWrapper) Type(typ ...string) *ESCountServiceWrapper {
	w.obj = w.obj.Type(typ...)
	return w
}

func (w *ESCountServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESDeleteByQueryServiceWrapper) AbortOnVersionConflict() *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.AbortOnVersionConflict()
	return w
}

func (w *ESDeleteByQueryServiceWrapper) AllowNoIndices(allow bool) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allow)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) AnalyzeWildcard(analyzeWildcard bool) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.AnalyzeWildcard(analyzeWildcard)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Analyzer(analyzer string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Analyzer(analyzer)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Body(body string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Body(body)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Conflicts(conflicts string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Conflicts(conflicts)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) DF(defaultField string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.DF(defaultField)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) DefaultField(defaultField string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.DefaultField(defaultField)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) DefaultOperator(defaultOperator string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.DefaultOperator(defaultOperator)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Do(ctx context.Context) (*elastic.BulkIndexByScrollResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.BulkIndexByScrollResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DeleteByQueryService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DeleteByQueryService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DeleteByQueryService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.DeleteByQueryService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.DeleteByQueryService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.DeleteByQueryService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.DeleteByQueryService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESDeleteByQueryServiceWrapper) DoAsync(ctx context.Context) (*elastic.StartTaskResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.StartTaskResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DeleteByQueryService.DoAsync", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DeleteByQueryService.DoAsync", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DeleteByQueryService.DoAsync", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.DeleteByQueryService.DoAsync")
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
			w.inflightMetric.WithLabelValues("elastic.DeleteByQueryService.DoAsync", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.DeleteByQueryService.DoAsync", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.DeleteByQueryService.DoAsync", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DoAsync(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESDeleteByQueryServiceWrapper) DocvalueFields(docvalueFields ...string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.DocvalueFields(docvalueFields...)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) ErrorTrace(errorTrace bool) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) ExpandWildcards(expand string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expand)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Explain(explain bool) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Explain(explain)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) FilterPath(filterPath ...string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) From(from int) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.From(from)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Header(name string, value string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Headers(headers http.Header) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Human(human bool) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) IgnoreUnavailable(ignore bool) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignore)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Index(index ...string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Index(index...)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Lenient(lenient bool) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Lenient(lenient)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) LowercaseExpandedTerms(lowercaseExpandedTerms bool) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.LowercaseExpandedTerms(lowercaseExpandedTerms)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Preference(preference string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Preference(preference)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Pretty(pretty bool) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) ProceedOnVersionConflict() *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.ProceedOnVersionConflict()
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Q(query string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Q(query)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Query(query elastic.Query) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Query(query)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) QueryString(query string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.QueryString(query)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Refresh(refresh string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) RequestCache(requestCache bool) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.RequestCache(requestCache)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) RequestsPerSecond(requestsPerSecond int) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.RequestsPerSecond(requestsPerSecond)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Routing(routing ...string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Routing(routing...)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Scroll(scroll string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Scroll(scroll)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) ScrollSize(scrollSize int) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.ScrollSize(scrollSize)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) SearchTimeout(searchTimeout string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.SearchTimeout(searchTimeout)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) SearchType(searchType string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.SearchType(searchType)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Size(size int) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Size(size)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Slices(slices interface{}) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Slices(slices)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Sort(sort ...string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Sort(sort...)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) SortByField(field string, ascending bool) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.SortByField(field, ascending)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Stats(stats ...string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Stats(stats...)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) StoredFields(storedFields ...string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.StoredFields(storedFields...)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) SuggestField(suggestField string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.SuggestField(suggestField)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) SuggestMode(suggestMode string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.SuggestMode(suggestMode)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) SuggestSize(suggestSize int) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.SuggestSize(suggestSize)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) SuggestText(suggestText string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.SuggestText(suggestText)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) TerminateAfter(terminateAfter int) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.TerminateAfter(terminateAfter)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Timeout(timeout string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) TimeoutInMillis(timeoutInMillis int) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.TimeoutInMillis(timeoutInMillis)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) TrackScores(trackScores bool) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.TrackScores(trackScores)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Type(typ ...string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Type(typ...)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESDeleteByQueryServiceWrapper) Version(version bool) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.Version(version)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) WaitForActiveShards(waitForActiveShards string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.WaitForActiveShards(waitForActiveShards)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) WaitForCompletion(waitForCompletion bool) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.WaitForCompletion(waitForCompletion)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) XSource(xSource ...string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.XSource(xSource...)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) XSourceExclude(xSourceExclude ...string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.XSourceExclude(xSourceExclude...)
	return w
}

func (w *ESDeleteByQueryServiceWrapper) XSourceInclude(xSourceInclude ...string) *ESDeleteByQueryServiceWrapper {
	w.obj = w.obj.XSourceInclude(xSourceInclude...)
	return w
}

func (w *ESDeleteScriptServiceWrapper) Do(ctx context.Context) (*elastic.DeleteScriptResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.DeleteScriptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DeleteScriptService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DeleteScriptService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DeleteScriptService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.DeleteScriptService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.DeleteScriptService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.DeleteScriptService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.DeleteScriptService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESDeleteScriptServiceWrapper) ErrorTrace(errorTrace bool) *ESDeleteScriptServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESDeleteScriptServiceWrapper) FilterPath(filterPath ...string) *ESDeleteScriptServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESDeleteScriptServiceWrapper) Header(name string, value string) *ESDeleteScriptServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESDeleteScriptServiceWrapper) Headers(headers http.Header) *ESDeleteScriptServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESDeleteScriptServiceWrapper) Human(human bool) *ESDeleteScriptServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESDeleteScriptServiceWrapper) Id(id string) *ESDeleteScriptServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESDeleteScriptServiceWrapper) MasterTimeout(masterTimeout string) *ESDeleteScriptServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESDeleteScriptServiceWrapper) Pretty(pretty bool) *ESDeleteScriptServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESDeleteScriptServiceWrapper) Timeout(timeout string) *ESDeleteScriptServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESDeleteScriptServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESDeleteServiceWrapper) Do(ctx context.Context) (*elastic.DeleteResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.DeleteResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.DeleteService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.DeleteService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.DeleteService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.DeleteService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.DeleteService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.DeleteService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.DeleteService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESDeleteServiceWrapper) ErrorTrace(errorTrace bool) *ESDeleteServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESDeleteServiceWrapper) FilterPath(filterPath ...string) *ESDeleteServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESDeleteServiceWrapper) Header(name string, value string) *ESDeleteServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESDeleteServiceWrapper) Headers(headers http.Header) *ESDeleteServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESDeleteServiceWrapper) Human(human bool) *ESDeleteServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESDeleteServiceWrapper) Id(id string) *ESDeleteServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESDeleteServiceWrapper) IfPrimaryTerm(primaryTerm int64) *ESDeleteServiceWrapper {
	w.obj = w.obj.IfPrimaryTerm(primaryTerm)
	return w
}

func (w *ESDeleteServiceWrapper) IfSeqNo(seqNo int64) *ESDeleteServiceWrapper {
	w.obj = w.obj.IfSeqNo(seqNo)
	return w
}

func (w *ESDeleteServiceWrapper) Index(index string) *ESDeleteServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESDeleteServiceWrapper) Parent(parent string) *ESDeleteServiceWrapper {
	w.obj = w.obj.Parent(parent)
	return w
}

func (w *ESDeleteServiceWrapper) Pretty(pretty bool) *ESDeleteServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESDeleteServiceWrapper) Refresh(refresh string) *ESDeleteServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESDeleteServiceWrapper) Routing(routing string) *ESDeleteServiceWrapper {
	w.obj = w.obj.Routing(routing)
	return w
}

func (w *ESDeleteServiceWrapper) Timeout(timeout string) *ESDeleteServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESDeleteServiceWrapper) Type(typ string) *ESDeleteServiceWrapper {
	w.obj = w.obj.Type(typ)
	return w
}

func (w *ESDeleteServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESDeleteServiceWrapper) Version(version interface{}) *ESDeleteServiceWrapper {
	w.obj = w.obj.Version(version)
	return w
}

func (w *ESDeleteServiceWrapper) VersionType(versionType string) *ESDeleteServiceWrapper {
	w.obj = w.obj.VersionType(versionType)
	return w
}

func (w *ESDeleteServiceWrapper) WaitForActiveShards(waitForActiveShards string) *ESDeleteServiceWrapper {
	w.obj = w.obj.WaitForActiveShards(waitForActiveShards)
	return w
}

func (w *ESExistsServiceWrapper) Do(ctx context.Context) (bool, error) {
	ctxOptions := FromContext(ctx)
	var res0 bool
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.ExistsService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.ExistsService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.ExistsService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.ExistsService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.ExistsService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.ExistsService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.ExistsService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESExistsServiceWrapper) ErrorTrace(errorTrace bool) *ESExistsServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESExistsServiceWrapper) FilterPath(filterPath ...string) *ESExistsServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESExistsServiceWrapper) Header(name string, value string) *ESExistsServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESExistsServiceWrapper) Headers(headers http.Header) *ESExistsServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESExistsServiceWrapper) Human(human bool) *ESExistsServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESExistsServiceWrapper) Id(id string) *ESExistsServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESExistsServiceWrapper) Index(index string) *ESExistsServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESExistsServiceWrapper) Parent(parent string) *ESExistsServiceWrapper {
	w.obj = w.obj.Parent(parent)
	return w
}

func (w *ESExistsServiceWrapper) Preference(preference string) *ESExistsServiceWrapper {
	w.obj = w.obj.Preference(preference)
	return w
}

func (w *ESExistsServiceWrapper) Pretty(pretty bool) *ESExistsServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESExistsServiceWrapper) Realtime(realtime bool) *ESExistsServiceWrapper {
	w.obj = w.obj.Realtime(realtime)
	return w
}

func (w *ESExistsServiceWrapper) Refresh(refresh string) *ESExistsServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESExistsServiceWrapper) Routing(routing string) *ESExistsServiceWrapper {
	w.obj = w.obj.Routing(routing)
	return w
}

func (w *ESExistsServiceWrapper) Type(typ string) *ESExistsServiceWrapper {
	w.obj = w.obj.Type(typ)
	return w
}

func (w *ESExistsServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESExplainServiceWrapper) AnalyzeWildcard(analyzeWildcard bool) *ESExplainServiceWrapper {
	w.obj = w.obj.AnalyzeWildcard(analyzeWildcard)
	return w
}

func (w *ESExplainServiceWrapper) Analyzer(analyzer string) *ESExplainServiceWrapper {
	w.obj = w.obj.Analyzer(analyzer)
	return w
}

func (w *ESExplainServiceWrapper) BodyJson(body interface{}) *ESExplainServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESExplainServiceWrapper) BodyString(body string) *ESExplainServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESExplainServiceWrapper) DefaultOperator(defaultOperator string) *ESExplainServiceWrapper {
	w.obj = w.obj.DefaultOperator(defaultOperator)
	return w
}

func (w *ESExplainServiceWrapper) Df(df string) *ESExplainServiceWrapper {
	w.obj = w.obj.Df(df)
	return w
}

func (w *ESExplainServiceWrapper) Do(ctx context.Context) (*elastic.ExplainResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.ExplainResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.ExplainService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.ExplainService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.ExplainService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.ExplainService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.ExplainService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.ExplainService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.ExplainService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESExplainServiceWrapper) ErrorTrace(errorTrace bool) *ESExplainServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESExplainServiceWrapper) Fields(fields ...string) *ESExplainServiceWrapper {
	w.obj = w.obj.Fields(fields...)
	return w
}

func (w *ESExplainServiceWrapper) FilterPath(filterPath ...string) *ESExplainServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESExplainServiceWrapper) Header(name string, value string) *ESExplainServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESExplainServiceWrapper) Headers(headers http.Header) *ESExplainServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESExplainServiceWrapper) Human(human bool) *ESExplainServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESExplainServiceWrapper) Id(id string) *ESExplainServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESExplainServiceWrapper) Index(index string) *ESExplainServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESExplainServiceWrapper) Lenient(lenient bool) *ESExplainServiceWrapper {
	w.obj = w.obj.Lenient(lenient)
	return w
}

func (w *ESExplainServiceWrapper) LowercaseExpandedTerms(lowercaseExpandedTerms bool) *ESExplainServiceWrapper {
	w.obj = w.obj.LowercaseExpandedTerms(lowercaseExpandedTerms)
	return w
}

func (w *ESExplainServiceWrapper) Parent(parent string) *ESExplainServiceWrapper {
	w.obj = w.obj.Parent(parent)
	return w
}

func (w *ESExplainServiceWrapper) Preference(preference string) *ESExplainServiceWrapper {
	w.obj = w.obj.Preference(preference)
	return w
}

func (w *ESExplainServiceWrapper) Pretty(pretty bool) *ESExplainServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESExplainServiceWrapper) Q(q string) *ESExplainServiceWrapper {
	w.obj = w.obj.Q(q)
	return w
}

func (w *ESExplainServiceWrapper) Query(query elastic.Query) *ESExplainServiceWrapper {
	w.obj = w.obj.Query(query)
	return w
}

func (w *ESExplainServiceWrapper) Routing(routing string) *ESExplainServiceWrapper {
	w.obj = w.obj.Routing(routing)
	return w
}

func (w *ESExplainServiceWrapper) Source(source string) *ESExplainServiceWrapper {
	w.obj = w.obj.Source(source)
	return w
}

func (w *ESExplainServiceWrapper) Type(typ string) *ESExplainServiceWrapper {
	w.obj = w.obj.Type(typ)
	return w
}

func (w *ESExplainServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESExplainServiceWrapper) XSource(xSource ...string) *ESExplainServiceWrapper {
	w.obj = w.obj.XSource(xSource...)
	return w
}

func (w *ESExplainServiceWrapper) XSourceExclude(xSourceExclude ...string) *ESExplainServiceWrapper {
	w.obj = w.obj.XSourceExclude(xSourceExclude...)
	return w
}

func (w *ESExplainServiceWrapper) XSourceInclude(xSourceInclude ...string) *ESExplainServiceWrapper {
	w.obj = w.obj.XSourceInclude(xSourceInclude...)
	return w
}

func (w *ESFieldCapsServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESFieldCapsServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESFieldCapsServiceWrapper) BodyJson(body interface{}) *ESFieldCapsServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESFieldCapsServiceWrapper) BodyString(body string) *ESFieldCapsServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESFieldCapsServiceWrapper) Do(ctx context.Context) (*elastic.FieldCapsResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.FieldCapsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.FieldCapsService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.FieldCapsService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.FieldCapsService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.FieldCapsService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.FieldCapsService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.FieldCapsService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.FieldCapsService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESFieldCapsServiceWrapper) ErrorTrace(errorTrace bool) *ESFieldCapsServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESFieldCapsServiceWrapper) ExpandWildcards(expandWildcards string) *ESFieldCapsServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESFieldCapsServiceWrapper) Fields(fields ...string) *ESFieldCapsServiceWrapper {
	w.obj = w.obj.Fields(fields...)
	return w
}

func (w *ESFieldCapsServiceWrapper) FilterPath(filterPath ...string) *ESFieldCapsServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESFieldCapsServiceWrapper) Header(name string, value string) *ESFieldCapsServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESFieldCapsServiceWrapper) Headers(headers http.Header) *ESFieldCapsServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESFieldCapsServiceWrapper) Human(human bool) *ESFieldCapsServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESFieldCapsServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESFieldCapsServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESFieldCapsServiceWrapper) Index(index ...string) *ESFieldCapsServiceWrapper {
	w.obj = w.obj.Index(index...)
	return w
}

func (w *ESFieldCapsServiceWrapper) Pretty(pretty bool) *ESFieldCapsServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESFieldCapsServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESGetScriptServiceWrapper) Do(ctx context.Context) (*elastic.GetScriptResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.GetScriptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.GetScriptService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.GetScriptService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.GetScriptService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.GetScriptService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.GetScriptService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.GetScriptService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.GetScriptService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESGetScriptServiceWrapper) ErrorTrace(errorTrace bool) *ESGetScriptServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESGetScriptServiceWrapper) FilterPath(filterPath ...string) *ESGetScriptServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESGetScriptServiceWrapper) Header(name string, value string) *ESGetScriptServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESGetScriptServiceWrapper) Headers(headers http.Header) *ESGetScriptServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESGetScriptServiceWrapper) Human(human bool) *ESGetScriptServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESGetScriptServiceWrapper) Id(id string) *ESGetScriptServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESGetScriptServiceWrapper) Pretty(pretty bool) *ESGetScriptServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESGetScriptServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESGetServiceWrapper) Do(ctx context.Context) (*elastic.GetResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.GetResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.GetService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.GetService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.GetService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.GetService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.GetService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.GetService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.GetService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESGetServiceWrapper) ErrorTrace(errorTrace bool) *ESGetServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESGetServiceWrapper) FetchSource(fetchSource bool) *ESGetServiceWrapper {
	w.obj = w.obj.FetchSource(fetchSource)
	return w
}

func (w *ESGetServiceWrapper) FetchSourceContext(fetchSourceContext *elastic.FetchSourceContext) *ESGetServiceWrapper {
	w.obj = w.obj.FetchSourceContext(fetchSourceContext)
	return w
}

func (w *ESGetServiceWrapper) FilterPath(filterPath ...string) *ESGetServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESGetServiceWrapper) Header(name string, value string) *ESGetServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESGetServiceWrapper) Headers(headers http.Header) *ESGetServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESGetServiceWrapper) Human(human bool) *ESGetServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESGetServiceWrapper) Id(id string) *ESGetServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESGetServiceWrapper) IgnoreErrorsOnGeneratedFields(ignore bool) *ESGetServiceWrapper {
	w.obj = w.obj.IgnoreErrorsOnGeneratedFields(ignore)
	return w
}

func (w *ESGetServiceWrapper) Index(index string) *ESGetServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESGetServiceWrapper) Parent(parent string) *ESGetServiceWrapper {
	w.obj = w.obj.Parent(parent)
	return w
}

func (w *ESGetServiceWrapper) Preference(preference string) *ESGetServiceWrapper {
	w.obj = w.obj.Preference(preference)
	return w
}

func (w *ESGetServiceWrapper) Pretty(pretty bool) *ESGetServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESGetServiceWrapper) Realtime(realtime bool) *ESGetServiceWrapper {
	w.obj = w.obj.Realtime(realtime)
	return w
}

func (w *ESGetServiceWrapper) Refresh(refresh string) *ESGetServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESGetServiceWrapper) Routing(routing string) *ESGetServiceWrapper {
	w.obj = w.obj.Routing(routing)
	return w
}

func (w *ESGetServiceWrapper) StoredFields(storedFields ...string) *ESGetServiceWrapper {
	w.obj = w.obj.StoredFields(storedFields...)
	return w
}

func (w *ESGetServiceWrapper) Type(typ string) *ESGetServiceWrapper {
	w.obj = w.obj.Type(typ)
	return w
}

func (w *ESGetServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESGetServiceWrapper) Version(version interface{}) *ESGetServiceWrapper {
	w.obj = w.obj.Version(version)
	return w
}

func (w *ESGetServiceWrapper) VersionType(versionType string) *ESGetServiceWrapper {
	w.obj = w.obj.VersionType(versionType)
	return w
}

func (w *ESIndexServiceWrapper) BodyJson(body interface{}) *ESIndexServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESIndexServiceWrapper) BodyString(body string) *ESIndexServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESIndexServiceWrapper) Do(ctx context.Context) (*elastic.IndexResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndexResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndexService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndexService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndexService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndexService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndexService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndexService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndexService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndexServiceWrapper) ErrorTrace(errorTrace bool) *ESIndexServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndexServiceWrapper) FilterPath(filterPath ...string) *ESIndexServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndexServiceWrapper) Header(name string, value string) *ESIndexServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndexServiceWrapper) Headers(headers http.Header) *ESIndexServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndexServiceWrapper) Human(human bool) *ESIndexServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndexServiceWrapper) Id(id string) *ESIndexServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESIndexServiceWrapper) IfPrimaryTerm(primaryTerm int64) *ESIndexServiceWrapper {
	w.obj = w.obj.IfPrimaryTerm(primaryTerm)
	return w
}

func (w *ESIndexServiceWrapper) IfSeqNo(seqNo int64) *ESIndexServiceWrapper {
	w.obj = w.obj.IfSeqNo(seqNo)
	return w
}

func (w *ESIndexServiceWrapper) Index(index string) *ESIndexServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESIndexServiceWrapper) OpType(opType string) *ESIndexServiceWrapper {
	w.obj = w.obj.OpType(opType)
	return w
}

func (w *ESIndexServiceWrapper) Parent(parent string) *ESIndexServiceWrapper {
	w.obj = w.obj.Parent(parent)
	return w
}

func (w *ESIndexServiceWrapper) Pipeline(pipeline string) *ESIndexServiceWrapper {
	w.obj = w.obj.Pipeline(pipeline)
	return w
}

func (w *ESIndexServiceWrapper) Pretty(pretty bool) *ESIndexServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndexServiceWrapper) Refresh(refresh string) *ESIndexServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESIndexServiceWrapper) Routing(routing string) *ESIndexServiceWrapper {
	w.obj = w.obj.Routing(routing)
	return w
}

func (w *ESIndexServiceWrapper) TTL(ttl string) *ESIndexServiceWrapper {
	w.obj = w.obj.TTL(ttl)
	return w
}

func (w *ESIndexServiceWrapper) Timeout(timeout string) *ESIndexServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIndexServiceWrapper) Timestamp(timestamp string) *ESIndexServiceWrapper {
	w.obj = w.obj.Timestamp(timestamp)
	return w
}

func (w *ESIndexServiceWrapper) Ttl(ttl string) *ESIndexServiceWrapper {
	w.obj = w.obj.Ttl(ttl)
	return w
}

func (w *ESIndexServiceWrapper) Type(typ string) *ESIndexServiceWrapper {
	w.obj = w.obj.Type(typ)
	return w
}

func (w *ESIndexServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndexServiceWrapper) Version(version interface{}) *ESIndexServiceWrapper {
	w.obj = w.obj.Version(version)
	return w
}

func (w *ESIndexServiceWrapper) VersionType(versionType string) *ESIndexServiceWrapper {
	w.obj = w.obj.VersionType(versionType)
	return w
}

func (w *ESIndexServiceWrapper) WaitForActiveShards(waitForActiveShards string) *ESIndexServiceWrapper {
	w.obj = w.obj.WaitForActiveShards(waitForActiveShards)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Analyzer(analyzer string) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.Analyzer(analyzer)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Attributes(attributes ...string) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.Attributes(attributes...)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) BodyJson(body interface{}) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) BodyString(body string) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) CharFilter(charFilter ...string) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.CharFilter(charFilter...)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Do(ctx context.Context) (*elastic.IndicesAnalyzeResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesAnalyzeResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesAnalyzeService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesAnalyzeService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesAnalyzeService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesAnalyzeService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesAnalyzeService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesAnalyzeService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesAnalyzeService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesAnalyzeServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Explain(explain bool) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.Explain(explain)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Field(field string) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.Field(field)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Filter(filter ...string) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.Filter(filter...)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) FilterPath(filterPath ...string) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Format(format string) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.Format(format)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Header(name string, value string) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Headers(headers http.Header) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Human(human bool) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Index(index string) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) PreferLocal(preferLocal bool) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.PreferLocal(preferLocal)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Pretty(pretty bool) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Request(request *elastic.IndicesAnalyzeRequest) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.Request(request)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Text(text ...string) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.Text(text...)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Tokenizer(tokenizer string) *ESIndicesAnalyzeServiceWrapper {
	w.obj = w.obj.Tokenizer(tokenizer)
	return w
}

func (w *ESIndicesAnalyzeServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesClearCacheServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesClearCacheServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesClearCacheServiceWrapper) Do(ctx context.Context) (*elastic.IndicesClearCacheResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesClearCacheResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesClearCacheService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesClearCacheService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesClearCacheService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesClearCacheService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesClearCacheService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesClearCacheService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesClearCacheService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesClearCacheServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesClearCacheServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesClearCacheServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesClearCacheServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesClearCacheServiceWrapper) FieldData(fieldData bool) *ESIndicesClearCacheServiceWrapper {
	w.obj = w.obj.FieldData(fieldData)
	return w
}

func (w *ESIndicesClearCacheServiceWrapper) Fields(fields string) *ESIndicesClearCacheServiceWrapper {
	w.obj = w.obj.Fields(fields)
	return w
}

func (w *ESIndicesClearCacheServiceWrapper) FilterPath(filterPath ...string) *ESIndicesClearCacheServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesClearCacheServiceWrapper) Header(name string, value string) *ESIndicesClearCacheServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesClearCacheServiceWrapper) Headers(headers http.Header) *ESIndicesClearCacheServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesClearCacheServiceWrapper) Human(human bool) *ESIndicesClearCacheServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesClearCacheServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesClearCacheServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesClearCacheServiceWrapper) Index(indices ...string) *ESIndicesClearCacheServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESIndicesClearCacheServiceWrapper) Pretty(pretty bool) *ESIndicesClearCacheServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesClearCacheServiceWrapper) Query(queryCache bool) *ESIndicesClearCacheServiceWrapper {
	w.obj = w.obj.Query(queryCache)
	return w
}

func (w *ESIndicesClearCacheServiceWrapper) Request(requestCache bool) *ESIndicesClearCacheServiceWrapper {
	w.obj = w.obj.Request(requestCache)
	return w
}

func (w *ESIndicesClearCacheServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesCloseServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesCloseServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesCloseServiceWrapper) Do(ctx context.Context) (*elastic.IndicesCloseResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesCloseResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesCloseService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesCloseService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesCloseService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesCloseService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesCloseService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesCloseService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesCloseService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesCloseServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesCloseServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesCloseServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesCloseServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesCloseServiceWrapper) FilterPath(filterPath ...string) *ESIndicesCloseServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesCloseServiceWrapper) Header(name string, value string) *ESIndicesCloseServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesCloseServiceWrapper) Headers(headers http.Header) *ESIndicesCloseServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesCloseServiceWrapper) Human(human bool) *ESIndicesCloseServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesCloseServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesCloseServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesCloseServiceWrapper) Index(index string) *ESIndicesCloseServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESIndicesCloseServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesCloseServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesCloseServiceWrapper) Pretty(pretty bool) *ESIndicesCloseServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesCloseServiceWrapper) Timeout(timeout string) *ESIndicesCloseServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIndicesCloseServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesCreateServiceWrapper) Body(body string) *ESIndicesCreateServiceWrapper {
	w.obj = w.obj.Body(body)
	return w
}

func (w *ESIndicesCreateServiceWrapper) BodyJson(body interface{}) *ESIndicesCreateServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESIndicesCreateServiceWrapper) BodyString(body string) *ESIndicesCreateServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESIndicesCreateServiceWrapper) Do(ctx context.Context) (*elastic.IndicesCreateResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesCreateResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesCreateService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesCreateService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesCreateService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesCreateService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesCreateService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesCreateService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesCreateService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesCreateServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesCreateServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesCreateServiceWrapper) FilterPath(filterPath ...string) *ESIndicesCreateServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesCreateServiceWrapper) Header(name string, value string) *ESIndicesCreateServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesCreateServiceWrapper) Headers(headers http.Header) *ESIndicesCreateServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesCreateServiceWrapper) Human(human bool) *ESIndicesCreateServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesCreateServiceWrapper) Index(index string) *ESIndicesCreateServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESIndicesCreateServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesCreateServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesCreateServiceWrapper) Pretty(pretty bool) *ESIndicesCreateServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesCreateServiceWrapper) Timeout(timeout string) *ESIndicesCreateServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIndicesDeleteIndexTemplateServiceWrapper) Do(ctx context.Context) (*elastic.IndicesDeleteIndexTemplateResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesDeleteIndexTemplateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesDeleteIndexTemplateService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesDeleteIndexTemplateService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesDeleteIndexTemplateService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesDeleteIndexTemplateService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesDeleteIndexTemplateService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesDeleteIndexTemplateService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesDeleteIndexTemplateService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesDeleteIndexTemplateServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesDeleteIndexTemplateServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesDeleteIndexTemplateServiceWrapper) FilterPath(filterPath ...string) *ESIndicesDeleteIndexTemplateServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesDeleteIndexTemplateServiceWrapper) Header(name string, value string) *ESIndicesDeleteIndexTemplateServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesDeleteIndexTemplateServiceWrapper) Headers(headers http.Header) *ESIndicesDeleteIndexTemplateServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesDeleteIndexTemplateServiceWrapper) Human(human bool) *ESIndicesDeleteIndexTemplateServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesDeleteIndexTemplateServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesDeleteIndexTemplateServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesDeleteIndexTemplateServiceWrapper) Name(name string) *ESIndicesDeleteIndexTemplateServiceWrapper {
	w.obj = w.obj.Name(name)
	return w
}

func (w *ESIndicesDeleteIndexTemplateServiceWrapper) Pretty(pretty bool) *ESIndicesDeleteIndexTemplateServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesDeleteIndexTemplateServiceWrapper) Timeout(timeout string) *ESIndicesDeleteIndexTemplateServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIndicesDeleteIndexTemplateServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesDeleteServiceWrapper) Do(ctx context.Context) (*elastic.IndicesDeleteResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesDeleteResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesDeleteService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesDeleteService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesDeleteService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesDeleteService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesDeleteService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesDeleteService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesDeleteService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesDeleteServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesDeleteServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesDeleteServiceWrapper) FilterPath(filterPath ...string) *ESIndicesDeleteServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesDeleteServiceWrapper) Header(name string, value string) *ESIndicesDeleteServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesDeleteServiceWrapper) Headers(headers http.Header) *ESIndicesDeleteServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesDeleteServiceWrapper) Human(human bool) *ESIndicesDeleteServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesDeleteServiceWrapper) Index(index []string) *ESIndicesDeleteServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESIndicesDeleteServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesDeleteServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesDeleteServiceWrapper) Pretty(pretty bool) *ESIndicesDeleteServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesDeleteServiceWrapper) Timeout(timeout string) *ESIndicesDeleteServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIndicesDeleteServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesDeleteTemplateServiceWrapper) Do(ctx context.Context) (*elastic.IndicesDeleteTemplateResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesDeleteTemplateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesDeleteTemplateService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesDeleteTemplateService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesDeleteTemplateService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesDeleteTemplateService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesDeleteTemplateService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesDeleteTemplateService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesDeleteTemplateService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesDeleteTemplateServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesDeleteTemplateServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesDeleteTemplateServiceWrapper) FilterPath(filterPath ...string) *ESIndicesDeleteTemplateServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesDeleteTemplateServiceWrapper) Header(name string, value string) *ESIndicesDeleteTemplateServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesDeleteTemplateServiceWrapper) Headers(headers http.Header) *ESIndicesDeleteTemplateServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesDeleteTemplateServiceWrapper) Human(human bool) *ESIndicesDeleteTemplateServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesDeleteTemplateServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesDeleteTemplateServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesDeleteTemplateServiceWrapper) Name(name string) *ESIndicesDeleteTemplateServiceWrapper {
	w.obj = w.obj.Name(name)
	return w
}

func (w *ESIndicesDeleteTemplateServiceWrapper) Pretty(pretty bool) *ESIndicesDeleteTemplateServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesDeleteTemplateServiceWrapper) Timeout(timeout string) *ESIndicesDeleteTemplateServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIndicesDeleteTemplateServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesExistsServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesExistsServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesExistsServiceWrapper) Do(ctx context.Context) (bool, error) {
	ctxOptions := FromContext(ctx)
	var res0 bool
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesExistsService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesExistsService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesExistsService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesExistsService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesExistsService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesExistsService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesExistsService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesExistsServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesExistsServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesExistsServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesExistsServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesExistsServiceWrapper) FilterPath(filterPath ...string) *ESIndicesExistsServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesExistsServiceWrapper) Header(name string, value string) *ESIndicesExistsServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesExistsServiceWrapper) Headers(headers http.Header) *ESIndicesExistsServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesExistsServiceWrapper) Human(human bool) *ESIndicesExistsServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesExistsServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesExistsServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesExistsServiceWrapper) Index(index []string) *ESIndicesExistsServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESIndicesExistsServiceWrapper) Local(local bool) *ESIndicesExistsServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESIndicesExistsServiceWrapper) Pretty(pretty bool) *ESIndicesExistsServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesExistsServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesExistsTemplateServiceWrapper) Do(ctx context.Context) (bool, error) {
	ctxOptions := FromContext(ctx)
	var res0 bool
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesExistsTemplateService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesExistsTemplateService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesExistsTemplateService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesExistsTemplateService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesExistsTemplateService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesExistsTemplateService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesExistsTemplateService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesExistsTemplateServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesExistsTemplateServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesExistsTemplateServiceWrapper) FilterPath(filterPath ...string) *ESIndicesExistsTemplateServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesExistsTemplateServiceWrapper) Header(name string, value string) *ESIndicesExistsTemplateServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesExistsTemplateServiceWrapper) Headers(headers http.Header) *ESIndicesExistsTemplateServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesExistsTemplateServiceWrapper) Human(human bool) *ESIndicesExistsTemplateServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesExistsTemplateServiceWrapper) Local(local bool) *ESIndicesExistsTemplateServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESIndicesExistsTemplateServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesExistsTemplateServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesExistsTemplateServiceWrapper) Name(name string) *ESIndicesExistsTemplateServiceWrapper {
	w.obj = w.obj.Name(name)
	return w
}

func (w *ESIndicesExistsTemplateServiceWrapper) Pretty(pretty bool) *ESIndicesExistsTemplateServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesExistsTemplateServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesFlushServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesFlushServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesFlushServiceWrapper) Do(ctx context.Context) (*elastic.IndicesFlushResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesFlushResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesFlushService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesFlushService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesFlushService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesFlushService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesFlushService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesFlushService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesFlushService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesFlushServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesFlushServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesFlushServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesFlushServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesFlushServiceWrapper) FilterPath(filterPath ...string) *ESIndicesFlushServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesFlushServiceWrapper) Force(force bool) *ESIndicesFlushServiceWrapper {
	w.obj = w.obj.Force(force)
	return w
}

func (w *ESIndicesFlushServiceWrapper) Header(name string, value string) *ESIndicesFlushServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesFlushServiceWrapper) Headers(headers http.Header) *ESIndicesFlushServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesFlushServiceWrapper) Human(human bool) *ESIndicesFlushServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesFlushServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesFlushServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesFlushServiceWrapper) Index(indices ...string) *ESIndicesFlushServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESIndicesFlushServiceWrapper) Pretty(pretty bool) *ESIndicesFlushServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesFlushServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesFlushServiceWrapper) WaitIfOngoing(waitIfOngoing bool) *ESIndicesFlushServiceWrapper {
	w.obj = w.obj.WaitIfOngoing(waitIfOngoing)
	return w
}

func (w *ESIndicesForcemergeServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesForcemergeServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesForcemergeServiceWrapper) Do(ctx context.Context) (*elastic.IndicesForcemergeResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesForcemergeResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesForcemergeService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesForcemergeService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesForcemergeService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesForcemergeService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesForcemergeService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesForcemergeService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesForcemergeService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesForcemergeServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesForcemergeServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesForcemergeServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesForcemergeServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesForcemergeServiceWrapper) FilterPath(filterPath ...string) *ESIndicesForcemergeServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesForcemergeServiceWrapper) Flush(flush bool) *ESIndicesForcemergeServiceWrapper {
	w.obj = w.obj.Flush(flush)
	return w
}

func (w *ESIndicesForcemergeServiceWrapper) Header(name string, value string) *ESIndicesForcemergeServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesForcemergeServiceWrapper) Headers(headers http.Header) *ESIndicesForcemergeServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesForcemergeServiceWrapper) Human(human bool) *ESIndicesForcemergeServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesForcemergeServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesForcemergeServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesForcemergeServiceWrapper) Index(index ...string) *ESIndicesForcemergeServiceWrapper {
	w.obj = w.obj.Index(index...)
	return w
}

func (w *ESIndicesForcemergeServiceWrapper) MaxNumSegments(maxNumSegments interface{}) *ESIndicesForcemergeServiceWrapper {
	w.obj = w.obj.MaxNumSegments(maxNumSegments)
	return w
}

func (w *ESIndicesForcemergeServiceWrapper) OnlyExpungeDeletes(onlyExpungeDeletes bool) *ESIndicesForcemergeServiceWrapper {
	w.obj = w.obj.OnlyExpungeDeletes(onlyExpungeDeletes)
	return w
}

func (w *ESIndicesForcemergeServiceWrapper) Pretty(pretty bool) *ESIndicesForcemergeServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesForcemergeServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesFreezeServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesFreezeServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesFreezeServiceWrapper) Do(ctx context.Context) (*elastic.IndicesFreezeResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesFreezeResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesFreezeService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesFreezeService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesFreezeService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesFreezeService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesFreezeService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesFreezeService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesFreezeService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesFreezeServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesFreezeServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesFreezeServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesFreezeServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesFreezeServiceWrapper) FilterPath(filterPath ...string) *ESIndicesFreezeServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesFreezeServiceWrapper) Header(name string, value string) *ESIndicesFreezeServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesFreezeServiceWrapper) Headers(headers http.Header) *ESIndicesFreezeServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesFreezeServiceWrapper) Human(human bool) *ESIndicesFreezeServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesFreezeServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesFreezeServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesFreezeServiceWrapper) Index(index string) *ESIndicesFreezeServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESIndicesFreezeServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesFreezeServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesFreezeServiceWrapper) Pretty(pretty bool) *ESIndicesFreezeServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesFreezeServiceWrapper) Timeout(timeout string) *ESIndicesFreezeServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIndicesFreezeServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesFreezeServiceWrapper) WaitForActiveShards(numShards string) *ESIndicesFreezeServiceWrapper {
	w.obj = w.obj.WaitForActiveShards(numShards)
	return w
}

func (w *ESIndicesGetFieldMappingServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesGetFieldMappingServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesGetFieldMappingServiceWrapper) Do(ctx context.Context) (map[string]interface{}, error) {
	ctxOptions := FromContext(ctx)
	var res0 map[string]interface{}
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesGetFieldMappingService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesGetFieldMappingService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesGetFieldMappingService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesGetFieldMappingService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesGetFieldMappingService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesGetFieldMappingService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesGetFieldMappingService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesGetFieldMappingServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesGetFieldMappingServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesGetFieldMappingServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesGetFieldMappingServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesGetFieldMappingServiceWrapper) Field(fields ...string) *ESIndicesGetFieldMappingServiceWrapper {
	w.obj = w.obj.Field(fields...)
	return w
}

func (w *ESIndicesGetFieldMappingServiceWrapper) FilterPath(filterPath ...string) *ESIndicesGetFieldMappingServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesGetFieldMappingServiceWrapper) Header(name string, value string) *ESIndicesGetFieldMappingServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesGetFieldMappingServiceWrapper) Headers(headers http.Header) *ESIndicesGetFieldMappingServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesGetFieldMappingServiceWrapper) Human(human bool) *ESIndicesGetFieldMappingServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesGetFieldMappingServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesGetFieldMappingServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesGetFieldMappingServiceWrapper) Index(indices ...string) *ESIndicesGetFieldMappingServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESIndicesGetFieldMappingServiceWrapper) Local(local bool) *ESIndicesGetFieldMappingServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESIndicesGetFieldMappingServiceWrapper) Pretty(pretty bool) *ESIndicesGetFieldMappingServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesGetFieldMappingServiceWrapper) Type(types ...string) *ESIndicesGetFieldMappingServiceWrapper {
	w.obj = w.obj.Type(types...)
	return w
}

func (w *ESIndicesGetFieldMappingServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesGetIndexTemplateServiceWrapper) Do(ctx context.Context) (*elastic.IndicesGetIndexTemplateResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesGetIndexTemplateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesGetIndexTemplateService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesGetIndexTemplateService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesGetIndexTemplateService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesGetIndexTemplateService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesGetIndexTemplateService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesGetIndexTemplateService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesGetIndexTemplateService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesGetIndexTemplateServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesGetIndexTemplateServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesGetIndexTemplateServiceWrapper) FilterPath(filterPath ...string) *ESIndicesGetIndexTemplateServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesGetIndexTemplateServiceWrapper) FlatSettings(flatSettings bool) *ESIndicesGetIndexTemplateServiceWrapper {
	w.obj = w.obj.FlatSettings(flatSettings)
	return w
}

func (w *ESIndicesGetIndexTemplateServiceWrapper) Header(name string, value string) *ESIndicesGetIndexTemplateServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesGetIndexTemplateServiceWrapper) Headers(headers http.Header) *ESIndicesGetIndexTemplateServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesGetIndexTemplateServiceWrapper) Human(human bool) *ESIndicesGetIndexTemplateServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesGetIndexTemplateServiceWrapper) Local(local bool) *ESIndicesGetIndexTemplateServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESIndicesGetIndexTemplateServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesGetIndexTemplateServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesGetIndexTemplateServiceWrapper) Name(name ...string) *ESIndicesGetIndexTemplateServiceWrapper {
	w.obj = w.obj.Name(name...)
	return w
}

func (w *ESIndicesGetIndexTemplateServiceWrapper) Pretty(pretty bool) *ESIndicesGetIndexTemplateServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesGetIndexTemplateServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesGetMappingServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesGetMappingServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesGetMappingServiceWrapper) Do(ctx context.Context) (map[string]interface{}, error) {
	ctxOptions := FromContext(ctx)
	var res0 map[string]interface{}
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesGetMappingService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesGetMappingService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesGetMappingService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesGetMappingService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesGetMappingService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesGetMappingService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesGetMappingService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesGetMappingServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesGetMappingServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesGetMappingServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesGetMappingServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesGetMappingServiceWrapper) FilterPath(filterPath ...string) *ESIndicesGetMappingServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesGetMappingServiceWrapper) Header(name string, value string) *ESIndicesGetMappingServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesGetMappingServiceWrapper) Headers(headers http.Header) *ESIndicesGetMappingServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesGetMappingServiceWrapper) Human(human bool) *ESIndicesGetMappingServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesGetMappingServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesGetMappingServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesGetMappingServiceWrapper) Index(indices ...string) *ESIndicesGetMappingServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESIndicesGetMappingServiceWrapper) Local(local bool) *ESIndicesGetMappingServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESIndicesGetMappingServiceWrapper) Pretty(pretty bool) *ESIndicesGetMappingServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesGetMappingServiceWrapper) Type(types ...string) *ESIndicesGetMappingServiceWrapper {
	w.obj = w.obj.Type(types...)
	return w
}

func (w *ESIndicesGetMappingServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesGetServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesGetServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesGetServiceWrapper) Do(ctx context.Context) (map[string]*elastic.IndicesGetResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 map[string]*elastic.IndicesGetResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesGetService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesGetService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesGetService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesGetService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesGetService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesGetService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesGetService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesGetServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesGetServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesGetServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesGetServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesGetServiceWrapper) Feature(features ...string) *ESIndicesGetServiceWrapper {
	w.obj = w.obj.Feature(features...)
	return w
}

func (w *ESIndicesGetServiceWrapper) FilterPath(filterPath ...string) *ESIndicesGetServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesGetServiceWrapper) Header(name string, value string) *ESIndicesGetServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesGetServiceWrapper) Headers(headers http.Header) *ESIndicesGetServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesGetServiceWrapper) Human(human bool) *ESIndicesGetServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesGetServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesGetServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesGetServiceWrapper) Index(indices ...string) *ESIndicesGetServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESIndicesGetServiceWrapper) Local(local bool) *ESIndicesGetServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESIndicesGetServiceWrapper) Pretty(pretty bool) *ESIndicesGetServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesGetServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesGetSettingsServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesGetSettingsServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesGetSettingsServiceWrapper) Do(ctx context.Context) (map[string]*elastic.IndicesGetSettingsResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 map[string]*elastic.IndicesGetSettingsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesGetSettingsService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesGetSettingsService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesGetSettingsService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesGetSettingsService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesGetSettingsService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesGetSettingsService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesGetSettingsService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesGetSettingsServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesGetSettingsServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesGetSettingsServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesGetSettingsServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesGetSettingsServiceWrapper) FilterPath(filterPath ...string) *ESIndicesGetSettingsServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesGetSettingsServiceWrapper) FlatSettings(flatSettings bool) *ESIndicesGetSettingsServiceWrapper {
	w.obj = w.obj.FlatSettings(flatSettings)
	return w
}

func (w *ESIndicesGetSettingsServiceWrapper) Header(name string, value string) *ESIndicesGetSettingsServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesGetSettingsServiceWrapper) Headers(headers http.Header) *ESIndicesGetSettingsServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesGetSettingsServiceWrapper) Human(human bool) *ESIndicesGetSettingsServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesGetSettingsServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesGetSettingsServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesGetSettingsServiceWrapper) Index(indices ...string) *ESIndicesGetSettingsServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESIndicesGetSettingsServiceWrapper) Local(local bool) *ESIndicesGetSettingsServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESIndicesGetSettingsServiceWrapper) Name(name ...string) *ESIndicesGetSettingsServiceWrapper {
	w.obj = w.obj.Name(name...)
	return w
}

func (w *ESIndicesGetSettingsServiceWrapper) Pretty(pretty bool) *ESIndicesGetSettingsServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesGetSettingsServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesGetTemplateServiceWrapper) Do(ctx context.Context) (map[string]*elastic.IndicesGetTemplateResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 map[string]*elastic.IndicesGetTemplateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesGetTemplateService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesGetTemplateService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesGetTemplateService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesGetTemplateService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesGetTemplateService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesGetTemplateService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesGetTemplateService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesGetTemplateServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesGetTemplateServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesGetTemplateServiceWrapper) FilterPath(filterPath ...string) *ESIndicesGetTemplateServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesGetTemplateServiceWrapper) FlatSettings(flatSettings bool) *ESIndicesGetTemplateServiceWrapper {
	w.obj = w.obj.FlatSettings(flatSettings)
	return w
}

func (w *ESIndicesGetTemplateServiceWrapper) Header(name string, value string) *ESIndicesGetTemplateServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesGetTemplateServiceWrapper) Headers(headers http.Header) *ESIndicesGetTemplateServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesGetTemplateServiceWrapper) Human(human bool) *ESIndicesGetTemplateServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesGetTemplateServiceWrapper) Local(local bool) *ESIndicesGetTemplateServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESIndicesGetTemplateServiceWrapper) Name(name ...string) *ESIndicesGetTemplateServiceWrapper {
	w.obj = w.obj.Name(name...)
	return w
}

func (w *ESIndicesGetTemplateServiceWrapper) Pretty(pretty bool) *ESIndicesGetTemplateServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesGetTemplateServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesOpenServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesOpenServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesOpenServiceWrapper) Do(ctx context.Context) (*elastic.IndicesOpenResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesOpenResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesOpenService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesOpenService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesOpenService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesOpenService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesOpenService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesOpenService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesOpenService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesOpenServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesOpenServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesOpenServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesOpenServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesOpenServiceWrapper) FilterPath(filterPath ...string) *ESIndicesOpenServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesOpenServiceWrapper) Header(name string, value string) *ESIndicesOpenServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesOpenServiceWrapper) Headers(headers http.Header) *ESIndicesOpenServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesOpenServiceWrapper) Human(human bool) *ESIndicesOpenServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesOpenServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesOpenServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesOpenServiceWrapper) Index(index string) *ESIndicesOpenServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESIndicesOpenServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesOpenServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesOpenServiceWrapper) Pretty(pretty bool) *ESIndicesOpenServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesOpenServiceWrapper) Timeout(timeout string) *ESIndicesOpenServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIndicesOpenServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesOpenServiceWrapper) WaitForActiveShards(waitForActiveShards string) *ESIndicesOpenServiceWrapper {
	w.obj = w.obj.WaitForActiveShards(waitForActiveShards)
	return w
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) BodyJson(body interface{}) *ESIndicesPutIndexTemplateServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) BodyString(body string) *ESIndicesPutIndexTemplateServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) Cause(cause string) *ESIndicesPutIndexTemplateServiceWrapper {
	w.obj = w.obj.Cause(cause)
	return w
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) Create(create bool) *ESIndicesPutIndexTemplateServiceWrapper {
	w.obj = w.obj.Create(create)
	return w
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) Do(ctx context.Context) (*elastic.IndicesPutIndexTemplateResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesPutIndexTemplateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesPutIndexTemplateService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesPutIndexTemplateService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesPutIndexTemplateService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesPutIndexTemplateService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesPutIndexTemplateService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesPutIndexTemplateService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesPutIndexTemplateService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesPutIndexTemplateServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) FilterPath(filterPath ...string) *ESIndicesPutIndexTemplateServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) Header(name string, value string) *ESIndicesPutIndexTemplateServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) Headers(headers http.Header) *ESIndicesPutIndexTemplateServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) Human(human bool) *ESIndicesPutIndexTemplateServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesPutIndexTemplateServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) Name(name string) *ESIndicesPutIndexTemplateServiceWrapper {
	w.obj = w.obj.Name(name)
	return w
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) Pretty(pretty bool) *ESIndicesPutIndexTemplateServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesPutIndexTemplateServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesPutMappingServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) BodyJson(mapping map[string]interface{}) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.BodyJson(mapping)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) BodyString(mapping string) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.BodyString(mapping)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) Do(ctx context.Context) (*elastic.PutMappingResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.PutMappingResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesPutMappingService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesPutMappingService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesPutMappingService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesPutMappingService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesPutMappingService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesPutMappingService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesPutMappingService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesPutMappingServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) FilterPath(filterPath ...string) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) Header(name string, value string) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) Headers(headers http.Header) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) Human(human bool) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) Index(indices ...string) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) Pretty(pretty bool) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) Timeout(timeout string) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) UpdateAllTypes(updateAllTypes bool) *ESIndicesPutMappingServiceWrapper {
	w.obj = w.obj.UpdateAllTypes(updateAllTypes)
	return w
}

func (w *ESIndicesPutMappingServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesPutSettingsServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesPutSettingsServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesPutSettingsServiceWrapper) BodyJson(body interface{}) *ESIndicesPutSettingsServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESIndicesPutSettingsServiceWrapper) BodyString(body string) *ESIndicesPutSettingsServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESIndicesPutSettingsServiceWrapper) Do(ctx context.Context) (*elastic.IndicesPutSettingsResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesPutSettingsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesPutSettingsService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesPutSettingsService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesPutSettingsService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesPutSettingsService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesPutSettingsService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesPutSettingsService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesPutSettingsService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesPutSettingsServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesPutSettingsServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesPutSettingsServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesPutSettingsServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesPutSettingsServiceWrapper) FilterPath(filterPath ...string) *ESIndicesPutSettingsServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesPutSettingsServiceWrapper) FlatSettings(flatSettings bool) *ESIndicesPutSettingsServiceWrapper {
	w.obj = w.obj.FlatSettings(flatSettings)
	return w
}

func (w *ESIndicesPutSettingsServiceWrapper) Header(name string, value string) *ESIndicesPutSettingsServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesPutSettingsServiceWrapper) Headers(headers http.Header) *ESIndicesPutSettingsServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesPutSettingsServiceWrapper) Human(human bool) *ESIndicesPutSettingsServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesPutSettingsServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesPutSettingsServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesPutSettingsServiceWrapper) Index(indices ...string) *ESIndicesPutSettingsServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESIndicesPutSettingsServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesPutSettingsServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesPutSettingsServiceWrapper) Pretty(pretty bool) *ESIndicesPutSettingsServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesPutSettingsServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesPutTemplateServiceWrapper) BodyJson(body interface{}) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) BodyString(body string) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) Cause(cause string) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.Cause(cause)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) Create(create bool) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.Create(create)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) Do(ctx context.Context) (*elastic.IndicesPutTemplateResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesPutTemplateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesPutTemplateService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesPutTemplateService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesPutTemplateService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesPutTemplateService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesPutTemplateService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesPutTemplateService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesPutTemplateService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesPutTemplateServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) FilterPath(filterPath ...string) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) FlatSettings(flatSettings bool) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.FlatSettings(flatSettings)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) Header(name string, value string) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) Headers(headers http.Header) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) Human(human bool) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) Name(name string) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.Name(name)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) Order(order interface{}) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.Order(order)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) Pretty(pretty bool) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) Timeout(timeout string) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIndicesPutTemplateServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesPutTemplateServiceWrapper) Version(version int) *ESIndicesPutTemplateServiceWrapper {
	w.obj = w.obj.Version(version)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) AddCondition(name string, value interface{}) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.AddCondition(name, value)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) AddMapping(typ string, mapping interface{}) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.AddMapping(typ, mapping)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) AddMaxIndexAgeCondition(time string) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.AddMaxIndexAgeCondition(time)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) AddMaxIndexDocsCondition(docs int64) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.AddMaxIndexDocsCondition(docs)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) AddSetting(name string, value interface{}) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.AddSetting(name, value)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) Alias(alias string) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.Alias(alias)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) BodyJson(body interface{}) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) BodyString(body string) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) Conditions(conditions map[string]interface{}) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.Conditions(conditions)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) Do(ctx context.Context) (*elastic.IndicesRolloverResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesRolloverResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesRolloverService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesRolloverService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesRolloverService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesRolloverService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesRolloverService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesRolloverService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesRolloverService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesRolloverServiceWrapper) DryRun(dryRun bool) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.DryRun(dryRun)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) FilterPath(filterPath ...string) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) Header(name string, value string) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) Headers(headers http.Header) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) Human(human bool) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) Mappings(mappings map[string]interface{}) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.Mappings(mappings)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) NewIndex(newIndex string) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.NewIndex(newIndex)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) Pretty(pretty bool) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) Settings(settings map[string]interface{}) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.Settings(settings)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) Timeout(timeout string) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIndicesRolloverServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesRolloverServiceWrapper) WaitForActiveShards(waitForActiveShards string) *ESIndicesRolloverServiceWrapper {
	w.obj = w.obj.WaitForActiveShards(waitForActiveShards)
	return w
}

func (w *ESIndicesSegmentsServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesSegmentsServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesSegmentsServiceWrapper) Do(ctx context.Context) (*elastic.IndicesSegmentsResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesSegmentsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesSegmentsService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesSegmentsService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesSegmentsService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesSegmentsService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesSegmentsService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesSegmentsService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesSegmentsService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesSegmentsServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesSegmentsServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesSegmentsServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesSegmentsServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesSegmentsServiceWrapper) FilterPath(filterPath ...string) *ESIndicesSegmentsServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesSegmentsServiceWrapper) Header(name string, value string) *ESIndicesSegmentsServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesSegmentsServiceWrapper) Headers(headers http.Header) *ESIndicesSegmentsServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesSegmentsServiceWrapper) Human(human bool) *ESIndicesSegmentsServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesSegmentsServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesSegmentsServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesSegmentsServiceWrapper) Index(indices ...string) *ESIndicesSegmentsServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESIndicesSegmentsServiceWrapper) OperationThreading(operationThreading interface{}) *ESIndicesSegmentsServiceWrapper {
	w.obj = w.obj.OperationThreading(operationThreading)
	return w
}

func (w *ESIndicesSegmentsServiceWrapper) Pretty(pretty bool) *ESIndicesSegmentsServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesSegmentsServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesSegmentsServiceWrapper) Verbose(verbose bool) *ESIndicesSegmentsServiceWrapper {
	w.obj = w.obj.Verbose(verbose)
	return w
}

func (w *ESIndicesShrinkServiceWrapper) BodyJson(body interface{}) *ESIndicesShrinkServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESIndicesShrinkServiceWrapper) BodyString(body string) *ESIndicesShrinkServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESIndicesShrinkServiceWrapper) Do(ctx context.Context) (*elastic.IndicesShrinkResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesShrinkResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesShrinkService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesShrinkService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesShrinkService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesShrinkService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesShrinkService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesShrinkService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesShrinkService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesShrinkServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesShrinkServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesShrinkServiceWrapper) FilterPath(filterPath ...string) *ESIndicesShrinkServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesShrinkServiceWrapper) Header(name string, value string) *ESIndicesShrinkServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesShrinkServiceWrapper) Headers(headers http.Header) *ESIndicesShrinkServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesShrinkServiceWrapper) Human(human bool) *ESIndicesShrinkServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesShrinkServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesShrinkServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesShrinkServiceWrapper) Pretty(pretty bool) *ESIndicesShrinkServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesShrinkServiceWrapper) Source(source string) *ESIndicesShrinkServiceWrapper {
	w.obj = w.obj.Source(source)
	return w
}

func (w *ESIndicesShrinkServiceWrapper) Target(target string) *ESIndicesShrinkServiceWrapper {
	w.obj = w.obj.Target(target)
	return w
}

func (w *ESIndicesShrinkServiceWrapper) Timeout(timeout string) *ESIndicesShrinkServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIndicesShrinkServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesShrinkServiceWrapper) WaitForActiveShards(waitForActiveShards string) *ESIndicesShrinkServiceWrapper {
	w.obj = w.obj.WaitForActiveShards(waitForActiveShards)
	return w
}

func (w *ESIndicesStatsServiceWrapper) CompletionFields(completionFields ...string) *ESIndicesStatsServiceWrapper {
	w.obj = w.obj.CompletionFields(completionFields...)
	return w
}

func (w *ESIndicesStatsServiceWrapper) Do(ctx context.Context) (*elastic.IndicesStatsResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesStatsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesStatsService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesStatsService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesStatsService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesStatsService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesStatsService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesStatsService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesStatsService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesStatsServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesStatsServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesStatsServiceWrapper) FielddataFields(fielddataFields ...string) *ESIndicesStatsServiceWrapper {
	w.obj = w.obj.FielddataFields(fielddataFields...)
	return w
}

func (w *ESIndicesStatsServiceWrapper) Fields(fields ...string) *ESIndicesStatsServiceWrapper {
	w.obj = w.obj.Fields(fields...)
	return w
}

func (w *ESIndicesStatsServiceWrapper) FilterPath(filterPath ...string) *ESIndicesStatsServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesStatsServiceWrapper) Groups(groups ...string) *ESIndicesStatsServiceWrapper {
	w.obj = w.obj.Groups(groups...)
	return w
}

func (w *ESIndicesStatsServiceWrapper) Header(name string, value string) *ESIndicesStatsServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesStatsServiceWrapper) Headers(headers http.Header) *ESIndicesStatsServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesStatsServiceWrapper) Human(human bool) *ESIndicesStatsServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesStatsServiceWrapper) Index(indices ...string) *ESIndicesStatsServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESIndicesStatsServiceWrapper) Level(level string) *ESIndicesStatsServiceWrapper {
	w.obj = w.obj.Level(level)
	return w
}

func (w *ESIndicesStatsServiceWrapper) Metric(metric ...string) *ESIndicesStatsServiceWrapper {
	w.obj = w.obj.Metric(metric...)
	return w
}

func (w *ESIndicesStatsServiceWrapper) Pretty(pretty bool) *ESIndicesStatsServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesStatsServiceWrapper) Type(types ...string) *ESIndicesStatsServiceWrapper {
	w.obj = w.obj.Type(types...)
	return w
}

func (w *ESIndicesStatsServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesSyncedFlushServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesSyncedFlushServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesSyncedFlushServiceWrapper) Do(ctx context.Context) (*elastic.IndicesSyncedFlushResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesSyncedFlushResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesSyncedFlushService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesSyncedFlushService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesSyncedFlushService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesSyncedFlushService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesSyncedFlushService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesSyncedFlushService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesSyncedFlushService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesSyncedFlushServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesSyncedFlushServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesSyncedFlushServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesSyncedFlushServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesSyncedFlushServiceWrapper) FilterPath(filterPath ...string) *ESIndicesSyncedFlushServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesSyncedFlushServiceWrapper) Header(name string, value string) *ESIndicesSyncedFlushServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesSyncedFlushServiceWrapper) Headers(headers http.Header) *ESIndicesSyncedFlushServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesSyncedFlushServiceWrapper) Human(human bool) *ESIndicesSyncedFlushServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesSyncedFlushServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesSyncedFlushServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesSyncedFlushServiceWrapper) Index(indices ...string) *ESIndicesSyncedFlushServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESIndicesSyncedFlushServiceWrapper) Pretty(pretty bool) *ESIndicesSyncedFlushServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesSyncedFlushServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesUnfreezeServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESIndicesUnfreezeServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESIndicesUnfreezeServiceWrapper) Do(ctx context.Context) (*elastic.IndicesUnfreezeResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IndicesUnfreezeResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IndicesUnfreezeService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IndicesUnfreezeService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IndicesUnfreezeService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IndicesUnfreezeService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IndicesUnfreezeService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IndicesUnfreezeService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IndicesUnfreezeService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIndicesUnfreezeServiceWrapper) ErrorTrace(errorTrace bool) *ESIndicesUnfreezeServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIndicesUnfreezeServiceWrapper) ExpandWildcards(expandWildcards string) *ESIndicesUnfreezeServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESIndicesUnfreezeServiceWrapper) FilterPath(filterPath ...string) *ESIndicesUnfreezeServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIndicesUnfreezeServiceWrapper) Header(name string, value string) *ESIndicesUnfreezeServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIndicesUnfreezeServiceWrapper) Headers(headers http.Header) *ESIndicesUnfreezeServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIndicesUnfreezeServiceWrapper) Human(human bool) *ESIndicesUnfreezeServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIndicesUnfreezeServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESIndicesUnfreezeServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESIndicesUnfreezeServiceWrapper) Index(index string) *ESIndicesUnfreezeServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESIndicesUnfreezeServiceWrapper) MasterTimeout(masterTimeout string) *ESIndicesUnfreezeServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIndicesUnfreezeServiceWrapper) Pretty(pretty bool) *ESIndicesUnfreezeServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIndicesUnfreezeServiceWrapper) Timeout(timeout string) *ESIndicesUnfreezeServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIndicesUnfreezeServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIndicesUnfreezeServiceWrapper) WaitForActiveShards(numShards string) *ESIndicesUnfreezeServiceWrapper {
	w.obj = w.obj.WaitForActiveShards(numShards)
	return w
}

func (w *ESIngestDeletePipelineServiceWrapper) Do(ctx context.Context) (*elastic.IngestDeletePipelineResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IngestDeletePipelineResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IngestDeletePipelineService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IngestDeletePipelineService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IngestDeletePipelineService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IngestDeletePipelineService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IngestDeletePipelineService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IngestDeletePipelineService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IngestDeletePipelineService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIngestDeletePipelineServiceWrapper) ErrorTrace(errorTrace bool) *ESIngestDeletePipelineServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIngestDeletePipelineServiceWrapper) FilterPath(filterPath ...string) *ESIngestDeletePipelineServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIngestDeletePipelineServiceWrapper) Header(name string, value string) *ESIngestDeletePipelineServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIngestDeletePipelineServiceWrapper) Headers(headers http.Header) *ESIngestDeletePipelineServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIngestDeletePipelineServiceWrapper) Human(human bool) *ESIngestDeletePipelineServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIngestDeletePipelineServiceWrapper) Id(id string) *ESIngestDeletePipelineServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESIngestDeletePipelineServiceWrapper) MasterTimeout(masterTimeout string) *ESIngestDeletePipelineServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIngestDeletePipelineServiceWrapper) Pretty(pretty bool) *ESIngestDeletePipelineServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIngestDeletePipelineServiceWrapper) Timeout(timeout string) *ESIngestDeletePipelineServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIngestDeletePipelineServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIngestGetPipelineServiceWrapper) Do(ctx context.Context) (elastic.IngestGetPipelineResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 elastic.IngestGetPipelineResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IngestGetPipelineService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IngestGetPipelineService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IngestGetPipelineService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IngestGetPipelineService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IngestGetPipelineService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IngestGetPipelineService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IngestGetPipelineService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIngestGetPipelineServiceWrapper) ErrorTrace(errorTrace bool) *ESIngestGetPipelineServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIngestGetPipelineServiceWrapper) FilterPath(filterPath ...string) *ESIngestGetPipelineServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIngestGetPipelineServiceWrapper) Header(name string, value string) *ESIngestGetPipelineServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIngestGetPipelineServiceWrapper) Headers(headers http.Header) *ESIngestGetPipelineServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIngestGetPipelineServiceWrapper) Human(human bool) *ESIngestGetPipelineServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIngestGetPipelineServiceWrapper) Id(id ...string) *ESIngestGetPipelineServiceWrapper {
	w.obj = w.obj.Id(id...)
	return w
}

func (w *ESIngestGetPipelineServiceWrapper) MasterTimeout(masterTimeout string) *ESIngestGetPipelineServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIngestGetPipelineServiceWrapper) Pretty(pretty bool) *ESIngestGetPipelineServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIngestGetPipelineServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIngestPutPipelineServiceWrapper) BodyJson(body interface{}) *ESIngestPutPipelineServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESIngestPutPipelineServiceWrapper) BodyString(body string) *ESIngestPutPipelineServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESIngestPutPipelineServiceWrapper) Do(ctx context.Context) (*elastic.IngestPutPipelineResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IngestPutPipelineResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IngestPutPipelineService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IngestPutPipelineService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IngestPutPipelineService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IngestPutPipelineService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IngestPutPipelineService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IngestPutPipelineService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IngestPutPipelineService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIngestPutPipelineServiceWrapper) ErrorTrace(errorTrace bool) *ESIngestPutPipelineServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIngestPutPipelineServiceWrapper) FilterPath(filterPath ...string) *ESIngestPutPipelineServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIngestPutPipelineServiceWrapper) Header(name string, value string) *ESIngestPutPipelineServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIngestPutPipelineServiceWrapper) Headers(headers http.Header) *ESIngestPutPipelineServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIngestPutPipelineServiceWrapper) Human(human bool) *ESIngestPutPipelineServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIngestPutPipelineServiceWrapper) Id(id string) *ESIngestPutPipelineServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESIngestPutPipelineServiceWrapper) MasterTimeout(masterTimeout string) *ESIngestPutPipelineServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESIngestPutPipelineServiceWrapper) Pretty(pretty bool) *ESIngestPutPipelineServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIngestPutPipelineServiceWrapper) Timeout(timeout string) *ESIngestPutPipelineServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESIngestPutPipelineServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIngestSimulatePipelineServiceWrapper) BodyJson(body interface{}) *ESIngestSimulatePipelineServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESIngestSimulatePipelineServiceWrapper) BodyString(body string) *ESIngestSimulatePipelineServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESIngestSimulatePipelineServiceWrapper) Do(ctx context.Context) (*elastic.IngestSimulatePipelineResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.IngestSimulatePipelineResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.IngestSimulatePipelineService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.IngestSimulatePipelineService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.IngestSimulatePipelineService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.IngestSimulatePipelineService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.IngestSimulatePipelineService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.IngestSimulatePipelineService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.IngestSimulatePipelineService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESIngestSimulatePipelineServiceWrapper) ErrorTrace(errorTrace bool) *ESIngestSimulatePipelineServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESIngestSimulatePipelineServiceWrapper) FilterPath(filterPath ...string) *ESIngestSimulatePipelineServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESIngestSimulatePipelineServiceWrapper) Header(name string, value string) *ESIngestSimulatePipelineServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESIngestSimulatePipelineServiceWrapper) Headers(headers http.Header) *ESIngestSimulatePipelineServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESIngestSimulatePipelineServiceWrapper) Human(human bool) *ESIngestSimulatePipelineServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESIngestSimulatePipelineServiceWrapper) Id(id string) *ESIngestSimulatePipelineServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESIngestSimulatePipelineServiceWrapper) Pretty(pretty bool) *ESIngestSimulatePipelineServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESIngestSimulatePipelineServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESIngestSimulatePipelineServiceWrapper) Verbose(verbose bool) *ESIngestSimulatePipelineServiceWrapper {
	w.obj = w.obj.Verbose(verbose)
	return w
}

func (w *ESMgetServiceWrapper) Add(items ...*elastic.MultiGetItem) *ESMgetServiceWrapper {
	w.obj = w.obj.Add(items...)
	return w
}

func (w *ESMgetServiceWrapper) Do(ctx context.Context) (*elastic.MgetResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.MgetResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MgetService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MgetService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MgetService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.MgetService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.MgetService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.MgetService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.MgetService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESMgetServiceWrapper) ErrorTrace(errorTrace bool) *ESMgetServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESMgetServiceWrapper) FilterPath(filterPath ...string) *ESMgetServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESMgetServiceWrapper) Header(name string, value string) *ESMgetServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESMgetServiceWrapper) Headers(headers http.Header) *ESMgetServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESMgetServiceWrapper) Human(human bool) *ESMgetServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESMgetServiceWrapper) Preference(preference string) *ESMgetServiceWrapper {
	w.obj = w.obj.Preference(preference)
	return w
}

func (w *ESMgetServiceWrapper) Pretty(pretty bool) *ESMgetServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESMgetServiceWrapper) Realtime(realtime bool) *ESMgetServiceWrapper {
	w.obj = w.obj.Realtime(realtime)
	return w
}

func (w *ESMgetServiceWrapper) Refresh(refresh string) *ESMgetServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESMgetServiceWrapper) Routing(routing string) *ESMgetServiceWrapper {
	w.obj = w.obj.Routing(routing)
	return w
}

func (w *ESMgetServiceWrapper) Source() (interface{}, error) {
	var res0 interface{}
	var err error
	res0, err = w.obj.Source()
	return res0, err
}

func (w *ESMgetServiceWrapper) StoredFields(storedFields ...string) *ESMgetServiceWrapper {
	w.obj = w.obj.StoredFields(storedFields...)
	return w
}

func (w *ESMultiSearchServiceWrapper) Add(requests ...*elastic.SearchRequest) *ESMultiSearchServiceWrapper {
	w.obj = w.obj.Add(requests...)
	return w
}

func (w *ESMultiSearchServiceWrapper) Do(ctx context.Context) (*elastic.MultiSearchResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.MultiSearchResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MultiSearchService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MultiSearchService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MultiSearchService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.MultiSearchService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.MultiSearchService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.MultiSearchService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.MultiSearchService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESMultiSearchServiceWrapper) ErrorTrace(errorTrace bool) *ESMultiSearchServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESMultiSearchServiceWrapper) FilterPath(filterPath ...string) *ESMultiSearchServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESMultiSearchServiceWrapper) Header(name string, value string) *ESMultiSearchServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESMultiSearchServiceWrapper) Headers(headers http.Header) *ESMultiSearchServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESMultiSearchServiceWrapper) Human(human bool) *ESMultiSearchServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESMultiSearchServiceWrapper) Index(indices ...string) *ESMultiSearchServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESMultiSearchServiceWrapper) MaxConcurrentSearches(max int) *ESMultiSearchServiceWrapper {
	w.obj = w.obj.MaxConcurrentSearches(max)
	return w
}

func (w *ESMultiSearchServiceWrapper) PreFilterShardSize(size int) *ESMultiSearchServiceWrapper {
	w.obj = w.obj.PreFilterShardSize(size)
	return w
}

func (w *ESMultiSearchServiceWrapper) Pretty(pretty bool) *ESMultiSearchServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Add(docs ...*elastic.MultiTermvectorItem) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Add(docs...)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) BodyJson(body interface{}) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) BodyString(body string) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Do(ctx context.Context) (*elastic.MultiTermvectorResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.MultiTermvectorResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.MultiTermvectorService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.MultiTermvectorService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.MultiTermvectorService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.MultiTermvectorService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.MultiTermvectorService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.MultiTermvectorService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.MultiTermvectorService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESMultiTermvectorServiceWrapper) ErrorTrace(errorTrace bool) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) FieldStatistics(fieldStatistics bool) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.FieldStatistics(fieldStatistics)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Fields(fields []string) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Fields(fields)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) FilterPath(filterPath ...string) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Header(name string, value string) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Headers(headers http.Header) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Human(human bool) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Ids(ids []string) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Ids(ids)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Index(index string) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Offsets(offsets bool) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Offsets(offsets)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Parent(parent string) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Parent(parent)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Payloads(payloads bool) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Payloads(payloads)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Positions(positions bool) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Positions(positions)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Preference(preference string) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Preference(preference)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Pretty(pretty bool) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Realtime(realtime bool) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Realtime(realtime)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Routing(routing string) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Routing(routing)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Source() interface{} {
	res0 := w.obj.Source()
	return res0
}

func (w *ESMultiTermvectorServiceWrapper) TermStatistics(termStatistics bool) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.TermStatistics(termStatistics)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Type(typ string) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Type(typ)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESMultiTermvectorServiceWrapper) Version(version interface{}) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.Version(version)
	return w
}

func (w *ESMultiTermvectorServiceWrapper) VersionType(versionType string) *ESMultiTermvectorServiceWrapper {
	w.obj = w.obj.VersionType(versionType)
	return w
}

func (w *ESNodesInfoServiceWrapper) Do(ctx context.Context) (*elastic.NodesInfoResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.NodesInfoResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.NodesInfoService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.NodesInfoService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.NodesInfoService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.NodesInfoService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.NodesInfoService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.NodesInfoService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.NodesInfoService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESNodesInfoServiceWrapper) ErrorTrace(errorTrace bool) *ESNodesInfoServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESNodesInfoServiceWrapper) FilterPath(filterPath ...string) *ESNodesInfoServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESNodesInfoServiceWrapper) FlatSettings(flatSettings bool) *ESNodesInfoServiceWrapper {
	w.obj = w.obj.FlatSettings(flatSettings)
	return w
}

func (w *ESNodesInfoServiceWrapper) Header(name string, value string) *ESNodesInfoServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESNodesInfoServiceWrapper) Headers(headers http.Header) *ESNodesInfoServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESNodesInfoServiceWrapper) Human(human bool) *ESNodesInfoServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESNodesInfoServiceWrapper) Metric(metric ...string) *ESNodesInfoServiceWrapper {
	w.obj = w.obj.Metric(metric...)
	return w
}

func (w *ESNodesInfoServiceWrapper) NodeId(nodeId ...string) *ESNodesInfoServiceWrapper {
	w.obj = w.obj.NodeId(nodeId...)
	return w
}

func (w *ESNodesInfoServiceWrapper) Pretty(pretty bool) *ESNodesInfoServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESNodesInfoServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESNodesStatsServiceWrapper) CompletionFields(completionFields ...string) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.CompletionFields(completionFields...)
	return w
}

func (w *ESNodesStatsServiceWrapper) Do(ctx context.Context) (*elastic.NodesStatsResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.NodesStatsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.NodesStatsService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.NodesStatsService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.NodesStatsService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.NodesStatsService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.NodesStatsService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.NodesStatsService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.NodesStatsService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESNodesStatsServiceWrapper) ErrorTrace(errorTrace bool) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESNodesStatsServiceWrapper) FielddataFields(fielddataFields ...string) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.FielddataFields(fielddataFields...)
	return w
}

func (w *ESNodesStatsServiceWrapper) Fields(fields ...string) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.Fields(fields...)
	return w
}

func (w *ESNodesStatsServiceWrapper) FilterPath(filterPath ...string) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESNodesStatsServiceWrapper) Groups(groups bool) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.Groups(groups)
	return w
}

func (w *ESNodesStatsServiceWrapper) Header(name string, value string) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESNodesStatsServiceWrapper) Headers(headers http.Header) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESNodesStatsServiceWrapper) Human(human bool) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESNodesStatsServiceWrapper) IndexMetric(indexMetric ...string) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.IndexMetric(indexMetric...)
	return w
}

func (w *ESNodesStatsServiceWrapper) Level(level string) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.Level(level)
	return w
}

func (w *ESNodesStatsServiceWrapper) Metric(metric ...string) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.Metric(metric...)
	return w
}

func (w *ESNodesStatsServiceWrapper) NodeId(nodeId ...string) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.NodeId(nodeId...)
	return w
}

func (w *ESNodesStatsServiceWrapper) Pretty(pretty bool) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESNodesStatsServiceWrapper) Timeout(timeout string) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESNodesStatsServiceWrapper) Types(types ...string) *ESNodesStatsServiceWrapper {
	w.obj = w.obj.Types(types...)
	return w
}

func (w *ESNodesStatsServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESPingServiceWrapper) Do(ctx context.Context) (*elastic.PingResult, int, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.PingResult
	var res1 int
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.PingService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.PingService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.PingService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.PingService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.PingService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.PingService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.PingService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, res1, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, res1, err
}

func (w *ESPingServiceWrapper) ErrorTrace(errorTrace bool) *ESPingServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESPingServiceWrapper) FilterPath(filterPath ...string) *ESPingServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESPingServiceWrapper) Header(name string, value string) *ESPingServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESPingServiceWrapper) Headers(headers http.Header) *ESPingServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESPingServiceWrapper) HttpHeadOnly(httpHeadOnly bool) *ESPingServiceWrapper {
	w.obj = w.obj.HttpHeadOnly(httpHeadOnly)
	return w
}

func (w *ESPingServiceWrapper) Human(human bool) *ESPingServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESPingServiceWrapper) Pretty(pretty bool) *ESPingServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESPingServiceWrapper) Timeout(timeout string) *ESPingServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESPingServiceWrapper) URL(url string) *ESPingServiceWrapper {
	w.obj = w.obj.URL(url)
	return w
}

func (w *ESPutScriptServiceWrapper) BodyJson(body interface{}) *ESPutScriptServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESPutScriptServiceWrapper) BodyString(body string) *ESPutScriptServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESPutScriptServiceWrapper) Context(context string) *ESPutScriptServiceWrapper {
	w.obj = w.obj.Context(context)
	return w
}

func (w *ESPutScriptServiceWrapper) Do(ctx context.Context) (*elastic.PutScriptResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.PutScriptResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.PutScriptService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.PutScriptService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.PutScriptService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.PutScriptService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.PutScriptService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.PutScriptService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.PutScriptService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESPutScriptServiceWrapper) ErrorTrace(errorTrace bool) *ESPutScriptServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESPutScriptServiceWrapper) FilterPath(filterPath ...string) *ESPutScriptServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESPutScriptServiceWrapper) Header(name string, value string) *ESPutScriptServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESPutScriptServiceWrapper) Headers(headers http.Header) *ESPutScriptServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESPutScriptServiceWrapper) Human(human bool) *ESPutScriptServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESPutScriptServiceWrapper) Id(id string) *ESPutScriptServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESPutScriptServiceWrapper) MasterTimeout(masterTimeout string) *ESPutScriptServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESPutScriptServiceWrapper) Pretty(pretty bool) *ESPutScriptServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESPutScriptServiceWrapper) Timeout(timeout string) *ESPutScriptServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESPutScriptServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESRefreshServiceWrapper) Do(ctx context.Context) (*elastic.RefreshResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.RefreshResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.RefreshService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.RefreshService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.RefreshService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.RefreshService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.RefreshService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.RefreshService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.RefreshService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESRefreshServiceWrapper) ErrorTrace(errorTrace bool) *ESRefreshServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESRefreshServiceWrapper) FilterPath(filterPath ...string) *ESRefreshServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESRefreshServiceWrapper) Header(name string, value string) *ESRefreshServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESRefreshServiceWrapper) Headers(headers http.Header) *ESRefreshServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESRefreshServiceWrapper) Human(human bool) *ESRefreshServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESRefreshServiceWrapper) Index(index ...string) *ESRefreshServiceWrapper {
	w.obj = w.obj.Index(index...)
	return w
}

func (w *ESRefreshServiceWrapper) Pretty(pretty bool) *ESRefreshServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESReindexServiceWrapper) AbortOnVersionConflict() *ESReindexServiceWrapper {
	w.obj = w.obj.AbortOnVersionConflict()
	return w
}

func (w *ESReindexServiceWrapper) Body(body interface{}) *ESReindexServiceWrapper {
	w.obj = w.obj.Body(body)
	return w
}

func (w *ESReindexServiceWrapper) Conflicts(conflicts string) *ESReindexServiceWrapper {
	w.obj = w.obj.Conflicts(conflicts)
	return w
}

func (w *ESReindexServiceWrapper) Destination(destination *elastic.ReindexDestination) *ESReindexServiceWrapper {
	w.obj = w.obj.Destination(destination)
	return w
}

func (w *ESReindexServiceWrapper) DestinationIndex(index string) *ESReindexServiceWrapper {
	w.obj = w.obj.DestinationIndex(index)
	return w
}

func (w *ESReindexServiceWrapper) DestinationIndexAndType(index string, typ string) *ESReindexServiceWrapper {
	w.obj = w.obj.DestinationIndexAndType(index, typ)
	return w
}

func (w *ESReindexServiceWrapper) Do(ctx context.Context) (*elastic.BulkIndexByScrollResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.BulkIndexByScrollResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.ReindexService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.ReindexService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.ReindexService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.ReindexService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.ReindexService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.ReindexService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.ReindexService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESReindexServiceWrapper) DoAsync(ctx context.Context) (*elastic.StartTaskResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.StartTaskResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.ReindexService.DoAsync", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.ReindexService.DoAsync", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.ReindexService.DoAsync", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.ReindexService.DoAsync")
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
			w.inflightMetric.WithLabelValues("elastic.ReindexService.DoAsync", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.ReindexService.DoAsync", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.ReindexService.DoAsync", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DoAsync(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESReindexServiceWrapper) ErrorTrace(errorTrace bool) *ESReindexServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESReindexServiceWrapper) FilterPath(filterPath ...string) *ESReindexServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESReindexServiceWrapper) Header(name string, value string) *ESReindexServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESReindexServiceWrapper) Headers(headers http.Header) *ESReindexServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESReindexServiceWrapper) Human(human bool) *ESReindexServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESReindexServiceWrapper) Pretty(pretty bool) *ESReindexServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESReindexServiceWrapper) ProceedOnVersionConflict() *ESReindexServiceWrapper {
	w.obj = w.obj.ProceedOnVersionConflict()
	return w
}

func (w *ESReindexServiceWrapper) Refresh(refresh string) *ESReindexServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESReindexServiceWrapper) RequestsPerSecond(requestsPerSecond int) *ESReindexServiceWrapper {
	w.obj = w.obj.RequestsPerSecond(requestsPerSecond)
	return w
}

func (w *ESReindexServiceWrapper) Script(script *elastic.Script) *ESReindexServiceWrapper {
	w.obj = w.obj.Script(script)
	return w
}

func (w *ESReindexServiceWrapper) Size(size int) *ESReindexServiceWrapper {
	w.obj = w.obj.Size(size)
	return w
}

func (w *ESReindexServiceWrapper) Slices(slices interface{}) *ESReindexServiceWrapper {
	w.obj = w.obj.Slices(slices)
	return w
}

func (w *ESReindexServiceWrapper) Source(source *elastic.ReindexSource) *ESReindexServiceWrapper {
	w.obj = w.obj.Source(source)
	return w
}

func (w *ESReindexServiceWrapper) SourceIndex(index string) *ESReindexServiceWrapper {
	w.obj = w.obj.SourceIndex(index)
	return w
}

func (w *ESReindexServiceWrapper) Timeout(timeout string) *ESReindexServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESReindexServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESReindexServiceWrapper) WaitForActiveShards(waitForActiveShards string) *ESReindexServiceWrapper {
	w.obj = w.obj.WaitForActiveShards(waitForActiveShards)
	return w
}

func (w *ESReindexServiceWrapper) WaitForCompletion(waitForCompletion bool) *ESReindexServiceWrapper {
	w.obj = w.obj.WaitForCompletion(waitForCompletion)
	return w
}

func (w *ESScrollServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESScrollServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESScrollServiceWrapper) Body(body interface{}) *ESScrollServiceWrapper {
	w.obj = w.obj.Body(body)
	return w
}

func (w *ESScrollServiceWrapper) Clear(ctx context.Context) error {
	var err error
	err = w.obj.Clear(ctx)
	return err
}

func (w *ESScrollServiceWrapper) Do(ctx context.Context) (*elastic.SearchResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.SearchResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.ScrollService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.ScrollService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.ScrollService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.ScrollService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.ScrollService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.ScrollService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.ScrollService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESScrollServiceWrapper) DocvalueField(docvalueField string) *ESScrollServiceWrapper {
	w.obj = w.obj.DocvalueField(docvalueField)
	return w
}

func (w *ESScrollServiceWrapper) DocvalueFieldWithFormat(docvalueField elastic.DocvalueField) *ESScrollServiceWrapper {
	w.obj = w.obj.DocvalueFieldWithFormat(docvalueField)
	return w
}

func (w *ESScrollServiceWrapper) DocvalueFields(docvalueFields ...string) *ESScrollServiceWrapper {
	w.obj = w.obj.DocvalueFields(docvalueFields...)
	return w
}

func (w *ESScrollServiceWrapper) DocvalueFieldsWithFormat(docvalueFields ...elastic.DocvalueField) *ESScrollServiceWrapper {
	w.obj = w.obj.DocvalueFieldsWithFormat(docvalueFields...)
	return w
}

func (w *ESScrollServiceWrapper) ErrorTrace(errorTrace bool) *ESScrollServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESScrollServiceWrapper) ExpandWildcards(expandWildcards string) *ESScrollServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESScrollServiceWrapper) FetchSource(fetchSource bool) *ESScrollServiceWrapper {
	w.obj = w.obj.FetchSource(fetchSource)
	return w
}

func (w *ESScrollServiceWrapper) FetchSourceContext(fetchSourceContext *elastic.FetchSourceContext) *ESScrollServiceWrapper {
	w.obj = w.obj.FetchSourceContext(fetchSourceContext)
	return w
}

func (w *ESScrollServiceWrapper) FilterPath(filterPath ...string) *ESScrollServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESScrollServiceWrapper) Header(name string, value string) *ESScrollServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESScrollServiceWrapper) Headers(headers http.Header) *ESScrollServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESScrollServiceWrapper) Highlight(highlight *elastic.Highlight) *ESScrollServiceWrapper {
	w.obj = w.obj.Highlight(highlight)
	return w
}

func (w *ESScrollServiceWrapper) Human(human bool) *ESScrollServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESScrollServiceWrapper) IgnoreThrottled(ignoreThrottled bool) *ESScrollServiceWrapper {
	w.obj = w.obj.IgnoreThrottled(ignoreThrottled)
	return w
}

func (w *ESScrollServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESScrollServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESScrollServiceWrapper) Index(indices ...string) *ESScrollServiceWrapper {
	w.obj = w.obj.Index(indices...)
	return w
}

func (w *ESScrollServiceWrapper) KeepAlive(keepAlive string) *ESScrollServiceWrapper {
	w.obj = w.obj.KeepAlive(keepAlive)
	return w
}

func (w *ESScrollServiceWrapper) MaxResponseSize(maxResponseSize int64) *ESScrollServiceWrapper {
	w.obj = w.obj.MaxResponseSize(maxResponseSize)
	return w
}

func (w *ESScrollServiceWrapper) NoStoredFields() *ESScrollServiceWrapper {
	w.obj = w.obj.NoStoredFields()
	return w
}

func (w *ESScrollServiceWrapper) PostFilter(postFilter elastic.Query) *ESScrollServiceWrapper {
	w.obj = w.obj.PostFilter(postFilter)
	return w
}

func (w *ESScrollServiceWrapper) Preference(preference string) *ESScrollServiceWrapper {
	w.obj = w.obj.Preference(preference)
	return w
}

func (w *ESScrollServiceWrapper) Pretty(pretty bool) *ESScrollServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESScrollServiceWrapper) Query(query elastic.Query) *ESScrollServiceWrapper {
	w.obj = w.obj.Query(query)
	return w
}

func (w *ESScrollServiceWrapper) RestTotalHitsAsInt(enabled bool) *ESScrollServiceWrapper {
	w.obj = w.obj.RestTotalHitsAsInt(enabled)
	return w
}

func (w *ESScrollServiceWrapper) Retrier(retrier elastic.Retrier) *ESScrollServiceWrapper {
	w.obj = w.obj.Retrier(retrier)
	return w
}

func (w *ESScrollServiceWrapper) Routing(routings ...string) *ESScrollServiceWrapper {
	w.obj = w.obj.Routing(routings...)
	return w
}

func (w *ESScrollServiceWrapper) Scroll(keepAlive string) *ESScrollServiceWrapper {
	w.obj = w.obj.Scroll(keepAlive)
	return w
}

func (w *ESScrollServiceWrapper) ScrollId(scrollId string) *ESScrollServiceWrapper {
	w.obj = w.obj.ScrollId(scrollId)
	return w
}

func (w *ESScrollServiceWrapper) SearchSource(searchSource *elastic.SearchSource) *ESScrollServiceWrapper {
	w.obj = w.obj.SearchSource(searchSource)
	return w
}

func (w *ESScrollServiceWrapper) Size(size int) *ESScrollServiceWrapper {
	w.obj = w.obj.Size(size)
	return w
}

func (w *ESScrollServiceWrapper) Slice(sliceQuery elastic.Query) *ESScrollServiceWrapper {
	w.obj = w.obj.Slice(sliceQuery)
	return w
}

func (w *ESScrollServiceWrapper) Sort(field string, ascending bool) *ESScrollServiceWrapper {
	w.obj = w.obj.Sort(field, ascending)
	return w
}

func (w *ESScrollServiceWrapper) SortBy(sorter ...elastic.Sorter) *ESScrollServiceWrapper {
	w.obj = w.obj.SortBy(sorter...)
	return w
}

func (w *ESScrollServiceWrapper) SortWithInfo(info elastic.SortInfo) *ESScrollServiceWrapper {
	w.obj = w.obj.SortWithInfo(info)
	return w
}

func (w *ESScrollServiceWrapper) StoredField(fieldName string) *ESScrollServiceWrapper {
	w.obj = w.obj.StoredField(fieldName)
	return w
}

func (w *ESScrollServiceWrapper) StoredFields(fields ...string) *ESScrollServiceWrapper {
	w.obj = w.obj.StoredFields(fields...)
	return w
}

func (w *ESScrollServiceWrapper) TrackTotalHits(trackTotalHits interface{}) *ESScrollServiceWrapper {
	w.obj = w.obj.TrackTotalHits(trackTotalHits)
	return w
}

func (w *ESScrollServiceWrapper) Type(types ...string) *ESScrollServiceWrapper {
	w.obj = w.obj.Type(types...)
	return w
}

func (w *ESScrollServiceWrapper) Version(version bool) *ESScrollServiceWrapper {
	w.obj = w.obj.Version(version)
	return w
}

func (w *ESSearchServiceWrapper) Aggregation(name string, aggregation elastic.Aggregation) *ESSearchServiceWrapper {
	w.obj = w.obj.Aggregation(name, aggregation)
	return w
}

func (w *ESSearchServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESSearchServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESSearchServiceWrapper) AllowPartialSearchResults(enabled bool) *ESSearchServiceWrapper {
	w.obj = w.obj.AllowPartialSearchResults(enabled)
	return w
}

func (w *ESSearchServiceWrapper) BatchedReduceSize(size int) *ESSearchServiceWrapper {
	w.obj = w.obj.BatchedReduceSize(size)
	return w
}

func (w *ESSearchServiceWrapper) CCSMinimizeRoundtrips(enabled bool) *ESSearchServiceWrapper {
	w.obj = w.obj.CCSMinimizeRoundtrips(enabled)
	return w
}

func (w *ESSearchServiceWrapper) Collapse(collapse *elastic.CollapseBuilder) *ESSearchServiceWrapper {
	w.obj = w.obj.Collapse(collapse)
	return w
}

func (w *ESSearchServiceWrapper) DefaultRescoreWindowSize(defaultRescoreWindowSize int) *ESSearchServiceWrapper {
	w.obj = w.obj.DefaultRescoreWindowSize(defaultRescoreWindowSize)
	return w
}

func (w *ESSearchServiceWrapper) Do(ctx context.Context) (*elastic.SearchResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.SearchResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.SearchService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.SearchService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.SearchService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.SearchService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.SearchService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.SearchService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.SearchService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESSearchServiceWrapper) DocvalueField(docvalueField string) *ESSearchServiceWrapper {
	w.obj = w.obj.DocvalueField(docvalueField)
	return w
}

func (w *ESSearchServiceWrapper) DocvalueFieldWithFormat(docvalueField elastic.DocvalueField) *ESSearchServiceWrapper {
	w.obj = w.obj.DocvalueFieldWithFormat(docvalueField)
	return w
}

func (w *ESSearchServiceWrapper) DocvalueFields(docvalueFields ...string) *ESSearchServiceWrapper {
	w.obj = w.obj.DocvalueFields(docvalueFields...)
	return w
}

func (w *ESSearchServiceWrapper) DocvalueFieldsWithFormat(docvalueFields ...elastic.DocvalueField) *ESSearchServiceWrapper {
	w.obj = w.obj.DocvalueFieldsWithFormat(docvalueFields...)
	return w
}

func (w *ESSearchServiceWrapper) ErrorTrace(errorTrace bool) *ESSearchServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESSearchServiceWrapper) ExpandWildcards(expandWildcards string) *ESSearchServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESSearchServiceWrapper) Explain(explain bool) *ESSearchServiceWrapper {
	w.obj = w.obj.Explain(explain)
	return w
}

func (w *ESSearchServiceWrapper) FetchSource(fetchSource bool) *ESSearchServiceWrapper {
	w.obj = w.obj.FetchSource(fetchSource)
	return w
}

func (w *ESSearchServiceWrapper) FetchSourceContext(fetchSourceContext *elastic.FetchSourceContext) *ESSearchServiceWrapper {
	w.obj = w.obj.FetchSourceContext(fetchSourceContext)
	return w
}

func (w *ESSearchServiceWrapper) FilterPath(filterPath ...string) *ESSearchServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESSearchServiceWrapper) From(from int) *ESSearchServiceWrapper {
	w.obj = w.obj.From(from)
	return w
}

func (w *ESSearchServiceWrapper) GlobalSuggestText(globalText string) *ESSearchServiceWrapper {
	w.obj = w.obj.GlobalSuggestText(globalText)
	return w
}

func (w *ESSearchServiceWrapper) Header(name string, value string) *ESSearchServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESSearchServiceWrapper) Headers(headers http.Header) *ESSearchServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESSearchServiceWrapper) Highlight(highlight *elastic.Highlight) *ESSearchServiceWrapper {
	w.obj = w.obj.Highlight(highlight)
	return w
}

func (w *ESSearchServiceWrapper) Human(human bool) *ESSearchServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESSearchServiceWrapper) IgnoreThrottled(ignoreThrottled bool) *ESSearchServiceWrapper {
	w.obj = w.obj.IgnoreThrottled(ignoreThrottled)
	return w
}

func (w *ESSearchServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESSearchServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESSearchServiceWrapper) Index(index ...string) *ESSearchServiceWrapper {
	w.obj = w.obj.Index(index...)
	return w
}

func (w *ESSearchServiceWrapper) Lenient(lenient bool) *ESSearchServiceWrapper {
	w.obj = w.obj.Lenient(lenient)
	return w
}

func (w *ESSearchServiceWrapper) MaxConcurrentShardRequests(max int) *ESSearchServiceWrapper {
	w.obj = w.obj.MaxConcurrentShardRequests(max)
	return w
}

func (w *ESSearchServiceWrapper) MaxResponseSize(maxResponseSize int64) *ESSearchServiceWrapper {
	w.obj = w.obj.MaxResponseSize(maxResponseSize)
	return w
}

func (w *ESSearchServiceWrapper) MinScore(minScore float64) *ESSearchServiceWrapper {
	w.obj = w.obj.MinScore(minScore)
	return w
}

func (w *ESSearchServiceWrapper) NoStoredFields() *ESSearchServiceWrapper {
	w.obj = w.obj.NoStoredFields()
	return w
}

func (w *ESSearchServiceWrapper) PostFilter(postFilter elastic.Query) *ESSearchServiceWrapper {
	w.obj = w.obj.PostFilter(postFilter)
	return w
}

func (w *ESSearchServiceWrapper) PreFilterShardSize(threshold int) *ESSearchServiceWrapper {
	w.obj = w.obj.PreFilterShardSize(threshold)
	return w
}

func (w *ESSearchServiceWrapper) Preference(preference string) *ESSearchServiceWrapper {
	w.obj = w.obj.Preference(preference)
	return w
}

func (w *ESSearchServiceWrapper) Pretty(pretty bool) *ESSearchServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESSearchServiceWrapper) Profile(profile bool) *ESSearchServiceWrapper {
	w.obj = w.obj.Profile(profile)
	return w
}

func (w *ESSearchServiceWrapper) Query(query elastic.Query) *ESSearchServiceWrapper {
	w.obj = w.obj.Query(query)
	return w
}

func (w *ESSearchServiceWrapper) RequestCache(requestCache bool) *ESSearchServiceWrapper {
	w.obj = w.obj.RequestCache(requestCache)
	return w
}

func (w *ESSearchServiceWrapper) Rescorer(rescore *elastic.Rescore) *ESSearchServiceWrapper {
	w.obj = w.obj.Rescorer(rescore)
	return w
}

func (w *ESSearchServiceWrapper) RestTotalHitsAsInt(enabled bool) *ESSearchServiceWrapper {
	w.obj = w.obj.RestTotalHitsAsInt(enabled)
	return w
}

func (w *ESSearchServiceWrapper) Routing(routings ...string) *ESSearchServiceWrapper {
	w.obj = w.obj.Routing(routings...)
	return w
}

func (w *ESSearchServiceWrapper) SearchAfter(sortValues ...interface{}) *ESSearchServiceWrapper {
	w.obj = w.obj.SearchAfter(sortValues...)
	return w
}

func (w *ESSearchServiceWrapper) SearchSource(searchSource *elastic.SearchSource) *ESSearchServiceWrapper {
	w.obj = w.obj.SearchSource(searchSource)
	return w
}

func (w *ESSearchServiceWrapper) SearchType(searchType string) *ESSearchServiceWrapper {
	w.obj = w.obj.SearchType(searchType)
	return w
}

func (w *ESSearchServiceWrapper) SeqNoPrimaryTerm(enabled bool) *ESSearchServiceWrapper {
	w.obj = w.obj.SeqNoPrimaryTerm(enabled)
	return w
}

func (w *ESSearchServiceWrapper) Size(size int) *ESSearchServiceWrapper {
	w.obj = w.obj.Size(size)
	return w
}

func (w *ESSearchServiceWrapper) Sort(field string, ascending bool) *ESSearchServiceWrapper {
	w.obj = w.obj.Sort(field, ascending)
	return w
}

func (w *ESSearchServiceWrapper) SortBy(sorter ...elastic.Sorter) *ESSearchServiceWrapper {
	w.obj = w.obj.SortBy(sorter...)
	return w
}

func (w *ESSearchServiceWrapper) SortWithInfo(info elastic.SortInfo) *ESSearchServiceWrapper {
	w.obj = w.obj.SortWithInfo(info)
	return w
}

func (w *ESSearchServiceWrapper) Source(source interface{}) *ESSearchServiceWrapper {
	w.obj = w.obj.Source(source)
	return w
}

func (w *ESSearchServiceWrapper) StoredField(fieldName string) *ESSearchServiceWrapper {
	w.obj = w.obj.StoredField(fieldName)
	return w
}

func (w *ESSearchServiceWrapper) StoredFields(fields ...string) *ESSearchServiceWrapper {
	w.obj = w.obj.StoredFields(fields...)
	return w
}

func (w *ESSearchServiceWrapper) Suggester(suggester elastic.Suggester) *ESSearchServiceWrapper {
	w.obj = w.obj.Suggester(suggester)
	return w
}

func (w *ESSearchServiceWrapper) TerminateAfter(terminateAfter int) *ESSearchServiceWrapper {
	w.obj = w.obj.TerminateAfter(terminateAfter)
	return w
}

func (w *ESSearchServiceWrapper) Timeout(timeout string) *ESSearchServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESSearchServiceWrapper) TimeoutInMillis(timeoutInMillis int) *ESSearchServiceWrapper {
	w.obj = w.obj.TimeoutInMillis(timeoutInMillis)
	return w
}

func (w *ESSearchServiceWrapper) TrackScores(trackScores bool) *ESSearchServiceWrapper {
	w.obj = w.obj.TrackScores(trackScores)
	return w
}

func (w *ESSearchServiceWrapper) TrackTotalHits(trackTotalHits interface{}) *ESSearchServiceWrapper {
	w.obj = w.obj.TrackTotalHits(trackTotalHits)
	return w
}

func (w *ESSearchServiceWrapper) Type(typ ...string) *ESSearchServiceWrapper {
	w.obj = w.obj.Type(typ...)
	return w
}

func (w *ESSearchServiceWrapper) TypedKeys(enabled bool) *ESSearchServiceWrapper {
	w.obj = w.obj.TypedKeys(enabled)
	return w
}

func (w *ESSearchServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESSearchServiceWrapper) Version(version bool) *ESSearchServiceWrapper {
	w.obj = w.obj.Version(version)
	return w
}

func (w *ESSearchShardsServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESSearchShardsServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESSearchShardsServiceWrapper) Do(ctx context.Context) (*elastic.SearchShardsResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.SearchShardsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.SearchShardsService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.SearchShardsService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.SearchShardsService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.SearchShardsService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.SearchShardsService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.SearchShardsService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.SearchShardsService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESSearchShardsServiceWrapper) ErrorTrace(errorTrace bool) *ESSearchShardsServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESSearchShardsServiceWrapper) ExpandWildcards(expandWildcards string) *ESSearchShardsServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESSearchShardsServiceWrapper) FilterPath(filterPath ...string) *ESSearchShardsServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESSearchShardsServiceWrapper) Header(name string, value string) *ESSearchShardsServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESSearchShardsServiceWrapper) Headers(headers http.Header) *ESSearchShardsServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESSearchShardsServiceWrapper) Human(human bool) *ESSearchShardsServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESSearchShardsServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESSearchShardsServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESSearchShardsServiceWrapper) Index(index ...string) *ESSearchShardsServiceWrapper {
	w.obj = w.obj.Index(index...)
	return w
}

func (w *ESSearchShardsServiceWrapper) Local(local bool) *ESSearchShardsServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESSearchShardsServiceWrapper) Preference(preference string) *ESSearchShardsServiceWrapper {
	w.obj = w.obj.Preference(preference)
	return w
}

func (w *ESSearchShardsServiceWrapper) Pretty(pretty bool) *ESSearchShardsServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESSearchShardsServiceWrapper) Routing(routing string) *ESSearchShardsServiceWrapper {
	w.obj = w.obj.Routing(routing)
	return w
}

func (w *ESSearchShardsServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) BodyJson(body interface{}) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) BodyString(body string) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) Do(ctx context.Context) (*elastic.SnapshotCreateRepositoryResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.SnapshotCreateRepositoryResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.SnapshotCreateRepositoryService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.SnapshotCreateRepositoryService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.SnapshotCreateRepositoryService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.SnapshotCreateRepositoryService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.SnapshotCreateRepositoryService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.SnapshotCreateRepositoryService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.SnapshotCreateRepositoryService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) ErrorTrace(errorTrace bool) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) FilterPath(filterPath ...string) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) Header(name string, value string) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) Headers(headers http.Header) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) Human(human bool) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) MasterTimeout(masterTimeout string) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) Pretty(pretty bool) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) Repository(repository string) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.Repository(repository)
	return w
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) Setting(name string, value interface{}) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.Setting(name, value)
	return w
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) Settings(settings map[string]interface{}) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.Settings(settings)
	return w
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) Timeout(timeout string) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) Type(typ string) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.Type(typ)
	return w
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESSnapshotCreateRepositoryServiceWrapper) Verify(verify bool) *ESSnapshotCreateRepositoryServiceWrapper {
	w.obj = w.obj.Verify(verify)
	return w
}

func (w *ESSnapshotCreateServiceWrapper) BodyJson(body interface{}) *ESSnapshotCreateServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESSnapshotCreateServiceWrapper) BodyString(body string) *ESSnapshotCreateServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESSnapshotCreateServiceWrapper) Do(ctx context.Context) (*elastic.SnapshotCreateResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.SnapshotCreateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.SnapshotCreateService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.SnapshotCreateService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.SnapshotCreateService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.SnapshotCreateService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.SnapshotCreateService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.SnapshotCreateService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.SnapshotCreateService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESSnapshotCreateServiceWrapper) ErrorTrace(errorTrace bool) *ESSnapshotCreateServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESSnapshotCreateServiceWrapper) FilterPath(filterPath ...string) *ESSnapshotCreateServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESSnapshotCreateServiceWrapper) Header(name string, value string) *ESSnapshotCreateServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESSnapshotCreateServiceWrapper) Headers(headers http.Header) *ESSnapshotCreateServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESSnapshotCreateServiceWrapper) Human(human bool) *ESSnapshotCreateServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESSnapshotCreateServiceWrapper) MasterTimeout(masterTimeout string) *ESSnapshotCreateServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESSnapshotCreateServiceWrapper) Pretty(pretty bool) *ESSnapshotCreateServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESSnapshotCreateServiceWrapper) Repository(repository string) *ESSnapshotCreateServiceWrapper {
	w.obj = w.obj.Repository(repository)
	return w
}

func (w *ESSnapshotCreateServiceWrapper) Snapshot(snapshot string) *ESSnapshotCreateServiceWrapper {
	w.obj = w.obj.Snapshot(snapshot)
	return w
}

func (w *ESSnapshotCreateServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESSnapshotCreateServiceWrapper) WaitForCompletion(waitForCompletion bool) *ESSnapshotCreateServiceWrapper {
	w.obj = w.obj.WaitForCompletion(waitForCompletion)
	return w
}

func (w *ESSnapshotDeleteRepositoryServiceWrapper) Do(ctx context.Context) (*elastic.SnapshotDeleteRepositoryResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.SnapshotDeleteRepositoryResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.SnapshotDeleteRepositoryService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.SnapshotDeleteRepositoryService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.SnapshotDeleteRepositoryService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.SnapshotDeleteRepositoryService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.SnapshotDeleteRepositoryService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.SnapshotDeleteRepositoryService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.SnapshotDeleteRepositoryService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESSnapshotDeleteRepositoryServiceWrapper) ErrorTrace(errorTrace bool) *ESSnapshotDeleteRepositoryServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESSnapshotDeleteRepositoryServiceWrapper) FilterPath(filterPath ...string) *ESSnapshotDeleteRepositoryServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESSnapshotDeleteRepositoryServiceWrapper) Header(name string, value string) *ESSnapshotDeleteRepositoryServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESSnapshotDeleteRepositoryServiceWrapper) Headers(headers http.Header) *ESSnapshotDeleteRepositoryServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESSnapshotDeleteRepositoryServiceWrapper) Human(human bool) *ESSnapshotDeleteRepositoryServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESSnapshotDeleteRepositoryServiceWrapper) MasterTimeout(masterTimeout string) *ESSnapshotDeleteRepositoryServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESSnapshotDeleteRepositoryServiceWrapper) Pretty(pretty bool) *ESSnapshotDeleteRepositoryServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESSnapshotDeleteRepositoryServiceWrapper) Repository(repositories ...string) *ESSnapshotDeleteRepositoryServiceWrapper {
	w.obj = w.obj.Repository(repositories...)
	return w
}

func (w *ESSnapshotDeleteRepositoryServiceWrapper) Timeout(timeout string) *ESSnapshotDeleteRepositoryServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESSnapshotDeleteRepositoryServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESSnapshotDeleteServiceWrapper) Do(ctx context.Context) (*elastic.SnapshotDeleteResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.SnapshotDeleteResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.SnapshotDeleteService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.SnapshotDeleteService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.SnapshotDeleteService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.SnapshotDeleteService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.SnapshotDeleteService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.SnapshotDeleteService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.SnapshotDeleteService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESSnapshotDeleteServiceWrapper) ErrorTrace(errorTrace bool) *ESSnapshotDeleteServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESSnapshotDeleteServiceWrapper) FilterPath(filterPath ...string) *ESSnapshotDeleteServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESSnapshotDeleteServiceWrapper) Header(name string, value string) *ESSnapshotDeleteServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESSnapshotDeleteServiceWrapper) Headers(headers http.Header) *ESSnapshotDeleteServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESSnapshotDeleteServiceWrapper) Human(human bool) *ESSnapshotDeleteServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESSnapshotDeleteServiceWrapper) Pretty(pretty bool) *ESSnapshotDeleteServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESSnapshotDeleteServiceWrapper) Repository(repository string) *ESSnapshotDeleteServiceWrapper {
	w.obj = w.obj.Repository(repository)
	return w
}

func (w *ESSnapshotDeleteServiceWrapper) Snapshot(snapshot string) *ESSnapshotDeleteServiceWrapper {
	w.obj = w.obj.Snapshot(snapshot)
	return w
}

func (w *ESSnapshotDeleteServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESSnapshotGetRepositoryServiceWrapper) Do(ctx context.Context) (elastic.SnapshotGetRepositoryResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 elastic.SnapshotGetRepositoryResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.SnapshotGetRepositoryService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.SnapshotGetRepositoryService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.SnapshotGetRepositoryService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.SnapshotGetRepositoryService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.SnapshotGetRepositoryService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.SnapshotGetRepositoryService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.SnapshotGetRepositoryService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESSnapshotGetRepositoryServiceWrapper) ErrorTrace(errorTrace bool) *ESSnapshotGetRepositoryServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESSnapshotGetRepositoryServiceWrapper) FilterPath(filterPath ...string) *ESSnapshotGetRepositoryServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESSnapshotGetRepositoryServiceWrapper) Header(name string, value string) *ESSnapshotGetRepositoryServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESSnapshotGetRepositoryServiceWrapper) Headers(headers http.Header) *ESSnapshotGetRepositoryServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESSnapshotGetRepositoryServiceWrapper) Human(human bool) *ESSnapshotGetRepositoryServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESSnapshotGetRepositoryServiceWrapper) Local(local bool) *ESSnapshotGetRepositoryServiceWrapper {
	w.obj = w.obj.Local(local)
	return w
}

func (w *ESSnapshotGetRepositoryServiceWrapper) MasterTimeout(masterTimeout string) *ESSnapshotGetRepositoryServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESSnapshotGetRepositoryServiceWrapper) Pretty(pretty bool) *ESSnapshotGetRepositoryServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESSnapshotGetRepositoryServiceWrapper) Repository(repositories ...string) *ESSnapshotGetRepositoryServiceWrapper {
	w.obj = w.obj.Repository(repositories...)
	return w
}

func (w *ESSnapshotGetRepositoryServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESSnapshotGetServiceWrapper) Do(ctx context.Context) (*elastic.SnapshotGetResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.SnapshotGetResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.SnapshotGetService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.SnapshotGetService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.SnapshotGetService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.SnapshotGetService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.SnapshotGetService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.SnapshotGetService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.SnapshotGetService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESSnapshotGetServiceWrapper) ErrorTrace(errorTrace bool) *ESSnapshotGetServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESSnapshotGetServiceWrapper) FilterPath(filterPath ...string) *ESSnapshotGetServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESSnapshotGetServiceWrapper) Header(name string, value string) *ESSnapshotGetServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESSnapshotGetServiceWrapper) Headers(headers http.Header) *ESSnapshotGetServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESSnapshotGetServiceWrapper) Human(human bool) *ESSnapshotGetServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESSnapshotGetServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESSnapshotGetServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESSnapshotGetServiceWrapper) MasterTimeout(masterTimeout string) *ESSnapshotGetServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESSnapshotGetServiceWrapper) Pretty(pretty bool) *ESSnapshotGetServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESSnapshotGetServiceWrapper) Repository(repository string) *ESSnapshotGetServiceWrapper {
	w.obj = w.obj.Repository(repository)
	return w
}

func (w *ESSnapshotGetServiceWrapper) Snapshot(snapshots ...string) *ESSnapshotGetServiceWrapper {
	w.obj = w.obj.Snapshot(snapshots...)
	return w
}

func (w *ESSnapshotGetServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESSnapshotGetServiceWrapper) Verbose(verbose bool) *ESSnapshotGetServiceWrapper {
	w.obj = w.obj.Verbose(verbose)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) BodyString(body string) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) Do(ctx context.Context) (*elastic.SnapshotRestoreResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.SnapshotRestoreResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.SnapshotRestoreService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.SnapshotRestoreService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.SnapshotRestoreService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.SnapshotRestoreService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.SnapshotRestoreService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.SnapshotRestoreService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.SnapshotRestoreService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESSnapshotRestoreServiceWrapper) ErrorTrace(errorTrace bool) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) FilterPath(filterPath ...string) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) Header(name string, value string) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) Headers(headers http.Header) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) Human(human bool) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) IncludeAliases(includeAliases bool) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.IncludeAliases(includeAliases)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) IncludeGlobalState(includeGlobalState bool) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.IncludeGlobalState(includeGlobalState)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) IndexSettings(indexSettings map[string]interface{}) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.IndexSettings(indexSettings)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) Indices(indices ...string) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.Indices(indices...)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) MasterTimeout(masterTimeout string) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) Partial(partial bool) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.Partial(partial)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) Pretty(pretty bool) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) RenamePattern(renamePattern string) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.RenamePattern(renamePattern)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) RenameReplacement(renameReplacement string) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.RenameReplacement(renameReplacement)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) Repository(repository string) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.Repository(repository)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) Snapshot(snapshot string) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.Snapshot(snapshot)
	return w
}

func (w *ESSnapshotRestoreServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESSnapshotRestoreServiceWrapper) WaitForCompletion(waitForCompletion bool) *ESSnapshotRestoreServiceWrapper {
	w.obj = w.obj.WaitForCompletion(waitForCompletion)
	return w
}

func (w *ESSnapshotStatusServiceWrapper) Do(ctx context.Context) (*elastic.SnapshotStatusResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.SnapshotStatusResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.SnapshotStatusService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.SnapshotStatusService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.SnapshotStatusService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.SnapshotStatusService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.SnapshotStatusService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.SnapshotStatusService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.SnapshotStatusService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESSnapshotStatusServiceWrapper) ErrorTrace(errorTrace bool) *ESSnapshotStatusServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESSnapshotStatusServiceWrapper) FilterPath(filterPath ...string) *ESSnapshotStatusServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESSnapshotStatusServiceWrapper) Header(name string, value string) *ESSnapshotStatusServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESSnapshotStatusServiceWrapper) Headers(headers http.Header) *ESSnapshotStatusServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESSnapshotStatusServiceWrapper) Human(human bool) *ESSnapshotStatusServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESSnapshotStatusServiceWrapper) MasterTimeout(masterTimeout string) *ESSnapshotStatusServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESSnapshotStatusServiceWrapper) Pretty(pretty bool) *ESSnapshotStatusServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESSnapshotStatusServiceWrapper) Repository(repository string) *ESSnapshotStatusServiceWrapper {
	w.obj = w.obj.Repository(repository)
	return w
}

func (w *ESSnapshotStatusServiceWrapper) Snapshot(snapshots ...string) *ESSnapshotStatusServiceWrapper {
	w.obj = w.obj.Snapshot(snapshots...)
	return w
}

func (w *ESSnapshotStatusServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESSnapshotVerifyRepositoryServiceWrapper) Do(ctx context.Context) (*elastic.SnapshotVerifyRepositoryResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.SnapshotVerifyRepositoryResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.SnapshotVerifyRepositoryService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.SnapshotVerifyRepositoryService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.SnapshotVerifyRepositoryService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.SnapshotVerifyRepositoryService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.SnapshotVerifyRepositoryService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.SnapshotVerifyRepositoryService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.SnapshotVerifyRepositoryService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESSnapshotVerifyRepositoryServiceWrapper) ErrorTrace(errorTrace bool) *ESSnapshotVerifyRepositoryServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESSnapshotVerifyRepositoryServiceWrapper) FilterPath(filterPath ...string) *ESSnapshotVerifyRepositoryServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESSnapshotVerifyRepositoryServiceWrapper) Header(name string, value string) *ESSnapshotVerifyRepositoryServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESSnapshotVerifyRepositoryServiceWrapper) Headers(headers http.Header) *ESSnapshotVerifyRepositoryServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESSnapshotVerifyRepositoryServiceWrapper) Human(human bool) *ESSnapshotVerifyRepositoryServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESSnapshotVerifyRepositoryServiceWrapper) MasterTimeout(masterTimeout string) *ESSnapshotVerifyRepositoryServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESSnapshotVerifyRepositoryServiceWrapper) Pretty(pretty bool) *ESSnapshotVerifyRepositoryServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESSnapshotVerifyRepositoryServiceWrapper) Repository(repository string) *ESSnapshotVerifyRepositoryServiceWrapper {
	w.obj = w.obj.Repository(repository)
	return w
}

func (w *ESSnapshotVerifyRepositoryServiceWrapper) Timeout(timeout string) *ESSnapshotVerifyRepositoryServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESSnapshotVerifyRepositoryServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESTasksCancelServiceWrapper) Actions(actions ...string) *ESTasksCancelServiceWrapper {
	w.obj = w.obj.Actions(actions...)
	return w
}

func (w *ESTasksCancelServiceWrapper) Do(ctx context.Context) (*elastic.TasksListResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.TasksListResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TasksCancelService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TasksCancelService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TasksCancelService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.TasksCancelService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.TasksCancelService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.TasksCancelService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.TasksCancelService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESTasksCancelServiceWrapper) ErrorTrace(errorTrace bool) *ESTasksCancelServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESTasksCancelServiceWrapper) FilterPath(filterPath ...string) *ESTasksCancelServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESTasksCancelServiceWrapper) Header(name string, value string) *ESTasksCancelServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESTasksCancelServiceWrapper) Headers(headers http.Header) *ESTasksCancelServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESTasksCancelServiceWrapper) Human(human bool) *ESTasksCancelServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESTasksCancelServiceWrapper) NodeId(nodeId ...string) *ESTasksCancelServiceWrapper {
	w.obj = w.obj.NodeId(nodeId...)
	return w
}

func (w *ESTasksCancelServiceWrapper) ParentTaskId(parentTaskId string) *ESTasksCancelServiceWrapper {
	w.obj = w.obj.ParentTaskId(parentTaskId)
	return w
}

func (w *ESTasksCancelServiceWrapper) ParentTaskIdFromNodeAndId(nodeId string, id int64) *ESTasksCancelServiceWrapper {
	w.obj = w.obj.ParentTaskIdFromNodeAndId(nodeId, id)
	return w
}

func (w *ESTasksCancelServiceWrapper) Pretty(pretty bool) *ESTasksCancelServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESTasksCancelServiceWrapper) TaskId(taskId string) *ESTasksCancelServiceWrapper {
	w.obj = w.obj.TaskId(taskId)
	return w
}

func (w *ESTasksCancelServiceWrapper) TaskIdFromNodeAndId(nodeId string, id int64) *ESTasksCancelServiceWrapper {
	w.obj = w.obj.TaskIdFromNodeAndId(nodeId, id)
	return w
}

func (w *ESTasksCancelServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESTasksGetTaskServiceWrapper) Do(ctx context.Context) (*elastic.TasksGetTaskResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.TasksGetTaskResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TasksGetTaskService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TasksGetTaskService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TasksGetTaskService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.TasksGetTaskService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.TasksGetTaskService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.TasksGetTaskService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.TasksGetTaskService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESTasksGetTaskServiceWrapper) ErrorTrace(errorTrace bool) *ESTasksGetTaskServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESTasksGetTaskServiceWrapper) FilterPath(filterPath ...string) *ESTasksGetTaskServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESTasksGetTaskServiceWrapper) Header(name string, value string) *ESTasksGetTaskServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESTasksGetTaskServiceWrapper) Headers(headers http.Header) *ESTasksGetTaskServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESTasksGetTaskServiceWrapper) Human(human bool) *ESTasksGetTaskServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESTasksGetTaskServiceWrapper) Pretty(pretty bool) *ESTasksGetTaskServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESTasksGetTaskServiceWrapper) TaskId(taskId string) *ESTasksGetTaskServiceWrapper {
	w.obj = w.obj.TaskId(taskId)
	return w
}

func (w *ESTasksGetTaskServiceWrapper) TaskIdFromNodeAndId(nodeId string, id int64) *ESTasksGetTaskServiceWrapper {
	w.obj = w.obj.TaskIdFromNodeAndId(nodeId, id)
	return w
}

func (w *ESTasksGetTaskServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESTasksGetTaskServiceWrapper) WaitForCompletion(waitForCompletion bool) *ESTasksGetTaskServiceWrapper {
	w.obj = w.obj.WaitForCompletion(waitForCompletion)
	return w
}

func (w *ESTasksListServiceWrapper) Actions(actions ...string) *ESTasksListServiceWrapper {
	w.obj = w.obj.Actions(actions...)
	return w
}

func (w *ESTasksListServiceWrapper) Detailed(detailed bool) *ESTasksListServiceWrapper {
	w.obj = w.obj.Detailed(detailed)
	return w
}

func (w *ESTasksListServiceWrapper) Do(ctx context.Context) (*elastic.TasksListResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.TasksListResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TasksListService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TasksListService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TasksListService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.TasksListService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.TasksListService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.TasksListService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.TasksListService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESTasksListServiceWrapper) ErrorTrace(errorTrace bool) *ESTasksListServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESTasksListServiceWrapper) FilterPath(filterPath ...string) *ESTasksListServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESTasksListServiceWrapper) GroupBy(groupBy string) *ESTasksListServiceWrapper {
	w.obj = w.obj.GroupBy(groupBy)
	return w
}

func (w *ESTasksListServiceWrapper) Header(name string, value string) *ESTasksListServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESTasksListServiceWrapper) Headers(headers http.Header) *ESTasksListServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESTasksListServiceWrapper) Human(human bool) *ESTasksListServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESTasksListServiceWrapper) NodeId(nodeId ...string) *ESTasksListServiceWrapper {
	w.obj = w.obj.NodeId(nodeId...)
	return w
}

func (w *ESTasksListServiceWrapper) ParentTaskId(parentTaskId string) *ESTasksListServiceWrapper {
	w.obj = w.obj.ParentTaskId(parentTaskId)
	return w
}

func (w *ESTasksListServiceWrapper) Pretty(pretty bool) *ESTasksListServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESTasksListServiceWrapper) TaskId(taskId ...string) *ESTasksListServiceWrapper {
	w.obj = w.obj.TaskId(taskId...)
	return w
}

func (w *ESTasksListServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESTasksListServiceWrapper) WaitForCompletion(waitForCompletion bool) *ESTasksListServiceWrapper {
	w.obj = w.obj.WaitForCompletion(waitForCompletion)
	return w
}

func (w *ESTermvectorsServiceWrapper) BodyJson(body interface{}) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESTermvectorsServiceWrapper) BodyString(body string) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESTermvectorsServiceWrapper) Dfs(dfs bool) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Dfs(dfs)
	return w
}

func (w *ESTermvectorsServiceWrapper) Do(ctx context.Context) (*elastic.TermvectorsResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.TermvectorsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.TermvectorsService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.TermvectorsService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.TermvectorsService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.TermvectorsService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.TermvectorsService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.TermvectorsService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.TermvectorsService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESTermvectorsServiceWrapper) Doc(doc interface{}) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Doc(doc)
	return w
}

func (w *ESTermvectorsServiceWrapper) ErrorTrace(errorTrace bool) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESTermvectorsServiceWrapper) FieldStatistics(fieldStatistics bool) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.FieldStatistics(fieldStatistics)
	return w
}

func (w *ESTermvectorsServiceWrapper) Fields(fields ...string) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Fields(fields...)
	return w
}

func (w *ESTermvectorsServiceWrapper) Filter(filter *elastic.TermvectorsFilterSettings) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Filter(filter)
	return w
}

func (w *ESTermvectorsServiceWrapper) FilterPath(filterPath ...string) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESTermvectorsServiceWrapper) Header(name string, value string) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESTermvectorsServiceWrapper) Headers(headers http.Header) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESTermvectorsServiceWrapper) Human(human bool) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESTermvectorsServiceWrapper) Id(id string) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESTermvectorsServiceWrapper) Index(index string) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Index(index)
	return w
}

func (w *ESTermvectorsServiceWrapper) Offsets(offsets bool) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Offsets(offsets)
	return w
}

func (w *ESTermvectorsServiceWrapper) Parent(parent string) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Parent(parent)
	return w
}

func (w *ESTermvectorsServiceWrapper) Payloads(payloads bool) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Payloads(payloads)
	return w
}

func (w *ESTermvectorsServiceWrapper) PerFieldAnalyzer(perFieldAnalyzer map[string]string) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.PerFieldAnalyzer(perFieldAnalyzer)
	return w
}

func (w *ESTermvectorsServiceWrapper) Positions(positions bool) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Positions(positions)
	return w
}

func (w *ESTermvectorsServiceWrapper) Preference(preference string) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Preference(preference)
	return w
}

func (w *ESTermvectorsServiceWrapper) Pretty(pretty bool) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESTermvectorsServiceWrapper) Realtime(realtime bool) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Realtime(realtime)
	return w
}

func (w *ESTermvectorsServiceWrapper) Routing(routing string) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Routing(routing)
	return w
}

func (w *ESTermvectorsServiceWrapper) TermStatistics(termStatistics bool) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.TermStatistics(termStatistics)
	return w
}

func (w *ESTermvectorsServiceWrapper) Type(typ string) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Type(typ)
	return w
}

func (w *ESTermvectorsServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESTermvectorsServiceWrapper) Version(version interface{}) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.Version(version)
	return w
}

func (w *ESTermvectorsServiceWrapper) VersionType(versionType string) *ESTermvectorsServiceWrapper {
	w.obj = w.obj.VersionType(versionType)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) AbortOnVersionConflict() *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.AbortOnVersionConflict()
	return w
}

func (w *ESUpdateByQueryServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) AnalyzeWildcard(analyzeWildcard bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.AnalyzeWildcard(analyzeWildcard)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Analyzer(analyzer string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Analyzer(analyzer)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Body(body string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Body(body)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Conflicts(conflicts string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Conflicts(conflicts)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) DF(df string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.DF(df)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) DefaultOperator(defaultOperator string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.DefaultOperator(defaultOperator)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Do(ctx context.Context) (*elastic.BulkIndexByScrollResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.BulkIndexByScrollResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.UpdateByQueryService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.UpdateByQueryService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.UpdateByQueryService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.UpdateByQueryService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.UpdateByQueryService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.UpdateByQueryService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.UpdateByQueryService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESUpdateByQueryServiceWrapper) DoAsync(ctx context.Context) (*elastic.StartTaskResult, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.StartTaskResult
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.UpdateByQueryService.DoAsync", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.UpdateByQueryService.DoAsync", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.UpdateByQueryService.DoAsync", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.UpdateByQueryService.DoAsync")
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
			w.inflightMetric.WithLabelValues("elastic.UpdateByQueryService.DoAsync", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.UpdateByQueryService.DoAsync", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.UpdateByQueryService.DoAsync", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.DoAsync(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESUpdateByQueryServiceWrapper) DocvalueFields(docvalueFields ...string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.DocvalueFields(docvalueFields...)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) ErrorTrace(errorTrace bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) ExpandWildcards(expandWildcards string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Explain(explain bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Explain(explain)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) FielddataFields(fielddataFields ...string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.FielddataFields(fielddataFields...)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) FilterPath(filterPath ...string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) From(from int) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.From(from)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Header(name string, value string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Headers(headers http.Header) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Human(human bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Index(index ...string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Index(index...)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Lenient(lenient bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Lenient(lenient)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) LowercaseExpandedTerms(lowercaseExpandedTerms bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.LowercaseExpandedTerms(lowercaseExpandedTerms)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Pipeline(pipeline string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Pipeline(pipeline)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Preference(preference string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Preference(preference)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Pretty(pretty bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) ProceedOnVersionConflict() *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.ProceedOnVersionConflict()
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Q(q string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Q(q)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Query(query elastic.Query) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Query(query)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Refresh(refresh string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) RequestCache(requestCache bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.RequestCache(requestCache)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) RequestsPerSecond(requestsPerSecond int) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.RequestsPerSecond(requestsPerSecond)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Routing(routing ...string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Routing(routing...)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Script(script *elastic.Script) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Script(script)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Scroll(scroll string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Scroll(scroll)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) ScrollSize(scrollSize int) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.ScrollSize(scrollSize)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) SearchTimeout(searchTimeout string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.SearchTimeout(searchTimeout)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) SearchType(searchType string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.SearchType(searchType)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Size(size int) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Size(size)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Slices(slices interface{}) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Slices(slices)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Sort(sort ...string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Sort(sort...)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) SortByField(field string, ascending bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.SortByField(field, ascending)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Stats(stats ...string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Stats(stats...)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) StoredFields(storedFields ...string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.StoredFields(storedFields...)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) SuggestField(suggestField string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.SuggestField(suggestField)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) SuggestMode(suggestMode string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.SuggestMode(suggestMode)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) SuggestSize(suggestSize int) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.SuggestSize(suggestSize)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) SuggestText(suggestText string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.SuggestText(suggestText)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) TerminateAfter(terminateAfter int) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.TerminateAfter(terminateAfter)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Timeout(timeout string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) TimeoutInMillis(timeoutInMillis int) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.TimeoutInMillis(timeoutInMillis)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) TrackScores(trackScores bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.TrackScores(trackScores)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Type(typ ...string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Type(typ...)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESUpdateByQueryServiceWrapper) Version(version bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.Version(version)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) VersionType(versionType bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.VersionType(versionType)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) WaitForActiveShards(waitForActiveShards string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.WaitForActiveShards(waitForActiveShards)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) WaitForCompletion(waitForCompletion bool) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.WaitForCompletion(waitForCompletion)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) XSource(xSource ...string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.XSource(xSource...)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) XSourceExclude(xSourceExclude ...string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.XSourceExclude(xSourceExclude...)
	return w
}

func (w *ESUpdateByQueryServiceWrapper) XSourceInclude(xSourceInclude ...string) *ESUpdateByQueryServiceWrapper {
	w.obj = w.obj.XSourceInclude(xSourceInclude...)
	return w
}

func (w *ESUpdateServiceWrapper) DetectNoop(detectNoop bool) *ESUpdateServiceWrapper {
	w.obj = w.obj.DetectNoop(detectNoop)
	return w
}

func (w *ESUpdateServiceWrapper) Do(ctx context.Context) (*elastic.UpdateResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.UpdateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.UpdateService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.UpdateService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.UpdateService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.UpdateService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.UpdateService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.UpdateService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.UpdateService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESUpdateServiceWrapper) Doc(doc interface{}) *ESUpdateServiceWrapper {
	w.obj = w.obj.Doc(doc)
	return w
}

func (w *ESUpdateServiceWrapper) DocAsUpsert(docAsUpsert bool) *ESUpdateServiceWrapper {
	w.obj = w.obj.DocAsUpsert(docAsUpsert)
	return w
}

func (w *ESUpdateServiceWrapper) ErrorTrace(errorTrace bool) *ESUpdateServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESUpdateServiceWrapper) FetchSource(fetchSource bool) *ESUpdateServiceWrapper {
	w.obj = w.obj.FetchSource(fetchSource)
	return w
}

func (w *ESUpdateServiceWrapper) FetchSourceContext(fetchSourceContext *elastic.FetchSourceContext) *ESUpdateServiceWrapper {
	w.obj = w.obj.FetchSourceContext(fetchSourceContext)
	return w
}

func (w *ESUpdateServiceWrapper) Fields(fields ...string) *ESUpdateServiceWrapper {
	w.obj = w.obj.Fields(fields...)
	return w
}

func (w *ESUpdateServiceWrapper) FilterPath(filterPath ...string) *ESUpdateServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESUpdateServiceWrapper) Header(name string, value string) *ESUpdateServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESUpdateServiceWrapper) Headers(headers http.Header) *ESUpdateServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESUpdateServiceWrapper) Human(human bool) *ESUpdateServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESUpdateServiceWrapper) Id(id string) *ESUpdateServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESUpdateServiceWrapper) IfPrimaryTerm(primaryTerm int64) *ESUpdateServiceWrapper {
	w.obj = w.obj.IfPrimaryTerm(primaryTerm)
	return w
}

func (w *ESUpdateServiceWrapper) IfSeqNo(seqNo int64) *ESUpdateServiceWrapper {
	w.obj = w.obj.IfSeqNo(seqNo)
	return w
}

func (w *ESUpdateServiceWrapper) Index(name string) *ESUpdateServiceWrapper {
	w.obj = w.obj.Index(name)
	return w
}

func (w *ESUpdateServiceWrapper) Parent(parent string) *ESUpdateServiceWrapper {
	w.obj = w.obj.Parent(parent)
	return w
}

func (w *ESUpdateServiceWrapper) Pretty(pretty bool) *ESUpdateServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESUpdateServiceWrapper) Refresh(refresh string) *ESUpdateServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESUpdateServiceWrapper) RetryOnConflict(retryOnConflict int) *ESUpdateServiceWrapper {
	w.obj = w.obj.RetryOnConflict(retryOnConflict)
	return w
}

func (w *ESUpdateServiceWrapper) Routing(routing string) *ESUpdateServiceWrapper {
	w.obj = w.obj.Routing(routing)
	return w
}

func (w *ESUpdateServiceWrapper) Script(script *elastic.Script) *ESUpdateServiceWrapper {
	w.obj = w.obj.Script(script)
	return w
}

func (w *ESUpdateServiceWrapper) ScriptedUpsert(scriptedUpsert bool) *ESUpdateServiceWrapper {
	w.obj = w.obj.ScriptedUpsert(scriptedUpsert)
	return w
}

func (w *ESUpdateServiceWrapper) Timeout(timeout string) *ESUpdateServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESUpdateServiceWrapper) Type(typ string) *ESUpdateServiceWrapper {
	w.obj = w.obj.Type(typ)
	return w
}

func (w *ESUpdateServiceWrapper) Upsert(doc interface{}) *ESUpdateServiceWrapper {
	w.obj = w.obj.Upsert(doc)
	return w
}

func (w *ESUpdateServiceWrapper) Version(version int64) *ESUpdateServiceWrapper {
	w.obj = w.obj.Version(version)
	return w
}

func (w *ESUpdateServiceWrapper) VersionType(versionType string) *ESUpdateServiceWrapper {
	w.obj = w.obj.VersionType(versionType)
	return w
}

func (w *ESUpdateServiceWrapper) WaitForActiveShards(waitForActiveShards string) *ESUpdateServiceWrapper {
	w.obj = w.obj.WaitForActiveShards(waitForActiveShards)
	return w
}

func (w *ESValidateServiceWrapper) AllShards(allShards *bool) *ESValidateServiceWrapper {
	w.obj = w.obj.AllShards(allShards)
	return w
}

func (w *ESValidateServiceWrapper) AllowNoIndices(allowNoIndices bool) *ESValidateServiceWrapper {
	w.obj = w.obj.AllowNoIndices(allowNoIndices)
	return w
}

func (w *ESValidateServiceWrapper) AnalyzeWildcard(analyzeWildcard bool) *ESValidateServiceWrapper {
	w.obj = w.obj.AnalyzeWildcard(analyzeWildcard)
	return w
}

func (w *ESValidateServiceWrapper) Analyzer(analyzer string) *ESValidateServiceWrapper {
	w.obj = w.obj.Analyzer(analyzer)
	return w
}

func (w *ESValidateServiceWrapper) BodyJson(body interface{}) *ESValidateServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESValidateServiceWrapper) BodyString(body string) *ESValidateServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESValidateServiceWrapper) DefaultOperator(defaultOperator string) *ESValidateServiceWrapper {
	w.obj = w.obj.DefaultOperator(defaultOperator)
	return w
}

func (w *ESValidateServiceWrapper) Df(df string) *ESValidateServiceWrapper {
	w.obj = w.obj.Df(df)
	return w
}

func (w *ESValidateServiceWrapper) Do(ctx context.Context) (*elastic.ValidateResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.ValidateResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.ValidateService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.ValidateService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.ValidateService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.ValidateService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.ValidateService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.ValidateService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.ValidateService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESValidateServiceWrapper) ErrorTrace(errorTrace bool) *ESValidateServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESValidateServiceWrapper) ExpandWildcards(expandWildcards string) *ESValidateServiceWrapper {
	w.obj = w.obj.ExpandWildcards(expandWildcards)
	return w
}

func (w *ESValidateServiceWrapper) Explain(explain *bool) *ESValidateServiceWrapper {
	w.obj = w.obj.Explain(explain)
	return w
}

func (w *ESValidateServiceWrapper) FilterPath(filterPath ...string) *ESValidateServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESValidateServiceWrapper) Header(name string, value string) *ESValidateServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESValidateServiceWrapper) Headers(headers http.Header) *ESValidateServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESValidateServiceWrapper) Human(human bool) *ESValidateServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESValidateServiceWrapper) IgnoreUnavailable(ignoreUnavailable bool) *ESValidateServiceWrapper {
	w.obj = w.obj.IgnoreUnavailable(ignoreUnavailable)
	return w
}

func (w *ESValidateServiceWrapper) Index(index ...string) *ESValidateServiceWrapper {
	w.obj = w.obj.Index(index...)
	return w
}

func (w *ESValidateServiceWrapper) Lenient(lenient bool) *ESValidateServiceWrapper {
	w.obj = w.obj.Lenient(lenient)
	return w
}

func (w *ESValidateServiceWrapper) Pretty(pretty bool) *ESValidateServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESValidateServiceWrapper) Q(q string) *ESValidateServiceWrapper {
	w.obj = w.obj.Q(q)
	return w
}

func (w *ESValidateServiceWrapper) Query(query elastic.Query) *ESValidateServiceWrapper {
	w.obj = w.obj.Query(query)
	return w
}

func (w *ESValidateServiceWrapper) Rewrite(rewrite *bool) *ESValidateServiceWrapper {
	w.obj = w.obj.Rewrite(rewrite)
	return w
}

func (w *ESValidateServiceWrapper) Type(typ ...string) *ESValidateServiceWrapper {
	w.obj = w.obj.Type(typ...)
	return w
}

func (w *ESValidateServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackIlmDeleteLifecycleServiceWrapper) Do(ctx context.Context) (*elastic.XPackIlmDeleteLifecycleResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackIlmDeleteLifecycleResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackIlmDeleteLifecycleService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackIlmDeleteLifecycleService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackIlmDeleteLifecycleService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackIlmDeleteLifecycleService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackIlmDeleteLifecycleService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackIlmDeleteLifecycleService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackIlmDeleteLifecycleService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackIlmDeleteLifecycleServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackIlmDeleteLifecycleServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackIlmDeleteLifecycleServiceWrapper) FilterPath(filterPath ...string) *ESXPackIlmDeleteLifecycleServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackIlmDeleteLifecycleServiceWrapper) FlatSettings(flatSettings bool) *ESXPackIlmDeleteLifecycleServiceWrapper {
	w.obj = w.obj.FlatSettings(flatSettings)
	return w
}

func (w *ESXPackIlmDeleteLifecycleServiceWrapper) Header(name string, value string) *ESXPackIlmDeleteLifecycleServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackIlmDeleteLifecycleServiceWrapper) Headers(headers http.Header) *ESXPackIlmDeleteLifecycleServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackIlmDeleteLifecycleServiceWrapper) Human(human bool) *ESXPackIlmDeleteLifecycleServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackIlmDeleteLifecycleServiceWrapper) MasterTimeout(masterTimeout string) *ESXPackIlmDeleteLifecycleServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESXPackIlmDeleteLifecycleServiceWrapper) Policy(policy string) *ESXPackIlmDeleteLifecycleServiceWrapper {
	w.obj = w.obj.Policy(policy)
	return w
}

func (w *ESXPackIlmDeleteLifecycleServiceWrapper) Pretty(pretty bool) *ESXPackIlmDeleteLifecycleServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackIlmDeleteLifecycleServiceWrapper) Timeout(timeout string) *ESXPackIlmDeleteLifecycleServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESXPackIlmDeleteLifecycleServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackIlmGetLifecycleServiceWrapper) Do(ctx context.Context) (map[string]*elastic.XPackIlmGetLifecycleResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 map[string]*elastic.XPackIlmGetLifecycleResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackIlmGetLifecycleService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackIlmGetLifecycleService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackIlmGetLifecycleService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackIlmGetLifecycleService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackIlmGetLifecycleService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackIlmGetLifecycleService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackIlmGetLifecycleService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackIlmGetLifecycleServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackIlmGetLifecycleServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackIlmGetLifecycleServiceWrapper) FilterPath(filterPath ...string) *ESXPackIlmGetLifecycleServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackIlmGetLifecycleServiceWrapper) FlatSettings(flatSettings bool) *ESXPackIlmGetLifecycleServiceWrapper {
	w.obj = w.obj.FlatSettings(flatSettings)
	return w
}

func (w *ESXPackIlmGetLifecycleServiceWrapper) Header(name string, value string) *ESXPackIlmGetLifecycleServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackIlmGetLifecycleServiceWrapper) Headers(headers http.Header) *ESXPackIlmGetLifecycleServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackIlmGetLifecycleServiceWrapper) Human(human bool) *ESXPackIlmGetLifecycleServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackIlmGetLifecycleServiceWrapper) MasterTimeout(masterTimeout string) *ESXPackIlmGetLifecycleServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESXPackIlmGetLifecycleServiceWrapper) Policy(policies ...string) *ESXPackIlmGetLifecycleServiceWrapper {
	w.obj = w.obj.Policy(policies...)
	return w
}

func (w *ESXPackIlmGetLifecycleServiceWrapper) Pretty(pretty bool) *ESXPackIlmGetLifecycleServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackIlmGetLifecycleServiceWrapper) Timeout(timeout string) *ESXPackIlmGetLifecycleServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESXPackIlmGetLifecycleServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) BodyJson(body interface{}) *ESXPackIlmPutLifecycleServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) BodyString(body string) *ESXPackIlmPutLifecycleServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) Do(ctx context.Context) (*elastic.XPackIlmPutLifecycleResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackIlmPutLifecycleResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackIlmPutLifecycleService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackIlmPutLifecycleService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackIlmPutLifecycleService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackIlmPutLifecycleService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackIlmPutLifecycleService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackIlmPutLifecycleService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackIlmPutLifecycleService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackIlmPutLifecycleServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) FilterPath(filterPath ...string) *ESXPackIlmPutLifecycleServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) FlatSettings(flatSettings bool) *ESXPackIlmPutLifecycleServiceWrapper {
	w.obj = w.obj.FlatSettings(flatSettings)
	return w
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) Header(name string, value string) *ESXPackIlmPutLifecycleServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) Headers(headers http.Header) *ESXPackIlmPutLifecycleServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) Human(human bool) *ESXPackIlmPutLifecycleServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) MasterTimeout(masterTimeout string) *ESXPackIlmPutLifecycleServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) Policy(policy string) *ESXPackIlmPutLifecycleServiceWrapper {
	w.obj = w.obj.Policy(policy)
	return w
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) Pretty(pretty bool) *ESXPackIlmPutLifecycleServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) Timeout(timeout string) *ESXPackIlmPutLifecycleServiceWrapper {
	w.obj = w.obj.Timeout(timeout)
	return w
}

func (w *ESXPackIlmPutLifecycleServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackInfoServiceWrapper) Do(ctx context.Context) (*elastic.XPackInfoServiceResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackInfoServiceResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackInfoService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackInfoService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackInfoService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackInfoService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackInfoService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackInfoService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackInfoService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackInfoServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackInfoServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackInfoServiceWrapper) FilterPath(filterPath ...string) *ESXPackInfoServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackInfoServiceWrapper) Header(name string, value string) *ESXPackInfoServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackInfoServiceWrapper) Headers(headers http.Header) *ESXPackInfoServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackInfoServiceWrapper) Human(human bool) *ESXPackInfoServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackInfoServiceWrapper) Pretty(pretty bool) *ESXPackInfoServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackInfoServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackSecurityChangePasswordServiceWrapper) Body(body interface{}) *ESXPackSecurityChangePasswordServiceWrapper {
	w.obj = w.obj.Body(body)
	return w
}

func (w *ESXPackSecurityChangePasswordServiceWrapper) Do(ctx context.Context) (*elastic.XPackSecurityChangeUserPasswordResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackSecurityChangeUserPasswordResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackSecurityChangePasswordService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackSecurityChangePasswordService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackSecurityChangePasswordService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackSecurityChangePasswordService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackSecurityChangePasswordService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackSecurityChangePasswordService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackSecurityChangePasswordService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackSecurityChangePasswordServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackSecurityChangePasswordServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackSecurityChangePasswordServiceWrapper) FilterPath(filterPath ...string) *ESXPackSecurityChangePasswordServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackSecurityChangePasswordServiceWrapper) Header(name string, value string) *ESXPackSecurityChangePasswordServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackSecurityChangePasswordServiceWrapper) Headers(headers http.Header) *ESXPackSecurityChangePasswordServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackSecurityChangePasswordServiceWrapper) Human(human bool) *ESXPackSecurityChangePasswordServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackSecurityChangePasswordServiceWrapper) Password(password string) *ESXPackSecurityChangePasswordServiceWrapper {
	w.obj = w.obj.Password(password)
	return w
}

func (w *ESXPackSecurityChangePasswordServiceWrapper) Pretty(pretty bool) *ESXPackSecurityChangePasswordServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackSecurityChangePasswordServiceWrapper) Refresh(refresh string) *ESXPackSecurityChangePasswordServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESXPackSecurityChangePasswordServiceWrapper) Username(username string) *ESXPackSecurityChangePasswordServiceWrapper {
	w.obj = w.obj.Username(username)
	return w
}

func (w *ESXPackSecurityChangePasswordServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackSecurityDeleteRoleMappingServiceWrapper) Do(ctx context.Context) (*elastic.XPackSecurityDeleteRoleMappingResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackSecurityDeleteRoleMappingResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackSecurityDeleteRoleMappingService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackSecurityDeleteRoleMappingService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackSecurityDeleteRoleMappingService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackSecurityDeleteRoleMappingService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackSecurityDeleteRoleMappingService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackSecurityDeleteRoleMappingService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackSecurityDeleteRoleMappingService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackSecurityDeleteRoleMappingServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackSecurityDeleteRoleMappingServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackSecurityDeleteRoleMappingServiceWrapper) FilterPath(filterPath ...string) *ESXPackSecurityDeleteRoleMappingServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackSecurityDeleteRoleMappingServiceWrapper) Header(name string, value string) *ESXPackSecurityDeleteRoleMappingServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackSecurityDeleteRoleMappingServiceWrapper) Headers(headers http.Header) *ESXPackSecurityDeleteRoleMappingServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackSecurityDeleteRoleMappingServiceWrapper) Human(human bool) *ESXPackSecurityDeleteRoleMappingServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackSecurityDeleteRoleMappingServiceWrapper) Name(name string) *ESXPackSecurityDeleteRoleMappingServiceWrapper {
	w.obj = w.obj.Name(name)
	return w
}

func (w *ESXPackSecurityDeleteRoleMappingServiceWrapper) Pretty(pretty bool) *ESXPackSecurityDeleteRoleMappingServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackSecurityDeleteRoleMappingServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackSecurityDeleteRoleServiceWrapper) Do(ctx context.Context) (*elastic.XPackSecurityDeleteRoleResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackSecurityDeleteRoleResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackSecurityDeleteRoleService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackSecurityDeleteRoleService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackSecurityDeleteRoleService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackSecurityDeleteRoleService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackSecurityDeleteRoleService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackSecurityDeleteRoleService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackSecurityDeleteRoleService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackSecurityDeleteRoleServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackSecurityDeleteRoleServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackSecurityDeleteRoleServiceWrapper) FilterPath(filterPath ...string) *ESXPackSecurityDeleteRoleServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackSecurityDeleteRoleServiceWrapper) Header(name string, value string) *ESXPackSecurityDeleteRoleServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackSecurityDeleteRoleServiceWrapper) Headers(headers http.Header) *ESXPackSecurityDeleteRoleServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackSecurityDeleteRoleServiceWrapper) Human(human bool) *ESXPackSecurityDeleteRoleServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackSecurityDeleteRoleServiceWrapper) Name(name string) *ESXPackSecurityDeleteRoleServiceWrapper {
	w.obj = w.obj.Name(name)
	return w
}

func (w *ESXPackSecurityDeleteRoleServiceWrapper) Pretty(pretty bool) *ESXPackSecurityDeleteRoleServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackSecurityDeleteRoleServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackSecurityDeleteUserServiceWrapper) Do(ctx context.Context) (*elastic.XPackSecurityDeleteUserResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackSecurityDeleteUserResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackSecurityDeleteUserService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackSecurityDeleteUserService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackSecurityDeleteUserService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackSecurityDeleteUserService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackSecurityDeleteUserService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackSecurityDeleteUserService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackSecurityDeleteUserService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackSecurityDeleteUserServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackSecurityDeleteUserServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackSecurityDeleteUserServiceWrapper) FilterPath(filterPath ...string) *ESXPackSecurityDeleteUserServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackSecurityDeleteUserServiceWrapper) Header(name string, value string) *ESXPackSecurityDeleteUserServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackSecurityDeleteUserServiceWrapper) Headers(headers http.Header) *ESXPackSecurityDeleteUserServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackSecurityDeleteUserServiceWrapper) Human(human bool) *ESXPackSecurityDeleteUserServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackSecurityDeleteUserServiceWrapper) Pretty(pretty bool) *ESXPackSecurityDeleteUserServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackSecurityDeleteUserServiceWrapper) Refresh(refresh string) *ESXPackSecurityDeleteUserServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESXPackSecurityDeleteUserServiceWrapper) Username(username string) *ESXPackSecurityDeleteUserServiceWrapper {
	w.obj = w.obj.Username(username)
	return w
}

func (w *ESXPackSecurityDeleteUserServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackSecurityDisableUserServiceWrapper) Do(ctx context.Context) (*elastic.XPackSecurityDisableUserResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackSecurityDisableUserResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackSecurityDisableUserService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackSecurityDisableUserService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackSecurityDisableUserService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackSecurityDisableUserService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackSecurityDisableUserService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackSecurityDisableUserService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackSecurityDisableUserService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackSecurityDisableUserServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackSecurityDisableUserServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackSecurityDisableUserServiceWrapper) FilterPath(filterPath ...string) *ESXPackSecurityDisableUserServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackSecurityDisableUserServiceWrapper) Header(name string, value string) *ESXPackSecurityDisableUserServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackSecurityDisableUserServiceWrapper) Headers(headers http.Header) *ESXPackSecurityDisableUserServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackSecurityDisableUserServiceWrapper) Human(human bool) *ESXPackSecurityDisableUserServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackSecurityDisableUserServiceWrapper) Pretty(pretty bool) *ESXPackSecurityDisableUserServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackSecurityDisableUserServiceWrapper) Refresh(refresh string) *ESXPackSecurityDisableUserServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESXPackSecurityDisableUserServiceWrapper) Username(username string) *ESXPackSecurityDisableUserServiceWrapper {
	w.obj = w.obj.Username(username)
	return w
}

func (w *ESXPackSecurityDisableUserServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackSecurityEnableUserServiceWrapper) Do(ctx context.Context) (*elastic.XPackSecurityEnableUserResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackSecurityEnableUserResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackSecurityEnableUserService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackSecurityEnableUserService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackSecurityEnableUserService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackSecurityEnableUserService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackSecurityEnableUserService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackSecurityEnableUserService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackSecurityEnableUserService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackSecurityEnableUserServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackSecurityEnableUserServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackSecurityEnableUserServiceWrapper) FilterPath(filterPath ...string) *ESXPackSecurityEnableUserServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackSecurityEnableUserServiceWrapper) Header(name string, value string) *ESXPackSecurityEnableUserServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackSecurityEnableUserServiceWrapper) Headers(headers http.Header) *ESXPackSecurityEnableUserServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackSecurityEnableUserServiceWrapper) Human(human bool) *ESXPackSecurityEnableUserServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackSecurityEnableUserServiceWrapper) Pretty(pretty bool) *ESXPackSecurityEnableUserServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackSecurityEnableUserServiceWrapper) Refresh(refresh string) *ESXPackSecurityEnableUserServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESXPackSecurityEnableUserServiceWrapper) Username(username string) *ESXPackSecurityEnableUserServiceWrapper {
	w.obj = w.obj.Username(username)
	return w
}

func (w *ESXPackSecurityEnableUserServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackSecurityGetRoleMappingServiceWrapper) Do(ctx context.Context) (*elastic.XPackSecurityGetRoleMappingResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackSecurityGetRoleMappingResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackSecurityGetRoleMappingService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackSecurityGetRoleMappingService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackSecurityGetRoleMappingService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackSecurityGetRoleMappingService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackSecurityGetRoleMappingService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackSecurityGetRoleMappingService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackSecurityGetRoleMappingService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackSecurityGetRoleMappingServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackSecurityGetRoleMappingServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackSecurityGetRoleMappingServiceWrapper) FilterPath(filterPath ...string) *ESXPackSecurityGetRoleMappingServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackSecurityGetRoleMappingServiceWrapper) Header(name string, value string) *ESXPackSecurityGetRoleMappingServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackSecurityGetRoleMappingServiceWrapper) Headers(headers http.Header) *ESXPackSecurityGetRoleMappingServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackSecurityGetRoleMappingServiceWrapper) Human(human bool) *ESXPackSecurityGetRoleMappingServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackSecurityGetRoleMappingServiceWrapper) Name(name string) *ESXPackSecurityGetRoleMappingServiceWrapper {
	w.obj = w.obj.Name(name)
	return w
}

func (w *ESXPackSecurityGetRoleMappingServiceWrapper) Pretty(pretty bool) *ESXPackSecurityGetRoleMappingServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackSecurityGetRoleMappingServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackSecurityGetRoleServiceWrapper) Do(ctx context.Context) (*elastic.XPackSecurityGetRoleResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackSecurityGetRoleResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackSecurityGetRoleService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackSecurityGetRoleService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackSecurityGetRoleService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackSecurityGetRoleService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackSecurityGetRoleService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackSecurityGetRoleService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackSecurityGetRoleService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackSecurityGetRoleServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackSecurityGetRoleServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackSecurityGetRoleServiceWrapper) FilterPath(filterPath ...string) *ESXPackSecurityGetRoleServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackSecurityGetRoleServiceWrapper) Header(name string, value string) *ESXPackSecurityGetRoleServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackSecurityGetRoleServiceWrapper) Headers(headers http.Header) *ESXPackSecurityGetRoleServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackSecurityGetRoleServiceWrapper) Human(human bool) *ESXPackSecurityGetRoleServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackSecurityGetRoleServiceWrapper) Name(name string) *ESXPackSecurityGetRoleServiceWrapper {
	w.obj = w.obj.Name(name)
	return w
}

func (w *ESXPackSecurityGetRoleServiceWrapper) Pretty(pretty bool) *ESXPackSecurityGetRoleServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackSecurityGetRoleServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackSecurityGetUserServiceWrapper) Do(ctx context.Context) (*elastic.XPackSecurityGetUserResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackSecurityGetUserResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackSecurityGetUserService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackSecurityGetUserService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackSecurityGetUserService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackSecurityGetUserService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackSecurityGetUserService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackSecurityGetUserService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackSecurityGetUserService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackSecurityGetUserServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackSecurityGetUserServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackSecurityGetUserServiceWrapper) FilterPath(filterPath ...string) *ESXPackSecurityGetUserServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackSecurityGetUserServiceWrapper) Header(name string, value string) *ESXPackSecurityGetUserServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackSecurityGetUserServiceWrapper) Headers(headers http.Header) *ESXPackSecurityGetUserServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackSecurityGetUserServiceWrapper) Human(human bool) *ESXPackSecurityGetUserServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackSecurityGetUserServiceWrapper) Pretty(pretty bool) *ESXPackSecurityGetUserServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackSecurityGetUserServiceWrapper) Usernames(usernames ...string) *ESXPackSecurityGetUserServiceWrapper {
	w.obj = w.obj.Usernames(usernames...)
	return w
}

func (w *ESXPackSecurityGetUserServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackSecurityPutRoleMappingServiceWrapper) Body(body interface{}) *ESXPackSecurityPutRoleMappingServiceWrapper {
	w.obj = w.obj.Body(body)
	return w
}

func (w *ESXPackSecurityPutRoleMappingServiceWrapper) Do(ctx context.Context) (*elastic.XPackSecurityPutRoleMappingResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackSecurityPutRoleMappingResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackSecurityPutRoleMappingService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackSecurityPutRoleMappingService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackSecurityPutRoleMappingService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackSecurityPutRoleMappingService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackSecurityPutRoleMappingService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackSecurityPutRoleMappingService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackSecurityPutRoleMappingService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackSecurityPutRoleMappingServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackSecurityPutRoleMappingServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackSecurityPutRoleMappingServiceWrapper) FilterPath(filterPath ...string) *ESXPackSecurityPutRoleMappingServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackSecurityPutRoleMappingServiceWrapper) Header(name string, value string) *ESXPackSecurityPutRoleMappingServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackSecurityPutRoleMappingServiceWrapper) Headers(headers http.Header) *ESXPackSecurityPutRoleMappingServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackSecurityPutRoleMappingServiceWrapper) Human(human bool) *ESXPackSecurityPutRoleMappingServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackSecurityPutRoleMappingServiceWrapper) Name(name string) *ESXPackSecurityPutRoleMappingServiceWrapper {
	w.obj = w.obj.Name(name)
	return w
}

func (w *ESXPackSecurityPutRoleMappingServiceWrapper) Pretty(pretty bool) *ESXPackSecurityPutRoleMappingServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackSecurityPutRoleMappingServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackSecurityPutRoleServiceWrapper) Body(body interface{}) *ESXPackSecurityPutRoleServiceWrapper {
	w.obj = w.obj.Body(body)
	return w
}

func (w *ESXPackSecurityPutRoleServiceWrapper) Do(ctx context.Context) (*elastic.XPackSecurityPutRoleResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackSecurityPutRoleResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackSecurityPutRoleService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackSecurityPutRoleService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackSecurityPutRoleService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackSecurityPutRoleService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackSecurityPutRoleService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackSecurityPutRoleService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackSecurityPutRoleService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackSecurityPutRoleServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackSecurityPutRoleServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackSecurityPutRoleServiceWrapper) FilterPath(filterPath ...string) *ESXPackSecurityPutRoleServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackSecurityPutRoleServiceWrapper) Header(name string, value string) *ESXPackSecurityPutRoleServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackSecurityPutRoleServiceWrapper) Headers(headers http.Header) *ESXPackSecurityPutRoleServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackSecurityPutRoleServiceWrapper) Human(human bool) *ESXPackSecurityPutRoleServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackSecurityPutRoleServiceWrapper) Name(name string) *ESXPackSecurityPutRoleServiceWrapper {
	w.obj = w.obj.Name(name)
	return w
}

func (w *ESXPackSecurityPutRoleServiceWrapper) Pretty(pretty bool) *ESXPackSecurityPutRoleServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackSecurityPutRoleServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackSecurityPutUserServiceWrapper) Body(body interface{}) *ESXPackSecurityPutUserServiceWrapper {
	w.obj = w.obj.Body(body)
	return w
}

func (w *ESXPackSecurityPutUserServiceWrapper) Do(ctx context.Context) (*elastic.XPackSecurityPutUserResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackSecurityPutUserResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackSecurityPutUserService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackSecurityPutUserService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackSecurityPutUserService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackSecurityPutUserService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackSecurityPutUserService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackSecurityPutUserService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackSecurityPutUserService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackSecurityPutUserServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackSecurityPutUserServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackSecurityPutUserServiceWrapper) FilterPath(filterPath ...string) *ESXPackSecurityPutUserServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackSecurityPutUserServiceWrapper) Header(name string, value string) *ESXPackSecurityPutUserServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackSecurityPutUserServiceWrapper) Headers(headers http.Header) *ESXPackSecurityPutUserServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackSecurityPutUserServiceWrapper) Human(human bool) *ESXPackSecurityPutUserServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackSecurityPutUserServiceWrapper) Pretty(pretty bool) *ESXPackSecurityPutUserServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackSecurityPutUserServiceWrapper) Refresh(refresh string) *ESXPackSecurityPutUserServiceWrapper {
	w.obj = w.obj.Refresh(refresh)
	return w
}

func (w *ESXPackSecurityPutUserServiceWrapper) User(user *elastic.XPackSecurityPutUserRequest) *ESXPackSecurityPutUserServiceWrapper {
	w.obj = w.obj.User(user)
	return w
}

func (w *ESXPackSecurityPutUserServiceWrapper) Username(username string) *ESXPackSecurityPutUserServiceWrapper {
	w.obj = w.obj.Username(username)
	return w
}

func (w *ESXPackSecurityPutUserServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackWatcherAckWatchServiceWrapper) ActionId(actionId ...string) *ESXPackWatcherAckWatchServiceWrapper {
	w.obj = w.obj.ActionId(actionId...)
	return w
}

func (w *ESXPackWatcherAckWatchServiceWrapper) Do(ctx context.Context) (*elastic.XPackWatcherAckWatchResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackWatcherAckWatchResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackWatcherAckWatchService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackWatcherAckWatchService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackWatcherAckWatchService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackWatcherAckWatchService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackWatcherAckWatchService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackWatcherAckWatchService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackWatcherAckWatchService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackWatcherAckWatchServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackWatcherAckWatchServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackWatcherAckWatchServiceWrapper) FilterPath(filterPath ...string) *ESXPackWatcherAckWatchServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackWatcherAckWatchServiceWrapper) Header(name string, value string) *ESXPackWatcherAckWatchServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackWatcherAckWatchServiceWrapper) Headers(headers http.Header) *ESXPackWatcherAckWatchServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackWatcherAckWatchServiceWrapper) Human(human bool) *ESXPackWatcherAckWatchServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackWatcherAckWatchServiceWrapper) MasterTimeout(masterTimeout string) *ESXPackWatcherAckWatchServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESXPackWatcherAckWatchServiceWrapper) Pretty(pretty bool) *ESXPackWatcherAckWatchServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackWatcherAckWatchServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackWatcherAckWatchServiceWrapper) WatchId(watchId string) *ESXPackWatcherAckWatchServiceWrapper {
	w.obj = w.obj.WatchId(watchId)
	return w
}

func (w *ESXPackWatcherActivateWatchServiceWrapper) Do(ctx context.Context) (*elastic.XPackWatcherActivateWatchResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackWatcherActivateWatchResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackWatcherActivateWatchService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackWatcherActivateWatchService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackWatcherActivateWatchService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackWatcherActivateWatchService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackWatcherActivateWatchService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackWatcherActivateWatchService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackWatcherActivateWatchService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackWatcherActivateWatchServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackWatcherActivateWatchServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackWatcherActivateWatchServiceWrapper) FilterPath(filterPath ...string) *ESXPackWatcherActivateWatchServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackWatcherActivateWatchServiceWrapper) Header(name string, value string) *ESXPackWatcherActivateWatchServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackWatcherActivateWatchServiceWrapper) Headers(headers http.Header) *ESXPackWatcherActivateWatchServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackWatcherActivateWatchServiceWrapper) Human(human bool) *ESXPackWatcherActivateWatchServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackWatcherActivateWatchServiceWrapper) MasterTimeout(masterTimeout string) *ESXPackWatcherActivateWatchServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESXPackWatcherActivateWatchServiceWrapper) Pretty(pretty bool) *ESXPackWatcherActivateWatchServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackWatcherActivateWatchServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackWatcherActivateWatchServiceWrapper) WatchId(watchId string) *ESXPackWatcherActivateWatchServiceWrapper {
	w.obj = w.obj.WatchId(watchId)
	return w
}

func (w *ESXPackWatcherDeactivateWatchServiceWrapper) Do(ctx context.Context) (*elastic.XPackWatcherDeactivateWatchResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackWatcherDeactivateWatchResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackWatcherDeactivateWatchService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackWatcherDeactivateWatchService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackWatcherDeactivateWatchService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackWatcherDeactivateWatchService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackWatcherDeactivateWatchService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackWatcherDeactivateWatchService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackWatcherDeactivateWatchService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackWatcherDeactivateWatchServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackWatcherDeactivateWatchServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackWatcherDeactivateWatchServiceWrapper) FilterPath(filterPath ...string) *ESXPackWatcherDeactivateWatchServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackWatcherDeactivateWatchServiceWrapper) Header(name string, value string) *ESXPackWatcherDeactivateWatchServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackWatcherDeactivateWatchServiceWrapper) Headers(headers http.Header) *ESXPackWatcherDeactivateWatchServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackWatcherDeactivateWatchServiceWrapper) Human(human bool) *ESXPackWatcherDeactivateWatchServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackWatcherDeactivateWatchServiceWrapper) MasterTimeout(masterTimeout string) *ESXPackWatcherDeactivateWatchServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESXPackWatcherDeactivateWatchServiceWrapper) Pretty(pretty bool) *ESXPackWatcherDeactivateWatchServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackWatcherDeactivateWatchServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackWatcherDeactivateWatchServiceWrapper) WatchId(watchId string) *ESXPackWatcherDeactivateWatchServiceWrapper {
	w.obj = w.obj.WatchId(watchId)
	return w
}

func (w *ESXPackWatcherDeleteWatchServiceWrapper) Do(ctx context.Context) (*elastic.XPackWatcherDeleteWatchResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackWatcherDeleteWatchResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackWatcherDeleteWatchService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackWatcherDeleteWatchService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackWatcherDeleteWatchService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackWatcherDeleteWatchService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackWatcherDeleteWatchService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackWatcherDeleteWatchService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackWatcherDeleteWatchService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackWatcherDeleteWatchServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackWatcherDeleteWatchServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackWatcherDeleteWatchServiceWrapper) FilterPath(filterPath ...string) *ESXPackWatcherDeleteWatchServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackWatcherDeleteWatchServiceWrapper) Header(name string, value string) *ESXPackWatcherDeleteWatchServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackWatcherDeleteWatchServiceWrapper) Headers(headers http.Header) *ESXPackWatcherDeleteWatchServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackWatcherDeleteWatchServiceWrapper) Human(human bool) *ESXPackWatcherDeleteWatchServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackWatcherDeleteWatchServiceWrapper) Id(id string) *ESXPackWatcherDeleteWatchServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESXPackWatcherDeleteWatchServiceWrapper) MasterTimeout(masterTimeout string) *ESXPackWatcherDeleteWatchServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESXPackWatcherDeleteWatchServiceWrapper) Pretty(pretty bool) *ESXPackWatcherDeleteWatchServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackWatcherDeleteWatchServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackWatcherExecuteWatchServiceWrapper) BodyJson(body interface{}) *ESXPackWatcherExecuteWatchServiceWrapper {
	w.obj = w.obj.BodyJson(body)
	return w
}

func (w *ESXPackWatcherExecuteWatchServiceWrapper) BodyString(body string) *ESXPackWatcherExecuteWatchServiceWrapper {
	w.obj = w.obj.BodyString(body)
	return w
}

func (w *ESXPackWatcherExecuteWatchServiceWrapper) Debug(debug bool) *ESXPackWatcherExecuteWatchServiceWrapper {
	w.obj = w.obj.Debug(debug)
	return w
}

func (w *ESXPackWatcherExecuteWatchServiceWrapper) Do(ctx context.Context) (*elastic.XPackWatcherExecuteWatchResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackWatcherExecuteWatchResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackWatcherExecuteWatchService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackWatcherExecuteWatchService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackWatcherExecuteWatchService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackWatcherExecuteWatchService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackWatcherExecuteWatchService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackWatcherExecuteWatchService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackWatcherExecuteWatchService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackWatcherExecuteWatchServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackWatcherExecuteWatchServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackWatcherExecuteWatchServiceWrapper) FilterPath(filterPath ...string) *ESXPackWatcherExecuteWatchServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackWatcherExecuteWatchServiceWrapper) Header(name string, value string) *ESXPackWatcherExecuteWatchServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackWatcherExecuteWatchServiceWrapper) Headers(headers http.Header) *ESXPackWatcherExecuteWatchServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackWatcherExecuteWatchServiceWrapper) Human(human bool) *ESXPackWatcherExecuteWatchServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackWatcherExecuteWatchServiceWrapper) Id(id string) *ESXPackWatcherExecuteWatchServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESXPackWatcherExecuteWatchServiceWrapper) Pretty(pretty bool) *ESXPackWatcherExecuteWatchServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackWatcherExecuteWatchServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackWatcherGetWatchServiceWrapper) Do(ctx context.Context) (*elastic.XPackWatcherGetWatchResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackWatcherGetWatchResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackWatcherGetWatchService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackWatcherGetWatchService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackWatcherGetWatchService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackWatcherGetWatchService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackWatcherGetWatchService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackWatcherGetWatchService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackWatcherGetWatchService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackWatcherGetWatchServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackWatcherGetWatchServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackWatcherGetWatchServiceWrapper) FilterPath(filterPath ...string) *ESXPackWatcherGetWatchServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackWatcherGetWatchServiceWrapper) Header(name string, value string) *ESXPackWatcherGetWatchServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackWatcherGetWatchServiceWrapper) Headers(headers http.Header) *ESXPackWatcherGetWatchServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackWatcherGetWatchServiceWrapper) Human(human bool) *ESXPackWatcherGetWatchServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackWatcherGetWatchServiceWrapper) Id(id string) *ESXPackWatcherGetWatchServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESXPackWatcherGetWatchServiceWrapper) Pretty(pretty bool) *ESXPackWatcherGetWatchServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackWatcherGetWatchServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackWatcherPutWatchServiceWrapper) Active(active bool) *ESXPackWatcherPutWatchServiceWrapper {
	w.obj = w.obj.Active(active)
	return w
}

func (w *ESXPackWatcherPutWatchServiceWrapper) Body(body interface{}) *ESXPackWatcherPutWatchServiceWrapper {
	w.obj = w.obj.Body(body)
	return w
}

func (w *ESXPackWatcherPutWatchServiceWrapper) Do(ctx context.Context) (*elastic.XPackWatcherPutWatchResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackWatcherPutWatchResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackWatcherPutWatchService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackWatcherPutWatchService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackWatcherPutWatchService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackWatcherPutWatchService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackWatcherPutWatchService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackWatcherPutWatchService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackWatcherPutWatchService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackWatcherPutWatchServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackWatcherPutWatchServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackWatcherPutWatchServiceWrapper) FilterPath(filterPath ...string) *ESXPackWatcherPutWatchServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackWatcherPutWatchServiceWrapper) Header(name string, value string) *ESXPackWatcherPutWatchServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackWatcherPutWatchServiceWrapper) Headers(headers http.Header) *ESXPackWatcherPutWatchServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackWatcherPutWatchServiceWrapper) Human(human bool) *ESXPackWatcherPutWatchServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackWatcherPutWatchServiceWrapper) Id(id string) *ESXPackWatcherPutWatchServiceWrapper {
	w.obj = w.obj.Id(id)
	return w
}

func (w *ESXPackWatcherPutWatchServiceWrapper) IfPrimaryTerm(primaryTerm int64) *ESXPackWatcherPutWatchServiceWrapper {
	w.obj = w.obj.IfPrimaryTerm(primaryTerm)
	return w
}

func (w *ESXPackWatcherPutWatchServiceWrapper) IfSeqNo(seqNo int64) *ESXPackWatcherPutWatchServiceWrapper {
	w.obj = w.obj.IfSeqNo(seqNo)
	return w
}

func (w *ESXPackWatcherPutWatchServiceWrapper) MasterTimeout(masterTimeout string) *ESXPackWatcherPutWatchServiceWrapper {
	w.obj = w.obj.MasterTimeout(masterTimeout)
	return w
}

func (w *ESXPackWatcherPutWatchServiceWrapper) Pretty(pretty bool) *ESXPackWatcherPutWatchServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackWatcherPutWatchServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackWatcherStartServiceWrapper) Do(ctx context.Context) (*elastic.XPackWatcherStartResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackWatcherStartResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackWatcherStartService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackWatcherStartService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackWatcherStartService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackWatcherStartService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackWatcherStartService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackWatcherStartService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackWatcherStartService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackWatcherStartServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackWatcherStartServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackWatcherStartServiceWrapper) FilterPath(filterPath ...string) *ESXPackWatcherStartServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackWatcherStartServiceWrapper) Header(name string, value string) *ESXPackWatcherStartServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackWatcherStartServiceWrapper) Headers(headers http.Header) *ESXPackWatcherStartServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackWatcherStartServiceWrapper) Human(human bool) *ESXPackWatcherStartServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackWatcherStartServiceWrapper) Pretty(pretty bool) *ESXPackWatcherStartServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackWatcherStartServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackWatcherStatsServiceWrapper) Do(ctx context.Context) (*elastic.XPackWatcherStatsResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackWatcherStatsResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackWatcherStatsService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackWatcherStatsService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackWatcherStatsService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackWatcherStatsService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackWatcherStatsService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackWatcherStatsService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackWatcherStatsService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackWatcherStatsServiceWrapper) EmitStacktraces(emitStacktraces bool) *ESXPackWatcherStatsServiceWrapper {
	w.obj = w.obj.EmitStacktraces(emitStacktraces)
	return w
}

func (w *ESXPackWatcherStatsServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackWatcherStatsServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackWatcherStatsServiceWrapper) FilterPath(filterPath ...string) *ESXPackWatcherStatsServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackWatcherStatsServiceWrapper) Header(name string, value string) *ESXPackWatcherStatsServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackWatcherStatsServiceWrapper) Headers(headers http.Header) *ESXPackWatcherStatsServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackWatcherStatsServiceWrapper) Human(human bool) *ESXPackWatcherStatsServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackWatcherStatsServiceWrapper) Metric(metric string) *ESXPackWatcherStatsServiceWrapper {
	w.obj = w.obj.Metric(metric)
	return w
}

func (w *ESXPackWatcherStatsServiceWrapper) Pretty(pretty bool) *ESXPackWatcherStatsServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackWatcherStatsServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}

func (w *ESXPackWatcherStopServiceWrapper) Do(ctx context.Context) (*elastic.XPackWatcherStopResponse, error) {
	ctxOptions := FromContext(ctx)
	var res0 *elastic.XPackWatcherStopResponse
	var err error
	err = w.retry.Do(func() error {
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.XPackWatcherStopService.Do", w.options.Name)); err != nil {
				return err
			}
		}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.XPackWatcherStopService.Do", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.XPackWatcherStopService.Do", w.options.Name), token)
			}
		}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "elastic.XPackWatcherStopService.Do")
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
			w.inflightMetric.WithLabelValues("elastic.XPackWatcherStopService.Do", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("elastic.XPackWatcherStopService.Do", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("elastic.XPackWatcherStopService.Do", ErrCode(err), ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
		res0, err = w.obj.Do(ctx)
		if err != nil && span != nil {
			span.SetTag("error", err.Error())
		}
		return err
	})
	return res0, err
}

func (w *ESXPackWatcherStopServiceWrapper) ErrorTrace(errorTrace bool) *ESXPackWatcherStopServiceWrapper {
	w.obj = w.obj.ErrorTrace(errorTrace)
	return w
}

func (w *ESXPackWatcherStopServiceWrapper) FilterPath(filterPath ...string) *ESXPackWatcherStopServiceWrapper {
	w.obj = w.obj.FilterPath(filterPath...)
	return w
}

func (w *ESXPackWatcherStopServiceWrapper) Header(name string, value string) *ESXPackWatcherStopServiceWrapper {
	w.obj = w.obj.Header(name, value)
	return w
}

func (w *ESXPackWatcherStopServiceWrapper) Headers(headers http.Header) *ESXPackWatcherStopServiceWrapper {
	w.obj = w.obj.Headers(headers)
	return w
}

func (w *ESXPackWatcherStopServiceWrapper) Human(human bool) *ESXPackWatcherStopServiceWrapper {
	w.obj = w.obj.Human(human)
	return w
}

func (w *ESXPackWatcherStopServiceWrapper) Pretty(pretty bool) *ESXPackWatcherStopServiceWrapper {
	w.obj = w.obj.Pretty(pretty)
	return w
}

func (w *ESXPackWatcherStopServiceWrapper) Validate() error {
	var err error
	err = w.obj.Validate()
	return err
}
