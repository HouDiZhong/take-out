package repository

import (
	"take-out/internal/model"
	"time"
)

type OrderRepo interface {
	CreateOrder(order *model.Order) error
	BatcheCreateOrder(orderDetail []model.OrderDetail) error
	QueryTimeoutOrders(state int8, t time.Time) ([]model.Order, error)
	UpdateTimeoutOrder(state, nstate int8, t time.Time) error
	UpdateOrderStatus(state int8, oids []uint64) error
}
