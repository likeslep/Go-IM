package message

import "server/global"

// AddOffline 添加离线消息
func AddOffline(uid, fromUID int64, content string, typ int8) error {
	offline := &OfflineMessage{
		UID:     uid,
		FromUID: fromUID,
		Content: content,
		Type:    typ,
	}
	return global.DB.Create(offline).Error
}

// ListOffline 获取用户离线消息
func ListOffline(uid int64) ([]OfflineMessage, error) {
	var list []OfflineMessage
	err := global.DB.Where("uid = ?", uid).Find(&list).Error
	return list, err
}

// DeleteOfflineByIDs 批量删除离线消息
func DeleteOfflineByIDs(ids []int64) error {
	return global.DB.Where("id IN (?)", ids).Delete(&OfflineMessage{}).Error
}

