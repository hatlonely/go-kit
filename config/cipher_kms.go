package config

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"

	"github.com/hatlonely/go-kit/alics"
)

func NewKMSCipherWithOptions(options *KMSCipherOptions) (*KMSCipher, error) {
	return NewKMSCipherWithAccessKey(options.AccessKeyID, options.AccessKeySecret, options.RegionID, options.KeyID)
}

func NewKMSCipherWithAccessKey(ak, sk, regionID string, keyID string) (*KMSCipher, error) {
	var kmsCli *kms.Client
	var err error
	if regionID == "" {
		if regionID, err = alics.ECSMetaDataRegionID(); err != nil {
			return nil, err
		}
	}
	if ak == "" {
		role, err := alics.ECSMetaDataRamSecurityCredentialsRole()
		if err != nil {
			return nil, err
		}
		if kmsCli, err = kms.NewClientWithEcsRamRole(regionID, role); err != nil {
			return nil, err
		}
	} else {
		if kmsCli, err = kms.NewClientWithAccessKey(ak, sk, regionID); err != nil {
			return nil, err
		}
	}
	return NewKMSCipher(kmsCli, keyID)
}

func NewKMSCipher(kmsCli *kms.Client, keyID string) (*KMSCipher, error) {
	return &KMSCipher{
		kmsCli: kmsCli,
		keyID:  keyID,
	}, nil
}

type KMSCipher struct {
	kmsCli *kms.Client
	keyID  string
}

func (c *KMSCipher) Encrypt(textToEncrypt []byte) ([]byte, error) {
	req := kms.CreateEncryptRequest()
	req.Scheme = "https"
	req.Plaintext = string(textToEncrypt)
	req.KeyId = c.keyID
	res, err := c.kmsCli.Encrypt(req)
	if err != nil {
		return nil, err
	}

	return []byte(res.CiphertextBlob), nil
}

func (c *KMSCipher) Decrypt(textToDecrypt []byte) ([]byte, error) {
	req := kms.CreateDecryptRequest()
	req.Scheme = "https"
	req.CiphertextBlob = string(textToDecrypt)
	res, err := c.kmsCli.Decrypt(req)
	if err != nil {
		return nil, err
	}

	return []byte(res.Plaintext), nil
}

func (c *KMSCipher) GenerateDataKey() (string, string, error) {
	req := kms.CreateGenerateDataKeyRequest()
	req.Scheme = "https"
	req.KeyId = c.keyID
	res, err := c.kmsCli.GenerateDataKey(req)
	if err != nil {
		return "", "", nil
	}
	return res.Plaintext, res.CiphertextBlob, nil
}

type KMSCipherOptions struct {
	AccessKeyID     string
	AccessKeySecret string
	RegionID        string
	KeyID           string
}
