package validator

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/strex"
)

func TestRule(t *testing.T) {
	type A struct {
		Key1 string `rule:"x in ['world', 'hello']"`
		Key2 int    `rule:"x>=5 && x<=6"`
		Key3 string `rule:"x =~ '^[0-9]{6}$'"`
		Key4 string `rule:"isEmail(x)"`
		Key5 int64  `rule:"x in [0, 1, 2]"`
	}

	RegisterFunction("isEmail", func(str string) (bool, error) {
		return strex.ReEmail.MatchString(str), nil
	})

	Convey("TestRule", t, func() {
		obj := &A{
			Key1: "hello",
			Key2: 5,
			Key3: "123456",
			Key4: "hatlonely@foxmail.com",
			Key5: 1,
		}

		Convey("case 0", func() {
			So(Validate(obj), ShouldBeNil)
		})

		Convey("case 1", func() {
			obj.Key1 = "abc"
			err := Validate(obj)
			So(err.(*Error).Code, ShouldEqual, ErrRuleNotMatch)
			So(err.(*Error).Key, ShouldEqual, "Key1")
			So(err.(*Error).Val, ShouldEqual, "abc")
		})

		Convey("case 2", func() {
			obj.Key2 = 100
			err := Validate(obj)
			So(err.(*Error).Code, ShouldEqual, ErrRuleNotMatch)
			So(err.(*Error).Key, ShouldEqual, "Key2")
			So(err.(*Error).Val, ShouldEqual, 100)
		})

		Convey("case 3", func() {
			obj.Key3 = "abcdef"
			err := Validate(obj)
			So(err.(*Error).Code, ShouldEqual, ErrRuleNotMatch)
			So(err.(*Error).Key, ShouldEqual, "Key3")
			So(err.(*Error).Val, ShouldEqual, "abcdef")
		})

		Convey("case 4", func() {
			obj.Key4 = "hello world"
			err := Validate(obj)
			So(err.(*Error).Code, ShouldEqual, ErrRuleNotMatch)
			So(err.(*Error).Key, ShouldEqual, "Key4")
			So(err.(*Error).Val, ShouldEqual, "hello world")
		})
	})
}
