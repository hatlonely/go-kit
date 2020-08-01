package config

import (
	"encoding/base64"
)

type Cipher interface {
	Encrypt(textToEncrypt []byte) ([]byte, error)
	Decrypt(textToDecrypt []byte) ([]byte, error)
}

func NewCipherWithConfig(conf *Config) (Cipher, error) {
	if conf == nil {
		return nil, nil
	}
	switch conf.GetString("type") {
	case "AES":
		buf, err := base64.StdEncoding.DecodeString(conf.GetString("key"))
		if err != nil {
			return nil, err
		}
		return NewAESCipher(buf)
	case "AESWithKMSKey":
		return NewAESWithKMSKeyCipherWithAccessKey(
			conf.GetString("accessKeyID"),
			conf.GetString("accessKeySecret"),
			conf.GetString("regionID"),
			conf.GetString("key"),
		)
	case "KMS":
		return NewKMSCipherWithAccessKey(
			conf.GetString("accessKeyID"),
			conf.GetString("accessKeySecret"),
			conf.GetString("regionID"),
			conf.GetString("keyID"),
		)
	case "Base64":
		return NewBase64Cipher(), nil
	case "Group":
		subs, err := conf.SubArr("ciphers")
		if err != nil {
			return nil, err
		}
		var ciphers []Cipher
		for _, sub := range subs {
			cipher, err := NewCipherWithConfig(sub)
			if err != nil {
				return nil, err
			}
			ciphers = append(ciphers, cipher)
		}
		return NewCipherGroup(ciphers...), nil
	}
	return nil, nil
}
