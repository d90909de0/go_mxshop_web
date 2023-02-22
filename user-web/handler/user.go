package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"go_mxshop_srvs/user_srv/global"
	"go_mxshop_srvs/user_srv/model"
	"go_mxshop_srvs/user_srv/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserServer struct {
}

var passwordOptions = &password.Options{16, 100, 32, sha512.New}

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

func encodePassword(s string) string {
	// Using custom options
	salt, encodedPwd := password.Encode(s, passwordOptions)
	// 通过$符号组合加密算法、盐值、加密后数据
	return fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
}

func verifyPassword(pwd, encodePwd string) bool {
	pwdOption := strings.Split(encodePwd, "$")
	return password.Verify(pwd, pwdOption[2], pwdOption[3], passwordOptions)
}

func modelToResponse(user model.User) *proto.UserInfoResponse {
	response := proto.UserInfoResponse{
		Id:       int32(user.ID),
		Password: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     uint32(user.Role),
	}

	if user.Birthday != nil {
		response.BrithDay = uint64(user.Birthday.Unix())
	}
	return &response
}

func (s UserServer) GetUserList(ctx context.Context, pageInfo *proto.PageInfo) (*proto.UserListResponse, error) {
	var count int64
	var users []model.User
	rsp := proto.UserListResponse{}

	global.DB.Count(&count)
	rsp.Total = int32(count)

	global.DB.Scopes(Paginate(int(pageInfo.Pn), int(pageInfo.PSize))).Find(&users)
	for _, user := range users {
		rsp.Data = append(rsp.Data, modelToResponse(user))
	}

	return &rsp, nil
}

func (s UserServer) GetUserById(ctx context.Context, request *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	response := modelToResponse(user)
	return response, nil
}

func (s UserServer) GetUserByMobile(ctx context.Context, request *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: request.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	response := modelToResponse(user)
	return response, nil
}

func (s UserServer) CreateUser(ctx context.Context, userInfo *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	user := model.User{
		Mobile:   userInfo.Mobile,
		NickName: userInfo.NickName,
		Password: encodePassword(userInfo.Password),
	}

	result := global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	response := modelToResponse(user)
	return response, nil
}

func (s UserServer) UpdateUser(ctx context.Context, userInfo *proto.UpdateUserInfo) (*empty.Empty, error) {
	var user model.User
	result := global.DB.First(&user, userInfo.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	birthday := time.Unix(int64(userInfo.BrithDay), 0)

	user.NickName = userInfo.NickName
	user.Gender = userInfo.Gender
	user.Birthday = &birthday

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	return &empty.Empty{}, nil
}

func (s UserServer) CheckPassword(ctx context.Context, req *proto.CheckPasswordInfo) (*proto.CheckPasswordResponse, error) {
	success := verifyPassword(req.Password, req.EncryptedPassword)
	response := proto.CheckPasswordResponse{
		Success: success,
	}
	return &response, nil
}
