package config

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOTSProvider(t *testing.T) {
	Convey("TestOTSProvider", t, func(c C) {
		patches := ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "PutRow", func(
			client *tablestore.TableStoreClient, req *tablestore.PutRowRequest) (*tablestore.PutRowResponse, error) {
			kvs := map[string]interface{}{}
			for _, pk := range req.PutRowChange.PrimaryKey.PrimaryKeys {
				kvs[pk.ColumnName] = pk.Value
			}
			fmt.Println(kvs)
			for _, col := range req.PutRowChange.Columns {
				kvs[col.ColumnName] = col.Value
			}
			fmt.Println(kvs)

			return nil, nil
		}).ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "DescribeTable", func(
			cli *tablestore.TableStoreClient, req *tablestore.DescribeTableRequest) (*tablestore.DescribeTableResponse, error) {
			So(req.TableName, ShouldEqual, "TestConfig")
			return nil, errors.New("does not exist")
		}).ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "CreateTable", func(
			cli *tablestore.TableStoreClient, req *tablestore.CreateTableRequest) (*tablestore.CreateTableResponse, error) {
			So(req.TableMeta.TableName, ShouldEqual, "TestConfig")
			So(*req.TableMeta.SchemaEntry[0].Name, ShouldEqual, "Key")
			So(*req.TableMeta.SchemaEntry[0].Type, ShouldEqual, tablestore.PrimaryKeyType_STRING)
			So(req.TableMeta.DefinedColumns[0].Name, ShouldEqual, "Val")
			So(req.TableMeta.DefinedColumns[0].ColumnType, ShouldEqual, tablestore.DefinedColumn_STRING)
			So(req.TableOption.MaxVersion, ShouldEqual, 1)
			return nil, nil
		}).ApplyMethod(reflect.TypeOf(&tablestore.TableStoreClient{}), "GetRow", func(
			client *tablestore.TableStoreClient, req *tablestore.GetRowRequest) (*tablestore.GetRowResponse, error) {
			So(req.SingleRowQueryCriteria.TableName, ShouldEqual, "TestConfig")
			kvs := map[string]interface{}{}
			for _, pk := range req.SingleRowQueryCriteria.PrimaryKey.PrimaryKeys {
				kvs[pk.ColumnName] = pk.Value
			}
			So(kvs, ShouldResemble, map[string]interface{}{"Key": "test"})

			return &tablestore.GetRowResponse{
				PrimaryKey: tablestore.PrimaryKey{
					PrimaryKeys: []*tablestore.PrimaryKeyColumn{
						{ColumnName: "Key", Value: "test"},
					},
				},
				Columns: []*tablestore.AttributeColumn{
					{ColumnName: "Val", Value: fmt.Sprintf("hello world")},
				},
			}, nil
		})
		defer patches.Reset()

		otsCli := tablestore.NewClient(
			"https://hatlonely.cn-shanghai.ots.aliyuncs.com",
			"hatlonely",
			"xx",
			"xx",
		)
		provider, err := NewOTSProvider(otsCli, "TestConfig", "test", 100*time.Millisecond)
		So(err, ShouldBeNil)
		{
			So(provider.Dump([]byte("hello world")), ShouldBeNil)
			buf, _, err := OTSGetRow(otsCli, "TestConfig", "test")
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, "hello world")
		}
		{
			ctx, cancel := context.WithCancel(context.Background())
			So(provider.EventLoop(ctx), ShouldBeNil)

			for len(provider.Events()) != 0 {
				<-provider.Events()
			}

			for i := 0; i < 5; i++ {
				So(OTSPutRow(otsCli, "TestConfig", "test", []byte(fmt.Sprintf("hello world %v", i))), ShouldBeNil)
				<-provider.Events()
				buf, err := provider.Load()
				So(err, ShouldBeNil)
				fmt.Println(string(buf))
				So(string(buf), ShouldEqual, fmt.Sprintf("hello world %v", i))
				time.Sleep(time.Second)
			}

			cancel()
		}
	})
}
