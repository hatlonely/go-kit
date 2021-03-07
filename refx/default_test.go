package refx

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReflectCopy(t *testing.T) {
	Convey("TestReflectCopy", t, func() {
		type A struct {
			Key1 string
			Key2 int
		}

		a1 := A{
			Key1: "val1",
			Key2: 2,
		}
		var a2 A

		reflect.ValueOf(&a2).Elem().Set(reflect.ValueOf(a1))
		So(a1.Key1, ShouldEqual, "val1")
		So(a1.Key2, ShouldEqual, 2)
		So(a2.Key1, ShouldEqual, "val1")
		So(a2.Key2, ShouldEqual, 2)
		a1.Key2 = 3
		So(a1.Key1, ShouldEqual, "val1")
		So(a1.Key2, ShouldEqual, 3)
		So(a2.Key1, ShouldEqual, "val1")
		So(a2.Key2, ShouldEqual, 2)
	})
}

func TestSetDefaultValueP(t *testing.T) {
	Convey("TestSetDefaultValueP", t, func() {
		type A struct {
			Key1 string `dft:"val1"`
			Key2 int    `dft:"2"`
		}

		var a A
		SetDefaultValueCopyP(&a)

		So(a.Key1, ShouldEqual, "val1")
		So(a.Key2, ShouldEqual, 2)
		a.Key1 = "hello"
		So(a.Key1, ShouldEqual, "hello")

		var b A
		SetDefaultValueCopyP(&b)
		So(b.Key1, ShouldEqual, "val1")
		So(b.Key2, ShouldEqual, 2)
	})
}

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
			So(SetDefaultValueCopy(&v), ShouldBeNil)

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
			So(SetDefaultValueCopy(v), ShouldBeNil)
		})

		Convey("case non point", func() {
			var v B
			So(SetDefaultValue(v), ShouldNotBeNil)
			So(func() {
				SetDefaultValueCopyP(v)
			}, ShouldPanic)
		})
	})
}

func TestDeepCopy(t *testing.T) {
	Convey("TestDeepCopy", t, func() {
		type A struct {
			Key1 *regexp.Regexp `dft:"val1"`
			Key2 *regexp.Regexp `dft:"val2"`
		}

		var a1 A
		SetDefaultValueP(&a1)
		So(a1.Key1.String(), ShouldEqual, "val1")
		So(a1.Key2.String(), ShouldEqual, "val2")

		a1.Key1 = regexp.MustCompile("val3")
		a1.Key2 = regexp.MustCompile("val4")

		var a2 A
		SetDefaultValueP(&a2)
		So(a2.Key1.String(), ShouldEqual, "val1")
		So(a2.Key2.String(), ShouldEqual, "val2")
	})
}

func BenchmarkSetDefaultValue(b *testing.B) {
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

	b.Run("set default with cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var v B
			_ = SetDefaultValueCopy(&v)
		}
	})

	b.Run("set default 1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var v B
			_ = SetDefaultValue(&v)
		}
	})
}

func BenchmarkSetDefaultValue2(b *testing.B) {
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

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var v B
			_ = SetDefaultValueCopy(&v)
		}
	})
}
