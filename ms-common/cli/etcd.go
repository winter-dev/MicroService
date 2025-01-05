package cli

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"ms-common/config"
	"ms-common/types"
)

type Etcd struct {
}

var prefix = "/services/"

func (*Etcd) Register(cli *clientv3.Client) {
	key := prefix + config.G_AppConfig.AppConfig.AppName + "/" + config.G_AppConfig.AppConfig.Address

	si := types.ServiceInfo{
		Name:    config.G_AppConfig.AppConfig.AppName,
		Address: config.G_AppConfig.AppConfig.Address,
		Meta: map[string]any{
			"env":     config.G_AppConfig.AppConfig.Env,
			"version": config.G_AppConfig.AppConfig.Version,
			"port":    config.G_AppConfig.AppConfig.GrpcPort,
		},
	}
	// 序列化为 JSON
	siData, err := json.Marshal(si)
	if err != nil {
		log.Printf("failed to marshal service info: %v", err)
		return
	}

	// 向 etcd 注册服务
	lease := clientv3.NewLease(cli)

	leaseRes, err := cli.Grant(context.Background(), 10)
	if err != nil {
		log.Printf("failed to create lease: %v\n", err)
		return
	}
	_, err = cli.Put(context.Background(), key, string(siData), clientv3.WithLease(leaseRes.ID))
	if err != nil {
		log.Printf("failed to register service: %v\n", err)
		return
	}

	// 保持心跳
	keepAlive, err := lease.KeepAlive(context.Background(), leaseRes.ID)
	if err != nil {
		log.Fatalf("failed to keep alive: %v", err)
		return
	}
	go func() {
		for {
			<-keepAlive
		}
	}()
}

func (*Etcd) ServiceList(cli *clientv3.Client) ([]types.ServiceInfo, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 获取所有服务
	rsp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		log.Printf("failed to get services: %v", err)
		return nil, err
	}
	var sis []types.ServiceInfo
	for _, kv := range rsp.Kvs {
		var si types.ServiceInfo
		err = json.Unmarshal(kv.Value, &si)
		if err != nil {
			log.Printf("failed to unmarshal service info: %v", err)
			continue
		}
		sis = append(sis, si)
	}
	return sis, nil
}

func (*Etcd) WatchService(cli *clientv3.Client, notifyCh chan<- struct{}) {
	wtChan := cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for {
		select {
		case resp, ok := <-wtChan:
			if !ok {
				log.Println("Watch channel closed, reconnecting...")
				wtChan = cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
				continue
			}
			for _, ev := range resp.Events {
				fmt.Println(string(ev.Kv.Key))
				switch ev.Type {
				case clientv3.EventTypeDelete:
					fmt.Println("service down")
					notifyCh <- struct{}{}
				case clientv3.EventTypePut:
					fmt.Println("service up")
					notifyCh <- struct{}{}
				}
			}
		}
	}
}
