package config

import (
	"encoding/base64"
)

type Base64Cipher struct {
	encoding *base64.Encoding
}

func NewBase64Cipher() Base64Cipher {
	return NewBase64CipherWithOptions(&Base64CipherOptions{})
}

func NewBase64CipherWithOptions(options *Base64CipherOptions) Base64Cipher {
	if options.URLEncoding {
		if options.NoPadding {
			return Base64Cipher{encoding: base64.URLEncoding.WithPadding(base64.NoPadding)}
		}
		if options.Padding != "" {
			return Base64Cipher{encoding: base64.URLEncoding.WithPadding(rune(options.Padding[0]))}
		}
		return Base64Cipher{encoding: base64.URLEncoding}
	}
	if options.StdEncoding || options.Encoding == "" {
		if options.NoPadding {
			return Base64Cipher{encoding: base64.StdEncoding.WithPadding(base64.NoPadding)}
		}
		if options.Padding != "" {
			return Base64Cipher{encoding: base64.StdEncoding.WithPadding(rune(options.Padding[0]))}
		}
		return Base64Cipher{encoding: base64.StdEncoding}
	}

	if options.NoPadding {
		return Base64Cipher{encoding: base64.NewEncoding(options.Encoding).WithPadding(base64.NoPadding)}
	}
	if options.Padding != "" {
		return Base64Cipher{encoding: base64.NewEncoding(options.Encoding).WithPadding(rune(options.Padding[0]))}
	}
	return Base64Cipher{encoding: base64.NewEncoding(options.Encoding)}
}

func (c Base64Cipher) Encrypt(textToEncrypt []byte) ([]byte, error) {
	return []byte(c.encoding.EncodeToString(textToEncrypt)), nil
}

func (c Base64Cipher) Decrypt(textToDecrypt []byte) ([]byte, error) {
	return c.encoding.DecodeString(string(textToDecrypt))
}

type Base64CipherOptions struct {
	Encoding    string
	Padding     string
	NoPadding   bool
	StdEncoding bool
	URLEncoding bool
}
