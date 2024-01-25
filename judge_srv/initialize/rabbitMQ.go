package initialize

import (
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"judge_srv/global"
)

func InitRabbitMQ() {
	rabbitMQInfo := global.ServerConfig.RabbitMQInfo
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", rabbitMQInfo.User, rabbitMQInfo.Password, rabbitMQInfo.Host, rabbitMQInfo.Port))
	if err != nil {
		panic(err)
	}
	global.RabbitMQChan, err = conn.Channel()
	if err != nil {
		panic(err)
	}
	// 创建判题队列
	_, err = global.RabbitMQChan.QueueDeclare(
		global.JudgeQueue, // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // noWait
		nil,               // arguments
	)
	if err != nil {
		panic(err)
	}
	zap.S().Info("初始化RabbitMQ连接成功")
}
