package main

import "go.uber.org/zap"

// 要实现的目标
// 		1.使用zap库，初始化一个logger
// 参考资料:https://www.liwenzhou.com/posts/Go/zap/

// InitLogger 初始化logger，并使用自己的logger替换全局logger，通过S()和L()函数可以安全地访问loggger
func InitLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

func main() {
	InitLogger()
	zap.S().Info("配置日志成功")
}
