package api

import (
	"server/internal/user"
	"server/pkg"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req user.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, pkg.Fail("参数错误: "+err.Error()))
		return
	}

	if err := user.Register(&req); err != nil {
		c.JSON(http.StatusOK, pkg.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, pkg.Success("注册成功"))
}

// Login 用户登录
func Login(c *gin.Context) {
	var req user.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, pkg.Fail("参数错误"))
		return
	}

	token, err := user.Login(&req)
	if err != nil {
		c.JSON(http.StatusOK, pkg.Fail(err.Error()))
		return
	}

	zap.S().Infof("用户登录成功 username: %s", req.Username)
	c.JSON(http.StatusOK, pkg.Success(token))
}
