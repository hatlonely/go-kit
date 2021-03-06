package config

import (
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/refx"
)

func CreateTestFile4() {
	fp, _ := os.Create("test.json")
	_, _ = fp.WriteString(`{
  "mysql": {
    "username": "root",
    "password": "123456",
    "database": "account",
    "host": "127.0.0.1",
    "port": 3306,
    "connMaxLifeTime": "60s",
    "maxIdleConns": 10,
    "maxOpenConns": 20
  },
}`)
	_ = fp.Close()
}

func CreateTestFile5() {
	fp, _ := os.Create("test.json")
	_, _ = fp.WriteString(`{
  "mysql": {
    "username": "hatlonely",
    "password": "654321",
    "database": "account2",
    "host": "11.46.23.89",
    "port": 3306,
    "connMaxLifeTime": "120s",
    "maxIdleConns": 15,
    "maxOpenConns": 25
  },
}`)
	_ = fp.Close()
}

func TestBindVar(t *testing.T) {
	Convey("TestBindVar", t, func() {
		CreateTestFile4()

		cfg, err := NewConfigWithSimpleFile("test.json")
		So(err, ShouldBeNil)
		So(cfg.Watch(), ShouldBeNil)
		defer cfg.Stop()

		type Options struct {
			Username        string `dft:"root"`
			Password        string
			Database        string
			Host            string        `dft:"localhost"`
			Port            int           `dft:"3306"`
			ConnMaxLifeTime time.Duration `dft:"60s"`
			MaxIdleConns    int           `dft:"10"`
			MaxOpenConns    int           `dft:"20"`
			LogMode         bool
		}

		var bindVal atomic.Value

		cfg.BindVar("mysql", Options{}, &bindVal, OnSucc(func(c *Config) {
			fmt.Println(c.ToString())
		}), WithUnmarshalOptions(refx.WithCamelName()))

		So(bindVal.Load().(Options).Username, ShouldEqual, "root")
		So(bindVal.Load().(Options).Password, ShouldEqual, "123456")
		So(bindVal.Load().(Options).Database, ShouldEqual, "account")
		So(bindVal.Load().(Options).Host, ShouldEqual, "127.0.0.1")
		So(bindVal.Load().(Options).Port, ShouldEqual, 3306)
		So(bindVal.Load().(Options).ConnMaxLifeTime, ShouldEqual, 60*time.Second)
		So(bindVal.Load().(Options).MaxIdleConns, ShouldEqual, 10)
		So(bindVal.Load().(Options).MaxOpenConns, ShouldEqual, 20)

		CreateTestFile5()
		time.Sleep(100 * time.Millisecond)

		So(bindVal.Load().(Options).Username, ShouldEqual, "hatlonely")
		So(bindVal.Load().(Options).Password, ShouldEqual, "654321")
		So(bindVal.Load().(Options).Database, ShouldEqual, "account2")
		So(bindVal.Load().(Options).Host, ShouldEqual, "11.46.23.89")
		So(bindVal.Load().(Options).Port, ShouldEqual, 3306)
		So(bindVal.Load().(Options).ConnMaxLifeTime, ShouldEqual, 120*time.Second)
		So(bindVal.Load().(Options).MaxIdleConns, ShouldEqual, 15)
		So(bindVal.Load().(Options).MaxOpenConns, ShouldEqual, 25)

		DeleteTestFile()
	})
}

func TestBind(t *testing.T) {
	Convey("TestBind", t, func() {
		CreateTestFile4()

		cfg, err := NewConfigWithSimpleFile("test.json")
		So(err, ShouldBeNil)
		So(cfg.Watch(), ShouldBeNil)
		defer cfg.Stop()

		type Options struct {
			Username        string `dft:"root"`
			Password        string
			Database        string
			Host            string        `dft:"localhost"`
			Port            int           `dft:"3306"`
			ConnMaxLifeTime time.Duration `dft:"60s"`
			MaxIdleConns    int           `dft:"10"`
			MaxOpenConns    int           `dft:"20"`
			LogMode         bool
		}

		bindVal := cfg.Bind("mysql", Options{}, OnSucc(func(c *Config) {
			fmt.Println(c.ToString())
		}), WithUnmarshalOptions(refx.WithCamelName()))

		So(bindVal.Load().(Options).Username, ShouldEqual, "root")
		So(bindVal.Load().(Options).Password, ShouldEqual, "123456")
		So(bindVal.Load().(Options).Database, ShouldEqual, "account")
		So(bindVal.Load().(Options).Host, ShouldEqual, "127.0.0.1")
		So(bindVal.Load().(Options).Port, ShouldEqual, 3306)
		So(bindVal.Load().(Options).ConnMaxLifeTime, ShouldEqual, 60*time.Second)
		So(bindVal.Load().(Options).MaxIdleConns, ShouldEqual, 10)
		So(bindVal.Load().(Options).MaxOpenConns, ShouldEqual, 20)

		CreateTestFile5()
		time.Sleep(100 * time.Millisecond)

		So(bindVal.Load().(Options).Username, ShouldEqual, "hatlonely")
		So(bindVal.Load().(Options).Password, ShouldEqual, "654321")
		So(bindVal.Load().(Options).Database, ShouldEqual, "account2")
		So(bindVal.Load().(Options).Host, ShouldEqual, "11.46.23.89")
		So(bindVal.Load().(Options).Port, ShouldEqual, 3306)
		So(bindVal.Load().(Options).ConnMaxLifeTime, ShouldEqual, 120*time.Second)
		So(bindVal.Load().(Options).MaxIdleConns, ShouldEqual, 15)
		So(bindVal.Load().(Options).MaxOpenConns, ShouldEqual, 25)

		DeleteTestFile()
	})
}
