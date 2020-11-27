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
		cipher, _ := NewAESCipher([]byte("123456"), "Zero")
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
		cipher, _ := NewAESCipher([]byte("123456"), "")
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
