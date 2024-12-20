package user

import (
	"take-out/global"
	"take-out/internal/api/controller"
	"take-out/internal/repository/dao"
	"take-out/internal/service"
	"take-out/middle"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (u *UserRouter) InitApiRouter(rg *gin.RouterGroup) {
	r := rg.Group("/user")
	authR := r.Group("/user")
	// 私有路由使用jwt验证
	authR.Use(middle.VerifyJWTUser())
	userCtl := controller.NewUserController(service.NewUserService(
		dao.NewUserRepo(global.DB),
	))
	{
		// 用户登录
		r.POST("/login", userCtl.Login)
		// 退出登录
		authR.POST("/logout", userCtl.Logout)
	}
}
