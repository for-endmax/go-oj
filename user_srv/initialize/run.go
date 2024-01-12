package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"user_srv/global"
	"user_srv/handler"
	"user_srv/proto"
)

// Run 启动服务
func Run() {
	//监听端口
	addr := fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port)
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("监听端口失败", err)
	}
	zap.S().Infof("启动服务: %s", addr)

	//grpc实例
	global.GrpcServer = grpc.NewServer()
	//注册
	proto.RegisterUserServer(global.GrpcServer, &handler.UserServer{})

	//启动服务
	go func() {
		err = global.GrpcServer.Serve(conn)
		if err != nil {
			zap.S().Errorw("fail server start for GRPC", err)
		}
	}()
	zap.S().Info("服务启动")
}
