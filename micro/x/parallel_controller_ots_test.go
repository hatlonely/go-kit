package microx

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/wrap"
)

// BenchmarkOTSParallelController_Acquire-12    	      87	  11784280 ns/op
func BenchmarkOTSParallelController_Acquire(b *testing.B) {
	ctl, _ := NewOTSParallelControllerWithOptions(&OTSParallelControllerOptions{
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
		Table:           "ParallelController",
		Prefix:          "redis",
		DefaultMaxToken: 999999999,
		Interval:        time.Second,
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token, _ := ctl.Acquire(context.Background(), "key1")
			_ = ctl.Release(context.Background(), "key1", token)
		}
	})
}

func TestOTSParallelController_Acquire_Release(t *testing.T) {
	Convey("TestOTSParallelController_Acquire_Release", t, func(c C) {
		for i := 3; i < 4; i++ {
			ctl, err := NewOTSParallelControllerWithOptions(&OTSParallelControllerOptions{
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
				Table:  "ParallelController",
				Prefix: "test",
				MaxToken: map[string]int{
					"key1": i,
				},
				Interval: time.Millisecond,
			})
			c.So(err, ShouldBeNil)
			fmt.Println(ctl.client)
			ctl.client.DeleteRow(context.Background(), &tablestore.DeleteRowRequest{
				DeleteRowChange: &tablestore.DeleteRowChange{
					TableName: ctl.options.Table,
					PrimaryKey: &tablestore.PrimaryKey{
						PrimaryKeys: []*tablestore.PrimaryKeyColumn{
							{ColumnName: "Key", Value: "test_key1"},
						},
					},
					Condition: &tablestore.RowCondition{RowExistenceExpectation: tablestore.RowExistenceExpectation_IGNORE},
				},
			})
			var wg sync.WaitGroup
			var m int64
			for j := 0; j < 10; j++ {
				wg.Add(1)
				go func(i int) {
					for k := 0; k < 10; k++ {
						res, err := ctl.Acquire(context.Background(), "key1")
						c.So(res, ShouldEqual, 0)
						c.So(err, ShouldBeNil)
						atomic.AddInt64(&m, 1)
						c.So(m, ShouldBeLessThanOrEqualTo, i)
						c.So(m, ShouldBeGreaterThanOrEqualTo, 0)
						atomic.AddInt64(&m, -1)
						c.So(m, ShouldBeLessThanOrEqualTo, i)
						c.So(m, ShouldBeGreaterThanOrEqualTo, 0)
						err = ctl.Release(context.Background(), "key1", res)
						c.So(err, ShouldBeNil)
					}
					wg.Done()
				}(i)
			}
			wg.Wait()
		}
	})
}

func TestOTSParallelController_TryAcquire(t *testing.T) {
	Convey("TestOTSParallelController_TryAcquire", t, func() {
		ctl, err := NewOTSParallelControllerWithOptions(&OTSParallelControllerOptions{
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
			Table:  "ParallelController",
			Prefix: "test",
			MaxToken: map[string]int{
				"key1": 2,
			},
			Interval: time.Second,
		})
		So(err, ShouldBeNil)
		ctl.client.DeleteRow(context.Background(), &tablestore.DeleteRowRequest{
			DeleteRowChange: &tablestore.DeleteRowChange{
				TableName: ctl.options.Table,
				PrimaryKey: &tablestore.PrimaryKey{
					PrimaryKeys: []*tablestore.PrimaryKeyColumn{
						{ColumnName: "Key", Value: "test_key1"},
					},
				},
				Condition: &tablestore.RowCondition{RowExistenceExpectation: tablestore.RowExistenceExpectation_IGNORE},
			},
		})
		res, err := ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldBeNil)
		So(res, ShouldEqual, 0)
		res, err = ctl.Acquire(context.Background(), "key1")
		So(err, ShouldBeNil)
		So(res, ShouldEqual, 0)
		res, err = ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldEqual, micro.ErrParallelControl)
		So(res, ShouldEqual, 0)
		So(ctl.Release(context.Background(), "key1", res), ShouldBeNil)
		res, err = ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldBeNil)
		So(res, ShouldEqual, 0)
		res, err = ctl.TryAcquire(context.Background(), "key1")
		So(err, ShouldEqual, micro.ErrParallelControl)
		So(res, ShouldEqual, 0)
	})
}

func TestOTSParallelController_ContextCancel(t *testing.T) {
	Convey("TestOTSParallelController_ContextCancel", t, func() {
		ctl, err := NewOTSParallelControllerWithOptions(&OTSParallelControllerOptions{
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
			Table:  "ParallelController",
			Prefix: "test",
			MaxToken: map[string]int{
				"key1": 2,
			},
			Interval: time.Second,
		})
		So(err, ShouldBeNil)
		ctl.client.DeleteRow(context.Background(), &tablestore.DeleteRowRequest{
			DeleteRowChange: &tablestore.DeleteRowChange{
				TableName: ctl.options.Table,
				PrimaryKey: &tablestore.PrimaryKey{
					PrimaryKeys: []*tablestore.PrimaryKeyColumn{
						{ColumnName: "Key", Value: "test_key1"},
					},
				},
				Condition: &tablestore.RowCondition{RowExistenceExpectation: tablestore.RowExistenceExpectation_IGNORE},
			},
		})
		Convey("key not match", func() {
			for i := 0; i < 10; i++ {
				res, err := ctl.Acquire(context.Background(), "key2")
				So(err, ShouldBeNil)
				So(res, ShouldEqual, 0)
			}
			for i := 0; i < 10; i++ {
				res, err := ctl.TryAcquire(context.Background(), "key2")
				So(err, ShouldBeNil)
				So(res, ShouldEqual, 0)
			}
			for i := 0; i < 10; i++ {
				err := ctl.Release(context.Background(), "key2", 0)
				So(err, ShouldBeNil)
			}
		})

		Convey("context cancel", func() {
			_, _ = ctl.Acquire(context.Background(), "key1")
			_, _ = ctl.Acquire(context.Background(), "key1")
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			_, err := ctl.Acquire(ctx, "key1")
			So(err, ShouldEqual, micro.ErrContextCancel)
		})
	})
}
