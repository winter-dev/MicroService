package service

import (
	"fmt"
	"log"
	"ms-common/cli"
	"user-service/global"
)

func Watch() {
	// 创建通知 channel
	notifyCh := make(chan struct{}, 10)

	etcd := cli.Etcd{}
	go func() {
		for {
			select {
			case <-notifyCh:
				var err error
				services, err := etcd.ServiceList(global.GLOBAL_ETCD)
				if err != nil {
					log.Printf("failed to update service list: %v", err)
					return
				}
				fmt.Println("Service changed, update service list")

				global.Lock.Lock()
				global.GLOBAL_SERVICES = services
				global.Lock.Unlock()
			}
		}
	}()
	etcd.WatchService(global.GLOBAL_ETCD, notifyCh)
}
