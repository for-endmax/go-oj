package initialize

import (
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"question_web/global"
)

// InitConfig 先从本地配置中找到服务名称和conusl的信息，接着到consul上读取相关配置
func InitConfig() {
	zap.S().Info("读取配置：")
	InitLocal()
	InitRemote()
}

func InitLocal() {
	vLocal := viper.New()
	vLocal.SetConfigName("question_web") // 配置文件名称(无扩展名)
	vLocal.SetConfigType("yaml")         // 如果配置文件的名称中没有扩展名，则需要配置此项
	vLocal.AddConfigPath("./config")     // 查找配置文件所在的路径
	vLocal.AddConfigPath(".")            // 还可以在工作目录中查找配置
	err := vLocal.ReadInConfig()         // 查找并读取配置文件
	if err != nil {                      // 处理读取配置文件的错误
		panic(fmt.Errorf("读取本地配置文件失败: %s \n", err))
	}

	if err := vLocal.Unmarshal(&global.LocalConfig); err != nil {
		panic(err)
	}
	//fmt.Println("本地配置：", global.LocalConfig)
	zap.S().Info(" 本地配置：", global.LocalConfig)

}

func InitRemote() {
	vRemote := viper.New()
	err := vRemote.AddRemoteProvider("consul", fmt.Sprintf("%s:%d", global.LocalConfig.Consul.Host, global.LocalConfig.Consul.Port), global.LocalConfig.Name)
	//err := vRemote.AddRemoteProvider("consul", "127.0.0.1:8500", "go-oj/user_srv")
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
	if err := vRemote.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Info(" 远程配置：", global.ServerConfig)
}
