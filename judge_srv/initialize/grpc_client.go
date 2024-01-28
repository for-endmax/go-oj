package initialize

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"judge_srv/global"
	"judge_srv/otgrpc"
	"judge_srv/proto"
)

// GetGrpcClient 获取下层服务的grpc客户端
func GetGrpcClient() {
	// 通过服务名称查询信息
	data, err := global.ConsulClient.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.QuestionSrvInfo.Name))
	if err != nil {
		zap.S().Errorw("向consul查询服务出错")
		return
	}
	recordSrvHost := ""
	recordSrvPort := 0
	for _, v := range data {
		recordSrvHost = v.Address
		recordSrvPort = v.Port
		break
	}

	if recordSrvHost == "" {
		zap.S().Errorw("record_srv 服务不存在")
		return
	}
	zap.S().Infof("获取 record_srv 服务：ip:%s, port:%d", recordSrvHost, recordSrvPort)
	//拨号连接grpc服务
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", recordSrvHost, recordSrvPort), grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	if err != nil {
		zap.S().Errorw("连接服务失败", err.Error())
		return
	}

	global.QuestionSrvClient = proto.NewQuestionClient(conn)
}
