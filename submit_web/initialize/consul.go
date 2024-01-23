package initialize

import (
	"fmt"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	"go.uber.org/zap"
	"submit_web/global"
)

// InitConsulClient 初始化consul客户端
func InitConsulClient() {
	cfg := consulApi.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.LocalConfig.Consul.Host, global.LocalConfig.Consul.Port)
	var err error
	global.ConsulClient, err = consulApi.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	zap.S().Info("创建Consul客户端成功")
}

// Register 服务注册
func Register() {
	//服务注册
	//生成对应的检查对象
	check := &consulApi.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", global.ServerConfig.CheckHost, global.ServerConfig.Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	global.ServeID, _ = uuid.GenerateUUID()
	srv := &consulApi.AgentServiceRegistration{
		ID:      global.ServeID,           // 服务唯一ID
		Name:    global.LocalConfig.Name,  // 服务名称
		Tags:    global.ServerConfig.Tags, // 服务标签
		Address: global.ServerConfig.CheckHost,
		Port:    global.ServerConfig.Port,
		Check:   check,
	}
	err := global.ConsulClient.Agent().ServiceRegister(srv)
	if err != nil {
		zap.S().Error("服务注册失败")
		panic(err)
	}
	zap.S().Info("服务注册成功")
}

// UnRegister 服务注销
func UnRegister() {
	//服务注销
	err := global.ConsulClient.Agent().ServiceDeregister(global.ServeID)
	if err != nil {
		zap.S().Info("注销失败", err)
	}
	zap.S().Info("注销成功")
}
