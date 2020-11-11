package config

import (
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

type Decoder interface {
	Decode(buf []byte) (*Storage, error)
	Encode(storage *Storage) ([]byte, error)
}

func NewDecoderWithOptions(options *DecoderOptions) (Decoder, error) {
	switch options.Type {
	case "Yaml":
		return &YamlDecoder{}, nil
	case "", "Json", "Json5":
		return &Json5Decoder{}, nil
	case "Toml":
		return &TomlDecoder{}, nil
	case "Ini":
		return &IniDecoder{}, nil
	case "Prop", "Properties":
		return &PropDecoder{}, nil
	}

	return nil, errors.Errorf("unsupported decoder type [%v]", options.Type)
}

func NewDecoderWithConfig(cfg *Config, opts ...refx.Option) (Decoder, error) {
	var options DecoderOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "cfg.Unmarshal failed.")
	}
	return NewDecoderWithOptions(&options)
}

func NewDecoder(typ string) (Decoder, error) {
	return NewDecoderWithOptions(&DecoderOptions{Type: typ})
}

type DecoderOptions struct {
	Type string
}
