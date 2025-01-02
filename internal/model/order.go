package model

import "time"

type Order struct {
	ID                    uint64    `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Number                string    `json:"number" gorm:"size:50"`
	Status                int8      `json:"status"`
	UserID                uint64    `json:"user_id"`
	AddressBookID         uint64    `json:"address_book_id"`
	OrderTime             time.Time `json:"order_time"`
	CheckoutTime          time.Time `json:"checkout_time"`
	PayMethod             int8      `json:"pay_method"`
	PayStatus             int8      `json:"pay_status"`
	Amount                float64   `json:"amount"`
	Remark                string    `json:"remark" gorm:"size:100"`
	Phone                 string    `json:"phone" gorm:"size:11"`
	Address               string    `json:"address" gorm:"size:255"`
	UserName              string    `json:"user_name" gorm:"size:32"`
	Consignee             string    `json:"consignee" gorm:"size:32"`
	CancelReason          string    `json:"cancel_reason" gorm:"size:255"`
	RejectionReason       string    `json:"rejection_reason" gorm:"size:255"`
	CancelTime            time.Time `json:"cancel_time"`
	EstimatedDeliveryTime time.Time `json:"estimated_delivery_time"`
	DeliveryStatus        int8      `json:"delivery_status"`
	DeliveryTime          time.Time `json:"delivery_time"`
	PackAmount            int       `json:"pack_amount"`
	TablewareNumber       int8      `json:"tableware_number"`
	TablewareStatus       int8      `json:"tableware_status"`
}

type OrderDetail struct {
	ID         uint64  `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Name       string  `json:"name" gorm:"size:32"`
	Image      string  `json:"image" gorm:"size:255"`
	OrderID    uint64  `json:"order_id"`
	DishID     uint64  `json:"dish_id"`
	SetmealID  uint64  `json:"setmeal_id"`
	DishFlavor string  `json:"dish_flavor" gorm:"size:50"`
	Number     int     `json:"number"`
	Amount     float64 `json:"amount"`
}

func (o *Order) TableName() string {
	return "orders"
}
