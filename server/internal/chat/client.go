package chat

import (
	"encoding/json"
	"server/global"
	"server/internal/conversation"
	"server/internal/group"
	"server/internal/message"

	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// 全局常量
const (
	// 服务端主动发送Ping间隔
	pingInterval = 5 * time.Second

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

		// 处理消息
		c.handleMessage(&msg)
	}
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


func (c *Client) handleMessage(msg *Message) {
	// 消息先入库
	msgRecord := &message.Message{
		Type:     int8(msg.Type),
		FromUid:  c.UserID,
		ToUid:    msg.ToUserId,
		Content:  msg.Content,
		MsgType:  msg.MsgType,
		IsRevoke: 0,
	}
	global.DB.Create(msgRecord)

	// 更新会话
	err := conversation.UpdateConversation(c.UserID, msg.ToUserId, msg.Content)
	if err != nil {
		zap.S().Errorf("更新会话失败: %v", err)
	}

	// 群聊
	if msg.Type == TypeGroup {
		c.handleGroupMessage(msg)
		return
	}

	// 单聊
	c.handlePrivateMessage(msg)
}

// 处理单聊消息
func (c *Client) handlePrivateMessage(msg *Message) {
	// 1. 判断对方是否在线
	toClient, isOnline := H.GetClient(msg.ToUserId)
	if !isOnline {
		// 对方不在线 → 存入离线消息
		zap.S().Infof("用户[%d] 离线，消息已存入离线表", msg.ToUserId)

		// 保存离线消息
		err := message.AddOffline(
			msg.ToUserId,
			c.UserID,
			string(msg.Content),
			1,
		)
		if err != nil {
			zap.S().Errorf("保存离线消息失败: %v", err)
		}
		return
	}

	// 2. 转发给目标用户
	jsonData, _ := json.Marshal(map[string]interface{}{
		"from_user_id": c.UserID,
		"content":      msg.Content,
		"msg_type":     msg.MsgType,
		"time":         time.Now().Format("15:04:05"),
	})

	toClient.Send <- jsonData
	zap.S().Infof("消息转发成功: %d → %d", c.UserID, msg.ToUserId)
}

// handleGroupMessage 群消息分发
func (c *Client) handleGroupMessage(msg *Message) {
	// 1. 检查是否在群内
	if !group.IsInGroup(msg.ToUserId, c.UserID) {
		zap.S().Warnf("用户 %d 不在群 %d 内", c.UserID, msg.ToUserId)
		return
	}

	// 2. 获取所有群成员
	members, err := group.GetGroupMembers(msg.ToUserId)
	if err != nil {
		zap.S().Errorf("获取群成员失败: %v", err)
		return
	}

	// 3. 组装消息体
	jsonData, _ := json.Marshal(map[string]interface{}{
		"from_user_id": c.UserID,
		"group_id":     msg.ToUserId,
		"content":      msg.Content,
		"type":         2,
		"time":         time.Now().Format("2006-01-02 15:04:05"),
	})

	// 4. 遍历成员：在线推送，不在线存离线
	for _, uid := range members {
		if uid == c.UserID {
			continue
		}

		client, ok := H.GetClient(uid)
		if ok {
			// 在线推送
			client.Send <- jsonData
		} else {
			// 不在线 → 存入离线消息
			err = message.AddOffline(uid, c.UserID, string(jsonData), 2)
			if err != nil {
				zap.S().Errorf("保存离线消息失败: %v", err)
			}
		}
	}

	zap.S().Infof("群消息发送成功 群ID:%d 成员数:%d", msg.ToUserId, len(members))
}
