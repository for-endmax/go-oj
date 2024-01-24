package initialize

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"submit_web/global"
)

func InitRedis() {
	redisInfo := global.ServerConfig.RedisInfo
	global.Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisInfo.Host, redisInfo.Port), // Redis 服务器地址
		Password: "",                                                   // Redis 密码
		DB:       0,                                                    // 使用默认的数据库
	})
	// 测试连接
	_, err := global.Redis.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	zap.S().Info("初始化redis成功")
}
