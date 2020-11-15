package flag

import (
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFlag_Usage(t *testing.T) {
	Convey("TestFlag_Usage", t, func() {
		flag := NewFlag("test")
		flag.AddFlag("bool-val", "bool val flag", Type(true))
		flag.AddFlag("int-val", "int val flag", DefaultValue("20"), Type(0))
		flag.AddFlag("int8-val", "int8 val flag", Type(int8(0)), Shorthand("i"))
		flag.AddFlag("int16-val", "int16 val flag", Type(int16(0)))
		flag.AddFlag("time-val", "time val flag", Type(time.Time{}))
		flag.AddFlag("duration-val", "", Type(time.Duration(0)))
		flag.AddFlag("ip-val", "", Type(net.IP{}))
		flag.AddArgument("p1", "pos flag", DefaultValue("123"))
		flag.AddArgument("p2", "pos flag")

		So(flag.Usage(), ShouldEqual, `usage: test [p1] [p2] [-bool-val bool] [-duration-val time.Duration] [-int-val int=20] [-int16-val int16] [-i,int8-val int8] [-ip-val net.IP] [-time-val time.Time]

arguments:
      p1              [string=123]     pos flag
      p2              [string]         pos flag

options:
    , --bool-val      [bool]           bool val flag
    , --duration-val  [time.Duration]  
    , --int-val       [int=20]         int val flag
    , --int16-val     [int16]          int16 val flag
  -i, --int8-val      [int8]           int8 val flag
    , --ip-val        [net.IP]         
    , --time-val      [time.Time]      time val flag
`)
	})
}

func TestFlag_Usage2(t *testing.T) {
	Convey("TestFlag_Usage2", t, func() {
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
		}

		flag := NewFlag("test")
		So(flag.Struct(&B{}), ShouldBeNil)

		So(flag.Usage(), ShouldEqual, `usage: test [Key8.Key9] [-key1 string] [-key2 int] [-a,key3-action string=hello] [-o,key3-operation int] <-key3-key6-key7 string>

arguments:
      key8-args                           [string]        key8 usage

options:
    , --key1, --Key1                      [string]        key1 usage
    , --key2, --Key2                      [int]           key2 usage
  -a, --key3-action, --Key3.Key4          [string=hello]  key4 usage
  -o, --key3-operation, --Key3.Key5       [int]           key5 usage
    , --key3-key6-key7, --Key3.Key6.Key7  [string]        key6 usage
`)
	})
}
