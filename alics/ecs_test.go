package alics

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestECSMetaDataRegionID(t *testing.T) {
	Convey("TestECSMetaDataRegionID", t, func() {
		patches := ApplyFunc(http.Get, func(url string) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`cn-zhangjiakou`))),
			}, nil
		})
		defer patches.Reset()

		res, err := ECSMetaDataRegionID()
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "cn-zhangjiakou")
	})
}

func TestECSMetaDataInstanceID(t *testing.T) {
	Convey("TestECSMetaDataInstanceID", t, func() {
		patches := ApplyFunc(http.Get, func(url string) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`i-8vb8p62ryuejka7tcz4p`))),
			}, nil
		})
		defer patches.Reset()

		res, err := ECSMetaDataInstanceID()
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "i-8vb8p62ryuejka7tcz4p")
	})
}

func TestECSMetaDataRamSecurityCredentialsRole(t *testing.T) {
	Convey("TestECSMetaDataRamSecurityCredentialsRole", t, func() {
		patches := ApplyFunc(http.Get, func(url string) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`KubernetesEcsRamRole`))),
			}, nil
		})
		defer patches.Reset()

		res, err := ECSMetaDataRamSecurityCredentialsRole()
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "KubernetesEcsRamRole")
	})
}

func TestECSUserData(t *testing.T) {
	Convey("TestECSUserData", t, func() {
		patches := ApplyFunc(http.Get, func(url string) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`test user data`))),
			}, nil
		})
		defer patches.Reset()

		res, err := ECSUserData()
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "test user data")
	})
}

func TestECSMetaDataImageID(t *testing.T) {
	Convey("TestECSMetaDataImageID", t, func() {
		patches := ApplyFunc(http.Get, func(url string) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`centos_7_8_x64_20G_alibase_20200914.vhd`))),
			}, nil
		})
		defer patches.Reset()

		res, err := ECSMetaDataImageID()
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "centos_7_8_x64_20G_alibase_20200914.vhd")
	})
}

func TestECSMaintenanceActiveSystemEvents(t *testing.T) {
	Convey("TestECSMaintenanceActiveSystemEvents", t, func() {
		patches := ApplyFunc(http.Get, func(url string) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`test event`))),
			}, nil
		})
		defer patches.Reset()

		res, err := ECSMaintenanceActiveSystemEvents()
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "test event")
	})
}

func TestECSMetaDataHostname(t *testing.T) {
	Convey("TestECSMetaDataHostname", t, func() {
		patches := ApplyFunc(http.Get, func(url string) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`imm-dev-hl`))),
			}, nil
		})
		defer patches.Reset()

		res, err := ECSMetaDataHostname()
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "imm-dev-hl")
	})
}

func TestECSMetaDataVPCID(t *testing.T) {
	Convey("TestECSMetaDataVPCID", t, func() {
		patches := ApplyFunc(http.Get, func(url string) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`vpc-8vbmtwbak0pi0ef8l7q0o`))),
			}, nil
		})
		defer patches.Reset()

		res, err := ECSMetaDataVPCID()
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "vpc-8vbmtwbak0pi0ef8l7q0o")
	})
}

func TestECSMetaDataZoneID(t *testing.T) {
	Convey("TestECSMetaDataZoneID", t, func() {
		patches := ApplyFunc(http.Get, func(url string) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`cn-zhangjiakou-a`))),
			}, nil
		})
		defer patches.Reset()

		res, err := ECSMetaDataZoneID()
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "cn-zhangjiakou-a")
	})
}

func TestECSMetaDataRamSecurityCredentials(t *testing.T) {
	Convey("TestECSMetaDataRamSecurityCredentials", t, func() {
		patches := ApplyFunc(http.Get, func(url string) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{
  "AccessKeyId" : "test-ak",
  "AccessKeySecret" : "test-sk",
  "Expiration" : "2020-11-27T20:08:18Z",
  "SecurityToken" : "test-token",
  "LastUpdated" : "2020-11-27T14:08:18Z",
  "Code" : "Success"
}`))),
			}, nil
		})
		defer patches.Reset()

		res, err := ECSMetaDataRamSecurityCredentials()
		So(err, ShouldBeNil)
		So(res.AccessKeyID, ShouldEqual, "test-ak")
		So(res.AccessKeySecret, ShouldEqual, "test-sk")
		So(res.Expiration, ShouldEqual, "2020-11-27T20:08:18Z")
		So(res.SecurityToken, ShouldEqual, "test-token")
		So(res.LastUpdated, ShouldEqual, "2020-11-27T14:08:18Z")
		So(res.Code, ShouldEqual, "Success")
	})
}

func TestECSDynamicInstanceIdentityDocument(t *testing.T) {
	Convey("TestECSDynamicInstanceIdentityDocument", t, func() {
		patches := ApplyFunc(http.Get, func(url string) (*http.Response, error) {
			So(url, ShouldEqual, "http://100.100.100.200/latest/dynamic/instance-identity/document")
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"zone-id":"cn-zhangjiakou-a","serial-number":"8e794d59-a809-453b-8f5a-9d0a0b29660a","instance-id":"i-8vb8p62ryuejka7tcz4p","region-id":"cn-zhangjiakou","private-ipv4":"192.168.1.215","owner-account-id":"1335926999564873","mac":"00:16:3e:13:9a:81","image-id":"centos_7_8_x64_20G_alibase_20200914.vhd","instance-type":"ecs.s6-c1m1.small"}`))),
			}, nil
		})
		defer patches.Reset()

		res, err := ECSDynamicInstanceIdentityDocument()
		So(err, ShouldBeNil)
		So(res.ZoneID, ShouldEqual, "cn-zhangjiakou-a")
		So(res.SerialNumber, ShouldEqual, "8e794d59-a809-453b-8f5a-9d0a0b29660a")
		So(res.InstanceID, ShouldEqual, "i-8vb8p62ryuejka7tcz4p")
		So(res.RegionID, ShouldEqual, "cn-zhangjiakou")
		So(res.PrivateIPv4, ShouldEqual, "192.168.1.215")
		So(res.OwnerAccountID, ShouldEqual, "1335926999564873")
		So(res.Mac, ShouldEqual, "00:16:3e:13:9a:81")
		So(res.ImageID, ShouldEqual, "centos_7_8_x64_20G_alibase_20200914.vhd")
		So(res.InstanceType, ShouldEqual, "ecs.s6-c1m1.small")
	})
}
