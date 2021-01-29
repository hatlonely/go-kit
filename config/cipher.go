package config

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

func init() {
	RegisterCipher("", func() (Cipher, error) {
		return EmptyCipher{}, nil
	})
	RegisterCipher("Empty", func() (Cipher, error) {
		return EmptyCipher{}, nil
	})
}

func RegisterCipher(key string, constructor interface{}) {
	if _, ok := cipherConstructorMap[key]; ok {
		panic(fmt.Sprintf("cipher type [%v] is already registered", key))
	}

	info, err := refx.NewConstructor(constructor, reflect.TypeOf((*Cipher)(nil)).Elem())
	refx.Must(err)

	cipherConstructorMap[key] = info
}

var cipherConstructorMap = map[string]*refx.Constructor{}

type Cipher interface {
	Encrypt(textToEncrypt []byte) ([]byte, error)
	Decrypt(textToDecrypt []byte) ([]byte, error)
}

func NewCipherWithConfig(cfg *Config, opts ...refx.Option) (Cipher, error) {
	var options CipherOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.WithMessage(err, "cfg.Unmarshal failed.")
	}
	return NewCipherWithOptions(&options, opts...)
}

func NewCipherWithOptions(options *CipherOptions, opts ...refx.Option) (Cipher, error) {
	constructor, ok := cipherConstructorMap[options.Type]
	if !ok {
		return nil, errors.Errorf("unsupported cipher type: [%v]", options.Type)
	}

	result, err := constructor.Call(options.Options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "constructor.Call failed")
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, errors.Wrapf(result[1].Interface().(error), "NewCipher failed. type: [%v]", options.Type)
		}
		return result[0].Interface().(Cipher), nil
	}

	return result[0].Interface().(Cipher), nil
}

type CipherOptions struct {
	Type    string
	Options interface{}
}

type EmptyCipher struct{}

func (c EmptyCipher) Encrypt(textToEncrypt []byte) ([]byte, error) {
	return textToEncrypt, nil
}

func (c EmptyCipher) Decrypt(textToDecrypt []byte) ([]byte, error) {
	return textToDecrypt, nil
}
