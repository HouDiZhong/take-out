package service

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
	"take-out/common"
	"take-out/common/enum"
	"take-out/common/utils"
	"take-out/global"
	"take-out/internal/api/request"
	"take-out/internal/api/response"
	"take-out/internal/model"
	"take-out/internal/repository"
	"take-out/internal/repository/dao"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(c *gin.Context, orderDTO request.OrderDTO) (response.Order, error)
	Reminder(c *gin.Context, orderId uint64) (response.Order, error)
	Broadcast(oc func(string), t int8, o response.Order)
	Rpetition(c *gin.Context, userid uint64) error
	HistoryOrders(c *gin.Context, query request.QueryDTO) (*common.PageResult, error)
	OrderDetail(c *gin.Context, userid uint64) (response.OrderDetail, error)
	Cancel(c *gin.Context, userid uint64) error
	AdminCancel(oid string, uid uint64, o request.OrderStatus) error
	StatisticsOrder(c *gin.Context) (*response.OrderStatusNumber, error)

	ConditionSearch(uid uint64, query request.OrderSearchDTO) (*common.PageResult, error)
}

type OrderServiceImpl struct {
	add      repository.AddressRepo
	shopcart repository.ShopCartRepo
	repo     repository.OrderRepo
	db       *gorm.DB
}

func NewOrderService() OrderService {
	return OrderServiceImpl{
		add:      dao.NewAddressDao(global.DB),
		shopcart: dao.NewShopCartRepo(global.DB),
		repo:     dao.NewOrderDao(),
		db:       global.DB,
	}
}

func (os OrderServiceImpl) CreateOrder(c *gin.Context, orderDTO request.OrderDTO) (response.Order, error) {
	var resOrder response.Order

	if uid, exists := c.Get(enum.CurrentId); exists {
		var addId = strconv.FormatUint(orderDTO.AddressBookId, 10)
		addInfo, _ := os.add.GetAddressById(uid.(uint64), addId)
		if addInfo.ID == 0 {
			return resOrder, errors.New("地址信息为空")
		}
		shopCart, _ := os.shopcart.GetShopCartAll(uid.(uint64))
		if len(shopCart) == 0 {
			return resOrder, errors.New("购物车为空")
		}
		// 开启事务
		return resOrder, os.db.Transaction(func(tx *gorm.DB) error {
			orderDao := dao.NewTXOrderDao(tx)
			EstimatedDeliveryTime, err := time.Parse(enum.TimeLayout, orderDTO.EstimatedDeliveryTime)
			if err != nil {
				slog.Error("预计送达时间格式错误", "Error", err.Error())
				return err
			}
			orderInfo := model.Order{
				UserID:                uid.(uint64),
				Consignee:             addInfo.Consignee,
				Phone:                 addInfo.Phone,
				Address:               addInfo.ProvinceName + addInfo.CityName + addInfo.DistrictName + addInfo.Detail,
				Number:                utils.TimeStampStr(),
				Status:                enum.OrderStatusUnpaid,
				Amount:                orderDTO.Amount,
				Remark:                orderDTO.Remark,
				AddressBookID:         orderDTO.AddressBookId,
				OrderTime:             time.Now(),
				CheckoutTime:          time.Now(),
				CancelTime:            time.Now(),
				DeliveryTime:          time.Now(),
				DeliveryStatus:        orderDTO.DeliveryStatus,
				EstimatedDeliveryTime: EstimatedDeliveryTime,
				PackAmount:            orderDTO.PackAmount,
				TablewareNumber:       orderDTO.TablewareNumber,
				TablewareStatus:       orderDTO.TablewareStatus,
			}
			// 创建订单
			if err := orderDao.CreateOrder(&orderInfo); err != nil {
				slog.Error("订单创建失败", "Error", err.Error())
				return err
			}

			// 创建订单详情
			var orderDetail []model.OrderDetail
			totalAmount := 0.0
			for _, shop := range shopCart {
				orderDetail = append(orderDetail, model.OrderDetail{
					OrderID:    orderInfo.ID,
					DishID:     shop.DishId,
					Number:     shop.Number,
					Amount:     shop.Amount,
					SetmealID:  shop.SetmealId,
					DishFlavor: shop.DishFlavor,
					Name:       shop.Name,
					Image:      shop.Image,
				})
				totalAmount += shop.Amount
			}

			if err := orderDao.BatcheCreateOrder(orderDetail); err != nil {
				slog.Error("批量创建订单详情失败", "Error", err.Error())
				return err
			}

			// 清空购物车
			if err := os.shopcart.CleanShopCart(uid.(uint64)); err != nil {
				slog.Error("购物车清空失败", "Error", err.Error())
				return err
			}

			// 前端数据封装
			resOrder = response.Order{
				ID:          orderInfo.ID,
				OrderAmount: totalAmount,
				OrderNumber: orderInfo.Number,
				OrderTime:   orderInfo.OrderTime.Format(enum.TimeLayout),
			}

			return nil
		})
	}

	return resOrder, nil
}

