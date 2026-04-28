package api

import (
	"net/http"
	"server/internal/chat"
	"server/pkg"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// 升级为 WebSocket
var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
}

// WebSocketConnect 建立连接
func WebSocketConnect(c *gin.Context) {
	// 1. 获取 Token
	token := c.Query("token")
	if token == "" {
		c.JSON(401, gin.H{"msg": "缺少token"})
		return
	}

	// 2. 解析 Token
	claims, err := pkg.ParseToken(token)
	if err != nil {
		c.JSON(401, gin.H{"msg": "token无效"})
		return
	}
	userID := claims.UserID

	// 3. 升级协议
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.S().Errorf("升级ws失败: %v", err)
		return
	}

	// 4. 创建客户端
	client := &chat.Client{
		UserID:    userID,
		Conn:      conn,
		Send:      make(chan []byte, 256),
		Heartbeat: time.Now(),
	}

	// 5. 注册到管理器
	chat.H.Register(userID, client)

	// 6. 启动读写协程
	go client.Write()
	go client.Read()

	// 用户上线，拉去离线消息
	go chat.PullOfflineMessage(client)
}
