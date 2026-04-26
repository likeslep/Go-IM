package main

import (
	"fmt"
	"server/api"
	"server/global"
	"server/initialize"
	"server/internal/friend"
	"server/internal/message"
	"server/internal/user"
	"server/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. 初始化日志
	initialize.InitLogger()

	// 2. 加载配置
	initialize.InitConfig()

	// 3. 初始化 MySQL
	initialize.InitMySQL()

	err := global.DB.AutoMigrate(
		&user.User{},
		&message.Message{},
		&friend.Friend{},
		&friend.FriendApply{},
	)
	if err != nil {
		zap.S().Fatalf("migrate failed: %v", err)
	}

	r := gin.Default()

	public := r.Group("/api")
	{
		public.POST("/register", api.Register)
		public.POST("/login", api.Login)
		public.GET("/ws", api.WebSocketConnect)
	}

	auth := r.Group("/auth")
	auth.Use(pkg.JWTMiddleWare())
	{
		auth.POST("/apply_friend", api.ApplyFriend)
		auth.POST("/agree_friend", api.AgreeFriend)
		auth.GET("/friend_list", api.FriendList)
		auth.GET("/friend_applies", api.FriendApplies)
	}

	port := global.Config.Server.Port
	zap.S().Info("IM服务启动成功，端口：", port)
	r.Run(fmt.Sprintf(":%d", port))
}
