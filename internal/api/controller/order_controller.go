package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"take-out/common"
	"take-out/common/e"
	"take-out/internal/api/request"
	"take-out/internal/api/response"
	"take-out/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	service service.OrderService
	notice  *service.NoticeService
}

func NewOrderController(hub *service.Hub) OrderController {
	return OrderController{
		service: service.NewOrderService(),
		notice:  service.NewNoticeService(hub),
	}
}

func (oc OrderController) Reminder(c *gin.Context) {
	// code := e.SUCCESS
	oInfo := response.Websocket{
		Type:    int8(2),
		OrderId: uint64(12),
		Content: "订单号: 2323434",
	}
	if jsInfo, err := json.Marshal(oInfo); err == nil {
		info := string(jsInfo)
		fmt.Printf("info---------------------------: %v\n", info)
		oc.notice.BroadcastMessage(string(jsInfo))
	}

	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: "成功"})
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
	oInfo := response.Websocket{
		Type:    int8(1),
		OrderId: uint64(12),
		Content: "订单号: 2323434",
	}
	if jsInfo, err := json.Marshal(oInfo); err == nil {
		oc.notice.BroadcastMessage(string(jsInfo))
	}
}
