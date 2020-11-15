package refx

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFormatKey(t *testing.T) {
	Convey("TestFormatKey", t, func() {
		So(FormatKey("hello-world", WithCamelName()), ShouldEqual, "helloWorld")
		So(FormatKey("hello-world", WithSnakeName()), ShouldEqual, "hello_world")
		So(FormatKey("helloWorld", WithKebabName()), ShouldEqual, "hello-world")
		So(FormatKey("helloWorld", WithPascalName()), ShouldEqual, "HelloWorld")
	})
}

func TestInterfaceGet(t *testing.T) {
	Convey("TestInterfaceGet", t, func() {
		v := map[interface{}]interface{}{
			"Key1": 1,
			"Key2": "val2",
			"Key3": []interface{}{
				map[string]interface{}{
					"Key4": "val4",
					"Key5": "val5",
					"Key6": map[interface{}]interface{}{
						"Key7": "val7",
						"Key8": []interface{}{1, 2, 3},
					},
				},
			},
			"Key9": []map[string]interface{}{{
				"Key10": "val10",
			}, {
				"Key11": "val11",
			}},
			"Key12": []int{1, 2, 3},
			"Key13": map[string]map[string]interface{}{
				"Key16": {
					"Key14": "val14",
					"Key15": 15,
				},
				"Key17": {
					"Key14": "val141",
					"Key15": 151,
				},
			},
			"Key18": map[interface{}]interface{}{
				"12": "24",
				13:   "26", // this key cannot get
			},
			"Key19": map[string]map[string]interface{}{
				"Key20": {
					"Key22": "val22",
					"Key23": 23,
				},
				"Key21": {
					"Key22": "val222",
					"Key23": 233,
				},
			},
		}

		for _, unit := range []struct {
			key string
			val interface{}
		}{
			{"Key1", 1},
			{"Key2", "val2"},
			{"Key3[0].Key4", "val4"},
			{"Key3[0].Key5", "val5"},
			{"Key3[0].Key6.Key7", "val7"},
			{"Key3[0].Key6.Key8[0]", 1},
			{"Key3[0].Key6.Key8[1]", 2},
			{"Key3[0].Key6.Key8[2]", 3},
			{"Key9[0].Key10", "val10"},
			{"Key9[1].Key11", "val11"},
			{"Key12[0]", 1},
			{"Key12[1]", 2},
			{"Key12[2]", 3},
			{"Key13.Key16.Key14", "val14"},
			{"Key13.Key16.Key15", 15},
			{"Key13.Key17.Key14", "val141"},
			{"Key13.Key17.Key15", 151},
			{"Key18.12", "24"},
			{"Key19.Key20.Key22", "val22"},
			{"Key19.Key20.Key23", 23},
			{"Key19.Key21.Key22", "val222"},
			{"Key19.Key21.Key23", 233},
		} {
			val, err := InterfaceGet(v, unit.key)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, unit.val)
		}

		for _, key := range []string{
			"Key3.Key4", "Key3[1]", "Key3[abc]", "[4]", "Key4", "Key18.13",
		} {
			_, err := InterfaceGet(v, key)
			So(err, ShouldNotBeNil)
		}
	})
}

func TestInterfaceSet(t *testing.T) {
	Convey("TestInterfaceSet", t, func() {
		var v interface{}

		for _, unit := range []struct {
			key string
			val interface{}
		}{
			{key: "key1", val: 1},
			{key: "key2", val: "val2"},
			{key: "key3[0].key4", val: "val4"},
			{key: "key3[0].key5", val: "val5"},
			{key: "key3[0].key6.key7", val: "val7"},
			{key: "key3[0].key6.key8[0]", val: 1},
			{key: "key3[0].key6.key8[1]", val: 2},
			{key: "key3[0].key6.key8[2]", val: 3},
		} {
			So(InterfaceSet(&v, unit.key, unit.val), ShouldBeNil)
		}

		So(v, ShouldResemble, map[string]interface{}{
			"key1": 1,
			"key2": "val2",
			"key3": []interface{}{
				map[string]interface{}{
					"key4": "val4",
					"key5": "val5",
					"key6": map[string]interface{}{
						"key7": "val7",
						"key8": []interface{}{1, 2, 3},
					},
				},
			},
		})

		for _, unit := range []struct {
			key string
			val interface{}
		}{
			{key: "key3[abc]", val: 1}, // parse key error
			{key: "[4]", val: 1},       // root is not a slice
			{key: "key3.key4", val: 1}, // key3 is not a map
			{key: "key3[2]", val: 1},   // index out of bounds
		} {
			So(InterfaceSet(&v, unit.key, unit.val), ShouldNotBeNil)
		}
	})
}

