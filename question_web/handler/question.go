package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"question_web/form"
	"question_web/global"
	"question_web/global/response"
	"question_web/proto"
	"strconv"
	"strings"
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
					"msg": "题目已存在",
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

// GetQuestionList 获取题目列表
func GetQuestionList(c *gin.Context) {
	//读取get方法的参数
	pNum, err := strconv.Atoi(c.Query("pn"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "pn 参数错误",
		})
		return
	}
	pSize, err := strconv.Atoi(c.Query("ps"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "ps 参数错误",
		})
		return
	}
	//调用srv方法
	rsp, err := global.QuestionSrvClient.GetQuestionList(context.Background(), &proto.PageInfoRequest{PNum: int32(pNum), PSize: int32(pSize)})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	var data []interface{}
	for _, v := range rsp.Data {
		data = append(data, &response.QuestionBrief{
			ID:   v.Id,
			Seq:  v.Seq,
			Name: v.Name,
		})
	}
	result := gin.H{
		"total": rsp.Total,
		"data":  data,
	}
	c.JSON(http.StatusOK, result)
}

// GetQuestionInfo 通过id获取题目信息
func GetQuestionInfo(c *gin.Context) {
	//获取参数
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "url参数错误",
		})
	}
	//调用srv
	rsp, err := global.QuestionSrvClient.GetQuestionInfo(context.Background(), &proto.GetQuestionInfoRequest{Id: int32(id)})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, response.QuestionResponse{
		ID:      rsp.Id,
		Seq:     rsp.Seq,
		Name:    rsp.Name,
		Content: rsp.Content,
	})

}

// AddQuestion 管理员添加题目
func AddQuestion(c *gin.Context) {
	//验证用户权限
	if !CheckRole(c) {
		return
	}
	//获取post参数
	addQuestionInfo := form.AddQuestionForm{}
	err := c.ShouldBindJSON(&addQuestionInfo)
	if err != nil {
		HandleValidatorErr(c, err)
		return
	}
	//调用rpc
	rsp, err := global.QuestionSrvClient.AddQuestion(context.Background(), &proto.AddQuestionRequest{
		Seq:     addQuestionInfo.Seq,
		Name:    addQuestionInfo.Name,
		Content: addQuestionInfo.Content,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, response.QuestionResponse{ID: rsp.Id,
		Seq:     rsp.Seq,
		Name:    rsp.Name,
		Content: rsp.Content,
	})
}

// DelQuestion 管理员删除题目
func DelQuestion(c *gin.Context) {
	//验证用户权限
	if !CheckRole(c) {
		return
	}
	var req form.DelQuestionForm
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleValidatorErr(c, err)
		return
	}
	_, err := global.QuestionSrvClient.DelQuestion(context.Background(), &proto.DelQuestionRequest{Id: req.ID})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

// UpdateQuestion 管理员修改题目
func UpdateQuestion(c *gin.Context) {
	//验证用户权限
	if !CheckRole(c) {
		return
	}
	//获取post参数
	var req form.UpdateQuestionForm
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleValidatorErr(c, err)
		return
	}
	_, err := global.QuestionSrvClient.UpdateQuestion(context.Background(), &proto.UpdateQuestionRequest{
		Id:      req.ID,
		Seq:     req.Seq,
		Name:    req.Name,
		Content: req.Content,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

////////////////////////////////////////////////

// GetTestInfo 获取qid对应的所有测试信息
func GetTestInfo(c *gin.Context) {
	// 获取get参数
	qID, err := strconv.Atoi(c.Query("q_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	// 调用rpc
	rsp, err := global.QuestionSrvClient.GetTestInfo(context.Background(), &proto.GetTestRequest{QId: int32(qID)})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	// 返回结果
	var testInfos []response.TestInfo
	for _, v := range rsp.Data {
		testInfos = append(testInfos, response.TestInfo{
			QID:       v.QId,
			TimeLimit: v.TimeLimit,
			MemLimit:  v.MemLimit,
			Input:     v.Input,
			Output:    v.Output,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  testInfos,
	})
}

// AddTestInfo 管理员增加测试信息
func AddTestInfo(c *gin.Context) {
	// 权限验证
	if !CheckRole(c) {
		return
	}
	// 读取表单参数
	var addTestForm form.AddTestForm
	if err := c.ShouldBindJSON(&addTestForm); err != nil {
		HandleValidatorErr(c, err)
		return
	}
	// 调用rpc
	_, err := global.QuestionSrvClient.AddTest(context.Background(), &proto.AddTestRequest{
		QId:       addTestForm.QID,
		TimeLimit: addTestForm.TimeLimit,
		MemLimit:  addTestForm.MemLimit,
		Input:     addTestForm.Input,
		Output:    addTestForm.Output,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "添加成功",
	})
}

// DelTestInfo 管理员删除测试信息
func DelTestInfo(c *gin.Context) {
	// 权限验证
	if !CheckRole(c) {
		return
	}
	// 获取form
	var delTestForm form.DelTestForm
	if err := c.ShouldBindJSON(&delTestForm); err != nil {
		HandleValidatorErr(c, err)
		return
	}
	// 调用rpc
	_, err := global.QuestionSrvClient.DelTest(context.Background(), &proto.DelTestRequest{Id: delTestForm.ID})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

// UpdateTestInfo 管理员修改测试信息
func UpdateTestInfo(c *gin.Context) {
	// 权限验证
	if !CheckRole(c) {
		return
	}
	// 获取form
	var updateTestForm form.UpdateTestForm
	if err := c.ShouldBindJSON(&updateTestForm); err != nil {
		HandleValidatorErr(c, err)
		return
	}
	// 调用rpc
	_, err := global.QuestionSrvClient.UpdateTest(context.Background(), &proto.UpdateTestRequest{
		Id:        updateTestForm.ID,
		QId:       updateTestForm.QID,
		TimeLimit: updateTestForm.TimeLimit,
		MemLimit:  updateTestForm.MemLimit,
		Input:     updateTestForm.Input,
		Output:    updateTestForm.Output,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
