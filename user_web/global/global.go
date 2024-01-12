package global

import (
	"github.com/gin-gonic/gin"
	consulApi "github.com/hashicorp/consul/api"
	"user_web/config"
	"user_web/proto"
)

var (
	LocalConfig   config.LocalConfig  // 本地配置
	ServerConfig  config.ServerConfig // 远程配置
	UserSrvClient proto.UserClient    // user_srv的grpc客户端
	GinEngine     *gin.Engine         // gin
	ConuslClient  *consulApi.Client   // consul客户端
	ServeID       string              // 服务id
)
