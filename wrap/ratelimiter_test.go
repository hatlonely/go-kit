package wrap

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/time/rate"

	"github.com/hatlonely/go-kit/refx"
)

func Test1(t *testing.T) {
	Convey("Test", t, func() {
		limiter := rate.NewLimiter(rate.Every(200*time.Millisecond), 100)
		limiter.Allow()
	})

	Convey("Test1", t, func() {
		options := RateLimiterGroupOptions{
			Type: "LocalGroup",
			RateLimiter: &LocalGroupRateLimiterOptions{
				"DB.First": {
					Interval: 1 * time.Second,
					Burst:    1,
				},
			},
		}

		rateLimiterGroup, err := NewRateLimiterGroupWithOptions(&options)
		So(err, ShouldBeNil)

		for i := 0; i < 5; i++ {
			rateLimiterGroup.Wait(context.Background(), "DB.First")
			fmt.Println("hello world")
		}
	})

	Convey("Test2", t, func() {
		v := map[string]interface{}{
			"type": "LocalGroup",
			"rateLimiter": map[string]interface{}{
				"DB.First": map[string]interface{}{
					"interval": "1s",
					"burst":    1,
				},
			},
		}

		var options RateLimiterGroupOptions

		So(refx.InterfaceToStruct(v, &options, refx.WithCamelName()), ShouldBeNil)

		rateLimiterGroup, err := NewRateLimiterGroupWithOptions(&options, refx.WithCamelName())
		So(err, ShouldBeNil)

		for i := 0; i < 5; i++ {
			rateLimiterGroup.Wait(context.Background(), "DB.First")
			fmt.Println("hello world")
		}
	})
}

func TestReflect(t *testing.T) {
	Convey("TestReflect", t, func() {
		rt := reflect.TypeOf(NewLocalGroupRateLimiterWithOptions)
		fmt.Println(rt.Kind() == reflect.Func)

		fmt.Println(rt.NumIn())
		fmt.Println(rt.NumOut())
		fmt.Println(rt.Out(0))
		fmt.Println(rt.Out(0).Implements(reflect.TypeOf((*RateLimiterGroup)(nil)).Elem()))
		fmt.Println(rt.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()))

		rv := reflect.ValueOf(NewLocalGroupRateLimiterWithOptions).Call([]reflect.Value{reflect.ValueOf(&LocalGroupRateLimiterOptions{})})
		fmt.Println(rv)
	})
}
