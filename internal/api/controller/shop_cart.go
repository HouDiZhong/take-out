package controller

import (
	"net/http"
	"take-out/common"
	"take-out/common/e"
	"take-out/common/enum"
	"take-out/internal/api/request"
	"take-out/internal/service"

	"github.com/gin-gonic/gin"
)

type ShopCartConteroller struct {
	service service.ShopCartService
}

func NewShopCartConteroller(service service.ShopCartService) ShopCartConteroller {
	return ShopCartConteroller{service: service}
}

func (s ShopCartConteroller) GetShopCart(c *gin.Context) {
	if uid, exists := c.Get(enum.CurrentId); exists {
		shopCarts, err := s.service.GetShopCartAll(uid.(uint64))
		if err != nil {
			c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: shopCarts})
		return
	}
	c.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR})
}

func (s ShopCartConteroller) AddShopCart(c *gin.Context) {
	var shopCartDTO request.ShopCartDTO
	if err := c.ShouldBind(&shopCartDTO); err != nil {
		c.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR, Msg: err.Error()})
		return
	}
	uid, _ := c.Get(enum.CurrentId)
	if err := s.service.AddShopCart(c, uid.(uint64), shopCartDTO); err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
}

func (s ShopCartConteroller) DelShopCart(c *gin.Context) {
	var shopCartDTO request.ShopCartDTO
	if err := c.ShouldBind(&shopCartDTO); err != nil {
		c.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR, Msg: err.Error()})
		return
	}
	if uid, exists := c.Get(enum.CurrentId); exists {
		if err := s.service.DelShopCart(uid.(uint64), shopCartDTO); err != nil {
			c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
		return
	}
}

func (s ShopCartConteroller) ClearnShopCart(c *gin.Context) {
	if uid, exists := c.Get(enum.CurrentId); exists {
		if err := s.service.ClearnShopCart(uid.(uint64)); err != nil {
			c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
		return
	}
	c.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR})
}
