package user

import (
	"take-out/global"
	"take-out/internal/api/controller"
	"take-out/internal/repository/dao"
	"take-out/internal/service"
	"take-out/middle"

	"github.com/gin-gonic/gin"
)

type AddressRouter struct{}

func (ar *AddressRouter) InitApiRouter(rg *gin.RouterGroup) {
	r := rg.Group("addressBook")
	r.Use(middle.VerifyJWTUser())
	ctl := controller.NewAddressConteroller(
		service.NewAddressService(
			dao.NewAddressDao(global.DB),
		),
	)
	{
		// 新增地址
		r.POST("", ctl.CreateAddress)
		// 当前用户所有地址
		r.GET("/list", ctl.AddressList)
		// 查询默认地址
		r.GET("/default", ctl.DefaultAddress)
		// 通过id修改地址
		r.PUT("", ctl.EditAddressById)
		// 通过id删除地址
		r.DELETE("", ctl.DeleteAddressById)
		// 通过id查询地址
		r.GET("/:id", ctl.GetAddressById)
		// 设置默认地址
		r.PUT("/default", ctl.SetDefaultAddress)
	}
}
