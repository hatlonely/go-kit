package config

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"

	"github.com/hatlonely/go-kit/alics"
)

func NewKMSCipherWithOptions(options *KMSCipherOptions) (*KMSCipher, error) {
	return NewKMSCipherWithAccessKey(options.AccessKeyID, options.AccessKeySecret, options.RegionID, options.KeyID)
}

func newKMSClient(accessKeyID string, accessKeySecret string, regionID string) (kmsCli *kms.Client, err error) {
	if regionID == "" {
		if regionID, err = alics.ECSMetaDataRegionID(); err != nil {
			return nil, err
		}
	}
	if accessKeyID == "" {
		role, err := alics.ECSMetaDataRamSecurityCredentialsRole()
		if err != nil {
			return nil, err
		}
		return kms.NewClientWithEcsRamRole(regionID, role)
	}
	return kms.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
}

func NewKMSCipherWithAccessKey(accessKeyID string, accessKeySecret string, regionID string, keyID string) (*KMSCipher, error) {
	kmsCli, err := newKMSClient(accessKeyID, accessKeySecret, regionID)
	if err != nil {
		return nil, err
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
