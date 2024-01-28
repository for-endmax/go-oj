package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user_web/global"
	"user_web/handler"
	"user_web/mid"
)

func InitRouter() {
	global.GinEngine = gin.Default()
	//跨域
	global.GinEngine.Use(mid.Cors())

	//健康检查
	global.GinEngine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "health",
		})
	})

	// 用户路由
	userGroup := global.GinEngine.Group("user").Use(mid.TraceIm())
	{
		userGroup.GET("/list", mid.JWTAuth(), handler.GetUserList)   //获取用户列表
		userGroup.POST("/login", handler.Login)                      //用户登录
		userGroup.POST("/signup", handler.SignUp)                    // 用户注册
		userGroup.POST("/update", mid.JWTAuth(), handler.UpdateUser) //修改用户信息
		userGroup.POST("/add", mid.JWTAuth(), handler.AddUser)       //管理员添加用户
	}
}
