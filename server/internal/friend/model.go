package friend

import "time"

// Friend 好友关系表（双方成为好友后存在）
type Friend struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UID       int64     `gorm:"index;not null"` // 用户A
	FID       int64     `gorm:"index;not null"` // 用户B
	CreatedAt time.Time `json:"created_at"`
}

// FriendApply 好友申请表
type FriendApply struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	FromUID   int64     `gorm:"index;not null"` // 申请人
	ToUID     int64     `gorm:"index;not null"` // 被申请人
	Status    int       `gorm:"default:0"`      // 0待处理 1已同意 2已拒绝
	CreatedAt time.Time `json:"created_at"`
}