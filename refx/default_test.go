package refx

import (
	"fmt"
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetDefaultValue(t *testing.T) {
	Convey("TestSetDefaultValue", t, func() {
		type A struct {
			Key1 string `dft:"val1"`
			Key2 int    `dft:"2"`
		}

		type B struct {
			A
			Key3 string `dft:"val3"`
			Key4 int    `dft:"4"`
			Key5 struct {
				Key6 string `dft:"val6"`
				Key7 struct {
					Key8 int `dft:"8"`
				}
			}
			Key9  time.Time     `dft:"2020-11-15T23:38:35+0800"`
			Key10 time.Duration `dft:"4s"`
			Key11 net.IP        `dft:"192.168.0.1"`
			Key12 float64       `dft:"6.4"`
			Key13 []int         `dft:"1,2,3"`
			A1    A
			A2    *A
			A3    *A
			a4    A
			A5    []A
		}

		Convey("case normal", func() {
			var v B
			v.A2 = &A{}
			So(SetDefaultValue(&v), ShouldBeNil)

			So(v.Key1, ShouldEqual, "val1")
			So(v.Key2, ShouldEqual, 2)
			So(v.Key3, ShouldEqual, "val3")
			So(v.Key4, ShouldEqual, 4)
			So(v.Key5.Key6, ShouldEqual, "val6")
			So(v.Key5.Key7.Key8, ShouldEqual, 8)
			So(v.A1.Key1, ShouldEqual, "val1")
			So(v.A1.Key2, ShouldEqual, 2)
			So(v.A3, ShouldBeNil)
			So(v.Key9, ShouldEqual, time.Unix(1605454715, 0))
			So(v.Key10, ShouldEqual, 4*time.Second)
			So(v.Key11, ShouldEqual, net.ParseIP("192.168.0.1"))
			So(v.Key12, ShouldAlmostEqual, 6.4)
			So(v.Key13, ShouldResemble, []int{1, 2, 3})
		})

		Convey("case nil", func() {
			var v *B
			fmt.Println(v)
			So(SetDefaultValue(v), ShouldBeNil)
		})
	})
}
