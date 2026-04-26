package user

import (
	"server/global"
)

// CreateUser 创建用户
func CreateUser(user *User) error {
	return global.DB.Create(user).Error
}

// GetUserByUsername 根据用户名查询用户
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := global.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID查询用户
func GetUserByID(id int64) (*User, error) {
	var user User
	err := global.DB.First(&user, id).Error
	return &user, err
}