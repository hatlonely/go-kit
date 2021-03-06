package wrap

import (
	_ "github.com/alibabacloud-go/pds-sdk/client"
	_ "github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	_ "github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	_ "github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	_ "github.com/aliyun/aliyun-mns-go-sdk"
	_ "github.com/aliyun/aliyun-oss-go-sdk/oss"
	_ "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	_ "github.com/go-redis/redis"
	_ "github.com/jinzhu/gorm"
	_ "github.com/nacos-group/nacos-sdk-go/clients/config_client"
	_ "github.com/olivere/elastic/v7"
	_ "go.mongodb.org/mongo-driver/mongo"
)
