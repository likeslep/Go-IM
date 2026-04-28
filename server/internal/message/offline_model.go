package message

import "time"

// OfflineMessage 离线消息表
// 用户不在线时，消息存入这里
// 用户上线后拉取，然后删除
type OfflineMessage struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UID       int64     `gorm:"index;not null"` // 接收者ID
	FromUID   int64     `gorm:"not null"`       // 发送者ID
	Content   string    `gorm:"type:text"`      // 消息内容JSON
	Type      int8      `gorm:"default:1"`      // 1单聊 2群聊
	CreatedAt time.Time `json:"created_at"`
}