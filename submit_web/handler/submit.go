package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"submit_web/form"
	"submit_web/global"
	"submit_web/model"
	"submit_web/proto"
	"submit_web/response"
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
		}
		rsp.Data = append(rsp.Data, record)
	}
	// 返回结果
	c.JSON(http.StatusOK, rsp)
}

// GetRecordByID 获取指定id的record的信息
func GetRecordByID(c *gin.Context) {
	//获取参数
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "参数错误")
		return
	}
	// 调用rpc
	recordInfo, err := global.RecordSrvClient.GetRecordByID(context.Background(), &proto.IDRequest{Id: int32(id)})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	//TODO
	// 长连接

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
	}
	c.JSON(http.StatusOK, rsp)
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
	//TODO
	// 分布式事务

	// 调用rpc,生成record
	record, err := global.RecordSrvClient.CreateRecord(context.Background(), &proto.CreateRecordRequest{
		UID:        submitForm.UID,
		QID:        submitForm.QID,
		Lang:       submitForm.Lang,
		SubmitCode: submitForm.SubmitCode,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	//TODO
	// 将recordID放到mq中
	zap.S().Infof("将record放到mq, id : %d", record.ID)

	c.JSON(http.StatusOK, gin.H{
		"msg":       "创建成功",
		"record_id": record.ID,
	})
}

// Done 接受状态改变通知
func Done(c *gin.Context) {
	zap.S().Info("判题结束,状态改变")
	//TODO
	// 接受judge_srv的请求，然后更新状态，使GetRecordByID返回结果
}
