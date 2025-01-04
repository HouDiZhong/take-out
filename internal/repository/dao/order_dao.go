package dao

import (
	"strconv"
	"take-out/common"
	"take-out/global"
	"take-out/internal/api/request"
	"take-out/internal/api/response"
	"take-out/internal/model"
	"take-out/internal/repository"
	"time"

	"gorm.io/gorm"
)

type OrderDao struct {
	db *gorm.DB
}

// 普通模式
func NewOrderDao() repository.OrderRepo {
	return &OrderDao{db: global.DB}
}

// 事务模式
func NewTXOrderDao(db *gorm.DB) repository.OrderRepo {
	return &OrderDao{db: db}
}

func (o *OrderDao) CreateOrder(order *model.Order) error {
	return o.db.Create(&order).Error
}

func (o *OrderDao) BatcheCreateOrder(orderDetail []model.OrderDetail) error {
	return o.db.CreateInBatches(&orderDetail, len(orderDetail)).Error
}

func (o *OrderDao) QueryOrderDetailById(oid string) ([]model.OrderDetail, error) {
	var orderDetail []model.OrderDetail
	err := o.db.Model(&orderDetail).Find(&orderDetail, "order_id = ?", oid).Error
	return orderDetail, err
}

func (o *OrderDao) QueryOrderList(uid uint64, query request.QueryDTO) (*common.PageResult, error) {
	var orderDetail []model.Order
	var pageResult common.PageResult
	db := o.db.Model(&orderDetail) //.Where("user_id = ?", uid)
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	// 查询总数
	if err := db.Count(&pageResult.Total).Error; err != nil {
		return nil, err
	}
	page, _ := strconv.Atoi(query.Page)
	psize, _ := strconv.Atoi(query.Pagesize)
	if err := db.Scopes(pageResult.Paginate(&page, &psize)).
		Order("order_time desc").
		Scan(&orderDetail).Error; err != nil {
		return nil, err
	}
	// 整合数据下发
	pageResult.Records = orderDetail
	return &pageResult, nil
}

func (o *OrderDao) ConditionSearch(uid uint64, query request.OrderSearchDTO) (*common.PageResult, error) {
	var orderDetail []model.Order
	var pageResult common.PageResult
	db := o.db.Model(&orderDetail)
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.BeginTime != "" {
		db = db.Where("order_time >= ?", query.BeginTime)
	}
	if query.EndTime != "" {
		db = db.Where("order_time <= ?", query.EndTime)
	}
	if query.Phone != "" {
		db = db.Where("phone = ?", query.Phone)
	}
	if query.Number != "" {
		db = db.Where("number = ?", query.Number)
	}
	// 查询总数
	if err := db.Count(&pageResult.Total).Error; err != nil {
		return nil, err
	}
	page, _ := strconv.Atoi(query.Page)
	psize, _ := strconv.Atoi(query.Pagesize)
	if err := db.Scopes(pageResult.Paginate(&page, &psize)).
		Order("order_time desc").
		Scan(&orderDetail).Error; err != nil {
		return nil, err
	}
	// 整合数据下发
	pageResult.Records = orderDetail
	return &pageResult, nil
}

func (o *OrderDao) QueryTimeoutOrders(state int8, t time.Time) ([]model.Order, error) {
	var orders []model.Order
	err := o.db.Model(&model.Order{}).
		Where("status = ? and order_time <= ?", state, t).
		Find(&orders).Error
	return orders, err
}

func (o *OrderDao) UpdateTimeoutOrder(state, nstate int8, t time.Time) error {
	var orders model.Order
	err := o.db.Model(&orders).
		Where("status = ? and order_time <= ?", state, t).
		Find(&orders).Update("status", nstate).Error
	return err
}

func (o *OrderDao) UpdateOrderStatus(oids []uint64, os request.OrderStatus) error {
	var orders model.Order
	return o.db.Model(&orders).Where("id in ?", oids).Updates(os).Error
}

func (o *OrderDao) GetOrderById(uid uint64, oid string) (model.Order, error) {
	var order model.Order
	err := o.db.Model(&order).First(&order, "id = ?", uid, oid).Error
	return order, err
}

func (o *OrderDao) GetStatusNumber(uid uint64) (*response.OrderStatusNumber, error) {
	var osn response.OrderStatusNumber
	type resultType struct {
		StatusAlias string
		Number      int
	}
	var results []resultType
	err := o.db.Raw(`
		SELECT 
			CASE 
				WHEN status = 3 THEN 'confirmed'
				WHEN status = 4 THEN 'deliveryInProgress'
				WHEN status = 2 THEN 'toBeConfirmed'
				ELSE 'Other Orders'
			END AS status_alias,
			COUNT(*) AS number
		FROM orders
		GROUP BY status_alias
	`).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		switch result.StatusAlias {
		case "confirmed":
			osn.Confirmed = result.Number
		case "deliveryInProgress":
			osn.DeliveryInProgress = result.Number
		case "toBeConfirmed":
			osn.ToBeConfirmed = result.Number
		}
	}

	return &osn, err
}
