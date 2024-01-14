package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"user_srv/proto"
)

var client proto.UserClient

func initClient() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("无法连接到gRPC服务: %v", err)
	}
	// 创建gRPC客户端
	client = proto.NewUserClient(conn)
}

func TestGetUserInfoList() {
	req := proto.PageInfo{
		PNum:  2,
		PSize: 2,
	}
	rsp, err := client.GetUserInfoList(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, v := range rsp.Data {
		fmt.Println(v.Id, v.Nickname, v.Gender, v.Role)
	}
}

func TestGetUserByNickname() {
	req := proto.NicknameRequest{Nickname: "somebody1"}
	rsp, err := client.GetUserByNickname(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id, rsp.Nickname, rsp.Gender, rsp.Role)
}

func TestGetUserById() {
	req := proto.IdRequest{Id: 2}
	rsp, err := client.GetUserById(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id, rsp.Nickname, rsp.Gender, rsp.Role)
}

func TestCreateUser() {
	req := proto.CreateUserInfo{
		Nickname: "Bob",
		Gender:   "male",
		Role:     1,
		Password: "12345",
	}
	rsp, err := client.CreateUser(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id, rsp.Nickname, rsp.Gender, rsp.Role)
}

func TestUpdateUser() {
	req := proto.UpdateUserInfo{
		Nickname: "Bob",
		Gender:   "female",
		Role:     1,
		Password: "123456",
	}
	_, err := client.UpdateUser(context.Background(), &req)
	if err != nil {
		panic(err)
	}
}

func TestCheckPassword() {
	req := proto.PasswordCheckInfo{
		Id:       2,
		Nickname: "somebody1",
		Password: "admin123",
	}
	res, err := client.CheckPassword(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Valid)
}

func main() {
	initClient()
	TestGetUserInfoList()
	//TestGetUserByNickname()
	//TestGetUserById()
	//TestCreateUser()
	//TestUpdateUser()
	//TestCheckPassword()
}
