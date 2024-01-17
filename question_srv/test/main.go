package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"question_srv/proto"
)

var client proto.QuestionClient

func initClient() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("无法连接到gRPC服务: %v", err)
	}
	// 创建gRPC客户端
	client = proto.NewQuestionClient(conn)
}

func TestGetQuestionList() {
	req := proto.PageInfoRequest{
		PNum:  2,
		PSize: 5,
	}
	rsp, err := client.GetQuestionList(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, v := range rsp.Data {
		fmt.Println(v.Id, v.Seq, v.Name, v.Content)
	}
}

func TestGetQuestionInfo() {
	req := proto.GetQuestionInfoRequest{Id: 1}
	rsp, err := client.GetQuestionInfo(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id, rsp.Seq, rsp.Name, rsp.Content)
}

func TestAddQuestion() {
	req := proto.AddQuestionRequest{
		Seq:     100,
		Name:    "测试题目",
		Content: "aaaaaaaaaaaaaaaaaaa",
	}
	question, err := client.AddQuestion(context.Background(), &req)
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(question.Id, question.Seq, question.Name, question.Content)
}

func TestDelQuestion() {
	req := proto.DelQuestionRequest{Id: 13}
	success, err := client.DelQuestion(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(success)
}

func TestUpdateQuestion() {
	req := proto.UpdateQuestionRequest{
		Id:      11,
		Seq:     1000,
		Name:    "修改题目",
		Content: "14312412412",
	}
	success, err := client.UpdateQuestion(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(success)
}

func TestGetTestInfo() {
	req := proto.GetTestRequest{QId: 5}
	info, err := client.GetTestInfo(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(info.Total)
	for _, v := range info.Data {
		fmt.Println(v.Id, v.QId, v.TimeLimit, v.MemLimit, v.Input, v.Output)
	}
}

func TestAddTest() {
	req := proto.AddTestRequest{
		QId:       5,
		TimeLimit: 100,
		MemLimit:  100,
		Input:     "111111",
		Output:    "222222",
	}
	info, err := client.AddTest(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(info.Id, info.QId, info.MemLimit, info.TimeLimit, info.Input, info.Output)
}

func TestDelTest() {
	req := proto.DelTestRequest{Id: 9}
	info, err := client.DelTest(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(info.Success)
}

func TestUpdateTest() {
	req := proto.UpdateTestRequest{
		Id:  3,
		QId: 5,
	}
	info, err := client.UpdateTest(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(info.Success)
}
func main() {
	initClient()
	//TestGetQuestionList()
	//TestGetQuestionInfo()
	//TestAddQuestion()
	//TestDelQuestion()
	//TestUpdateQuestion()

	//TestGetTestInfo()
	//TestAddTest()
	//TestDelTest()
	//TestUpdateTest()
}
