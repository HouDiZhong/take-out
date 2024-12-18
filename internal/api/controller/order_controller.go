package controller

import (
	"net/http"
	"take-out/common"
	"take-out/common/e"

	"github.com/gin-gonic/gin"
)

type OrderController struct{}

func NewOrderController() OrderController {
	return OrderController{}
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

}

func (oc OrderController) Payment(c *gin.Context) {

}
