package biz

import (
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"ms-common/pb/pb"
	"time"
	"user-service/global"
	"user-service/model"
)

type UserBiz struct {
	pb.UnimplementedUserServiceServer
}

/*
	GetUser(context.Context, *UserRequest) (*UserResponse, error)
	ListUsers(context.Context, *emptypb.Empty) (*UserListResponse, error)
	CreateUser(context.Context, *User) (*BaseResponse, error)
	UpdateUser(context.Context, *User) (*BaseResponse, error)
	DeleteUser(context.Context, *UserRequest) (*BaseResponse, error)
*/

func (UserBiz) UpdateUser(_ context.Context, u *pb.User) (*pb.BaseResponse, error) {
	if u.Id == 0 {
		return &pb.BaseResponse{Code: 500, Message: "user id is required"}, errors.New("user id is required")
	}
	if err := global.GLOBAL_DB.Model(&model.User{}).Where("id=?", u.Id).
		Updates(model.User{Name: u.Name, Phone: u.Phone, Address: u.Address, Email: u.Email,
			Password: u.Password, UpdatedAt: time.Now().Unix()}).Error; err != nil {
		log.Println(err)
		return &pb.BaseResponse{Code: 500, Message: "failed to update user"}, err
	}
	return &pb.BaseResponse{Code: 200, Message: "success"}, nil
}

func (UserBiz) DeleteUser(_ context.Context, u *pb.UserRequest) (*pb.BaseResponse, error) {
	if err := global.GLOBAL_DB.Delete(&model.User{}, u.Id).Error; err != nil {
		log.Println(err)
		return &pb.BaseResponse{Code: 500, Message: "failed to delete user"}, err
	}

	return &pb.BaseResponse{Code: 200, Message: "success"}, nil
}

func (UserBiz) CreateUser(_ context.Context, u *pb.User) (*pb.BaseResponse, error) {
	var count int64
	var user *model.User
	if err := global.GLOBAL_DB.Model(user).Where("email=?", u.Email).Count(&count).Error; err != nil {
		log.Println(err)
		return &pb.BaseResponse{Code: 500, Message: "failed to create user"}, err
	}

	if count > 0 {
		return &pb.BaseResponse{Code: 500, Message: "email already exists"}, errors.New("user already email")
	}

	if err := global.GLOBAL_DB.Model(user).Where("phone=?", u.Phone).Count(&count).Error; err != nil {
		log.Println(err)
		return &pb.BaseResponse{Code: 500, Message: "failed to create user"}, err
	}

	if count > 0 {
		return &pb.BaseResponse{Code: 500, Message: "phone already exists"}, errors.New("user already email")
	}

	if err := global.GLOBAL_DB.Model(user).Where("name=?", u.Name).Count(&count).Error; err != nil {
		log.Println(err)
		return &pb.BaseResponse{Code: 500, Message: "failed to create user"}, err
	}

	if count > 0 {
		return &pb.BaseResponse{Code: 500, Message: "user name already exists"}, errors.New("user already email")
	}

	user = &model.User{Id: u.Id, Name: u.Name, Phone: u.Phone, Address: u.Address, Email: u.Email, Password: u.Password,
		CreatedAt: time.Now().Unix()}

	if err := global.GLOBAL_DB.Create(user).Error; err != nil {
		log.Println(err)
		return &pb.BaseResponse{Code: 500, Message: "failed to create user"}, err
	}
	return &pb.BaseResponse{Code: 200, Message: "success"}, nil
}

func (UserBiz) ListUsers(_ context.Context, _ *emptypb.Empty) (*pb.UserListResponse, error) {
	var users []model.User
	if err := global.GLOBAL_DB.Find(users).Error; err != nil {
		log.Println(err)
		return &pb.UserListResponse{Response: &pb.BaseResponse{Code: 500, Message: "failed to get users"}}, err
	}
	usersPb := make([]*pb.User, len(users))
	for i, user := range users {
		usersPb[i] = &pb.User{Id: user.Id, Name: user.Name, Phone: user.Phone, Address: user.Address,
			Email: user.Email, Password: user.Password, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}
	}
	return &pb.UserListResponse{Response: &pb.BaseResponse{Code: 200, Message: "success"}, Users: usersPb}, nil
}

func (UserBiz) GetUser(_ context.Context, ur *pb.UserRequest) (*pb.UserResponse, error) {
	var user model.User
	if err := global.GLOBAL_DB.Where("id=?", ur.Id).First(&user).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	if user.Id == 0 {
		return nil, errors.New("user not found")
	}
	return &pb.UserResponse{Response: &pb.BaseResponse{Code: 200, Message: "success"},
		User: &pb.User{Id: user.Id, Name: user.Name,
			Phone: user.Phone, Address: user.Address,
			Email: user.Email, Password: user.Password,
			CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}}, nil
}
