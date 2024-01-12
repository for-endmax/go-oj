package handler

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"time"
	"user_web/form"
	"user_web/global"
	"user_web/global/response"
	"user_web/mid"
	"user_web/model"
	"user_web/proto"
)

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
					"msg": "用户已存在",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
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

// Login 用户登录
func Login(c *gin.Context) {
	// 解析用户名和密码
	loginForm := form.LoginForm{}
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		zap.S().Info("解析表单出错")
		HandleValidatorErr(c, err)
		return
	}
	// 查询用户是否存在
	rsp1, err := global.UserSrvClient.GetUserByNickname(context.Background(), &proto.NicknameRequest{Nickname: loginForm.NickName})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "登录失败",
				})
			}
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "登录失败",
			})
		}
	}
	// 验证用户名和密码
	rsp, err := global.UserSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
		Nickname: loginForm.NickName,
		Password: loginForm.Password,
	})
	if err != nil {
		zap.S().Info("验证用户名密码时，user_srv服务出错")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	if !rsp.Valid {
		zap.S().Info("用户名或密码错误")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "用户名或密码错误",
		})
		return
	}
	zap.S().Info("验证成功")
	//登录成功，返回token
	j := mid.NewJWT()
	//负载内容
	Claims := model.CustomClaims{
		uint(rsp1.Id),
		rsp1.Nickname,
		rsp1.Role,
		jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
			Issuer:    "endmax",
		},
	}
	token, err := j.CreateToken(Claims)
	if err != nil {
		zap.S().Infof("[CreateToken] 生成token失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":        rsp1.Id,
		"nickname":  rsp1.Nickname,
		"gender":    rsp1.Gender,
		"token":     token,
		"expiresAt": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	// 验证用户权限
	role, ok := c.Get("userRole")
	if !ok {
		zap.S().Error("从context中获取值出错")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	value, ok := role.(int32)
	if !ok {
		zap.S().Error("从context中获取值出错")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	if int(value) != 2 {
		zap.S().Info("无权限")
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "无权限",
		})
		return
	}
	// 获取参数
	var page string
	var size string
	page = c.DefaultQuery("page", "1")
	size = c.DefaultQuery("size", "5")
	pageInt, _ := strconv.Atoi(page)
	sizeInt, _ := strconv.Atoi(size)
	zap.S().Info("获取用户列表页")

	// 调用接口
	//调用接口
	rsp, err := global.UserSrvClient.GetUserInfoList(context.Background(), &proto.PageInfo{
		PNum:  int32(pageInt),
		PSize: int32(sizeInt),
	})
	if err != nil {
		zap.S().Errorw("GetUserList 调用失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		user := response.UserResponse{
			ID:       value.Id,
			Nickname: value.Nickname,
			Gender:   value.Gender,
			Role:     value.Role,
		}

		result = append(result, user)
	}
	c.JSON(http.StatusOK, result)
}

// SignUp 用户注册
func SignUp(c *gin.Context) {
	// 解析表单
	signupForm := form.SignupForm{}
	if err := c.ShouldBindJSON(&signupForm); err != nil {
		zap.S().Info("解析表单出错")
		HandleValidatorErr(c, err)
		return
	}
	zap.S().Info(signupForm)
	rsp, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Nickname: signupForm.NickName,
		Gender:   signupForm.Gender,
		Role:     1,
		Password: signupForm.Password,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":       rsp.Id,
		"nickname": rsp.Nickname,
		"gender":   rsp.Gender,
		"role":     rsp.Role,
	})
}

// AddUser 管理员添加用户
func AddUser(c *gin.Context) {
	// 验证用户权限
	role, ok := c.Get("userRole")
	if !ok {
		zap.S().Error("从context中获取值出错")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	value, ok := role.(int32)
	if !ok {
		zap.S().Error("从context中获取值出错")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	if int(value) != 2 {
		zap.S().Info("无权限")
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "无权限",
		})
		return
	}

	// 解析信息
	signupForm := form.SignupForm{}
	if err := c.ShouldBindJSON(&signupForm); err != nil {
		zap.S().Info("解析表单出错")
		HandleValidatorErr(c, err)
		return
	}
	zap.S().Info(signupForm)
	rsp, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Nickname: signupForm.NickName,
		Gender:   signupForm.Gender,
		Role:     1,
		Password: signupForm.Password,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":       rsp.Id,
		"nickname": rsp.Nickname,
		"gender":   rsp.Gender,
		"role":     rsp.Role,
	})
}

// UpdateUser 修改用户信息
func UpdateUser(c *gin.Context) {
	// 解析表单
	signupForm := form.SignupForm{}
	if err := c.ShouldBindJSON(&signupForm); err != nil {
		zap.S().Info("解析表单出错")
		HandleValidatorErr(c, err)
		return
	}

	// 验证用户权限
	value1, ok1 := c.Get("userRole")
	value2, ok2 := c.Get("userNickname")
	if !ok1 || !ok2 {
		zap.S().Error("从context中获取值出错")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	role, ok1 := value1.(int32)
	nickname, ok2 := value2.(string)
	if !ok1 || !ok2 {
		zap.S().Error("从context中获取值出错")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}

	if int(role) != 2 && nickname != signupForm.NickName {
		zap.S().Info("无权限")
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "无权限",
		})
		return
	}

	_, err := global.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Nickname: signupForm.NickName,
		Gender:   signupForm.Gender,
		Password: signupForm.Password,
		Role:     role,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新信息成功",
	})

}
