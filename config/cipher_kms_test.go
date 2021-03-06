package config

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/alics"
)

func TestNewKMSCipher(t *testing.T) {
	Convey("TestNewKMSCipher", t, func() {
		patches := ApplyFunc(alics.ECSMetaDataRegionID, func() (string, error) {
			return "cn-beijing", nil
		}).ApplyFunc(alics.ECSMetaDataRamSecurityCredentialsRole, func() (string, error) {
			return "test-role", nil
		})
		defer patches.Reset()

		kmsCli, err := NewKMSCipherWithOptions(&KMSCipherOptions{
			KeyID: "test-id",
		})
		So(err, ShouldBeNil)
		So(kmsCli, ShouldNotBeNil)
	})

	Convey("case region id error", t, func() {
		patches := ApplyFunc(alics.ECSMetaDataRegionID, func() (string, error) {
			return "", errors.New("error")
		}).ApplyFunc(alics.ECSMetaDataRamSecurityCredentialsRole, func() (string, error) {
			return "test-role", nil
		})
		defer patches.Reset()

		kmsCli, err := NewKMSCipherWithOptions(&KMSCipherOptions{
			KeyID: "test-id",
		})
		So(err, ShouldNotBeNil)
		So(kmsCli, ShouldBeNil)
	})

	Convey("case credential error", t, func() {
		patches := ApplyFunc(alics.ECSMetaDataRegionID, func() (string, error) {
			return "cn-beijing", nil
		}).ApplyFunc(alics.ECSMetaDataRamSecurityCredentialsRole, func() (string, error) {
			return "", errors.New("error")
		})
		defer patches.Reset()

		kmsCli, err := NewKMSCipherWithOptions(&KMSCipherOptions{
			KeyID: "test-id",
		})
		So(err, ShouldNotBeNil)
		So(kmsCli, ShouldBeNil)
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

		cipher, err := NewKMSCipherWithOptions(&KMSCipherOptions{
			AccessKeyID:     "xx",
			AccessKeySecret: "xx",
			RegionID:        "cn-shanghai",
			KeyID:           "9f2d041b-2fb1-46a6-b37f-1f53edcf8414",
		})
		So(err, ShouldBeNil)

		buf, err := cipher.Encrypt([]byte("hello world"))
		So(err, ShouldBeNil)
		fmt.Println(string(buf))
		buf, err = cipher.Decrypt(buf)
		So(err, ShouldBeNil)
		So(string(buf), ShouldEqual, "hello world")

		plain, blob, err := cipher.GenerateDataKey()
		So(err, ShouldBeNil)
		So(plain, ShouldEqual, "plain")
		So(blob, ShouldEqual, "blob")
	})
}