func TestInterfaceDel(t *testing.T) {
	Convey("TestInterfaceDel", t, func() {
		var v interface{}
		v = map[interface{}]interface{}{
			"key1": 1,
			"key2": "val2",
			"key3": []interface{}{
				map[string]interface{}{
					"key4": "val4",
					"key5": "val5",
					"key6": map[interface{}]interface{}{
						"key7": "val7",
						"key8": []interface{}{1, 2, 3},
					},
				},
			},
		}

		for _, key := range []string{
			"key3[abc]", "[4]",
		} {
			So(InterfaceDel(&v, key), ShouldNotBeNil)
		}

		for _, key := range []string{
			"key4",
			"key3[1]",  // not exist key
			"key33[1]", // not exist key
			"key3[0].key6.key8[1]",
			"key3[0].key5",
		} {
			So(InterfaceDel(&v, key), ShouldBeNil)
		}

		So(v, ShouldResemble, map[interface{}]interface{}{
			"key1": 1,
			"key2": "val2",
			"key3": []interface{}{
				map[string]interface{}{
					"key4": "val4",
					"key6": map[interface{}]interface{}{
						"key7": "val7",
						"key8": []interface{}{1, 3},
					},
				},
			},
		})
	})
}

func TestInterfaceToStruct(t *testing.T) {
	v := map[interface{}]interface{}{
		"Key1": 1,
		"Key2": "val2",
		"Key3": []interface{}{
			map[string]interface{}{
				"Key4": "val4",
				"Key5": "val5",
				"Key6": map[interface{}]interface{}{
					"Key7": "val7",
					"Key8": []interface{}{1, 2, 3},
				},
			},
		},
		"Key9": []map[string]interface{}{{
			"Key10": "val10",
		}, {
			"Key11": "val11",
		}},
		"Key12": []int{1, 2, 3},
		"Key13": map[string]map[string]interface{}{
			"Key16": {
				"Key14": "val14",
				"Key15": 15,
			},
			"Key17": {
				"Key14": "val141",
				"Key15": 151,
			},
		},
		"Key18": map[interface{}]interface{}{
			"12": "24",
			13:   "26",
		},
		"Key19": map[string]map[string]interface{}{
			"Key20": {
				"Key22": "val22",
				"Key23": 23,
			},
			"Key21": {
				"Key22": "val222",
				"Key23": 233,
			},
		},
	}

	Convey("TestInterfaceToStruct 1", t, func() {
		type Options struct {
			Key1 int
			Key2 string
			Key3 []struct {
				Key4 string
				Key5 string
				Key6 struct {
					Key7        string
					Key8        []int64
					KeyNotExist string
				}
			}
			Key9  []map[string]interface{}
			Key12 []interface{}
			Key13 map[string]struct {
				Key14 string
				Key15 int
			}
			Key18 map[int]int
			Key19 struct {
				Key20 struct {
					Key22 string
					Key23 int
					Key24 string
				}
				Key21 map[string]interface{}
			}
		}
		var options Options
		So(InterfaceToStruct(v, &options), ShouldBeNil)
		So(options.Key1, ShouldEqual, 1)
		So(options.Key2, ShouldEqual, "val2")
		So(options.Key3[0].Key4, ShouldEqual, "val4")
		So(options.Key3[0].Key5, ShouldEqual, "val5")
		So(options.Key3[0].Key6.Key7, ShouldEqual, "val7")
		So(options.Key3[0].Key6.Key8[0], ShouldEqual, 1)
		So(options.Key3[0].Key6.Key8[1], ShouldEqual, 2)
		So(options.Key3[0].Key6.Key8[2], ShouldEqual, 3)
		So(options.Key9[0]["Key10"], ShouldEqual, "val10")
		So(options.Key9[1]["Key11"], ShouldEqual, "val11")
		So(options.Key12, ShouldResemble, []interface{}{1, 2, 3})
		So(options.Key13["Key16"].Key14, ShouldEqual, "val14")
		So(options.Key13["Key16"].Key15, ShouldEqual, 15)
		So(options.Key13["Key17"].Key14, ShouldEqual, "val141")
		So(options.Key13["Key17"].Key15, ShouldEqual, 151)
		So(options.Key18[12], ShouldEqual, 24)
		So(options.Key18[13], ShouldEqual, 26)
		So(options.Key19.Key20.Key22, ShouldEqual, "val22")
		So(options.Key19.Key20.Key23, ShouldEqual, 23)
		So(options.Key19.Key20.Key24, ShouldEqual, "")
		So(options.Key19.Key21["Key22"], ShouldEqual, "val222")
		So(options.Key19.Key21["Key23"], ShouldEqual, 233)
	})

	Convey("TestInterfaceToStruct 2", t, func() {
		type Options struct {
			Key1 int
			Key2 string
			Key3 []struct {
				Key4 string
				Key5 string
				Key6 map[string]interface{}
			}
		}
		var options Options
		So(InterfaceToStruct(v, &options), ShouldBeNil)
		So(options.Key1, ShouldEqual, 1)
		So(options.Key2, ShouldEqual, "val2")
		So(options.Key3[0].Key4, ShouldEqual, "val4")
		So(options.Key3[0].Key5, ShouldEqual, "val5")
		So(options.Key3[0].Key6["Key7"], ShouldEqual, "val7")
		So(options.Key3[0].Key6["Key8"], ShouldResemble, []interface{}{1, 2, 3})
	})

	Convey("TestInterfaceToStruct 3", t, func() {
		type Options struct {
			Key1 int
			Key2 *string
			Key3 []*struct {
				Key4 string
				Key5 string
				Key6 *struct {
					Key7 *string
					Key8 []*int64
				}
			}
		}
		var options Options
		So(InterfaceToStruct(v, &options), ShouldBeNil)
		So(options.Key1, ShouldEqual, 1)
		So(*options.Key2, ShouldEqual, "val2")
		So(options.Key3[0].Key4, ShouldEqual, "val4")
		So(options.Key3[0].Key5, ShouldEqual, "val5")
		So(*options.Key3[0].Key6.Key7, ShouldEqual, "val7")
		So(*options.Key3[0].Key6.Key8[0], ShouldEqual, 1)
		So(*options.Key3[0].Key6.Key8[1], ShouldEqual, 2)
		So(*options.Key3[0].Key6.Key8[2], ShouldEqual, 3)
	})

	Convey("TestInterfaceToStruct 4", t, func() {
		type Options struct {
			Key3 []map[string]interface{}
		}
		var options Options
		So(InterfaceToStruct(v, &options), ShouldBeNil)
		So(options.Key3[0]["Key4"], ShouldEqual, "val4")
		So(options.Key3[0]["Key5"], ShouldEqual, "val5")
		So(options.Key3[0]["Key6"], ShouldResemble, map[interface{}]interface{}{
			"Key7": "val7",
			"Key8": []interface{}{1, 2, 3},
		})
	})
}

