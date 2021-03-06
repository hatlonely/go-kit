package cli

import (
	"context"
	"time"

	"github.com/olivere/elastic/v7"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

func NewElasticSearch(opts ...ElasticSearchOption) (*elastic.Client, error) {
	options := defaultElasticSearchOptions
	for _, opt := range opts {
		opt(&options)
	}

	return NewElasticSearchWithOptions(&options)
}

func NewElasticSearchWithConfig(cfg *config.Config, opts ...refx.Option) (*elastic.Client, error) {
	options := defaultElasticSearchOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, err
	}
	return NewElasticSearchWithOptions(&options)
}

func NewElasticSearchWithOptions(options *ElasticSearchOptions) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(options.URI),
		elastic.SetSniff(options.EnableSniff),
		elastic.SetBasicAuth(options.Username, options.Password),
		elastic.SetHealthcheck(options.EnableHealthCheck),
		elastic.SetHealthcheckInterval(options.HealthCheckInterval),
		elastic.SetHealthcheckTimeout(options.HealthCheckTimeout),
		elastic.SetHealthcheckTimeoutStartup(options.HealthCheckTimeoutStartUp),
	)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()
	if _, _, err := client.Ping(options.URI).Do(ctx); err != nil {
		return nil, err
	}

	return client, err
}

type ElasticSearchOptions struct {
	URI                       string `dft:"http://elasticsearch:9200"`
	EnableSniff               bool
	Username                  string
	Password                  string
	EnableHealthCheck         bool          `dft:"true"`
	HealthCheckInterval       time.Duration `dft:"60s"`
	HealthCheckTimeout        time.Duration `dft:"5s"`
	HealthCheckTimeoutStartUp time.Duration `dft:"5s"`
}

var defaultElasticSearchOptions = ElasticSearchOptions{
	URI:                       "http://elasticsearch:9200",
	EnableSniff:               false,
	EnableHealthCheck:         true,
	HealthCheckInterval:       60 * time.Second,
	HealthCheckTimeout:        5 * time.Second,
	HealthCheckTimeoutStartUp: 5 * time.Second,
}

type ElasticSearchOption func(options *ElasticSearchOptions)

func WithElasticSearchURI(uri string) ElasticSearchOption {
	return func(options *ElasticSearchOptions) {
		options.URI = uri
	}
}

func WithElasticSearchEnableSniff() ElasticSearchOption {
	return func(options *ElasticSearchOptions) {
		options.EnableSniff = true
	}
}

func WithElasticSearchDisableHealthCheck() ElasticSearchOption {
	return func(options *ElasticSearchOptions) {
		options.EnableHealthCheck = false
	}
}

func WithElasticSearchHealthCheck(interval time.Duration, timeout time.Duration, timeoutStartUp time.Duration) ElasticSearchOption {
	return func(options *ElasticSearchOptions) {
		options.HealthCheckInterval = interval
		options.HealthCheckTimeout = timeout
		options.HealthCheckTimeoutStartUp = timeoutStartUp
	}
}
