package global

import (
	"server/config"

	"gorm.io/gorm"
)

// 全局配置
var Config *config.AppConfig

// 全局数据库
var DB *gorm.DB
