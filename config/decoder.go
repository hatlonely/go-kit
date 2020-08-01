package config

import (
	"fmt"
)

type Decoder interface {
	Decode(buf []byte) (*Storage, error)
	Encode(storage *Storage) ([]byte, error)
}

func NewDecoder(name string) (Decoder, error) {
	switch name {
	case "yaml":
		return &YamlDecoder{}, nil
	case "json", "json5":
		return &Json5Decoder{}, nil
	case "toml":
		return &TomlDecoder{}, nil
	case "ini":
		return &IniDecoder{}, nil
	case "prop", "properties":
		return &PropDecoder{}, nil
	default:
		return nil, fmt.Errorf("unsupport decoder type. type: [%v]", name)
	}
}

func NewDecoderWithConfig(conf *Config) (Decoder, error) {
	return NewDecoder(conf.GetStringD("name", "json"))
}
