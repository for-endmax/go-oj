package global

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	consulApi "github.com/hashicorp/consul/api"
	"submit_web/config"
	"submit_web/proto"
)

var (
	LocalConfig     config.LocalConfig  // 本地配置
	ServerConfig    config.ServerConfig // 远程配置
	RecordSrvClient proto.RecordClient  // record_srv的grpc客户端
	GinEngine       *gin.Engine         // gin
	ConsulClient    *consulApi.Client   // consul客户端
	ServeID         string              // 服务id
	Trans           ut.Translator       //声明一个全局翻译器
)
var JWTSigningKey string = "endmax" //JWT签名
