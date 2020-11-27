package config

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/pkg/errors"
)

const AESMaxKeyLen = 32

func NewAESWithKMSKeyCipherWithAccessKey(accessKeyID string, accessKeySecret string, regionID string, cipherText string) (*AESCipher, error) {
	kmsCli, err := newKMSClient(accessKeyID, accessKeySecret, regionID)
	if err != nil {
		return nil, err
	}
	return NewAESWithKMSKeyCipher(kmsCli, cipherText)
}

func NewAESWithKMSKeyCipher(kmsCli *kms.Client, cipherText string) (*AESCipher, error) {
	req := kms.CreateDecryptRequest()
	req.Scheme = "https"
	req.CiphertextBlob = cipherText
	res, err := kmsCli.Decrypt(req)
	if err != nil {
		return nil, err
	}

	buf, err := base64.StdEncoding.DecodeString(res.Plaintext)
	if err != nil {
		return nil, err
	}

	return NewAESCipher(buf)
}

func NewAESCipher(key []byte) (*AESCipher, error) {
	if len(key) > AESMaxKeyLen {
		return nil, fmt.Errorf("key len should less or equal to %v", AESMaxKeyLen)
	}
	buf := make([]byte, AESMaxKeyLen)
	copy(buf, key)
	block, err := aes.NewCipher(buf)
	if err != nil {
		return nil, err
	}

	return &AESCipher{
		cb: block,
	}, nil
}

func NewAESCipherWithOptions(options *AESCipherOptions) (*AESCipher, error) {
	if options.Key != "" {
		return NewAESCipher([]byte(options.Key))
	}
	if options.Base64Key != "" {
		buf, err := base64.StdEncoding.DecodeString(options.Base64Key)
		if err != nil {
			return nil, errors.Wrap(err, "base64 decode key failed.")
		}
		return NewAESCipher(buf)
	}
	if options.KMSKey != "" {
		return NewAESWithKMSKeyCipherWithAccessKey(options.KMS.AccessKeyID, options.KMS.AccessKeySecret, options.KMS.RegionID, options.KMSKey)
	}
	return nil, errors.New("no key found")
}

type AESCipher struct {
	cb cipher.Block
}

func (c *AESCipher) Encrypt(textToEncrypt []byte) ([]byte, error) {
	plainText := make([]byte, len(textToEncrypt)+aes.BlockSize+(-len(textToEncrypt)%aes.BlockSize))
	copy(plainText, textToEncrypt)
	encryptText := make([]byte, aes.BlockSize+len(plainText))
	iv := encryptText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(c.cb, iv)
	mode.CryptBlocks(encryptText[aes.BlockSize:], plainText)

	return encryptText, nil
}

func (c *AESCipher) Decrypt(textToDecrypt []byte) ([]byte, error) {
	if len(textToDecrypt) < aes.BlockSize {
		return nil, errors.New("cipher text too short")
	}

	iv := textToDecrypt[:aes.BlockSize]
	textToDecrypt = textToDecrypt[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(c.cb, iv)
	mode.CryptBlocks(textToDecrypt, textToDecrypt)

	return bytes.TrimRight(textToDecrypt, "\x00"), nil
}

type AESCipherOptions struct {
	Key       string
	Base64Key string
	KMSKey    string
	KMS       struct {
		AccessKeyID     string
		AccessKeySecret string
		RegionID        string
	}
}
