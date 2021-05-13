package config

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey"
	"github.com/aliyun/aliyun-tablestore-go-sdk/v5/tablestore"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOTSLegacyProvider(t *testing.T) {
	Convey("TestOTSProvider", t, func(c C) {
		patches := ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "GetRange",
			func(client *tablestore.TableStoreClient, req *tablestore.GetRangeRequest) (*tablestore.GetRangeResponse, error) {
				c.So(req.RangeRowQueryCriteria.TableName, ShouldEqual, "TestConfig")
				fmt.Println(req.RangeRowQueryCriteria.TableName)
				skvs := map[string]interface{}{}
				for _, pk := range req.RangeRowQueryCriteria.StartPrimaryKey.PrimaryKeys {
					skvs[pk.ColumnName] = pk.Value
				}
				fmt.Println(skvs)
				ekvs := map[string]interface{}{}
				for _, pk := range req.RangeRowQueryCriteria.EndPrimaryKey.PrimaryKeys {
					ekvs[pk.ColumnName] = pk.Value
				}
				fmt.Println(ekvs)

				return &tablestore.GetRangeResponse{
					Rows: []*tablestore.Row{
						{
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key1"},
									{ColumnName: "PK2", Value: "Key2"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "_", Value: "val2", Timestamp: 1604980294775},
							},
						}, {
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key1"},
									{ColumnName: "PK2", Value: "Key3"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "Key4", Value: "val4", Timestamp: 1604980294803},
								{ColumnName: "Key5", Value: "val5", Timestamp: 1604980294803},
							},
						}, {
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key1"},
									{ColumnName: "PK2", Value: "Key6"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "Key7", Value: "{\"Key8\":\"val8\",\"Key9\":9}", Timestamp: 1604980294830},
							},
						}, {
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key10"},
									{ColumnName: "PK2", Value: "_"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "_", Value: 10, Timestamp: 1604980294666},
							},
						}, {
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key11"},
									{ColumnName: "PK2", Value: "_"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "_", Value: "val11", Timestamp: 1604980294694},
							},
						}, {
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key12"},
									{ColumnName: "PK2", Value: "_"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "_", Value: "[{\"Key13\":\"val13\",\"Key14\":{\"Key15\":\"val15\"}}]", Timestamp: 1604980294721},
							},
						}, {
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key16"},
									{ColumnName: "PK2", Value: "_"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "_", Value: "val16", Timestamp: 1604980294721},
							},
						}, {
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key17"},
									{ColumnName: "PK2", Value: "_"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "_", Value: 1.23, Timestamp: 1604980294721},
							},
						},
					},
				}, nil
			},
		).ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "DescribeTable", func(
			cli *tablestore.TableStoreClient, req *tablestore.DescribeTableRequest) (*tablestore.DescribeTableResponse, error) {
			So(req.TableName, ShouldEqual, "TestConfig")
			res := &tablestore.DescribeTableResponse{
				TableMeta: &tablestore.TableMeta{},
			}
			res.TableMeta.AddPrimaryKeyColumn("PK1", tablestore.PrimaryKeyType_STRING)
			res.TableMeta.AddPrimaryKeyColumn("PK2", tablestore.PrimaryKeyType_STRING)
			return res, nil
		})
		defer patches.Reset()

		provider, err := NewOTSLegacyProviderWithOptions(&OTSLegacyProviderOptions{
			Endpoint:        "https://hatlonely.cn-shanghai.ots.aliyuncs.com",
			AccessKeyID:     "xx",
			AccessKeySecret: "xx",
			Instance:        "hatlonely",
			Table:           "TestConfig",
			PrimaryKeys:     []string{"PK1", "PK2"},
			Interval:        200 * time.Millisecond,
		})
		So(err, ShouldBeNil)

		ctx, cancel := context.WithCancel(context.Background())
		So(provider.EventLoop(ctx), ShouldBeNil)
		defer cancel()

		for len(provider.Events()) != 0 {
			<-provider.Events()
		}

		time.Sleep(400 * time.Millisecond)
	})
}

