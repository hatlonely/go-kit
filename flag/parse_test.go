package flag

import (
	"net"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFlag_ParseArgs(t *testing.T) {
	Convey("TestFlag_ParseArgs", t, func() {
		flag := NewFlag("test")

		x := flag.Bool("x", false, "x")
		y := flag.Bool("y", false, "y")
		z := flag.Bool("z", false, "z")
		u := flag.Bool("u", false, "u")
		v := flag.Bool("v", false, "v")
		i1 := flag.Int("i1", 1, "i1")
		i2 := flag.Int("i2", 2, "i2")
		i3 := flag.Int("i3", 3, "i3")
		i4 := flag.Int("i4", 4, "i4")
		i5 := flag.Int("i5", 5, "i5")
		s1 := flag.String("s1", "1", "s1")
		s2 := flag.String("s2", "2", "s2")
		s3 := flag.String("s3", "3", "s3")
		s4 := flag.String("s4", "4", "s4")
		s5 := flag.String("s5", "5", "s5")
		p := flag.String("p", "0", "p")
		q := flag.Int("q", 0, "q")
		vi := flag.IntSlice("int-slice", []int{1, 2, 3}, "vi")
		ip := flag.IP("ip", nil, "ip flag")
		ti := flag.Time("time", time.Now(), "time flag")
		flag.AddArgument("p1", "p1")
		if err := flag.ParseArgs(strings.Split("-xy -z -u true -v=true pos1 "+
			"-i1 111 -i2=222 --i3 333 --i4=444 "+
			"-s1 111 -s2=222 --s3 333 --s4=444 "+
			"-p12345 pos2 -q67890 "+
			"-ip=192.168.0.1 --int-slice 4,5,6", " ")); err != nil {
			panic(err)
		}
		So(*x, ShouldBeTrue)
		So(*y, ShouldBeTrue)
		So(*z, ShouldBeTrue)
		So(*u, ShouldBeTrue)
		So(*v, ShouldBeTrue)
		So(*i1, ShouldEqual, 111)
		So(*i2, ShouldEqual, 222)
		So(*i3, ShouldEqual, 333)
		So(*i4, ShouldEqual, 444)
		So(*i5, ShouldEqual, 5)
		So(*s1, ShouldEqual, "111")
		So(*s2, ShouldEqual, "222")
		So(*s3, ShouldEqual, "333")
		So(*s4, ShouldEqual, "444")
		So(*s5, ShouldEqual, "5")
		So(*p, ShouldEqual, "12345")
		So(*q, ShouldEqual, 67890)
		So(*vi, ShouldResemble, []int{4, 5, 6})
		So(*ip, ShouldResemble, net.ParseIP("192.168.0.1"))
		So(time.Now().Sub(*ti), ShouldBeLessThan, 1*time.Second)
		So(flag.GetString("p1"), ShouldEqual, "pos1")
	})
}
