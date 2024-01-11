package main

import (
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"user_srv/initialize"
)

func main() {
	initialize.InitLogger()              //初始化日志
	initialize.InitConfig()              //读取配置信息
	initialize.InitDB()                  //初始化MySQL
	initialize.GetPort()                 //获取端口
	client, serverID := initialize.Run() //启动服务，并注册

	//接收终止信号
	quit := make(chan os.Signal)

	//接收control+c
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	err := client.Agent().ServiceDeregister(serverID)
	if err != nil {
		zap.S().Info("注销失败", err)
	}

	zap.S().Info("注销成功")
}
