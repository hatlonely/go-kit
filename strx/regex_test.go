package strx

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReEmail(t *testing.T) {
	Convey("TestReEmail", t, func() {
		So(ReEmail.MatchString("xx@yy.zz.com"), ShouldBeTrue)
		So(ReEmail.MatchString("abc@def.com"), ShouldBeTrue)
	})
}

func TestReIdentifier(t *testing.T) {
	Convey("TestReIdentifier", t, func() {
		So(ReIdentifier.MatchString("abc"), ShouldBeTrue)
		So(ReIdentifier.MatchString("abc1"), ShouldBeTrue)
		So(ReIdentifier.MatchString("abc1abc"), ShouldBeTrue)
		So(ReIdentifier.MatchString("Abc"), ShouldBeTrue)
		So(ReIdentifier.MatchString("ABC"), ShouldBeTrue)
		So(ReIdentifier.MatchString("abcAbc"), ShouldBeTrue)
		So(ReIdentifier.MatchString("abc1Abc"), ShouldBeTrue)
		So(ReIdentifier.MatchString("abc1Abc1"), ShouldBeTrue)
		So(ReIdentifier.MatchString("@abc"), ShouldBeFalse)
		So(ReIdentifier.MatchString("$abc"), ShouldBeFalse)
		So(ReIdentifier.MatchString("0abc"), ShouldBeFalse)
	})
}

func TestRePhone(t *testing.T) {
	Convey("TestRePhone", t, func() {
		So(RePhone.MatchString("13112345678"), ShouldBeTrue)
		So(RePhone.MatchString("1311234567"), ShouldBeFalse)
		So(RePhone.MatchString("12312345678"), ShouldBeFalse)
	})
}
