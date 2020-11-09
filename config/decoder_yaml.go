package config

import (
	"bytes"

	"gopkg.in/yaml.v3"
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
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err := enc.Encode(storage.Interface())
	return buf.Bytes(), err
}
