package message

import (
	"time"
)

// Message 和数据库 message 表对应
type Message struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Type      int8      `gorm:"type:tinyint;not null" json:"type"` // 1单聊 2群聊
	FromUid   int64     `gorm:"not null" json:"from_uid"`
	ToUid     int64     `gorm:"not null" json:"to_uid"`
	Content   string    `gorm:"type:text" json:"content"`
	MsgType   int      `gorm:"type:tinyint;not null" json:"msg_type"`
	IsRevoke  int8      `gorm:"type:tinyint;default:0" json:"is_revoke"`
	CreatedAt time.Time `json:"created_at"`
}