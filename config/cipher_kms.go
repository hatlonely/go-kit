package config

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/alics"
)

func init() {
	RegisterCipher("KMS", NewKMSCipherWithOptions)
}

func NewKMSCipherWithOptions(options *KMSCipherOptions) (*KMSCipher, error) {
	var err error
	if options.RegionID == "" {
		if options.RegionID, err = alics.ECSMetaDataRegionID(); err != nil {
			return nil, errors.Wrap(err, "alics.ECSMetaDataRegionID failed")
		}
	}
	var client *kms.Client
	if options.AccessKeyID == "" {
		role, err := alics.ECSMetaDataRamSecurityCredentialsRole()
		if err != nil {
			return nil, errors.Wrap(err, "alics.ECSMetaDataRamSecurityCredentialsRole failed")
		}
		client, err = kms.NewClientWithEcsRamRole(options.RegionID, role)
		if err != nil {
			return nil, errors.Wrap(err, "kms.NewClientWithEcsRamRole failed")
		}
	} else {
		client, err = kms.NewClientWithAccessKey(options.RegionID, options.AccessKeyID, options.AccessKeySecret)
		if err != nil {
			return nil, errors.Wrap(err, "kms.NewClientWithAccessKey failed")
		}
	}

	return NewKMSCipher(client, options.KeyID)
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
		return nil, errors.Wrap(err, "kmsCli.Encrypt failed")
	}

	return []byte(res.CiphertextBlob), nil
}

func (c *KMSCipher) Decrypt(textToDecrypt []byte) ([]byte, error) {
	req := kms.CreateDecryptRequest()
	req.Scheme = "https"
	req.CiphertextBlob = string(textToDecrypt)
	res, err := c.kmsCli.Decrypt(req)
	if err != nil {
		return nil, errors.Wrap(err, "kmsCli.Decrypt failed")
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
