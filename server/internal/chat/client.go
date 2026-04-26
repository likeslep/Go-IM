package chat

import (
	"encoding/json"
	"server/global"
	"server/internal/message"

	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// 全局常量
const (
	// 服务端主动发送Ping间隔
	pingInterval = 5*time.Second 

	// 读写超时时间：超过该时间无任何响应则断开
	readWriteTimeout = 14 * time.Second 
)

// 客户端连接
type Client struct {
	UserID    int64
	Conn      *websocket.Conn
	Send      chan []byte
	Heartbeat time.Time
}

// 读消息：客户端 ——> 服务端
func (c *Client) Read() {
	defer func() {
		_ = c.Conn.Close()
		H.UnRegister(c.UserID)
		zap.S().Infof("客户端关闭: %d", c.UserID)
	}()
	
	// 客户端回复Pong帧时自动触发，刷新超时时间
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(readWriteTimeout))
		_ = c.Conn.SetWriteDeadline(time.Now().Add(readWriteTimeout))
		return nil 
	})


	// 设置读取配置
	c.Conn.SetReadLimit(1024 * 4)
	c.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	for {
		_, msgBytes, err := c.Conn.ReadMessage()
		if err != nil {
			zap.S().Errorf("读取消息错误 user:%d, err:%v", c.UserID, err)
			break
		}

		// 刷新心跳超时时间
		c.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		// 解析消息
		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			zap.S().Errorf("消息解析失败: %v", err)
			continue
		}

		// 处理单聊消息
		c.handlePrivateMessage(&msg)
	}
}

// 处理单聊消息
func (c *Client) handlePrivateMessage(msg *Message) {
	// 1. 保存消息到数据库
	msgRecord := &message.Message{
		Type:     1, // 1单聊
		FromUid:  c.UserID,
		ToUid:    msg.ToUserId,
		Content:  msg.Content,
		MsgType:  msg.MsgType,
		IsRevoke: 0,
	}
	global.DB.Create(&msgRecord)

	// 2. 判断对方是否在线
	toClient, isOnline := H.GetClient(msg.ToUserId)
	if !isOnline {
		zap.S().Infof("用户 %d 不在线，消息已存库", msg.ToUserId)
		return
	}

	// 3. 转发给目标用户
	jsonData, _ := json.Marshal(map[string]interface{}{
		"from_user_id": c.UserID,
		"content":      msg.Content,
		"msg_type":     msg.MsgType,
		"time":         time.Now().Format("15:04:05"),
	})

	toClient.Send <- jsonData
	zap.S().Infof("消息转发成功: %d → %d", c.UserID, msg.ToUserId)
}

// 写消息：服务端 ——> 客户端
func (c *Client) Write() {
	defer func() {
		_ = c.Conn.Close()
	}()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 通道关闭，关闭连接
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					zap.S().Errorf("通道关闭，关闭连接失败：%v", err)
					return 
				}
			}
			_ = c.Conn.SetWriteDeadline(time.Now().Add(readWriteTimeout))
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				zap.S().Errorf("写入客户端消息失败：%v", err)
				return 
			}

		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(readWriteTimeout))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				zap.S().Errorf("用户[%d] 发送Ping心跳失败: %v", c.UserID, err)
				return
			}
		}

	}
}
