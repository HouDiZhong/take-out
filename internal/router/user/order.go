package user

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
	r.Use(middle.VerifyJWTUser())
	orderCtl := controller.NewOrderController(hub)
	{
		// 催单
		r.GET("/reminder/:id", orderCtl.Reminder)
		// 再来一单
		r.POST("/repetition/:id", orderCtl.Rpetition)
		// 历史订单查询
		r.GET("/historyOrders", orderCtl.HistoryOrders)
		// 取消订单
		r.PUT("/cancel/:id", orderCtl.Cancel)
		// 订单详情
		r.GET("/orderDetail/:id", orderCtl.OrderDetail)
		// 用户下单
		r.POST("/submit", orderCtl.Submit)
		// 支付
		r.PUT("/payment", orderCtl.Payment)
	}
}
