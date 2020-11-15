package flag

import (
	"net"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/refx"
)

func TestBind(t *testing.T) {
	Convey("TestBind", t, func() {
		type A struct {
			Key1 string `flag:"usage: key1 usage"`
			Key2 int    `flag:"usage: key2 usage"`
		}

		type B struct {
			A
			Key3 struct {
				Key4 string `flag:"--action, -a; usage: key4 usage; default: hello"`
				Key5 int    `flag:"-o, --operation; usage: key5 usage"`
				Key6 struct {
					Key7 string `flag:"usage: key6 usage; required"`
				}
			}
			Key8 *struct {
				Key9 string `flag:"args; usage: key8 usage; required"`
			}
			Key10 struct {
				Key11 string `flag:"isArgument"`
				Key12 int
			}
		}

		flag := NewFlag("test")
		So(flag.Struct(&B{}, refx.WithCamelName()), ShouldBeNil)

		So(flag.shorthandKeyMap, ShouldResemble, map[string]string{
			"a": "key3.key4",
			"o": "key3.key5",
		})
		So(flag.nameKeyMap, ShouldResemble, map[string]string{
			"key1":           "key1",
			"action":         "key3.key4",
			"operation":      "key3.key5",
			"key3.key6.key7": "key3.key6.key7",
			"args":           "key8.key9",
			"key10.key11":    "key10.key11",
			"key10.key12":    "key10.key12",
			"key2":           "key2",
		})
		So(flag.arguments, ShouldResemble, []string{"args", "key10.key11"})
	})
}

func TestBind2(t *testing.T) {
	Convey("TestBind", t, func() {
		type Sub struct {
			F1 int
			F2 string
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
			Sub3     Sub
			Sub
			ignoreF string
		}

		var options Options
		flag := NewFlag("test")
		if err := flag.Struct(&options); err != nil {
			panic(err)
		}

		if err := flag.ParseArgs(strings.Split("-i 456 -str abc -ip 192.168.0.1 --int-slice 1,2,3 posflag "+
			"-Sub1.F1 30 -F1 40 --Sub2.F2 hello "+
			"--Sub3.F1 50 -Sub3.F2=world "+
			"--not-exist-key val", " ")); err != nil {
			panic(err)
		}

		So(flag.shorthandKeyMap, ShouldResemble, map[string]string{
			"i": "I",
			"s": "S",
		})
		So(flag.nameKeyMap, ShouldResemble, map[string]string{
			"ip":        "IP",
			"time":      "Time",
			"Sub2.F1":   "Sub2.F1",
			"Sub3.F1":   "Sub3.F1",
			"F1":        "F1",
			"str":       "S",
			"pos":       "Pos",
			"Sub1.F2":   "Sub1.F2",
			"Sub2.F2":   "Sub2.F2",
			"Sub3.F2":   "Sub3.F2",
			"F2":        "F2",
			"int":       "I",
			"int-slice": "IntSlice",
			"Sub1.F1":   "Sub1.F1",
		})

		So(flag.root, ShouldResemble, map[string]interface{}{
			"Time": "2020-11-14T23:46:14+08:00",
			"S":    "abc",
			"F1":   "40",
			"Sub3": map[string]interface{}{
				"F1": "50",
				"F2": "world",
			},
			"not-exist-key": "val",
			"I":             "456",
			"IntSlice":      "1,2,3",
			"IP":            "192.168.0.1",
			"Sub1": map[string]interface{}{
				"F1": "30",
			},
			"Sub2": map[string]interface{}{
				"F2": "hello",
			},
			"Pos": "posflag",
		})

		So(options.I, ShouldEqual, 456)
		So(options.S, ShouldEqual, "abc")
		So(options.IP, ShouldResemble, net.ParseIP("192.168.0.1"))
		So(options.IntSlice, ShouldResemble, []int{1, 2, 3})
		So(options.Time, ShouldEqual, time.Unix(1605368774, 0))
		So(options.Pos, ShouldEqual, "posflag")
		So(options.Sub1.F1, ShouldEqual, 30)
		So(options.Sub2.F2, ShouldEqual, "hello")
		So(options.F1, ShouldEqual, 40)
		So(options.Sub3.F1, ShouldEqual, 50)
		So(options.Sub3.F2, ShouldEqual, "world")

		So(flag.GetString("not-exist-key"), ShouldEqual, "val")
	})

}
