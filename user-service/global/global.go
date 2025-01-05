package global

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"gorm.io/gorm"
	"ms-common/common"
	"ms-common/types"
	"sync"
)

var (
	GLOBAL_DB       *gorm.DB
	GLOBAL_ETCD     *clientv3.Client
	GLOBAL_SERVICES []types.ServiceInfo
	Lock            sync.RWMutex
)

type Gl struct {
}

func (Gl) GetSvcMeta(sn common.ServiceName) types.ServiceInfo {
	if sn == "" {
		return types.ServiceInfo{}
	}
	snName := sn.String()
	for _, v := range GLOBAL_SERVICES {
		if v.Name == snName {
			return v
		}
	}
	return types.ServiceInfo{}
}
