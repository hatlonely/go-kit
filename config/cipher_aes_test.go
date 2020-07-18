package config

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAESCipher(t *testing.T) {
	Convey("TestAESCipher", t, func() {
		cipher, err := NewAESCipher([]byte("123456"))
		So(err, ShouldBeNil)
		So(cipher, ShouldNotBeNil)

		buf, err := cipher.Encrypt([]byte("hello world"))
		So(err, ShouldBeNil)
		buf, err = cipher.Decrypt(buf)
		So(err, ShouldBeNil)
		So(string(buf), ShouldEqual, "hello world")
	})
}

func TestAESWithKMSKeyCipher(t *testing.T) {
	Convey("TestAESWithKMSKeyCipher", t, func() {
		kmsCli, err := kms.NewClientWithAccessKey(
			"cn-shanghai",
			"xx",
			"xx",
		)
		So(err, ShouldBeNil)
		cipher, err := NewAESWithKMSKeyCipher(kmsCli, "NWMzYmNjODQtNTgxMC00NGZmLTkwMTAtNWIwMGY1NzhiNTg129Uj83I4hoqFOFsKrx/SSiuSn+zOHr/vUVdi8t7z1Bw/swRjHwE5NoBV6wn8RMG5rM1pvgg70bZwEYjUHdzP+NS+AgiWmy/t")
		So(err, ShouldBeNil)
		So(cipher, ShouldNotBeNil)

		buf, err := cipher.Encrypt([]byte("hello world"))
		So(err, ShouldBeNil)
		buf, err = cipher.Decrypt(buf)
		So(err, ShouldBeNil)
		So(string(buf), ShouldEqual, "hello world")
	})
}
