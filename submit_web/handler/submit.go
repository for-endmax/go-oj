package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/go-uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"strconv"
	"strings"
	"submit_web/form"
	"submit_web/global"
	"submit_web/model"
	"submit_web/proto"
	"submit_web/response"
	"time"
)

// HandleGrpcErrorToHttp 错误处理
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	// 将grpc的code转换为http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.AlreadyExists:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "已存在",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "题目服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "其他错误",
			})
		}
		return
	}
}

// RemoveTopStruct 去除以"."及其左部分内容
func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, value := range fields {
		res[field[strings.Index(field, ".")+1:]] = value
	}
	return res
}

// HandleValidatorErr 表单验证错误处理返回
func HandleValidatorErr(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": RemoveTopStruct(errs.Translate(global.Trans)),
	})
}

// CheckRole 验证用户权限
func CheckRole(c *gin.Context) bool {
	// 验证用户权限
	role, ok := c.Get("userRole")
	if !ok {
		zap.S().Error("从context中获取值出错")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "无权限",
		})
		return false
	}
	value, ok := role.(int32)
	if !ok {
		zap.S().Error("从context中获取值出错")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "无权限",
		})
		return false
	}
	if int(value) != 2 {
		zap.S().Info("无权限")
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "无权限",
		})
		return false
	}
	return true
}

////////////////////////////////////////////////

// GetRecordListByUID 通过uid获取全部记录
func GetRecordListByUID(c *gin.Context) {
	// 获取get参数
	uID, err := strconv.Atoi(c.Query("u_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "参数错误")
		return
	}
	pNum, err := strconv.Atoi(c.Query("pn"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "参数错误")
		return
	}
	pSize, err := strconv.Atoi(c.Query("ps"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "参数错误")
		return
	}

	// 调用rpc
	recordInfoList, err := global.RecordSrvClient.GetAllRecordByUID(context.Background(), &proto.UIDRequest{
		Uid:   int32(uID),
		PNum:  int32(pNum),
		PSize: int32(pSize),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	var rsp response.RecordInfoListResponse
	rsp.Total = recordInfoList.Total
	for _, v := range recordInfoList.Data {
		record := response.RecordInfoResponse{
			ID:         v.ID,
			UID:        v.UID,
			QID:        v.QID,
			Lang:       v.Lang,
			Status:     v.Status,
			ErrCode:    v.ErrCode,
			ErrMsg:     v.ErrMsg,
			TimeLimit:  v.TimeLimit,
			MemLimit:   v.MemLimit,
			SubmitCode: v.SubmitCode,
			MemUsage:   v.MemUsage,
			TimeUsage:  v.TimeUsage,
		}
		rsp.Data = append(rsp.Data, record)
	}
	// 返回结果
	c.JSON(http.StatusOK, rsp)
}

// GetRecordByID 获取指定id的record的信息(长连接读取最新状态)
func GetRecordByID(c *gin.Context) {
	//获取参数
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "参数错误")
		return
	}
	// 首先查询一次，判断是否要等待
	recordInfo, err := global.RecordSrvClient.GetRecordByID(context.Background(), &proto.IDRequest{Id: int32(id)})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	rsp := response.RecordInfoResponse{
		ID:         recordInfo.ID,
		UID:        recordInfo.UID,
		QID:        recordInfo.QID,
		Lang:       recordInfo.Lang,
		Status:     recordInfo.Status,
		ErrCode:    recordInfo.ErrCode,
		ErrMsg:     recordInfo.ErrMsg,
		TimeLimit:  recordInfo.TimeLimit,
		MemLimit:   recordInfo.MemLimit,
		SubmitCode: recordInfo.SubmitCode,
		MemUsage:   recordInfo.MemUsage,
		TimeUsage:  recordInfo.TimeUsage,
	}
	if recordInfo.Status != 0 {
		// 返回状态
		c.JSON(http.StatusOK, rsp)
		return
	}
	// 仍在判题,等待
	zap.S().Info("查询到判题未结束,等待状态更新")
	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	// 每2s读取一次redis，若有值说明状态更新了，超时时间为60s
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.Tick(time.Second * 2): //每两秒读取一次
				zap.S().Info("读取redis中的值")
				_, err := global.Redis.Get(context.Background(), strconv.Itoa(id)).Result()
				if err != nil {
					if !errors.Is(err, redis.Nil) {
						zap.S().Infof("读取redis中的记录id: %s失败", strconv.Itoa(id))
					}
					zap.S().Info("状态未更新")
					continue
				}
				//如果有，说明状态更新
				close(done)
				zap.S().Info("读取到最新状态")
				return
			}
		}
	}(ctx)

	select {
	case <-time.After(time.Second * 60):
		//超时
		cancel()
		zap.S().Info("超时,协程退出")
		c.JSON(http.StatusNotModified, gin.H{})
		return
	case <-done:
		cancel()
		// 返回最新状态
		recordInfo, err = global.RecordSrvClient.GetRecordByID(context.Background(), &proto.IDRequest{Id: int32(id)})
		if err != nil {
			HandleGrpcErrorToHttp(err, c)
			return
		}
		rsp = response.RecordInfoResponse{
			ID:         recordInfo.ID,
			UID:        recordInfo.UID,
			QID:        recordInfo.QID,
			Lang:       recordInfo.Lang,
			Status:     recordInfo.Status,
			ErrCode:    recordInfo.ErrCode,
			ErrMsg:     recordInfo.ErrMsg,
			TimeLimit:  recordInfo.TimeLimit,
			MemLimit:   recordInfo.MemLimit,
			SubmitCode: recordInfo.SubmitCode,
			TimeUsage:  recordInfo.TimeUsage,
			MemUsage:   recordInfo.MemUsage,
		}
		c.JSON(http.StatusOK, rsp)
	}
}

