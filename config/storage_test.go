package config

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestJsonUnmarshal(t *testing.T) {
	Convey("TestJsonUnmarshal", t, func() {
		a := struct {
			I  int
			S  string
			PS *string
			D  time.Duration // cannot unmarshal string into Go struct field A.D of type time.Duration
		}{I: 10}

		So(json.Unmarshal([]byte(`{"S": "str", "PS": "pstr", "D": "5m"}`), &a), ShouldNotBeNil)
	})
}

func TestStorage_Encrypt(t *testing.T) {
	Convey("TestStorage_Encrypt", t, func() {
		patches := ApplyMethod(reflect.TypeOf(&AESCipher{}), "Encrypt", func(aesCipher *AESCipher, text []byte) ([]byte, error) {
			return []byte("encrypt-" + string(text)), nil
		})
		defer patches.Reset()

		storage, err := NewStorage(map[interface{}]interface{}{
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
		})
		So(err, ShouldBeNil)
		cipher, _ := NewAESCipher([]byte("123456"))
		storage.SetCipherKeys([]string{"Key2", "Key3[0].Key4", "Key3[0].Key5", "Key3[0].Key6.Key7"})
		So(storage.Encrypt(NewCipherGroup(cipher, NewBase64Cipher())), ShouldBeNil)
		So(storage.root, ShouldResemble, map[interface{}]interface{}{
			"Key1":  1,
			"@Key2": base64.StdEncoding.EncodeToString([]byte("encrypt-val2")),
			"Key3": []interface{}{
				map[string]interface{}{
					"@Key4": base64.StdEncoding.EncodeToString([]byte("encrypt-val4")),
					"@Key5": base64.StdEncoding.EncodeToString([]byte("encrypt-val5")),
					"Key6": map[interface{}]interface{}{
						"@Key7": base64.StdEncoding.EncodeToString([]byte("encrypt-val7")),
						"Key8":  []interface{}{1, 2, 3},
					},
				},
			},
		})
	})
}

func TestStorage_Decrypt(t *testing.T) {
	Convey("TestStorage_Decrypt", t, func() {
		patches := ApplyMethod(reflect.TypeOf(&AESCipher{}), "Decrypt", func(aesCipher *AESCipher, encryptText []byte) ([]byte, error) {
			return bytes.Trim(encryptText, "encrypt-"), nil
		})
		defer patches.Reset()

		storage, err := NewStorage(map[interface{}]interface{}{
			"Key1":  1,
			"@Key2": base64.StdEncoding.EncodeToString([]byte("encrypt-val2")),
			"Key3": []interface{}{
				map[string]interface{}{
					"@Key4": base64.StdEncoding.EncodeToString([]byte("encrypt-val4")),
					"@Key5": base64.StdEncoding.EncodeToString([]byte("encrypt-val5")),
					"Key6": map[interface{}]interface{}{
						"@Key7": base64.StdEncoding.EncodeToString([]byte("encrypt-val7")),
						"Key8":  []interface{}{1, 2, 3},
					},
				},
			},
		})
		So(err, ShouldBeNil)
		cipher, _ := NewAESCipher([]byte("123456"))
		So(storage.Decrypt(NewCipherGroup(cipher, NewBase64Cipher())), ShouldBeNil)
		So(storage.root, ShouldResemble, map[interface{}]interface{}{
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
		})
		for _, key := range []string{
			"Key2", "Key3[0].Key4", "Key3[0].Key5", "Key3[0].Key6.Key7",
		} {
			So(storage.cipherKeySet[key], ShouldBeTrue)
		}
	})
}

func TestPrefixAppendKey(t *testing.T) {
	Convey("TestPrefixAppendKey", t, func() {
		So(prefixAppendKey("key1.key2", "key3"), ShouldEqual, "key1.key2.key3")
		So(prefixAppendKey("", "key1"), ShouldEqual, "key1")
	})
}

func TestPrefixAppendIdx(t *testing.T) {
	Convey("TestPrefixAppendIdx", t, func() {
		So(prefixAppendIdx("key1.key2", 3), ShouldEqual, "key1.key2[3]")
		So(prefixAppendIdx("", 3), ShouldEqual, "[3]")
	})
}

