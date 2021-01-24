package wrap

import (
	"context"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type ESOptions struct {
	URI                       string `dft:"http://elasticsearch:9200"`
	EnableSniff               bool
	Username                  string
	Password                  string
	EnableHealthCheck         bool          `dft:"true"`
	HealthCheckInterval       time.Duration `dft:"60s"`
	HealthCheckTimeout        time.Duration `dft:"5s"`
	HealthCheckTimeoutStartUp time.Duration `dft:"5s"`
}

type ESClientWrapperOptions struct {
	Retry   RetryOptions
	Wrapper WrapperOptions
	ES      ESOptions
}

func NewESClientWrapperWithOptions(options *ESClientWrapperOptions) (*ESClientWrapper, error) {
	retry, err := NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}

	client, err := elastic.NewClient(
		elastic.SetURL(options.ES.URI),
		elastic.SetSniff(options.ES.EnableSniff),
		elastic.SetBasicAuth(options.ES.Username, options.ES.Password),
		elastic.SetHealthcheck(options.ES.EnableHealthCheck),
		elastic.SetHealthcheckInterval(options.ES.HealthCheckInterval),
		elastic.SetHealthcheckTimeout(options.ES.HealthCheckTimeout),
		elastic.SetHealthcheckTimeoutStartup(options.ES.HealthCheckTimeoutStartUp),
	)
	if err != nil {
		return nil, errors.Wrap(err, "elastic.NewClient failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()
	if _, _, err := client.Ping(options.ES.URI).Do(ctx); err != nil {
		return nil, errors.Wrap(err, "elastic.Client.Ping failed")
	}

	return &ESClientWrapper{
		obj:     client,
		retry:   retry,
		options: &options.Wrapper,
	}, nil
}

func NewESClientWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*ESClientWrapper, error) {
	var options ESClientWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewESClientWrapperWithOptions(&options)
	if err != nil {
		return nil, errors.Wrap(err, "NewMongoClientWrapperWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Wrapper"), w.OnWrapperChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Retry"), w.OnRetryChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("ES"), func(cfg *config.Config) error {
		var options ESOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}

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
			return errors.Wrap(err, "elastic.NewClient failed")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
		defer cancel()
		if _, _, err := client.Ping(options.URI).Do(ctx); err != nil {
			return errors.Wrap(err, "elastic.Client.Ping failed")
		}

		w.obj = client
		return nil
	})

	return w, err
}
