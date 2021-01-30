package wrap

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

func init() {
	RegisterErrCode(mongo.MongocryptError{}, func(err error) string {
		e := err.(mongo.MongocryptError)
		return fmt.Sprintf("mongo_%v", e.Code)
	})
}

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

func NewMongoClientWithOptions(options *MongoOptions) (*mongo.Client, error) {
	client, err := mongo.NewClient(mopt.Client().ApplyURI(options.URI))
	if err != nil {
		return nil, errors.Wrap(err, "mongo.NewClient failed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), options.ConnectTimeout)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "mongo.Client.Connect failed")
	}
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.Wrap(err, "mongo.Client.Ping failed")
	}
	return client, nil
}

func NewMongoClientWrapperWithOptions(options *MongoClientWrapperOptions, opts ...refx.Option) (*MongoClientWrapper, error) {
	var w MongoClientWrapper
	var err error

	w.options = &options.Wrapper
	w.retry, err = NewRetryWithOptions(&options.Retry)
	if err != nil {
		return nil, errors.Wrap(err, "NewRetryWithOptions failed")
	}
	w.rateLimiterGroup, err = NewRateLimiterGroupWithOptions(&options.RateLimiterGroup, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "NewRateLimiterGroupWithOptions failed")
	}
	if w.options.EnableMetric {
		w.CreateMetric(w.options)
	}
	client, err := NewMongoClientWithOptions(&options.Mongo)
	if err != nil {
		return nil, errors.Wrap(err, "NewMongoClientWithOptions failed")
	}
	w.obj = client

	return &w, nil
}

func NewMongoClientWrapperWithConfig(cfg *config.Config, opts ...refx.Option) (*MongoClientWrapper, error) {
	var options MongoClientWrapperOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "config.Config.Unmarshal failed")
	}
	w, err := NewMongoClientWrapperWithOptions(&options, opts...)
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
		client, err := NewMongoClientWithOptions(&options)
		if err != nil {
			return errors.Wrap(err, "NewMongoClientWithOptions failed")
		}
		w.obj = client
		return nil
	})

	return w, err
}
