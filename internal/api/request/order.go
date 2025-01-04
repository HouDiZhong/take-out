package request

import (
	"time"
)

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

type QueryDTO struct {
	Page     string `json:"page"`
	Pagesize string `json:"pagesize"`
	Status   string `json:"status"`
}

type ConfirmDTO struct {
	Id int64 `json:"orderId"`
}

type CancelDTO struct {
	Id           int64  `json:"orderId"`
	CancelReason string `json:"cancelReason"`
}

type RejectionDTO struct {
	Id              int64  `json:"orderId"`
	RejectionReason string `json:"RejectionReason"`
}

type OrderStatus struct {
	Type            int8      `json:"type"`
	CancelTime      time.Time `json:"cancelTime"`
	CancelReason    string    `json:"cancelReason"`    // 取消原因
	RejectionReason string    `json:"rejectionReason"` // 拒绝原因
}

type OrderSearchDTO struct {
	BeginTime string `json:"beginTime"`
	EndTime   string `json:"endTime"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	Number    string `json:"number"`
	Page      string `json:"page"`
	Pagesize  string `json:"pagesize"`
}
