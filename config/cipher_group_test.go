package config

import (
	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCipher(t *testing.T) {
	Convey("TestCipher", t, func() {
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

		for _, options := range []*CipherOptions{{
			Type: "Empty",
		}, {}, {
			Type: "Base64",
			Options: &Base64CipherOptions{
				URLEncoding: true,
			},
		}, {
			Type: "AES",
			Options: &AESCipherOptions{
				Key: "123456",
			},
		}, {
			Type: "KMS",
			Options: &KMSCipherOptions{
				AccessKeyID:     "xx",
				AccessKeySecret: "xx",
				RegionID:        "cn-shanghai",
				KeyID:           "xx",
			},
		}} {
			cipher, err := NewCipherWithOptions(options)
			So(err, ShouldBeNil)
			buf1, err := cipher.Encrypt([]byte(`hello world`))
			So(err, ShouldBeNil)
			buf2, err := cipher.Decrypt(buf1)
			So(err, ShouldBeNil)
			So(string(buf2), ShouldEqual, `hello world`)
		}
	})
}