func TestInterfaceToStructWithCamelName(t *testing.T) {
	Convey("TestInterfaceToStruct snake name", t, func() {
		v := map[interface{}]interface{}{
			"key1Key1": 1,
			"key2Key2": "val2",
			"key3Key3": []interface{}{
				map[string]interface{}{
					"key4Key4": "val4",
					"key5Key5": "val5",
					"key6Key6": map[interface{}]interface{}{
						"key7Key7": "val7",
						"key8Key8": []interface{}{1, 2, 3},
					},
				},
			},
		}

		type Options struct {
			Key1Key1 int
			Key2Key2 string
			Key3Key3 []struct {
				Key4Key4 string
				Key5Key5 string
				Key6Key6 struct {
					Key7Key7 string
					Key8Key8 []int64
				}
			}
		}

		var options Options
		So(InterfaceToStruct(v, &options, WithCamelName()), ShouldBeNil)
		So(options.Key1Key1, ShouldEqual, 1)
		So(options.Key2Key2, ShouldEqual, "val2")
		So(options.Key3Key3[0].Key4Key4, ShouldEqual, "val4")
		So(options.Key3Key3[0].Key5Key5, ShouldEqual, "val5")
		So(options.Key3Key3[0].Key6Key6.Key7Key7, ShouldEqual, "val7")
		So(options.Key3Key3[0].Key6Key6.Key8Key8[0], ShouldEqual, 1)
		So(options.Key3Key3[0].Key6Key6.Key8Key8[1], ShouldEqual, 2)
		So(options.Key3Key3[0].Key6Key6.Key8Key8[2], ShouldEqual, 3)
	})
}

