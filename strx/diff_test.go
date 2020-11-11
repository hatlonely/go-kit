package strx

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDiff(t *testing.T) {
	Convey("TestDiff", t, func() {
		Convey("case normal", func() {
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
					"key4": "val5",
				}),
				JsonMarshalIndentSortKeys(map[string]interface{}{
					"key1": "val1",
					"key2": "val3",
					"key3": "val3",
					"key4": "val4",
				}),
			))
		})

		Convey("case null", func() {
			fmt.Println(Diff("null", "null"))
			fmt.Println(Diff("null", JsonMarshalIndentSortKeys(map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
			})))
			fmt.Println(Diff(JsonMarshalIndentSortKeys(map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
			}), "null"))
		})

		Convey("case empty", func() {
			fmt.Println(Diff("", ""))
			fmt.Println(Diff("", "test"))
			fmt.Println(Diff("test", ""))
		})
	})
}
