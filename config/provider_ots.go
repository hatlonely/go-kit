package config

import (
	"context"
	"strings"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
	"github.com/hatlonely/go-kit/strx"
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

func NewOTSClient(options *OTSProviderOptions) (*tablestore.TableStoreClient, error) {
	if options.AccessKeyID != "" {
		return tablestore.NewClient(options.Endpoint, options.Instance, options.AccessKeyID, options.AccessKeySecret), nil
	}
	res, err := alics.ECSMetaDataRamSecurityCredentials()
	if err != nil {
		return nil, err
	}
	return tablestore.NewClientWithConfig(options.Endpoint, options.Instance, res.AccessKeyID, res.AccessKeySecret, res.SecurityToken, nil), nil
}

func NewOTSProviderWithOptions(options *OTSProviderOptions) (*OTSProvider, error) {
	otsCli, err := NewOTSClient(options)
	if err != nil {
		return nil, errors.Wrap(err, "NewOTSClient failed")
	}

	if res, err := otsCli.DescribeTable(&tablestore.DescribeTableRequest{
		TableName: options.Table,
	}); err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			return nil, errors.Wrap(err, "otsCli.DescribeTable failed")
		}
		req := &tablestore.CreateTableRequest{
			TableMeta: &tablestore.TableMeta{
				TableName: options.Table,
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
			return nil, errors.Wrap(err, "otsCli.CreateTable failed")
		}
	} else {
		if len(res.TableMeta.SchemaEntry) != 1 {
			return nil, errors.Errorf("ots primary key [%v] is not match options [Key]", strx.JsonMarshal(res.TableMeta.SchemaEntry))
		}
		if *res.TableMeta.SchemaEntry[0].Name != "Key" {
			return nil, errors.Errorf("ots primary key [%v] is not match options [Key]", strx.JsonMarshal(res.TableMeta.SchemaEntry))
		}
		if *res.TableMeta.SchemaEntry[0].Type != tablestore.PrimaryKeyType_STRING {
			return nil, errors.Errorf("table [%v] primary key should be string", options.Table)
		}
	}

	provider := &OTSProvider{
		options: options,
		events:  make(chan struct{}, 10),
		errors:  make(chan error, 10),
	}
	buf, ts, err := provider.otsGetRow()
	if err != nil {
		return nil, errors.Wrap(err, "otsGetRow failed")
	}
	provider.buf = buf
	provider.ts = ts

	return provider, nil
}

type OTSProvider struct {
	events  chan struct{}
	errors  chan error
	buf     []byte
	ts      int64
	options *OTSProviderOptions
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

func (p *OTSProvider) otsGetRow() ([]byte, int64, error) {
	otsCli, err := NewOTSClient(p.options)
	if err != nil {
		return nil, 0, errors.Wrap(err, "NewOTSClient failed")
	}
	res, err := otsCli.GetRow(&tablestore.GetRowRequest{
		SingleRowQueryCriteria: &tablestore.SingleRowQueryCriteria{
			TableName: p.options.Table,
			PrimaryKey: &tablestore.PrimaryKey{
				PrimaryKeys: []*tablestore.PrimaryKeyColumn{
					{ColumnName: "Key", Value: p.options.Key},
				},
			},
			MaxVersion: 1,
		},
	})
	if err != nil {
		return nil, 0, errors.Wrap(err, "ots.GetRow failed")
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

func (p *OTSProvider) otsPutRow(buf []byte) error {
	otsCli, err := NewOTSClient(p.options)
	if err != nil {
		return errors.Wrap(err, "NewOTSClient failed")
	}
	_, err = otsCli.PutRow(&tablestore.PutRowRequest{
		PutRowChange: &tablestore.PutRowChange{
			TableName: p.options.Table,
			PrimaryKey: &tablestore.PrimaryKey{
				PrimaryKeys: []*tablestore.PrimaryKeyColumn{
					{ColumnName: "Key", Value: p.options.Key},
				},
			},
			Columns: []tablestore.AttributeColumn{
				{ColumnName: "Val", Value: string(buf)},
			},
			Condition: &tablestore.RowCondition{RowExistenceExpectation: tablestore.RowExistenceExpectation_IGNORE},
		},
	})

	return errors.Wrap(err, "ots.PutRow failed")
}

func (p *OTSProvider) Dump(buf []byte) error {
	return p.otsPutRow(buf)
}

func (p *OTSProvider) EventLoop(ctx context.Context) error {
	go func() {
		ticker := time.NewTicker(p.options.Interval)
		defer ticker.Stop()

	out:
		for {
			select {
			case <-ticker.C:
				buf, ts, err := p.otsGetRow()
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