// Submit 提交代码
func Submit(c *gin.Context) {

	// 读取表单
	var submitForm form.SubmitForm
	if err := c.ShouldBindJSON(&submitForm); err != nil {
		HandleValidatorErr(c, err)
		return
	}

	// 验证身份,只能以自己的uid来提交
	claims, exist := c.Get("claims")
	if !exist {
		c.JSON(http.StatusForbidden, "没有权限")
		return
	}
	customClaims := claims.(*model.CustomClaims)
	if customClaims.ID != uint(submitForm.UID) {
		c.JSON(http.StatusForbidden, "没有权限")
		return
	}
	//提交
	////////////////////////////////////////////

	// 调用rpc,生成record
	record, err := global.RecordSrvClient.CreateRecord(context.WithValue(context.Background(), "ginContext", c), &proto.CreateRecordRequest{
		UID:        submitForm.UID,
		QID:        submitForm.QID,
		Lang:       submitForm.Lang,
		SubmitCode: submitForm.SubmitCode,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	// 将 record信息 放到mq中并启动协程监听回调
	zap.S().Infof("将record放到mq, id : %d", record.ID)

	recordMsg := global.MsgSend{
		ID:         record.ID,
		Lang:       record.Lang,
		SubmitCode: record.SubmitCode,
		MemLimit:   record.MemLimit,
		TimeLimit:  record.TimeLimit,
		QID:        record.QID,
	}
	err = Send2MQ(c, recordMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":       "内部错误",
			"record_id": record.ID,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":       "创建成功",
		"record_id": record.ID,
	})
}

// Send2MQ 发送消息并监听回调
func Send2MQ(c *gin.Context, recordMsg global.MsgSend) error {
	span := c.Value("parentSpan")
	var ok bool
	var parentSpan opentracing.Span
	if parentSpan, ok = span.(opentracing.Span); !ok {
		c.JSON(http.StatusInternalServerError, "内部错误")
		return fmt.Errorf("获取span错误")
	}
	sendSpan := opentracing.StartSpan("send2MQ", opentracing.ChildOf(parentSpan.Context()))
	defer sendSpan.Finish()

	msgID, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}

	msg, err := json.Marshal(recordMsg)
	if err != nil {
		return err
	}
	//创建回调队列,每次调用方法会生成一个
	replyQueue, err := global.RabbitMQChan.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)

	//创建消费者
	msgs, err := global.RabbitMQChan.Consume(
		replyQueue.Name, // queue
		"",              // consumer
		false,           // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)

	zap.S().Info("发送信息 ", string(msg))
	//发送信息
	err = global.RabbitMQChan.Publish(
		"",                // exchange
		global.JudgeQueue, // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: msgID,
			ReplyTo:       replyQueue.Name,
			Body:          msg,
		})
	if err != nil {
		return err
	}

	//监听回调
	go func() {
		timer := time.After(time.Second * 60)
		v := c.Value("closer")
		var ok bool
		var closer io.Closer
		if closer, ok = v.(io.Closer); !ok {
			return
		}
		defer closer.Close() //回调结束关闭tracer

		for {
			select {
			case <-timer:
				//调用超时，更新record状态为-1
				zap.S().Info("调用超时,内部错误,更新record状态为-1")
				_, err := global.RecordSrvClient.UpdateRecord(context.WithValue(context.Background(), "ginContext", c), &proto.UpdateRecordRequest{
					ID:     recordMsg.ID,
					Status: -1,
				})
				if err != nil {
					zap.S().Info("更新record状态为-1出错")
					return
				}
				return
			case d, ok := <-msgs:
				if !ok {
					zap.S().Info("通道关闭")
					return
				}

				if d.CorrelationId == msgID {
					zap.S().Info("接收到回调,判题已完成,向redis写入当前recordID，向数据库写入redis更新信息")
					// 解析回调信息
					var msgReply global.MsgReply
					err := json.Unmarshal(d.Body, &msgReply)
					if err != nil {
						zap.S().Info("解析回调信息出错")
						continue
					}
					// 更新record信息
					_, err = global.RecordSrvClient.UpdateRecord(context.WithValue(context.Background(), "ginContext", c), &proto.UpdateRecordRequest{
						ID:        msgReply.ID,
						Status:    msgReply.Status,
						ErrCode:   msgReply.ErrCode,
						ErrMsg:    msgReply.ErrMsg,
						MemUsage:  msgReply.MemUsage,
						TimeUsage: msgReply.TimeUsage,
					})
					if err != nil {
						zap.S().Infof("record信息更新失败  %s", err.Error())
						return
					}
					// 向redis写值，通知信息更新
					redisSpan := opentracing.StartSpan("updateRedis", opentracing.ChildOf(parentSpan.Context()))
					_, err = global.Redis.Set(context.Background(), strconv.Itoa(int(msgReply.ID)), 1, time.Second*60).Result()
					redisSpan.Finish()
					if err != nil {
						return
					}
					_ = d.Ack(false)
					return
				}
			}
		}
	}()
	return nil
}

