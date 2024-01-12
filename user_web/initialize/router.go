package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user_web/global"
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
	userGroup := global.GinEngine.Group("user")
	{
		userGroup.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"ping": "pong",
			})
		})
	}
}