func TestGetToken(t *testing.T) {
	Convey("TestGetToken", t, func() {
		Convey("success", func() {
			for _, unit := range []struct {
				key  string
				info KeyInfo
				next string
			}{
				{key: "key1.key2", info: KeyInfo{key: "key1", mod: MapMod}, next: "key2"},
				{key: "[1].key", info: KeyInfo{idx: 1, mod: ArrMod}, next: "key"},
				{key: "[1][2]", info: KeyInfo{idx: 1, mod: ArrMod}, next: "[2]"},
				{key: "key", info: KeyInfo{key: "key", mod: MapMod}, next: ""},
				{key: "key[0]", info: KeyInfo{key: "key", mod: MapMod}, next: "[0]"},
			} {
				info, next, err := getToken(unit.key)
				So(err, ShouldBeNil)
				So(info.key, ShouldEqual, unit.info.key)
				So(info.mod, ShouldEqual, unit.info.mod)
				So(info.idx, ShouldEqual, unit.info.idx)
				So(next, ShouldEqual, unit.next)
			}
		})

		Convey("error", func() {
			for _, key := range []string{
				"[123", "[]", "[abc]", ".key1.key2",
			} {
				_, _, err := getToken(key)
				So(err, ShouldNotBeNil)
			}
		})
	})
}

func TestGetLastToken(t *testing.T) {
	Convey("TestGetLastToken", t, func() {
		Convey("success", func() {
			for _, unit := range []struct {
				key  string
				info KeyInfo
				prev string
			}{
				{key: "key[3]", info: KeyInfo{idx: 3, mod: ArrMod}, prev: "key"},
				{key: "key", info: KeyInfo{key: "key", mod: MapMod}, prev: ""},
				{key: "key1[3].key2", info: KeyInfo{key: "key2", mod: MapMod}, prev: "key1[3]"},
			} {
				info, next, err := getLastToken(unit.key)
				So(err, ShouldBeNil)
				So(info.key, ShouldEqual, unit.info.key)
				So(info.mod, ShouldEqual, unit.info.mod)
				So(info.idx, ShouldEqual, unit.info.idx)
				So(next, ShouldEqual, unit.prev)
			}
		})

		Convey("error", func() {
			for _, key := range []string{
				"123]", "[]", "[abc]", "key1.key2.",
			} {
				_, _, err := getLastToken(key)
				So(err, ShouldNotBeNil)
			}
		})
	})
}

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
			val, err := interfaceGet(v, unit.key)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, unit.val)
		}

		for _, key := range []string{
			"key3.key4", "key3[1]", "key3[abc]", "[4]", "key4",
		} {
			_, err := interfaceGet(v, key)
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
			So(interfaceSet(&v, unit.key, unit.val), ShouldBeNil)
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
			So(interfaceSet(&v, unit.key, unit.val), ShouldNotBeNil)
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
			So(interfaceDel(&v, key), ShouldNotBeNil)
		}

		for _, key := range []string{
			"key4", "key3[1]", // not exist key
			"key3[0].key6.key8[1]",
			"key3[0].key5",
		} {
			So(interfaceDel(&v, key), ShouldBeNil)
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
		}
		var opt Option
		So(interfaceToStruct(v, &opt), ShouldBeNil)
		So(opt.Key1, ShouldEqual, 1)
		So(opt.Key2, ShouldEqual, "val2")
		So(opt.Key3[0].Key4, ShouldEqual, "val4")
		So(opt.Key3[0].Key5, ShouldEqual, "val5")
		So(opt.Key3[0].Key6.Key7, ShouldEqual, "val7")
		So(opt.Key3[0].Key6.Key8[0], ShouldEqual, 1)
		So(opt.Key3[0].Key6.Key8[1], ShouldEqual, 2)
		So(opt.Key3[0].Key6.Key8[2], ShouldEqual, 3)
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
		So(interfaceToStruct(v, &opt), ShouldBeNil)
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
		So(interfaceToStruct(v, &opt), ShouldBeNil)
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
		So(interfaceToStruct(v, &opt), ShouldBeNil)
		So(opt.Key3[0]["Key4"], ShouldEqual, "val4")
		So(opt.Key3[0]["Key5"], ShouldEqual, "val5")
		So(opt.Key3[0]["Key6"], ShouldResemble, map[interface{}]interface{}{
			"Key7": "val7",
			"Key8": []interface{}{1, 2, 3},
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
		}

		kvs := map[string]interface{}{}
		err := interfaceTravel(v, func(key string, val interface{}) error {
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
			keys, err := interfaceDiff(v1, v2)
			So(err, ShouldBeNil)
			keySet := sliceToSet(keys)
			for _, key := range []string{
				"Key2", "Key3[0].Key6.Key8[2]", "Key3[0].Key6.Key8[3]", "Key3[0].Key9",
			} {
				So(keySet[key], ShouldBeTrue)
			}
		}
		{
			keys, err := interfaceDiff(v2, v1)
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
