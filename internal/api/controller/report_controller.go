package controller

import (
	"errors"
	"net/http"
	"strings"
	"take-out/common"
	"take-out/common/e"
	"take-out/global"
	"take-out/internal/api/request"
	"take-out/internal/api/response"
	"take-out/internal/repository"
	"take-out/internal/repository/dao"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	user  repository.UserRepo
	order repository.OrderRepo
}

func NewReportController() ReportController {
	return ReportController{
		user:  dao.NewUserRepo(global.DB),
		order: dao.NewOrderDao(),
	}
}

func (rt ReportController) ProduceDate(c *gin.Context) ([]string, error) {
	startDate, err := time.Parse("2006-01-02", c.Query("begin"))
	if err != nil {
		return nil, errors.New("解析开始日期失败")
	}

	endDate, err := time.Parse("2006-01-02", c.Query("end"))
	if err != nil {
		return nil, errors.New("解析结束日期失败")
	}

	if endDate.Before(startDate) {
		return nil, errors.New("结束日期必须在开始日期之后")
	}

	var dates []string
	for current := startDate; !current.After(endDate); current = current.AddDate(0, 0, 1) {
		dates = append(dates, current.Format("2006-01-02"))
	}

	return dates, nil
}

func (rt ReportController) DateRange(c *gin.Context) request.ReportQuestDTO {
	dto := request.ReportQuestDTO{
		Begin: c.Query("begin"),
		End:   c.Query("end") + " 23:59:59", // 要添加上不然查不到
	}
	return dto
}

func (rt ReportController) Top(c *gin.Context) {
	dto := rt.DateRange(c)
	orders, _ := rt.order.OrderTop(dto)
	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: orders})

}

func (rt ReportController) UserStatistics(c *gin.Context) {
	dto := rt.DateRange(c)
	everyDay, _ := rt.user.UserReport(dto)
	dates, err := rt.ProduceDate(c)
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
		return
	}

	var result response.UserReportVO
	result.DateList = strings.Join(dates, ",")
	for _, date := range dates {
		o, ok := GetLocalOrderVO[response.LocalUsertVO](date, everyDay)
		if ok {
			result.NewUserList += o.NewUserCount + ","
			result.TotalUserList += o.TotalUserCount + ","
		} else {
			result.NewUserList += "0,"
			result.TotalUserList += "0,"
		}
	}

	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: result})
}

func (rt ReportController) TurnoverStatistics(c *gin.Context) {
	dto := rt.DateRange(c)
	everyDay, _ := rt.order.OrderTurnover(dto)
	dates, err := rt.ProduceDate(c)
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
		return
	}

	var result response.TurnoverReportVO
	result.DateList = strings.Join(dates, ",")
	for _, date := range dates {
		o, ok := GetLocalOrderVO[response.LocalTurnoverVO](date, everyDay)
		if ok {
			result.TurnoverList += o.TurnoverCount + ","
		} else {
			result.TurnoverList += "0,"
		}
	}

	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: result})
}

func GetLocalOrderVO[T response.ILocalOrder](date string, ds []T) (T, bool) {
	var lv T
	for _, d := range ds {
		if d.GetDate() == date {
			lv = d
			return lv, true
		}
	}
	return lv, false
}

func (rt ReportController) OrdersStatistics(c *gin.Context) {
	dto := rt.DateRange(c)
	orderNuberReportVO, _ := rt.order.QueryOrderNumber(dto)
	dates, err := rt.ProduceDate(c)
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
		return
	}
	var result response.OrderReportVO
	result.DateList = strings.Join(dates, ",")
	result.OrderNuberReportVO = orderNuberReportVO

	everyDay, _ := rt.order.OrderReport(dto)
	for _, date := range dates {
		o, ok := GetLocalOrderVO[response.LocalOrderVO](date, everyDay)
		if ok {
			result.OrderCountList += o.TotalOrderCount + ","
			result.ValidOrderCountList += o.ValidOrderCount + ","
		} else {
			result.OrderCountList += "0,"
			result.ValidOrderCountList += "0,"
		}
	}
	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: result})
}
