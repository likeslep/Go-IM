package chat

import (
	"server/internal/message"

	"go.uber.org/zap"
)

// PullOfflineMessage 用户上线拉取离线消息
func PullOfflineMessage(c *Client) {
	// 1. 获取离线消息
	list, err := message.ListOffline(c.UserID)
	if err != nil {
		zap.S().Errorf("拉取离线消息失败 user:%d, err:%v", c.UserID, err)
		return
	}

	if len(list) == 0 {
		zap.S().Infof("用户[%d] 无离线消息", c.UserID)
		return
	}

	zap.S().Infof("用户[%d] 拉取离线消息 %d 条", c.UserID, len(list))

	// 2. 逐条发送给客户端
	var ids []int64
	for _, msg := range list {
		select {
		case c.Send <- []byte(msg.Content):
		default:
			zap.S().Warnf("离线消息发送失败，客户端阻塞 user:%d", c.UserID)
		}
		ids = append(ids, msg.ID)
	}

	// 3. 发送成功，删除离线消息
	_ = message.DeleteOfflineByIDs(ids)
}
