package strex

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAllAny(t *testing.T) {
	Convey("test all and any", t, func() {
		So(All("1234567890", IsDigit), ShouldBeTrue)
		So(All("12345|67890", IsDigit), ShouldBeFalse)
		So(All("abcdefghijklmnopqrstuvwxyz", IsLower), ShouldBeTrue)
		So(All("abcdefghijklmnOpqrstuvwxyz", IsLower), ShouldBeFalse)
		So(Any("abcdefghijklmnOpqrstuvwxyz", IsUpper), ShouldBeTrue)
		So(Any("1234567890@", func(ch uint8) bool {
			return ch == '$' || ch == '@'
		}), ShouldBeTrue)
	})
}
