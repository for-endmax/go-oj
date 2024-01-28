package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"question_web/global"
	"question_web/handler"
	"question_web/mid"
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

	// 题目路由
	questionGroup := global.GinEngine.Group("question").Use(mid.JWTAuth()).Use(mid.TraceIm())
	{
		questionGroup.GET("/list", handler.GetQuestionList)
		questionGroup.POST("/add", handler.AddQuestion)
		questionGroup.GET("/info", handler.GetQuestionInfo)
		questionGroup.POST("/del", handler.DelQuestion)
		questionGroup.POST("/update", handler.UpdateQuestion)
	}

	// 测试信息路由
	TestGroup := global.GinEngine.Group("test").Use(mid.JWTAuth()).Use(mid.TraceIm())
	{
		TestGroup.GET("/list", handler.GetTestInfo)
		TestGroup.POST("/add", handler.AddTestInfo)
		TestGroup.POST("/del", handler.DelTestInfo)
		TestGroup.POST("/update", handler.UpdateTestInfo)
	}
}
