package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"mxshop/user_srv/global"
	"mxshop/user_srv/model"
	"mxshop/user_srv/proto"
	"strings"
	"time"
)

type UserServer struct {
	proto.UnimplementedUserServer //避免未实现的方法进行编译，从而导致的报错
}

func ModelToResponse(user model.User) *proto.UserInfoResponse {
	userInfoRsp := &proto.UserInfoResponse{
		Id:       user.ID,
		PassWord: user.PassWord,
		NickName: user.Nickname,
		Gender:   user.Gender,
		Mobile:   user.Mobile,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
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

func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	//获取用户列表
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println("用户列表")
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users) //使用scopes方法进行分页查询

	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, userInfoRsp)
	}
	return rsp, nil
}

func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	//通过手机号码查询用户
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	userInfoRsp := ModelToResponse(user)
	return userInfoRsp, nil
}

func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	// 通过ID查询用户
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	userInfoRsp := ModelToResponse(user)
	return userInfoRsp, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	// 新建用户
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Error(codes.AlreadyExists, "用户已存在")
	}
	user.Mobile = req.Mobile
	user.Nickname = req.NickName
	// 密码加密
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(req.PassWord, options)
	user.PassWord = fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodedPwd)

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	userInfoRsp := ModelToResponse(user)
	return userInfoRsp, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	// 个人更新用户
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	birthday := time.Unix(int64(req.BirthDay), 0)
	user.Nickname = req.NickName
	user.Birthday = &birthday
	user.Gender = req.Gender
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *UserServer) CheckPassWord(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	// 校验密码
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword, "$") // 存储格式为: 算法$salt$hash值
	check := password.Verify(req.Password, passwordInfo[1], passwordInfo[2], options)
	return &proto.CheckResponse{Success: check}, nil
}
