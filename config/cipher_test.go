package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/refx"
)

func TestNewCipherWithOptions(t *testing.T) {
	Convey("TestNewCipherWithOptions", t, func() {
		Convey("case 1", func() {
			cipher, err := NewCipherWithOptions(&CipherOptions{
				Type: "Base64",
				Options: &Base64CipherOptions{
					Encoding: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_",
				},
			}, refx.WithCamelName())
			So(err, ShouldBeNil)
			So(cipher, ShouldNotBeNil)

			buf, err := cipher.Encrypt([]byte("hello world"))
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, "AgvSBg8GD29YBgq=")
		})

		Convey("case 2", func() {
			cipher, err := NewCipherWithOptions(&CipherOptions{
				Type: "Base64",
				Options: map[string]interface{}{
					"encoding": "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_",
				},
			}, refx.WithCamelName())
			So(err, ShouldBeNil)
			So(cipher, ShouldNotBeNil)

			buf, err := cipher.Encrypt([]byte("hello world"))
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, "AgvSBg8GD29YBgq=")
		})
	})
}
