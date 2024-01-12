package initialize

import "go.uber.org/zap"

// InitLogger 初始化logger，并使用自己的logger替换全局logger，通过S()和L()函数可以安全地访问loggger
func InitLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	zap.S().Info("日志初始化")
}
