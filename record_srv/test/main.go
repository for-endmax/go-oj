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
		UID:        111,
		QID:        111,
		Lang:       "c++",
		TimeLimit:  500,
		MemLimit:   500,
		SubmitCode: "hello",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(record.ID)
}

func TestGetAllRecordByUID() {
	records, err := client.GetAllRecordByUID(context.Background(), &proto.UIDRequest{Uid: 111, PNum: 0, PSize: 2})
	if err != nil {
		panic(err)
	}
	fmt.Println(records.Total)
	for _, v := range records.Data {
		fmt.Println(v)
	}
}

func TestGetRecordByID() {
	record, err := client.GetRecordByID(context.Background(), &proto.IDRequest{Id: 11})
	if err != nil {
		panic(err)
	}
	fmt.Println(record.ID, record.QID, record.UID, record.SubmitCode, record.ErrMsg, record.TimeUsage, record.MemUsage)
}

func TestUpdateRecord() {
	record, err := client.UpdateRecord(context.Background(), &proto.UpdateRecordRequest{
		ID:        11,
		Status:    -1,
		ErrCode:   0,
		ErrMsg:    "err!!!",
		MemUsage:  12312,
		TimeUsage: 1431123,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(record.ID, record.QID, record.UID, record.SubmitCode, record.ErrMsg)
}
func main() {
	initClient()
	//TestCreateRecord()
	//TestGetAllRecordByUID()
	//TestGetRecordByID()
	//TestUpdateRecord()
}
