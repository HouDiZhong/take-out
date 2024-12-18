package initialize

import (
	"take-out/internal/router"

	"github.com/gin-gonic/gin"
)

func routerInit() *gin.Engine {
	r := gin.Default()
	// allRouter := router.AllRouter
	// admin
	/* admin := r.Group("/admin")
	router.InitAdminRouterGroup(admin)
	// user
	user := r.Group("/user")
	router.InitUserRouterGroup(user) */

	/* {
		allRouter.EmployeeRouter.InitApiRouter(admin)
		allRouter.CategoryRouter.InitApiRouter(admin)
		allRouter.DishRouter.InitApiRouter(admin)
		allRouter.CommonRouter.InitApiRouter(admin)
		allRouter.SetMealRouter.InitApiRouter(admin)
	} */
	return router.InitRouterGroup(r)
}
