package admin

import (
	"take-out/internal/api/controller"
	"take-out/internal/service"
	"take-out/middle"

	"github.com/gin-gonic/gin"
)

type OrderRouter struct{}

func (o *OrderRouter) InitApiRouter(rg *gin.RouterGroup, hub *service.Hub) {
	r := rg.Group("order")
	// 私有路由使用jwt验证
	r.Use(middle.VerifyJWTAdmin())
	orderCtl := controller.NewOrderController(hub)
	{
		// 取消订单
		r.PUT("/cancel", orderCtl.AdminCancel)
		// 各个状态订单数量统计
		r.GET("statistics", orderCtl.Statistics)
		// 完成订单
		r.PUT("complete/:id", orderCtl.Complete)
		// 拒单
		r.PUT("/rejection", orderCtl.Rejection)
		// 接单
		r.PUT("/confirm", orderCtl.Confirm)
		// 订单详情
		r.GET("/details/:id", orderCtl.OrderDetail)
		// 派送订单
		r.PUT("/delivery/:id", orderCtl.Delivery)
		// 订单搜索
		r.GET("/conditionSearch", orderCtl.ConditionSearch)
	}
}
