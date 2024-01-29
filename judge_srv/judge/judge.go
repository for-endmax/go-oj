package judge

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/streadway/amqp"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"judge_srv/global"
	"judge_srv/message"
	"judge_srv/proto"
	"strconv"
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
		global.Qos, // prefetch count
		0,          // prefetch size
		false,      // global
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

	//global.JudgeDone = make(chan struct{})
	//close(global.JudgeDone)
}

// Close 关闭连接
func (m *MQ) Close() {
	m.Conn.Close()
	m.Chan.Close()
}

// Run 解析
func (m *MQ) Run() {
	for msg := range m.Msgs {
		d := msg
		go func() { //读取消息头
			cfg := jaegercfg.Configuration{
				Sampler: &jaegercfg.SamplerConfig{
					Type:  jaeger.SamplerTypeConst,
					Param: 1,
				},
				Reporter: &jaegercfg.ReporterConfig{
					LogSpans:           true,
					LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port),
				},
				ServiceName: "go-oj/judge_srv",
			}
			tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
			defer closer.Close()
			if err != nil {
				zap.S().Info("创建tracer失败")
				return
			}
			opentracing.SetGlobalTracer(tracer)
			headers := amqpTableCarrier{headers: d.Headers}
			spanContext, err := tracer.Extract(opentracing.TextMap, headers)
			judgeSpan := opentracing.StartSpan("judge", opentracing.ChildOf(spanContext))
			defer judgeSpan.Finish()
			ctx := opentracing.ContextWithSpan(context.Background(), judgeSpan)
			if err != nil {
				zap.S().Info("读取消息头错误")
				return
			}

			//读取消息体
			var msgSend message.MsgSend
			err = json.Unmarshal(d.Body, &msgSend)
			if err != nil {
				zap.S().Info("从mq中解析信息出错")
				return
			}
			zap.S().Infof("读取消息msg:%s", string(d.Body))

			//base64解码
			realCode := make([]byte, base64.StdEncoding.DecodedLen(len(msgSend.SubmitCode)))
			_, err = base64.StdEncoding.Decode(realCode, []byte(msgSend.SubmitCode))
			if err != nil {
				zap.S().Info("代码解码错误")
				//返回回调信息
				err = m.Reply(ctx, d, &message.MsgReply{
					ID:        msgSend.ID,
					Status:    1,
					ErrCode:   1,
					ErrMsg:    "代码编码错误",
					TimeUsage: 0,
					MemUsage:  0,
				})
				if err != nil {
					return
				}
				return
			}
			msgSend.SubmitCode = string(realCode)
			zap.S().Infof("解码后的提交代码:\n%s", msgSend.SubmitCode)

			// 进行判题
			//global.JudgeDone = make(chan struct{})
			msgReply, err := Judge(ctx, msgSend)
			//close(global.JudgeDone)

			if err != nil {
				zap.S().Info("判题失败")
				return
			}
			//返回回调信息
			err = m.Reply(ctx, d, msgReply)
			if err != nil {
				return
			}
			judgeSpan.LogKV("recordID", msgReply.ID, "status", msgReply.Status, "errCode", msgReply.ErrCode, "errMsg", msgReply.ErrMsg)
		}()
	}
}

// Reply 返回信息
func (m *MQ) Reply(ctx context.Context, d amqp.Delivery, msgReply *message.MsgReply) error {
	judgeSpan := opentracing.SpanFromContext(ctx)
	replySpan := opentracing.StartSpan("reply", opentracing.ChildOf(judgeSpan.Context()))
	defer replySpan.Finish()
	msg, err := json.Marshal(msgReply)
	if err != nil {
		zap.S().Info("生成回调信息出错")
		return err
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
		return err
	}
	d.Ack(false)
	return nil
}

// Judge 判题
func Judge(ctx context.Context, msgSend message.MsgSend) (*message.MsgReply, error) {
	judgeSpan := opentracing.SpanFromContext(ctx)

	createTaskSpan := opentracing.StartSpan("CreateTask", opentracing.ChildOf(judgeSpan.Context()))
	// 创建task
	task, err := CreateTask(msgSend)
	if err != nil {
		return nil, err
	}
	defer task.Clean()
	createTaskSpan.Finish()
	// 查询test信息
	testInfos, err := global.QuestionSrvClient.GetTestInfo(ctx, &proto.GetTestRequest{
		QId: task.msgSend.QID,
	})
	if err != nil {
		zap.S().Infof("查询record信息出错,QID:%d", task.msgSend.QID)
		return nil, err
	}
	if testInfos.Total == 0 {
		zap.S().Info("没有测试信息")
		return nil, fmt.Errorf("题目：%d,没有测试信息", task.msgSend.QID)
	}
	zap.S().Infof("查询到的测试信息total:%d", testInfos.Total)

	//判题
	var totalTimeUsage int32 = 0
	var totalMemUsage int32 = 0
	for i, v := range testInfos.Data {
		RunSpan := opentracing.StartSpan("runRecord:"+strconv.Itoa(i), opentracing.ChildOf(judgeSpan.Context()))
		err, result := task.Run(v.Input, v.Output, i)
		if err != nil {
			zap.S().Info("判题出错")
			return nil, err
		}

		zap.S().Infof("测试用例：%d ,判题状态码：%d ,运行时间：%d ms,运行内存: %d KB, 错误信息：%s", i, result.ErrCode, result.runTime, result.runMem, result.ErrMsg)
		totalTimeUsage += result.runTime
		totalMemUsage += result.runMem
		if result.ErrCode != 0 {
			zap.S().Info("判题未通过")
			RunSpan.Finish()
			return &message.MsgReply{
				ID:      task.msgSend.ID,
				Status:  1,
				ErrCode: result.ErrCode,
				ErrMsg:  result.ErrMsg,
			}, nil
		}
		RunSpan.Finish()
	}
	avgMemUsage, avgTimeUsage := totalMemUsage/testInfos.Total, totalTimeUsage/testInfos.Total
	//判断是否超时或超内存
	if avgMemUsage > task.msgSend.MemLimit {
		//超内存
		return &message.MsgReply{
			ID:        task.msgSend.ID,
			Status:    1,
			ErrCode:   4,
			ErrMsg:    "超内存",
			TimeUsage: avgTimeUsage,
			MemUsage:  avgTimeUsage,
		}, nil
	}
	if avgTimeUsage > task.msgSend.TimeLimit {
		//超时
		return &message.MsgReply{
			ID:        task.msgSend.ID,
			Status:    1,
			ErrCode:   2,
			ErrMsg:    "超时",
			TimeUsage: avgTimeUsage,
			MemUsage:  avgTimeUsage,
		}, nil
	}

	return &message.MsgReply{
		ID:        task.msgSend.ID,
		Status:    1,
		ErrCode:   0,
		ErrMsg:    "",
		TimeUsage: avgTimeUsage,
		MemUsage:  avgMemUsage,
	}, nil
}