func TestInterfaceToStructWithSnakeName(t *testing.T) {
	Convey("TestInterfaceToStruct snake name", t, func() {
		v := map[interface{}]interface{}{
			"key1_key1": 1,
			"key2_key2": "val2",
			"key3_key3": []interface{}{
				map[string]interface{}{
					"key4_key4": "val4",
					"key5_key5": "val5",
					"key6_key6": map[interface{}]interface{}{
						"key7_key7": "val7",
						"key8_key8": []interface{}{1, 2, 3},
					},
				},
			},
		}

		type Options struct {
			Key1Key1 int
			Key2Key2 string
			Key3Key3 []struct {
				Key4Key4 string
				Key5Key5 string
				Key6Key6 struct {
					Key7Key7 string
					Key8Key8 []int64
				}
			}
		}

		var options Options
		So(InterfaceToStruct(v, &options, WithSnakeName()), ShouldBeNil)
		So(options.Key1Key1, ShouldEqual, 1)
		So(options.Key2Key2, ShouldEqual, "val2")
		So(options.Key3Key3[0].Key4Key4, ShouldEqual, "val4")
		So(options.Key3Key3[0].Key5Key5, ShouldEqual, "val5")
		So(options.Key3Key3[0].Key6Key6.Key7Key7, ShouldEqual, "val7")
		So(options.Key3Key3[0].Key6Key6.Key8Key8[0], ShouldEqual, 1)
		So(options.Key3Key3[0].Key6Key6.Key8Key8[1], ShouldEqual, 2)
		So(options.Key3Key3[0].Key6Key6.Key8Key8[2], ShouldEqual, 3)
	})
}

func TestInterfaceToStructWithKebabName(t *testing.T) {
	Convey("TestInterfaceToStruct snake name", t, func() {
		v := map[interface{}]interface{}{
			"Key1Key1": 1,
			"Key2Key2": "val2",
			"Key3Key3": []interface{}{
				map[string]interface{}{
					"Key4Key4": "val4",
					"Key5Key5": "val5",
					"Key6Key6": map[interface{}]interface{}{
						"Key7Key7": "val7",
						"Key8Key8": []interface{}{1, 2, 3},
					},
				},
			},
		}

		type Options struct {
			Key1_Key1 int
			Key2_Key2 string
			Key3_Key3 []struct {
				Key4Key4 string
				Key5Key5 string
				Key6Key6 struct {
					Key7Key7 string
					Key8Key8 []int64
				}
			}
		}

		var options Options
		So(InterfaceToStruct(v, &options, WithPascalName()), ShouldBeNil)
		So(options.Key1_Key1, ShouldEqual, 1)
		So(options.Key2_Key2, ShouldEqual, "val2")
		So(options.Key3_Key3[0].Key4Key4, ShouldEqual, "val4")
		So(options.Key3_Key3[0].Key5Key5, ShouldEqual, "val5")
		So(options.Key3_Key3[0].Key6Key6.Key7Key7, ShouldEqual, "val7")
		So(options.Key3_Key3[0].Key6Key6.Key8Key8[0], ShouldEqual, 1)
		So(options.Key3_Key3[0].Key6Key6.Key8Key8[1], ShouldEqual, 2)
		So(options.Key3_Key3[0].Key6Key6.Key8Key8[2], ShouldEqual, 3)
	})
}

func TestInterfaceToStructWithPascalName(t *testing.T) {
	Convey("TestInterfaceToStruct snake name", t, func() {
		v := map[interface{}]interface{}{
			"key1-key1": 1,
			"key2-key2": "val2",
			"key3-key3": []interface{}{
				map[string]interface{}{
					"key4-key4": "val4",
					"key5-key5": "val5",
					"key6-key6": map[interface{}]interface{}{
						"key7-key7": "val7",
						"key8-key8": []interface{}{1, 2, 3},
					},
				},
			},
		}

		type Options struct {
			Key1Key1 int
			Key2Key2 string
			Key3Key3 []struct {
				Key4Key4 string
				Key5Key5 string
				Key6Key6 struct {
					Key7Key7 string
					Key8Key8 []int64
				}
			}
		}

		var options Options
		So(InterfaceToStruct(v, &options, WithKebabName()), ShouldBeNil)
		So(options.Key1Key1, ShouldEqual, 1)
		So(options.Key2Key2, ShouldEqual, "val2")
		So(options.Key3Key3[0].Key4Key4, ShouldEqual, "val4")
		So(options.Key3Key3[0].Key5Key5, ShouldEqual, "val5")
		So(options.Key3Key3[0].Key6Key6.Key7Key7, ShouldEqual, "val7")
		So(options.Key3Key3[0].Key6Key6.Key8Key8[0], ShouldEqual, 1)
		So(options.Key3Key3[0].Key6Key6.Key8Key8[1], ShouldEqual, 2)
		So(options.Key3Key3[0].Key6Key6.Key8Key8[2], ShouldEqual, 3)
	})
}

