package flag

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

type MySubFlags struct {
	F1 int    `flag:"--f1; default:20; usage:f1 flag"`
	F2 string `flag:"--f2; default:hatlonely; usage:f2 flag"`
}

type MyFlags struct {
	I        int        `flag:"--int, -i; required; default: 123; usage: int flag"`
	S        string     `flag:"--str, -s; required; usage: str flag"`
	IntSlice []int      `flag:"--int-slice; default: 1,2,3; usage: int slice flag"`
	IP       net.IP     `flag:"--ip; usage: ip flag"`
	Time     time.Time  `flag:"--time; usage: time flag; default: 2019-11-27"`
	Pos      string     `flag:"pos; usage: pos flag"`
	Sub      MySubFlags `flag:"sub"`
}

func TestBind(t *testing.T) {
	mf := &MyFlags{}
	f := NewFlag("hello")
	if err := f.Struct(mf); err != nil {
		panic(err)
	}
	fmt.Println(f.Usage())
	if err := f.ParseArgs(strings.Split("-str abc -ip 192.168.0.1 --int-slice 1,2,3 posflag -sub-f1 140", " ")); err != nil {
		fmt.Println(f.Usage())
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
