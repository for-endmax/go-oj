package initialize

import (
	"fmt"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"user_srv/global"
)

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
	//注册健康检查
	grpc_health_v1.RegisterHealthServer(global.GrpcServer, health.NewServer())

	//生成对应的检查对象
	check := &consulApi.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.CheckHost, global.ServerConfig.Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	//生成注册对象
	registration := new(consulApi.AgentServiceRegistration)
	registration.Name = global.LocalConfig.Name
	global.ServeID, _ = uuid.GenerateUUID()
	registration.ID = global.ServeID
	registration.Port = global.ServerConfig.Port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = global.ServerConfig.CheckHost
	registration.Check = check

	err := global.ConsulClient.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	zap.S().Info("服务注册成功")

}

func UnRegister() {
	err := global.ConsulClient.Agent().ServiceDeregister(global.ServeID)
	if err != nil {
		zap.S().Info("注销失败", err)
	}
	zap.S().Info("注销成功")
}
