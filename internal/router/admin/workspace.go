package admin

import (
	"take-out/internal/api/controller"
	"take-out/middle"

	"github.com/gin-gonic/gin"
)

type WorkspaceRouter struct{}

func (ws *WorkspaceRouter) InitApiRouter(rg *gin.RouterGroup) {
	r := rg.Group("workspace")
	// 私有路由使用jwt验证
	r.Use(middle.VerifyJWTAdmin())
	wsCtl := controller.NewWorkSpaceConteroler()
	{
		// 今日运营数据
		r.GET("/businessData", wsCtl.BusinessData)
		// 套餐总览
		r.GET("/overviewSetmeals", wsCtl.OverviewSetmeals)
		// 菜品总量
		r.GET("/overviewDishes", wsCtl.OverviewDishes)
		// 订单数据
		r.GET("/overviewOrders", wsCtl.OverviewOrders)
	}
}
