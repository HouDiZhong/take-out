package admin

import (
	"take-out/global"
	"take-out/internal/api/controller"
	"take-out/internal/repository/dao"
	"take-out/internal/service"
	"take-out/middle"

	"github.com/gin-gonic/gin"
)

type EmployeeRouter struct{}

func (er *EmployeeRouter) InitApiRouter(router *gin.RouterGroup) {
	publicRouter := router.Group("employee")
	privateRouter := router.Group("employee")
	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifyJWTAdmin())
	employeeCtl := controller.NewEmployeeController(
		service.NewEmployeeService(
			dao.NewEmployeeDao(global.DB),
		),
	)
	{
		publicRouter.POST("/login", employeeCtl.Login)
		privateRouter.POST("/logout", employeeCtl.Logout)
		privateRouter.POST("", employeeCtl.AddEmployee)
		privateRouter.POST("/status/:status", employeeCtl.OnOrOff)
		privateRouter.PUT("/editPassword", employeeCtl.EditPassword)
		privateRouter.PUT("", employeeCtl.UpdateEmployee)
		privateRouter.GET("/page", employeeCtl.PageQuery)
		privateRouter.GET("/:id", employeeCtl.GetById)
	}
}
