package service

import (
	"take-out/common/e"
	"take-out/common/enum"
	"take-out/common/utils"
	"take-out/global"
	"take-out/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	Login(openId string) (string, error)
	Logout(c *gin.Context) error
}

type UserServiceImpl struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) UserService {
	return UserServiceImpl{repo: repo}
}

func (u UserServiceImpl) Login(openId string) (string, error) {
	// 跟据openId查询用户信息
	user, err := u.repo.FindByOpenId(openId)
	if err != nil {
		global.Log.Error("查询用户信息失败", "Error", err.Error())
	}
	if user.OpenID == "" {
		err := u.repo.CreateUser(openId)
		if err != nil {
			global.Log.Error("创建用户失败", "Error", err.Error())
			return "", err
		}
		u.Login(openId)
	}
	jwtConfig := global.Config.Jwt.User
	token, err := utils.GenerateToken(uint64(user.ID), jwtConfig)

	return token, err
}

func (u UserServiceImpl) Logout(c *gin.Context) error {
	id, exists := c.Get(enum.CurrentId)

	if exists {
		_, err := utils.DeleteRedisToken(id.(uint64), global.Config.Jwt.Admin.Secret)
		return err
	}

	return e.Error_ACCOUNT_NOT_FOUND
}
