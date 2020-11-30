package config

import (
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

type Cipher interface {
	Encrypt(textToEncrypt []byte) ([]byte, error)
	Decrypt(textToDecrypt []byte) ([]byte, error)
}

func NewCipherWithConfig(cfg *Config, opts ...refx.Option) (Cipher, error) {
	var options CipherOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "cfg.Unmarshal failed.")
	}
	return NewCipherWithOptions(&options)
}

func NewCipherWithOptions(options *CipherOptions) (Cipher, error) {
	switch options.Type {
	case "", "Empty":
		return EmptyCipher{}, nil
	case "AES":
		return NewAESCipherWithOptions(&options.AESCipher)
	case "Base64":
		return NewBase64CipherWithOptions(&options.Base64Cipher), nil
	case "KMS":
		return NewKMSCipherWithOptions(&options.KMSCipher)
	case "Group":
		return NewCipherGroupWithOptions(&options.CipherGroup)
	}
	return nil, errors.Errorf("unsupported cipher type [%v]", options.Type)
}

type CipherOptions struct {
	Type         string `dft:"Empty" rule:"x in ['Empty', 'AES', 'Base64', 'KMS', 'Group']"`
	AESCipher    AESCipherOptions
	KMSCipher    KMSCipherOptions
	CipherGroup  CipherGroupOptions
	Base64Cipher Base64CipherOptions
}

type EmptyCipher struct{}

func (c EmptyCipher) Encrypt(textToEncrypt []byte) ([]byte, error) {
	return textToEncrypt, nil
}

func (c EmptyCipher) Decrypt(textToDecrypt []byte) ([]byte, error) {
	return textToDecrypt, nil
}
