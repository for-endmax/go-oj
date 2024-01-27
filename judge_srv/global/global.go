package global

import (
	consulApi "github.com/hashicorp/consul/api"
	"github.com/streadway/amqp"
	"judge_srv/config"
	"judge_srv/proto"
)

var (
	LocalConfig       config.LocalConfig   // 本地配置
	ServerConfig      config.ServerConfig  // 远程配置
	QuestionSrvClient proto.QuestionClient // record_srv的grpc客户端
	ConsulClient      *consulApi.Client    // consul客户端
	ServeID           string               // 服务id
	RabbitMQChan      *amqp.Channel        //RabbitMQ channel
)
var JudgeQueue string = "judge_queue"

var JudgeDone chan struct{}
