package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/strx"
)

func TestTomlDecoder(t *testing.T) {
	Convey("TestYamlDecoder", t, func() {
		decoder, err := NewDecoderWithOptions(&DecoderOptions{
			Type: "Toml",
		})
		So(err, ShouldBeNil)

		s, err := NewStorage(map[string]interface{}{
			"Key1": "val1",
			"Key2": map[string]interface{}{
				"Key3": "val3",
				"Key4": 4,
				"Key5": map[string]interface{}{
					"Key6": "val6",
					"Key7": "val7",
				},
			},
			"Key8": map[string]interface{}{
				"Key9": map[string]interface{}{
					"Key10": 10,
					"Key11": 11,
				},
			},
			"Key12": []interface{}{
				map[string]interface{}{
					"Key13": "val13",
				},
				map[string]interface{}{
					"Key14": "val14",
				},
			},
		})
		So(err, ShouldBeNil)

		{
			buf, err := decoder.Encode(s)
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, `Key1 = "val1"

[[Key12]]
  Key13 = "val13"

[[Key12]]
  Key14 = "val14"

[Key2]
  Key3 = "val3"
  Key4 = 4
  [Key2.Key5]
    Key6 = "val6"
    Key7 = "val7"

[Key8]
  [Key8.Key9]
    Key10 = 10
    Key11 = 11
`)
		}

		{
			s, err := decoder.Decode([]byte(`Key1 = "val1"
# comment 1

[[Key12]]
  Key13 = "val13"  # comment 2

[[Key12]]
  Key14 = "val14"

[Key2]
  Key3 = "val3"
  Key4 = 4
  [Key2.Key5]
    Key6 = "val6"
    Key7 = "val7"

[Key8]
  [Key8.Key9]
    Key10 = 10
    Key11 = 11`))
			So(err, ShouldBeNil)
			So(strx.JsonMarshal(s.Interface()), ShouldResemble, strx.JsonMarshal(map[string]interface{}{
				"Key1": "val1",
				"Key2": map[string]interface{}{
					"Key3": "val3",
					"Key4": 4,
					"Key5": map[string]interface{}{
						"Key6": "val6",
						"Key7": "val7",
					},
				},
				"Key8": map[string]interface{}{
					"Key9": map[string]interface{}{
						"Key10": 10,
						"Key11": 11,
					},
				},
				"Key12": []interface{}{
					map[string]interface{}{
						"Key13": "val13",
					},
					map[string]interface{}{
						"Key14": "val14",
					},
				},
			}))
		}
	})
}
