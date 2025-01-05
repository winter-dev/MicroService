package biz

import (
	"context"
	"errors"
	"log"
	"ms-common/pb/pb"
	"order-service/global"
	"order-service/model"
	"time"
)

type OrderBiz struct {
	pb.UnimplementedOrderServiceServer
}

func (o OrderBiz) OrderList(_ context.Context, uor *pb.UserOrderRequest) (*pb.OrderListResponse, error) {
	if uor.UserId == 0 {
		return &pb.OrderListResponse{Response: &pb.BaseResponse{Code: 0, Message: "fail"},
			Orders: nil,
		}, errors.New("userId is null")
	}
	var orders []model.Order
	if err := global.GLOBAL_DB.Where("user_id=?", uor.UserId).Find(&orders).Error; err != nil {
		log.Println(err)
		return &pb.OrderListResponse{Response: &pb.BaseResponse{Code: 0, Message: "fail"},
			Orders: nil,
		}, err
	}

	if len(orders) == 0 {
		return &pb.OrderListResponse{Response: &pb.BaseResponse{Code: 0, Message: "fail"},
			Orders: nil,
		}, errors.New("order not found")
	}

	orderList := make([]*pb.Order, len(orders))
	for i, v := range orders {
		orderList[i] = &pb.Order{
			Id:          v.Id,
			UserId:      v.UserId,
			Name:        v.Name,
			Price:       v.Price,
			Description: v.Description,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}
	return &pb.OrderListResponse{Response: &pb.BaseResponse{Code: 200, Message: "success"},
		Orders: orderList,
	}, nil

}

func (OrderBiz) GetOrder(_ context.Context, or *pb.OrderRequest) (*pb.OrderResponse, error) {
	var order model.Order
	if err := global.GLOBAL_DB.Where("id=?", or.Id).First(&order).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	if order.Id == 0 {
		return nil, errors.New("order not found")
	}

	return &pb.OrderResponse{Base: &pb.BaseResponse{Code: 200, Message: "success"}, Order: &pb.Order{
		Id:          order.Id,
		UserId:      order.UserId,
		Name:        order.Name,
		Price:       order.Price,
		Description: order.Description,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}}, nil
}

func (OrderBiz) DeleteOrder(_ context.Context, order *pb.OrderRequest) (*pb.BaseResponse, error) {
	global.GLOBAL_DB.Delete(&model.Order{}, order.Id)
	return &pb.BaseResponse{Code: 200, Message: "success"}, nil
}

func (OrderBiz) CreateOrder(_ context.Context, o *pb.Order) (*pb.BaseResponse, error) {
	var order model.Order
	order.Id = o.Id
	order.UserId = o.UserId
	order.Name = o.Name
	order.Price = o.Price
	order.Description = o.Description
	order.CreatedAt = time.Now().Unix()
	global.GLOBAL_DB.Create(&order)
	return &pb.BaseResponse{Code: 200, Message: "success"}, nil
}
