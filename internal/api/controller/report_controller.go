package controller

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
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
	"github.com/xuri/excelize/v2"
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

func (rt ReportController) ProduceReportQuest(c *gin.Context) request.ReportQuestDTO {
	return request.ReportQuestDTO{
		Begin: c.Query("begin"),
		End:   c.Query("end"),
	}
}

func (rt ReportController) ProduceDate(dto request.ReportQuestDTO) ([]string, error) {
	startDate, err := time.Parse("2006-01-02", dto.Begin)
	if err != nil {
		return nil, errors.New("解析开始日期失败")
	}

	endDate, err := time.Parse("2006-01-02", dto.End)
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
	dates, err := rt.ProduceDate(rt.ProduceReportQuest(c))
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
	dates, err := rt.ProduceDate(rt.ProduceReportQuest(c))
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
	dates, err := rt.ProduceDate(rt.ProduceReportQuest(c))
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

func (rt ReportController) Overview(users []response.EveryUserVO, orders []response.ExcelVO) (response.ExcelsVO, []response.ExcelsVO) {
	var result []response.ExcelsVO
	totalTurnover, totalValidOrder, totalNewUser, totalOrder := 0.0, 0, 0, 0
	for _, user := range users {
		totalNewUser += user.NewUsers
		order, ok := GetLocalOrderVO[response.ExcelVO](user.Times[:10], orders)
		if ok {
			result = append(result, response.ExcelsVO{
				Date:    order.Times[:10],
				ExcelVO: order,
				EveryUserVO: response.EveryUserVO{
					NewUsers: user.NewUsers,
				},
			})
		} else {
			result = append(result, response.ExcelsVO{
				Date:    user.Times[:10],
				ExcelVO: response.ExcelVO{},
				EveryUserVO: response.EveryUserVO{
					NewUsers: user.NewUsers,
				},
			})
		}
	}

	for _, order := range orders {
		totalOrder += order.TotalOrders
		totalTurnover += order.Turnovers
		totalValidOrder += order.ValidOrders
	}

	return response.ExcelsVO{
		ExcelVO: response.ExcelVO{
			Turnovers:          totalTurnover,
			ValidOrders:        totalValidOrder,
			UnitPrices:         totalTurnover / float64(totalValidOrder),
			OrderStatusNumbers: float64(totalValidOrder) / float64(totalOrder),
		},
		EveryUserVO: response.EveryUserVO{
			NewUsers: totalNewUser,
		},
	}, result
}

func (rt ReportController) ExportExcel(c *gin.Context) {
	f, err := excelize.OpenFile("template/运营数据报表模板.xlsx")
	if err != nil {
		slog.Error("打开文件失败", "error", err.Error())
		return
	}
	dates := request.ReportQuestDTO{
		// 一个月前
		Begin: time.Now().AddDate(0, -1, 0).Format("2006-01-02"),
		// 一天前
		End: time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
	}
	orders, _ := rt.order.BatchBusinessOrder(dates)
	users, _ := rt.user.EveryUserReport(dates)

	var results []response.ExcelsVO
	overview, result := rt.Overview(users, orders)
	dateList, _ := rt.ProduceDate(dates)

	for _, date := range dateList {
		relt, ok := GetLocalOrderVO[response.ExcelsVO](date, result)
		if ok {
			results = append(results, relt)
		} else {
			results = append(results, response.ExcelsVO{
				Date:        date,
				ExcelVO:     response.ExcelVO{},
				EveryUserVO: response.EveryUserVO{},
			})
		}
	}

	// 概览
	f.SetCellValue("Sheet1", "B2", "时间："+dates.Begin+"至"+dates.End)
	f.SetCellValue("Sheet1", "C4", overview.Turnovers)
	f.SetCellValue("Sheet1", "E4", overview.OrderStatusNumbers)
	f.SetCellValue("Sheet1", "G4", overview.NewUsers)
	f.SetCellValue("Sheet1", "C5", overview.ValidOrders)
	f.SetCellValue("Sheet1", "E5", overview.UnitPrices)

	for i, result := range results {
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+8), result.Date)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+8), result.Turnovers)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+8), result.ValidOrders)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+8), result.OrderStatusNumbers)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+8), result.UnitPrices)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(i+8), result.NewUsers)
	}
	fileName := "template/" + dates.Begin + "至" + dates.End + "运营数据报表模板.xlsx"
	if err := f.SaveAs(fileName); err != nil {
		slog.Error("保存文件失败", "error", err.Error())
		return
	}

	c.File(fileName)
}
