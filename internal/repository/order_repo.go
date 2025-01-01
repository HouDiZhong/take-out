package repository

import (
	"take-out/internal/model"
)

type OrderRepo interface {
	CreateOrder(order *model.Order) error
	BatcheCreateOrder(orderDetail []model.OrderDetail) error
	// UpdateOrder(Order model.Order) error
}
