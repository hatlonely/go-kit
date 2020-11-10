package config

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/strx"
)

type OTSLegacyProviderOptions struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	Instance        string
	Table           string
	PrimaryKeys     []string
	Interval        time.Duration
}

type OTSLegacyProvider struct {
	otsCli      *tablestore.TableStoreClient
	events      chan struct{}
	errors      chan error
	table       string
	primaryKeys []string
	interval    time.Duration

	buf []byte
	ts  int64
}

func NewOTSLegacyProviderWithOptions(options *OTSLegacyProviderOptions) (*OTSLegacyProvider, error) {
	otsCli := tablestore.NewClient(options.Endpoint, options.Instance, options.AccessKeyID, options.AccessKeySecret)
	if _, err := otsCli.DescribeTable(&tablestore.DescribeTableRequest{
		TableName: options.Table,
	}); err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			return nil, err
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
		for _, key := range options.PrimaryKeys {
			req.TableMeta.AddPrimaryKeyColumn(key, tablestore.PrimaryKeyType_STRING)
		}
		if _, err := otsCli.CreateTable(req); err != nil {
			return nil, err
		}
	}

	buf, ts, err := OTSLegacyGetRange(otsCli, options.Table, options.PrimaryKeys)
	if err != nil {
		return nil, err
	}

	return &OTSLegacyProvider{
		otsCli:      otsCli,
		table:       options.Table,
		primaryKeys: options.PrimaryKeys,
		interval:    options.Interval,
		buf:         buf,
		ts:          ts,
		events:      make(chan struct{}, 10),
		errors:      make(chan error, 10),
	}, nil
}

func (p *OTSLegacyProvider) Events() <-chan struct{} {
	return p.events
}

func (p *OTSLegacyProvider) Errors() <-chan error {
	return p.errors
}

func (p *OTSLegacyProvider) Load() ([]byte, error) {
	return p.buf, nil
}

func (p *OTSLegacyProvider) Dump(buf []byte) error {
	return OTSLegacyPutRow(p.otsCli, p.table, p.primaryKeys, buf)
}

func (p *OTSLegacyProvider) EventLoop(ctx context.Context) error {
	go func() {
		ticker := time.NewTicker(p.interval)
		defer ticker.Stop()

	out:
		for {
			select {
			case <-ticker.C:
				buf, ts, err := OTSLegacyGetRange(p.otsCli, p.table, p.primaryKeys)
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

func OTSLegacyGetRange(otsCli *tablestore.TableStoreClient, table string, primaryKeys []string) ([]byte, int64, error) {
	req := &tablestore.GetRangeRequest{
		RangeRowQueryCriteria: &tablestore.RangeRowQueryCriteria{
			TableName:       table,
			MaxVersion:      1,
			Direction:       tablestore.FORWARD,
			Limit:           20,
			EndPrimaryKey:   &tablestore.PrimaryKey{},
			StartPrimaryKey: &tablestore.PrimaryKey{},
		},
	}
	for _, key := range primaryKeys {
		req.RangeRowQueryCriteria.StartPrimaryKey.AddPrimaryKeyColumnWithMinValue(key)
		req.RangeRowQueryCriteria.EndPrimaryKey.AddPrimaryKeyColumnWithMaxValue(key)
	}

	var ts int64
	v := map[string]interface{}{}
	p := map[string]interface{}{
		"$": v,
	}
	key := "$"
	for req.RangeRowQueryCriteria.StartPrimaryKey != nil {
		res, err := otsCli.GetRange(req)
		if err != nil {
			return nil, 0, errors.Wrap(err, "otsCli.GetRange failed")
		}

		for _, row := range res.Rows {
			m := v
			for _, pks := range row.PrimaryKey.PrimaryKeys {
				if pks.Value.(string) == "_" {
					continue
				}
				key = pks.Value.(string)
				if _, ok := m[key]; !ok {
					m[key] = map[string]interface{}{}
				}
				p = m
				m = m[key].(map[string]interface{})
			}
			if len(row.Columns) == 1 && row.Columns[0].ColumnName == "_" {
				col := row.Columns[0]
				switch val := col.Value.(type) {
				case string:
					var obj interface{}
					if err := json.Unmarshal([]byte(val), &obj); err != nil {
						p[key] = val
					} else {
						p[key] = obj
					}
				default:
					p[key] = val
				}
				if ts < col.Timestamp {
					ts = col.Timestamp
				}
			} else {
				for _, col := range row.Columns {
					switch val := col.Value.(type) {
					case string:
						var obj interface{}
						if err := json.Unmarshal([]byte(val), &obj); err != nil {
							m[col.ColumnName] = val
						} else {
							m[col.ColumnName] = obj
						}
					default:
						m[col.ColumnName] = col.Value
					}
					if ts < col.Timestamp {
						ts = col.Timestamp
					}
				}
			}
		}

		req.RangeRowQueryCriteria.StartPrimaryKey = res.NextStartPrimaryKey
	}

	buf, _ := json.Marshal(v)
	return buf, ts, nil
}

func OTSLegacyPutRow(otsCli *tablestore.TableStoreClient, table string, primaryKeys []string, buf []byte) error {
	var v interface{}
	d := json.NewDecoder(bytes.NewBuffer(buf))
	d.UseNumber()
	if err := d.Decode(&v); err != nil {
		return errors.Wrap(err, "unmarshal failed")
	}
	return otsLegacyPutRowRecursive(otsCli, table, primaryKeys, []string{}, v, len(primaryKeys))
}

func otsLegacyPutRowRecursive(otsCli *tablestore.TableStoreClient, table string, primaryKeys []string, values []string, v interface{}, idx int) error {
	if idx == 0 {
		req := &tablestore.PutRowRequest{
			PutRowChange: &tablestore.PutRowChange{
				TableName:  table,
				PrimaryKey: &tablestore.PrimaryKey{},
				Columns:    []tablestore.AttributeColumn{},
				Condition:  &tablestore.RowCondition{RowExistenceExpectation: tablestore.RowExistenceExpectation_IGNORE},
			},
		}
		for i, key := range primaryKeys {
			req.PutRowChange.PrimaryKey.AddPrimaryKeyColumn(key, values[i])
		}
		m, ok := v.(map[string]interface{})
		if !ok {
			switch n := v.(type) {
			case float64, int64, string, bool:
				req.PutRowChange.AddColumn("_", n)
			case int:
				req.PutRowChange.AddColumn("_", int64(n))
			default:
				req.PutRowChange.AddColumn("_", strx.JsonMarshal(n))
			}
		} else {
			for key, val := range m {
				switch n := val.(type) {
				case float64, int64, string, bool:
					req.PutRowChange.AddColumn(key, n)
				case int:
					req.PutRowChange.AddColumn(key, int64(n))
				default:
					req.PutRowChange.AddColumn(key, strx.JsonMarshal(n))
				}
			}
		}
		if _, err := otsCli.PutRow(req); err != nil {
			return err
		}
		return nil
	}

	m, ok := v.(map[string]interface{})
	if !ok {
		return otsLegacyPutRowRecursive(otsCli, table, primaryKeys, append(values, "_"), v, idx-1)
	}

	for key, val := range m {
		if err := otsLegacyPutRowRecursive(otsCli, table, primaryKeys, append(values, key), val, idx-1); err != nil {
			return err
		}
	}

	return nil
}
