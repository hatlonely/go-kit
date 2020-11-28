package bind

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/strx"
)

func TestBind(t *testing.T) {
	Convey("TestBind", t, func() {
		type A struct {
			Key1 string
			Key2 string
			Key4 time.Duration `dft:"4s"`
			Key6 []string      `dft:"hello,world"`
		}

		type B struct {
			A
			A1   A
			A2   *A
			A3   []A
			A4   []*A
			A5   map[string]A // 不支持，因为无法知道有哪些 key
			A6   []A
			Key3 string
			key5 int // ignore unexported field
		}

		var b B
		err := Bind(&b, []Getter{MapGetter(map[string]interface{}{
			"A1.Key1":       "val1",
			"A1.Key2":       "val2",
			"A2.Key1":       "val3",
			"A2.Key2":       "val4",
			"Key3":          "val5",
			"A3[0].Key1":    "val6",
			"A3[0].Key2":    "val7",
			"A4[0].Key1":    "val8",
			"A4[0].Key2":    "val9",
			"A4[1].Key1":    "val10",
			"A4[1].Key2":    "val11",
			"A5.hello.Key1": "val12",
			"A5.hello.Key2": "val13",
			"Key1":          "val14",
			"Key2":          "val15",
			"Key4":          "1m",
		})})
		So(err, ShouldBeNil)
		fmt.Println(strx.JsonMarshalIndent(b))
		So(b.A1.Key1, ShouldEqual, "val1")
		So(b.A1.Key2, ShouldEqual, "val2")
		So(b.A2.Key1, ShouldEqual, "val3")
		So(b.A2.Key2, ShouldEqual, "val4")
		So(b.Key3, ShouldEqual, "val5")
		So(b.A3[0].Key1, ShouldEqual, "val6")
		So(b.A3[0].Key2, ShouldEqual, "val7")
		So(b.A4[0].Key1, ShouldEqual, "val8")
		So(b.A4[0].Key2, ShouldEqual, "val9")
		So(b.A4[1].Key1, ShouldEqual, "val10")
		So(b.A4[1].Key2, ShouldEqual, "val11")
		So(b.Key1, ShouldEqual, "val14")
		So(b.Key2, ShouldEqual, "val15")
		So(b.Key4, ShouldEqual, 1*time.Minute)
		So(b.A1.Key4, ShouldEqual, 4*time.Second)
		So(b.A1.Key6, ShouldResemble, []string{"hello", "world"})
	})
}
