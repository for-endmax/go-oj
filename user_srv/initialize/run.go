package initialize

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"user_srv/global"
	"user_srv/handler"
	"user_srv/proto"
)

// Run 启动服务
func Run() (*consulapi.Client, string) {
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

	//注册健康检查
	//将user_srv服务注册到consul中，让web层可获取其配置信息
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	//DefaultConfig 返回客户端的默认配置
	cfg := consulapi.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.LocalConfig.Consul.Host, global.LocalConfig.Consul.Port)

	client, err := consulapi.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &consulapi.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.CheckHost, global.ServerConfig.Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	//生成注册对象
	registration := new(consulapi.AgentServiceRegistration)
	registration.Name = global.LocalConfig.Name
	serverID, _ := uuid.GenerateUUID()
	registration.ID = serverID
	registration.Port = global.ServerConfig.Port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = global.ServerConfig.CheckHost
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	zap.S().Info("服务注册成功")

	//启动服务
	go func() {
		err = s.Serve(conn)
		if err != nil {
			zap.S().Errorw("fail server start for GRPC", err)
		}
	}()
	zap.S().Info("服务启动")
	return client, serverID
}
