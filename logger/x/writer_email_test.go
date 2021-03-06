package loggerx

import (
	"net/smtp"
	"testing"
	"time"

	"github.com/hatlonely/go-kit/logger"

	. "github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEmailWriter_Write(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339Nano, "2020-11-07T16:07:36.65798+08:00")
	info := &logger.Info{
		Time:   ts,
		Level:  logger.LevelInfo,
		File:   "formatter_test.go:18",
		Caller: "TextFormatter",
		Fields: map[string]interface{}{
			"meta1": "val1",
			"meta2": "val2",
		},
		Data: map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
		},
	}

	Convey("TestEmailWriter_Write", t, func(c C) {
		patches := ApplyFunc(smtp.SendMail, func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			c.So(addr, ShouldEqual, "smtp.qq.com:25")
			c.So(a, ShouldResemble, smtp.PlainAuth("", "hatlonely@foxmail.com", "xx", "smtp.qq.com"))
			c.So(from, ShouldEqual, "hatlonely@foxmail.com")
			c.So(to, ShouldResemble, []string{"hatlonely@foxmail.com"})
			c.So(string(msg), ShouldStartWith, `From: hatlonely@foxmail.com
To: hatlonely@foxmail.com
Subject: 告警测试
Content-Type: text/html; charset=UTF-8;`)
			return nil
		})
		defer patches.Reset()

		w, err := NewEmailWriterWithOptions(&EmailWriterOptions{
			Level:      "Info",
			Server:     "smtp.qq.com",
			Port:       25,
			Username:   "hatlonely@foxmail.com",
			Password:   "xx",
			ToList:     []string{"hatlonely@foxmail.com"},
			Subject:    "告警测试",
			WorkerNum:  1,
			MsgChanLen: 200,
			Formatter: logger.FormatterOptions{
				Type: "Html",
				Options: &HtmlFormatterOptions{
					DefaultTitle: "测试告警",
				},
			},
		})
		So(err, ShouldBeNil)

		So(w.Write(info), ShouldBeNil)
		So(w.Close(), ShouldBeNil)
	})
}
