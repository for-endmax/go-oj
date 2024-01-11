package main

import (
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

// 要实现的目标
// 		1.viper读取本地配置文件(本地配置文件记录了consul的信息)
//      2.viper读取consul上的配置信息

// docker安装consul
/*
docker run -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp consul consul agent -dev -client=0.0.0.0
*/
type LocalConfig struct {
	Name   string       `mapstructure:"name"`
	Consul ConsulConfig `mapstructure:"consul"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	MysqlHost string `mapstructure:"mysql_host"`
	MysqlPort int    `mapstructure:"mysql_port"`
}

var localConfig LocalConfig
var serverConfig ServerConfig

func InitLocal() {
	vLocal := viper.New()
	vLocal.SetConfigName("config")         // 配置文件名称(无扩展名)
	vLocal.SetConfigType("yaml")           // 如果配置文件的名称中没有扩展名，则需要配置此项
	vLocal.AddConfigPath("./02viper_init") // 查找配置文件所在的路径
	vLocal.AddConfigPath(".")              // 还可以在工作目录中查找配置
	err := vLocal.ReadInConfig()           // 查找并读取配置文件
	if err != nil {                        // 处理读取配置文件的错误
		panic(fmt.Errorf("读取本地配置文件失败: %s \n", err))
	}

	if err := vLocal.Unmarshal(&localConfig); err != nil {
		panic(err)
	}
	fmt.Println("本地配置：", localConfig)

}

func InitRemote() {
	vRemote := viper.New()
	//err := vRemote.AddRemoteProvider("consul", fmt.Sprintf("%s:%d", localConfig.Consul.Host, localConfig.Consul.Port), localConfig.Name)
	err := vRemote.AddRemoteProvider("consul", "127.0.0.1:8500", "02viper_init")
	if err != nil {
		panic(fmt.Errorf("viper读取consul配置失败: %s \n", err))
		return
	}
	vRemote.SetConfigType("yaml")
	err = vRemote.ReadRemoteConfig()
	if err != nil {
		panic(err)
		return
	}
	if err := vRemote.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	fmt.Println("远程配置：", serverConfig)
}

func InitViper() {
	InitLocal()
	InitRemote()
}

func main() {
	InitViper()
}
