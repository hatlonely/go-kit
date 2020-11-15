package flag

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFlag(t *testing.T) {
	Convey("TestFlag", t, func() {
		flag := Flag{
			shorthandKeyMap: map[string]string{
				"k": "opt.key",
				"a": "opt.action",
			},
			nameKeyMap: map[string]string{
				"key":    "opt.key",
				"action": "opt.action",
			},
			root: map[string]interface{}{
				"opt": map[string]interface{}{
					"key":    "val1",
					"action": "add",
				},
			},
		}
		Convey("Get", func() {
			{
				val, ok := flag.Get("k")
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, "val1")
			}
			{
				val, ok := flag.Get("key")
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, "val1")
			}
			{
				val, ok := flag.Get("opt.key")
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, "val1")
			}
		})

		Convey("Set", func() {
			{
				So(flag.Set("k", "val2"), ShouldBeNil)
				val, ok := flag.Get("k")
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, "val2")
			}
			{
				So(flag.Set("opt.key3", "val3"), ShouldBeNil)
				val, ok := flag.Get("opt.key3")
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, "val3")
			}
		})
	})
}
