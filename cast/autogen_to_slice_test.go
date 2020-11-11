package cast

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestToBoolSlice(t *testing.T) {
	Convey("TestToBoolSlice", t, func() {
		So(ToBoolSlice([]bool{true, false, true}), ShouldResemble, []bool{true, false, true})
		So(ToBoolSlice([]interface{}{"true", false, 1, 0}), ShouldResemble, []bool{true, false, true, false})
		So(ToBoolSlice([]string{"true", "1", "false"}), ShouldResemble, []bool{true, true, false})
		So(ToBoolSlice("true,1,false"), ShouldResemble, []bool{true, true, false})
	})
}
