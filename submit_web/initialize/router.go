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

	// 提交记录路由
	recordGroup := global.GinEngine.Group("/record").Use(mid.JWTAuth()).Use(mid.TraceIm())
	{
		recordGroup.GET("/list", handler.GetRecordListByUID)
		recordGroup.GET("/info", handler.GetRecordByID)
	}
	// 提交路由
	submitGroup := global.GinEngine.Group("/record").Use(mid.JWTAuth()).Use(mid.Trace())
	{
		submitGroup.POST("/submit", handler.Submit)
		submitGroup.GET("/retry", handler.Retry)
	}
}
