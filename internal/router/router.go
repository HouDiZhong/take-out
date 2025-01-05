package router

import (
	"take-out/internal/api/websocket"
	"take-out/internal/router/admin"
	"take-out/internal/router/user"
	"take-out/internal/service"

	"github.com/gin-gonic/gin"
)

type RouterGroup interface {
	InitApiRouter(rg *gin.RouterGroup)
}

var AdminRouters = []RouterGroup{
	&admin.EmployeeRouter{},
	&admin.CategoryRouter{},
	&admin.DishRouter{},
	&admin.CommonRouter{},
	&admin.SetMealRouter{},
	&admin.WorkspaceRouter{},
}

var UserRouters = []RouterGroup{
	&user.AddressRouter{},
	// &user.OrderRouter{}, // 由于 OrderRouter 依赖了 websocket 所以在这里不初始化
	&user.UserRouter{},
	&user.CommonRouter{},
	&user.ShopCartRouter{},
}

func InitRouterGroup(r *gin.Engine, hub *service.Hub) *gin.Engine {
	admins := r.Group("/admin")
	InitApiRouterFun(admins, AdminRouters)
	users := r.Group("/user")
	InitApiRouterFun(users, UserRouters)

	// 初始化websocket
	wsCtl := websocket.NewWebSocketHub(hub)
	r.GET("/ws/:id", wsCtl.HandleWebSocket)

	orderRouter := &user.OrderRouter{}
	adminOrderRouter := &admin.OrderRouter{}
	orderRouter.InitApiRouter(users, wsCtl.Hub)
	adminOrderRouter.InitApiRouter(admins, wsCtl.Hub)

	return r
}

func InitApiRouterFun(rg *gin.RouterGroup, routers []RouterGroup) {
	for _, router := range routers {
		router.InitApiRouter(rg)
	}
}
