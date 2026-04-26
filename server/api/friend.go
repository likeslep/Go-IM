package api

import (
	"net/http"
	"server/internal/friend"
	"server/pkg"

	"github.com/gin-gonic/gin"
)

// ApplyFriend 发送好友申请
func ApplyFriend(c *gin.Context) {
	uid := pkg.GetUserID(c)
	var req struct {
		ToUID int64 `json:"to_uid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	err := friend.Apply(uid, req.ToUID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "好友申请已发送",
	})
}

// AgreeFriend 同意好友
func AgreeFriend(c *gin.Context) {
	var req struct {
		ApplyID int64 `json:"apply_id"`
		FromUID int64 `json:"from_uid"`
		ToUID   int64 `json:"to_uid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	err := friend.Agree(req.ApplyID, req.FromUID, req.ToUID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "同意成功，现已成为好友",
	})
}

// FriendList 好友列表
func FriendList(c *gin.Context) {
	uid := pkg.GetUserID(c)
	list, err := friend.List(uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "获取失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"list": list,
	})
}

// FriendApplies 申请列表
func FriendApplies(c *gin.Context) {
	uid := pkg.GetUserID(c)
	list, err := friend.GetApplies(uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "获取失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"list": list,
	})
}
