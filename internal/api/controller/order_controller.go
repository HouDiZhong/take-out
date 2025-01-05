package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"take-out/common"
	"take-out/common/e"
	"take-out/common/enum"
	"take-out/internal/api/request"
	"take-out/internal/service"
	"time"

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
	if uid, exists := c.Get(enum.CurrentId); exists {
		order, err := oc.service.Reminder(c, uid.(uint64))
		if err != nil {
			c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		oc.service.Broadcast(oc.notice.BroadcastMessage, enum.BroadcastRemind, order)
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS})
	}
}

func (oc OrderController) Rpetition(c *gin.Context) {
	if uid, exists := c.Get(enum.CurrentId); exists {
		err := oc.service.Rpetition(c, uid.(uint64))
		if err != nil {
			c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS})
	}
}

func (oc OrderController) HistoryOrders(c *gin.Context) {
	var query = request.QueryDTO{
		Page:     c.DefaultQuery("page", "1"),
		Pagesize: c.DefaultQuery("pagesize", "10"),
		Status:   c.DefaultQuery("status", ""),
	}
	result, err := oc.service.HistoryOrders(c, query)
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: result})
}

func (oc OrderController) Cancel(c *gin.Context) {
	if uid, exists := c.Get(enum.CurrentId); exists {
		if err := oc.service.Cancel(c, uid.(uint64)); err != nil {
			c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: "取消成功"})
	}
}

func (oc OrderController) Complete(c *gin.Context) {
	if uid, exists := c.Get(enum.CurrentId); exists {
		oid := c.Param("id")
		o := request.OrderStatus{
			Status: enum.OrderStatusFinish,
		}
		if err := oc.service.AdminCancel(oid, uid.(uint64), o); err != nil {
			c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: "订单成功"})
	}
}

func (oc OrderController) Delivery(c *gin.Context) {
	if uid, exists := c.Get(enum.CurrentId); exists {
		oid := c.Param("id")
		o := request.OrderStatus{
			Status: enum.OrderStatusSend,
		}
		if err := oc.service.AdminCancel(oid, uid.(uint64), o); err != nil {
			c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: "订单派送中"})
	}
}

func (oc OrderController) Confirm(c *gin.Context) {
	if uid, exists := c.Get(enum.CurrentId); exists {
		var cDTO request.ConfirmDTO
		if err := c.ShouldBind(&cDTO); err != nil {
			c.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		o := request.OrderStatus{
			Status: enum.OrderStatusAccept,
		}
		oid := strconv.Itoa(int(cDTO.Id))
		if err := oc.service.AdminCancel(oid, uid.(uint64), o); err != nil {
			c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: "成功接单"})
	}
}

func (oc OrderController) AdminCancel(c *gin.Context) {
	if uid, exists := c.Get(enum.CurrentId); exists {
		var cDTO request.CancelDTO
		if err := c.ShouldBind(&cDTO); err != nil {
			c.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		o := request.OrderStatus{
			Status:       enum.OrderStatusCancel,
			CancelTime:   time.Now(),
			CancelReason: cDTO.CancelReason,
		}
		oid := strconv.Itoa(int(cDTO.Id))
		if err := oc.service.AdminCancel(oid, uid.(uint64), o); err != nil {
			c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: "取消成功"})
	}
}

func (oc OrderController) Statistics(c *gin.Context) {
	reslut, err := oc.service.StatisticsOrder(c)
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: reslut})
}

func (oc OrderController) ConditionSearch(c *gin.Context) {
	if uid, exists := c.Get(enum.CurrentId); exists {
		oDTO := request.OrderSearchDTO{
			BeginTime: c.DefaultQuery("beginTime", ""),
			EndTime:   c.DefaultQuery("endTime", ""),
			Phone:     c.DefaultQuery("phone", ""),
			Status:    c.DefaultQuery("status", ""),
			Number:    c.DefaultQuery("number", ""),
			Page:      c.DefaultQuery("page", "1"),
			Pagesize:  c.DefaultQuery("pagesize", "10"),
		}
		result, err := oc.service.ConditionSearch(uid.(uint64), oDTO)
		fmt.Printf("----------------------------------------------uid: %v\n", uid)
		if err != nil {
			c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: result})
	}
}

func (oc OrderController) Rejection(c *gin.Context) {
	if uid, exists := c.Get(enum.CurrentId); exists {
		var rDTO request.RejectionDTO
		if err := c.ShouldBind(&rDTO); err != nil {
			c.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		o := request.OrderStatus{
			Status:          enum.OrderStatusCancel,
			CancelTime:      time.Now(),
			RejectionReason: rDTO.RejectionReason,
		}
		oid := strconv.Itoa(int(rDTO.Id))
		if err := oc.service.AdminCancel(oid, uid.(uint64), o); err != nil {
			c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: "取消成功"})
	}
}

func (oc OrderController) OrderDetail(c *gin.Context) {
	if uid, exists := c.Get(enum.CurrentId); exists {
		result, err := oc.service.OrderDetail(c, uid.(uint64))
		if err != nil {
			c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: result})
	}
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
	/* oInfo := response.Websocket{
		Type:    int8(1),
		OrderId: uint64(12),
		Content: "订单号: 2323434",
	}
	if jsInfo, err := json.Marshal(oInfo); err == nil {
		oc.notice.BroadcastMessage(string(jsInfo))
	} */
}