func TestInterfaceToStructWithDefaultValue(t *testing.T) {
	Convey("TestInterfaceToStruct snake name", t, func() {
		type Options struct {
			Key1Key1 int    `dft:"1"`
			Key2Key2 string `dft:"val2"`
			Key3Key3 []struct {
				Key4Key4 string `dft:"val4"`
				Key5Key5 string `dft:"val5"`
				Key6Key6 struct {
					Key7Key7 string  `dft:"val7"`
					Key8Key8 []int64 `dft:"1,2,3"`
				}
			}
		}

		Convey("test default 1", func() {
			v := map[interface{}]interface{}{
				"key3Key3": []interface{}{
					map[string]interface{}{},
				},
			}

			var options Options
			So(InterfaceToStruct(v, &options, WithCamelName()), ShouldBeNil)
			So(options.Key1Key1, ShouldEqual, 1)
			So(options.Key2Key2, ShouldEqual, "val2")
			So(options.Key3Key3[0].Key4Key4, ShouldEqual, "val4")
			So(options.Key3Key3[0].Key5Key5, ShouldEqual, "val5")
			So(options.Key3Key3[0].Key6Key6.Key7Key7, ShouldEqual, "val7")
			So(options.Key3Key3[0].Key6Key6.Key8Key8, ShouldResemble, []int64{1, 2, 3})
		})

		Convey("test default disable default value", func() {
			v := map[interface{}]interface{}{
				"key3Key3": []interface{}{
					map[string]interface{}{},
				},
			}

			var options Options
			So(InterfaceToStruct(v, &options, WithCamelName(), WithDisableDefaultValue()), ShouldBeNil)
			So(options.Key1Key1, ShouldEqual, 0)
			So(options.Key2Key2, ShouldEqual, "")
			So(options.Key3Key3[0].Key4Key4, ShouldEqual, "")
			So(options.Key3Key3[0].Key5Key5, ShouldEqual, "")
			So(options.Key3Key3[0].Key6Key6.Key7Key7, ShouldEqual, "")
			So(options.Key3Key3[0].Key6Key6.Key8Key8, ShouldBeEmpty)
		})

		Convey("test default 2", func() {
			v := map[interface{}]interface{}{
				"key1Key1": 11,
				"key2Key2": "val22",
				"key3Key3": []interface{}{
					map[string]interface{}{
						"key4Key4": "val44",
						"key5Key5": "val55",
						"key6Key6": map[interface{}]interface{}{
							"key7Key7": "val77",
							"key8Key8": []interface{}{11, 22, 33},
						},
					},
					map[string]interface{}{},
				},
			}

			var options Options
			So(InterfaceToStruct(v, &options, WithCamelName()), ShouldBeNil)
			So(options.Key1Key1, ShouldEqual, 11)
			So(options.Key2Key2, ShouldEqual, "val22")
			So(options.Key3Key3[0].Key4Key4, ShouldEqual, "val44")
			So(options.Key3Key3[0].Key5Key5, ShouldEqual, "val55")
			So(options.Key3Key3[0].Key6Key6.Key7Key7, ShouldEqual, "val77")
			So(options.Key3Key3[0].Key6Key6.Key8Key8, ShouldResemble, []int64{11, 22, 33})
			So(options.Key3Key3[1].Key4Key4, ShouldEqual, "val4")
			So(options.Key3Key3[1].Key5Key5, ShouldEqual, "val5")
			So(options.Key3Key3[1].Key6Key6.Key7Key7, ShouldEqual, "val7")
			So(options.Key3Key3[1].Key6Key6.Key8Key8, ShouldResemble, []int64{1, 2, 3})
		})
	})
}

