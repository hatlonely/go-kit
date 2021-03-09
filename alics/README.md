# alics

阿里云服务的一些辅助函数集合


## Quick Start

```go
func ECSMetaDataRegionID() (string, error)
func ECSMetaDataRamSecurityCredentialsRole() (string, error)
func ECSUserData() (string, error)
func ECSMetaDataInstanceID() (string, error)
func ECSMetaDataImageID() (string, error)
func ECSMaintenanceActiveSystemEvents() (string, error)
func ECSMetaDataHostname() (string, error)
func ECSMetaDataVPCID() (string, error)
func ECSMetaDataZoneID() (string, error)
func ECSDynamicInstanceIdentityDocument() (*ECSDynamicInstanceIdentityDocumentRes, error)
func ECSMetaDataRamSecurityCredentials() (*ECSMetaDataRamSecurityCredentialsRes, error)

func ParseOSSUri(uri string) (*OSSUriInfo, error)
func ParsePDSUri(uri string) (*PDSUriInfo, error)
```