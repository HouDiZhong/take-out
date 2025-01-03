package response

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
