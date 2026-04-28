package conversation

import "time"

// Conversation 会话表
type Conversation struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	UID        int64     `gorm:"index;not null"` // 所属用户
	ToUID      int64     `gorm:"not null"`        // 对方用户ID
	Type       int       `gorm:"default:1"`      // 1单聊
	LastMsg    string    `gorm:"type:varchar(512)"` // 最后一条消息
	Unread     int       `gorm:"default:0"`      // 未读数量
	UpdatedAt  time.Time
}