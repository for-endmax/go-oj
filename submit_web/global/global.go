package global

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v8"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/streadway/amqp"
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
	Redis           *redis.Client       //Redis连接客户端
	RabbitMQChan    *amqp.Channel       //RabbitMQ channel
)
var JWTSigningKey string = "endmax" //JWT签名
var JudgeQueue string = "judge_queue"

// MsgSend 通过mq发送的记录信息
type MsgSend struct {
	ID         int32  `json:"id"`
	Lang       string `json:"lang,omitempty"`
	SubmitCode string `json:"submit_code,omitempty"`
	//新增
	TimeLimit int32 `json:"time_limit"`
	MemLimit  int32 `json:"mem_limit"`
}

// MsgReply 接收的mq回调消息
type MsgReply struct {
	ID        int32  `json:"id"`
	Status    int32  `json:"status,omitempty"`
	ErrCode   int32  `json:"err_code,omitempty"`
	ErrMsg    string `json:"err_msg,omitempty"`
	TimeUsage int32  `json:"time_usage"`
	MemUsage  int32  `json:"mem_usage"`
}
