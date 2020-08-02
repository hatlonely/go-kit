package flag_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/cli"
	"github.com/hatlonely/go-kit/flag"
)

type Options struct {
	Help    bool   `flag:"--help,-h; default: false; usage: show help info"`
	Version bool   `flag:"--version,-v; default: false; usage: show version"`
	Port    string `bind:"port" flag:"--port; usage: service port"`

	Redis cli.RedisOptions `bind:"redis"`
	Mysql cli.MySQLOptions `bind:"mysql"`
}

func TestFlag(t *testing.T) {
	Convey("TestFlag", t, func() {
		options := &Options{}
		flag.Struct(options)
		fmt.Println(flag.Usage())
	})
}