func TestNewOTSLegacyProviderWithOptions(t *testing.T) {
	Convey("TestNewOTSLegacyProviderWithOptions", t, func() {
		patches := ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "DescribeTable", func(
			cli *tablestore.TableStoreClient, req *tablestore.DescribeTableRequest) (*tablestore.DescribeTableResponse, error) {
			So(req.TableName, ShouldEqual, "TestConfig")
			res := &tablestore.DescribeTableResponse{
				TableMeta: &tablestore.TableMeta{},
			}
			res.TableMeta.AddPrimaryKeyColumn("PK1", tablestore.PrimaryKeyType_STRING)
			res.TableMeta.AddPrimaryKeyColumn("PK2", tablestore.PrimaryKeyType_STRING)
			res.TableMeta.AddPrimaryKeyColumn("PK3", tablestore.PrimaryKeyType_STRING)
			return res, nil
		}).ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "GetRange",
			func(client *tablestore.TableStoreClient, req *tablestore.GetRangeRequest) (*tablestore.GetRangeResponse, error) {
				return &tablestore.GetRangeResponse{}, nil
			},
		)
		defer patches.Reset()

		_, err := NewOTSLegacyProviderWithOptions(&OTSLegacyProviderOptions{
			Endpoint:        "https://hatlonely.cn-shanghai.ots.aliyuncs.com",
			AccessKeyID:     "xx",
			AccessKeySecret: "xx",
			Instance:        "hatlonely",
			Table:           "TestConfig",
			PrimaryKeys:     []string{"PK1", "PK2", "PK3"},
			Interval:        200 * time.Millisecond,
		})
		So(err, ShouldBeNil)
	})

	Convey("TestNewOTSLegacyProviderWithOptions CreateTable", t, func() {
		patches := ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "DescribeTable", func(
			cli *tablestore.TableStoreClient, req *tablestore.DescribeTableRequest) (*tablestore.DescribeTableResponse, error) {
			So(req.TableName, ShouldEqual, "TestConfig")
			return nil, errors.New("does not exist")
		}).ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "GetRange", func(
			client *tablestore.TableStoreClient, req *tablestore.GetRangeRequest) (*tablestore.GetRangeResponse, error) {
			return &tablestore.GetRangeResponse{}, nil
		}).ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "CreateTable", func(
			cli *tablestore.TableStoreClient, req *tablestore.CreateTableRequest) (*tablestore.CreateTableResponse, error) {
			So(req.TableMeta.TableName, ShouldEqual, "TestConfig")
			So(len(req.TableMeta.SchemaEntry), ShouldEqual, 3)
			So(*req.TableMeta.SchemaEntry[0].Name, ShouldEqual, "PK1")
			So(*req.TableMeta.SchemaEntry[0].Type, ShouldEqual, tablestore.PrimaryKeyType_STRING)
			So(*req.TableMeta.SchemaEntry[1].Name, ShouldEqual, "PK2")
			So(*req.TableMeta.SchemaEntry[1].Type, ShouldEqual, tablestore.PrimaryKeyType_STRING)
			So(*req.TableMeta.SchemaEntry[2].Name, ShouldEqual, "PK3")
			So(*req.TableMeta.SchemaEntry[2].Type, ShouldEqual, tablestore.PrimaryKeyType_STRING)
			return nil, nil
		})
		defer patches.Reset()

		_, err := NewOTSLegacyProviderWithOptions(&OTSLegacyProviderOptions{
			Endpoint:        "https://hatlonely.cn-shanghai.ots.aliyuncs.com",
			AccessKeyID:     "xx",
			AccessKeySecret: "xx",
			Instance:        "hatlonely",
			Table:           "TestConfig",
			PrimaryKeys:     []string{"PK1", "PK2", "PK3"},
			Interval:        200 * time.Millisecond,
		})
		So(err, ShouldBeNil)
	})
}

