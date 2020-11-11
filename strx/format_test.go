package strx

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJsonMarshal(t *testing.T) {
	Convey("TestJsonMarshal", t, func() {
		So(JsonMarshal(nil), ShouldEqual, "null")
		So(JsonMarshalIndent(nil), ShouldEqual, "null")
		So(JsonMarshalSortKeys(nil), ShouldEqual, "null")
		So(JsonMarshalIndentSortKeys(nil), ShouldEqual, "null")

		v := map[string]interface{}{
			"key1": "val1",
			"key2": 2,
			"key3": "val3",
			"key4": map[string]interface{}{
				"key5": "val5",
				"key6": 6,
				"key7": []int{1, 2, 3},
				"key8": []map[interface{}]interface{}{
					{
						"key9":  9,
						"key10": "val10",
					}, {
						"key11": 11,
						12:      "val12",
					},
				},
			},
		}
		fmt.Println(JsonMarshal(v))
		fmt.Println(JsonMarshalIndent(v))
		fmt.Println(JsonMarshalSortKeys(v))
		fmt.Println(JsonMarshalIndentSortKeys(v))
	})
}

func TestFormatSpace(t *testing.T) {
	Convey("TestFormatSpace", t, func() {
		So(FormatSpace("  hello  \t\n  world \r\n  hello  golang\n"), ShouldEqual, "hello world hello golang")
	})
}
