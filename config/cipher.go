package config

import (
	"encoding/base64"
)

type Cipher interface {
	Encrypt(textToEncrypt []byte) ([]byte, error)
	Decrypt(textToDecrypt []byte) ([]byte, error)
}

func NewCipherWithConfig(conf *Config) (Cipher, error) {
	switch conf.GetString("Type") {
	case "AES":
		buf, err := base64.StdEncoding.DecodeString(conf.GetString("Key"))
		if err != nil {
			return nil, err
		}
		return NewAESCipher(buf)
	case "AESWithKMSKey":
		return NewAESWithKMSKeyCipherWithAccessKey(
			conf.GetString("AccessKeyID"),
			conf.GetString("AccessKeySecret"),
			conf.GetString("regionID"),
			conf.GetString("Key"),
		)
	case "KMS":
		return NewKMSCipherWithAccessKey(
			conf.GetString("AccessKeyID"),
			conf.GetString("AccessKeySecret"),
			conf.GetString("regionID"),
			conf.GetString("KeyID"),
		)
	case "Base64":
		return NewBase64Cipher(), nil
	case "Group":
		subs, err := conf.SubArr("Ciphers")
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
