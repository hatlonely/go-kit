package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/v5/tablestore"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
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
	events  chan struct{}
	errors  chan error
	buf     []byte
	ts      int64
	options *OTSLegacyProviderOptions
}

func NewOTSLegacyProviderWithOptions(options *OTSLegacyProviderOptions) (*OTSLegacyProvider, error) {
	if options.Instance == "" {
		return nil, errors.New("OTSProviderOptions.Instance is required")
	}
	if options.Endpoint == "" {
		regionID, err := alics.ECSMetaDataRegionID()
		if err != nil {
			return nil, errors.Wrap(err, "ECSMetaDataRegionID failed")
		}
		options.Endpoint = fmt.Sprintf("https://%s.%s.ots.aliyuncs.com", options.Instance, regionID)
	}

	otsCli, err := newOTSClient(options.Endpoint, options.Instance, options.AccessKeyID, options.AccessKeySecret)
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
		for _, key := range options.PrimaryKeys {
			req.TableMeta.AddPrimaryKeyColumn(key, tablestore.PrimaryKeyType_STRING)
		}
		if _, err := otsCli.CreateTable(req); err != nil {
			return nil, errors.Wrap(err, "otsCli.CreateTable failed")
		}
	} else {
		if len(res.TableMeta.SchemaEntry) != len(options.PrimaryKeys) {
			return nil, errors.Errorf("ots primary key [%v] is not match options [%v]", strx.JsonMarshal(res.TableMeta.SchemaEntry), options.PrimaryKeys)
		}
		for i, pk := range res.TableMeta.SchemaEntry {
			if *pk.Name != options.PrimaryKeys[i] {
				return nil, errors.Errorf("ots primary key [%v] is not match options [%v]", strx.JsonMarshal(res.TableMeta.SchemaEntry), options.PrimaryKeys)
			}
			if *pk.Type != tablestore.PrimaryKeyType_STRING {
				return nil, errors.Errorf("table [%v] primary key should be string", options.Table)
			}
		}
	}

	provider := &OTSLegacyProvider{
		options: options,
		events:  make(chan struct{}, 10),
		errors:  make(chan error, 10),
	}
	buf, ts, err := provider.otsGetRange()
	if err != nil {
		return nil, errors.Wrap(err, "otsGetRange failed")
	}
	provider.buf = buf
	provider.ts = ts

	return provider, nil
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
	return p.otsPutRow(buf)
}

func (p *OTSLegacyProvider) EventLoop(ctx context.Context) error {
	go func() {
		ticker := time.NewTicker(p.options.Interval)
		defer ticker.Stop()

	out:
		for {
			select {
			case <-ticker.C:
				buf, ts, err := p.otsGetRange()
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

func (p *OTSLegacyProvider) otsGetRange() ([]byte, int64, error) {
	otsCli, err := newOTSClient(p.options.Endpoint, p.options.Instance, p.options.AccessKeyID, p.options.AccessKeySecret)
	if err != nil {
		return nil, 0, errors.Wrap(err, "newOTSClient failed")
	}

	req := &tablestore.GetRangeRequest{
		RangeRowQueryCriteria: &tablestore.RangeRowQueryCriteria{
			TableName:       p.options.Table,
			MaxVersion:      1,
			Direction:       tablestore.FORWARD,
			Limit:           20,
			EndPrimaryKey:   &tablestore.PrimaryKey{},
			StartPrimaryKey: &tablestore.PrimaryKey{},
		},
	}
	for _, key := range p.options.PrimaryKeys {
		req.RangeRowQueryCriteria.StartPrimaryKey.AddPrimaryKeyColumnWithMinValue(key)
		req.RangeRowQueryCriteria.EndPrimaryKey.AddPrimaryKeyColumnWithMaxValue(key)
	}

	var ts int64
	v := map[string]interface{}{}
	prev := map[string]interface{}{
		"$": v,
	}
	key := "$"
	for req.RangeRowQueryCriteria.StartPrimaryKey != nil {
		res, err := otsCli.GetRange(req)
		if err != nil {
			return nil, 0, errors.Wrap(err, "otsCli.GetRange failed")
		}

		for _, row := range res.Rows {
			curr := v
			for _, pks := range row.PrimaryKey.PrimaryKeys {
				if pks.Value.(string) == "_" {
					continue
				}
				key = pks.Value.(string)
				if _, ok := curr[key]; !ok {
					curr[key] = map[string]interface{}{}
				}
				prev = curr
				curr = curr[key].(map[string]interface{})
			}
			if len(row.Columns) == 1 && row.Columns[0].ColumnName == "_" {
				col := row.Columns[0]
				switch val := col.Value.(type) {
				case string:
					var obj interface{}
					if err := json.Unmarshal([]byte(val), &obj); err != nil {
						prev[key] = val
					} else {
						prev[key] = obj
					}
				default:
					prev[key] = val
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
							curr[col.ColumnName] = val
						} else {
							curr[col.ColumnName] = obj
						}
					default:
						curr[col.ColumnName] = col.Value
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

func (p *OTSLegacyProvider) otsPutRow(buf []byte) error {
	otsCli, err := newOTSClient(p.options.Endpoint, p.options.Instance, p.options.AccessKeyID, p.options.AccessKeySecret)
	if err != nil {
		return errors.Wrap(err, "newOTSClient failed")
	}
	var v interface{}
	d := json.NewDecoder(bytes.NewBuffer(buf))
	d.UseNumber()
	if err := d.Decode(&v); err != nil {
		return errors.Wrap(err, "unmarshal failed")
	}
	return p.otsPutRowRecursive(otsCli, p.options.Table, p.options.PrimaryKeys, []string{}, v, len(p.options.PrimaryKeys))
}

func (p *OTSLegacyProvider) otsPutRowRecursive(otsCli *tablestore.TableStoreClient, table string, primaryKeys []string, values []string, v interface{}, idx int) error {
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
				req.PutRowChange.AddColumn("_", strx.JsonMarshalSortKeys(n))
			}
		} else {
			for key, val := range m {
				switch n := val.(type) {
				case float64, int64, string, bool:
					req.PutRowChange.AddColumn(key, n)
				case int:
					req.PutRowChange.AddColumn(key, int64(n))
				default:
					req.PutRowChange.AddColumn(key, strx.JsonMarshalSortKeys(n))
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
		return p.otsPutRowRecursive(otsCli, table, primaryKeys, append(values, "_"), v, idx-1)
	}

	for key, val := range m {
		if err := p.otsPutRowRecursive(otsCli, table, primaryKeys, append(values, key), val, idx-1); err != nil {
			return err
		}
	}

	return nil
}
