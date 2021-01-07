package cli

import (
	"testing"

	. "github.com/agiledragon/gomonkey"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewMysqlWithConfig(t *testing.T) {
	patches := ApplyFunc(gorm.Open, func(dialect string, args ...interface{}) (db *gorm.DB, err error) {
		So(dialect, ShouldEqual, "mysql")
		So(args, ShouldResemble, []interface{}{"root:@tcp(127.0.0.1:3307)/hatlonely?charset=utf8&parseTime=True&loc=Local"})
		return &gorm.DB{}, nil
	})
	defer patches.Reset()

	Convey("TestNewMysqlWithConfig", t, func() {
		client, err := NewMysqlWithOptions(&MySQLOptions{
			Username: "root",
			Database: "hatlonely",
			Host:     "127.0.0.1",
			Port:     3307,
		})
		So(err, ShouldBeNil)
		So(client, ShouldNotBeNil)
	})
}
