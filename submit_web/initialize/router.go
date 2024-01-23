package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"submit_web/global"
	"submit_web/handler"
	"submit_web/mid"
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

	// 提交路由
	recordGroup := global.GinEngine.Group("/record").Use(mid.JWTAuth())
	{
		recordGroup.GET("/list", handler.GetRecordListByUID)
		recordGroup.GET("/info", handler.GetRecordByID)
		recordGroup.POST("/submit", handler.Submit)

		// 接受状态改变通知
		recordGroup.GET("/done", handler.Done)
	}

}
