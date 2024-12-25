package user

import (
	"take-out/global"
	"take-out/internal/api/controller"
	"take-out/internal/repository/dao"
	"take-out/internal/service"
	"take-out/middle"

	"github.com/gin-gonic/gin"
)

type ShopCartRouter struct{}

func (s ShopCartRouter) InitApiRouter(rg *gin.RouterGroup) {
	r := rg.Group("shoppingCart")
	// 私有路由使用jwt验证
	r.Use(middle.VerifyJWTUser())
	shCtl := controller.NewShopCartConteroller(
		service.NewShopCartService(
			dao.NewDishRepo(global.DB),
			dao.NewSetMealDao(global.DB),
			dao.NewShopCartRepo(global.DB),
		),
	)
	{
		// 添加购物车
		r.POST("/add", shCtl.AddShopCart)
		// 获取购物车列表
		r.GET("/list", shCtl.GetShopCart)
		// 删除购物车
		r.POST("/sub", shCtl.DelShopCart)
		// 清空购物车
		r.DELETE("/clean", shCtl.ClearnShopCart)
	}
}
