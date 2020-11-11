package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/strx"
)

func TestDecoder(t *testing.T) {
	Convey("TestIniDecoder", t, func() {
		for _, option := range []*DecoderOptions{
			{Type: "Ini"},
			{Type: "Json"},
			{Type: "Json5"},
			{Type: "Yaml"},
			{Type: "Toml"},
			{Type: "Prop"},
		} {
			decoder, err := NewDecoderWithOptions(option)
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

			buf, err := decoder.Encode(s)
			So(err, ShouldBeNil)

			s, err = decoder.Decode(buf)
			So(err, ShouldBeNil)
			So(strx.JsonMarshalSortKeys(s.Interface()), ShouldBeIn, []string{
				"{\"Key1\":\"val1\",\"Key12\":[{\"Key13\":\"val13\"},{\"Key14\":\"val14\"}],\"Key2\":{\"Key3\":\"val3\",\"Key4\":\"4\",\"Key5\":{\"Key6\":\"val6\",\"Key7\":\"val7\"}},\"Key8\":{\"Key9\":{\"Key10\":\"10\",\"Key11\":\"11\"}}}",
				"{\"Key1\":\"val1\",\"Key12\":[{\"Key13\":\"val13\"},{\"Key14\":\"val14\"}],\"Key2\":{\"Key3\":\"val3\",\"Key4\":4,\"Key5\":{\"Key6\":\"val6\",\"Key7\":\"val7\"}},\"Key8\":{\"Key9\":{\"Key10\":10,\"Key11\":11}}}",
			})
		}
	})
}
