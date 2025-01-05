package response

type BusinessDataVO struct {
	NewUsers int64 `json:"newUsers"` // 新增用户数
	BusinessOrderVO
}

type BusinessOrderVO struct {
	OrderStatusNumber float64 `json:"orderCompletionRate"` // 订单完成率
	Turnover          float64 `json:"turnover"`            // 营业额
	UnitPrice         float64 `json:"unitPrice"`           // 平均客单价
	ValidOrderCount   int64   `json:"validOrderCount"`     // 有效订单数
}

type SetmealAndDishVO struct {
	Discontinued int `json:"discontinued"`
	Sold         int `json:"sold"`
}

type OrderNumberVO struct {
	AllOrders       int `json:"allOrders"`       // 全部订单
	CancelledOrders int `json:"cancelledOrders"` // 已取消数量
	CompletedOrders int `json:"completedOrders"` // 已完成数量
	DeliveredOrders int `json:"deliveredOrders"` // 待派送数量
	WaitingOrders   int `json:"waitingOrders"`   // 待接单数量
}
