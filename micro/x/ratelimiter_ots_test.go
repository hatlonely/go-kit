package microx

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/wrap"
)

// BenchmarkOTSRateLimiter_Wait-12    	     262	   3858786 ns/op
func BenchmarkOTSRateLimiter_Wait(b *testing.B) {
	r, _ := NewOTSRateLimiterWithOptions(&OTSRateLimiterOptions{
		OTS: wrap.OTSTableStoreClientWrapperOptions{
			OTS: wrap.OTSOptions{
				Endpoint:        "https://xx.cn-shanghai.ots.aliyuncs.com",
				AccessKeyID:     "xx",
				AccessKeySecret: "xx",
				InstanceName:    "xx",
			},
			Retry: micro.RetryOptions{
				Attempts:      3,
				Delay:         time.Millisecond * 500,
				LastErrorOnly: true,
			},
		},
		Table:      "RateLimiter",
		DefaultQPS: 9999999999,
		Window:     time.Second * 1,
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Wait(context.Background(), "key1")
		}
	})
}

func TestOTSRateLimiter(t *testing.T) {
	Convey("TestOTSRateLimiter", t, func() {
		r, err := NewOTSRateLimiterWithOptions(&OTSRateLimiterOptions{
			OTS: wrap.OTSTableStoreClientWrapperOptions{
				OTS: wrap.OTSOptions{
					Endpoint:        "https://test-name.cn-shanghai.ots.aliyuncs.com",
					AccessKeyID:     "test-ak",
					AccessKeySecret: "test-sk",
					InstanceName:    "test-name",
				},
				Retry: micro.RetryOptions{
					Attempts:      3,
					Delay:         time.Millisecond * 500,
					LastErrorOnly: true,
				},
			},
			Table:      "RateLimiter",
			DefaultQPS: 3,
			QPS: map[string]int{
				"key1": 2,
				"key2": 0,
			},
			Window: time.Second * 5,
		})
		So(err, ShouldBeNil)

		{
			time.Sleep(time.Until(time.Now().Truncate(time.Second).Add(time.Second)))
			now := time.Now()
			So(r.Allow(context.Background(), "key1"), ShouldBeNil)
			So(r.Allow(context.Background(), "key1"), ShouldBeNil)
			So(r.Allow(context.Background(), "key1"), ShouldEqual, micro.ErrFlowControl)
			So(r.Wait(context.Background(), "key1"), ShouldBeNil)
			So(time.Now().Sub(now), ShouldBeGreaterThan, 20*time.Millisecond)
		}
		{
			time.Sleep(time.Until(time.Now().Truncate(time.Second).Add(time.Second)))
			now := time.Now()
			So(r.Allow(context.Background(), "key3"), ShouldBeNil)
			So(r.Allow(context.Background(), "key3"), ShouldBeNil)
			So(r.Allow(context.Background(), "key3"), ShouldBeNil)
			So(r.Allow(context.Background(), "key3"), ShouldEqual, micro.ErrFlowControl)
			So(r.Wait(context.Background(), "key3"), ShouldBeNil)
			So(time.Now().Sub(now), ShouldBeGreaterThan, 20*time.Millisecond)
		}

		for i := 0; i < 10; i++ {
			So(r.Allow(context.Background(), "key2"), ShouldBeNil)
		}
		now := time.Now()
		for i := 0; i < 10; i++ {
			So(r.Wait(context.Background(), "key2"), ShouldBeNil)
		}
		So(time.Now().Sub(now), ShouldBeLessThan, 20*time.Millisecond)

		for i := 0; i < 30; i++ {
			So(r.Wait(context.Background(), "key3"), ShouldBeNil)
			fmt.Println(time.Now(), "hello world")
		}
	})
}
