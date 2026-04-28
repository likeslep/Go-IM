package conversation

import (
	"errors"
	"server/global"

	"gorm.io/gorm"
)

// List 获取会话列表
func List(uid int64) ([]Conversation, error) {
	var list []Conversation
	err := global.DB.
		Where("uid = ?", uid).
		Order("updated_at desc").
		Find(&list).Error
	return list, err
}

// UpdateConversation 发送消息时，更新会话
func UpdateConversation(fromUID, toUID int64, content string) error {
	// 给发送方更新
	if err := upsertOne(fromUID, toUID, content); err != nil {
		return err
	}
	// 给接收方更新 + 未读+1
	return upsertOneUnread(toUID, fromUID, content)
}

// upsertOne 不存在则创建，存在则更新
func upsertOne(uid, toUID int64, lastMsg string) error {
	var conv Conversation
	err := global.DB.
		Where("uid = ? AND to_uid = ?", uid, toUID).
		First(&conv).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 不存在 → 创建
		return global.DB.Create(&Conversation{
			UID:     uid,
			ToUID:   toUID,
			LastMsg: lastMsg,
			Unread:  0,
		}).Error
	}

	// 存在 → 更新最后一条消息
	return global.DB.Model(&conv).
		Updates(map[string]interface{}{
			"last_msg":   lastMsg,
			"updated_at": gorm.Expr("now()"),
		}).Error
}

// upsertOneUnread 接收方会话：未读+1
func upsertOneUnread(uid, toUID int64, lastMsg string) error {
	var conv Conversation
	err := global.DB.
		Where("uid = ? AND to_uid = ?", uid, toUID).
		First(&conv).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return global.DB.Create(&Conversation{
			UID:     uid,
			ToUID:   toUID,
			LastMsg: lastMsg,
			Unread:  1,
		}).Error
	}

	// 未读 +1
	return global.DB.Model(&conv).
		Updates(map[string]interface{}{
			"last_msg":   lastMsg,
			"unread":     gorm.Expr("unread + 1"),
			"updated_at": gorm.Expr("now()"),
		}).Error
}

// ClearUnread 进入聊天清空未读
func ClearUnread(uid, toUID int64) error {
	return global.DB.Model(&Conversation{}).
		Where("uid = ? AND to_uid = ?", uid, toUID).
		Update("unread", 0).Error
}