func (os OrderServiceImpl) Reminder(c *gin.Context, userid uint64) (response.Order, error) {
	// 通过订单id查询订单信息
	orderId := c.Param("id")
	orderInfo, err := os.repo.GetOrderById(userid, orderId)
	if orderInfo.ID == 0 {
		return response.Order{}, errors.New("订单信息为空")
	}
	orders := response.Order{
		ID:          orderInfo.ID,
		OrderNumber: orderInfo.Number,
	}
	return orders, err
}

func (os OrderServiceImpl) Rpetition(c *gin.Context, userid uint64) error {
	// 通过订单id查询订单信息
	orderId := c.Param("id")
	orderInfo, err := os.repo.GetOrderById(userid, orderId)
	if orderInfo.ID == 0 || err != nil {
		return errors.New("订单信息为空")
	}

	err = os.db.Transaction(func(tx *gorm.DB) error {
		orderDao := dao.NewTXOrderDao(tx)
		// 重新下单
		orderInfo.ID = 0
		orderInfo.Number = utils.TimeStampStr()
		orderInfo.Status = enum.OrderStatusUnpaid
		orderInfo.OrderTime = time.Now()

		// 创建订单
		if err := orderDao.CreateOrder(&orderInfo); err != nil {
			slog.Error("订单创建失败", "Error", err.Error())
			return err
		}

		// 创建订单详情
		orderDetail, _ := os.repo.QueryOrderDetailById(orderId)
		for i := range orderDetail {
			orderDetail[i].ID = 0
			orderDetail[i].OrderID = orderInfo.ID
		}

		// 批量创建订单详情
		if err := orderDao.BatcheCreateOrder(orderDetail); err != nil {
			slog.Error("批量创建订单详情失败", "Error", err.Error())
			return err
		}

		return nil
	})

	return err
}

func (os OrderServiceImpl) StatisticsOrder(c *gin.Context) (*response.OrderStatusNumber, error) {
	if uid, exists := c.Get(enum.CurrentId); exists {
		return os.repo.GetStatusNumber(uid.(uint64))
	}
	return nil, nil
}

func (os OrderServiceImpl) HistoryOrders(c *gin.Context, query request.QueryDTO) (*common.PageResult, error) {
	if uid, exists := c.Get(enum.CurrentId); exists {
		return os.repo.QueryOrderList(uid.(uint64), query)
	}
	return nil, nil
}

func (os OrderServiceImpl) ConditionSearch(uid uint64, query request.OrderSearchDTO) (*common.PageResult, error) {
	return os.repo.ConditionSearch(uid, query)
}

func (os OrderServiceImpl) QueryIdNil(oid string, userid uint64) ([]uint64, error) {
	order, err := os.repo.GetOrderById(userid, oid)
	if order.ID == 0 || err != nil {
		return nil, errors.New("订单信息为空")
	}
	return []uint64{order.ID}, nil
}

func (os OrderServiceImpl) Cancel(c *gin.Context, userid uint64) error {
	// 拿到参数订单id
	orderId := c.Param("id")
	// 查询订单是否为空
	ids, err := os.QueryIdNil(orderId, userid)
	if err != nil {
		return err
	}
	o := request.OrderStatus{
		Type:       enum.OrderStatusCancel,
		CancelTime: time.Now(),
	}
	// 取消订单
	return os.repo.UpdateOrderStatus(ids, o)
}

func (os OrderServiceImpl) AdminCancel(oid string, uid uint64, o request.OrderStatus) error {
	ids, err := os.QueryIdNil(oid, uid)
	if err != nil {
		return err
	}
	return os.repo.UpdateOrderStatus(ids, o)
}

func (os OrderServiceImpl) OrderDetail(c *gin.Context, userid uint64) (response.OrderDetail, error) {
	// 通过订单id查询订单信息
	orderId := c.Param("id")
	orderInfo, err := os.repo.GetOrderById(userid, orderId)
	if orderInfo.ID == 0 || err != nil {
		return response.OrderDetail{}, errors.New("订单信息为空")
	}

	// 通过订单id查询订单详情
	oDetail, err := os.repo.QueryOrderDetailById(orderId)
	if err != nil {
		return response.OrderDetail{}, err
	}

	var orderDetail = response.OrderDetail{
		Order:           orderInfo,
		OrderDetailList: oDetail,
	}

	return orderDetail, nil
}

func (os OrderServiceImpl) Broadcast(broadcast func(string), t int8, o response.Order) {
	oInfo := response.Websocket{
		Type:    t,
		OrderId: o.ID,
		Content: "订单号:" + o.OrderNumber,
	}
	if jsInfo, err := json.Marshal(oInfo); err == nil {
		broadcast(string(jsInfo))
	}
}
