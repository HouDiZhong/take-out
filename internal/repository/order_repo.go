package repository

import (
	"take-out/common"
	"take-out/internal/api/request"
	"take-out/internal/api/response"
	"take-out/internal/model"
	"time"
)

type OrderRepo interface {
	// 创建订单
	CreateOrder(order *model.Order) error
	// 批量创建订单详情
	BatcheCreateOrder(orderDetail []model.OrderDetail) error
	// 通过订单id查询订单详情
	QueryOrderDetailById(oid string) ([]model.OrderDetail, error)
	// 查询超时订单
	QueryTimeoutOrders(state int8, t time.Time) ([]model.Order, error)
	// 查询订单列表
	QueryOrderList(uid uint64, query request.QueryDTO) (*common.PageResult, error)
	// 条件查询
	ConditionSearch(uid uint64, query request.OrderSearchDTO) (*common.PageResult, error)
	// 更新超时订单状态(查询，更新二合一接口)
	UpdateTimeoutOrder(state, nstate int8, t time.Time) error
	// 更新订单状态
	// UpdateOrderStatus(state int8, oids []uint64) error
	// 更新订单状态
	UpdateOrderStatus(oids []uint64, os request.OrderStatus) error
	// 通过订单id查询订单信息
	GetOrderById(uid uint64, oid string) (model.Order, error)
	// 获取订单状态数量
	GetStatusNumber(uid uint64) (*response.OrderStatusNumber, error)
}
