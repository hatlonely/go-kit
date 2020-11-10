package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBase64Cipher(t *testing.T) {
	Convey("TestBase64Cipher", t, func() {
		Convey("case 1", func() {
			cipher, err := NewCipherWithOptions(&CipherOptions{
				Type: "Base64",
				Base64Cipher: Base64CipherOptions{
					Padding:     "=",
					URLEncoding: true,
				},
			})
			So(err, ShouldBeNil)
			{
				buf, err := cipher.Encrypt([]byte("hello world"))
				So(err, ShouldBeNil)
				So(string(buf), ShouldEqual, "aGVsbG8gd29ybGQ=")
			}
			{
				buf, err := cipher.Decrypt([]byte("aGVsbG8gd29ybGQ="))
				So(err, ShouldBeNil)
				So(string(buf), ShouldEqual, "hello world")
			}
		})

		Convey("case 2", func() {
			cipher, err := NewCipherWithOptions(&CipherOptions{
				Type: "Base64",
				Base64Cipher: Base64CipherOptions{
					Padding:     "=",
					StdEncoding: true,
				},
			})
			So(err, ShouldBeNil)
			{
				buf, err := cipher.Encrypt([]byte("hello world"))
				So(err, ShouldBeNil)
				So(string(buf), ShouldEqual, "aGVsbG8gd29ybGQ=")
			}
			{
				buf, err := cipher.Decrypt([]byte("aGVsbG8gd29ybGQ="))
				So(err, ShouldBeNil)
				So(string(buf), ShouldEqual, "hello world")
			}
		})

		Convey("case 3", func() {
			cipher, err := NewCipherWithOptions(&CipherOptions{
				Type: "Base64",
				Base64Cipher: Base64CipherOptions{
					Encoding: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_",
				},
			})
			So(err, ShouldBeNil)
			{
				buf, err := cipher.Encrypt([]byte("hello world"))
				So(err, ShouldBeNil)
				So(string(buf), ShouldEqual, "AgvSBg8GD29YBgq=")
			}
			{
				buf, err := cipher.Decrypt([]byte("AgvSBg8GD29YBgq="))
				So(err, ShouldBeNil)
				So(string(buf), ShouldEqual, "hello world")
			}
		})
	})
}
