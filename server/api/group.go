package api

import (
	"github.com/gin-gonic/gin"
	"server/internal/group"
	"server/pkg"
)

// CreateGroup 创建群
func CreateGroup(c *gin.Context) {
	uid := pkg.GetUserID(c)
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	gid, err := group.CreateGroup(req.Name, uid)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "group_id": gid})
}

// JoinGroup 加入群
func JoinGroup(c *gin.Context) {
	uid := pkg.GetUserID(c)
	var req struct {
		GroupID int64 `json:"group_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	err := group.JoinGroup(req.GroupID, uid)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "加入成功"})
}