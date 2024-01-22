package main

import (
	"os"
	"os/signal"
	"question_web/initialize"
	"syscall"
)

func main() {
	initialize.InitLogger()       // 初始化日志
	initialize.InitConfig()       // 初始化配置
	initialize.GetPort()          // 获取端口
	initialize.InitTrans("zh")    //初始化validator翻译
	initialize.InitConsulClient() // 初始化consul客户端
	initialize.GetGrpcClient()    // 获取question_srv的grpc客户端
	initialize.InitRouter()       // 初始化路由
	initialize.Run()              // 启动服务
	initialize.Register()         // 服务注册
	//接收终止信号
	quit := make(chan os.Signal)
	//接收control+c
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	initialize.UnRegister() // 服务注销
}
