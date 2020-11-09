package config

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
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

func TestKMSCipher(t *testing.T) {
	Convey("TestKMSCipher", t, func() {
		patches := ApplyMethod(reflect.TypeOf(&kms.Client{}), "Decrypt", func(
			cli *kms.Client, request *kms.DecryptRequest) (*kms.DecryptResponse, error) {
			return &kms.DecryptResponse{
				Plaintext: "hello world",
			}, nil
		}).ApplyMethod(reflect.TypeOf(&kms.Client{}), "Encrypt", func(
			cli *kms.Client, request *kms.EncryptRequest) (*kms.EncryptResponse, error) {
			return &kms.EncryptResponse{
				CiphertextBlob: "NWMzYmNjODQtNTgxMC00NGZmLTkwMTAtNWIwMGY1NzhiNTg1qUxIaKYQ+GSXgLaPbl/XXrENJLEX4xqIezxb3+qM23+THE8pcs9u",
			}, nil
		})
		defer patches.Reset()

		cipher, err := NewCipherWithOptions(&CipherOptions{
			Type: "KMS",
			KMSCipher: KMSCipherOptions{
				AccessKeyID:     "xx",
				AccessKeySecret: "xx",
				RegionID:        "cn-shanghai",
				KeyID:           "9f2d041b-2fb1-46a6-b37f-1f53edcf8414",
			},
		})
		So(err, ShouldBeNil)

		buf, err := cipher.Encrypt([]byte("hello world"))
		So(err, ShouldBeNil)
		fmt.Println(string(buf))
		buf, err = cipher.Decrypt(buf)
		So(err, ShouldBeNil)
		So(string(buf), ShouldEqual, "hello world")
	})
}
