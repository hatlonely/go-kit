package config

import (
	"github.com/BurntSushi/toml"
)

type TomlDecoder struct{}

func (d *TomlDecoder) Decode(buf []byte) (*Storage, error) {
	var data interface{}
	if _, err := toml.Decode(string(buf), &data); err != nil {
		return nil, err
	}
	return NewStorage(data)
}

func (d *TomlDecoder) Encode(storage *Storage) ([]byte, error) {
	panic("toml encode not support yet")
	return nil, nil
}
