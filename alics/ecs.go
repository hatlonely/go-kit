package alics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// @see https://help.aliyun.com/document_detail/108460.html

func ECSMetaDataRegionID() (string, error) {
	return httpGet("http://100.100.100.200/latest/meta-data/region-id")
}

func ECSMetaDataRamSecurityCredentialsRole() (string, error) {
	return httpGet("http://100.100.100.200/latest/meta-data/Ram/security-credentials/")
}

func ECSUserData() (string, error) {
	return httpGet("http://100.100.100.200/latest/user-data")
}

func ECSMetaDataInstanceID() (string, error) {
	return httpGet("http://100.100.100.200/latest/meta-data/instance-id")
}

func ECSMetaDataImageID() (string, error) {
	return httpGet("http://100.100.100.200/latest/meta-data/image-id")
}

func ECSMaintenanceActiveSystemEvents() (string, error) {
	return httpGet("http://100.100.100.200/latest/maintenance/active-system-events")
}

func ECSMetaDataHostname() (string, error) {
	return httpGet("http://100.100.100.200/latest/meta-data/hostname")
}

func ECSMetaDataVPCID() (string, error) {
	return httpGet("http://100.100.100.200/latest/meta-data/vpc-id")
}

func ECSMetaDataZoneID() (string, error) {
	return httpGet("http://100.100.100.200/latest/meta-data/zone-id")
}

type ECSDynamicInstanceIdentityDocumentRes struct {
	ZoneID         string `json:"zone-id,omitempty"`
	SerialNumber   string `json:"serial-number,omitempty"`
	InstanceID     string `json:"instance-id,omitempty"`
	RegionID       string `json:"region-id,omitempty"`
	PrivateIPv4    string `json:"private-ipv4,omitempty"`
	OwnerAccountID string `json:"owner-account-id,omitempty"`
	Mac            string `json:"mac,omitempty"`
	ImageID        string `json:"image-id,omitempty"`
	InstanceType   string `json:"instance-type,omitempty"`
}

func ECSDynamicInstanceIdentityDocument() (*ECSDynamicInstanceIdentityDocumentRes, error) {
	var res ECSDynamicInstanceIdentityDocumentRes
	if err := httpGetJson("http://100.100.100.200/latest/dynamic/instance-identity/document", &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type ECSMetaDataRamSecurityCredentialsRes struct {
	AccessKeyID     string `json:"AccessKeyId,omitempty"`
	AccessKeySecret string
	Expiration      string
	SecurityToken   string
	LastUpdated     string
	Code            string
}

// @see https://help.aliyun.com/document_detail/127171.html
func ECSMetaDataRamSecurityCredentials() (*ECSMetaDataRamSecurityCredentialsRes, error) {
	role, err := ECSMetaDataRamSecurityCredentialsRole()
	if err != nil {
		return nil, err
	}

	var res ECSMetaDataRamSecurityCredentialsRes
	if err := httpGetJson(
		fmt.Sprintf("http://100.100.100.200/latest/meta-data/Ram/security-credentials/%s", role),
		&res,
	); err != nil {
		return nil, err
	}

	return &res, nil
}

func httpGetJson(url string, v interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(buf, v); err != nil {
		return err
	}

	return nil
}

func httpGet(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
