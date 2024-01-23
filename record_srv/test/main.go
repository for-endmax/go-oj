package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"record_srv/proto"
)

var client proto.RecordClient

func initClient() {
	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("无法连接到gRPC服务: %v", err)
	}
	// 创建gRPC客户端
	client = proto.NewRecordClient(conn)
}

func TestCreateRecord() {
	record, err := client.CreateRecord(context.Background(), &proto.CreateRecordRequest{
		UID:        10,
		QID:        10,
		Lang:       "c++",
		Status:     0,
		ErrCode:    0,
		ErrMsg:     "",
		TimeLimit:  0,
		MemLimit:   0,
		SubmitCode: "",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(record.ID)
}

func TestGetAllRecordByUID() {
	records, err := client.GetAllRecordByUID(context.Background(), &proto.UIDRequest{Uid: 4, PNum: 0, PSize: 2})
	if err != nil {
		panic(err)
	}
	fmt.Println(records.Total)
	for _, v := range records.Data {
		fmt.Println(v)
	}
}

func TestGetRecordByID() {
	record, err := client.GetRecordByID(context.Background(), &proto.IDRequest{Id: 4})
	if err != nil {
		panic(err)
	}
	fmt.Println(record.ID, record.QID, record.UID, record.SubmitCode, record.ErrMsg)
}

func TestUpdateRecord() {
	record, err := client.UpdateRecord(context.Background(), &proto.UpdateRecordRequest{
		ID:      4,
		Status:  0,
		ErrCode: 0,
		ErrMsg:  "err!!!",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(record.ID, record.QID, record.UID, record.SubmitCode, record.ErrMsg)
}
func main() {
	initClient()
	//TestCreateRecord()
	TestGetAllRecordByUID()
	//TestGetRecordByID()
	//TestUpdateRecord()
}
