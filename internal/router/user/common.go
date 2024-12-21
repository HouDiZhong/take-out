package user

import (
	"take-out/global"
	"take-out/internal/api/controller"
	"take-out/internal/repository/dao"
	"take-out/internal/service"
	"take-out/middle"

	"github.com/gin-gonic/gin"
)

type CommonRouter struct{}

func (c *CommonRouter) InitApiRouter(r *gin.RouterGroup) {
	r.Use(middle.VerifyJWTUser())
	{
		// 获取店铺状态
		shopCtrl := controller.NewShopController(service.NewShopService())
		r.GET("/shop/status", shopCtrl.GetShopStatus)
	}
	{
		// 获取分类 1.菜品分类 2.套餐分类
		categoryCtrl := controller.NewCategoryController(
			service.NewCategoryService(dao.NewCategoryDao(global.DB), dao.NewDishRepo(global.DB)),
		)
		r.GET("/category/list", categoryCtrl.List)
	}
	{
		// 依赖注入
		sCtrl := controller.NewSetMealController(
			service.NewSetMealService(dao.NewSetMealDao(global.DB), dao.NewSetMealDishDao()),
		)
		// 根据分类id查询套餐
		r.GET("/setmeal/list", sCtrl.QueryListById)
		// 根据套餐id查询包含的菜品
		r.GET("/setmeal/dish/:id", sCtrl.SetMealDishById)
	}
	{
		// 依赖注入
		dishCtrl := controller.NewDishController(
			service.NewDishService(
				dao.NewDishRepo(global.DB),
				dao.NewDishFlavorDao(),
				dao.NewSetMealDishDao(),
			),
		)
		// 根据分类id查询菜品
		r.GET("/dish/list", dishCtrl.List)
	}
}
