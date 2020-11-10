package config

import (
	"context"
	"strings"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

type OTSProviderOptions struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	Instance        string
	Table           string
	Key             string
	Interval        time.Duration
}

func NewOTSProviderWithOptions(options *OTSProviderOptions) (*OTSProvider, error) {
	otsCli := tablestore.NewClient(options.Endpoint, options.Instance, options.AccessKeyID, options.AccessKeySecret)
	return NewOTSProvider(otsCli, options.Table, options.Key, options.Interval)
}

func NewOTSProvider(otsCli *tablestore.TableStoreClient, table string, key string, interval time.Duration) (*OTSProvider, error) {
	if _, err := otsCli.DescribeTable(&tablestore.DescribeTableRequest{
		TableName: table,
	}); err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			return nil, err
		}
		req := &tablestore.CreateTableRequest{
			TableMeta: &tablestore.TableMeta{
				TableName: table,
			},
			TableOption: &tablestore.TableOption{
				TimeToAlive: -1,
				MaxVersion:  1,
			},
			ReservedThroughput: &tablestore.ReservedThroughput{},
		}
		req.TableMeta.AddPrimaryKeyColumn("Key", tablestore.PrimaryKeyType_STRING)
		req.TableMeta.AddDefinedColumn("Val", tablestore.DefinedColumn_STRING)
		if _, err := otsCli.CreateTable(req); err != nil {
			return nil, err
		}
	}

	provider := &OTSProvider{
		otsCli:   otsCli,
		table:    table,
		key:      key,
		interval: interval,
		events:   make(chan struct{}, 10),
		errors:   make(chan error, 10),
	}
	buf, ts, err := provider.otsGetRow(otsCli, table, key)
	if err != nil {
		return nil, err
	}
	provider.buf = buf
	provider.ts = ts

	return provider, nil
}

type OTSProvider struct {
	otsCli   *tablestore.TableStoreClient
	events   chan struct{}
	errors   chan error
	table    string
	key      string
	interval time.Duration

	buf []byte
	ts  int64
}

func (p *OTSProvider) Events() <-chan struct{} {
	return p.events
}

func (p *OTSProvider) Errors() <-chan error {
	return p.errors
}

func (p *OTSProvider) Load() ([]byte, error) {
	return p.buf, nil
}

func (p *OTSProvider) otsGetRow(otsCli *tablestore.TableStoreClient, table string, key string) ([]byte, int64, error) {
	res, err := otsCli.GetRow(&tablestore.GetRowRequest{
		SingleRowQueryCriteria: &tablestore.SingleRowQueryCriteria{
			TableName: table,
			PrimaryKey: &tablestore.PrimaryKey{
				PrimaryKeys: []*tablestore.PrimaryKeyColumn{
					{ColumnName: "Key", Value: key},
				},
			},
			MaxVersion: 1,
		},
	})
	if err != nil {
		return nil, 0, err
	}

	var val string
	var ts int64
	for _, col := range res.Columns {
		if col.ColumnName == "Val" {
			val = col.Value.(string)
			ts = col.Timestamp
		}
	}

	return []byte(val), ts, nil
}

func (p *OTSProvider) otsPutRow(otsCli *tablestore.TableStoreClient, table string, key string, buf []byte) error {
	_, err := otsCli.PutRow(&tablestore.PutRowRequest{
		PutRowChange: &tablestore.PutRowChange{
			TableName: table,
			PrimaryKey: &tablestore.PrimaryKey{
				PrimaryKeys: []*tablestore.PrimaryKeyColumn{
					{ColumnName: "Key", Value: key},
				},
			},
			Columns: []tablestore.AttributeColumn{
				{ColumnName: "Val", Value: string(buf)},
			},
			Condition: &tablestore.RowCondition{RowExistenceExpectation: tablestore.RowExistenceExpectation_IGNORE},
		},
	})

	return err
}

func (p *OTSProvider) Dump(buf []byte) error {
	return p.otsPutRow(p.otsCli, p.table, p.key, buf)
}

func (p *OTSProvider) EventLoop(ctx context.Context) error {
	go func() {
		ticker := time.NewTicker(p.interval)
		defer ticker.Stop()

	out:
		for {
			select {
			case <-ticker.C:
				buf, ts, err := p.otsGetRow(p.otsCli, p.table, p.key)
				if err != nil {
					p.errors <- err
					continue
				}
				if ts == p.ts {
					continue
				}
				p.ts = ts
				p.buf = buf
				p.events <- struct{}{}
			case <-ctx.Done():
				break out
			}
		}
	}()

	return nil
}
