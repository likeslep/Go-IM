package user

import (
	"errors"
	"server/pkg"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Register 注册业务
func Register(req *RegisterReq) error {
	// 1. 检查用户名是否存在
	_, err := GetUserByUsername(req.Username)
	if err == nil {
		return errors.New("用户名已存在")
	}

	// 2. 密码加密
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		zap.S().Errorf("密码加密失败: %v", err)
		return err
	}

	// 3. 构建用户
	user := &User{
		Username: req.Username,
		Password: string(hashPwd),
		Nickname: req.Nickname,
		Status:   2, // 默认离线
	}

	// 4. 入库
	return CreateUser(user)
}

// Login 登录业务，返回token
func Login(req *LoginReq) (string, error) {
	// 1. 查询用户
	user, err := GetUserByUsername(req.Username)
	if err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 2. 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 3. 生成JWT
	token, err := pkg.GenerateToken(user.ID)
	if err != nil {
		zap.S().Errorf("生成token失败: %v", err)
		return "", err
	}

	return token, nil
}