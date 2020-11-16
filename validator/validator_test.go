package validator

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/go-playground/validator.v9"

	"github.com/hatlonely/go-kit/strx"
)

func TestRule(t *testing.T) {
	type A struct {
		Key1 string `rule:"x in ['world', 'hello']"`
		Key2 int    `rule:"x>=5 && x<=6"`
		Key3 string `rule:"x =~ '^[0-9]{6}$'"`
		Key4 string `rule:"isEmail(x)"`
		Key5 int64  `rule:"x in [0, 1, 2]"`
		key6 int
		Key7 struct {
			Key8 int `rule:"x in [5, 7]"`
		}
	}

	RegisterFunction("isEmail", func(str string) (bool, error) {
		return strx.ReEmail.MatchString(str), nil
	})

	Convey("TestRule", t, func() {
		obj := &A{
			Key1: "hello",
			Key2: 5,
			Key3: "123456",
			Key4: "hatlonely@foxmail.com",
			Key5: 1,
		}
		obj.Key7.Key8 = 7

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

		Convey("case 5", func() {
			obj.Key7.Key8 = 6
			err := Validate(obj)
			So(err.(*Error).Code, ShouldEqual, ErrRuleNotMatch)
			So(err.(*Error).Key, ShouldEqual, "Key7.Key8")
			So(err.(*Error).Val, ShouldEqual, 6)
		})
	})
}

func BenchmarkValidate(b *testing.B) {
	type A struct {
		Key1 string `rule:"x in ['world', 'hello']"`
		Key2 int    `rule:"x>=5 && x<=6"`
		Key3 string `rule:"x =~ '^[0-9]{6}$'"`
		Key4 string `rule:"isEmail(x)"`
		Key5 int64  `rule:"x in [0, 1, 2]"`
		key6 int
		Key7 struct {
			Key8 int `rule:"x in [5, 7]"`
		}
	}

	obj := &A{
		Key1: "hello",
		Key2: 5,
		Key3: "123456",
		Key4: "hatlonely@foxmail.com",
		Key5: 1,
	}
	obj.Key7.Key8 = 7

	b.Run("validate 1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := Validate(obj)
			_ = err
		}
	})

	b.Run("validate 2", func(b *testing.B) {
		v := MustCompile(A{})
		for i := 0; i < b.N; i++ {
			err := v.Validate(obj)
			_ = err
		}
	})
}

// https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
//BenchmarkValidatorV10/github.com/go-playground/validator.v9-12         	 3848056	       308 ns/op
//BenchmarkValidatorV10/github.com/hatlonely/go-kit/validator-12         	  579918	      1944 ns/op
//BenchmarkValidatorV10/raw_operation-12                                 	16793824	        71.0 ns/op
func BenchmarkValidatorV10(b *testing.B) {
	type User struct {
		Age  uint8  `validate:"gte=0,lte=130" rule:"x >= 0 && x <= 130"`
		Name string `validate:"oneof=hello world" rule:"x in ['hello', 'world']"`
	}

	user := &User{
		Age:  130,
		Name: "hello",
	}

	b.Run("github.com/go-playground/validator.v9", func(b *testing.B) {
		validate := validator.New()

		for i := 0; i < b.N; i++ {
			_ = validate.Struct(user)
		}
	})

	b.Run("github.com/hatlonely/go-kit/validator", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Validate(user)
		}
	})

	b.Run("raw operation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b1 := user.Age > 0 && user.Age <= 130
			_, b2 := map[string]struct{}{"hello": {}, "world": {}}[user.Name]
			b := b1 && b2
			_ = b
		}
	})
}
