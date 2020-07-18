package config

import (
	"encoding/base64"
)

type Base64Cipher struct{}

func NewBase64Cipher() Base64Cipher {
	return Base64Cipher{}
}

func (c Base64Cipher) Encrypt(textToEncrypt []byte) ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(textToEncrypt)), nil
}

func (c Base64Cipher) Decrypt(textToDecrypt []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(textToDecrypt))
}
