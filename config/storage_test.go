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
