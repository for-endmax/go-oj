package judge

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"judge_srv/global"
	"judge_srv/message"
	"time"
)

type MQ struct {
	Conn  *amqp.Connection
	Chan  *amqp.Channel
	Queue amqp.Queue
	Msgs  <-chan amqp.Delivery
}

// InitMQ 初始化mq
func (m *MQ) InitMQ() {
	mq := global.ServerConfig.RabbitMQInfo
	var err error
	m.Conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", mq.User, mq.Password, mq.Host, mq.Port))
	if err != nil {
		panic(err)
	}
	m.Chan, err = m.Conn.Channel()
	if err != nil {
		panic(err)
	}

	m.Queue, err = m.Chan.QueueDeclare(
		"judge_queue", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		panic(err)
	}

	err = m.Chan.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		panic(err)
	}

	m.Msgs, err = m.Chan.Consume(
		m.Queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		panic(err)
	}
}

// Close 关闭连接
func (m *MQ) Close() {
	m.Conn.Close()
	m.Chan.Close()
}

// Run 解析
func (m *MQ) Run() {
	go func() {
		for d := range m.Msgs {
			var msgSend message.MsgSend
			err := json.Unmarshal(d.Body, &msgSend)
			if err != nil {
				zap.S().Info("从mq中解析信息出错")
				continue
			}
			zap.S().Infof("读取消息msg:%s\n", string(d.Body))

			// 进行判题，返回回调信息
			msgReply, err := Judge(msgSend)
			if err != nil {
				zap.S().Info("判题失败")
				continue
			}
			msg, err := json.Marshal(&msgReply)
			if err != nil {
				zap.S().Info("生成回调信息出错")
				continue
			}

			zap.S().Info("回调消息 %s", string(msg))
			err = m.Chan.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          msg,
				})
			if err != nil {
				zap.S().Info("返回回调信息出错")
				continue
			}
			d.Ack(false)
		}
	}()
}

// Judge 判题
func Judge(msgSend message.MsgSend) (message.MsgReply, error) {
	time.Sleep(time.Second * 2)
	return message.MsgReply{
		ID:        msgSend.ID,
		Status:    1,
		ErrCode:   0,
		ErrMsg:    "no err",
		MemUsage:  100,
		TimeUsage: 10,
	}, nil
}
