package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mxshop/user_srv/proto"
)

var userClient proto.UserClient
var conn *grpc.ClientConn
var err error

func Init() {
	conn, err = grpc.NewClient("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	userClient = proto.NewUserClient(conn)
}

func TestGetUserList() {
	rsp, err_1 := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    0,
		PSize: 0,
	})
	if err_1 != nil {
		panic(err_1)
	}
	for _, user := range rsp.Data {
		fmt.Println(user)
		checkRsp, err2 := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.PassWord,
		})
		if err2 != nil {
			panic(err2)
		}
		fmt.Println(checkRsp.Success)
	}
}

func TestCreateUser() {
	for i := 0; i < 10; i++ {
		rsp, err_3 := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: fmt.Sprintf("zjy%d", i),
			Mobile:   fmt.Sprintf("1950360120%d", i),
			PassWord: "admin123",
		})
		if err_3 != nil {
			panic(err_3)
		}
		fmt.Println(rsp.Id)
	}
}

func main() {
	Init()
	//TestCreateUser()
	TestGetUserList()
	err = conn.Close()
	if err != nil {
		return
	}
}
