package admin

import (
	"take-out/internal/api/controller"
	"take-out/middle"

	"github.com/gin-gonic/gin"
)

type ReportRouter struct{}

func (ws *ReportRouter) InitApiRouter(rg *gin.RouterGroup) {
	r := rg.Group("report")
	// 私有路由使用jwt验证
	r.Use(middle.VerifyJWTAdmin())
	rtCtl := controller.NewReportController()
	{
		// 销量前十
		r.GET("/top10", rtCtl.Top)
		// 用户统计
		r.GET("/userStatistics", rtCtl.UserStatistics)
		// 营业额统计
		r.GET("/turnoverStatistics", rtCtl.TurnoverStatistics)
		// 订单统计
		r.GET("/ordersStatistics", rtCtl.OrdersStatistics)
		// 导出excel
		r.GET("/export", rtCtl.ExportExcel)
	}
}
