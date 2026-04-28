package api

import (
	"github.com/gin-gonic/gin"
	"server/internal/conversation"
	"server/pkg"
	"net/http"
)

// ConversationList 会话列表（首页最近聊天）
func ConversationList(c *gin.Context) {
	uid := pkg.GetUserID(c)

	list, err := conversation.List(uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "获取会话失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"list": list,
	})
}

// ClearConversationUnread 清空未读
func ClearConversationUnread(c *gin.Context) {
	uid := pkg.GetUserID(c)

	var req struct {
		ToUID int64 `json:"to_uid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	err := conversation.ClearUnread(uid, req.ToUID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "清空失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "清空未读成功",
	})
}