package strx

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDiff(t *testing.T) {
	Convey("TestDiff", t, func() {
		fmt.Println(Diff(
			JsonMarshalSortKeys(map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
			}),
			JsonMarshalSortKeys(map[string]interface{}{
				"key1": "val1",
				"key2": "val3",
				"key4": "val4",
			}),
		))

		fmt.Println(Diff(
			JsonMarshalIndentSortKeys(map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
			}),
			JsonMarshalIndentSortKeys(map[string]interface{}{
				"key1": "val1",
				"key2": "val3",
				"key4": "val4",
			}),
		))
	})
}

func TestJsonDiff(t *testing.T) {
	Convey("TestJsonDiff", t, func() {
		fmt.Println(JsonDiff(
			JsonMarshal(map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
			}),
			JsonMarshal(map[string]interface{}{
				"key1": "val1",
				"key2": "val3",
				"key4": "val4",
			}),
		))
	})
}
