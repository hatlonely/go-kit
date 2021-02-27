package microx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/wrap"
)

type OTSRateLimiterOptions struct {
	OTS wrap.OTSTableStoreClientWrapperOptions
	// 窗口长度
	Window time.Duration `dft:"1s"`
	// OTS 表名
	Table string
	// key 前缀，可当成命名空间使用
	Prefix string
	// QPS 计算规则
	// 1. key 在 map 中，直接返回 key 对应的 qps
	// 2. key 按 '|' 分割，第 0 个字符串作为 key，如果在 map 中，返回 qps
	// 3. 返回 DefaultQPS
	QPS map[string]int
	// QPS 中未匹配到，使用 DefaultQPS，默认为 0，不限流
	DefaultQPS int
}

type OTSRateLimiter struct {
	client  *wrap.OTSTableStoreClientWrapper
	options *OTSRateLimiterOptions
}

func NewOTSRateLimiterWithConfig(cfg *config.Config, opts ...refx.Option) (*OTSRateLimiter, error) {
	var options OTSRateLimiterOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.WithMessage(err, "cfg.Unmarshal failed")
	}

	refxOptions := refx.NewOptions(opts...)
	client, err := wrap.NewOTSTableStoreClientWrapperWithConfig(cfg.Sub(refxOptions.FormatKey("OTS")), opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "NewOTSRateLimiterWithConfig failed")
	}
	r := &OTSRateLimiter{client: client, options: &options}

	cfg.AddOnChangeHandler(func(cfg *config.Config) error {
		var options OTSRateLimiterOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.WithMessage(err, "cfg.Unmarshal failed")
		}
		r.options = &options
		return nil
	})

	return r, nil
}

func NewOTSRateLimiterWithOptions(options *OTSRateLimiterOptions) (*OTSRateLimiter, error) {
	client, err := wrap.NewOTSTableStoreClientWrapperWithOptions(&options.OTS)
	if err != nil {
		return nil, errors.Wrap(err, "NewOTSClientWrapperWithOptions failed")
	}

	_, err = client.DescribeTable(context.Background(), &tablestore.DescribeTableRequest{
		TableName: options.Table,
	})
	if err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			return nil, err
		}
		req := &tablestore.CreateTableRequest{
			TableMeta: &tablestore.TableMeta{
				TableName: options.Table,
			},
			TableOption: &tablestore.TableOption{
				TimeToAlive: 86400,
				MaxVersion:  1,
			},
			ReservedThroughput: &tablestore.ReservedThroughput{},
		}
		req.TableMeta.AddPrimaryKeyColumn("Key", tablestore.PrimaryKeyType_STRING)
		req.TableMeta.AddDefinedColumn("Val", tablestore.DefinedColumn_INTEGER)

		if _, err := client.CreateTable(context.Background(), req); err != nil {
			return nil, errors.Wrap(err, "client.CreateTable failed")
		}
	}

	return &OTSRateLimiter{client: client, options: options}, nil
}

func (r *OTSRateLimiter) Allow(ctx context.Context, key string) error {
	qps := r.calculateQPS(key)
	if qps == 0 {
		return nil
	}

	now := time.Now().Truncate(r.options.Window)
	tsKey := fmt.Sprintf("%s_%d", r.generateKey(key), now.Unix())
	_, err := r.client.UpdateRow(ctx, &tablestore.UpdateRowRequest{
		UpdateRowChange: &tablestore.UpdateRowChange{
			TableName: r.options.Table,
			PrimaryKey: &tablestore.PrimaryKey{
				PrimaryKeys: []*tablestore.PrimaryKeyColumn{
					{ColumnName: "Key", Value: tsKey},
				},
			},
			Columns: []tablestore.ColumnToUpdate{
				{ColumnName: "Val", Value: int64(1), Type: tablestore.INCREMENT, HasType: true, IgnoreValue: false},
			},
			Condition: &tablestore.RowCondition{
				RowExistenceExpectation: tablestore.RowExistenceExpectation_IGNORE,
				ColumnCondition:         tablestore.NewSingleColumnCondition("Val", tablestore.CT_LESS_THAN, int64(qps)),
			},
		},
	})
	if err != nil {
		switch e := err.(type) {
		case *tablestore.OtsError:
			if e.Code != "OTSConditionCheckFail" {
				return errors.Wrap(err, "query ots failed")
			}
		default:
			return errors.Wrap(err, "query ots failed")
		}
		return micro.ErrFlowControl
	}

	return nil
}

func (r *OTSRateLimiter) Wait(ctx context.Context, key string) error {
	return r.WaitN(ctx, key, 1)
}

func (r *OTSRateLimiter) WaitN(ctx context.Context, key string, n int) error {
	qps := r.calculateQPS(key)
	if qps == 0 {
		return nil
	}

	for {
		now := time.Now().Truncate(r.options.Window)
		tsKey := fmt.Sprintf("%s_%d", r.generateKey(key), now.Unix())
		res, err := r.client.UpdateRow(ctx, &tablestore.UpdateRowRequest{
			UpdateRowChange: &tablestore.UpdateRowChange{
				TableName: r.options.Table,
				PrimaryKey: &tablestore.PrimaryKey{
					PrimaryKeys: []*tablestore.PrimaryKeyColumn{
						{ColumnName: "Key", Value: tsKey},
					},
				},
				Columns: []tablestore.ColumnToUpdate{
					{ColumnName: "Val", Value: int64(1), Type: tablestore.INCREMENT, HasType: true, IgnoreValue: false},
				},
				Condition: &tablestore.RowCondition{
					RowExistenceExpectation: tablestore.RowExistenceExpectation_IGNORE,
					ColumnCondition:         tablestore.NewSingleColumnCondition("Val", tablestore.CT_LESS_THAN, int64(qps)),
				},
				ColumnNamesToReturn: []string{"Val"},
				ReturnType:          tablestore.ReturnType_RT_AFTER_MODIFY,
			},
		})
		if err != nil {
			switch e := err.(type) {
			case *tablestore.OtsError:
				if e.Code != "OTSConditionCheckFail" {
					return errors.Wrap(err, "query ots failed")
				}
			default:
				return errors.Wrap(err, "query ots failed")
			}
		} else {
			val := int(res.Columns[0].Value.(int64))
			if val <= qps {
				break
			}
			if val-qps <= n {
				n = val - qps
			}
		}

		select {
		case <-time.After(time.Until(now.Add(time.Second))):
		case <-ctx.Done():
			return errors.New("cancel by ctx.Done")
		}
	}
	return nil
}

func (r *OTSRateLimiter) generateKey(key string) string {
	if r.options.Prefix == "" {
		return key
	}
	return fmt.Sprintf("%s_%s", r.options.Prefix, key)
}

func (r *OTSRateLimiter) calculateQPS(key string) int {
	if qps, ok := r.options.QPS[key]; ok {
		return qps
	}

	if idx := strings.Index(key, "|"); idx != -1 {
		if qps, ok := r.options.QPS[key[:idx]]; ok {
			return qps
		}
	}

	return r.options.DefaultQPS
}
