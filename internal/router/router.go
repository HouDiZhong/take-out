package router

import (
	"take-out/internal/router/admin"
	"take-out/internal/router/user"

	"github.com/gin-gonic/gin"
)

type RouterGroup interface {
	InitApiRouter(rg *gin.RouterGroup)
}

var AdminRouter = []RouterGroup{
	&admin.EmployeeRouter{},
	&admin.CategoryRouter{},
	&admin.DishRouter{},
	&admin.CommonRouter{},
	&admin.SetMealRouter{},
}

var UserRouter = []RouterGroup{
	&user.AddressRouter{},
	&user.OrderRouter{},
	&user.UserRouter{},
}

func InitRouterGroup(r *gin.Engine) *gin.Engine {
	admin := r.Group("/admin")
	InitApiRouterFun(admin, AdminRouter)
	user := r.Group("/user")
	InitApiRouterFun(user, UserRouter)

	return r
}

func InitApiRouterFun(rg *gin.RouterGroup, routers []RouterGroup) {
	for _, router := range routers {
		router.InitApiRouter(rg)
	}
}

/* type RouterGroup struct {
	admin.EmployeeRouter
	admin.CategoryRouter
	admin.DishRouter
	admin.CommonRouter
	admin.SetMealRouter

	user.AddressRouter
	user.OrderRouter
}

var AllRouter = new(RouterGroup) */
