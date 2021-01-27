package wrap

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type MongoOptions struct {
	URI            string        `dft:"mongodb://localhost:27017"`
	ConnectTimeout time.Duration `dft:"3s"`
	PingTimeout    time.Duration `dft:"2s"`
}

type MongoClientWrapperOptions struct {
	Retry            RetryOptions
	Wrapper          WrapperOptions
	Mongo            MongoOptions
	RateLimiterGroup RateLimiterGroupOptions
}

func NewMongoClientWrapperWithOptions(options *MongoClientWrapperOptions) (*MongoClientWrapper, error) {
	retry, err := NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}
	rateLimiterGroup, err := NewRateLimiterGroup(&options.RateLimiterGroup)
	if err != nil {
		return nil, errors.Wrap(err, "NewRateLimiterGroup failed")
	}

	client, err := mongo.NewClient(mopt.Client().ApplyURI(options.Mongo.URI))
	if err != nil {
		return nil, errors.Wrap(err, "mongo.NewClient failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), options.Mongo.ConnectTimeout)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "mongo.Client.Connect failed")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.Wrap(err, "mongo.Client.Ping failed")
	}

	w := &MongoClientWrapper{
		obj:              client,
		retry:            retry,
		options:          &options.Wrapper,
		rateLimiterGroup: rateLimiterGroup,
	}

	if w.options.EnableMetric {
		w.CreateMetric(w.options)
	}

	return w, nil
}

func NewMongoClientWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*MongoClientWrapper, error) {
	var options MongoClientWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewMongoClientWrapperWithOptions(&options)
	if err != nil {
		return nil, errors.Wrap(err, "NewMongoClientWrapperWithOptions failed")
	}

	refxOptions := refx.NewOptions(opts...)
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Wrapper"), w.OnWrapperChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Retry"), w.OnRetryChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("RateLimiterGroup"), w.OnRateLimiterGroupChange(opts...))
	cfg.AddOnItemChangeHandler(refxOptions.FormatKey("Mongo"), func(cfg *config.Config) error {
		var options MongoOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}

		client, err := mongo.NewClient(mopt.Client().ApplyURI(options.URI))
		if err != nil {
			return errors.Wrap(err, "mongo.NewClient failed")
		}

		ctx, cancel := context.WithTimeout(context.Background(), options.ConnectTimeout)
		defer cancel()
		if err := client.Connect(ctx); err != nil {
			return errors.Wrap(err, "mongo.Client.Connect failed")
		}

		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			return errors.Wrap(err, "mongo.Client.Ping failed")
		}

		w.obj = client
		return nil
	})

	return w, err
}
