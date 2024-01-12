package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"user_web/global"
)

// Run 启动服务
func Run() {
	go func() {
		err := global.GinEngine.Run(fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port))
		if err != nil {
			return
		}
		zap.S().Info("启动服务")
	}()
}
