package response

type Order struct {
	ID          uint64  `json:"id"`
	OrderAmount float64 `json:"orderAmount"`
	OrderNumber string  `json:"orderNumber"`
	OrderTime   string  `json:"orderTime"`
}
