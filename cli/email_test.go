package cli

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEmail(t *testing.T) {
	Convey("TestEmail", t, func() {
		cli := NewEmail("hatlonely@foxmail.com", "xuckndegounrbfhf", WithEmailQQServer())

		So(cli.Send("hatlonely@foxmail.com", "hello world", "hello world"), ShouldBeNil)
	})
}
