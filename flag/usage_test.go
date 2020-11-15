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
