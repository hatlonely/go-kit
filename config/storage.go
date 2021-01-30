package config

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

func NewStorage(root interface{}) (*Storage, error) {
	if err := refx.InterfaceTravel(root, func(key string, val interface{}) error { return nil }); err != nil {
		return nil, errors.WithMessage(err, "NewStorage failed")
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
	s.cipherKeySet = map[string]bool{}
	for _, key := range keys {
		s.cipherKeySet[key] = true
	}
}

func (s *Storage) AddCipherKeys(keys []string) {
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
		info, prev, err := refx.GetLastToken(key)
		if err != nil {
			return err
		}
		if err := s.Set(refx.PrefixAppendKey(prev, "@"+info.Key), string(blob)); err != nil {
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
		info, prev, err := refx.GetLastToken(key)
		if err != nil {
			return err
		}
		if info.Mod == refx.ArrMod {
			return nil
		}
		if info.Key[0] != '@' {
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
		cipherKeyValMap[refx.PrefixAppendKey(prev, info.Key[1:])] = string(text)
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

func (s *Storage) Interface() interface{} {
	if s == nil {
		return nil
	}
	return s.root
}

func (s Storage) Diff(s2 Storage) ([]string, error) {
	return refx.InterfaceDiff(s.root, s2.root)
}

func (s Storage) Travel(fun func(key string, val interface{}) error) error {
	return refx.InterfaceTravel(s.root, fun)
}
