package friend

import (
	"errors"
	"os/user"
	"server/global"

	"gorm.io/gorm"
)

// Apply 发送好友申请
func Apply(fromUID, toUID int64) error {
	// 1. 不能自己加自己
	if fromUID == toUID {
		return errors.New("不能添加自己为好友")
	}

	// 2. 检查对方用户是否存在
	var count int64
	err := global.DB.Model(&user.User{}).Where("id = ?", toUID).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("对方用户不存在")
	}

	// 3. 检查是否已经是好友
	var friendCount int64
	err = global.DB.Model(&Friend{}).
		Where("uid = ? AND f_id = ?", fromUID, toUID).
		Count(&friendCount).Error
	if err != nil {
		return err
	}
	if friendCount > 0 {
		return errors.New("你们已经是好友")
	}

	// 4. 是否已经发送过申请？
	var cnt int64
	err = global.DB.Model(&FriendApply{}).
		Where("from_uid = ? AND to_uid = ?", fromUID, toUID).
		Count(&cnt).Error
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("已经发送过好友申请")
	}

	// 5. 创建申请
	return global.DB.Create(&FriendApply{
		FromUID: fromUID,
		ToUID:   toUID,
		Status:  0,
	}).Error
}

// Agree 同意好友申请
func Agree(applyID int64, fromUID, toUID int64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 查询申请是否存在 & 状态正确
		var apply FriendApply
		err := tx.Where("id = ? AND status = 0", applyID).First(&apply).Error
		if err != nil {
			return errors.New("申请不存在或已处理")
		}

		// 2. 校验申请人和被申请人是否匹配
		if apply.FromUID != fromUID || apply.ToUID != toUID {
			return errors.New("申请信息不匹配")
		}

		// 3. 更新申请状态
		err = tx.Model(&FriendApply{}).
			Where("id = ?", applyID).
			Update("status", 1).Error
		if err != nil {
			return err
		}

		// 4. 建立双向好友关系
		friends := []Friend{
			{UID: fromUID, FID: toUID},
			{UID: toUID, FID: fromUID},
		}
		return tx.Create(&friends).Error
	})
}

// List 获取好友列表
func List(uid int64) ([]int64, error) {
	var list []int64
	err := global.DB.Model(&Friend{}).
		Where("uid = ?", uid).
		Pluck("f_id", &list).Error
	return list, err
}

// GetApplies 获取收到的好友申请
func GetApplies(uid int64) ([]FriendApply, error) {
	var list []FriendApply
	err := global.DB.Where("to_uid = ? AND status  = 0", uid).
		Order("id desc").Find(&list).Error
	return list, err
}
