package main

import (
	"fmt"
	"server/api"
	"server/global"
	"server/initialize"
	"server/internal/group"
	"server/internal/conversation"
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
		&conversation.Conversation{},
		&message.OfflineMessage{},
		&group.Group{},
		&group.GroupMember{},
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
		// 好友
		auth.POST("/apply_friend", api.ApplyFriend)		// 好友申请
		auth.POST("/agree_friend", api.AgreeFriend)		// 通过好友
		auth.GET("/friend_list", api.FriendList)		// 获取好友列表
		auth.GET("/friend_applies", api.FriendApplies)	// 获取好友申请列表
		
		// 会话列表
		auth.GET("/conversation_list", api.ConversationList)		// 获取会话列表
		auth.POST("/clear_unread", api.ClearConversationUnread)		// 清空未读记录

		// 群聊
		auth.POST("/create_group", api.CreateGroup)		// 创建群聊
		auth.POST("/join_group", api.JoinGroup)			// 加入群聊 
	}

	port := global.Config.Server.Port
	zap.S().Info("IM服务启动成功，端口：", port)
	r.Run(fmt.Sprintf(":%d", port))
}
