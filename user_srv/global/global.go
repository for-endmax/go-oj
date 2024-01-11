package global

import (
	"gorm.io/gorm"
	"user_srv/config"
)

var (
	LocalConfig  config.LocalConfig  // 本地配置
	ServerConfig config.ServerConfig // 远程配置
	DB           *gorm.DB            // MySQL对象
)
