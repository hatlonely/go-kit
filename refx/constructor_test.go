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
