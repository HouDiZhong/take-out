package service

import (
	"errors"
	"strconv"
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
				EstimatedDeliveryTime: time.Now(),
				PackAmount:            orderDTO.PackAmount,
				TablewareNumber:       orderDTO.TablewareNumber,
				TablewareStatus:       orderDTO.TablewareStatus,
			}
			// 创建订单
			if err := orderDao.CreateOrder(&orderInfo); err != nil {
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
				return err
			}

			// 清空购物车
			if err := os.shopcart.CleanShopCart(uid.(uint64)); err != nil {
				return err
			}

			// 返回前端数据
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
