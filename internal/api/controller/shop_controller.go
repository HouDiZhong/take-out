package controller

import (
	"net/http"
	"take-out/common"
	"take-out/common/e"
	"take-out/internal/service"

	"github.com/gin-gonic/gin"
)

type Shop struct {
	shopService service.ShopService
}

func NewShopController(service service.ShopService) Shop {
	return Shop{shopService: service}
}

func (s Shop) GetShopStatus(c *gin.Context) {
	code := e.SUCCESS
	status := s.shopService.GetShopService(&code)

	c.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: status,
		Msg:  e.GetMsg(code),
	})
}

func (s Shop) SetShopStatus(c *gin.Context) {
	code := e.SUCCESS
	status := c.Param("status")
	s.shopService.SetShopService(status, &code)

	c.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}
