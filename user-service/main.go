package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"ms-common/cli"
	"ms-common/config"
	"ms-common/pb/pb"
	"net"
	"net/http"
	"time"
	"user-service/global"
	"user-service/internal/biz"
	"user-service/model"
	"user-service/router"
	"user-service/service"
)

func Init() {
	config.Init()
	var err error

	global.GLOBAL_ETCD, err = clientv3.New(clientv3.Config{
		Endpoints:   config.G_AppConfig.AppConfig.EcdAddress,
		DialTimeout: 5. * time.Second,
	})

	global.GLOBAL_DB, err = gorm.Open(sqlite.Open("./user.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect to database")
	}
	global.GLOBAL_DB.AutoMigrate(&model.User{})
}

func grpcServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.G_AppConfig.AppConfig.GrpcPort))
	if err != nil {
		log.Fatal("failed to listen:", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &biz.UserBiz{})
	log.Println("grpc server start on port:", config.G_AppConfig.AppConfig.GrpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}

func main() {
	Init()

	go func() {
		etcd := cli.Etcd{}
		etcd.Register(global.GLOBAL_ETCD)
	}()

	go grpcServer()
	go service.Watch()
	go func() {
		etcd := cli.Etcd{}
		global.GLOBAL_SERVICES, _ = etcd.ServiceList(global.GLOBAL_ETCD)
	}()

	gin.SetMode("debug")
	app := gin.Default()
	app.Use(gin.Recovery())

	app.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, "ok")
	})
	app.GET("/orderList/:userId", router.OrderList)

	app.GET("/services", func(context *gin.Context) {
		context.JSON(http.StatusOK, global.GLOBAL_SERVICES)
	})
	app.Run(fmt.Sprintf(":%d", config.G_AppConfig.AppConfig.HttpPort))
}
