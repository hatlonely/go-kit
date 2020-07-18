package config

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKMSCipher(t *testing.T) {
	Convey("TestKMSCipher", t, func() {
		kmsCli, err := kms.NewClientWithAccessKey(
			"cn-shanghai",
			"xx",
			"xx",
		)
		So(err, ShouldBeNil)

		cipher, err := NewKMSCipher(kmsCli, "9f2d041b-2fb1-46a6-b37f-1f53edcf8414")
		So(err, ShouldBeNil)

		buf, err := cipher.Encrypt([]byte("hello world"))
		So(err, ShouldBeNil)
		fmt.Println(string(buf))
		buf, err = cipher.Decrypt(buf)
		So(err, ShouldBeNil)
		So(string(buf), ShouldEqual, "hello world")

		fmt.Println(cipher.GenerateDataKey())
	})
}