func TestOTSLegacyProvider_Load(t *testing.T) {
	Convey("TestOTSLegacyProvider_Load", t, func() {
		patches := ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "GetRange",
			func(client *tablestore.TableStoreClient, req *tablestore.GetRangeRequest) (*tablestore.GetRangeResponse, error) {
				So(req.RangeRowQueryCriteria.TableName, ShouldEqual, "TestConfig")
				fmt.Println(req.RangeRowQueryCriteria.TableName)
				skvs := map[string]interface{}{}
				for _, pk := range req.RangeRowQueryCriteria.StartPrimaryKey.PrimaryKeys {
					skvs[pk.ColumnName] = pk.Value
				}
				fmt.Println(skvs)
				ekvs := map[string]interface{}{}
				for _, pk := range req.RangeRowQueryCriteria.EndPrimaryKey.PrimaryKeys {
					ekvs[pk.ColumnName] = pk.Value
				}
				fmt.Println(ekvs)

				return &tablestore.GetRangeResponse{
					Rows: []*tablestore.Row{
						{
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key1"},
									{ColumnName: "PK2", Value: "Key2"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "_", Value: "val2", Timestamp: 1604980294775},
							},
						}, {
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key1"},
									{ColumnName: "PK2", Value: "Key3"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "Key4", Value: "val4", Timestamp: 1604980294803},
								{ColumnName: "Key5", Value: "val5", Timestamp: 1604980294803},
							},
						}, {
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key1"},
									{ColumnName: "PK2", Value: "Key6"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "Key7", Value: "{\"Key8\":\"val8\",\"Key9\":9}", Timestamp: 1604980294830},
							},
						}, {
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key10"},
									{ColumnName: "PK2", Value: "_"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "_", Value: 10, Timestamp: 1604980294666},
							},
						}, {
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key11"},
									{ColumnName: "PK2", Value: "_"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "_", Value: "val11", Timestamp: 1604980294694},
							},
						}, {
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key12"},
									{ColumnName: "PK2", Value: "_"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "_", Value: "[{\"Key13\":\"val13\",\"Key14\":{\"Key15\":\"val15\"}}]", Timestamp: 1604980294721},
							},
						}, {
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key16"},
									{ColumnName: "PK2", Value: "_"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "_", Value: "val16", Timestamp: 1604980294721},
							},
						}, {
							PrimaryKey: &tablestore.PrimaryKey{
								PrimaryKeys: []*tablestore.PrimaryKeyColumn{
									{ColumnName: "PK1", Value: "Key17"},
									{ColumnName: "PK2", Value: "_"},
								},
							},
							Columns: []*tablestore.AttributeColumn{
								{ColumnName: "_", Value: 1.23, Timestamp: 1604980294721},
							},
						},
					},
				}, nil
			},
		).ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "DescribeTable", func(
			cli *tablestore.TableStoreClient, req *tablestore.DescribeTableRequest) (*tablestore.DescribeTableResponse, error) {
			So(req.TableName, ShouldEqual, "TestConfig")
			res := &tablestore.DescribeTableResponse{
				TableMeta: &tablestore.TableMeta{},
			}
			res.TableMeta.AddPrimaryKeyColumn("PK1", tablestore.PrimaryKeyType_STRING)
			res.TableMeta.AddPrimaryKeyColumn("PK2", tablestore.PrimaryKeyType_STRING)
			return res, nil
		})
		defer patches.Reset()

		provider, err := NewOTSLegacyProviderWithOptions(&OTSLegacyProviderOptions{
			Endpoint:        "https://hatlonely.cn-shanghai.ots.aliyuncs.com",
			AccessKeyID:     "xx",
			AccessKeySecret: "xx",
			Instance:        "hatlonely",
			Table:           "TestConfig",
			PrimaryKeys:     []string{"PK1", "PK2"},
			Interval:        200 * time.Millisecond,
		})
		So(err, ShouldBeNil)
		buf, err := provider.Load()
		So(err, ShouldBeNil)
		So(string(buf), ShouldEqual, "{\"Key1\":{\"Key2\":\"val2\",\"Key3\":{\"Key4\":\"val4\",\"Key5\":\"val5\"},\"Key6\":{\"Key7\":{\"Key8\":\"val8\",\"Key9\":9}}},\"Key10\":10,\"Key11\":\"val11\",\"Key12\":[{\"Key13\":\"val13\",\"Key14\":{\"Key15\":\"val15\"}}],\"Key16\":\"val16\",\"Key17\":1.23}")
	})
}

