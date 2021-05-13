package microx

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/v5/tablestore"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/micro"
	"github.com/hatlonely/go-kit/wrap"
)

type OTSParallelControllerOptions struct {
	OTS wrap.OTSTableStoreClientWrapperOptions
	// OTS 表名
	Table string
	// key 前缀，可当成命名空间使用
	Prefix string
	// MaxToken 计算规则
	// 1. key 在 map 中，直接返回 key 对应的 qps
	// 2. key 按 '|' 分割，第 0 个字符串作为 key，如果在 map 中，返回 qps
	// 3. 返回 DefaultQPS
	MaxToken map[string]int
	// MaxToken 中未匹配到，使用 DefaultMaxToken，默认为 0，不限流
	DefaultMaxToken int
	// 获取 token 失败时重试时间间隔最大值
	Interval time.Duration
}

type OTSParallelController struct {
	client            *wrap.OTSTableStoreClientWrapper
	options           *OTSParallelControllerOptions
	acquireScriptSha1 string
	releaseScriptSha1 string
}

func NewOTSParallelControllerWithOptions(options *OTSParallelControllerOptions) (*OTSParallelController, error) {
	if options.Interval == 0 {
		options.Interval = time.Microsecond
	}

	client, err := wrap.NewOTSTableStoreClientWrapperWithOptions(&options.OTS)
	if err != nil {
		return nil, errors.WithMessage(err, "wrap.NewOTSTableStoreClientWrapperWithOptions")
	}

	_, err = client.DescribeTable(context.Background(), &tablestore.DescribeTableRequest{
		TableName: options.Table,
	})
	if err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			return nil, errors.Wrap(err, "client.DescribeTable failed")
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

	return &OTSParallelController{client: client, options: options}, nil
}

func (c *OTSParallelController) TryAcquire(ctx context.Context, key string) (int, error) {
	maxToken := c.calculateMaxToken(key)
	if maxToken == 0 {
		return 0, nil
	}

	key = c.generateKey(key)
	_, err := c.client.UpdateRow(ctx, &tablestore.UpdateRowRequest{
		UpdateRowChange: &tablestore.UpdateRowChange{
			TableName: c.options.Table,
			PrimaryKey: &tablestore.PrimaryKey{
				PrimaryKeys: []*tablestore.PrimaryKeyColumn{
					{ColumnName: "Key", Value: key},
				},
			},
			Columns: []tablestore.ColumnToUpdate{
				{ColumnName: "Val", Value: int64(1), Type: tablestore.INCREMENT, HasType: true, IgnoreValue: false},
			},
			Condition: &tablestore.RowCondition{
				RowExistenceExpectation: tablestore.RowExistenceExpectation_IGNORE,
				ColumnCondition:         tablestore.NewSingleColumnCondition("Val", tablestore.CT_LESS_THAN, int64(maxToken)),
			},
		},
	})
	if err != nil {
		switch e := err.(type) {
		case *tablestore.OtsError:
			if e.Code != "OTSConditionCheckFail" {
				return 0, errors.Wrap(err, "query ots failed")
			}
		default:
			return 0, errors.Wrap(err, "query ots failed")
		}
		return 0, micro.ErrParallelControl
	}

	return 0, nil
}

func (c *OTSParallelController) Acquire(ctx context.Context, key string) (int, error) {
	maxToken := c.calculateMaxToken(key)
	if maxToken == 0 {
		return 0, nil
	}

	key = c.generateKey(key)
	for {
		_, err := c.client.UpdateRow(ctx, &tablestore.UpdateRowRequest{
			UpdateRowChange: &tablestore.UpdateRowChange{
				TableName: c.options.Table,
				PrimaryKey: &tablestore.PrimaryKey{
					PrimaryKeys: []*tablestore.PrimaryKeyColumn{
						{ColumnName: "Key", Value: key},
					},
				},
				Columns: []tablestore.ColumnToUpdate{
					{ColumnName: "Val", Value: int64(1), Type: tablestore.INCREMENT, HasType: true, IgnoreValue: false},
				},
				Condition: &tablestore.RowCondition{
					RowExistenceExpectation: tablestore.RowExistenceExpectation_IGNORE,
					ColumnCondition:         tablestore.NewSingleColumnCondition("Val", tablestore.CT_LESS_THAN, int64(maxToken)),
				},
			},
		})
		if err != nil {
			switch e := err.(type) {
			case *tablestore.OtsError:
				if e.Code != "OTSConditionCheckFail" {
					return 0, errors.Wrap(err, "query ots failed")
				}
			default:
				return 0, errors.Wrap(err, "query ots failed")
			}
		} else {
			return 0, nil
		}
		select {
		case <-ctx.Done():
			return 0, micro.ErrContextCancel
		case <-time.After(time.Duration(rand.Int63n(c.options.Interval.Microseconds())) * time.Microsecond):
		}
	}
}

func (c *OTSParallelController) Release(ctx context.Context, key string, token int) error {
	maxToken := c.calculateMaxToken(key)
	if maxToken == 0 {
		return nil
	}

	key = c.generateKey(key)
	_, err := c.client.UpdateRow(ctx, &tablestore.UpdateRowRequest{
		UpdateRowChange: &tablestore.UpdateRowChange{
			TableName: c.options.Table,
			PrimaryKey: &tablestore.PrimaryKey{
				PrimaryKeys: []*tablestore.PrimaryKeyColumn{
					{ColumnName: "Key", Value: key},
				},
			},
			Columns: []tablestore.ColumnToUpdate{
				{ColumnName: "Val", Value: int64(-1), Type: tablestore.INCREMENT, HasType: true, IgnoreValue: false},
			},
			Condition: &tablestore.RowCondition{
				RowExistenceExpectation: tablestore.RowExistenceExpectation_IGNORE,
				ColumnCondition:         tablestore.NewSingleColumnCondition("Val", tablestore.CT_GREATER_THAN, int64(0)),
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
		return nil
	}

	return nil
}

func (c *OTSParallelController) generateKey(key string) string {
	if c.options.Prefix == "" {
		return key
	}
	return fmt.Sprintf("%s_%s", c.options.Prefix, key)
}

func (c *OTSParallelController) calculateMaxToken(key string) int {
	if maxToken, ok := c.options.MaxToken[key]; ok {
		return maxToken
	}

	if idx := strings.Index(key, "|"); idx != -1 {
		if maxToken, ok := c.options.MaxToken[key[:idx]]; ok {
			return maxToken
		}
	}

	return c.options.DefaultMaxToken
}
