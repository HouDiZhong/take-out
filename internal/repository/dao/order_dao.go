package dao

import (
	"take-out/global"
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

func (o *OrderDao) UpdateOrderStatus(state int8, oids []uint64) error {
	var orders model.Order
	return o.db.Model(&orders).Where("id in ?", oids).Update("status", state).Error
}
