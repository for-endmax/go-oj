package main

import (
	"user_srv/initialize"
)

func main() {
	initialize.InitLogger() //初始化日志
	initialize.InitConfig() //读取配置信息
	initialize.InitDB()     //初始化MySQL
}
