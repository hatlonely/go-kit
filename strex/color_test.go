package strex

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConsole(t *testing.T) {
	Convey("test console", t, func() {
		fmt.Println(Render("hello world", FormatSetBold, FormatSetUnderline, ForegroundRed, BackgroundWhite))
		fmt.Println(Render("hello world"))
		fmt.Println(Render("hello world", ForegroundGreen, FormatSetBold, BackgroundLightYellow))
	})
}
