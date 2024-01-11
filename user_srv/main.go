package main

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"user_srv/global"
	"user_srv/handler"
	"user_srv/initialize"
	"user_srv/proto"
)

func main() {
	initialize.InitLogger() //初始化日志
	initialize.InitConfig() //读取配置信息
	initialize.InitDB()     //初始化MySQL
	initialize.GetPort()    //获取端口

	//监听端口
	addr := fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port)
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("监听端口失败", err)
	}
	zap.S().Infof("启动服务: %s", addr)

	//grpc实例
	s := grpc.NewServer()

	//注册
	proto.RegisterUserServer(s, &handler.UserServer{})

	go func() {
		err = s.Serve(conn)
		if err != nil {
			zap.S().Errorw("fail server start for GRPC", err)
		}
	}()

	//接收终止信号
	qiut := make(chan os.Signal)
	//接收control+c
	signal.Notify(qiut, syscall.SIGINT, syscall.SIGTERM)
	<-qiut
	zap.S().Info("注销成功")
}
