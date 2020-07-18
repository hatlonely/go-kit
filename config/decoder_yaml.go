package config

import (
	"gopkg.in/yaml.v2"
)

type YamlDecoder struct{}

func (d *YamlDecoder) Decode(buf []byte) (*Storage, error) {
	var data interface{}
	if err := yaml.Unmarshal(buf, &data); err != nil {
		return nil, err
	}
	return NewStorage(data)
}

func (d *YamlDecoder) Encode(storage *Storage) ([]byte, error) {
	return yaml.Marshal(storage.Interface())
}
