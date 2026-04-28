package group

import (
	"errors"
	"server/global"
)

// CreateGroup 创建群
func CreateGroup(name string, creatorID int64) (int64, error) {
	gr := &Group{
		Name:      name,
		CreatorID: creatorID,
	}
	err := global.DB.Create(gr).Error
	if err != nil {
		return 0, err
	}

	// 创建者自动入群
	err = global.DB.Create(&GroupMember{
		GroupID: gr.ID,
		UID:     creatorID,
	}).Error
	return gr.ID, err
}

// JoinGroup 加入群
func JoinGroup(groupID, uid int64) error {
	// 判断是否已经在群内
	if IsInGroup(groupID, uid) {
		return errors.New("已在群内")
	}

	return global.DB.Create(&GroupMember{
		GroupID: groupID,
		UID:     uid,
	}).Error
}

// GetGroupMembers 获取群所有成员ID
func GetGroupMembers(groupID int64) ([]int64, error) {
	var uids []int64
	err := global.DB.Model(&GroupMember{}).
		Where("group_id = ?", groupID).
		Pluck("uid", &uids).Error
	return uids, err
}

// IsInGroup 判断用户是否在群里
func IsInGroup(groupID, uid int64) bool {
	var cnt int64
	global.DB.Model(&GroupMember{}).
		Where("group_id = ? AND uid = ?", groupID, uid).
		Count(&cnt)
	return cnt > 0
}