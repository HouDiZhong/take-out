package response

import (
	"take-out/internal/model"
)

type Order struct {
	ID          uint64  `json:"id"`
	OrderAmount float64 `json:"orderAmount"`
	OrderNumber string  `json:"orderNumber"`
	OrderTime   string  `json:"orderTime"`
}

type Websocket struct {
	Type    int8   `json:"type"` // 1 来单 2 催单
	OrderId uint64 `json:"orderId"`
	Content string `json:"content"`
}

type OrderDetail struct {
	model.Order
	OrderDetailList []model.OrderDetail
}

type OrderStatusNumber struct {
	Confirmed          int `json:"confirmed"`          // 待派送数量
	DeliveryInProgress int `json:"deliveryInProgress"` // 派送中数量
	ToBeConfirmed      int `json:"toBeConfirmed"`      // 待接单数量
}
