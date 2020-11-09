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
		}).ApplyMethod(reflect.TypeOf(&kms.Client{}), "GenerateDataKey", func(
			cli *kms.Client, request *kms.GenerateDataKeyRequest) (*kms.GenerateDataKeyResponse, error) {
			return &kms.GenerateDataKeyResponse{
				CiphertextBlob: "blob",
				Plaintext:      "plain",
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

		plain, blob, err := cipher.(*KMSCipher).GenerateDataKey()
		So(err, ShouldBeNil)
		So(plain, ShouldEqual, "plain")
		So(blob, ShouldEqual, "blob")

	})
}

func TestAESCipher(t *testing.T) {
	Convey("TestAESCipher 1", t, func() {
		cipher, err := NewCipherWithOptions(&CipherOptions{
			Type: "AES",
			AESCipher: AESCipherOptions{
				Key: "123456",
			},
		})
		So(err, ShouldBeNil)

		buf, err := cipher.Encrypt([]byte("hello world"))
		So(err, ShouldBeNil)
		buf, err = cipher.Decrypt(buf)
		So(err, ShouldBeNil)
		So(string(buf), ShouldEqual, "hello world")
	})

	Convey("TestAESCipher 2", t, func() {
		patches := ApplyMethod(reflect.TypeOf(&kms.Client{}), "Decrypt", func(
			cli *kms.Client, request *kms.DecryptRequest) (*kms.DecryptResponse, error) {
			return &kms.DecryptResponse{
				Plaintext: "nV+j39GExjRZCX2kUkYB/gDMTE9F7H15bZ3faddUMB4=",
			}, nil
		})
		defer patches.Reset()

		cipher, err := NewCipherWithOptions(&CipherOptions{
			Type: "AES",
			AESCipher: AESCipherOptions{
				KMSKey: "NWMzYmNjODQtNTgxMC00NGZmLTkwMTAtNWIwMGY1NzhiNTg129Uj83I4hoqFOFsKrx/SSiuSn+zOHr/vUVdi8t7z1Bw/swRjHwE5NoBV6wn8RMG5rM1pvgg70bZwEYjUHdzP+NS+AgiWmy/t",
				KMS: struct {
					AccessKeyID     string
					AccessKeySecret string
					RegionID        string
				}{
					AccessKeyID:     "xx",
					AccessKeySecret: "xx",
					RegionID:        "cn-shanghai",
				},
			},
		})
		So(err, ShouldBeNil)

		buf, err := cipher.Encrypt([]byte("hello world"))
		So(err, ShouldBeNil)
		buf, err = cipher.Decrypt(buf)
		So(err, ShouldBeNil)
		So(string(buf), ShouldEqual, "hello world")
	})
}
