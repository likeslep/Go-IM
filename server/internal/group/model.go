package group

import "time"

// Group 群表
type Group struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:64;not null"`  // 群名称
	CreatorID int64     `gorm:"not null"`          // 创建者ID
	CreatedAt time.Time `json:"created_at"`
}

// GroupMember 群成员表
type GroupMember struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	GroupID   int64     `gorm:"index;not null"`
	UID       int64     `gorm:"index;not null"`
	CreatedAt time.Time `json:"created_at"`
}