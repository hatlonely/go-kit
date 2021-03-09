# alics

阿里云服务的一些辅助函数集合


## ECS

- `ECSMetaDataRegionID`:获取 region id
- `ECSMetaDataRamSecurityCredentialsRole`:获取 ecs role 角色名
- `ECSUserData`:获取用户自定义数据
- `ECSMetaDataInstanceID`:获取 ecs instance id
- `ECSMetaDataImageID`:获取镜像 id
- `ECSMaintenanceActiveSystemEvents`:获取系统事件
- `ECSMetaDataHostname`:获取 hostname
- `ECSMetaDataVPCID`:获取 vpc id
- `ECSMetaDataZoneID`:获取 zone id
- `ECSDynamicInstanceIdentityDocument`:获取机器实例信息
- `ECSMetaDataRamSecurityCredentials`:获取临时的授权信息

## OSS

- `ParseOSSUri`:解析 oss uri，获取 bucket 和 object

## PDS

- `ParsePDSUri`:解析 pds uri