package main

import (
	"os"
	"os/signal"
	"syscall"
	"user_srv/initialize"
)

func main() {
	initialize.InitLogger()       // 初始化日志
	initialize.InitConfig()       // 读取配置信息
	initialize.InitDB()           // 初始化MySQL
	initialize.GetPort()          // 获取端口
	initialize.InitConsulClient() // 初始化consul客户端
	initialize.Register()         // 服务注册
	initialize.Run()              // 启动服务

	//接收终止信号
	quit := make(chan os.Signal)
	//接收control+c
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	initialize.UnRegister() // 服务注销
}
