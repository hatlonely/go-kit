package strx

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDiff(t *testing.T) {
	Convey("TestDiff", t, func() {
		fmt.Println(Diff(
			MustJsonMarshal(map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
			}),
			MustJsonMarshal(map[string]interface{}{
				"key1": "val1",
				"key2": "val3",
				"key4": "val4",
			}),
		))

		fmt.Println(Diff(
			MustJsonMarshalIndent(map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
			}),
			MustJsonMarshalIndent(map[string]interface{}{
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
			MustJsonMarshal(map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
			}),
			MustJsonMarshal(map[string]interface{}{
				"key1": "val1",
				"key2": "val3",
				"key4": "val4",
			}),
		))
	})
}
