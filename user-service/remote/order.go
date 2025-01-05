package remote

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"ms-common/common"
	"ms-common/pb/pb"
	"time"
	"user-service/global"
)

func Call(svc string) []*pb.Order {
	var service common.ServiceName
	var gl global.Gl
	si := gl.GetSvcMeta(service.GetService(svc))
	if si.Name == "" {
		log.Println("service not online")
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(si.GetGrpcURI(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	client := pb.NewOrderServiceClient(conn)

	// 调用 gRPC 方法
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.OrderList(ctx, &pb.UserOrderRequest{UserId: 1})
	if err != nil {
		log.Fatalf("failed to call OrderList: %v", err)
	}
	return resp.Orders
}
