package refx

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInterfaceGet(t *testing.T) {
	Convey("TestInterfaceGet", t, func() {
		v := map[interface{}]interface{}{
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

		for _, unit := range []struct {
			key string
			val interface{}
		}{
			{"key1", 1},
			{"key2", "val2"},
			{"key3[0].key4", "val4"},
			{"key3[0].key5", "val5"},
			{"key3[0].key6.key7", "val7"},
			{"key3[0].key6.key8[0]", 1},
			{"key3[0].key6.key8[1]", 2},
			{"key3[0].key6.key8[2]", 3},
		} {
			val, err := InterfaceGet(v, unit.key)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, unit.val)
		}

		for _, key := range []string{
			"key3.key4", "key3[1]", "key3[abc]", "[4]", "key4",
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
			"key4", "key3[1]", // not exist key
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
	}

	Convey("TestInterfaceToStruct 1", t, func() {
		type Option struct {
			Key1 int
			Key2 string
			Key3 []struct {
				Key4 string
				Key5 string
				Key6 struct {
					Key7 string
					Key8 []int64
				}
			}
			Key9  []map[string]interface{}
			Key12 []interface{}
			Key13 map[string]struct {
				Key14 string
				Key15 int
			}
		}
		var opt Option
		So(InterfaceToStruct(v, &opt), ShouldBeNil)
		So(opt.Key1, ShouldEqual, 1)
		So(opt.Key2, ShouldEqual, "val2")
		So(opt.Key3[0].Key4, ShouldEqual, "val4")
		So(opt.Key3[0].Key5, ShouldEqual, "val5")
		So(opt.Key3[0].Key6.Key7, ShouldEqual, "val7")
		So(opt.Key3[0].Key6.Key8[0], ShouldEqual, 1)
		So(opt.Key3[0].Key6.Key8[1], ShouldEqual, 2)
		So(opt.Key3[0].Key6.Key8[2], ShouldEqual, 3)
		So(opt.Key9[0]["Key10"], ShouldEqual, "val10")
		So(opt.Key9[1]["Key11"], ShouldEqual, "val11")
		So(opt.Key12, ShouldResemble, []interface{}{1, 2, 3})
		So(opt.Key13["Key16"].Key14, ShouldEqual, "val14")
		So(opt.Key13["Key16"].Key15, ShouldEqual, 15)
		So(opt.Key13["Key17"].Key14, ShouldEqual, "val141")
		So(opt.Key13["Key17"].Key15, ShouldEqual, 151)
	})

	Convey("TestInterfaceToStruct 2", t, func() {
		type Option struct {
			Key1 int
			Key2 string
			Key3 []struct {
				Key4 string
				Key5 string
				Key6 map[string]interface{}
			}
		}
		var opt Option
		So(InterfaceToStruct(v, &opt), ShouldBeNil)
		So(opt.Key1, ShouldEqual, 1)
		So(opt.Key2, ShouldEqual, "val2")
		So(opt.Key3[0].Key4, ShouldEqual, "val4")
		So(opt.Key3[0].Key5, ShouldEqual, "val5")
		So(opt.Key3[0].Key6["Key7"], ShouldEqual, "val7")
		So(opt.Key3[0].Key6["Key8"], ShouldResemble, []interface{}{1, 2, 3})
	})

	Convey("TestInterfaceToStruct 3", t, func() {
		type Option struct {
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
		var opt Option
		So(InterfaceToStruct(v, &opt), ShouldBeNil)
		So(opt.Key1, ShouldEqual, 1)
		So(*opt.Key2, ShouldEqual, "val2")
		So(opt.Key3[0].Key4, ShouldEqual, "val4")
		So(opt.Key3[0].Key5, ShouldEqual, "val5")
		So(*opt.Key3[0].Key6.Key7, ShouldEqual, "val7")
		So(*opt.Key3[0].Key6.Key8[0], ShouldEqual, 1)
		So(*opt.Key3[0].Key6.Key8[1], ShouldEqual, 2)
		So(*opt.Key3[0].Key6.Key8[2], ShouldEqual, 3)
	})

	Convey("TestInterfaceToStruct 4", t, func() {
		type Option struct {
			Key3 []map[string]interface{}
		}
		var opt Option
		So(InterfaceToStruct(v, &opt), ShouldBeNil)
		So(opt.Key3[0]["Key4"], ShouldEqual, "val4")
		So(opt.Key3[0]["Key5"], ShouldEqual, "val5")
		So(opt.Key3[0]["Key6"], ShouldResemble, map[interface{}]interface{}{
			"Key7": "val7",
			"Key8": []interface{}{1, 2, 3},
		})
	})

	Convey("TestInterfaceToStruct camel name", t, func() {
		v := map[interface{}]interface{}{
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

		type Option struct {
			Key1 int
			Key2 string
			Key3 []struct {
				Key4 string
				Key5 string
				Key6 struct {
					Key7 string
					Key8 []int64
				}
			}
		}
		var opt Option
		So(InterfaceToStruct(v, &opt, WithCamelName()), ShouldBeNil)
		So(opt.Key1, ShouldEqual, 1)
		So(opt.Key2, ShouldEqual, "val2")
		So(opt.Key3[0].Key4, ShouldEqual, "val4")
		So(opt.Key3[0].Key5, ShouldEqual, "val5")
		So(opt.Key3[0].Key6.Key7, ShouldEqual, "val7")
		So(opt.Key3[0].Key6.Key8[0], ShouldEqual, 1)
		So(opt.Key3[0].Key6.Key8[1], ShouldEqual, 2)
		So(opt.Key3[0].Key6.Key8[2], ShouldEqual, 3)
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
		}

		kvs := map[string]interface{}{}
		err := InterfaceTravel(v, func(key string, val interface{}) error {
			kvs[key] = val
			return nil
		})
		So(err, ShouldBeNil)
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
			keySet := sliceToSet(keys)
			for _, key := range []string{
				"Key2", "Key3[0].Key6.Key8[2]", "Key3[0].Key6.Key8[3]", "Key3[0].Key9",
			} {
				So(keySet[key], ShouldBeTrue)
			}
		}
		{
			keys, err := InterfaceDiff(v2, v1)
			So(err, ShouldBeNil)
			keySet := sliceToSet(keys)
			for _, key := range []string{
				"Key2", "Key3[0].Key6.Key8[2]", "Key3[0].Key10",
			} {
				So(keySet[key], ShouldBeTrue)
			}
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
