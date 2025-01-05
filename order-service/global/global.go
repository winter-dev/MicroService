package global

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"gorm.io/gorm"
)

var (
	GLOBAL_DB   *gorm.DB
	GLOBAL_ETCD *clientv3.Client
)