func TestOTSLegacyProvider_Dump(t *testing.T) {
	Convey("TestOTSLegacyProvider_Dump", t, func() {
		patches := ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "PutRow", func(
			client *tablestore.TableStoreClient, req *tablestore.PutRowRequest) (*tablestore.PutRowResponse, error) {
			So(req.PutRowChange.TableName, ShouldEqual, "TestConfig")
			kvs := map[string]interface{}{}
			for _, pk := range req.PutRowChange.PrimaryKey.PrimaryKeys {
				kvs[pk.ColumnName] = pk.Value
			}
			for _, col := range req.PutRowChange.Columns {
				kvs[col.ColumnName] = col.Value
			}
			So(kvs, ShouldBeIn, []map[string]interface{}{
				{"PK1": "Key10", "PK2": "_", "_": "10"},
				{"PK1": "Key11", "PK2": "_", "_": "val11"},
				{"PK1": "Key12", "PK2": "_", "_": "[{\"Key13\":\"val13\",\"Key14\":{\"Key15\":\"val15\"}}]"},
				{"PK1": "Key16", "PK2": "_", "_": "val16"},
				{"PK1": "Key17", "PK2": "_", "_": "1.23"},
				{"Key4": "val4", "Key5": "val5", "PK1": "Key1", "PK2": "Key3"},
				{"Key7": "{\"Key8\":\"val8\",\"Key9\":9}", "PK1": "Key1", "PK2": "Key6"},
				{"PK1": "Key1", "PK2": "Key2", "_": "val2"},
			})
			return nil, nil
		}).ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "DescribeTable", func(
			cli *tablestore.TableStoreClient, req *tablestore.DescribeTableRequest) (*tablestore.DescribeTableResponse, error) {
			So(req.TableName, ShouldEqual, "TestConfig")
			res := &tablestore.DescribeTableResponse{
				TableMeta: &tablestore.TableMeta{},
			}
			res.TableMeta.AddPrimaryKeyColumn("PK1", tablestore.PrimaryKeyType_STRING)
			res.TableMeta.AddPrimaryKeyColumn("PK2", tablestore.PrimaryKeyType_STRING)
			return res, nil
		}).ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "GetRange",
			func(client *tablestore.TableStoreClient, req *tablestore.GetRangeRequest) (*tablestore.GetRangeResponse, error) {
				return &tablestore.GetRangeResponse{}, nil
			},
		)
		defer patches.Reset()

		provider, err := NewOTSLegacyProviderWithOptions(&OTSLegacyProviderOptions{
			Endpoint:        "https://hatlonely.cn-shanghai.ots.aliyuncs.com",
			AccessKeyID:     "xx",
			AccessKeySecret: "xx",
			Instance:        "hatlonely",
			Table:           "TestConfig",
			PrimaryKeys:     []string{"PK1", "PK2"},
			Interval:        200 * time.Millisecond,
		})
		So(err, ShouldBeNil)
		So(provider.Dump([]byte(`{
  "Key1": {
    "Key2": "val2",
    "Key3": {
      "Key4": "val4",
      "Key5": "val5"
    },
    "Key6": {
      "Key7": {
        "Key8": "val8",
        "Key9": 9
      }
    }
  },
  "Key10": 10,
  "Key11": "val11",
  "Key12": [{
	"Key13": "val13",
    "Key14": {
      "Key15": "val15"
    }
  }],
  "Key16": "val16",
  "Key17": 1.23
}`)), ShouldBeNil)
	})
}
