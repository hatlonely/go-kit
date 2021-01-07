package cli

import (
	"net/smtp"
	"testing"

	. "github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEmail(t *testing.T) {
	patches := ApplyFunc(smtp.SendMail, func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		So(addr, ShouldEqual, "smtp.qq.com:25")
		So(a, ShouldResemble, smtp.PlainAuth("", "hatlonely@foxmail.com", "123456", "smtp.qq.com"))
		So(from, ShouldEqual, "hatlonely@foxmail.com")
		So(to, ShouldResemble, []string{"hatlonely@foxmail.com"})
		So(string(msg), ShouldEqual, `From: hatlonely@foxmail.com
To: hatlonely@foxmail.com
Subject: hello world
Content-Type: text/html; charset=UTF-8;
hello world
`)
		return nil
	})
	defer patches.Reset()

	Convey("TestEmail", t, func() {
		cli := NewEmail("hatlonely@foxmail.com", "123456", WithEmailQQServer())

		So(cli.Send("hatlonely@foxmail.com", "hello world", "hello world"), ShouldBeNil)
	})
}
