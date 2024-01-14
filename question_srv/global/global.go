package global

import (
	consulApi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"question_srv/config"
)

var (
	LocalConfig  config.LocalConfig  // 本地配置
	ServerConfig config.ServerConfig // 远程配置
	DB           *gorm.DB            // MySQL对象
	ConsulClient *consulApi.Client   // consul客户端
	ServeID      string              // 服务id
	GrpcServer   *grpc.Server        // grpcServer实例
)
