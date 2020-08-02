package binding

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/strex"
)

func TestParseTag(t *testing.T) {
	Convey("TestParseTag", t, func() {
		for _, unit := range []struct {
			fkey     string
			bTag     string
			dTag     string
			rt       interface{}
			required bool
			dftVal   interface{}
			key      string
		}{
			{fkey: "Key1", bTag: `key1;required`, dTag: `123`, rt: 0, required: true, dftVal: 123, key: "key1"},
			{fkey: "Key1", bTag: `key2;required`, dTag: `123`, rt: "", required: true, dftVal: "123", key: "key2"},
		} {
			info, err := parseTag(unit.fkey, unit.bTag, unit.dTag, reflect.TypeOf(unit.rt))
			So(err, ShouldBeNil)
			So(info.required, ShouldEqual, unit.required)
			So(info.dftVal, ShouldEqual, unit.dftVal)
			So(info.key, ShouldEqual, unit.key)
		}
	})
}

func TestCompile(t *testing.T) {
	Convey("TestCompile", t, func() {
		type A struct {
			Key1 string  `bind:"key1;required" dft:"val1"`
			Key2 int     `dft:"123"`
			Key3 *string `dft:"val3"`
		}

		type B struct {
			Key4 string `bind:"key3"`
			A1   A      `bind:"a1"`
			A2   *A     `bind:"a2"`
		}

		b, err := Compile(&B{})
		So(err, ShouldBeNil)
		fmt.Println(b.infos)
		So(b.infos["Key4"].key, ShouldEqual, "key3")
		So(b.infos["A1"].key, ShouldEqual, "a1")
		So(b.infos["A1.Key1"].key, ShouldEqual, "key1")
		So(b.infos["A1.Key1"].required, ShouldBeTrue)
		So(b.infos["A1.Key1"].dftVal, ShouldEqual, "val1")
		So(b.infos["A1.Key2"].key, ShouldEqual, "key2")
		So(b.infos["A1.Key2"].required, ShouldBeFalse)
		So(b.infos["A1.Key2"].dftVal, ShouldEqual, 123)
		So(b.infos["A1.Key3"].dftVal, ShouldEqual, "val3")
		So(b.infos["A2"].key, ShouldEqual, "a2")
		So(b.infos["A2.Key1"].key, ShouldEqual, "key1")
		So(b.infos["A2.Key1"].required, ShouldBeTrue)
		So(b.infos["A2.Key1"].dftVal, ShouldEqual, "val1")
		So(b.infos["A2.Key2"].key, ShouldEqual, "key2")
		So(b.infos["A2.Key2"].required, ShouldBeFalse)
		So(b.infos["A2.Key2"].dftVal, ShouldEqual, 123)
		So(b.infos["A2.Key3"].dftVal, ShouldEqual, "val3")
	})
}

func TestBind(t *testing.T) {
	Convey("TestBind", t, func() {
		type A struct {
			Key1 string  `bind:"key1;required"`
			Key2 int     `dft:"123"`
			Key3 *string `dft:"val3"`
			Key4 *int
		}

		type B struct {
			Key4 string `bind:"key4"`
			A1   A      `bind:"a1"`
			A2   *A     `bind:"a2"`
		}
		bind, err := Compile(&B{})
		So(err, ShouldBeNil)

		b := &B{}
		Convey("normal", func() {
			So(bind.Bind(b, MapGetter(map[string]interface{}{
				"a1.key1": "val1-1",
				"a2.key1": "val1-2",
				"a1.Key2": 456,
				"a2.Key3": "val3-2",
			})), ShouldBeNil)

			buf, _ := json.MarshalIndent(b, "  ", "  ")
			fmt.Printf(string(buf))
		})

		Convey("error", func() {
			So(bind.Bind(b, MapGetter(map[string]interface{}{
				"a1.key1": "val1-1",
			})), ShouldNotBeNil)
		})

		os.Setenv("A1_KEY1", "val1-1")
		os.Setenv("A2_KEY1", "val1-2")
		os.Setenv("A1_KEY2", "456")
		os.Setenv("A2_KEY3", "val3-2")
		Convey("env", func() {
			So(bind.Bind(b, NewEnvGetter()), ShouldBeNil)
			fmt.Println(strex.MustJsonMarshal(b))
		})
	})
}
