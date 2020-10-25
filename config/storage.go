package config

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/hatlonely/go-kit/refx"
)

func NewStorage(root interface{}) (*Storage, error) {
	if err := refx.Validate(root); err != nil {
		return nil, err
	}

	return &Storage{
		root:         root,
		cipherKeySet: map[string]bool{},
	}, nil
}

type Storage struct {
	root         interface{}
	cipherKeySet map[string]bool
}

func (s *Storage) SetCipherKeys(keys []string) {
	if s.cipherKeySet == nil {
		s.cipherKeySet = map[string]bool{}
	}
	for _, key := range keys {
		s.cipherKeySet[key] = true
	}
}

func (s *Storage) GetCipherKeys() []string {
	var keys []string
	for key, _ := range s.cipherKeySet {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (s Storage) Get(key string) (interface{}, error) {
	return refx.InterfaceGet(s.root, key)
}

func (s *Storage) Del(key string) error {
	return refx.InterfaceDel(&s.root, key)
}

func (s *Storage) Set(key string, val interface{}) error {
	return refx.InterfaceSet(&s.root, key, val)
}

func (s *Storage) Encrypt(cipher Cipher) error {
	if cipher == nil {
		return nil
	}
	for key := range s.cipherKeySet {
		val, err := s.Get(key)
		if err != nil {
			return err
		}
		text, ok := val.(string)
		if !ok {
			return fmt.Errorf("encrypt value should be a string. key: [%v]", key)
		}
		blob, err := cipher.Encrypt([]byte(text))
		if err != nil {
			return fmt.Errorf("encrypt failed. key: [%v], err: [%v]", key, err)
		}
		info, prev, err := getLastToken(key)
		if err != nil {
			return err
		}
		if err := s.Set(prefixAppendKey(prev, "@"+info.key), string(blob)); err != nil {
			return err
		}
		if err := s.Del(key); err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) Decrypt(cipher Cipher) error {
	if cipher == nil {
		return nil
	}
	cipherKeyValMap := map[string]string{}
	var toDeleteKeys []string
	if err := s.Travel(func(key string, val interface{}) error {
		info, prev, err := getLastToken(key)
		if err != nil {
			return err
		}
		if info.mod == ArrMod {
			return nil
		}
		if info.key[0] != '@' {
			return nil
		}

		blob, ok := val.(string)
		if !ok {
			return fmt.Errorf("decrypt value should be a string. key: [%v]", key)
		}
		text, err := cipher.Decrypt([]byte(blob))
		if err != nil {
			return err
		}
		cipherKeyValMap[prefixAppendKey(prev, info.key[1:])] = string(text)
		toDeleteKeys = append(toDeleteKeys, key)

		return nil
	}); err != nil {
		return err
	}

	for key, val := range cipherKeyValMap {
		if err := s.Set(key, val); err != nil {
			return err
		}
		s.cipherKeySet[key] = true
	}

	for _, key := range toDeleteKeys {
		if err := s.Del(key); err != nil {
			return err
		}
	}

	return nil
}

func (s Storage) Unmarshal(v interface{}, opts ...refx.Option) error {
	return refx.InterfaceToStruct(s.root, v, opts...)
}

func (s Storage) Sub(key string) *Storage {
	v, err := s.Get(key)
	if err != nil {
		return nil
	}
	return &Storage{root: v}
}

func (s Storage) SubArr(key string) ([]*Storage, error) {
	val, err := s.Get(key)
	if err != nil {
		return nil, err
	}
	vs, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unsupport slice type. key: [%v], type: [%v]", key, reflect.TypeOf(val))
	}
	var res []*Storage
	for _, v := range vs {
		storage, _ := NewStorage(v)
		res = append(res, storage)
	}
	return res, nil
}

func (s Storage) SubMap(key string) (map[string]*Storage, error) {
	val, err := s.Get(key)
	if err != nil {
		return nil, err
	}
	res := map[string]*Storage{}
	switch val.(type) {
	case map[string]interface{}:
		for k, v := range val.(map[string]interface{}) {
			res[k], _ = NewStorage(v)
		}
	case map[interface{}]interface{}:
		for k, v := range val.(map[interface{}]interface{}) {
			res[k.(string)], _ = NewStorage(v)
		}
	default:
		return nil, fmt.Errorf("unsupport map type. key: [%v], type: [%v]", key, reflect.TypeOf(val))
	}
	return res, nil
}

func (s Storage) Interface() interface{} {
	return s.root
}

func (s Storage) Diff(s2 Storage) ([]string, error) {
	return refx.InterfaceDiff(s.root, s2.root)
}

func (s Storage) Travel(fun func(key string, val interface{}) error) error {
	return refx.InterfaceTravelRecursive(s.root, fun, "")
}

func prefixAppendKey(prefix string, key string) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%v.%v", prefix, key)
}

func prefixAppendIdx(prefix string, idx int) string {
	if prefix == "" {
		return fmt.Sprintf("[%v]", idx)
	}
	return fmt.Sprintf("%v[%v]", prefix, idx)
}

func getLastToken(key string) (info KeyInfo, prev string, err error) {
	if key[len(key)-1] == ']' {
		pos := strings.LastIndex(key, "[")
		// "123]" => error
		if pos == -1 {
			return info, "", fmt.Errorf("miss '[' in key. key: [%v]", key)
		}
		sub := key[pos+1 : len(key)-1]
		// "[]" => error
		if sub == "" {
			return info, "", fmt.Errorf("idx should not be empty. key: [%v]", key)
		}
		// "[abc]" => error
		idx, err := strconv.Atoi(sub)
		if err != nil {
			return info, "", fmt.Errorf("idx to int fail. key: [%v], sub: [%v]", key, sub)
		}
		// "key[3]" => 3, "key"
		return KeyInfo{idx: idx, mod: ArrMod}, key[:pos], nil
	}
	pos := strings.LastIndex(key, ".")
	// "key" => "key", ""
	if pos == -1 {
		return KeyInfo{key: key, mod: MapMod}, "", nil
	}
	// "key1.key2." => error
	if key[pos+1:] == "" {
		return info, "", fmt.Errorf("key should not be empty. key: [%v]", key)
	}
	// "key1[3].key2" => "key2", "key1[3]"
	return KeyInfo{key: key[pos+1:], mod: MapMod}, key[:pos], nil
}

func getToken(key string) (info KeyInfo, next string, err error) {
	if key[0] == '[' {
		pos := strings.Index(key, "]")
		// "[123" => error
		if pos == -1 {
			return info, next, fmt.Errorf("miss ']' in key. key: [%v]", key)
		}
		// "[]" => error
		if key[1:pos] == "" {
			return info, next, fmt.Errorf("idx should not be empty. key: [%v]", key)
		}
		idx, err := strconv.Atoi(key[1:pos])
		// "[abc]" => error
		if err != nil {
			return info, next, fmt.Errorf("idx to int fail. key: [%v], sub: [%v]", key, key[1:pos])
		}
		// "[1].key" => "1", "key"
		if pos+1 < len(key) && key[pos+1] == '.' {
			return KeyInfo{idx: idx, mod: ArrMod}, key[pos+2:], nil
		}
		// "[1][2]" => 1, "[2]"
		return KeyInfo{idx: idx, mod: ArrMod}, key[pos+1:], nil
	}
	pos := strings.IndexAny(key, ".[")
	// "key" => "key", ""
	if pos == -1 {
		return KeyInfo{key: key, mod: MapMod}, "", nil
	}
	// "key[0]" => "key", "[0]"
	if key[pos] == '[' {
		return KeyInfo{key: key[:pos], mod: MapMod}, key[pos:], nil
	}
	// ".key1.key2" => error
	if key[:pos] == "" {
		return info, "", fmt.Errorf("key should not be empty. key: [%v]", key)
	}
	// "key1.key2.key3" => "key1", "key2.key3"
	return KeyInfo{key: key[:pos], mod: MapMod}, key[pos+1:], nil
}

const MapMod = 1
const ArrMod = 2

type KeyInfo struct {
	key string
	idx int
	mod int
}