// Retry 重试（提交后记录创建成功但是未判题，status=-1）
func Retry(c *gin.Context) {
	span := c.Value("parentSpan")
	var ok bool
	var parentSpan opentracing.Span
	if parentSpan, ok = span.(opentracing.Span); !ok {
		c.JSON(http.StatusInternalServerError, "内部错误")
		return
	}
	// 解析id
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "参数错误")
		return
	}

	// 检查记录是否存在
	record, err := global.RecordSrvClient.GetRecordByID(context.Background(), &proto.IDRequest{Id: int32(id)})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	//记录的uid必须和当前用户的uid相等
	// 验证身份,只能以自己的uid来提交
	claims, exist := c.Get("claims")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "没有权限",
		})
		return
	}
	customClaims := claims.(*model.CustomClaims)
	if int32(customClaims.ID) != record.UID {
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "没有权限",
		})
		return
	}

	// 记录状态只能是-1（超时）
	if record.Status != -1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "该记录无需重试",
		})
		return
	}

	// 重新将该recordID放进mq中
	recordMsg := global.MsgSend{
		ID:         record.ID,
		Lang:       record.Lang,
		SubmitCode: record.SubmitCode,
		MemLimit:   record.MemLimit,
		TimeLimit:  record.TimeLimit,
		QID:        record.QID,
	}
	sendSpan := opentracing.StartSpan("send2MQ", opentracing.ChildOf(parentSpan.Context()))
	err = Send2MQ(c, recordMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "内部错误")
		return
	}
	sendSpan.Finish()
	c.JSON(http.StatusOK, gin.H{
		"msg": "重新尝试提交该记录",
	})
}
