package handler

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"user_srv/global"
	"user_srv/model"
	"user_srv/proto"
	"user_srv/utils"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

// Paginate 将数据进行分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// GetUserInfoList 获取用户信息列表
func (s *UserServer) GetUserInfoList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	var total int64
	global.DB.Model(&model.User{}).Count(&total)
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(total)
	global.DB.Scopes(Paginate(int(req.PNum), int(req.PSize))).Find(&users)

	for _, user := range users {
		userInfoRsp := proto.UserInfoResponse{
			Id:       user.ID,
			Nickname: user.Nickname,
			Gender:   user.Gender,
			Role:     user.Role,
		}
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}

// GetUserByNickname 通过用户名获取用户信息
func (s *UserServer) GetUserByNickname(ctx context.Context, req *proto.NicknameRequest) (*proto.UserInfoResponse, error) {
	user := model.User{}
	result := global.DB.Where(&model.User{Nickname: req.Nickname}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	rsp := &proto.UserInfoResponse{
		Id:       user.ID,
		Nickname: user.Nickname,
		Gender:   user.Gender,
		Role:     user.Role,
	}
	return rsp, nil
}

// GetUserById 通过id获取用户信息
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	user := model.User{}
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	rsp := &proto.UserInfoResponse{
		Id:       user.ID,
		Nickname: user.Nickname,
		Gender:   user.Gender,
		Role:     user.Role,
	}
	return rsp, nil
}

// CreateUser 创建用户
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	result := global.DB.Where(&model.User{Nickname: req.Nickname}).Find(&model.User{})
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户名已经存在")
	}
	salt, encodedPassword := utils.EncodePassword(req.Password)
	user := model.User{
		Password: salt + ":" + encodedPassword,
		Nickname: req.Nickname,
		Gender:   req.Gender,
		Role:     req.Role,
	}
	global.DB.Create(&user)
	global.DB.Where(&model.User{Nickname: user.Nickname}).Find(&user)
	rsp := &proto.UserInfoResponse{
		Id:       user.ID,
		Nickname: user.Nickname,
		Gender:   user.Gender,
		Role:     user.Role,
	}
	return rsp, nil
}

// UpdateUser 更新用户信息
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	//个人中心更新用户
	var user model.User
	result := global.DB.Where(&model.User{Nickname: req.Nickname}).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	user.Nickname = req.Nickname
	user.Gender = req.Gender
	user.Role = req.Role
	salt, encodedPwd := utils.EncodePassword(req.Password)
	user.Password = salt + ":" + encodedPwd

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &empty.Empty{}, nil
}

// CheckPassword 检查密码是否正确
func (s *UserServer) CheckPassword(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Nickname: req.Nickname}).Where(&model.User{BaseModel: model.BaseModel{ID: req.Id}}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	res := utils.VerifyPassword(req.Password, user.Password)
	return &proto.CheckResponse{Valid: res}, nil
}
