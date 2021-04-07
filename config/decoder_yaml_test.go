package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/strx"
)

func TestYamlDecoder(t *testing.T) {
	Convey("TestYamlDecoder", t, func() {
		decoder, err := NewDecoderWithOptions(&DecoderOptions{
			Type: "Yaml",
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
					"Key12": "Hello\nWorld",
				},
			},
		})
		So(err, ShouldBeNil)

		{
			buf, err := decoder.Encode(s)
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, `Key1: val1
Key2:
  Key3: val3
  Key4: 4
  Key5:
    Key6: val6
    Key7: val7
Key8:
  Key9:
    Key10: 10
    Key11: 11
    Key12: |-
      Hello
      World
`)
		}

		{
			s, err := decoder.Decode([]byte(`Key1: val1
Key2:
  Key3: val3
  Key4: 4
  Key5:
    Key6: val6
    Key7: val7
Key8:
  Key9:
    Key10: 10
    Key11: 11
    Key12: |-
      Hello
      World
    Key13: "
Hello
World
"
`))
			So(err, ShouldBeNil)
			So(strx.JsonMarshalSortKeys(s.Interface()), ShouldResemble, strx.JsonMarshalSortKeys(map[string]interface{}{
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
						"Key12": "Hello\nWorld",
						"Key13": " Hello World ",
					},
				},
			}))
		}
	})
}
