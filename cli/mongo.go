package cli

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoOptions struct {
	URI            string        `dft:"mongodb://localhost:27017"`
	ConnectTimeout time.Duration `dft:"3s"`
	PingTimeout    time.Duration `dft:"2s"`
}

func NewMongoWithOptions(options *MongoOptions) (*mongo.Client, error) {
	client, err := mongo.NewClient(mopt.Client().ApplyURI(options.URI))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), options.ConnectTimeout)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		return nil, errors.WithMessage(err, "mongo.Client.Connect failed")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.WithMessage(err, "mongo.Client.Ping failed")
	}

	return client, nil
}
