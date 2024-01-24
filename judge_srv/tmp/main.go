package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// MsgSend 通过mq发送的记录信息
type MsgSend struct {
	ID         int32  `json:"id"`
	Lang       string `json:"lang,omitempty"`
	SubmitCode string `json:"submit_code,omitempty"`
}

// MsgReply 接收的mq回调消息
type MsgReply struct {
	ID      int32  `json:"id"`
	Status  int32  `json:"status,omitempty"`
	ErrCode int32  `json:"err_code,omitempty"`
	ErrMsg  string `json:"err_msg,omitempty"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"judge_queue", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			failOnError(err, "Failed to convert body to integer")
			var msgSend MsgSend
			err := json.Unmarshal(d.Body, &msgSend)
			if err != nil {
				log.Println("解析信息出错")
				return
			}
			fmt.Printf("读取消息msg:%s\n", string(d.Body))
			// 模拟耗时操作
			time.Sleep(time.Second * 1)

			msgReply := MsgReply{
				ID:      msgSend.ID,
				Status:  1,
				ErrCode: 0,
				ErrMsg:  "",
			}
			msg, err := json.Marshal(&msgReply)
			if err != nil {
				log.Println("生成回调信息出错")
				return
			}
			log.Printf("回调消息 %s", msg)
			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(msg),
				})
			failOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}
