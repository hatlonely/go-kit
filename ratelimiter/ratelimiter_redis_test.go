package ratelimiter

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/wrap"
)

func TestRedisRateLimiter(t *testing.T) {
	Convey("TestRedisRateLimiter", t, func() {
		r, err := NewRedisRateLimiterWithOptions(&RedisRateLimiterOptions{
			Redis: wrap.RedisClientWrapperOptions{
				Retry: wrap.RetryOptions{
					Attempts: 3,
					Delay:    time.Millisecond * 500,
				},
			},
			QPS: map[string]int{
				"Key1": 1,
			},
		})
		So(err, ShouldBeNil)

		for i := 0; i < 10; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
			err := r.WaitN(ctx, "Key1", 2)
			cancel()
			fmt.Println("hello world")
			So(err, ShouldBeNil)
		}
	})
}

func TestRedisRateLimiter_WaitN_Parallel(t *testing.T) {
	r, _ := NewRedisRateLimiterWithOptions(&RedisRateLimiterOptions{
		Redis: wrap.RedisClientWrapperOptions{
			Retry: wrap.RetryOptions{
				Attempts: 3,
				Delay:    time.Millisecond * 500,
			},
		},
		QPS: map[string]int{
			"Key1": 10,
		},
	})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			for i := 0; i < 10; i++ {
				err := r.WaitN(context.Background(), "Key1", 1)
				fmt.Println("hello world", err, i)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestNewRedisRateLimiterWithConfig(t *testing.T) {
	Convey("t", t, func() {
		cfg, err := config.NewConfigWithOptions(&config.Options{
			Provider: config.ProviderOptions{
				Type: "Memory",
				Options: &config.MemoryProviderOptions{
					Buffer: `{
"rateLimiter": {
    "redis": {
      "wrapper": {
        "enableTrace": true,
        "enableMetric": true,
      }
    },
    "qps": {
      "DB.First": 1
    }
  },
}`,
				},
			},
		})

		So(err, ShouldBeNil)
		v, ok := cfg.Get("")
		fmt.Println(v, ok)

		r, err := NewRedisRateLimiterWithConfig(cfg.Sub("rateLimiter"), refx.WithCamelName())
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)
	})
}
