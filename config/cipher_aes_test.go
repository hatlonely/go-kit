package config

import (
	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	. "github.com/smartystreets/goconvey/convey"
)

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
