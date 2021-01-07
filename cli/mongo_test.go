package cli

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestNewMongoWithOptions(t *testing.T) {
	patches := ApplyFunc(mongo.NewClient, func(opts ...*mopt.ClientOptions) (*mongo.Client, error) {
		return &mongo.Client{}, nil
	}).ApplyMethod(reflect.TypeOf(&mongo.Client{}), "Connect", func(client *mongo.Client, ctx context.Context) error {
		return nil
	}).ApplyMethod(reflect.TypeOf(&mongo.Client{}), "Ping", func(client *mongo.Client, ctx context.Context, rp *readpref.ReadPref) error {
		return nil
	}).ApplyMethod(reflect.TypeOf(&mongo.Collection{}), "InsertOne", func(client *mongo.Collection, ctx context.Context, document interface{},
		opts ...*mopt.InsertOneOptions) (*mongo.InsertOneResult, error) {
		return &mongo.InsertOneResult{InsertedID: "123"}, nil
	}).ApplyMethod(reflect.TypeOf(&mongo.Collection{}), "FindOne", func(client *mongo.Collection, ctx context.Context, filter interface{},
		opts ...*mopt.FindOneOptions) *mongo.SingleResult {
		return &mongo.SingleResult{}
	}).ApplyMethod(reflect.TypeOf(&mongo.SingleResult{}), "Decode", func(result *mongo.SingleResult, v interface{}) error {
		return nil
	})
	defer patches.Reset()

	Convey("TestNewMongoWithOptions", t, func() {
		cli, err := NewMongoWithOptions(&MongoOptions{
			URI:            "mongodb://localhost:27016",
			ConnectTimeout: 3 * time.Second,
			PingTimeout:    3 * time.Second,
		})
		So(err, ShouldBeNil)
		So(cli, ShouldNotBeNil)

		collection := cli.Database("hatlonely").Collection("task")

		{
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
			So(err, ShouldBeNil)
			fmt.Println(res.InsertedID)
		}
		{
			var result struct {
				Name  string
				Value float64
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := collection.FindOne(ctx, bson.M{"name": "pi"}).Decode(&result)
			So(err, ShouldBeNil)
			fmt.Println(result)
		}
	})
}
