package config

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

func init() {
	RegisterDecoder("Toml", func() *TomlDecoder {
		return &TomlDecoder{}
	})
}

// @see https://github.com/toml-lang/toml
type TomlDecoder struct{}

func (d *TomlDecoder) Decode(buf []byte) (*Storage, error) {
	var data interface{}
	if _, err := toml.Decode(string(buf), &data); err != nil {
		return nil, err
	}
	return NewStorage(data)
}

func (d *TomlDecoder) Encode(storage *Storage) ([]byte, error) {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(storage.Interface()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
