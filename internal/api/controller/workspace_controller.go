package controller

import (
	"net/http"
	"take-out/common"
	"take-out/common/e"
	"take-out/global"
	"take-out/internal/api/response"
	"take-out/internal/repository"
	"take-out/internal/repository/dao"

	"github.com/gin-gonic/gin"
)

type WorkSpaceConteroler struct {
	user    repository.UserRepo
	order   repository.OrderRepo
	dish    repository.DishRepo
	setMeal repository.SetMealRepo
}

func NewWorkSpaceConteroler() *WorkSpaceConteroler {
	return &WorkSpaceConteroler{
		user:    dao.NewUserRepo(global.DB),
		order:   dao.NewOrderDao(),
		dish:    dao.NewDishRepo(global.DB),
		setMeal: dao.NewSetMealDao(global.DB),
	}
}

func (ws WorkSpaceConteroler) BusinessData(c *gin.Context) {
	// 订单相关
	orders, _ := ws.order.BusinessOrder()
	// 新增用户
	newUsers, _ := ws.user.GetNewUserNumber()
	reslut := response.BusinessDataVO{
		NewUsers:        newUsers,
		BusinessOrderVO: orders,
	}

	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: reslut})
}

func (ws WorkSpaceConteroler) OverviewSetmeals(c *gin.Context) {
	result, err := ws.setMeal.QuerySetMealDesStatusNumber()
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: result})
}

func (ws WorkSpaceConteroler) OverviewDishes(c *gin.Context) {
	result, err := ws.dish.QuerySetMealDesStatusNumber()
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: result})
}

func (ws WorkSpaceConteroler) OverviewOrders(c *gin.Context) {
	result, err := ws.order.OrderStatusNumber()
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: result})
}
