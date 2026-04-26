package user

import (
	"gorm.io/gorm"
	"time"
)

// User 用户表模型（对应数据库 users 表）
type User struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:100;not null" json:"-"` // 不返回密码
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Email     string         `gorm:"size:50" json:"email"`
	Status    int8           `gorm:"default:1" json:"status"` // 1在线 2离线
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// RegisterReq 注册请求参数
type RegisterReq struct {
	Username string `json:"username" binding:"required,min=4,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Nickname string `json:"nickname" binding:"required"`
}

// LoginReq 登录请求参数
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}