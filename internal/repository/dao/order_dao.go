package dao

import (
	"take-out/global"
	"take-out/internal/model"
	"take-out/internal/repository"

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
