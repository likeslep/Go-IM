package initialize

import (
	"log"
	"server/config"
	"server/global"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("配置文件加载失败：%v", err)
	}

	var cfg config.AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("配置解析失败：%v", err)
	}

	global.Config = &cfg

	log.Println("配置文件加载成功")
}
