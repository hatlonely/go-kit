package flag

import (
	"net"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBind(t *testing.T) {
	Convey("TestBind", t, func() {
		type Sub struct {
			F1 int    `flag:"--f1; default:20; usage:f1 flag"`
			F2 string `flag:"--f2; default:hatlonely; usage:f2 flag"`
		}

		type Options struct {
			I        int       `flag:"--int, -i; required; default: 123; usage: int flag"`
			S        string    `flag:"--str, -s; required; usage: str flag"`
			IntSlice []int     `flag:"--int-slice; default: 1,2,3; usage: int slice flag"`
			IP       net.IP    `flag:"--ip; usage: ip flag"`
			Time     time.Time `flag:"--time; usage: time flag; default: 2020-11-14T23:46:14+08:00"`
			Pos      string    `flag:"pos; usage: pos flag"`
			Sub1     Sub
			Sub2     *Sub
			Sub
			ignoreF string
		}

		var options Options
		flag := NewFlag("hello")
		if err := flag.Struct(&options); err != nil {
			panic(err)
		}
		if err := flag.ParseArgs(strings.Split("-i 456 -str abc -ip 192.168.0.1 --int-slice 1,2,3 posflag "+
			"-sub1-f1 30 -f1 40 --sub2-f2 hello  --not-exist-key val", " ")); err != nil {
			panic(err)
		}

		So(options.I, ShouldEqual, 456)
		So(options.S, ShouldEqual, "abc")
		So(options.IP, ShouldResemble, net.ParseIP("192.168.0.1"))
		So(options.IntSlice, ShouldResemble, []int{1, 2, 3})
		So(options.Time, ShouldEqual, time.Unix(1605368774, 0))
		So(options.Pos, ShouldEqual, "posflag")
		So(options.Sub1.F1, ShouldEqual, 30)
		So(options.Sub1.F2, ShouldEqual, "hatlonely")
		So(options.Sub2.F1, ShouldEqual, 20)
		So(options.Sub2.F2, ShouldEqual, "hello")
		So(options.F1, ShouldEqual, 40)

		So(flag.GetString("not-exist-key"), ShouldEqual, "val")
	})

}