func TestInterfaceTravel(t *testing.T) {
	Convey("TestInterfaceTravel", t, func() {
		v := map[interface{}]interface{}{
			"Key1": 1,
			"Key2": "val2",
			"Key3": []interface{}{
				map[string]interface{}{
					"Key4": "val4",
					"Key5": "val5",
					"Key6": map[interface{}]interface{}{
						"Key7": "val7",
						"Key8": []interface{}{1, 2, 3},
					},
				},
			},
			"Key9": []map[string]interface{}{{
				"Key10": "val10",
			}, {
				"Key11": "val11",
			}},
			"Key12": []int{1, 2, 3},
			"Key13": map[string]map[string]interface{}{
				"Key16": {
					"Key14": "val14",
					"Key15": 15,
				},
				"Key17": {
					"Key14": "val141",
					"Key15": 151,
				},
			},
			"Key18": map[interface{}]interface{}{
				"12": "24",
				13:   "26",
			},
			"Key19": map[string]map[string]interface{}{
				"Key20": {
					"Key22": "val22",
					"Key23": 23,
				},
				"Key21": {
					"Key22": "val222",
					"Key23": 233,
				},
			},
		}

		kvs := map[string]interface{}{}
		err := InterfaceTravel(v, func(key string, val interface{}) error {
			kvs[key] = val
			return nil
		})
		So(err, ShouldBeNil)
		So(len(kvs), ShouldEqual, 23)
		So(kvs["Key1"], ShouldEqual, 1)
		So(kvs["Key2"], ShouldEqual, "val2")
		So(kvs["Key3[0].Key4"], ShouldEqual, "val4")
		So(kvs["Key3[0].Key5"], ShouldEqual, "val5")
		So(kvs["Key3[0].Key6.Key7"], ShouldEqual, "val7")
		So(kvs["Key3[0].Key6.Key8[0]"], ShouldEqual, 1)
		So(kvs["Key3[0].Key6.Key8[1]"], ShouldEqual, 2)
		So(kvs["Key3[0].Key6.Key8[2]"], ShouldEqual, 3)
		So(kvs["Key9[0].Key10"], ShouldEqual, "val10")
		So(kvs["Key9[1].Key11"], ShouldEqual, "val11")
		So(kvs["Key12[0]"], ShouldEqual, 1)
		So(kvs["Key12[1]"], ShouldEqual, 2)
		So(kvs["Key12[2]"], ShouldEqual, 3)
		So(kvs["Key13.Key16.Key14"], ShouldEqual, "val14")
		So(kvs["Key13.Key16.Key15"], ShouldEqual, 15)
		So(kvs["Key18.12"], ShouldEqual, "24")
		So(kvs["Key18.13"], ShouldEqual, "26")
		So(kvs["Key19.Key20.Key22"], ShouldEqual, "val22")
		So(kvs["Key19.Key20.Key23"], ShouldEqual, 23)
		So(kvs["Key19.Key21.Key22"], ShouldEqual, "val222")
		So(kvs["Key19.Key21.Key23"], ShouldEqual, 233)
	})
}

func TestInterfaceDiff(t *testing.T) {
	Convey("TestInterfaceDiff", t, func() {
		v1 := map[interface{}]interface{}{
			"Key1": 1,
			"Key2": "val2",
			"Key3": []interface{}{
				map[string]interface{}{
					"Key4": "val4",
					"Key5": "val5",
					"Key6": map[string]interface{}{
						"Key7": "val7",
						"Key8": []interface{}{1, 2, 4, 3},
					},
					"Key9": "val9",
				},
			},
		}
		v2 := map[string]interface{}{
			"Key1": 1,
			"Key2": "val3",
			"Key3": []interface{}{
				map[string]interface{}{
					"Key4": "val4",
					"Key5": "val5",
					"Key6": map[interface{}]interface{}{
						"Key7": "val7",
						"Key8": []interface{}{1, 2, 3},
					},
					"Key10": "val10",
				},
			},
		}
		{
			keys, err := InterfaceDiff(v1, v2)
			So(err, ShouldBeNil)
			So(sliceToSet(keys), ShouldResemble, sliceToSet([]string{
				"Key2", "Key3[0].Key6.Key8[2]", "Key3[0].Key6.Key8[3]", "Key3[0].Key9",
			}))
		}
		{
			keys, err := InterfaceDiff(v2, v1)
			So(err, ShouldBeNil)
			So(sliceToSet(keys), ShouldResemble, sliceToSet([]string{
				"Key2", "Key3[0].Key6.Key8[2]", "Key3[0].Key10",
			}))
		}
	})
}

func sliceToSet(keys []string) map[string]bool {
	set := map[string]bool{}
	for _, key := range keys {
		set[key] = true
	}
	return set
}
