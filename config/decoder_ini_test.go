package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnescape(t *testing.T) {
	Convey("TestUnescape", t, func() {
		for _, unit := range []struct {
			str          string
			unescapedStr string
		}{
			{str: `\;\#\n\r`, unescapedStr: ";#\n\r"},
			{str: `\"`, unescapedStr: `"`},
			{str: `\'`, unescapedStr: `'`},
		} {
			So(unescape(unit.str), ShouldEqual, unit.unescapedStr)
		}
	})
}
