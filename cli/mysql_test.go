package cli

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

func TestNewMysqlWithConfig(t *testing.T) {
	Convey("TestNewMysqlWithConfig", t, func() {
		So(ioutil.WriteFile("1.json", []byte(`{
		  "username": "root",
		  "password": "",
		  "database": "account",
		  "host": "127.0.0.1",
		  "port": 3306,
		  "connMaxLifeTime": "40s",
		  "maxIdleConns": 12,
		  "maxOpenConns": 23
		}`), 0644), ShouldBeNil)
		cfg, err := config.NewConfigWithSimpleFile("1.json")
		So(err, ShouldBeNil)
		client, err := NewMysqlWithConfig(cfg, refx.WithCamelName())
		So(err, ShouldBeNil)
		So(client, ShouldNotBeNil)
		So(os.RemoveAll("1.json"), ShouldBeNil)
	})
}
