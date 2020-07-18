package config

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func NewStorage(root interface{}) (*Storage, error) {
	if err := validate(root); err != nil {
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
	return interfaceGet(s.root, key)
}

func (s *Storage) Del(key string) error {
	return interfaceDel(&s.root, key)
}

func (s *Storage) Set(key string, val interface{}) error {
	return interfaceSet(&s.root, key, val)
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

func (s Storage) Unmarshal(v interface{}) error {
	return interfaceToStruct(s.root, v)
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
	return interfaceDiff(s.root, s2.root)
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

func (s Storage) Travel(fun func(key string, val interface{}) error) error {
	return interfaceTravelRecursive(s.root, fun, "")
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

func interfaceDiff(v1 interface{}, v2 interface{}) ([]string, error) {
	var keys []string
	if err := interfaceTravel(v1, func(key string, val1 interface{}) error {
		val2, err := interfaceGet(v2, key)
		if err != nil || val1 != val2 {
			keys = append(keys, key)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return keys, nil
}

// todo remove value if map or slice is empty, use recursive implement
func interfaceDel(pv *interface{}, key string) error {
	info, prev, err := getLastToken(key)
	if err != nil {
		return err
	}
	v, err := interfaceGet(*pv, prev)
	if err != nil {
		return nil
	}
	if v == nil {
		return nil
	}
	if info.mod == ArrMod {
		val, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("unsupport slice type. prefix: [%v], type: [%v]", prev, reflect.TypeOf(*pv))
		}
		if info.idx >= len(val) {
			return nil
		}
		val = append(val[:info.idx], val[info.idx+1:]...)
		return interfaceSet(pv, prev, val)
	}

	switch v.(type) {
	case map[string]interface{}:
		delete(v.(map[string]interface{}), info.key)
	case map[interface{}]interface{}:
		delete(v.(map[interface{}]interface{}), info.key)
	default:
		return fmt.Errorf("unsupport slice type. prefix: [%v], type: [%v]", prev, reflect.TypeOf(*pv))
	}

	return nil
}

func interfaceGet(v interface{}, key string) (interface{}, error) {
	return interfaceGetRecursive(v, key, "")
}

func interfaceSet(pv *interface{}, key string, val interface{}) error {
	return interfaceSetRecursive(pv, key, val, "")
}

func interfaceSetRecursive(pv *interface{}, key string, val interface{}, prefix string) error {
	if key == "" {
		*pv = val
		return nil
	}
	info, next, err := getToken(key)
	if err != nil {
		return fmt.Errorf("get token failed. prefix: [%v], err: [%v]", prefix, err)
	}
	if info.mod == ArrMod {
		if *pv == nil {
			*pv = []interface{}{}
		}
		v, ok := (*pv).([]interface{})
		if !ok {
			return fmt.Errorf("unsupport slice type. prefix: [%v], type: [%v]", prefix, reflect.TypeOf(*pv))
		}
		if info.idx > len(v) {
			return fmt.Errorf("index out of bounds. prefix: [%v], index: [%v]", prefix, info.idx)
		}
		if info.idx < len(v) {
			sub := v[info.idx]
			if err := interfaceSetRecursive(&sub, next, val, prefixAppendIdx(prefix, info.idx)); err != nil {
				return err
			}
			v[info.idx] = sub
			return nil
		}

		var sub interface{}
		if err := interfaceSetRecursive(&sub, next, val, prefixAppendIdx(prefix, info.idx)); err != nil {
			return err
		}
		v = append(v, sub)
		*pv = v
		return nil
	}

	if *pv == nil {
		*pv = map[string]interface{}{}
	}
	switch (*pv).(type) {
	case map[string]interface{}:
		sub := (*pv).(map[string]interface{})[info.key]
		if err := interfaceSetRecursive(&sub, next, val, prefixAppendKey(prefix, info.key)); err != nil {
			return err
		}
		(*pv).(map[string]interface{})[info.key] = sub
	case map[interface{}]interface{}:
		sub := (*pv).(map[interface{}]interface{})[info.key]
		if err := interfaceSetRecursive(&sub, next, val, prefixAppendKey(prefix, info.key)); err != nil {
			return err
		}
		(*pv).(map[interface{}]interface{})[info.key] = sub
	default:
		return fmt.Errorf("unsupport map type. prefix: [%v], type: [%v]", prefix, reflect.TypeOf(*pv))
	}

	return nil
}

func interfaceGetRecursive(v interface{}, key string, prefix string) (interface{}, error) {
	if v == nil {
		return nil, fmt.Errorf("no such key. prefix: [%v], key: [%v]", prefix, key)
	}
	if key == "" {
		return v, nil
	}
	info, next, err := getToken(key)
	if err != nil {
		return nil, fmt.Errorf("get token failed. prefix: [%v], err: [%v]", prefix, err)
	}
	if info.mod == ArrMod {
		val, ok := v.([]interface{})
		if !ok {
			return nil, fmt.Errorf("node is not a slice. prefix: [%v], type: [%v]", prefix, reflect.TypeOf(v))
		}
		if info.idx >= len(val) {
			return nil, fmt.Errorf("index out of bounds. prefix: [%v], index: [%v]", prefix, info.idx)
		}
		return interfaceGetRecursive(val[info.idx], next, prefixAppendIdx(prefix, info.idx))
	}

	switch v.(type) {
	case map[string]interface{}:
		return interfaceGetRecursive(v.(map[string]interface{})[info.key], next, prefixAppendKey(prefix, info.key))
	case map[interface{}]interface{}:
		return interfaceGetRecursive(v.(map[interface{}]interface{})[info.key], next, prefixAppendKey(prefix, info.key))
	default:
		return nil, fmt.Errorf("node is not a map. prefix: [%v], type: [%v]", prefix, reflect.TypeOf(v))
	}
}

func interfaceToStruct(src interface{}, dst interface{}) error {
	return interfaceToStructRecursive(src, dst, "")
}

func interfaceToStructRecursive(src interface{}, dst interface{}, prefix string) error {
	if reflect.ValueOf(dst).Kind() != reflect.Ptr || dst == nil {
		return fmt.Errorf("invalid dst type or dst is nil. dst: [%v]", reflect.TypeOf(dst))
	}

	rv := reflect.ValueOf(dst).Elem()
	rt := reflect.TypeOf(dst).Elem()
	switch rt.Kind() {
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			key := rt.Field(i).Name
			var val interface{}
			switch src.(type) {
			case map[string]interface{}:
				val = src.(map[string]interface{})[key]
			case map[interface{}]interface{}:
				val = src.(map[interface{}]interface{})[key]
			default:
				return fmt.Errorf("convert src to map[string]interface{} or map[interface{}]interface{} failed. prefix: [%v], type: [%v]", prefix, reflect.TypeOf(src))
			}

			switch rv.Field(i).Type().Kind() {
			case reflect.Ptr:
				if rv.Field(i).IsNil() {
					nv := reflect.New(rv.Field(i).Type().Elem())
					rv.Field(i).Set(nv)
				}
				if err := interfaceToStructRecursive(val, rv.Field(i).Interface(), prefixAppendKey(prefix, key)); err != nil {
					return err
				}
			case reflect.Interface:
				rv.Field(i).Set(reflect.ValueOf(val))
			default:
				if err := interfaceToStructRecursive(val, rv.Field(i).Addr().Interface(), prefixAppendKey(prefix, key)); err != nil {
					return err
				}
			}
		}
	case reflect.Slice:
		eles, ok := src.([]interface{})
		if !ok {
			return fmt.Errorf("convert src to []interface{} failed. prefix: [%v], type: [%v]", prefix, reflect.TypeOf(src))
		}
		rv.Set(reflect.MakeSlice(rt, 0, len(eles)))
		for idx, ele := range eles {
			switch rt.Elem().Kind() {
			case reflect.Ptr:
				nv := reflect.New(rt.Elem().Elem())
				if err := interfaceToStructRecursive(ele, nv.Interface(), prefixAppendIdx(prefix, idx)); err != nil {
					return err
				}
				rv.Set(reflect.Append(rv, nv.Elem().Addr()))
			case reflect.Interface:
				rv.Set(reflect.Append(rv, reflect.ValueOf(ele)))
			default:
				nv := reflect.New(rt.Elem())
				if err := interfaceToStructRecursive(ele, nv.Interface(), prefixAppendIdx(prefix, idx)); err != nil {
					return err
				}
				rv.Set(reflect.Append(rv, nv.Elem()))
			}
		}
	case reflect.Map:
		if rt.Key().Kind() != reflect.String {
			return fmt.Errorf("unsupport dst type. prefix: [%v], type: [%v]", prefix, rt)
		}
		if rv.IsNil() {
			rv.Set(reflect.MakeMap(rt))
		}
		switch src.(type) {
		case map[string]interface{}:
			for key, val := range src.(map[string]interface{}) {
				switch rt.Elem().Kind() {
				case reflect.Ptr:
					nv := reflect.New(rt.Elem().Elem())
					if err := interfaceToStructRecursive(val, nv.Interface(), prefixAppendKey(prefix, key)); err != nil {
						return err
					}
					rv.SetMapIndex(reflect.ValueOf(key), nv.Elem().Addr())
				case reflect.Interface:
					rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
				default:
					nv := reflect.New(rt.Elem())
					if err := interfaceToStructRecursive(val, nv.Interface(), prefixAppendKey(prefix, key)); err != nil {
						return err
					}
					rv.SetMapIndex(reflect.ValueOf(key), nv.Elem())
				}
			}
		case map[interface{}]interface{}:
			for key, val := range src.(map[interface{}]interface{}) {
				switch rt.Elem().Kind() {
				case reflect.Ptr:
					nv := reflect.New(rt.Elem().Elem())
					if err := interfaceToStructRecursive(val, nv.Interface(), prefixAppendKey(prefix, key.(string))); err != nil {
						return err
					}
					rv.SetMapIndex(reflect.ValueOf(key), nv.Elem().Addr())
				case reflect.Interface:
					rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
				default:
					nv := reflect.New(rt.Elem())
					if err := interfaceToStructRecursive(val, nv.Interface(), prefixAppendKey(prefix, key.(string))); err != nil {
						return err
					}
					rv.SetMapIndex(reflect.ValueOf(key), nv.Elem())
				}
			}
		default:
			return fmt.Errorf("convert src to map[string]interface{} or map[interface{}]interface{} failed. prefix: [%v], type: [%v]", prefix, reflect.TypeOf(src))
		}
	default:
		if err := SetInterface(dst, src); err != nil {
			return fmt.Errorf("set interface failed. prefix: [%v], err: [%v]", prefix, err)
		}
	}

	return nil
}

func validate(v interface{}) error {
	return interfaceTravel(v, func(key string, val interface{}) error { return nil })
}

func interfaceTravel(v interface{}, fun func(key string, val interface{}) error) error {
	return interfaceTravelRecursive(v, fun, "")
}

// return error if fun is error immediately
func interfaceTravelRecursive(v interface{}, fun func(key string, val interface{}) error, prefix string) error {
	if v == nil {
		return nil
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.Ptr:
		return fmt.Errorf("key [%v], unsupport type [%v]", prefix, reflect.TypeOf(v))
	case reflect.Map:
		switch v.(type) {
		case map[interface{}]interface{}:
			for key, val := range v.(map[interface{}]interface{}) {
				if _, ok := key.(string); !ok {
					return fmt.Errorf("key [%v.%v], unsupport type [%v]", prefix, key, reflect.TypeOf(v))
				}
				if err := interfaceTravelRecursive(val, fun, prefixAppendKey(prefix, key.(string))); err != nil {
					return err
				}
			}
		case map[string]interface{}:
			for key, val := range v.(map[string]interface{}) {
				if err := interfaceTravelRecursive(val, fun, prefixAppendKey(prefix, key)); err != nil {
					return err
				}
			}
		}
	case reflect.Slice:
		switch v.(type) {
		case []interface{}:
			for idx, ele := range v.([]interface{}) {
				if err := interfaceTravelRecursive(ele, fun, prefixAppendIdx(prefix, idx)); err != nil {
					return err
				}
			}
			return nil
		default:
			return fmt.Errorf("key [%v], unsupport type [%v]", prefix, reflect.TypeOf(v))
		}
	default:
		if err := fun(prefix, v); err != nil {
			return err
		}
	}

	return nil
}

const MapMod = 1
const ArrMod = 2

type KeyInfo struct {
	key string
	idx int
	mod int
}
