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
	Id int64 `json:"id"`
}

type CancelDTO struct {
	Id           int64  `json:"id"`
	CancelReason string `json:"cancelReason"`
}

type RejectionDTO struct {
	Id              int64  `json:"id"`
	RejectionReason string `json:"RejectionReason"`
}

type OrderStatus struct {
	Status          int8      `json:"status"`          // 订单状态码
	CancelTime      time.Time `json:"cancelTime"`      // 取消时间
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
