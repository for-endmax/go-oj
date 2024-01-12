package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"user_web/global"
	"user_web/proto"
)

// GetGrpcClient 获取下层服务的grpc客户端
func GetGrpcClient() {
	// 通过服务名称查询信息
	data, err := global.ConuslClient.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		zap.S().Errorw("向consul查询服务出错")
		return
	}
	userSrvHost := ""
	userSrvPort := 0
	for _, v := range data {
		userSrvHost = v.Address
		userSrvPort = v.Port
		break
	}

	if userSrvHost == "" {
		zap.S().Errorw("user_srv 服务不存在")
		return
	}
	zap.S().Infof("获取 user_srv 服务：ip:%s, port:%d", userSrvHost, userSrvPort)
	//拨号连接grpc服务
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("连接服务失败", err.Error())
		return
	}

	global.UserSrvClient = proto.NewUserClient(conn)
}
