package micro

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/time/rate"

	"github.com/hatlonely/go-kit/refx"
)

func TestNewRateLimiterWithOptions(t *testing.T) {
	Convey("TestNewRateLimiterWithOptions", t, func() {
		Convey("empty options case", func() {
			r, err := NewRateLimiterWithOptions(&RateLimiterOptions{})
			So(err, ShouldBeNil)
			So(r, ShouldBeNil)
		})

		r, err := NewRateLimiterWithOptions(&RateLimiterOptions{
			Type: "LocalGroup",
			Options: &LocalGroupRateLimiterOptions{
				"key1": {
					Interval: time.Second,
					Burst:    2,
				},
			},
		})
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)

		So(r.Allow(context.Background(), "key1"), ShouldBeNil)
	})
}

func Test123(t *testing.T) {
	Convey("Test", t, func() {
		limiter := rate.NewLimiter(rate.Every(200*time.Millisecond), 100)
		limiter.Allow()
	})

	Convey("Test1", t, func() {
		options := RateLimiterOptions{
			Type: "LocalGroup",
			Options: &LocalGroupRateLimiterOptions{
				"DB.First": {
					Interval: 20 * time.Millisecond,
					Burst:    1,
				},
			},
		}

		rateLimiterGroup, err := NewRateLimiterWithOptions(&options)
		So(err, ShouldBeNil)

		for i := 0; i < 5; i++ {
			rateLimiterGroup.Wait(context.Background(), "DB.First")
			fmt.Println("hello world")
		}
	})

	Convey("Test2", t, func() {
		options := RateLimiterOptions{
			Type: "LocalGroup",
			Options: map[string]interface{}{
				"DB.First": map[string]interface{}{
					"interval": "1s",
					"burst":    1,
				},
			},
		}

		rateLimiterGroup, err := NewRateLimiterWithOptions(&options, refx.WithCamelName())
		So(err, ShouldBeNil)

		for i := 0; i < 5; i++ {
			rateLimiterGroup.Wait(context.Background(), "DB.First")
			fmt.Println("hello world")
		}
	})
}

func TestRegex(t *testing.T) {
	Convey("TestRegex", t, func() {
		re := regexp.MustCompile(`^DB$`)
		fmt.Println(re.MatchString("*gorm.DB"))
	})
}

func TestReflect(t *testing.T) {
	Convey("TestReflect", t, func() {
		rt := reflect.TypeOf(NewLocalGroupRateLimiterWithOptions)
		fmt.Println(rt.Kind() == reflect.Func)

		fmt.Println(rt.NumIn())
		fmt.Println(rt.NumOut())
		fmt.Println(rt.Out(0))
		fmt.Println(rt.Out(0).Implements(reflect.TypeOf((*RateLimiter)(nil)).Elem()))
		fmt.Println(rt.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()))

		rv := reflect.ValueOf(NewLocalGroupRateLimiterWithOptions).Call([]reflect.Value{reflect.ValueOf(&LocalGroupRateLimiterOptions{})})
		fmt.Println(rv)
	})
}

func TestNewRateLimiterWithOptions1(t *testing.T) {
	Convey("TestNewRateLimiterWithOptions", t, func() {
		r, err := NewLocalGroupRateLimiterWithOptions(&LocalGroupRateLimiterOptions{
			"DB.First": {
				Interval: time.Second,
				Burst:    1,
			},
		})
		So(err, ShouldBeNil)
		RegisterRateLimiter("Test", r)

		rateLimiterGroup, err := NewRateLimiterWithOptions(&RateLimiterOptions{
			Type: "Test",
		})
		So(err, ShouldBeNil)

		for i := 0; i < 5; i++ {
			rateLimiterGroup.Wait(context.Background(), "DB.First")
			fmt.Println("hello world")
		}
	})
}
