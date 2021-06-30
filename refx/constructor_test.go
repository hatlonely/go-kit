package refx

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewConstructorInfo(t *testing.T) {
	Convey("TestNewConstructorInfo", t, func() {
		type A interface{}
		type Options struct {
			Key1 string
		}
		constructor, err := NewConstructor(func(options *Options, opts ...Option) A {
			return "hello world"
		}, reflect.TypeOf((*A)(nil)).Elem())

		So(err, ShouldBeNil)
		So(constructor, ShouldNotBeNil)
	})
}

type B interface {
	Fun() int
}

type C1 struct{}
type C1Options struct{}

func NewC1WithOptions(options *C1Options, opts ...Option) (*C1, error) {
	return &C1{}, nil
}
func (c *C1) Fun() int {
	return 1
}

type C2 struct {
	Options *C2Options
}
type C2Options struct {
	Val int
}

func NewC2WithOptions(options *C2Options) (*C2, error) {
	return &C2{
		Options: options,
	}, nil
}
func (c *C2) Fun() int {
	return c.Options.Val
}

func TestNew(t *testing.T) {
	Convey("TestNew", t, func() {
		Register(reflect.TypeOf((*B)(nil)).Elem(), "C1", NewC1WithOptions)
		Register(reflect.TypeOf((*B)(nil)).Elem(), "C2", NewC2WithOptions)
		Register(reflect.TypeOf((*B)(nil)).Elem(), "C3", &C2{Options: &C2Options{Val: 3}})

		{
			v, err := New(reflect.TypeOf((*B)(nil)).Elem(), &TypeOptions{
				Type: "C1",
			})
			So(err, ShouldBeNil)
			So(v.(B).Fun(), ShouldEqual, 1)
		}

		{
			v, err := New(reflect.TypeOf((*B)(nil)).Elem(), &TypeOptions{
				Type: "C2",
				Options: &C2Options{
					Val: 2,
				},
			})
			So(err, ShouldBeNil)
			So(v.(B).Fun(), ShouldEqual, 2)
		}

		{
			v, err := New(reflect.TypeOf((*B)(nil)).Elem(), &TypeOptions{
				Type: "C3",
			})
			So(err, ShouldBeNil)
			So(v.(B).Fun(), ShouldEqual, 3)
		}
	})
}
