package controller

import (
	"net/http"
	"take-out/common"
	"take-out/common/e"
	"take-out/internal/api/request"
	"take-out/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	service service.OrderService
}

func NewOrderController() OrderController {
	return OrderController{service: service.NewOrderService()}
}

func (oc OrderController) Reminder(c *gin.Context) {
	code := e.SUCCESS

	c.JSON(http.StatusOK, common.Result{Code: code})
}

func (oc OrderController) Rpetition(c *gin.Context) {

}

func (oc OrderController) HistoryOrders(c *gin.Context) {

}

func (oc OrderController) Cancel(c *gin.Context) {

}

func (oc OrderController) OrderDetail(c *gin.Context) {

}

func (oc OrderController) Submit(c *gin.Context) {
	var orderDTO request.OrderDTO
	if err := c.ShouldBind(&orderDTO); err != nil {
		c.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR, Msg: err.Error()})
		return
	}

	if info, err := oc.service.CreateOrder(c, orderDTO); err != nil {
		c.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR, Msg: err.Error()})
	} else {
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: info})
	}
}

func (oc OrderController) Payment(c *gin.Context) {

}
