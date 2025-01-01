package request

/*
	"addressBookId": 0,

"amount": 0,
"deliveryStatus": 0,
"estimatedDeliveryTime": "yyyy-MM-dd HH:mm:ss",
"packAmount": 0,
"payMethod": 0,
"remark": "string",
"tablewareNumber": 0,
"tablewareStatus": 0
*/
type OrderDTO struct {
	AddressBookId         uint64  `json:"addressBookId"`
	Amount                float64 `json:"amount"`
	DeliveryStatus        int8    `json:"deliveryStatus"`
	EstimatedDeliveryTime string  `json:"estimatedDeliveryTime"`
	PackAmount            int     `json:"packAmount"`
	PayMethod             int8    `json:"payMethod"`
	Remark                string  `json:"remark"`
	TablewareNumber       int8    `json:"tablewareNumber"`
	TablewareStatus       int8    `json:"tablewareStatus"`
}
