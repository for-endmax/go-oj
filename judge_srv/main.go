package main

import (
	"judge_srv/initialize"
	"judge_srv/judge"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initialize.InitLogger()       // 初始化日志
	initialize.InitConfig()       // 初始化配置
	initialize.GetHost()          //获取ip地址
	initialize.GetPort()          // 获取端口
	initialize.InitRabbitMQ()     // 初始化rabbitMQ连接
	initialize.InitConsulClient() // 初始化consul客户端
	initialize.GetGrpcClient()    // 获取question_srv的grpc客户端
	initialize.Register()         // 服务注册

	// 初始化mq
	var mq judge.MQ
	mq.InitMQ() //初始化mq连接和队列
	mq.Run()    //从队列中获取信息，判题，并返回回调信息
	//接收终止信号
	quit := make(chan os.Signal)
	//接收control+c
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	initialize.UnRegister() // 服务注销
	mq.Close()
}
