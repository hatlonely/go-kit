package cli

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

func TestNewMongoWithOptions(t *testing.T) {
	Convey("TestNewMongoWithOptions", t, func() {
		cli, err := NewMongoWithOptions(&MongoOptions{
			URI:            "mongodb://localhost:27017",
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
