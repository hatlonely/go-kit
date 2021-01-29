package config

import (
	"bytes"

	"github.com/yosuke-furukawa/json5/encoding/json5"
)

func NewJson5Decoder() Json5Decoder {
	return Json5Decoder{}
}

type Json5Decoder struct{}

func (d Json5Decoder) Decode(buf []byte) (*Storage, error) {
	var data interface{}
	if err := json5.NewDecoder(bytes.NewReader(buf)).Decode(&data); err != nil {
		return nil, err
	}
	return NewStorage(data)
}

func (d Json5Decoder) Encode(storage *Storage) ([]byte, error) {
	return json5.MarshalIndent(storage.Interface(), "", "  ")
}
