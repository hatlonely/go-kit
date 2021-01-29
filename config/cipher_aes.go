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

	"github.com/hatlonely/go-kit/alics"
)

func init() {
	RegisterCipher("AES", NewAESCipherWithOptions)
}

const AESMaxKeyLen = 32

func NewAESCipherWithOptions(options *AESCipherOptions) (*AESCipher, error) {
	if options.Key != "" {
		return NewAESCipher([]byte(options.Key), options.PaddingType)
	}
	if options.Base64Key != "" {
		buf, err := base64.StdEncoding.DecodeString(options.Base64Key)
		if err != nil {
			return nil, errors.Wrap(err, "base64.StdEncoding.DecodeString failed.")
		}
		return NewAESCipher(buf, options.PaddingType)
	}
	if options.KMSKey != "" {
		return NewAESWithKMSKeyCipherWithAccessKey(
			options.KMS.AccessKeyID, options.KMS.AccessKeySecret,
			options.KMS.RegionID, options.KMSKey, options.PaddingType,
		)
	}
	return nil, errors.New("no key found")
}

func NewAESWithKMSKeyCipherWithAccessKey(accessKeyID string, accessKeySecret string, regionID string, cipherText string, paddingType string) (*AESCipher, error) {
	var client *kms.Client
	var err error
	if accessKeyID == "" {
		role, err := alics.ECSMetaDataRamSecurityCredentialsRole()
		if err != nil {
			return nil, errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentialsRole failed")
		}
		client, err = kms.NewClientWithEcsRamRole(regionID, role)
		if err != nil {
			return nil, errors.Wrap(err, "kms.NewClientWithEcsRamRole failed")
		}
	} else {
		client, err = kms.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
		if err != nil {
			return nil, errors.Wrap(err, "kms.NewClientWithAccessKey failed")
		}
	}
	return NewAESWithKMSKeyCipher(client, cipherText, paddingType)
}

func NewAESWithKMSKeyCipher(kmsCli *kms.Client, cipherText string, paddingType string) (*AESCipher, error) {
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

	return NewAESCipher(buf, paddingType)
}

func NewAESCipher(key []byte, paddingType string) (*AESCipher, error) {
	if len(key) > AESMaxKeyLen {
		return nil, fmt.Errorf("key len should less or equal to %v", AESMaxKeyLen)
	}
	buf := make([]byte, AESMaxKeyLen)
	copy(buf, key)
	block, err := aes.NewCipher(buf)
	if err != nil {
		return nil, err
	}

	if paddingType == "Zero" {
		return &AESCipher{
			cb:        block,
			padding:   zeroPadding,
			unPadding: zeroUnPadding,
		}, nil
	}
	return &AESCipher{
		cb:        block,
		padding:   pkcs7Padding,
		unPadding: pkcs7UnPadding,
	}, nil
}

type AESCipher struct {
	cb cipher.Block

	padding   func(textToPadding []byte, blockSize int) []byte
	unPadding func(textToUnPadding []byte) []byte
}

func (c *AESCipher) Encrypt(textToEncrypt []byte) ([]byte, error) {
	plainText := c.padding(textToEncrypt, aes.BlockSize)
	encryptText := make([]byte, aes.BlockSize+len(plainText))
	iv := encryptText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, errors.Wrap(err, "io.ReadFull failed")
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

	return c.unPadding(textToDecrypt), nil
}

type AESCipherOptions struct {
	Key         string
	Base64Key   string
	KMSKey      string
	PaddingType string
	KMS         struct {
		AccessKeyID     string
		AccessKeySecret string
		RegionID        string
	}
}

func pkcs7Padding(textToPadding []byte, blockSize int) []byte {
	n := blockSize - len(textToPadding)%blockSize
	return append(textToPadding, bytes.Repeat([]byte{byte(n)}, n)...)
}

func pkcs7UnPadding(textToUnPadding []byte) []byte {
	n := len(textToUnPadding)
	if n == 0 {
		return textToUnPadding
	}
	return textToUnPadding[:n-int(textToUnPadding[n-1])]
}

func zeroPadding(textToPadding []byte, blockSize int) []byte {
	n := len(textToPadding)
	if n%blockSize == 0 {
		return textToPadding
	}
	return append(textToPadding, bytes.Repeat([]byte{'\x00'}, blockSize-n%blockSize)...)
}

func zeroUnPadding(textToUnPadding []byte) []byte {
	return bytes.TrimRight(textToUnPadding, "\x00")
}
