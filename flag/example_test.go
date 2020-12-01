package flag_test

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/hatlonely/go-kit/flag"
	"github.com/hatlonely/go-kit/refx"
)

func TestExample1(t *testing.T) {
	type MySubFlags struct {
		F1 int    `flag:"--f1; default:20; usage:f1 flag"`
		F2 string `flag:"--f2; default:hatlonely; usage:f2 flag"`
	}

	type MyFlags struct {
		I        int       `flag:"--int, -i; required; default: 123; usage: int flag"`
		S        string    `flag:"--str, -s; required; usage: str flag"`
		IntSlice []int     `flag:"--int-slice; default: 1,2,3; usage: int slice flag"`
		IP       net.IP    `flag:"--ip; usage: ip flag"`
		Time     time.Time `flag:"--time; usage: time flag; default: 2019-11-27"`
		Pos      string    `flag:"pos; usage: pos flag"`
		Sub      MySubFlags
		Sub2     *MySubFlags
	}

	mf := &MyFlags{}
	flag1 := flag.NewFlag("test")
	if err := flag1.Struct(mf, refx.WithCamelName()); err != nil {
		panic(err)
	}
	if err := flag1.ParseArgs(strings.Split("-str abc -ip 192.168.0.1 --int-slice 4,5,6 posflag -sub-f1 140", " ")); err != nil {
		fmt.Println(flag1.Usage())
		panic(err)
	}

	fmt.Println("int =>", mf.I)
	fmt.Println("str =>", mf.S)
	fmt.Println("int-slice =>", mf.IntSlice)
	fmt.Println("ip =>", mf.IP)
	fmt.Println("time =>", mf.Time)
	fmt.Println("sub.f1 =>", mf.Sub.F1)
	fmt.Println("sub.f2 =>", mf.Sub.F2)
}

func TestExample2(t *testing.T) {
	flag2 := flag.NewFlag("test")

	i := flag2.Int("int", 123, "int flag")
	s := flag2.String("str", "", "str flag")
	vi := flag2.IntSlice("int-slice", []int{1, 2, 3}, "int slice flag")
	ip := flag2.IP("ip", nil, "ip flag")
	ti := flag2.Time("time", time.Now(), "time flag")
	if err := flag2.ParseArgs(strings.Split("-str abc -ip 192.168.0.1 --int-slice 4,5,6 posflag", " ")); err != nil {
		fmt.Println(flag2.Usage())
		panic(err)
	}

	fmt.Println("int =>", *i)
	fmt.Println("str =>", *s)
	fmt.Println("int-slice =>", *vi)
	fmt.Println("ip =>", *ip)
	fmt.Println("time =>", *ti)
}

func TestExample3(t *testing.T) {
	flag3 := flag.NewFlag("test")

	flag3.AddOption("int", "int flag", flag.Required(), flag.Shorthand("i"), flag.Type(0), flag.DefaultValue("123"))
	flag3.AddOption("str", "str flag", flag.Shorthand("s"), flag.Required())
	flag3.AddOption("int-slice", "int slice flag", flag.Type([]int{}), flag.DefaultValue("1,2,3"))
	flag3.AddOption("ip", "ip flag", flag.Type(net.IP{}))
	flag3.AddOption("time", "time flag", flag.Type(time.Time{}), flag.DefaultValue("2019-11-27"))
	flag3.AddArgument("pos", "pos flag")
	if err := flag3.ParseArgs(strings.Split("-str abc -ip 192.168.0.1 --int-slice 4,5,6 posflag", " ")); err != nil {
		panic(err)
	}

	fmt.Println("int =>", flag3.GetInt("int"))
	fmt.Println("str =>", flag3.GetString("str"))
	fmt.Println("int-slice =>", flag3.GetIntSlice("intSlice"))
	fmt.Println("int-slice =>", flag3.GetIntSlice("int.slice"))
	fmt.Println("ip =>", flag3.GetIP("ip"))
	fmt.Println("time =>", flag3.GetTime("time"))
	fmt.Println("pos =>", flag3.GetString("pos"))
}
