package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"question_web/global"
	"question_web/proto"
)

// GetGrpcClient 获取下层服务的grpc客户端
func GetGrpcClient() {
	// 通过服务名称查询信息
	data, err := global.ConuslClient.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.QuestionSrvInfo.Name))
	if err != nil {
		zap.S().Errorw("向consul查询服务出错")
		return
	}
	questionSrvHost := ""
	questionSrvPort := 0
	for _, v := range data {
		questionSrvHost = v.Address
		questionSrvPort = v.Port
		break
	}

	if questionSrvHost == "" {
		zap.S().Errorw("question_srv 服务不存在")
		return
	}
	zap.S().Infof("获取 question_srv 服务：ip:%s, port:%d", questionSrvHost, questionSrvPort)
	//拨号连接grpc服务
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", questionSrvHost, questionSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("连接服务失败", err.Error())
		return
	}

	global.QuestionSrvClient = proto.NewQuestionClient(conn)
}
